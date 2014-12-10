package main

import (
	"encoding/csv"
	"errors"
	"os"
	"runtime"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/cryptix/go/logging"
	"golang.org/x/crypto/ssh/terminal"
)

var log = logging.Logger("gombrella")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app := cli.NewApp()
	app.Name = "gombrella"
	app.Version = "0.0.2"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "file,f", Value: "bookmarks.csv", Usage: "csv filename to write bookmarks to"},
	}

	app.Run(os.Args)
}

func run(ctx *cli.Context) {
	u := ctx.Args().First()
	if len(u) == 0 {
		logging.CheckFatal(errors.New("no user supplied"))
	}

	fname := ctx.String("file")
	if len(fname) == 0 {
		logging.CheckFatal(errors.New("filename can't be empty"))
	}

	// prepare the output file
	file, err := os.Create(fname)
	logging.CheckFatal(err)

	log.Notice("Enter Password:")
	pw, err := terminal.ReadPassword(0)
	logging.CheckFatal(err)

	colls, client, err := loginAndGetCollections(u, string(pw))
	logging.CheckFatal(err)

	log.Notice("Logged in.")

	// each collection is requested asynchronously with a  worker
	workerCnt := 3 * runtime.NumCPU()
	workerChans := make([]<-chan *Bookmark, workerCnt)
	for i := 0; i < workerCnt; i++ {
		workerChans[i] = bookmarkWorker(client, colls)
	}

	bookmarks := mergeWorkers(workerChans...)

	bookmarkCsv := csv.NewWriter(file)

	for a := range bookmarks {
		rec := []string{
			a.Title,
			a.Link,
		}
		logging.CheckFatal(bookmarkCsv.Write(rec))

	}
	bookmarkCsv.Flush()
	logging.CheckFatal(file.Close())
	log.Notice("Done. Written to ", fname)
}

func mergeWorkers(cs ...<-chan *Bookmark) <-chan *Bookmark {
	var wg sync.WaitGroup
	out := make(chan *Bookmark)

	output := func(c <-chan *Bookmark) {
		for a := range c {
			out <- a
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

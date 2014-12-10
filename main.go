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
	app.Action = run
	app.Flags = []cli.Flag{
	}

	app.Run(os.Args)
}

func run(ctx *cli.Context) {
	u := ctx.Args().First()
	if len(u) == 0 {
		logging.CheckFatal(errors.New("no user supplied"))
	}

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

	// prepare the output file
	file, err := os.Create("bookmarks.csv")
	logging.CheckFatal(err)

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
	log.Notice("Done.")
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

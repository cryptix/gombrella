package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func bookmarkWorker(client *http.Client, colls <-chan *CollectionMeta) <-chan *Bookmark {

	out := make(chan *Bookmark)
	go func() {
		for coll := range colls {
			log.Debugf("Requesting bookmarks for %s", coll)

			url := fmt.Sprintf("https://raindrop.io/api/raindrops/%.0f", coll.ID)
			resp, err := client.Get(url)
			if err != nil {
				log.Criticalf("Get of bookmarks from %s failed.\n%s", coll, err)
				break
			}
			defer resp.Body.Close()

			var collResp CollectionResp
			err = json.NewDecoder(resp.Body).Decode(&collResp)
			if err != nil {
				log.Criticalf("Decoding of data for %s failed.\n%s", coll, err)
				break
			}

			if !collResp.Result {
				log.Criticalf("%s - result in collResp is false", coll)
				break
			}

			for _, bookmark := range collResp.Items {
				out <- bookmark
			}
		}
		close(out)
	}()

	return out
}

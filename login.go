package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
)

func loginAndGetCollections(user, passw string) (<-chan *CollectionMeta, *http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, nil, err
	}

	client := &http.Client{Jar: jar}
	var loginCreds = struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{user, passw}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(loginCreds)
	if err != nil {
		return nil, nil, err
	}
	liResp, err := client.Post("https://raindrop.io/api/auth/login", "application/json", &buf)
	if err != nil {
		return nil, nil, err
	}
	defer liResp.Body.Close()

	if liResp.StatusCode != http.StatusOK {
		return nil, nil, errors.New("login failed")
	}

	var result struct {
		Result bool `json:"result"`
	}
	if err = json.NewDecoder(liResp.Body).Decode(&result); err != nil {
		return nil, nil, err
	}

	if !result.Result {
		return nil, nil, errors.New("login failed")
	}

	log.Debug("Logged in")

	req, err := http.NewRequest("GET", "https://raindrop.io/api/collections", nil)
	if err != nil {
		return nil, nil, err
	}

	collResp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer collResp.Body.Close()

	var colls CollectionsMetaResp
	err = json.NewDecoder(collResp.Body).Decode(&colls)
	if err != nil {
		return nil, nil, err
	}

	if !colls.Result {
		return nil, nil, errors.New("not a result?!")
	}
	log.Debugf("collections requested. #%d", len(colls.Items))

	collchan := make(chan *CollectionMeta)
	go func() {
		for _, c := range colls.Items {
			collchan <- c
		}
		close(collchan)
	}()

	return collchan, client, nil
}

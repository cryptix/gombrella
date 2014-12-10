package main

import "fmt"

type CollectionMeta struct {
	ID         float64 `json:"_id"`
	Background string  `json:"background"`
	Count      float64 `json:"count"`
	Created    string  `json:"created"`
	Excerpt    string  `json:"excerpt"`
	Lang       string  `json:"lang"`
	LastUpdate string  `json:"lastUpdate"`
	Public     bool    `json:"public"`
	ShortLink  string  `json:"shortLink"`
	Title      string  `json:"title"`
	User       struct {
		_Db  string  `json:"$db"`
		_Id  float64 `json:"$id"`
		_Ref string  `json:"$ref"`
	} `json:"user"`
	View string `json:"view"`
}

func (c CollectionMeta) String() string {
	return fmt.Sprintf("CollectionMeta(%.0f) %s", c.ID, c.Title)
}

type CollectionsMetaResp struct {
	Items  []*CollectionMeta `json:"items"`
	Result bool              `json:"result"`
}

type Bookmark struct {
	ID         float64 `json:"_id"`
	Cover      string  `json:"cover"`
	CoverId    float64 `json:"coverId"`
	Domain     string  `json:"domain"`
	Excerpt    string  `json:"excerpt"`
	LastUpdate string  `json:"lastUpdate"`
	Link       string  `json:"link"`
	Media      []struct {
		Link string `json:"link"`
		Type string `json:"type"`
	} `json:"media"`
	Removed bool          `json:"removed"`
	Sort    float64       `json:"sort"`
	Tags    []interface{} `json:"tags"`
	Title   string        `json:"title"`
	Type    string        `json:"type"`
	User    struct {
		_Db  string  `json:"$db"`
		_Id  float64 `json:"$id"`
		_Ref string  `json:"$ref"`
	} `json:"user"`
}

type CollectionResp struct {
	Author bool        `json:"author"`
	Count  float64     `json:"count"`
	Items  []*Bookmark `json:"items"`
	Result bool        `json:"result"`
}

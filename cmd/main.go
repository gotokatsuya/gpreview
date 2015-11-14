package main

import (
	gp "github.com/gotokatsuya/gpreview"
	"log"
)

func init() {
	gp.Load()
}

func main() {
	log.Println(gp.InsertReviews("gpreviews.db", "review_2015-11-14.csv"))
}

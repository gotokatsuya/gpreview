package main

import (
	"flag"
	"log"

	gp "github.com/gotokatsuya/gpreview-go/gpreview"
)

func init() {
	gp.Load()
}

func main() {
	flag.Usage = func() {
		flag.PrintDefaults()
	}
	db := flag.String("db", "", "A path of writable database")
	file := flag.String("file", "", "A path of review csv file")
	from := flag.String("from", "", "A language you want to translate")
	to := flag.String("to", "", "A language you can understand")
	flag.Parse()
	if *db == "" {
		log.Println("Specify a path of db, please")
		return
	}
	if *file == "" {
		log.Println("Specify a path of file, please")
		return
	}
	if err := gp.NotifyTranslatedReviews(*db, *file, *from, *to); err != nil {
		log.Println(err)
	}
	log.Println("bye")
}

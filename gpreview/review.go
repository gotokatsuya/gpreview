package gpreview

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func OpenReviewDatabase(dbName string) *sql.DB {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}

	createReviewTableStmt := `
	CREATE TABLE IF NOT EXISTS review (package_name text, app_version text, review_submit_datetime text, star_rating text, review_title text, review_text text, UNIQUE(package_name, review_submit_datetime) ON CONFLICT ABORT);
	`
	_, err = db.Exec(createReviewTableStmt)
	if err != nil {
		panic(err)
	}
	return db
}

// Review ...
type Review struct {
	PackageName string
	AppVersion  string
	// ReviewerLanguage                 string
	// ReviewerHardwareModel            string
	ReviewSubmitDateAndTime string
	// ReviewSubmitMillisSinceEpoch     string
	// ReviewLastUpdateDateAndTime      string
	// ReviewLastUpdateMillisSinceEpoch string
	StarRating  string
	ReviewTitle string
	ReviewText  string
	// DeveloperReplyDateAndTime        string
	// DeveloperReplyMillisSinceEpoch   string
	// DeveloperReplyText               string
	// ReviewLink                       string
}

func (r *Review) Insert(db *sql.DB) error {
	insertReviewTableStmt := `
	INSERT into review values('%s', '%s', '%s', '%s', '%s', '%s');
	`
	_, err := db.Exec(fmt.Sprintf(insertReviewTableStmt, r.PackageName, r.AppVersion, r.ReviewSubmitDateAndTime, r.StarRating, r.ReviewTitle, r.ReviewText))
	if err != nil {
		return err
	}
	return nil
}

func NotifyTranslatedReviews(dbName, fileName, from, to string) error {
	db := OpenReviewDatabase(dbName)

	defer db.Close()

	atmc := new(MsAccessTokenMessageCache)

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if len(record) <= 0 {
			continue
		}

		log.Println(fmt.Sprintf("%#v", record))

		review := Review{
			PackageName: record[0],
			AppVersion:  record[1],
			// ReviewerLanguage:                 record[2],
			// ReviewerHardwareModel:            record[3],
			ReviewSubmitDateAndTime: record[4],
			// ReviewSubmitMillisSinceEpoch:     record[5],
			// ReviewLastUpdateDateAndTime:      record[6],
			// ReviewLastUpdateMillisSinceEpoch: record[7],
			StarRating:  record[8],
			ReviewTitle: record[9],
			ReviewText:  record[10],
			// DeveloperReplyDateAndTime:        record[11],
			// DeveloperReplyMillisSinceEpoch:   record[12],
			// DeveloperReplyText:               record[13],
			// ReviewLink:                       record[14],
		}
		if err := review.Insert(db); err != nil {
			log.Println(err)
			continue
		}

		title := review.ReviewTitle
		if len(title) > 0 && len(from) > 0 && len(to) > 0 {
			title, err := Translate(title, from, to, atmc)
			if err != nil {
				log.Println(err)
			}
			log.Println(title)
		}

		text := review.ReviewText
		if len(text) > 0 && len(from) > 0 && len(to) > 0 {
			text, err := Translate(text, from, to, atmc)
			if err != nil {
				log.Println(err)
			}
			log.Println(text)
		}

		if len(GPReview.SlackURL) > 0 {
			if err := PostSlack(SlackData{
				text + "\n" + review.StarRating + "\n" + review.AppVersion + " " + review.ReviewSubmitDateAndTime,
				title,
				":santa:"},
				GPReview.SlackURL); err != nil {
				log.Println(err)
				continue
			}
		}
	}
	return nil
}

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

func OpenDatabase(dbName string) *sql.DB {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}

	createReviewTableStmt := `
	CREATE TABLE IF NOT EXISTS review (package_name text, review_submit_datetime text, UNIQUE(package_name, review_submit_datetime) ON CONFLICT ABORT);
	`
	_, err = db.Exec(createReviewTableStmt)
	if err != nil {
		panic(err)
	}
	return db
}

// Review ...
type Review struct {
	PackageName                      string // パッケージ名
	AppVersion                       string // バージョンナンバー
	ReviewerLanguage                 string // レビューの言語
	ReviewerHardwareModel            string // レビュー者の利用デバイス
	ReviewSubmitDateAndTime          string // 投稿日時(YYYY-MM-DDThh:mm:ssTZ)
	ReviewSubmitMillisSinceEpoch     string // 投稿日時(エポックタイムms)
	ReviewLastUpdateDateAndTime      string // 更新日時(YYYY-MM-DDThh:mm:ssTZ)
	ReviewLastUpdateMillisSinceEpoch string // 更新日時(エポックタイムms)
	StarRating                       string // 星の数
	ReviewTitle                      string // レビュータイトル
	ReviewText                       string // レビュー本文
	DeveloperReplyDateAndTime        string // 返信日時(YYYY-MM-DDThh:mm:ssTZ)
	DeveloperReplyMillisSinceEpoch   string // 返信日時(エポックタイムms)
	DeveloperReplyText               string // 返信内容
	ReviewLink                       string // レビューURL
}

func (r *Review) Insert(db *sql.DB) error {
	insertReviewTableStmt := `
	INSERT into review values('%s', '%s');
	`
	_, err := db.Exec(fmt.Sprintf(insertReviewTableStmt, r.PackageName, r.ReviewSubmitDateAndTime))
	if err != nil {
		return err
	}
	return nil
}

func InsertReviews(dbName, fileName string, tranlate bool) error {
	db := OpenDatabase(dbName)

	defer db.Close()

	atmc := new(AccessTokenMessageCache)

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	reader := csv.NewReader(file)

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

		log.Printf("%#v", record)

		review := Review{
			PackageName:                      record[0],
			AppVersion:                       record[1],
			ReviewerLanguage:                 record[2],
			ReviewerHardwareModel:            record[3],
			ReviewSubmitDateAndTime:          record[4],
			ReviewSubmitMillisSinceEpoch:     record[5],
			ReviewLastUpdateDateAndTime:      record[6],
			ReviewLastUpdateMillisSinceEpoch: record[7],
			StarRating:                       record[8],
			ReviewTitle:                      record[9],
			ReviewText:                       record[10],
			DeveloperReplyDateAndTime:        record[11],
			DeveloperReplyMillisSinceEpoch:   record[12],
			DeveloperReplyText:               record[13],
			ReviewLink:                       record[14],
		}
		if err := review.Insert(db); err != nil {
			log.Println(err)
			continue
		}
		postWord := review.ReviewText
		if len(postWord) <= 0 {
			continue
		}
		if tranlate {
			postWord = Translate(postWord, review.ReviewerLanguage, "ja", atmc)
			if err != nil {
				log.Println(err)
				continue
			}
		}
		if err := PostSlack(SlackData{postWord, review.PackageName, ":santa:"}, GPReview.SlackURL); err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}

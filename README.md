# gpreview
Get google play app's reviews and notify slack.

#### Feature
- Translatable


## Install
```
go get github.com/gotokatsuya/gpreview-go/cmd/gpreview
```


## Execution
```
 gpreview -db=gpreviews.db -file=reviews__201406.csv
```

#### Translation Option
```
 gpreview -db=gpreviews.db -file=reviews__201406.csv -from=en -to=ja
```


### Database

Create a database.
```
sqlite3 gpreviews.db
```

### Download Review files

Use [gsutil](https://cloud.google.com/storage/docs/gsutil), please.


## Config

### Microsoft

[Translator API](http://www.microsoft.com/en-us/translator/translatorapi.aspx)

```
ms_tranlator_client_id=""
ms_tranlator_client_secret=""
```

### Slack

[Webhooks](https://api.slack.com/incoming-webhooks)

```
slack_url=""
```

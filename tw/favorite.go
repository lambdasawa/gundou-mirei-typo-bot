package tw

import (
	"net/url"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/lambdasawa/gundou-mirei-typo-bot/validate"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	dryRun bool
)

type (
	TwitterConfig struct {
		ConsumerKey    string
		ConsumerSecret string
		AccessToken    string
		AccessSecret   string
	}
)

func Favorite(conf TwitterConfig) error {
	log.WithField("conf", conf).Info("config")

	client := twitter.NewClient(
		oauth1.
			NewConfig(conf.ConsumerKey, conf.ConsumerSecret).
			Client(
				oauth1.NoContext,
				oauth1.NewToken(conf.AccessToken, conf.AccessSecret),
			),
	)

	invalidTexts := validate.GetInvalidTexts()
	for _, text := range invalidTexts {
		searchParams := twitter.SearchTweetParams{
			Query:  url.QueryEscape(text),
			Lang:   "ja",
			Locale: "ja",
			Count:  100,
		}
		search, _, err := client.Search.Tweets(&searchParams)
		if err != nil {
			return errors.Wrap(err, "failed to call search API")
		}
		log.WithFields(log.Fields{"request": searchParams, "response": search}).Debugf("search tweets")

		for _, status := range search.Statuses {
			if err := validateTweet(client, status); err != nil {
				if isBlockedError(err) {
					continue
				}
				return errors.Wrapf(err, "failed to favorite tweet. %+v", status)
			}
		}
	}

	return nil
}

func validateTweet(client *twitter.Client, tweet twitter.Tweet) error {
	if validate.IsValidText(tweet.Text) {
		log.WithField("tweet", tweet).Info("valid tweet")
		return nil
	}
	log.WithField("tweet", tweet).Info("invalid tweet")

	statusDetail, resp, err := client.Statuses.Show(tweet.ID, &twitter.StatusShowParams{
		ID: tweet.ID,
	})
	if err != nil {
		return errors.Wrapf(err, "failed show status. https://twitter.com/_/status/%v", tweet.ID)
	}
	log.Debugf("status detail. %+v %+v", statusDetail, resp)
	if statusDetail.Favorited {
		log.WithField("tweet", tweet).Info("already favorited")
		return nil
	}

	if dryRun {
		log.WithField("tweet", tweet).Info("skip POST /favorites/create")
		return nil
	}

	if _, _, err := client.Favorites.Create(&twitter.FavoriteCreateParams{
		ID: tweet.ID,
	}); err != nil {
		return errors.Wrap(err, "failed create favorites")
	}
	log.WithField("tweet", tweet).Infof("favorite")

	return nil
}

func isBlockedError(err error) bool {
	return strings.Contains(err.Error(), "You have been blocked from the author of this tweet.")
}

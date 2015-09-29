package api

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/rendon/anaconda"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"identities/app/config"
	"identities/app/models"
)

var (
	ErrUnknownNet = errors.New("Unknown network.")
)

var getFrom = map[string]func(string, string) (*models.Identity, error){
	"twitter": getFromTwitter,
}

var addTo = map[string]func(*models.Identity, *mgo.Collection) error{
	"twitter": addToTwitter,
}

var ta *anaconda.TwitterApi // Twitter API

func init() {
	var ck = os.Getenv("TWITTER_CONSUMER_KEY")
	var cs = os.Getenv("TWITTER_CONSUMER_SECRET")
	var at = os.Getenv("TWITTER_ACCESS_TOKEN")
	var ats = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	anaconda.SetConsumerKey(ck)
	anaconda.SetConsumerSecret(cs)
	ta = anaconda.NewTwitterApi(at, ats)
}

func WipeIdentitiesDatabase() error {
	var session, err = mgo.Dial(config.MongoDSN)
	if err != nil {
		return err
	}
	defer session.Close()
	return session.DB("identities").DropDatabase()
}

func getFromTwitter(id, username string) (*models.Identity, error) {
	if id == "" && username == "" {
		return nil, fmt.Errorf("Id and username are both empty.")
	}

	var user anaconda.User
	var err error
	if id != "" {
		nid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}
		user, err = ta.GetUsersShowById(nid, nil)
	} else {
		user, err = ta.GetUsersShow(username, nil)
	}

	var i = models.Identity{
		Network:         "twitter",
		Id:              user.IdStr,
		Username:        strings.ToLower(user.ScreenName),
		ProfileImageURL: user.ProfileImageURL,
	}
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func ensureIdentitySchema(c *mgo.Collection) error {
	var err error
	// Ensure <network, id> to be unique
	// This also will speed up queries on this two fields.
	err = c.EnsureIndex(mgo.Index{
		Key:      []string{"network", "id"},
		Unique:   true,
		DropDups: true,
		Sparse:   true,
	})

	if err != nil {
		return err
	}

	// Ensure <network, username> to be unique.
	// This also will speed up queries on this two fields.
	err = c.EnsureIndex(mgo.Index{
		Key:      []string{"network", "username"},
		Unique:   true,
		DropDups: true,
		Sparse:   true,
	})
	return nil
}

func addToTwitter(item *models.Identity, c *mgo.Collection) error {
	if item == nil {
		return errors.New("item is nil.")
	}
	item.Status = "cached"
	var err error
	if err = ensureIdentitySchema(c); err != nil {
		return err
	}
	return c.Insert(item)
}

func GetIdentity(network, id, username string) (*models.Identity, error) {
	var ok bool
	if _, ok = config.Networks[network]; !ok {
		return nil, ErrUnknownNet
	}
	var session, err = mgo.Dial(config.MongoDSN)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	var identities = session.DB("identities").C(network)
	var one = new(models.Identity)
	var two models.Identity
	var selector bson.M
	if id != "" {
		selector = bson.M{
			"network": network,
			"id":      id,
		}
	} else {
		selector = bson.M{
			"network":  network,
			"username": username,
		}
	}

	err = identities.Find(selector).One(&two)
	one = &two

	if err != nil && err.Error() != "not found" {
		return nil, err
	}

	if err != nil && err.Error() == "not found" {
		log.Printf("Not found!")
		one, err = getFrom[network](id, username)
		if err != nil {
			return nil, err
		}
		if err = addTo[network](one, identities); err != nil {
			log.Printf("==>Error: %#v", one)
			return nil, err
		}
		one.Status = "new"
	}
	return one, nil
}

func WipeDatabase() error {
	var session, err = mgo.Dial(config.MongoDSN)
	if err != nil {
		return err
	}
	defer session.Close()
	return session.DB("identities").DropDatabase()
}

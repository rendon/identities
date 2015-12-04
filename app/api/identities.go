package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rendon/anaconda"
	"github.com/rendon/kb"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"identities/app/models"
)

var (
	ErrUnknownNet = errors.New("Unknown network.")
	MongoDSN      = "mongodb-server"
)

var ExpirationTime = map[string]time.Duration{
	"twitter":   30 * 24 * time.Hour,
	"instagram": 30 * 24 * time.Hour,
}

var Networks = map[string]bool{
	"twitter":   true,
	"instagram": true,
}

var getFrom = map[string]func(string, string) (*models.Identity, error){
	"twitter": getFromTwitter,
}

var addTo = map[string]func(*models.Identity, *mgo.Collection) error{
	"twitter": addToTwitter,
}

var ta *anaconda.TwitterApi // Twitter API
var limit = 0

func init() {
	key_file := os.Getenv("TWITTER_KEYS_FILE")
	buf, err := ioutil.ReadFile(key_file)
	if err != nil {
		log.Fatalf("Failed to read keys file: %s", err)
	}
	for _, line := range strings.Split(string(buf), "\n") {
		if line != "" {
			kb.AddKey(kb.Key{Value: line})
		}
	}
	rotateKeys()
}

func rotateKeys() {
	k := kb.NextKey()
	tokens := strings.Split(k.Value.(string), " ")
	if len(tokens) != 4 {
		log.Fatalf("Failed to set keys: %q", k.Value.(string))
	}
	ck := tokens[0]
	cs := tokens[1]
	at := tokens[2]
	ats := tokens[3]
	anaconda.SetConsumerKey(ck)
	anaconda.SetConsumerSecret(cs)
	ta = anaconda.NewTwitterApi(at, ats)
	ta.ReturnRateLimitError(true)
}

func WipeIdentitiesDatabase() error {
	var session, err = mgo.Dial(MongoDSN)
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
	if err != nil || user.IdStr == "" {
		limit++
		if limit == 10 {
			limit = 0
			return nil, errors.New("Not found")
		}
		rotateKeys()
		return getFromTwitter(id, username)
	}

	var i = models.Identity{
		Network:         "twitter",
		Id:              user.IdStr,
		Username:        strings.ToLower(user.ScreenName),
		ProfileImageURL: user.ProfileImageURL,
		ProfileURL:      "https://twitter.com/" + user.ScreenName,
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
	network = strings.ToLower(network)
	username = strings.ToLower(username)
	id = strings.ToLower(id) // ID is not necessarily an integer

	if strings.HasPrefix(username, "@") {
		username = username[1:]
	}

	var ok bool
	if _, ok = Networks[network]; !ok {
		return nil, ErrUnknownNet
	}
	var session, err = mgo.Dial(MongoDSN)
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
		one, err = getFrom[network](id, username)
		if err != nil {
			return nil, err
		}
		if err = addTo[network](one, identities); err != nil {
			log.Printf("Failed to obtain profile: %s", err)
			return nil, errors.New("Failed to obtain profile")
		}
		one.Status = "new"
	}
	return one, nil
}

func WipeDatabase() error {
	var session, err = mgo.Dial(MongoDSN)
	if err != nil {
		return err
	}
	defer session.Close()
	return session.DB("identities").DropDatabase()
}

func Identify(user string) (*models.Identity, error) {
	// TODO: Generalize for all social networks
	if _, err := strconv.ParseInt(user, 10, 64); err == nil {
		return GetIdentity("twitter", user, "")
	} else {
		return GetIdentity("twitter", "", user)
	}
}

package tests

import (
	"encoding/json"
	"log"
	"time"

	"identities/app/models"
)

func (t *AppTest) TestGetUncachedUserId() {
	t.Get("/ids/twitter/twitterdev")
	log.Printf("Response body: %s\n", t.ResponseBody)
	t.AssertOk()
	var res models.IdentityResponse
	t.Assert(json.Unmarshal(t.ResponseBody, &res) == nil)
	t.Assert(res.Data.Id == "2244994945")
	t.Assert(res.Data.Status == "new")
}

func (t *AppTest) TestGetCachedUserId() {
	t.Get("/ids/twitter/twitterdev")
	log.Printf("Response body: %s\n", t.ResponseBody)
	t.AssertOk()

	time.Sleep(0 * time.Second)
	t.Get("/ids/twitter/twitterdev")
	log.Printf("Response body 2: %s\n", t.ResponseBody)
	t.AssertOk()
	var res models.IdentityResponse
	t.Assert(json.Unmarshal(t.ResponseBody, &res) == nil)
	t.Assert(res.Data.Id == "2244994945")
	t.Assert(res.Data.Status == "cached")
}

func (t *AppTest) TestGetUncachedUserName() {
	t.Get("/usernames/twitter/2244994945")
	log.Printf("Response body: %s\n", t.ResponseBody)
	t.AssertOk()
	var res models.IdentityResponse
	t.Assert(json.Unmarshal(t.ResponseBody, &res) == nil)
	t.Assert(res.Data.Username == "twitterdev")
	t.Assert(res.Data.Status == "new")
}

func (t *AppTest) TestGetCachedUserName() {
	t.Get("/usernames/twitter/2244994945")
	t.AssertOk()

	t.Get("/usernames/twitter/2244994945")
	log.Printf("Response body 2: %s\n", t.ResponseBody)
	t.AssertOk()
	var res models.IdentityResponse
	t.Assert(json.Unmarshal(t.ResponseBody, &res) == nil)
	t.Assert(res.Data.Username == "twitterdev")
	t.Assert(res.Data.Status == "cached")
}

func (t *AppTest) TestGetAtUsername() {
	t.Get("/ids/twitter/twitterdev")
	t.AssertOk()

	t.Get("/ids/twitter/@twitterdev")
	t.AssertOk()
	var res models.IdentityResponse
	t.Assert(json.Unmarshal(t.ResponseBody, &res) == nil)
	t.Assert(res.Data.Id == "2244994945")
	t.Assert(res.Data.Username == "twitterdev")
	t.Assert(res.Data.Status == "cached")
}

func (t *AppTest) TestGetIdentityWithUsername() {
	t.Get("/identities/twitterdev")
	t.AssertOk()

	var res models.IdentityResponse
	t.Assert(json.Unmarshal(t.ResponseBody, &res) == nil)
	t.Assert(res.Data.Id == "2244994945")
	t.Assert(res.Data.Username == "twitterdev")
}

func (t *AppTest) TestGetIdentityWithID() {
	t.Get("/identities/2244994945")
	t.AssertOk()

	var res models.IdentityResponse
	t.Assert(json.Unmarshal(t.ResponseBody, &res) == nil)
	t.Assert(res.Data.Id == "2244994945")
	t.Assert(res.Data.Username == "twitterdev")
}

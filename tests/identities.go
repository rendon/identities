package tests

import (
	"encoding/json"
	"time"

	"identities/app/models"
)

func (t *AppTest) TestGetUncachedUserId() {
	t.Get("/v1/ids/twitter/twitterdev")
	t.AssertOk()
	var res models.IdentityResponse
	t.Assert(json.Unmarshal(t.ResponseBody, &res) == nil)
	t.Assert(res.Data.ID == "2244994945")
	t.Assert(res.Data.Status == "new")
}

func (t *AppTest) TestGetCachedUserId() {
	t.Get("/v1/ids/twitter/twitterdev")
	t.AssertOk()

	time.Sleep(0 * time.Second)
	t.Get("/v1/ids/twitter/twitterdev")
	t.AssertOk()
	var res models.IdentityResponse
	t.Assert(json.Unmarshal(t.ResponseBody, &res) == nil)
	t.Assert(res.Data.ID == "2244994945")
	t.Assert(res.Data.Status == "cached")
}

func (t *AppTest) TestGetUncachedUserName() {
	t.Get("/v1/usernames/twitter/2244994945")
	t.AssertOk()
	var res models.IdentityResponse
	t.Assert(json.Unmarshal(t.ResponseBody, &res) == nil)
	t.Assert(res.Data.Username == "twitterdev")
	t.Assert(res.Data.Status == "new")
}

func (t *AppTest) TestGetCachedUserName() {
	t.Get("/v1/usernames/twitter/2244994945")
	t.AssertOk()

	t.Get("/v1/usernames/twitter/2244994945")
	t.AssertOk()
	var res models.IdentityResponse
	t.Assert(json.Unmarshal(t.ResponseBody, &res) == nil)
	t.Assert(res.Data.Username == "twitterdev")
	t.Assert(res.Data.Status == "cached")
}

func (t *AppTest) TestGetAtUsername() {
	t.Get("/v1/ids/twitter/twitterdev")
	t.AssertOk()

	t.Get("/v1/ids/twitter/@twitterdev")
	t.AssertOk()
	var res models.IdentityResponse
	t.Assert(json.Unmarshal(t.ResponseBody, &res) == nil)
	t.Assert(res.Data.ID == "2244994945")
	t.Assert(res.Data.Username == "twitterdev")
	t.Assert(res.Data.Status == "cached")
}

func (t *AppTest) TestGetIdentityWithUsername() {
	t.Get("/v1/identities/twitterdev")
	t.AssertOk()

	var res models.IdentityResponse
	t.Assert(json.Unmarshal(t.ResponseBody, &res) == nil)
	t.Assert(res.Data.ID == "2244994945")
	t.Assert(res.Data.Username == "twitterdev")
}

func (t *AppTest) TestGetIdentityWithID() {
	t.Get("/v1/identities/2244994945")
	t.AssertOk()

	var res models.IdentityResponse
	t.Assert(json.Unmarshal(t.ResponseBody, &res) == nil)
	t.Assert(res.Data.ID == "2244994945")
	t.Assert(res.Data.Username == "twitterdev")
}

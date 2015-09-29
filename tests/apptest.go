package tests

import (
	"bytes"
	"encoding/json"
	"github.com/revel/revel/testing"
	"identities/app/api"
	"io"
	"log"
)

type AppTest struct {
	testing.TestSuite
}

const (
	jsonType = "application/json"
)

func serialize(m interface{}) io.Reader {
	var result, err = json.Marshal(m)
	if err != nil {
		// Unlikely
	}
	return bytes.NewReader([]byte(result))
}

func (t *AppTest) Before() {
	if err := api.WipeDatabase(); err != nil {
		log.Fatalf("Couldn't setup database")
	}

	if err := api.WipeIdentitiesDatabase(); err != nil {
		log.Fatalf("Couldn't setup database: %s\n", err)
	}
}

func (t *AppTest) After() {
}

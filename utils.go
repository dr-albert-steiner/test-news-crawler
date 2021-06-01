package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func decodeJSON(data io.Reader, v interface{}) error {
	body, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}
	return nil
}

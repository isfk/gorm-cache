package main

import (
	"encoding/json"
	"testing"
)

func TestJson(t *testing.T) {
	jstr := "{\"ID\":80,\"Username\":\"sfk\",\"CreatedAt\":1713512474,\"UpdatedAt\":1713512474}"

	type User struct {
		ID        int64
		Username  string
		CreatedAt int64
		UpdatedAt int64
	}

	user := &User{}
	err := json.Unmarshal([]byte(jstr), &user)
	if err != nil {
		t.Errorf("err %v", err)
		return
	}

	t.Logf("user %v", user)

}

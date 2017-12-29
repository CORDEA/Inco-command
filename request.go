package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	BaseUrl       = ""
	HistoriesPath = "/histories"
	LoginPath     = "/login"
)

type Request struct {
	token string
}

func NewRequest() *Request {
	return &Request{}
}

func (r *Request) DeleteHistories(histories []History) error {
	client := &http.Client{}
	query := "?id="
	for i := range histories {
		id := strconv.FormatInt(histories[i].ID, 10)
		if len(histories) <= i+1 {
			query = query + id
			continue
		}
		query = query + id + ","
	}
	req, err := http.NewRequest("DELETE", BaseUrl+HistoriesPath+query, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+r.token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		return nil
	}
	return errors.New("Failed")
}

func (r *Request) GetHistories() ([]History, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", BaseUrl+HistoriesPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+r.token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var histories []History
	if err := json.Unmarshal(body, &histories); err != nil {
		return nil, err
	}
	return histories, nil
}

func (r *Request) Login(user string, pass string) error {
	resp, err := http.Get(BaseUrl + LoginPath + "?user=" + user + "&pass=" + pass)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var login Login
	if err := json.Unmarshal(body, &login); err != nil {
		return err
	}

	r.token = login.Token
	return nil
}

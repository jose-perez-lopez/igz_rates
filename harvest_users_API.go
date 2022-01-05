package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func getHarvestActiveUsers() []interface{} {

	rows := 100
	var listSlice []interface{}
	for i := 1; rows == 100; i++ {
		results := getHarvestActiveUsersPage(i)
		rows = len(results)
		listSlice = append(listSlice, results...)
	}
	return listSlice
}

func getHarvestActiveUsersPage(page int) []interface{} {
	url := "https://api.harvestapp.com/v2/users?is_active=true&page=" + strconv.Itoa(page)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Go Harvest API Sample")
	req.Header.Set("Harvest-Account-ID", os.Getenv("HarvestAccountID"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("Authorization"))
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var jsonResponse map[string]interface{}
	json.Unmarshal(body, &jsonResponse)
	listSlice, _ := jsonResponse["users"].([]interface{})

	return listSlice
}

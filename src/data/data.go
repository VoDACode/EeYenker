package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func GetIdsFromFile() []int {
	file, err := os.Open("data/ids.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var ids []int
	err = json.Unmarshal(content, &ids)
	if err != nil {
		panic(err)
	}

	return ids
}

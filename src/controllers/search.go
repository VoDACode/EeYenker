package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type SearchResult struct {
	Items []struct {
		Type  string `json:"type"`
		Name  string `json:"name"`
		ID    int    `json:"id"`
		Price struct {
			Currency           string  `json:"currency"`
			Initial            float32 `json:"initial"`
			Final              float32 `json:"final"`
			DiscountPercentage int
		} `json:"price"`
		TinyImage string `json:"tiny_image"`
		Metascore string `json:"metascore"`
		Platforms struct {
			Windows bool `json:"windows"`
			Mac     bool `json:"mac"`
			Linux   bool `json:"linux"`
		} `json:"platforms"`
		StreamingVideo    bool   `json:"streamingvideo"`
		ControllerSupport string `json:"controller_support"`
	} `json:"items"`
}

func Search(c *gin.Context) {

	query := c.Query("name")
	searchURL := fmt.Sprintf("https://store.steampowered.com/api/storesearch/?term=%s&cc=ua", url.QueryEscape(query))

	resp, err := http.Get(searchURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}

	for i := range result.Items {
		result.Items[i].Price.Initial /= 100
		result.Items[i].Price.Final /= 100
		result.Items[i].Price.DiscountPercentage = (int)(100 - (result.Items[i].Price.Final * 100 / result.Items[i].Price.Initial))
	}

	c.HTML(http.StatusOK, "search.result", gin.H{
		"result": result,
	})
}

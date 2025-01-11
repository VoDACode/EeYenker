package main

import (
	"EeYenker/src/data"
	"EeYenker/src/models/responses"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

var db *sql.DB
var mu sync.Mutex

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "steam_data.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Apps (
		Id INTEGER PRIMARY KEY,
		Name TEXT
	)`)
	if err != nil {
		log.Fatalf("Failed to create Apps table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS OnlineHistory (
		AppId INTEGER,
		Count INTEGER,
		Datetime DATETIME,
		FOREIGN KEY(AppId) REFERENCES Apps(Id)
	)`)
	if err != nil {
		log.Fatalf("Failed to create OnlineHistory table: %v", err)
	}
}

func fetchGameStats(appId int) {
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUserStats/GetNumberOfCurrentPlayers/v1/?appid=%d", appId)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch data for app %d: %v", appId, err)
		return
	}
	defer resp.Body.Close()

	var steamResp responses.SteamResponse
	if err := json.NewDecoder(resp.Body).Decode(&steamResp); err != nil {
		log.Printf("Failed to decode response for app %d: %v", appId, err)
		return
	}

	if steamResp.Response.Result == 1 {
		mu.Lock()
		_, err = db.Exec("INSERT INTO OnlineHistory (AppId, Count, Datetime) VALUES (?, ?, ?)",
			appId, steamResp.Response.PlayerCount, time.Now())
		mu.Unlock()
		if err != nil {
			log.Printf("Failed to insert data for app %d: %v", appId, err)
		}
	}
}

func fetchAndStoreData() {
	rows, err := db.Query("SELECT Id FROM Apps")
	if err != nil {
		log.Printf("Failed to query apps: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var appId int
		if err := rows.Scan(&appId); err != nil {
			log.Printf("Failed to scan app ID: %v", err)
			continue
		}

		fetchGameStats(appId)
	}
}

func startScheduler() {
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for range ticker.C {
			fetchAndStoreData()
		}
	}()
}

func getAggregatedData(c *gin.Context) {
	appId := c.Query("appid")
	from := c.Query("from")
	to := c.Query("to")
	detail := c.Query("detail")

	layout := "2006-01-02T15:04:05"
	fromTime, err := time.Parse(layout, from)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'from' time format"})
		return
	}

	toTime, err := time.Parse(layout, to)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'to' time format"})
		return
	}

	groupFormat := "%Y-%m-%dT%H:00:00" // Default to hourly
	switch detail {
	case "12h":
		groupFormat = "%Y-%m-%dT%H:00:00"
	case "day":
		groupFormat = "%Y-%m-%d"
	case "week":
		groupFormat = "%Y-%W"
	case "month":
		groupFormat = "%Y-%m"
	case "all":
		groupFormat = "%Y"
	}

	query := fmt.Sprintf(`
		SELECT AppId, AVG(Count) as AvgCount, strftime('%s', Datetime) as TimeGroup
		FROM OnlineHistory
		WHERE AppId = ? AND Datetime BETWEEN ? AND ?
		GROUP BY TimeGroup
	`, groupFormat)

	rows, err := db.Query(query, appId, fromTime, toTime)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query data"})
		return
	}
	defer rows.Close()

	var results []gin.H
	for rows.Next() {
		// print results
		var appId int
		var avgCount float64
		var timeGroup sql.NullString
		if err := rows.Scan(&appId, &avgCount, &timeGroup); err != nil {
			log.Printf("Failed to scan data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan data"})
			return
		}
		results = append(results, gin.H{
			"appid":     appId,
			"avg_count": avgCount,
			"time":      timeGroup.String,
		})
	}

	c.JSON(http.StatusOK, results)
}

func MainPage(c *gin.Context) {
	url := "https://store.steampowered.com/api/appdetails?appids=%d&cc=ua"
	var result []responses.SteamGameResponse
	gameIds := data.GetIdsFromFile()
	for _, id := range gameIds {
		resp, err := http.Get(fmt.Sprintf(url, id))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var responseJson map[int]responses.SteamGameResponse
		if err := json.NewDecoder(resp.Body).Decode(&responseJson); err != nil {
			panic(err)
		}

		result = append(result, responseJson[id])
		result[len(result)-1].Data.CardImage = "https://cdn.cloudflare.steamstatic.com/steam/apps/" + fmt.Sprint(id) + "/header.jpg"
	}

	// query := c.Query("name")
	// searchURL := fmt.Sprintf("https://store.steampowered.com/api/storesearch/?term=%s&cc=ua", url.QueryEscape(query))

	// resp, err := http.Get(searchURL)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	// var result responses.SearchResult
	// if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	// 	panic(err)
	// }

	// for i := range result.Items {
	// 	result.Items[i].HeaderImage = "https://cdn.cloudflare.steamstatic.com/steam/apps/" + fmt.Sprint(result.Items[i].ID) + "/header.jpg"
	// 	result.Items[i].LibraryImage = "https://cdn.cloudflare.steamstatic.com/steam/apps/" + fmt.Sprint(result.Items[i].ID) + "/library_600x900_2x.jpg"
	// 	result.Items[i].BigImage = "https://cdn.cloudflare.steamstatic.com/steam/apps/" + fmt.Sprint(result.Items[i].ID) + "/capsule_616x353.jpg"
	// 	result.Items[i].Price.Initial /= 100
	// 	result.Items[i].Price.Final /= 100
	// 	result.Items[i].Price.DiscountPercentage = (int)(100 - (result.Items[i].Price.Final * 100 / result.Items[i].Price.Initial))

	// 	var exists bool
	// 	mu.Lock()
	// 	row := db.QueryRow("SELECT EXISTS(SELECT 1 FROM Apps WHERE Id = ?)", result.Items[i].ID)
	// 	if err := row.Scan(&exists); err != nil {
	// 		log.Printf("Failed to check existence for app %d: %v", result.Items[i].ID, err)
	// 		mu.Unlock()
	// 		continue
	// 	}
	// 	mu.Unlock()

	// 	if !exists {
	// 		mu.Lock()
	// 		_, err := db.Exec("INSERT INTO Apps (Id, Name) VALUES (?, ?)", result.Items[i].ID, result.Items[i].Name)
	// 		if err != nil {
	// 			log.Printf("Failed to insert app %d: %v", result.Items[i].ID, err)
	// 			mu.Unlock()
	// 		} else {
	// 			log.Printf("Inserted app %d: %s", result.Items[i].ID, result.Items[i].Name)
	// 			mu.Unlock()
	// 			fetchGameStats(result.Items[i].ID)
	// 		}
	// 	}
	// }

	c.HTML(http.StatusOK, "search.result", gin.H{
		"games": result,
	})
}

func main() {
	initDB()
	startScheduler()

	r := gin.Default()

	// Завантаження шаблонів
	r.LoadHTMLGlob("src/templates/**/*")

	r.Static("/static", "./src/static")
	r.Static("/assets", "./src/assets")

	// Реєстрація маршрутів
	r.GET("/", MainPage)
	r.GET("/api/aggregated_data", getAggregatedData)
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("favicon.ico")
	})

	// Виведення адреси сервера в консоль із підсвіткою
	port := "8070"
	url := fmt.Sprintf("http://localhost:%s/", port)

	// Використовуємо кольоровий текст
	c := color.New(color.FgHiGreen, color.Bold) // Високий зелений текст, жирний
	c.Printf("Сервер запущено на: %s\n", url)

	// Запуск сервера
	r.Run("0.0.0.0:" + port)
}

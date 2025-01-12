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

var mu sync.Mutex

func createDBConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "steam_data_1.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	_, err = db.Exec("PRAGMA busy_timeout = 5000;")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set busy timeout: %w", err)
	}

	return db, nil
}

func initDB() {
	db, err := createDBConnection()
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Apps (
		Id INTEGER PRIMARY KEY,
		Name TEXT
	)`) // Create Apps table
	if err != nil {
		log.Fatalf("Failed to create Apps table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS OnlineHistory (
		AppId INTEGER,
		Count INTEGER,
		Datetime DATETIME,
		FOREIGN KEY(AppId) REFERENCES Apps(Id)
	)`) // Create OnlineHistory table
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
		db, err := createDBConnection()
		if err != nil {
			log.Printf("Failed to create DB connection: %v", err)
			return
		}
		defer db.Close()

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
	db, err := createDBConnection()
	if err != nil {
		log.Printf("Failed to create DB connection: %v", err)
		return
	}
	defer db.Close()

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
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			log.Println("Fetching online data...")
			fetchAndStoreData()
			log.Println("Done fetching online data.")
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
	case "h":
		groupFormat = "%Y-%m-%d %H:00:00"
	case "d":
		groupFormat = "%Y-%m-%d"
	case "w":
		groupFormat = "%Y-%W"
	case "m":
		groupFormat = "%Y-%m"
	case "*":
		groupFormat = "%Y"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'detail' value"})
		return
	}

	db, err := createDBConnection()
	if err != nil {
		log.Printf("Failed to create DB connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query data"})
		return
	}
	defer db.Close()

	query := fmt.Sprintf(`
	SELECT strftime('%s', Datetime) AS Period,
	AVG(Count) AS AvgCount
	FROM OnlineHistory
	WHERE AppId = ? AND Datetime BETWEEN ? AND ?
	GROUP BY Period
	ORDER BY Period
	LIMIT 10000
	`, groupFormat)

	rows, err := db.Query(query, appId, fromTime, toTime)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query data"})
		return
	}
	defer rows.Close()

	var results []gin.H
	for rows.Next() {
		var avgCount float64
		var period sql.NullString
		if err := rows.Scan(&period, &avgCount); err != nil {
			log.Printf("Failed to scan data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan data"})
			return
		}
		avgCount = float64(int(avgCount*100)) / 100
		results = append(results, gin.H{
			"avg_count": avgCount,
			"period":    period.String,
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
			log.Printf("Failed to fetch game details for app %d: %v", id, err)
			continue
		}
		defer resp.Body.Close()

		var responseJson map[int]responses.SteamGameResponse
		if err := json.NewDecoder(resp.Body).Decode(&responseJson); err != nil {
			log.Printf("Failed to decode game details for app %d: %v", id, err)
			continue
		}

		result = append(result, responseJson[id])
		result[len(result)-1].Data.CardImage = "https://cdn.cloudflare.steamstatic.com/steam/apps/" + fmt.Sprint(id) + "/header.jpg"
		if result[len(result)-1].Data.Price != nil {
			result[len(result)-1].Data.Price.Initial /= 100
			result[len(result)-1].Data.Price.Final /= 100
			result[len(result)-1].Data.Price.DiscountPercentage = (int)(100 - (result[len(result)-1].Data.Price.Final * 100 / result[len(result)-1].Data.Price.Initial))
		}
	}

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
	r.GET("/api/stats", getAggregatedData)
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

package tests_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"shin/src/app"
	"shin/src/config"
	"shin/src/database"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var configPath string

func init() {
	// We back to root dir on execute tests
	os.Chdir("../")
	// Define a flag for the config path
	flag.StringVar(&configPath, "c", "test.config.yml", "Path to the configuration file")
}

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Suite")
}

func replaceAny(a, b gin.H) {
	for key, valueA := range a {
		if valueB, exists := b[key]; exists {
			// If the value in a is a map, recurse into it
			if mapA, ok := valueA.(gin.H); ok {
				if mapB, ok := valueB.(gin.H); ok {
					replaceAny(mapA, mapB)
				}
			} else if valueA == "<ANY>" {
				// Replace "<ANY>" with the corresponding value from b
				a[key] = valueB
			}
		}
	}
}

func bodyExpect(responseBody io.Reader, expect gin.H) {
	body := gin.H{}
	decoder := json.NewDecoder(responseBody)
	decoder.Decode(&body)
	replaceAny(expect, body)
	Expect(body).To(Equal(expect))
}

func setupTestEnvironment() (*sqlx.DB, *gin.Engine) {
	config.Init(configPath)
	db := database.Connect(&database.ConnectOption{
		URL:         config.Config.Database.URL,
		SqlDir:      config.Config.Database.SqlDir,
		MaxRequests: 5,
		Interval:    30 * time.Second,
		Timeout:     5 * time.Second,
	})
	m, err := migrate.New(
		fmt.Sprintf("file://%s", config.Config.Database.Migrations),
		config.Config.Database.URL,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	log.Println("Migrations applied successfully!")
	router := app.Init()

	return db, router
}

func teardownTestEnvironment(db *sqlx.DB) {
	db.Close()
	if err := database.DropDatabase(config.Config.Database.URL); err != nil {
		log.Fatalf("Dropping database %v", err)
	}
}

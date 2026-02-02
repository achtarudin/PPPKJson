package main

import (
	"context"
	"cutbray/pppk-json/cmd/config"
	"cutbray/pppk-json/internal/adapters/db_adapter"
	"cutbray/pppk-json/internal/adapters/logger"
	"cutbray/pppk-json/internal/repositories/models"
	"cutbray/pppk-json/internal/utils"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

// QuestionData represents the JSON structure from data files
type QuestionData struct {
	ID           string       `json:"id"`
	Category     string       `json:"category"`
	QuestionText string       `json:"question_text"`
	Options      []OptionData `json:"options"`
}

type OptionData struct {
	OptionText string `json:"option_text"`
	Score      int    `json:"score"`
}

func main() {
	logger.New()

	err := config.LoadEnvFile()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := utils.GetEnvOrDefault("DB_HOST", "localhost")
	dbUser := utils.GetEnvOrDefault("DB_USER", "encang_cutbray")
	dbPassword := utils.GetEnvOrDefault("DB_PASSWORD", "encang_cutbray")
	dbName := utils.GetEnvOrDefault("DB_NAME", "togotestgo")
	dbPort := utils.GetEnvOrDefault("DB_PORT", "5432")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	dbAdapter := db_adapter.New(dsn, &gorm.Config{})

	connectManagers := []config.ConnectManager{
		{Name: "Postgres DB", Adapter: dbAdapter},
	}

	if err := config.ConnectAdapters(context.Background(), connectManagers...); err != nil {
		log.Fatalf("%v", err)
	}

	db, ok := dbAdapter.Value().(*gorm.DB)
	if !ok {
		log.Fatalf("Database adapter is not properly initialized")
	}

	questions, err := readJsonFile("./migrations/data")
	if err != nil {
		log.Fatalf("Failed to read JSON files: %v", err)
	}

	log.Printf("Total questions read from JSON files: %d", len(questions))
	if err := seedQuestions(db, questions); err != nil {
		log.Fatalf("Failed to seed questions: %v", err)
	}

	log.Println("Seeding completed successfully!")
	config.DisconnectAdapters(connectManagers...)
}

func readJsonFile(filePath string) ([]QuestionData, error) {

	var allQuestions []QuestionData

	err := filepath.WalkDir(filePath, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".json" {
			var questions []QuestionData

			byteValue, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			err = json.Unmarshal(byteValue, &questions)
			if err != nil {
				return err
			}

			allQuestions = append(allQuestions, questions...)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return allQuestions, nil
}

func seedQuestions(db *gorm.DB, questions []QuestionData) error {

	return db.WithContext(context.Background()).Transaction(func(tx *gorm.DB) error {

		// Truncate existing questions and options
		if err := tx.Exec("TRUNCATE TABLE questions RESTART IDENTITY CASCADE").Error; err != nil {
			return fmt.Errorf("failed to truncate table: %v", err)
		}

		for _, q := range questions {
			question := models.Question{
				Category:     q.Category,
				QuestionText: q.QuestionText,
			}

			options := make([]models.QuestionOption, len(q.Options))
			for i, opt := range q.Options {
				if opt.OptionText == "" {
					log.Printf("[Error Skipping] option with empty text for question: %s (category: %s, text: %s)", q.QuestionText, question.Category, question.QuestionText)
					return fmt.Errorf("option text cannot be empty")
				}
				options[i] = models.QuestionOption{
					OptionText: opt.OptionText,
					Score:      opt.Score,
				}
			}

			question.Options = options

			result := tx.Create(&question)

			if err := result.Error; err != nil {
				return fmt.Errorf("failed to insert question: %v", err)
			}
		}
		return nil
	})

}

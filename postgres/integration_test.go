package postgres

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"os"
	"testing"
)

type Config struct {
	DbUser     string
	DbPassword string
	DbPort     string
	DbHost     string
	DbName     string
	ApiSecret  string
}

type TestSuite struct {
	config Config
	db *DB
}

func (s *TestSuite) GetDb(t *testing.T) *DB {
	if s.db != nil {
		return s.db
	}

	config, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}

	s.config = config

	dbUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", s.config.DbHost, s.config.DbPort, s.config.DbUser, s.config.DbName, s.config.DbPassword)

	db, err := Open(dbUrl)
	if err != nil {
		t.Fatal("cannot connect to db", err)
	}

	s.db = db
	return db
}

func (s *TestSuite) CleanDb(t *testing.T) {
	_, err := s.db.Exec("DELETE FROM articles WHERE true")
	if err != nil {
		t.Fatal("error deleting articles", err)
	}
	_, err = s.db.Exec("DELETE FROM users WHERE true")
	if err != nil {
		t.Fatal("error deleting users", err)
	}
}

func loadConfig() (Config,error) {

	if os.Getenv("ENV_FILE") == "" {
		return Config{}, errors.New("env file is required")
	}

	viper.SetConfigFile(fmt.Sprintf("../%s", os.Getenv("ENV_FILE")))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	return Config{
		DbUser:     viper.GetString("DB_USER"),
		DbPassword: viper.GetString("DB_PASSWORD"),
		DbPort:     viper.GetString("DB_PORT"),
		DbHost:     viper.GetString("DB_HOST"),
		DbName:     viper.GetString("DB_NAME"),
		ApiSecret:  viper.GetString("API_SECRET"),
	}, nil
}


var Suite = TestSuite{}


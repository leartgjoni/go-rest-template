package main

import (
	"fmt"
	"github.com/leartgjoni/go-rest-template/http"
	"github.com/leartgjoni/go-rest-template/postgres"
	"github.com/spf13/viper"
	"io"
	"os"
	"os/signal"
)

func main() {
	m := NewMain()

	// Parse command line flags.
	// Todo

	// Load configuration.
	if err := m.LoadConfig(); err != nil {
		fmt.Fprintln(m.Stderr, err)
		os.Exit(1)
	}

	// Execute program.
	if err := m.Run(); err != nil {
		fmt.Fprintln(m.Stderr, err)
		os.Exit(1)
	}

	// Shutdown on SIGINT (CTRL-C).
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Fprintln(m.Stdout, "received interrupt, shutting down...")
}

// Main represents the main program execution.
type Main struct {
	ConfigPath string
	Config     Config

	// Input/output streams
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	closeFn func() error
}

// NewMain returns a new instance of Main.
func NewMain() *Main {
	return &Main{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,

		closeFn: func() error { return nil },
	}
}

// Close cleans up the program.
func (m *Main) Close() error { return m.closeFn() }

// LoadConfig parses the configuration file.
func (m *Main) LoadConfig() error {

	if os.Getenv("CONFIG_PATH") != "" {
		m.ConfigPath = os.Getenv("CONFIG_PATH")
	} else {
		m.ConfigPath = ".env"
	}

	viper.SetConfigFile(m.ConfigPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	m.Config = Config{
		DbUser: viper.GetString("DB_USER"),
		DbPassword: viper.GetString("DB_PASSWORD"),
		DbPort: viper.GetString("DB_PORT"),
		DbHost: viper.GetString("DB_HOST"),
		DbName: viper.GetString("DB_NAME"),
		ApiSecret: viper.GetString("API_SECRET"),
	}

	return nil
}

func (m *Main) Run() error {
	dbUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", m.Config.DbHost, m.Config.DbPort, m.Config.DbUser, m.Config.DbName, m.Config.DbPassword)
	db, err := postgres.Open(dbUrl)
	if err != nil {
		fmt.Println(m.Stderr, err)
		os.Exit(1)
	}

	// Initialize postgres services.
	userService := postgres.NewUserService(db, m.Config.ApiSecret)
	articleService := postgres.NewArticleService(db)

	// Initialize Http server.
	httpServer := http.NewServer()
	httpServer.Addr = ":8080"

	httpServer.UserService = userService
	httpServer.ArticleService = articleService

	// Open HTTP server.
	if err := httpServer.Open(); err != nil {
		return err
	}
	fmt.Fprintf(m.Stdout, "Listening on port: %s\n", httpServer.Addr)

	// Assign close function.
	m.closeFn = func() error {
		httpServer.Close()
		db.Close()
		return nil
	}

	return nil
}

type Config struct {
	DbUser string
	DbPassword string
	DbPort string
	DbHost string
	DbName string
	ApiSecret string
}
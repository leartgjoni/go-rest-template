package postgres_test

import (
	app "github.com/leartgjoni/go-rest-template"
	"github.com/leartgjoni/go-rest-template/postgres"
	"testing"
	"time"
)

func TestUserServiceIntegration_GetById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := Suite.GetDb(t)
	Suite.CleanDb(t)

	timeNow := time.Now()
	// expected user
	eUser := app.User{
		ID: 0,
		Username: "test",
		Email: "test@test.com",
		Password: "hashed_password",
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	row := db.QueryRow("INSERT INTO users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at", eUser.Username, eUser.Email, eUser.Password, eUser.CreatedAt, eUser.UpdatedAt)
	if err := row.Scan(&eUser.ID, &eUser.CreatedAt, &eUser.UpdatedAt); err != nil {
		t.Fatal("error while inserting user", err)
	}

	us := postgres.NewUserService(&postgres.DB{DB: db}, "random-api-string")
	// actual user
	aUser, err := us.GetById(eUser.ID)
	if err != nil {
		t.Fatal("error getting the user", err)
	}

	if aUser.ID != eUser.ID || aUser.Username != eUser.Username || aUser.Email != eUser.Email || aUser.Password != eUser.Password || aUser.CreatedAt != eUser.CreatedAt || aUser.UpdatedAt != eUser.UpdatedAt {
		t.Errorf("Expected %v but got %v", eUser, aUser)
	}
	// TODO: figure out time thing with db and coverage
}


package postgres

import (
	app "github.com/leartgjoni/go-rest-template"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestUserServiceIntegration_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("normal save", func(t *testing.T) {
		db := Suite.GetDb(t)
		Suite.CleanDb(t)

		timeNow := time.Now()
		// expected user
		user := app.User{
			ID:        0,
			Username:  "test",
			Email:     "test@test.com",
			Password:  "hashed_password",
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		}

		us := NewUserService(db, "random-api-string")
		if err := us.Save(&user); err != nil {
			t.Fatal("cannot save user", err)
		}

		if user.ID == 0 {
			t.Fatal("user id still zero")
		}

		var dbUser app.User
		if err := db.QueryRow("SELECT * FROM users WHERE id = $1", user.ID).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.Password, &dbUser.CreatedAt, &dbUser.UpdatedAt); err != nil {
			t.Fatal("cannot read user from db", err)
		}

		if user.ID != dbUser.ID || user.Username != dbUser.Username || user.Email != dbUser.Email || user.Password != dbUser.Password || !user.CreatedAt.Equal(dbUser.CreatedAt) || !user.UpdatedAt.Equal(dbUser.UpdatedAt) {
			t.Errorf("Expected %v but got %v", user, dbUser)
		}
	})

	t.Run("email already exists", func(t *testing.T) {
		db := Suite.GetDb(t)
		Suite.CleanDb(t)

		email := "test@test.com"

		if _, err := db.Exec("INSERT INTO users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", "random username", email, "random-password", time.Now(), time.Now()); err != nil {
			t.Fatal("cannot insert user", err)
		}

		user := app.User{
			ID:        0,
			Username:  "test",
			Email:     email,
			Password:  "hashed_password",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		us := NewUserService(db, "random-api-string")
		err := us.Save(&user)
		if err != app.ErrEmailAlreadyUsed {
			t.Fatal("incorrect error", err)
		}
	})
}

func TestUserServiceIntegration_GetById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := Suite.GetDb(t)
	Suite.CleanDb(t)

	timeNow := time.Now()
	// expected user
	eUser := app.User{
		ID:        0,
		Username:  "test",
		Email:     "test@test.com",
		Password:  "hashed_password",
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	row := db.QueryRow("INSERT INTO users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id", eUser.Username, eUser.Email, eUser.Password, eUser.CreatedAt, eUser.UpdatedAt)
	if err := row.Scan(&eUser.ID); err != nil {
		t.Fatal("error while inserting user", err)
	}

	us := NewUserService(db, "random-api-string")
	// actual user
	aUser, err := us.GetById(eUser.ID)
	if err != nil {
		t.Fatal("error getting the user", err)
	}

	if aUser.ID != eUser.ID || aUser.Username != eUser.Username || aUser.Email != eUser.Email || aUser.Password != eUser.Password || !aUser.CreatedAt.Equal(eUser.CreatedAt) || !aUser.UpdatedAt.Equal(eUser.UpdatedAt) {
		t.Errorf("Expected %v but got %v", eUser, aUser)
	}
}

func TestUserServiceIntegration_Login(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("normal login", func(t *testing.T) {
		db := Suite.GetDb(t)
		Suite.CleanDb(t)

		timeNow := time.Now()
		// expected user
		user := app.User{
			ID:        0,
			Username:  "test",
			Email:     "test@test.com",
			Password:  "password",
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			t.Fatal("cannot hash passowrd", err)
		}

		if _, err := db.Exec("INSERT INTO users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", user.Username, user.Email, hashedPassword, user.CreatedAt, user.UpdatedAt); err != nil {
			t.Fatal("cannot insert user", err)
		}

		us := NewUserService(db, "random-api-string")
		token, err := us.Login(&user)
		if !(token != "" && err == nil) {
			t.Fatal("err with login", err)
		}
	})

	t.Run("wrong credentials", func(t *testing.T) {
		db := Suite.GetDb(t)
		Suite.CleanDb(t)

		timeNow := time.Now()
		// expected user
		user := app.User{
			ID:        0,
			Username:  "test",
			Email:     "test@test.com",
			Password:  "password",
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			t.Fatal("cannot hash passowrd", err)
		}

		if _, err := db.Exec("INSERT INTO users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", user.Username, user.Email, hashedPassword, user.CreatedAt, user.UpdatedAt); err != nil {
			t.Fatal("cannot insert user", err)
		}

		us := NewUserService(db, "random-api-string")
		// change password
		user.Password = "password-edit"
		_, err = us.Login(&user)
		if err != app.ErrWrongCredentials {
			t.Fatal("error incorrect", err)
		}
	})
}

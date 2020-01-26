package postgres

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	app "github.com/leartgjoni/go-rest-template"
	"testing"
)

func TestUserService_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name string
		countResult *sqlmock.Rows
		insertResult *sqlmock.Rows
		expected error
	}{
		{
			name: "Email already used",
			countResult: sqlmock.NewRows([]string{"count"}).
				AddRow(1),
			insertResult: nil,
			expected: errors.New("email already in use"),
		},
		{
			name: "User without ID after saving",
			countResult: sqlmock.NewRows([]string{"count"}).
				AddRow(0),
			insertResult: sqlmock.NewRows([]string{"id"}).
				AddRow(0),
			expected: errors.New("unable to save"),
		},
		{
			name: "Success",
			countResult: sqlmock.NewRows([]string{"count"}).
				AddRow(0),
			insertResult: sqlmock.NewRows([]string{"id"}).
				AddRow(1),
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock.ExpectQuery("^SELECT (.+) FROM users*").WillReturnRows(test.countResult)
			if test.insertResult != nil {
				mock.ExpectQuery("^INSERT INTO users *").WillReturnRows(test.insertResult)
			}

			us := NewUserService(&DB{db})

			err = us.Save(&app.User{})

			// we make sure that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			if err != nil {
				if err.Error() != test.expected.Error() {
					t.Fatalf("wrong error. expected %s but got %s", test.expected, err)
				}
			} else {
				if err != test.expected {
					t.Fatalf("wrong error. expected %s but got %s", test.expected, err)
				}
			}
		})
	}
}
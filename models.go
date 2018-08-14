package courier

import (
	"time"

	"github.com/jmoiron/sqlx"
	// test
	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE account (
	id SERIAL PRIMARY KEY,
	phone VARCHAR(12) NOT NULL
);

CREATE TABLE message (
	id SERIAL PRIMARY KEY,
	sender_id INTEGER REFERENCES account(id),
	message TEXT NOT NULL,
	status SMALLINT NOT NULL,
	scheduledto TIMESTAMP
);

CREATE TABLE license (
	id SERIAL PRIMARY KEY,
	account_id INTEGER REFERENCES account(id),
	kind INTEGER NOT NULL,
	status INTEGER NOT NULL,
	expiration TIMESTAMP NOT NULL
);`

type AccountModel struct {
	ID    uint
	Phone string
}

type MessageModel struct {
	ID          uint
	SenderID    uint
	Message     string
	Status      int
	ScheduledTo time.Time
}

type LicenseModel struct {
	ID         uint
	AccountID  uint
	Kind       int
	Status     int
	Expiration time.Time
}

func CreateAccount(db *sqlx.DB, account AccountModel) {
	tx := db.MustBegin()

	tx.MustExec(`
		INSERT INTO
			account (phone)
		VALUES
			($1)
	`, account.Phone)

	tx.Commit()
}

func Migrate() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=wpcourier sslmode=disable")
	if err != nil {
		return nil, err
	}

	// TODO:
	// Check if the schema was running
	// Tip: To use migrations
	db.MustExec(schema)

	return db, nil
}

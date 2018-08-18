package courier

import (
	"time"

	"github.com/jmoiron/sqlx"
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
	receiver_phone VARCHAR(12) NOT NULL,
	message TEXT NOT NULL,
	status SMALLINT NOT NULL DEFAULT 0,
	scheduledto TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE license (
	id SERIAL PRIMARY KEY,
	account_id INTEGER REFERENCES account(id),
	kind INTEGER NOT NULL,
	status INTEGER NOT NULL,
	expiration TIMESTAMP NOT NULL
);`

type AccountModel struct {
	ID       uint
	Phone    string
	License  *LicenseModel  `db:"-"`
	Messages []MessageModel `db:"-"`
}

type LicenseModel struct {
	ID         uint
	AccountID  uint
	Kind       int
	Status     int
	Expiration time.Time
	Account    *AccountModel `db:"-"`
}

type MessageModel struct {
	ID            uint
	SenderID      uint   `db:"sender_id"`
	ReceiverPhone string `db:"receiver_phone"`
	Message       string
	Status        int
	ScheduledTo   time.Time
	account       *AccountModel `db:"-"`
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

func (m *MessageModel) GetAccount(db *sqlx.DB) (AccountModel, error) {
	if m.account != nil {
		return *m.account, nil
	}

	m.account = &AccountModel{}
	err := db.Get(m.account, "SELECT * FROM account WHERE id = $1", m.SenderID)

	if err != nil {
		return AccountModel{}, err
	}

	return *m.account, nil
}

func (m *MessageModel) SetSent(db *sqlx.DB) error {
	query := "UPDATE message SET status = 1 WHERE id = :id"
	_, err := db.NamedExec(query, m)

	if err != nil {
		return err
	}

	return nil
}

func (m *MessageModel) SetError(db *sqlx.DB) error {
	query := "UPDATE message SET status = 2 WHERE id = :id"
	_, err := db.NamedExec(query, m)

	if err != nil {
		return err
	}

	return nil
}

func Migrate() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=wpcourier sslmode=disable")
	if err != nil {
		return nil, err
	}

	// TODO:
	// Verify that migration has already been performed
	// Tip: To use migrations
	// db.MustExec(schema)

	return db, nil
}

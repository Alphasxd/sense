package db

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func CreateDatabase() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", "./orders.db")
    if err != nil {
        return nil, err
    }

    createTableSQL := `CREATE TABLE IF NOT EXISTS orders (
        "id" INTEGER PRIMARY KEY AUTOINCREMENT,
        "uid" INTEGER,
        "weight" REAL,
        "cost" INTEGER,
        "created_at" DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
    _, err = db.Exec(createTableSQL)
    if err != nil {
        return nil, err
    }

    return db, nil
}
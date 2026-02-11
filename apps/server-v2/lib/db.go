package lib

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Product struct {
	ID          int64    `json:"id"`
	Name        *string  `json:"name"`
	Manufacture *string  `json:"manufacture"`
	Output      *float64 `json:"output"`
	Price       *int64   `json:"price"`
	Width       *float64 `json:"width"`
	Height      *float64 `json:"height"`
}

func ScanProduct(scanner interface{ Scan(...any) error }) (Product, error) {
	var p Product
	err := scanner.Scan(&p.ID, &p.Name, &p.Manufacture, &p.Output, &p.Price, &p.Width, &p.Height)
	return p, err
}

func InitDb(dbPath string) error {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return err
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			manufacture TEXT,
			output REAL,
			price INTEGER,
			width REAL,
			height REAL
		);
	`)
	if err != nil {
		return err
	}

	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM products;").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = DB.Exec(`
			INSERT INTO products (name, manufacture, output, price, width, height) VALUES
				('Air Source Heat Pump', 'Generic Manufacturer', 12.5, 45000, 80, 120);
		`)
		if err != nil {
			return err
		}
		fmt.Println("Database seeded with initial data")
	}

	fmt.Println("Database initialized")
	return nil
}

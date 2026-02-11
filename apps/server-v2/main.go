package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"server-v2/lib"
)

// GET /products
func getProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := lib.DB.Query("SELECT id, name, manufacture, output, price, width, height FROM products ORDER BY id;")
	if err != nil {
		lib.WriteJSON(w, 500, map[string]string{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	products := []lib.Product{}
	for rows.Next() {
		p, err := lib.ScanProduct(rows)
		if err != nil {
			lib.WriteJSON(w, 500, map[string]string{"error": "Internal server error"})
			return
		}
		products = append(products, p)
	}

	lib.WriteJSON(w, 200, products)
}

// GET /products/{id}
func getProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	row := lib.DB.QueryRow("SELECT id, name, manufacture, output, price, width, height FROM products WHERE id = ?;", id)
	p, err := lib.ScanProduct(row)
	if err == sql.ErrNoRows {
		lib.WriteJSON(w, 404, map[string]string{"error": "Product not found"})
		return
	}
	if err != nil {
		lib.WriteJSON(w, 500, map[string]string{"error": "Internal server error"})
		return
	}

	lib.WriteJSON(w, 200, p)
}

// POST /products
func createProduct(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        *string  `json:"name"`
		Manufacture *string  `json:"manufacture"`
		Output      *float64 `json:"output"`
		Price       *int64   `json:"price"`
		Width       *float64 `json:"width"`
		Height      *float64 `json:"height"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		lib.WriteJSON(w, 400, map[string]string{"error": "Invalid request body"})
		return
	}

	result, err := lib.DB.Exec(
		"INSERT INTO products (name, manufacture, output, price, width, height) VALUES (?, ?, ?, ?, ?, ?);",
		body.Name, body.Manufacture, body.Output, body.Price, body.Width, body.Height,
	)
	if err != nil {
		lib.WriteJSON(w, 500, map[string]string{"error": "Internal server error"})
		return
	}

	id, _ := result.LastInsertId()
	row := lib.DB.QueryRow("SELECT id, name, manufacture, output, price, width, height FROM products WHERE id = ?;", id)
	p, err := lib.ScanProduct(row)
	if err != nil {
		lib.WriteJSON(w, 500, map[string]string{"error": "Internal server error"})
		return
	}

	lib.WriteJSON(w, 201, p)
}

// PUT /products/{id}
func updateProduct(w http.ResponseWriter, r *http.Request) {
	// Fail randomly 50% of the time
	if rand.Float64() < 0.5 {
		lib.WriteJSON(w, 500, map[string]string{"error": "Internal server error"})
		return
	}

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var body struct {
		Name        *string  `json:"name"`
		Manufacture *string  `json:"manufacture"`
		Output      *float64 `json:"output"`
		Price       *int64   `json:"price"`
		Width       *float64 `json:"width"`
		Height      *float64 `json:"height"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		lib.WriteJSON(w, 400, map[string]string{"error": "Invalid request body"})
		return
	}

	result, err := lib.DB.Exec(
		"UPDATE products SET name = ?, manufacture = ?, output = ?, price = ?, width = ?, height = ? WHERE id = ?;",
		id, body.Name, body.Manufacture, body.Output, body.Price, body.Width, body.Height,
	)
	if err != nil {
		lib.WriteJSON(w, 500, map[string]string{"error": "Internal server error"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		lib.WriteJSON(w, 404, map[string]string{"error": "Product not found"})
		return
	}

	row := lib.DB.QueryRow("SELECT id, name, manufacture, output, price, width, height FROM products WHERE id = ?;", id)
	p, err := lib.ScanProduct(row)
	if err != nil {
		lib.WriteJSON(w, 500, map[string]string{"error": "Internal server error"})
		return
	}

	lib.WriteJSON(w, 200, p)
}

// DELETE /products/{id}
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	result, err := lib.DB.Exec("DELETE FROM products WHERE id = ? RETURNING *;", id)
	if err != nil {
		lib.WriteJSON(w, 500, map[string]string{"error": "Internal server error"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		lib.WriteJSON(w, 404, map[string]string{"error": "Product not found"})
		return
	}

	lib.WriteJSON(w, 200, map[string]string{"message": "Product deleted"})
}

func main() {
	if err := lib.InitDb("./data/products.db"); err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Get("/products", getProducts)
	r.Post("/products", createProduct)
	r.Get("/products/{id}", getProduct)
	r.Put("/products/{id}", updateProduct)
	r.Delete("/products/{id}", deleteProduct)

	port := 3002
	fmt.Printf("Server running on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// CalculateShippingCost 计算快递费用
func CalculateShippingCost(weight float64) int {
	chargeableWeight := math.Ceil(weight)
	if chargeableWeight > 100 {
		chargeableWeight = 100
	}

	if chargeableWeight <= 1 {
		return 18
	}

	cost := 18
	for i := 2.0; i <= chargeableWeight; i++ {
		cost += 5
		cost = int(math.Round(float64(cost) * 1.01))
	}

	return cost
}

func createDatabase() (*sql.DB, error) {
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

func generateTestData(db *sql.DB) error {
	rand.Seed(time.Now().UnixNano())

	// 生成1000个用户id
	userIDs := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		userIDs[i] = i + 1
	}

	// 生成100000条订单记录
	for i := 0; i < 100000; i++ {
		uid := userIDs[rand.Intn(len(userIDs))]
		weight := generateWeight()
		cost := CalculateShippingCost(weight)

		_, err := db.Exec("INSERT INTO orders (uid, weight, cost) VALUES (?, ?, ?)", uid, weight, cost)
		if err != nil {
			return err
		}
	}

	return nil
}

// generateWeight 生成随机重量，保证计费重量的分布权重大致为 1/W
func generateWeight() float64 {
    r := rand.Float64()
    weight := 1 / math.Pow(r, 2)
	// 限制最大重量为 100KG
    if weight > 100 {
        weight = 100
    }
	// 保留两位小数
    return math.Ceil(weight)
}

func queryOrders(db *sql.DB, uid int) {
	rows, err := db.Query("SELECT id, weight, cost, created_at FROM orders WHERE uid = ?", uid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var totalCost int
	fmt.Printf("Orders for user %d:\n", uid)
	for rows.Next() {
		var id int
		var weight float64
		var cost int
		var createdAt string
		err := rows.Scan(&id, &weight, &cost, &createdAt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Order ID: %d, Weight: %.2f KG, Cost: %d 元, Created At: %s\n", id, weight, cost, createdAt)
		totalCost += cost
	}
	fmt.Printf("Total Cost: %d 元\n", totalCost)
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go run main.go generate  - 生成测试数据")
	fmt.Println("  go run main.go query <user_id>  - 查询指定用户的订单")
}

func main() {
	db, err := createDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	switch command {
	case "generate":
		err = generateTestData(db)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Test data generated successfully.")
	case "query":
		if len(os.Args) < 3 {
			printUsage()
			return
		}
		uid, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		queryOrders(db, uid)
	default:
		printUsage()
	}
}

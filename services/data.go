package services

import (
    "database/sql"
    "log"
    "math/rand"
    "sense/models"
    "sense/utils"
    "time"
)

func GenerateTestData(db *sql.DB) error {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    // 生成1000个用户id
    userIDs := make([]int, 1000)
    for i := 0; i < 1000; i++ {
        userIDs[i] = i + 1
    }

    // 开始事务
    tx, err := db.Begin()
    if err != nil {
        return err
    }

    stmt, err := tx.Prepare("INSERT INTO orders (uid, weight, cost) VALUES (?, ?, ?)")
    if err != nil {
        return err
    }
    defer func() {
        if err := stmt.Close(); err != nil {
            log.Fatal(err)
        }
    }()

    // 生成100000条订单记录
    for i := 0; i < 100000; i++ {
        uid := userIDs[r.Intn(len(userIDs))]
        weight := utils.GenerateWeight()
        cost := CalculateShippingCost(weight)

        _, err := stmt.Exec(uid, weight, cost)
        if err != nil {
            tx.Rollback()
            return err
        }
    }

    // 提交事务
    err = tx.Commit()
    if err != nil {
        return err
    }

    return nil
}

func QueryOrders(db *sql.DB, uid int) {
    rows, err := db.Query("SELECT id, uid, weight, cost, created_at FROM orders WHERE uid = ?", uid)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var totalCost int
    log.Printf("Orders for user %d:\n", uid)
    for rows.Next() {
        var order models.Order
        err := rows.Scan(&order.ID, &order.UID, &order.Weight, &order.Cost, &order.CreatedAt)
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("Order ID: %d, Weight: %.2f KG, Cost: %d 元, Created At: %s\n", order.ID, order.Weight, order.Cost, order.CreatedAt)
        totalCost += order.Cost
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
    log.Printf("Total Cost: %d 元\n", totalCost)
}
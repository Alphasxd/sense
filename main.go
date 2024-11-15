package main

import (
    "flag"
    "fmt"
    "log"
    "sense/db"
    "sense/services"
)

func printUsage() {
    fmt.Println("Usage:")
    fmt.Println("  -generate  生成测试数据")
    fmt.Println("  -query <user_id>  查询指定用户的订单")
}

func main() {
    generateFlag := flag.Bool("generate", false, "生成测试数据")
    queryFlag := flag.Int("query", 0, "查询指定用户的订单")
    flag.Usage = printUsage
    flag.Parse()

    database, err := db.CreateDatabase()
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()

    if *generateFlag {
        err = services.GenerateTestData(database)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println("Test data generated successfully.")
    } else if *queryFlag > 0 {
        services.QueryOrders(database, *queryFlag)
    } else {
        printUsage()
    }
}
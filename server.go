package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/takahiro-1227/go-orders-graphql-api/graph"
	"github.com/takahiro-1227/go-orders-graphql-api/graph/generated"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB;

func initDB() {
	var err error
	dataSourceName := "root:@tcp(localhost:3306)/?parseTime=True"
	db, err = gorm.Open("mysql", dataSourceName)

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	db.LogMode(true)

	db.Exec("CREATE DATABASE test_db")
	db.Exec("USE test_db")

	db.AutoMigrate(&models.Order{}, &models.Item{})
}

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	initDB();
	srv := handler.GraphQL(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		DB: db,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}



package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/egaevan/merchant/config"
	"github.com/egaevan/merchant/delivery/rest"
	"github.com/egaevan/merchant/repository"
	"github.com/egaevan/merchant/usecase"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func main() {
	// Init config
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Init echo framework
	e := echo.New()

	// Init DB
	mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)

	db, err := sql.Open("mysql", mysqlInfo)
	if err != nil {
		fmt.Println("Error connect to db")
		log.Fatal(err)
	}

	defer db.Close()

	// Init repository
	productRepo := repository.NewProductRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Init usecase
	productUsecae := usecase.NewProduct(productRepo)
	userUsecae := usecase.NewUser(userRepo)

	// Init handler
	rest.NewHandler(e, productUsecae, userUsecae)

	e.Logger.Fatal(e.Start(":8080"))
}

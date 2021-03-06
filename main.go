package main

import (
	"allbooks/routing"
	"flag"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"log"
	"net/http"
)

func main() {
	port := flag.Int("port", 8081, "specify the server port")
	bind := flag.String("bind", "", "specify IP to bind to")
	address := fmt.Sprintf("%s:%d", *bind, *port)

	router := routing.Router(AppRoutes())

	log.Println("Starting listening on ", address)

	err := http.ListenAndServe(address, router)
	log.Fatal(err)
}

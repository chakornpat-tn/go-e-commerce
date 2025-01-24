package main

import (
	"os"

	"github.com/chakornpat-tn/go-rest-api/config"
	"github.com/chakornpat-tn/go-rest-api/modules/servers"
	"github.com/chakornpat-tn/go-rest-api/pkg/database"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := database.DbConnect(cfg.DB())
	defer db.Close()

	servers.NewServer(cfg, db).Start()

}

package main

import (
	"os"

	"github.com/TurterDev/basicgolang/config"
	"github.com/TurterDev/basicgolang/modules/servers"
	"github.com/TurterDev/basicgolang/pkg/databases"
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
	// fmt.Println(cfg.App())
	// fmt.Println(cfg.Db())
	// fmt.Println(cfg.Jwt())
	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	// fmt.Println(db)
	servers.NewServer(cfg, db).Start()
}

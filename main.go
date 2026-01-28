package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/jwoodsiii/blogator/internal/config"
	"github.com/jwoodsiii/blogator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	dbQueries := database.New(db)

	s := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", handlerFollow)
	cmds.register("following", handlerFollowing)

	if len(os.Args) < 2 {
		log.Fatal("Not enough args, require command to execute")
	}
	userInput := os.Args[1:]

	cmdName := userInput[0]
	cmdArgs := userInput[1:]

	if err := cmds.run(s, command{Name: cmdName, Args: cmdArgs}); err != nil {
		log.Fatalf("Error running command: %v", err)
	}
	//fmt.Printf("Current config: %v\n", cfg)
}

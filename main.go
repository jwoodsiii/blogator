package main

import (
	"os"
	"log"
	"github.com/jwoodsiii/blogator/internal/config"
)

type state struct {
	cfg		*config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// cfg.SetUser("jwoodsiii")
	s := &state{
		cfg: &cfg,
	}

	cmds := commands{
			handlers: make(map[string]func(*state, command) error),
		}
	cmds.register("login", handlerLogin)

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

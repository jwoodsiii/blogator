package main

import (
	"context"
	"fmt"
	"os"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	user := cmd.Args[0]
	dbUser, err := s.db.GetUser(context.Background(), user)
	if err != nil {
		fmt.Errorf("Error attempting to login as user: %s, %v", user, err)
		os.Exit(1)
	}

	if err := s.cfg.SetUser(dbUser.Name); err != nil {
		return err
	}

	fmt.Printf("user has been set to: %s", dbUser.Name)
	return nil
}

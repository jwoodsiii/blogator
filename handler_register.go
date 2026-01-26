package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jwoodsiii/blogator/internal/database"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: name})
	if err != nil {
		return fmt.Errorf("Error getting user: %v", err)
	}
	// TODO: add check if user already exists

	s.cfg.CurrentUserName = user.Name
	fmt.Printf("User: %s was created\n", user.Name)
	log.Printf("New user: %s was created at: %s with uuid: %s\n", user.Name, user.CreatedAt, user.ID)

	return nil
}

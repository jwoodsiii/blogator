package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jwoodsiii/blogator/internal/database"
)

func handlerUsers(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage %s", cmd.Name)
	}

	currUser := s.cfg.CurrentUserName
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error listing users: %v", err)
	}

	for _, user := range users {
		out := user.Name
		if user.Name == currUser {
			out += " (current)"
		}
		fmt.Printf("* %s\n", out)
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	if err := s.db.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("Failed to delete users: %v", err)
	}
	fmt.Println("Users deleted...")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Name: name})
	if err != nil {
		return fmt.Errorf("Error creating user: %v", err)
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Printf("User: %s was created\n", user.Name)
	log.Printf("New user: %s was created at: %s with uuid: %s\n", user.Name, user.CreatedAt, user.ID)

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	user := cmd.Args[0]
	dbUser, err := s.db.GetUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("Error attempting to login as user: %s, %v", user, err)
	}

	if err := s.cfg.SetUser(dbUser.Name); err != nil {
		return err
	}

	fmt.Printf("user has been set to: %s", dbUser.Name)
	return nil
}

func printUser(user database.User) {
	fmt.Printf("User Id: %s\n", user.ID)
	fmt.Printf("User Name: %s\n", user.Name)
}

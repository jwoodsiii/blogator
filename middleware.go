package main

import (
	"context"
	"fmt"

	"github.com/jwoodsiii/blogator/internal/database"
)

func middlewareLoggedIn(handler func(*state, command, database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		// 1. Get the user here using s.db and s.cfg.CurrentUserName
		//    (the same code you used to have at the top of addfeed/follow/following)
		ctx := context.Background()
		user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Error getting user from database: %v", err)
		}

		// 2. Call the handler and pass that user in
		//    handler(s, cmd, user)
		// 3. Return whatever the handler returns
		return handler(s, cmd, user)
	}
}

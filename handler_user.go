package main
import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	user := cmd.Args[0]
	if err := s.cfg.SetUser(user); err != nil {
		return err
	}
	fmt.Printf("user has been set to: %s", user)
	return nil
}

package main

import (
	"fmt"
	"github.com/jwoodsiii/blogator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}
	cfg.SetUser()
	fmt.Printf("Current config: %v\n", cfg)
}

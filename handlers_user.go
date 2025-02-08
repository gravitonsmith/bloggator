package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gravitonsmith/bloggator/internal/database"
)

func usersHandler(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Too many arguments for reset")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		name := user.Name
		sp := fmt.Sprintf("* %s", name)
		if name == s.config.CurrentUser {
			sp += " (current)"
		}
		fmt.Println(sp)
	}
	return nil
}

func registerHandler(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No arguments were found for register")
	}
	name := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == nil {
		return err
	}

	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	user, err := s.db.CreateUser(context.Background(), args)
	if err != nil {
		return err
	}
	log.Printf("New user data: %v\n", user)

	if err := s.config.SetUser(user.Name); err != nil {
		return err
	}

	return nil
}

func loginHandler(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No arguments were found for login")
	}
	name := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return err
	}
	if err := s.config.SetUser(name); err != nil {
		return err
	}
	fmt.Printf("User has been set to %s\n", cmd.args[0])
	return nil
}

func resetHandler(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Too many arguments for reset")
	}
	err := s.db.DeleteUser(context.Background())
	if err != nil {
		return err
	}

	return nil
}

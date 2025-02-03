package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gravitonsmith/bloggator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

func loginHandler(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No arguments were found for login")
	}
	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return err
	}
	fmt.Printf("User has been set to %s\n", cmd.args[0])
	return nil
}

type commands struct {
	registered map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registered[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	fun, exists := c.registered[cmd.name]
	if !exists {
		return fmt.Errorf("Function to run not found")
	}
	if err := fun(s, cmd); err != nil {
		return err
	}
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error with reading json file: %v", err)
	}

	s := &state{config: &cfg}
	cmds := commands{make(map[string]func(*state, command) error)}

	cmds.register("login", loginHandler)

	args := os.Args
	if len(args) < 2 {
		log.Fatalln("Please provide a command to run")
	}

	cmd := command{name: args[1], args: args[2:]}
	if err := cmds.run(s, cmd); err != nil {
		log.Fatalf("Error occured while running command: %v", err)
	}
}

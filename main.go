package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gravitonsmith/bloggator/internal/config"
	"github.com/gravitonsmith/bloggator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

type command struct {
	name string
	args []string
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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUser)
		if err != nil {
			return err
		}
		if err := handler(s, cmd, user); err != nil {
			return err
		}
		return nil
	}
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error with reading json file: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("Error opening db connection: %v\n", err)
	}

	dbQueries := database.New(db)

	s := &state{
		db:     dbQueries,
		config: &cfg,
	}
	cmds := commands{make(map[string]func(*state, command) error)}

	cmds.register("login", loginHandler)
	cmds.register("register", registerHandler)
	cmds.register("reset", resetHandler)
	cmds.register("users", usersHandler)
	cmds.register("agg", aggHandler)
	cmds.register("addfeed", middlewareLoggedIn(addFeedHandler))
	cmds.register("feeds", middlewareLoggedIn(getAllFeedsHandler))
	cmds.register("follow", middlewareLoggedIn(followFeedHandler))
	cmds.register("following", middlewareLoggedIn(currentUserFeeds))
	cmds.register("unfollow", middlewareLoggedIn(deleteFeedFollow))
	cmds.register("browse", middlewareLoggedIn(browseHandler))

	args := os.Args
	if len(args) < 2 {
		log.Fatalln("Please provide a command to run")
	}

	cmd := command{name: args[1], args: args[2:]}
	if err := cmds.run(s, cmd); err != nil {
		log.Fatalf("Error occured while running command: %v", err)
	}
}

package main

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	user, err := s.db.GetUser(context.Background(), s.config.CurrentUser)
	if err != nil {
		return err
	}
	return handler
}

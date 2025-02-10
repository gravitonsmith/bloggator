# bloggator
This is a project to scrape RSS feeds from urls provide by users so that details can be easily viewed about those feeds by users.

## Dependencies
Postgres will need to be installed to handle the storage of all the different feeds and users
Go will need to be installed to run build and then run the application

## Installation
Once you have go installed, use git clone on this repo and then navigate to it and run `go install` to ensure all necessary packages are installed to run the program.

Manually create a config file in your home directory, ~/.gatorconfig.json, with the following content:
`{
  "db_url": "connection_string_goes_here",
  "current_user_name": "username_goes_here"
}`

Replace the values with your credentials and the username can be whatever you would like to start.

### Commands
- `login <username>` - sets the current user in the config
- `register <username>` - adds a new user to the database
- `users` - lists all the users in the database
- `addfeed <name> <url>` - adds a url to be scraped from with a name
- `agg <time-duration>` - start the aggregator that wil ping on the interval set
- `browse [limit]` - optional limit woth a default of 2

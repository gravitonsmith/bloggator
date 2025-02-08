package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/gravitonsmith/bloggator/internal/database"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func currentUserFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Following command does not have the correct number of args")
	}

	user, err := s.db.GetUser(context.Background(), s.config.CurrentUser)
	if err != nil {
		return err
	}

	feeds, err := s.db.GetFeedsByUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("Break it down for me:\n")
	fmt.Printf("User: %s\n", feeds[0].UserName)
	for _, feed := range feeds {
		fmt.Printf("Feed name: %s\n", feed.FeedName)
	}

	return nil
}

func followFeedHandler(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Follow command does not have the correct number of args")
	}

	user, err := s.db.GetUser(context.Background(), s.config.CurrentUser)
	if err != nil {
		return err
	}
	feed, err := s.db.GetFeedbyUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	args := database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), args)
	if err != nil {
		return err
	}

	fmt.Printf("Follow created!\nFeed name: %s\n User name: %s\n", follow.FeedName, follow.UserName)
	return nil
}

func getAllFeedsHandler(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Too many arguments added to feeds command")
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Here are all stored feeds")
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Feed name: %s\nURL: %s\nUser name: %s\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

func addFeedHandler(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Add command does not have the correct number of args")
	}
	name := cmd.args[0]
	url := cmd.args[1]
	currentUser := s.config.CurrentUser

	user, err := s.db.GetUser(context.Background(), currentUser)
	if err != nil {
		return err
	}
	params := database.CreateFeedParams{
		Name:   name,
		Url:    url,
		UserID: user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}

	args := database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	follow, err := s.db.CreateFeedFollow(context.Background(), args)
	fmt.Printf("Feed created and followed: %s", follow.FeedName)
	return nil
}

func aggHandler(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	feedData, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	fmt.Println(feedData)
	return nil
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	feed := &RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return feed, err
	}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return feed, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err := xml.Unmarshal(body, feed); err != nil {
		return feed, err
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return feed, nil
}

package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gravitonsmith/bloggator/internal/database"
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

func deleteFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Unfollow command does not have the correct number of args")
	}

	feed, err := s.db.GetFeedbyUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	args := database.DeleteFollowByUserParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	if err = s.db.DeleteFollowByUser(context.Background(), args); err != nil {
		return err
	}
	return nil
}

func currentUserFeeds(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Following command does not have the correct number of args")
	}

	feeds, err := s.db.GetFeedsByUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("Break it down for me:\n")
	fmt.Printf("User: %s\n", user.Name)
	for _, feed := range feeds {
		fmt.Printf("Feed name: %s\n", feed.FeedName)
	}

	return nil
}

func followFeedHandler(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Follow command does not have the correct number of args")
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

func getAllFeedsHandler(s *state, cmd command, user database.User) error {
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

func addFeedHandler(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Add command does not have the correct number of args")
	}

	name := cmd.args[0]
	url := cmd.args[1]

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
	if len(cmd.args) != 1 {
		return fmt.Errorf("Agg command needs a duration to run")
	}

	timeBetweenRefresh, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	log.Printf("Collecting feeds every %s", timeBetweenRefresh)

	ticker := time.NewTicker(timeBetweenRefresh)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeed(context.Background())
	if err != nil {
		log.Println("No feed to fetch: ", err)
		return
	}
	scrapeFeed(s, feed)
}

func scrapeFeed(s *state, feed database.Feed) {
	_, err := s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Could not mark feed as fetched: ", err)
		return
	}
	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Println("Could not get feed data: ", err)
		return
	}

	for _, item := range feedData.Channel.Item {
		fmt.Printf("Feed Item Title: %s\n", item.Title)
	}
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

package actions

import (
	"fmt"
	"io"
	"time"

	"github.com/bscott/golangflow/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gorilla/feeds"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	posts := &models.Posts{}
	q := tx.PaginateFromParams(c.Request().URL.Query())
	// You can order your list here. Just change
	err := q.Order("created_at desc").All(posts)
	// to:
	// err := tx.Order("create_at desc").All(posts)
	if err != nil {
		return errors.WithStack(err)
	}
	// getAvatar helper retrieves the post user's avatar link

	// Make posts available inside the html template
	c.Set("posts", posts)
	c.Set("pagination", q.Paginator)
	return c.Render(200, r.HTML("index.html"))
}

func RSSFeed(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	posts := models.Posts{}
	err := tx.Order("created_at desc").All(&posts)
	if err != nil {
		return errors.WithStack(err)
	}

	feed := feeds.Feed{
		Title:       "Golang Flow",
		Link:        &feeds.Link{Href: App().Host},
		Description: "All the Go news that's fit to print!",
		Author:      &feeds.Author{Name: "Brian Scott"},
		Created:     time.Now(),
		Copyright:   "This work is copyright © Brian Scott",
		Items:       make([]*feeds.Item, len(posts), len(posts)),
	}

	for i, p := range posts {
		u := &models.User{}
		err := tx.Find(u, p.UserID)
		if err != nil {
			return errors.WithStack(err)
		}
		feed.Items[i] = &feeds.Item{
			Title:       p.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("/posts/%s", p.ID)},
			Description: p.Content,
			Author:      &feeds.Author{Name: u.Name},
			Created:     p.CreatedAt,
		}
	}

	return c.Render(200, r.Func("application/rss+xml", func(w io.Writer, d render.Data) error {
		s, err := feed.ToRss()
		if err != nil {
			return errors.WithStack(err)
		}
		w.Write([]byte(s))
		return nil
	}))
}

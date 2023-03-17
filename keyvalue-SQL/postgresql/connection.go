package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Client struct {
	User       string
	Pass       string
	Host       string
	DB         string
	Connection *pgx.Conn
}

func (c *Client) Init(user, pass, db, host string) error {
	c.User = user
	c.Pass = pass
	c.Host = host
	c.DB = db
	var err error
	url := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", user, pass, host, 5432, db)
	fmt.Println(url)
	c.Connection, err = pgx.Connect(context.Background(), url)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	fmt.Println(c.Connection.Config().Database)
	return nil
}

func (c *Client) Test() error {
	greeting := ""
	err := c.Connection.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		return err
	}

	fmt.Println(greeting)
	return nil
}

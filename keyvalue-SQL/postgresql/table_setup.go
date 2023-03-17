package postgresql

import "context"

func (c *Client) CreateTable() error {
	_, err := c.Connection.Exec(context.Background(), "create table kvstore (key varchar(128) primary key not null,value text, ttl bigserial)")
	return err
}

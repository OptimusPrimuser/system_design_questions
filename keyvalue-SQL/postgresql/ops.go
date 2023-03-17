package postgresql

import (
	"context"
	"fmt"
	"time"
)

type RetVal struct {
	Key          string        `json:"key"`
	Value        string        `json:"value"`
	TtlRemaining time.Duration `json:"ttl remaining"`
}

type TTL struct {
	TTL    time.Duration
	remove bool
}

func (c *Client) PutKey(key, val string, ttl TTL) error {
	timeTill := "null"
	if ttl.remove {
		timeTill = "-1"
	} else if (ttl != TTL{}) {
		timeTill = fmt.Sprintf("%v", time.Now().Add(ttl.TTL).Unix())
	}
	query := fmt.Sprintf(
		"INSERT INTO kvstore VALUES('%s','%s',%s) ON CONFLICT (key) DO UPDATE Set value=EXCLUDED.value,ttl=EXCLUDED.ttl;",
		key, val, timeTill,
	)
	fmt.Println(query)
	_, err := c.Connection.Exec(context.Background(), query)
	return err
}

func (c *Client) GetKey(key string) (RetVal, error) {
	retVal := RetVal{}
	var ttl int64
	query := fmt.Sprintf("select * from kvstore where key='%s' and ttl != -1", key)
	err := c.Connection.QueryRow(context.Background(), query).Scan(&key, &retVal.Value, &ttl)
	if err != nil {
		return RetVal{}, err
	}
	endTime := time.Unix(ttl, 0)
	if endTime.Unix() != 0 {
		retVal.TtlRemaining = time.Until(endTime)
	}
	retVal.Key = key
	return retVal, nil
}

func (c *Client) RemoveKey(key string) error {
	retval, err := c.GetKey(key)
	if err != nil {
		return err
	}
	err = c.PutKey(retval.Key, retval.Value, TTL{remove: true})
	return err
}

func (c *Client) BatchDelete() error {
	query := "DELETE FROM kvstore WHERE  ttl < 0;"
	_, err := c.Connection.Exec(context.Background(), query)
	return err
}

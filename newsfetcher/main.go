package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
)

func setNews(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Newsfetcher: SENDING REQUEST FOR NEWS...")

	// newPool returns a pointer to a redis.Pool
	pool := newPool()
	// get a connection from the pool (redis.Conn)
	conn := pool.Get()
	// use defer to close the connection when the function completes
	defer conn.Close()

	// call Redis PING command to test connectivity
	err := ping(conn)
	if err != nil {
		fmt.Println(err)
	}

	err = set(conn)
	if err != nil {
		fmt.Println(err)
	}
}

func set(c redis.Conn) error {
	_, err := c.Do("SET", "news", "No news! Is good news ;-)")
	if err != nil {
		return err
	}
	return nil
}

func main() {

	http.HandleFunc("/", setNews)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// ping tests connectivity for redis (PONG should be returned)
func ping(c redis.Conn) error {
	// Send PING command to Redis
	// PING command returns a Redis "Simple String"
	// Use redis.String to convert the interface type to string
	s, err := redis.String(c.Do("PING"))
	if err != nil {
		return err
	}

	fmt.Printf("PING Response = %s\n", s)
	// Output: PONG

	return nil
}

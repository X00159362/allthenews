/* custom file for eades microservices api first lab */
package swagger

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
)

func GetAllNews(w http.ResponseWriter, r *http.Request) {
	urls := [2]string{"http://127.0.0.1:8888", "http://127.0.0.1:9999"}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// PROCESSING STAGE 1
	// Get information from news and weather services
	infoTypes := [2]string{"news", "weather"}
	var fetchedStrings [2]string = [2]string{"", ""}
	for i := 0; i < 2; i++ {
		resp, err := http.Get(urls[i])
		if err != nil {
			fmt.Fprintln(w, "allthenews[ERROR]: Couldn't get "+infoTypes[i]+" from site. "+err.Error()+"<br/>")
			log.Printf("allthenews[ERROR]: Couldn't get " + infoTypes[i] + " from site. " + err.Error())
		} else {
			if resp.StatusCode == http.StatusOK {
				// bodyBytes, err2 := ioutil.ReadAll(resp.Body)
				// if err2 != nil {
				// 	fmt.Fprintln(w, "allthenews[ERROR]: Couldn't get "+infoTypes[i]+" from response."+err2.Error()+"<br/>")
				// 	log.Printf("allthenews[ERROR]: Couldn't get " + infoTypes[i] + " from response." + err2.Error())
				// } else {

				//fetchedStrings[i] = string(bodyBytes)

				log.Printf("allthenews[INFO]: getting news from redis")

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

				time.Sleep(2 * time.Second)

				key := "news"
				s, err := redis.String(conn.Do("GET", key))
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("%s = %s\n", key, s)
				fetchedStrings[0] = s

				time.Sleep(2 * time.Second)

				key = "weather"
				s, err = redis.String(conn.Do("GET", key))
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("%s = %s\n", key, s)
				fetchedStrings[1] = s
				//}
			} else {
				fmt.Fprintln(w, "allthenews[ERROR]: HTTP returned status "+string(resp.StatusCode)+"<br/>")
				log.Printf("allthenews[ERROR]: HTTP returned status " + string(resp.StatusCode))
			}
		}
	}

	// Create the inserts for the HTML file, which is only a skeleton with no information.
	inserts := struct {
		News    string
		Weather string
	}{
		fetchedStrings[0],
		fetchedStrings[1],
	}

	// PROCESSING STAGE 2
	// Read the query parameter "style" and check it against the different allowed values, assigning
	// the appropriate template name to variable templateName.
	var templateName = ""
	switch r.URL.Query().Get("style") {
	case "plain":
		templateName = "plain.html"
	case "colourful":
		templateName = "colour.html"
	case "blackandwhite":
		templateName = "bandw.html"
	}

	// PROCESSING STAGE 3
	// We are using the template handling library html/template to insert the information fetches from the other
	// services into the page with the requested style (via parameter 'style').
	if templateName != "" {
		// Now put together an HTML page. The template.ParseFiles() function inserts the values from structure
		// 'inserts' into the chosen template.
		t, _ := template.ParseFiles(templateName)
		t.Execute(w, inserts)
	} else {
		fmt.Fprintln(w, "allthenews[ERROR]: Invalid style parameter.<br>")
		log.Printf("allthenews[ERROR]: Invalid style parameter.")
	}
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

// get executes the redis LPOP command
func get(c redis.Conn) error {

	key := "news"
	s, err := redis.String(c.Do("LPOP", key))
	if err != nil {
		return (err)
	}
	fmt.Printf("%s = %s\n", key, s)

	return nil
}

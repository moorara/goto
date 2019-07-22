package main

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	xhttp "github.com/moorara/goto/http"
	"github.com/moorara/goto/log"
)

type client struct {
	logger *log.Logger
	mid    *xhttp.ClientMiddleware
}

func (c *client) call() {
	// A random delay between 1s to 5s
	d := 1 + rand.Intn(4)
	time.Sleep(time.Duration(d) * time.Second)

	// Create an http client
	client := http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{},
	}

	// Create an http request
	req, _ := http.NewRequest("GET", serverAddress+"/", nil)

	// Make the request to http server
	doer := c.mid.Metrics(c.mid.RequestID(c.mid.Tracing(c.mid.Logging(client.Do))))
	res, err := doer(req)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	c.logger.Info("message", string(b))
}

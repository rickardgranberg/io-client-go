package main

/*
A sanity-check sketch for inspecting the raw request before it's sent to
Adafruit IO. To use, change the API calls inside the `CallAPI` method to the
one / ones you want to see.

No requests generated by this sketch will be sent to io.adafruit.com, so it's
safe to use a bogus secret key.

For example:

    $ go run examples/debug/request_viewer.go -key "12345ABC"
    2016/05/26 09:10:07 -- received request --
    ---
    POST /api/v2/feeds/beta-test/data HTTP/1.1
    Host: 127.0.0.1:53626
    Accept: application/json
    Accept-Encoding: gzip
    Content-Type: application/json
    User-Agent: Adafruit IO Go Client v0.1
    X-Aio-Key: 12345ABC

    {"value":22}
    ---

*/

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"

	"github.com/adafruit/io-client-go/v2/pkg/adafruitio"
)

// Add the API call you want to examine here to see it output at the command line.
func CallAPI(client *adafruitio.Client) {
	client.SetFeed(&adafruitio.Feed{Key: "beta-test"})
	client.Data.Send(&adafruitio.Data{Value: "22"})
}

func main() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		log.Println("received request")
		fmt.Println("---")
		fmt.Printf("%v", string(dump))
		fmt.Println("---")
		fmt.Fprint(w, "{}")
	}))
	defer ts.Close()

	var key string
	var username string
	flag.StringVar(&key, "key", "", "your Adafruit IO key")
	flag.StringVar(&username, "user", "", "your Adafruit IO user name")
	flag.Parse()

	if key == "" {
		key = os.Getenv("ADAFRUIT_IO_KEY")
	}

	if username == "" {
		username = os.Getenv("ADAFRUIT_IO_USERNAME")
	}

	client := adafruitio.NewClient(username, key)

	if ts.URL != "" {
		client.SetBaseURL(ts.URL)
	}

	CallAPI(client)
}

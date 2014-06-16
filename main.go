// © 2014 Steve McCoy under the MIT license. See LICENSE for details.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
)

func main() {
	u, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Who are you? %v\n", err)
		os.Exit(1)
	}

	keyloc := flag.String("k", u.HomeDir+"/lib/weather_key", "weather key file")
	flag.Parse()

	k, err := os.Open(*keyloc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open %q: %v\n", *keyloc, err)
		os.Exit(1)
	}

	key, err := ioutil.ReadAll(k)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read %q: %v\n", *keyloc, err)
		os.Exit(1)
	}

	resp, err := http.Get("http://api.wunderground.com/api/" + string(key) + "/conditions/q/03801.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to query wunderground: %v\n", err)
		os.Exit(1)
	}

	j := json.NewDecoder(resp.Body)
	var weather conds
	err = j.Decode(&weather)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse response: %v\n", err)
		os.Exit(1)
	}

	_, err = fmt.Println(weather)
}

type conds struct {
	Current current `json:"current_observation"`
}

type current struct {
	Weather  string  `json:"weather"`
	TempF    float64 `json:"temp_f"`
	TempC    float64 `json:"temp_c"`
	Humidity string  `json:"relative_humidity"`
}

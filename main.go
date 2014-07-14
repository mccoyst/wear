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

type layer int

const (
	Primary   layer = 1 + iota
	Secondary       = iota << 1
	Top = iota << 1
)

type clothing struct {
	name string
	layer
}

func (c clothing) String() string {
	return c.name
}

var (
	shirt          = clothing{"shirt", Primary|Top}
	tshirt         = clothing{"t-shirt", Primary | Secondary|Top}
	longundershirt = clothing{"long undershirt", Secondary|Top}
	hoodie         = clothing{"hoodie", Secondary|Top}
	jacket         = clothing{"jacket", Secondary|Top}
	coat           = clothing{"coat", Secondary|Top}
	trousers       = clothing{"trousers", Primary}
	shorts         = clothing{"shorts", Primary}
	leggings       = clothing{"leggings", Secondary}
)

func possibilities(t float64) []clothing {
	var c []clothing
	if t < 30 {
		c = append(c, longundershirt, leggings)
	}
	if t < 40 {
		c = append(c, coat)
	}
	if t < 50 {
		c = append(c, hoodie)
	}
	if t < 60 {
		c = append(c, jacket)
	}
	if t < 70 {
		c = append(c, shirt)
	}
	if t >= 70 {
		c = append(c, tshirt)
	}
	if t >= 85 {
		c = append(c, shorts)
	} else {
		c = append(c, trousers)
	}

	return c
}

func clothes(t float64) [][]clothing {
	return [][]clothing{possibilities(t)}
}

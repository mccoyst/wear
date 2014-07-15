// © 2014 Steve McCoy under the MIT license. See LICENSE for details.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/user"
	"sort"
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
	Top             = iota << 1
)

type clothing struct {
	name  string
	score int
}

func (c clothing) String() string {
	return fmt.Sprintf("%s(%d)", c.name, c.score)
}

type clothes []clothing

func (c clothes) Len() int {
	return len(c)
}

func (c clothes) Less(i, j int) bool {
	return c[i].name < c[j].name
}

func (c clothes) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c clothes) Rand() clothing {
	return c[rand.Intn(len(c))]
}

var primaryTops = clothes{
	clothing{"shirt", 10},
	clothing{"t-shirt", 5},
}

var secondaryTops = clothes{
	clothing{"t-shirt", 5},
	clothing{"long undershirt", 10},
	clothing{"hoodie", 5},
	clothing{"jacket", 5},
	clothing{"coat", 10},
}

var primaryBottoms = clothes{
	clothing{"trousers", 10},
	clothing{"shorts", 5},
}

var secondaryBottoms = clothes{
	clothing{"leggings", 10},
}

func wear(t float64) [][]clothing {
	var combos [][]clothing

	tier := int(math.Ceil(t))
	goal := 35 - tier
	fmt.Fprintln(os.Stderr, "wearing", int(t), goal)

	sort.Sort(primaryTops)
	sort.Sort(secondaryTops)
	sort.Sort(primaryBottoms)
	sort.Sort(secondaryBottoms)

	//	bot0 := primaryBottoms.Rand()
	fmt.Fprintln(os.Stderr, "primary tops =", len(primaryTops))
	for _, top0 := range primaryTops {
		fmt.Fprintln(os.Stderr, "-------------------------")
		for {
			k := len(secondaryTops) - 2
			for ; k >= 0; k-- {
				if secondaryTops[k].name < secondaryTops[k+1].name {
					break
				}
			}

			if k < 0 || secondaryTops[k].name >= secondaryTops[k+1].name {
				break
			}

			i := len(secondaryTops) - 1
			for ; i > k; i-- {
				if secondaryTops[k].name < secondaryTops[i].name {
					break
				}
			}

			if i == k {
				break
			}

			secondaryTops.Swap(k, i)
			sort.Sort(sort.Reverse(secondaryTops[k+1:]))

			score := top0.score
			tops := make([]clothing, 0, len(secondaryTops)+1)
			tops = append(tops, top0)
			for _, c := range secondaryTops {
				if c.score+score > goal {
					break
				}
				tops = append(tops, c)
				score += c.score
			}

			fmt.Fprintln(os.Stderr, tops)
			combos = append(combos, tops)
		}
		sort.Sort(secondaryTops)
	}

	return combos
}

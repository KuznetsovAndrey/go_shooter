package shooter

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)

// Victim for shooter
type Victim struct {
	Headers []string
	Method  string
	Host    string
	Port    string
	Scheme  string
}

// Gun for shooter
type Gun struct {
	Shots    int
	Delay    int
	Parallel int
}

// Shooter structure
type Shooter struct {
	codes  map[int]int
	times  []float64
	victim Victim
	gun    Gun
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func (g *Gun) bang(done chan *http.Response, err chan error, client http.Client, req *http.Request) {
	resp, requestErr := client.Do(req)
	if requestErr != nil {
		err <- requestErr
	}
	done <- resp
}

// Shoot the victim
func (sh *Shooter) Shoot() {
	client := &http.Client{}

	req, err := http.NewRequest(sh.victim.Method, fmt.Sprintf("%s://%s%s", sh.victim.Scheme, sh.victim.Host, sh.victim.Port), nil)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(sh.victim.Headers); i++ {
		if sh.victim.Headers[i] == "" {
			continue
		}

		header := strings.Split(sh.victim.Headers[i], ":")
		header[0] = strings.TrimSpace(header[0])
		req.Header.Add(header[0], header[1])
	}

	go spinner(100 * time.Millisecond)

	counts := math.Ceil(float64(sh.gun.Shots) / float64(sh.gun.Parallel))

	for i := 0; i < int(counts); i++ {
		start := time.Now()

		done := make(chan *http.Response, sh.gun.Parallel)
		err := make(chan error)

		defer close(done)
		defer close(err)

		qCount := sh.gun.Parallel
		if (i+1)*sh.gun.Parallel > sh.gun.Shots {
			qCount = sh.gun.Parallel - ((i+1)*sh.gun.Parallel - sh.gun.Shots)
		}

		for i := 0; i < qCount; i++ {
			go sh.gun.bang(done, err, *client, req)
		}

		for i := 0; i < qCount; i++ {
			select {
			case requestErr := <-err:
				fmt.Println(requestErr)
			case resp := <-done:
				defer resp.Body.Close()
				t := time.Now()
				elapsed := t.Sub(start) / time.Millisecond

				sh.times = append(sh.times, float64(elapsed))
				sh.codes[resp.StatusCode]++
			}
		}
		time.Sleep(time.Duration(sh.gun.Delay) * time.Millisecond)
	}
}

// Report - answer shooter about his shooting
func (sh *Shooter) Report() string {
	max := sh.times[0]
	min := sh.times[0]
	sum := float64(0)

	for i := 0; i < len(sh.times); i++ {
		if max < sh.times[i] {
			max = sh.times[i]
		}
		if min > sh.times[i] {
			min = sh.times[i]
		}

		sum += sh.times[i]
	}

	resString := fmt.Sprintf("\rTotal - %dms, max - %dms, min - %dms", int(sum), int(max), int(min))
	resString += "\nCodes:"

	for key, value := range sh.codes {
		resString += fmt.Sprintf("\n%d: %d", key, value)
	}

	return resString
}

// HireShooter - hire your own shooter for victim and give him a gun
func HireShooter(gun Gun, victim Victim) *Shooter {
	shooter := &Shooter{
		gun:    gun,
		victim: victim,
		codes:  make(map[int]int),
		times:  []float64{},
	}

	return shooter
}

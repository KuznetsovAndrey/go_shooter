package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"./shooter"
)

func main() {
	host := flag.String("h", "", "host")
	port := flag.Int("p", -1, "port")
	method := flag.String("X", "GET", "method")
	headers := flag.String("H", "", "headers, separated by semicolon")
	shots := flag.Int("n", 1, "count of shots")
	delay := flag.Int("d", 0, "delay in ms")
	parallel := flag.Int("c", 1, "parallel shoots in time")
	https := flag.Int("ssl", 0, "use https schema")

	flag.Parse()

	victimPort := ""
	if *port > -1 {
		victimPort = ":" + strconv.Itoa(*port)
	}

	if *parallel <= 0 {
		*parallel = 1
	}

	scheme := "http"
	if *https == 1 {
		scheme = "https"
	}

	victim := shooter.Victim{
		Headers: strings.Split(*headers, ";"),
		Port:    victimPort,
		Method:  *method,
		Host:    *host,
		Scheme:  scheme,
	}

	gun := shooter.Gun{
		Shots:    *shots,
		Delay:    *delay,
		Parallel: *parallel,
	}

	hiredShooter := shooter.HireShooter(gun, victim)

	hiredShooter.Shoot()
	fmt.Println(hiredShooter.Report())
}

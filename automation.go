package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	challenge_number := get_challenge_number()
	println(challenge_number)
}

// This function finds the latest challenge number by scraping the tester's github page
func get_challenge_number() int {
	const URL = "https://github.com/Pomroka/TWT_Challenges_Tester/"
	resp, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// looking for the latest challenge number from the page
	body_array := strings.Split(string(body), "\n")
	challenge_number := 0

	for _, line := range body_array {
		if strings.Contains(line, "Challenge_") {
			index := strings.Index(line, "Challenge_")
			curr_challenge_number, err := strconv.Atoi(line[index+len("Challenge_") : index+len("Challenge_")+3])
			if err != nil {
				continue
			}
			if challenge_number < curr_challenge_number {
				challenge_number = curr_challenge_number
			}
		}
	}

	return challenge_number
}

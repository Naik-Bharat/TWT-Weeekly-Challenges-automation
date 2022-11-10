package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	challenge_number := get_challenge_number()
	println(challenge_number)
	wget("https://raw.githubusercontent.com/Naik-Bharat/TWT-Weekly-Challenges/master/Challenge_100/Solution.py", "test.py")
}

// This function finds the latest challenge number by scraping the tester's github page
func get_challenge_number() int {
	const URL = "https://github.com/Pomroka/TWT_Challenges_Tester/"
	resp, err := http.Get(URL)
	if err != nil {
		println("Error downloading tester's webpage!")
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println("Error reading content from tester's webpage!")
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

	println("Challenge Number: {}", challenge_number)
	return challenge_number
}

func wget(link string, file_path string) {
	cmd := exec.Command("wget", link, "-O", file_path)
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	println(string(stdout))
}

func move_directory(old_path string, new_path string) {
	err := os.Rename(old_path, new_path)
	if err != nil {
		println("Error moving {}", old_path)
		panic(err)
	}
}

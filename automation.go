package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type git_page struct {
	link string
}

// this function returns a slice of files in a git page
func (page git_page) get_git_files() []string {
	// list of all files
	files := []string{}

	// content in the page as a slice
	body_array := get_url(page.link)
	for index, line := range body_array {
		// if it is a file
		if strings.Contains(line, "rowheader") {
			words_array := strings.Split(body_array[index+1], " ")
			for _, word := range words_array {
				// if this word contains the title
				if strings.Contains(word, "title") {
					file := word[len("title")+2 : len(word)-1]
					files = append(files, file)
				}
			}
		}
	}

	return files
}

func main() {
	challenge_number := get_challenge_number()
	fmt.Printf("Challenge Number: %v\n", challenge_number)
	new_directory(fmt.Sprintf("Challenge_%v", challenge_number))
}

// This function finds the latest challenge number by scraping the tester's github page
func get_challenge_number() int {
	var challenge_number int

	page := git_page{"https://github.com/Pomroka/TWT_Challenges_Tester/"}
	files := page.get_git_files()
	for _, file := range files {
		if strings.Contains(file, "Challenge_") {
			curr_challenge_number, err := strconv.Atoi(file[len("Challenge_"):])
			if err != nil {
				println("Error converting challenge number to integer")
				panic(err)
			}
			if challenge_number < curr_challenge_number {
				challenge_number = curr_challenge_number
			}
		}
	}

	return challenge_number
}

func get_url(link string) []string {
	resp, err := http.Get(link)
	if err != nil {
		println("Error downloading", link)
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println("Error reading content from", link)
		panic(err)
	}
	return strings.Split(string(body), "\n")
}

// function to wget files
func wget(link string, file_path string) {
	cmd := exec.Command("wget", link, "-O", file_path)
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}
}

// function to handle moving directories
func move_directory(old_path string, new_path string) {
	err := os.Rename(old_path, new_path)
	if err != nil {
		println("Error moving {}", old_path)
		panic(err)
	}
}

func new_directory(path string) {
	err := os.Mkdir(path, 0750)
	if err == nil {
		fmt.Printf("New Directory %v created\n", path)
	}
}

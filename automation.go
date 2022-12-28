package main

import (
	"fmt"
	"io"
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
	folder_restructure(challenge_number)
	download_files(challenge_number)
}

func download_files(challenge_number int) {
	page := git_page{fmt.Sprintf("https://github.com/Pomroka/TWT_Challenges_Tester/tree/main/Challenge_%v", challenge_number)}
	files := page.get_git_files()
	for _, file := range files {
		if file != "G" {
			wget(fmt.Sprintf("https://raw.githubusercontent.com/Pomroka/TWT_Challenges_Tester/main/Challenge_%v/%v", challenge_number, file), fmt.Sprintf("Challenge_%v/%v", challenge_number, file))
		}
	}
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

	fmt.Printf("Challenge Number: %v\n", challenge_number)
	return challenge_number
}

// function to read a url and return a slice of its content
func get_url(link string) []string {
	resp, err := http.Get(link)
	if err != nil {
		println("Error downloading", link)
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		println("Error reading content from", link)
		panic(err)
	}
	return strings.Split(string(body), "\n")
}

// function to wget files
func wget(link string, file_path string) {
	_, err := os.OpenFile(file_path, os.O_RDONLY, 0750)
	if err == nil {
		var choice string
		fmt.Printf("%v already exists. Do you want to rewrite? (y/N) ", link)
		fmt.Scanln(&choice)
		if choice != "y" && choice != "Y" {
			return
		}
	}
	cmd := exec.Command("wget", link, "-O", file_path)
	_, err = cmd.Output()
	if err != nil {
		panic(err)
	}
	println(fmt.Sprintf("Downloaded %v", link))
}

// this function returns whether a folder should be moved to previous challenges folder
func old_challenge(file string, challenge_number int) bool {
	file_challenge_number, err := strconv.Atoi(file[len("Challenge_"):])
	if err != nil {
		println("Cannot convert challenge number to int in challenge tester's git repo")
		panic(err)
	}

	challenge_diff := challenge_number - file_challenge_number
	return challenge_diff > 2
}

// function to handle folder re-structuring
func folder_restructure(challenge_number int) {
	new_directory(fmt.Sprintf("Challenge_%v", challenge_number))
	files, err := os.ReadDir("./")
	if err != nil {
		println("Cannot read current directory")
		panic(err)
	}
	for _, file := range files {
		if strings.Contains(file.Name(), "Challenge_") {
			if old_challenge(file.Name(), challenge_number) {
				move_directory(file.Name(), fmt.Sprintf("Previous_Challenges/%v", file.Name()))
				fmt.Printf("Moved %v to %v\n", file.Name(), fmt.Sprintf("Previous_Challenges/%v", file.Name()))
			}
		}
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

// function to create new directories
func new_directory(path string) {
	err := os.Mkdir(path, 0750)
	if err == nil {
		fmt.Printf("New Directory %v created\n", path)
	}
}

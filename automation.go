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

type gitPage struct {
	link string
}

// this function returns a slice of files in a git page
func (page gitPage) getGitFiles() []string {
	// list of all files
	files := []string{}

	// content in the page as a slice
	bodyArray := getUrl(page.link)
	for index, line := range bodyArray {
		// if it is a file
		if strings.Contains(line, "rowheader") {
			wordsArray := strings.Split(bodyArray[index+1], " ")
			for _, word := range wordsArray {
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
	challengeNumber := getChallengeNumber()
	folderRestructure(challengeNumber)
	downloadFiles(challengeNumber)
}

func downloadFiles(challengeNumber int) {
	page := gitPage{fmt.Sprintf("https://github.com/Pomroka/TWT_Challenges_Tester/tree/main/Challenge_%v", challengeNumber)}
	files := page.getGitFiles()
	for _, file := range files {
		if file != "G" {
			wget(fmt.Sprintf("https://raw.githubusercontent.com/Pomroka/TWT_Challenges_Tester/main/Challenge_%v/%v", challengeNumber, file), fmt.Sprintf("Challenge_%v/%v", challengeNumber, file))
		}
	}
}

// This function finds the latest challenge number by scraping the tester's github page
func getChallengeNumber() int {
	var challengeNumber int

	page := gitPage{"https://github.com/Pomroka/TWT_Challenges_Tester/"}
	files := page.getGitFiles()
	for _, file := range files {
		if strings.Contains(file, "Challenge_") {
			currChallengeNumber, err := strconv.Atoi(file[len("Challenge_"):])
			if err != nil {
				println("Error converting challenge number to integer")
				panic(err)
			}
			if challengeNumber < currChallengeNumber {
				challengeNumber = currChallengeNumber
			}
		}
	}

	fmt.Println("Challenge Number:", challengeNumber)
	return challengeNumber
}

// function to read a url and return a slice of its content
func getUrl(link string) []string {
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
func wget(link string, filePath string) {
	_, err := os.OpenFile(filePath, os.O_RDONLY, 0750)
	if err == nil {
		var choice string
		fmt.Printf("%v already exists. Do you want to rewrite? (y/N) ", link)
		fmt.Scanln(&choice)
		if choice != "y" && choice != "Y" {
			return
		}
	}
	cmd := exec.Command("wget", link, "-O", filePath)
	_, err = cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("Downloaded", link)
}

// this function returns whether a folder should be moved to previous challenges folder
func oldChallenge(file string, challengeNumber int) bool {
	fileChallengeNumber, err := strconv.Atoi(file[len("Challenge_"):])
	if err != nil {
		println("Cannot convert challenge number to int in challenge tester's git repo")
		panic(err)
	}

	challengeDiff := challengeNumber - fileChallengeNumber
	return challengeDiff > 2
}

// function to handle folder re-structuring
func folderRestructure(challengeNumber int) {
	newDirectory(fmt.Sprintf("Challenge_%v", challengeNumber))
	files, err := os.ReadDir("./")
	if err != nil {
		println("Cannot read current directory")
		panic(err)
	}
	for _, file := range files {
		if strings.Contains(file.Name(), "Challenge_") {
			if oldChallenge(file.Name(), challengeNumber) {
				newDirectory("Previous_Challenges")
				moveDirectory(file.Name(), fmt.Sprintf("Previous_Challenges/%v", file.Name()))
				fmt.Printf("Moved %v to %v\n", file.Name(), fmt.Sprintf("Previous_Challenges/%v", file.Name()))
			}
		}
	}
}

// function to handle moving directories
func moveDirectory(oldPath string, newPath string) {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		println("Error moving ", oldPath)
		panic(err)
	}
}

// function to create new directories
func newDirectory(path string) {
	err := os.Mkdir(path, 0750)
	if err == nil {
		fmt.Printf("New Directory %v created\n", path)
	}
}

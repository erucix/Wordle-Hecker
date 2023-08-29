package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Solution struct {
	ID              int    `json:"id"`
	Solution        string `json:"solution"`
	PrintDate       string `json:"print_date"`
	DaysSinceLaunch int    `json:"days_since_launch"`
	Editor          string `json:"editor"`
}

func main() {
	var passedArguments = os.Args
	var currentTime = time.Now()

	var userMentionedDate string
	var dateToBeSolved string

	if len(passedArguments) > 2 {
		for elementIndex, elementValue := range passedArguments {
			if elementValue == "--solve" {
				userMentionedDate = passedArguments[elementIndex+1]
			}
		}

		if userMentionedDate == "today" {
			dateToBeSolved = currentTime.Format("2006-01-02")
		} else if userMentionedDate == "yesterday" {
			dateToBeSolved = currentTime.AddDate(0, 0, -1).Format("2006-01-02")
		} else if userMentionedDate == "tomorrow" {
			dateToBeSolved = currentTime.AddDate(0, 0, 1).Format("2006-01-02")
		} else {
			dateToBeSolved = userMentionedDate
		}

		getHttpRequest(dateToBeSolved)

	} else {
		fmt.Printf(`
./wordle [OPTIONS] [DATE]

Available options:

--solve			Solves the wordle of given date

Available alias for date:

today			Assigns today's date
tomorrow		Assigns tomorrow's date
yesterday		Assigns yesterday's date
%v  		Solves the wordle of 2023-07-22
			You can specify your own here

Note: Date should be in this format: YYYY:MM:DD (Ex: 2023-07-22)

`, currentTime.Format("2006-01-02"))
	}
}

func getHttpRequest(dateToBeSolved string) {
	var solution Solution
	var requestURL string = "https://www.nytimes.com/svc/wordle/v2/" + dateToBeSolved + ".json"

	response, err := http.Get(requestURL)

	if err != nil || response.StatusCode != 200 {
		errorMessage := fmt.Sprintf("Server responded with status code %v", response.StatusCode)

		panic(errorMessage)
	}

	defer response.Body.Close()

	responseContent, _ := io.ReadAll((response.Body))

	if err := json.Unmarshal([]byte(string(responseContent)), &solution); err != nil {
		panic(err)
	}

	fmt.Println("\n\033[1;34mWordle ID\033[0m \033[1;32m:\033[0m ", solution.ID)
	fmt.Println("\033[1;34mAnswer\033[0m    \033[1;32m:\033[0m \033[1;32m", solution.Solution)
	fmt.Println("\033[1;34mDate\033[0m      \033[1;32m:\033[0m ", solution.PrintDate)
	fmt.Println("\033[1;34mEditor\033[0m    \033[1;32m:\033[0m ", solution.Editor)
}

package quizz

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type Problem struct {
	question string
	answer   string
}

func BuildQuiz() {
	csvFileName, timeLimit := getInput()

	fileContent := readCsv(csvFileName)

	problems := readProblems(fileContent)

	correct := 0
	fmt.Println("Time limit is ", timeLimit, " seconds")
	timer := time.NewTimer(time.Duration(timeLimit * int(time.Second)))

problemLoop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)
		ansCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			ansCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("!!! Time's up !!!\n")
			break problemLoop
		case answer := <-ansCh:
			if answer == p.answer {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func getInput() (string, int) {
	csvFileName := flag.String(
		"csv",
		"utils/problems.csv",
		"a csv file in the format of 'question,answer' ",
	)
	timeLimit := flag.Int(
		"limit",
		30,
		"the time limit for the quiz in seconds",
	)
	flag.Parse()
	return *csvFileName, *timeLimit
}

func readCsv(csvFileName string) [][]string {
	csvFile, err := os.Open(csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Error opening CSV file: %v", err))
	}

	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Error parsing CSV file: %v", err))
	}
	return lines
}

func readProblems(lines [][]string) []Problem {
	problems := make([]Problem, len(lines))
	for i, line := range lines {
		problems[i] = Problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return problems
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	var timeoutFlag int
	var helpFlag bool

	flag.BoolVar(&helpFlag, "h", false, "show usage")
	flag.IntVar(&timeoutFlag, "t", 30, "amount of time allowed for completion of the quiz (in seconds)")
	flag.Parse()

	if helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	fmt.Printf("You will have %d seconds to complete the quiz. Press Enter when you are ready to begin...", timeoutFlag)
	fmt.Scanln()

	answers := make(chan bool)
	completed := make(chan bool)
	timeout := time.After(time.Second * time.Duration(timeoutFlag))

	f, err := os.Open("problems.csv")
	if err != nil {
		handleError(err)
	}

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		handleError(err)
	}

	go presentQuestions(records, answers, completed)

	correct := 0

Abort:
	for {
		select {
		case answer := <-answers:
			if answer {
				correct++
			}
		case <-timeout:
			fmt.Printf("\nYou're out of time!\n")
			break Abort
		case <-completed:
			break Abort
		}
	}

	fmt.Printf("You answered %d out of %d questions correctly.\n", correct, len(records))
}

func handleError(err error) {
	fmt.Printf("Error: %s\n", err)
	os.Exit(1)
}

func presentQuestions(questions [][]string, answers chan bool, completed chan bool) {
	for _, r := range questions {
		presentQuestion(r, answers)
	}
	completed <- true
}

func presentQuestion(question []string, answers chan bool) {
	fmt.Printf("%s: ", question[0])
	var got string
	_, err := fmt.Scanln(&got)
	if err != nil {
		handleError(err)
	}
	isCorrect := got == question[1]
	answers <- isCorrect
}

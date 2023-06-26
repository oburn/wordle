package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"log"
)

//go:embed valid-words.txt
var validWordsAsBytes []byte

func loadWords() ([]string, error) {
	var words []string
	scanner := bufio.NewScanner(bytes.NewReader(validWordsAsBytes))
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return words, nil
}

func main() {
	var guess string
	var outcome string
	state := newState(5)
	words, err := loadWords()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		fmt.Print("  Enter guess: ")
		_, err := fmt.Scanln(&guess)
		if err != nil {
			log.Fatalf("Unable to read the guess: %v\n", err)
		}
		fmt.Print("Enter outcome: ")
		_, err = fmt.Scanln(&outcome)
		if err != nil {
			log.Fatalf("Unable to read the outcome: %v\n", err)
		}
		state = state.Next(guess, outcome)
		scored := state.ScoreWords(words)
		fmt.Println("---------------------------------------------------------------------------")
		fmt.Print(Tabulate(scored, 11, 10))
		fmt.Println("---------------------------------------------------------------------------")
	}
}

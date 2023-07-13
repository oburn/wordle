package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"
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
	words, err := loadWords()
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 1 {
		batch(words, os.Args[1:])
	} else {
		interactive(words)
	}
}

func batch(words []string, args []string) {
	if len(args)%2 != 0 {
		log.Fatalf("Expected even number of arguments, got %d\n", len(args))
	}
	num := len(args) / 2
	state := newState(5)
	for i := 0; i < num; i++ {
		guess := args[2*i]
		outcome := args[2*i+1]
		fmt.Printf("Guess: %s, Outcome: %s\n", guess, outcome)
		state = state.Next(guess, outcome)
	}
	scored := state.ScoreWords(words)
	fmt.Println("---------------------------------------------------------------------------")
	fmt.Print(Tabulate(scored, 11, 10))
	fmt.Println("---------------------------------------------------------------------------")
}

func interactive(words []string) {
	var guess string
	var outcome string
	state := newState(5)

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

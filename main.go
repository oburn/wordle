package main

import (
	"bufio"
	"fmt"
	"os"
)

func loadWords(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
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
	words, err := loadWords("valid-words.txt")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		fmt.Print("  Enter guess: ")
		fmt.Scanln(&guess)
		fmt.Print("Enter outcome: ")
		fmt.Scanln(&outcome)
		state = state.Next(guess, outcome)
		scored := state.ScoreWords(words)
		fmt.Println("---------------------------------------------------------------------------")
		fmt.Print(Tabulate(scored, 11, 10))
		fmt.Println("---------------------------------------------------------------------------")
	}
}

package main

import (
	"regexp"
	"sort"
	"strings"
)

type ScoredWord struct {
	word  string
	score int
}

type State struct {
	knownLetters string
	candidates   []string
}

func newState(size int) State {
	var result = State{knownLetters: "", candidates: []string{}}
	for i := 0; i < size; i++ {
		result.candidates = append(result.candidates, "abcdefghijklmnopqrstuvwxyz")
	}
	return result
}

func (s State) String() string {
	result := "{knownLetters: " + s.knownLetters + ", candidates: ["

	for i, c := range s.candidates {
		if i > 0 {
			result += ", "
		}
		result += c
	}
	return result + "]}"
}

func (s State) Clone() State {
	var result = newState(len(s.candidates))
	result.knownLetters = s.knownLetters
	copy(result.candidates, s.candidates)
	return result
}

func (s State) Regex() string {
	var result = "^"
	for _, c := range s.candidates {
		result += "[" + c + "]"
	}
	return result + "$"
}

func (s State) Grep(wordsFile string) string {
	var result = "grep '" + s.Regex() + "' " + wordsFile
	for _, c := range s.knownLetters {
		result += " | grep " + string(c)
	}
	return result
}

func (s State) Next(guess string, outcome string) State {
	var result = s.Clone()

	if len(guess) != len(s.candidates) {
		panic("guess and state must have same length")
	}
	if len(guess) != len(outcome) {
		panic("guess and outcome must have same length")
	}

	for i, o := range outcome {
		if o == 'x' {
			for j, candy := range result.candidates {
				if len(candy) > 1 {
					result.candidates[j] = strings.Replace(candy, string(guess[i]), "", -1)
				}
			}
		} else if o == 'c' {
			if strings.IndexByte(result.knownLetters, guess[i]) == -1 {
				result.knownLetters += string(guess[i])
			}
			result.candidates[i] = string(guess[i])
		} else if o == 'e' {
			if strings.IndexByte(result.knownLetters, guess[i]) == -1 {
				result.knownLetters += string(guess[i])
			}
			result.candidates[i] = strings.Replace(result.candidates[i], string(guess[i]), "", -1)
		} else {
			panic("unknown outcome: " + string(o))
		}
	}

	return result
}

func NumUniqueChars(str string) int {
	uniqueChars := make(map[rune]bool)
	for _, ch := range str {
		uniqueChars[ch] = true
	}
	return len(uniqueChars)
}

func lettersContained(letters, word string) bool {
	for _, letter := range letters {
		if !strings.ContainsRune(word, letter) {
			return false
		}
	}
	return true
}

func scoreWord(word string) int {
	return NumUniqueChars(word)
}

func (s State) ScoreWords(words []string) []ScoredWord {
	result := []ScoredWord{}
	regex := regexp.MustCompile(s.Regex())
	for _, word := range words {
		if regex.MatchString(word) && lettersContained(s.knownLetters, word) {
			result = append(result, ScoredWord{word, scoreWord(word)})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].score > result[j].score
	})
	return result
}

func Tabulate(scored []ScoredWord, width, height int) string {
	total := width * height
	result := ""
	for i := 0; i < total && i < len(scored); i++ {
		result += scored[i].word
		if (i+1)%width == 0 {
			result += "\n"
		} else {
			result += "  "
		}
	}
	return strings.TrimSpace(result) + "\n"
}

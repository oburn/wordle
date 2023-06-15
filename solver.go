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

var (
	CharToRating = map[rune]int{
		'a': 52,
		'e': 48,
		's': 41,
		'o': 32,
		'r': 31,
		'i': 30,
		'l': 26,
		't': 25,
		'n': 25,
		'u': 21,
		'd': 17,
		'c': 16,
		'y': 15,
		'm': 15,
		'p': 14,
		'h': 14,
		'b': 13,
		'g': 12,
		'k': 11,
		'f': 8,
		'w': 7,
		'v': 6,
		'z': 3,
		'j': 2,
		'x': 2,
		'q': 1,
	}
)

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

func uniqueChars(str string) []rune {
	seen := make(map[rune]bool)
	unique := []rune{}

	for _, char := range str {
		if !seen[char] {
			seen[char] = true
			unique = append(unique, char)
		}
	}

	return unique
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
	unique := uniqueChars(word)
	result := 1_000 * len(unique)
	for _, ch := range unique {
		result += CharToRating[ch]
	}
	return result
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

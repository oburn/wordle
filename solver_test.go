package main

import (
	"reflect"
	"testing"
)

func TestRegex(t *testing.T) {
	var s1 = newState(3)
	s1.candidates[0] = "abc"
	s1.candidates[1] = "def"
	s1.candidates[2] = "xyz"
	var r1 = s1.Regex()
	if r1 != "^[abc][def][xyz]$" {
		t.Error("Expected ^[abc][def][xyz]$, got ", r1)
	}
}

func TestGrep(t *testing.T) {
	state := State{knownLetters: "aez", candidates: []string{"abc", "def", "xyz"}}
	got := state.Grep("words")
	if want := "grep '^[abc][def][xyz]$' words | grep a | grep e | grep z"; want != got {
		t.Error("Wanted ", want, ", got ", got)
	}
}

func TestString(t *testing.T) {
	state := State{knownLetters: "aez", candidates: []string{"abc", "def", "xyz"}}
	got := state.String()
	if want := "{knownLetters: aez, candidates: [abc, def, xyz]}"; want != got {
		t.Error("Wanted ", want, ", got ", got)
	}
}

func TestClone(t *testing.T) {
	var s1 = newState(3)
	s1.candidates[0] = "abc"
	s1.candidates[1] = "def"
	s1.candidates[2] = "xyz"
	s1.knownLetters = "aez"
	var r1 = s1.String()
	const e1 = "{knownLetters: aez, candidates: [abc, def, xyz]}"
	if r1 != e1 {
		t.Error("Expected ", e1, ", got ", r1)
	}

	var s2 = s1.Clone()
	var r2 = s2.String()
	if r2 != e1 {
		t.Error("Expected ", e1, ", got ", r2)
	}

	s2.knownLetters = "new"
	s2.candidates[2] = "thre"
	var r3 = s2.String()
	const e2 = "{knownLetters: new, candidates: [abc, def, thre]}"
	if r3 != e2 {
		t.Error("Expected ", e2, ", got ", r3)
	}
}

func TestNextState(t *testing.T) {
	var state = newState(5)
	var next = state.Next("vodka", "xcxxe")
	const expect = "{knownLetters: oa, candidates: [abcefghijlmnopqrstuwxyz, o, abcefghijlmnopqrstuwxyz, abcefghijlmnopqrstuwxyz, bcefghijlmnopqrstuwxyz]}"
	var got = next.String()
	if expect != got {
		t.Error("Expected ", expect, ", got ", got)
	}
}

func TestNextState_WRONG(t *testing.T) {
	var s1 = newState(5)

	var s2 = s1.Next("tinea", "xxexx")
	const e2string = "{knownLetters: n, candidates: [bcdfghjklmnopqrsuvwxyz, bcdfghjklmnopqrsuvwxyz, bcdfghjklmopqrsuvwxyz, bcdfghjklmnopqrsuvwxyz, bcdfghjklmnopqrsuvwxyz]}"
	if e2string != s2.String() {
		t.Fatal("Expected ", e2string, ", got ", s2.String())
	}
	const e2grep = "grep '^[bcdfghjklmnopqrsuvwxyz][bcdfghjklmnopqrsuvwxyz][bcdfghjklmopqrsuvwxyz][bcdfghjklmnopqrsuvwxyz][bcdfghjklmnopqrsuvwxyz]$' valid-words.txt | grep n"
	if e2grep != s2.Grep("valid-words.txt") {
		t.Fatal("Expected ", e2grep, ", got ", s2.Grep("valid-words.txt"))
	}

	var s3 = s2.Next("blond", "xxccx")
	const e3string = "{knownLetters: no, candidates: [cfghjkmnopqrsuvwxyz, cfghjkmnopqrsuvwxyz, o, n, cfghjkmnopqrsuvwxyz]}"
	if e3string != s3.String() {
		t.Fatal("Expected ", e3string, ", got ", s3.String())
	}
	const e3grep = "grep '^[cfghjkmnopqrsuvwxyz][cfghjkmnopqrsuvwxyz][o][n][cfghjkmnopqrsuvwxyz]$' valid-words.txt | grep n | grep o"
	if e3grep != s3.Grep("valid-words.txt") {
		t.Fatal("Expected ", e3grep, ", got ", s3.Grep("valid-words.txt"))
	}

	var s4 = s3.Next("phony", "xxccx")
	const e4string = "{knownLetters: no, candidates: [cfgjkmnoqrsuvwxz, cfgjkmnoqrsuvwxz, o, n, cfgjkmnoqrsuvwxz]}"
	if e4string != s4.String() {
		t.Fatal("Expected ", e4string, ", got ", s4.String())
	}
	const e4grep = "grep '^[cfgjkmnoqrsuvwxz][cfgjkmnoqrsuvwxz][o][n][cfgjkmnoqrsuvwxz]$' valid-words.txt | grep n | grep o"
	if e4grep != s4.Grep("valid-words.txt") {
		t.Fatal("Expected ", e4grep, ", got ", s4.Grep("valid-words.txt"))
	}
}

func TestUniqueChars(t *testing.T) {
	r1 := uniqueChars("abc")
	if want := ([]rune{'a', 'b', 'c'}); !reflect.DeepEqual(r1, want) {
		t.Fatalf("Wanted %v, got %v", want, r1)
	}
	r1 = uniqueChars("abcabc")
	if want := ([]rune{'a', 'b', 'c'}); !reflect.DeepEqual(r1, want) {
		t.Fatalf("Wanted %v, got %v", want, r1)
	}
	r1 = uniqueChars("abcxabc")
	if want := ([]rune{'a', 'b', 'c', 'x'}); !reflect.DeepEqual(r1, want) {
		t.Fatalf("Wanted %v, got %v", want, r1)
	}
}

func TestScoreWords(t *testing.T) {
	s1 := newState(4)
	s1.candidates[0] = "abcdef"
	s1.candidates[1] = "abcdef"
	s1.candidates[2] = "abcdef"
	s1.candidates[3] = "abcdef"
	s1.knownLetters = "ab"

	inputs := []string{"aaaa", "bbbb", "aabb", "gaba", "abcd", "ddba"}

	r1 := s1.ScoreWords(inputs)
	if len(r1) != 3 {
		t.Fatal("Expected 3, got ", len(r1))
	}
	if want := (ScoredWord{word: "abcd", score: 4}); r1[0] != want {
		t.Fatalf("Expected %v, got %v", want, r1[0])
	}
	if want := (ScoredWord{word: "ddba", score: 3}); r1[1] != want {
		t.Fatalf("Expected %v, got %v", want, r1[1])
	}
	if want := (ScoredWord{word: "aabb", score: 2}); r1[2] != want {
		t.Fatalf("Expected %v, got %v", want, r1[2])
	}
}

func TestTabulate_Excess(t *testing.T) {
	scored := []ScoredWord{
		{word: "aaaa", score: 1},
		{word: "bbbb", score: 2},
		{word: "aabb", score: 3},
		{word: "gaba", score: 4},
		{word: "abcd", score: 5},
		{word: "ddba", score: 6},
		{word: "xxxx", score: 7},
	}
	r1 := Tabulate(scored, 3, 2)
	e1 := `aaaa  bbbb  aabb
gaba  abcd  ddba
`
	if r1 != e1 {
		t.Fatal("Expected ", e1, ", got ", r1)
	}

	r2 := Tabulate(scored, 2, 2)
	e2 := `aaaa  bbbb
aabb  gaba
`
	if r2 != e2 {
		t.Fatal("Expected ", e2, ", got ", r2)
	}

	r3 := Tabulate(scored, 4, 4)
	e3 := `aaaa  bbbb  aabb  gaba
abcd  ddba  xxxx
`
	if r3 != e3 {
		t.Fatal("Expected ", e3, ", got ", r3, "<<<")
	}
}

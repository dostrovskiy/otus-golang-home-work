// Package hw03frequencyanalysis provides top 10 most frequent words
package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var simpleWordRegExp = regexp.MustCompile(`[\p{L}\p{P}]+`) // one or more letters and punctuation

var hardWordRegExp = regexp.MustCompile(
	`\p{L}[\p{L}\p{P}]*\p{L}` + // words with punctuation inside, 2-letter or more
		`|\p{L}` + // 1-letter words: articles, conjunctions, etc.
		`|\p{Pd}{2,}`) // dashes-only word, 2-letter or more

type wordCount struct {
	word  string
	count int
}

func countWords(words []string, hard bool) []wordCount {
	wordCountMap := map[string]int{}
	for _, word := range words {
		if hard {
			word = strings.ToLower(word)
		}
		if _, found := wordCountMap[word]; !found {
			wordCountMap[word] = 1
		} else {
			wordCountMap[word]++
		}
	}
	wordCounts := make([]wordCount, 0, len(wordCountMap))
	for word, count := range wordCountMap {
		wordCounts = append(wordCounts, wordCount{word, count})
	}
	return wordCounts
}

/*
Top10 returns top 10 most frequent words from text. Words are sets of characters separated by spaces.
*/
func Top10(text string) []string {
	return Top10Hard(text, false)
}

/*
Top10Hard returns top 10 most frequent words from text. Words are sets of characters separated by spaces.
hard - if true, the case of words is ignored and punctuation marks after a word are not included.
*/
func Top10Hard(text string, hard bool) []string {
	if text == "" {
		return []string{}
	}
	var textWords []string
	if hard {
		textWords = hardWordRegExp.FindAllString(text, -1)
	} else {
		textWords = simpleWordRegExp.FindAllString(text, -1)
	}
	words := countWords(textWords, hard)
	sort.Slice(words, func(i, j int) bool {
		return words[i].count > words[j].count ||
			words[i].count == words[j].count && words[i].word < words[j].word
	})
	topWords := make([]string, 0, 10)
	for i := 0; i < 10 && i < len(words); i++ {
		topWords = append(topWords, words[i].word)
	}
	return topWords
}

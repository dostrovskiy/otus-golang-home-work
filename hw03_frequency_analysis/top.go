// Package hw03frequencyanalysis provides top 10 most frequent words
package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var hardWordRegExp = regexp.MustCompile(
	`\p{L}[\p{L}\p{P}]*\p{L}` + // words with punctuation inside, 2-letter or more
		`|\p{L}` + // 1-letter words: articles, conjunctions, etc.
		`|\p{Pd}{2,}`) // dashes-only word, 2-letter or more

type wordCount struct {
	word  string
	count int
}

func countWords(words []string) []wordCount {
	wordCountMap := map[string]int{}
	for _, word := range words {
		word = strings.ToLower(word)
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
	if text == "" {
		return []string{}
	}
	textWords := hardWordRegExp.FindAllString(text, -1)
	words := countWords(textWords)
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

package rake

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"unicode"
)

// Score : (Word, Score) pair
type Score struct {
	word  string
	score float64
}

type byScore []Score

func (s byScore) Len() int {
	return len(s)
}

func (s byScore) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byScore) Less(i, j int) bool {
	return s[i].score > s[j].score
}

func getTextFromFile(filename string) string {
	content, _ := ioutil.ReadFile(filename)
	return string(content)
}

func getLinesFromFile(filename string) []string {
	content, _ := ioutil.ReadFile(filename)
	return strings.Split(string(content), "\n")
}

func splitIntoWords(text string) []string {
	words := []string{}
	splitWords := strings.Fields(text)
	for _, word := range splitWords {
		currentWord := strings.ToLower(strings.TrimSpace(word))
		if currentWord != "" {
			words = append(words, currentWord)
		}
	}
	return words
}

func getStopwords() map[string]bool {
	stopwords := getLinesFromFile(stopwordFilename)
	dict := map[string]bool{}
	for _, word := range stopwords {
		dict[word] = true
	}
	return dict
}

func getStopWordRegex() string {
	stopwords := getLinesFromFile(stopwordFilename)
	stopwordRegexPattern := []string{}
	for _, word := range stopwords {
		wordRegex := fmt.Sprintf(forStopWordDetection, word)
		stopwordRegexPattern = append(stopwordRegexPattern, wordRegex)
	}
	return `(?i)` + strings.Join(stopwordRegexPattern, "|")
}

func generateCandidatePhrases(text string) []string {
	stopwords := getStopwords()
	words := splitIntoWords(text)
	acceptedWords := []string{}
	for _, word := range words {
		if !stopwords[word] {
			acceptedWords = append(acceptedWords, word)
		} else {
			acceptedWords = append(acceptedWords, "|")
		}
	}

	phraseList := []string{}
	phrase := ""
	for _, word := range acceptedWords {
		if word == "|" {
			phraseList = append(phraseList, phrase)
			phrase = ""
		} else {
			phrase = phrase + " " + word
		}
	}
	return phraseList
}

func splitIntoSentences(text string) []string {
	splitFunc := func(c rune) bool {
		return unicode.IsPunct(c)
	}
	return strings.FieldsFunc(text, splitFunc)
}

func combineScores(phraseList []string, scores map[string]float64) map[string]float64 {
	candidateScores := map[string]float64{}
	for _, phrase := range phraseList {
		words := splitIntoWords(phrase)
		candidateScore := float64(0.0)

		for _, word := range words {
			candidateScore += scores[word]
		}
		candidateScores[phrase] = candidateScore
	}
	return candidateScores
}

func calculateWordScores(phraseList []string) map[string]float64 {
	frequencies := map[string]int{}
	degrees := map[string]int{}
	for _, phrase := range phraseList {
		words := splitIntoWords(phrase)
		length := len(words)
		degree := length - 1

		for _, word := range words {
			frequencies[word]++
			degrees[word] += degree
		}
	}
	for key := range frequencies {
		degrees[key] = degrees[key] + frequencies[key]
	}

	score := map[string]float64{}

	for key := range frequencies {
		score[key] += (float64(degrees[key]) / float64(frequencies[key]))
	}

	return score
}

func sortScores(scores map[string]float64) []Score {
	rakeScores := []Score{}
	for k, v := range scores {
		rakeScores = append(rakeScores, Score{k, v})
	}
	sort.Sort(byScore(rakeScores))
	return rakeScores
}

func rake(text string) []Score {
	sentences := splitIntoSentences(text)
	phraseList := []string{}
	for _, sentence := range sentences {
		phraseList = append(phraseList, generateCandidatePhrases(sentence)...)
	}
	wordScores := calculateWordScores(phraseList)
	candidateScores := combineScores(phraseList, wordScores)
	sortedScores := sortScores(candidateScores)
	return sortedScores
}

// WithFile : Run rake with text from a file
func WithFile(filename string) map[string]float64 {
	text := getTextFromFile(filename)
	scores := rake(text)
	scoreDict := make(map[string]float64)
	for _, score := range scores {
		scoreDict[score.word] = score.score
	}
	return scoreDict
}

// WithText : Run rake directly from text
func WithText(text string) map[string]float64 {
	scores := rake(text)
	scoreDict := make(map[string]float64)
	for _, score := range scores {
		scoreDict[score.word] = score.score
	}
	return scoreDict
}

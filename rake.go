package rake

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
)

type rakeScore struct {
	word  string
	score float64
}

type byScore []rakeScore

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
	wordSplitRegex := regexp.MustCompile(ForWordSplit)
	splitWords := wordSplitRegex.FindAllString(text, -1)
	for _, word := range splitWords {
		currentWord := strings.ToLower(strings.TrimSpace(word))
		if currentWord != "" {
			words = append(words, currentWord)
		}
	}
	return words
}

func getStopWordRegex() string {
	stopwords := getLinesFromFile(StopwordFilename)
	stopwordRegexPattern := []string{}
	for _, word := range stopwords {
		wordRegex := fmt.Sprintf(ForStopWordDetection, word)
		stopwordRegexPattern = append(stopwordRegexPattern, wordRegex)
	}
	return `(?i)` + strings.Join(stopwordRegexPattern, "|")
}

func generateCandidatePhrases(text string) []string {
	stopWordRegex := regexp.MustCompile(getStopWordRegex())
	temp := stopWordRegex.ReplaceAllString(text, "|")
	multipleWhitespaceRegex := regexp.MustCompile(`\s\s+`)
	temp = multipleWhitespaceRegex.ReplaceAllString(strings.TrimSpace(temp), " ")

	phraseList := []string{}
	phrases := strings.Split(temp, "|")
	for _, phrase := range phrases {
		phrase = strings.ToLower(phrase)
		if phrase != "" {
			phraseList = append(phraseList, phrase)
		}
	}
	return phraseList
}

func splitIntoSentences(text string) []string {
	splitPattern := regexp.MustCompile(ForSplittingSentences)
	return splitPattern.Split(text, -1)
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

func sortScores(scores map[string]float64) []rakeScore {
	rakeScores := []rakeScore{}
	for k, v := range scores {
		rakeScores = append(rakeScores, rakeScore{k, v})
	}
	sort.Sort(byScore(rakeScores))
	return rakeScores
}

func rake() {
	sentences := splitIntoSentences(getTextFromFile(TextFilename))
	phraseList := []string{}
	for _, sentence := range sentences {
		phraseList = append(phraseList, generateCandidatePhrases(sentence)...)
	}
	wordScores := calculateWordScores(phraseList)
	candidateScores := combineScores(phraseList, wordScores)
	sortedScores := sortScores(candidateScores)
	for _, rakeScore := range sortedScores {
		fmt.Println(rakeScore.word, rakeScore.score)
	}
}

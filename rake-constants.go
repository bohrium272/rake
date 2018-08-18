package main

// ForWordSplit - Splitting Words
const ForWordSplit = "[\\p{L}\\d_]+"

// ForStopWordDetection - Filtering stop words
const ForStopWordDetection = `(?:\A|\z|\s)%s(?:\A|\z|\s)`

// ForSplittingSentences - Splitting Sentences
const ForSplittingSentences = `[.,\/#!$%\^&\*;:{}=\-_~()]`

// TextFilename - The file from which text is read for analysis
const TextFilename = "text.txt"

// StopwordFilename - The file of which each line is a stop word
const StopwordFilename = "SmartStoplist.txt"

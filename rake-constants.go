package rake

// forWordSplit - Splitting Words
const forWordSplit = "[^a-zA-Z0-9_\\+\\-/]"

// forStopWordDetection - Filtering stop words
const forStopWordDetection = `\b%s(?![\w-])`

// forSplittingSentences - Splitting Sentences
const forSplittingSentences = `[.!?,;:\t\\\\"\\(\\)\\\'\u2019\u2013]|\\s\\-\\s`

// textFilename - The file from which text is read for analysis
const textFilename = "text.txt"

// stopwordFilename - The file of which each line is a stop word
const stopwordFilename = "SmartStoplist.txt"

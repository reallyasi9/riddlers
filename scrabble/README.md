# [June 28, 2019](https://fivethirtyeight.com/features/whats-your-best-scrabble-string/)
## What's Your Best Scrabble String?

### Riddler Classic

From Benjamin Danard, the Superstring Scrabble Challenge:

The game of Scrabble has 100 tiles — 98 of these tiles contain a letter and a score, and two of them are wildcards worth zero points. At home on a lazy summer day with a bag of these tiles, you decide to play the Superstring Scrabble Challenge. Using only the 100 tiles, you lay them out into one long 100-letter string of your choosing. You look through the string. For each word you find, you earn points equal to its score. Once you find a word, you don’t get any points for finding it again. The same tile may be used in multiple, overlapping words. So "theater" includes "the," "heat," "heater," "eat," "eater," "ate," etc.

The super challenge: What order of tiles gives you the biggest score? (The blank tiles are locked into the letter they represent once you've picked it.)

The winner, and inaugural Wordsmith Extraordinaire of Riddler Nation, will be the solver whose string generates the most points. You should use [this word list](https://norvig.com/ngrams/enable1.txt) to determine whether a word is valid.

For reference, this is the distribution of letter [tiles in the bag](https://en.wikipedia.org/wiki/Scrabble_letter_distributions#English), by their point value:

0: ?×2
1: E×12 A×9 I×9 O×8 N×6 R×6 T×6 L×4 S×4 U×4
2: D×4 G×3
3: B×2 C×2 M×2 P×2
4: F×2 H×2 V×2 W×2 Y×2
5: K
8: J X
10: Q Z

[Submit your string](https://docs.google.com/forms/d/e/1FAIpQLSdAVV14Sc5Mdsd9w7XkPnzJRf81MgL5XWmD2XYKB7hITRciKg/viewform?usp=sf_link)

### Note

I settled on a form of simulated annealing to solve what is essentially the TSP, but with a computationally-intensive distance function.  To make the computations as fast as possible, I stored the word/score list in a custom-built prefix tree, giving me score lookup operations of O(N) at worst (as opposed to O(N) amortized), where N is the length of the string to lookup.  It also has a short-circuit mechanism to stop searching for words of length M+1 if the string of length M is not a prefix of any words in the dictionary.

On a reasonably modern AWS machine, this algorithm can accurately score approximately 100,000 100-letter Scrabble Superstrings per second.

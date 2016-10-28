# Riddler, October 21, 2016

What arrangement of any letters on a Boggle board has the most points attainable? Boggle is played with a 4-by-4 grid of letters. Points are scored by finding strings of letters — connected in any direction, horizontally, vertically or diagonally — that form valid words at least three letters long. Words 3, 4, 5, 6, 7 or 8 or more letters long score 1, 1, 2, 3, 5 and 11 points, respectively. (You can find the full [official rules here](http://www.hasbro.com/common/instruct/boggle.pdf).)

Extra credit: What if you limit the hypothetical configurations to only those that are possible using the actual [letter cubes](http://www.bananagrammer.com/2013/10/the-boggle-cube-redesign-and-its-effect.html) included with the game?

(If you need a word list to aid in your quest, feel free to use [the public-domain ENABLE list](https://storage.googleapis.com/google-code-archive-downloads/v2/code.google.com/dotnetperls-controls/enable1.txt) — a modified version of which is used in Words With Friends.)

[Submit your answer](https://docs.google.com/forms/d/e/1FAIpQLScsrBwAzes_mi2NZjs7NeeE3NGx1-IdbPoBgxVLaSZNfx5Dfw/viewform)

## Author's note

The solution I use is a re-worked solution from the Coursera course [Algorithms, Part II](https://www.coursera.org/learn/java-data-structures-algorithms-2).  I ported my solution to Go and re-optimized the routines for the new language as best I could in the time I had.  The included word lists are probably in the public domain, and the included Boggle boards were given as test cases for us to check our code.  I am assuming that posting these here is is fair use...

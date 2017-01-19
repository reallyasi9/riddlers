# Riddler Express, January 13, 2017

From Conor McMeel, a birthday party puzzle:

It’s your 30th birthday (congrats, by the way), and your friends bought you a cake with 30 candles on it. You make a wish and try to blow them out. Every time you blow, you blow out a random number of candles between one and the number that remain, including one and that other number. How many times do you blow before all the candles are extinguished, on average?

[Submit your answer](https://docs.google.com/forms/d/e/1FAIpQLSff0FMNPEM3wVEo-x__FL2tOJZ0FVx-BkCZIDPh6RpZ3-Kvng/viewform)

## Author's note

This is pretty simple:  you are splitting the set of remaining candles in half each time on average, so the number of times you can do this on average for N candles will look like log(N).  The notebook here just runs a simulation to demonstrate this.

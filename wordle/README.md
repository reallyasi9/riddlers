# [January 14, 2022](https://fivethirtyeight.com/features/when-the-riddler-met-wordle/)
## When The Riddler Met Wordle
### Riddler Classic

Over the last few weeks, [Wordle](https://www.powerlanguage.co.uk/wordle/) has taken the puzzling world by storm. Thousands of people (including yours truly) play daily, and the story of its creation has been [well documented](https://www.nytimes.com/2022/01/03/technology/wordle-word-game-creator.html).

Wordle closely resembles the classic game show [Lingo](https://www.imdb.com/title/tt0423691/). In Wordle, you have six guesses to determine a five-letter mystery word. For each word that you guess, you are told which letters are correct and in the correct position (marked in green), which among the remaining letters are in the mystery word but are in the incorrect position (marked in yellow) and which letters are incorrect altogether.

This sounds straightforward enough. But things get a little hairier when the mystery word or one of your guesses has a letter that appears more than once. To brush up on the rules, you may want to check out the following example in which the mystery word is MISOS, taken from last year’s [Lingo-inspired Riddler Classic](https://fivethirtyeight.com/features/can-you-guess-the-mystery-word/):

[MAGIC (M in green, I in yellow) MAIMS (M and S in green, I in yellow) SUMPS (first S in yellow, M in yellow, last S in green) MOSSO (M in green, first O in yellow, first S in green, second S in yellow) MISOS (all in green)](https://fivethirtyeight.com/wp-content/uploads/2022/01/lingo.png?w=602)

In addition to the many people who play Wordle daily, some folks — including Friends-of-The-Riddler™ [Laurent Lessard](https://twitter.com/LaurentLessard/status/1479981744837308422) and [Tyler Barron](https://www.twitch.tv/videos/1253897976) — have generated approaches for winning Wordle in relatively few guesses, no matter the mystery word. And this is closely related to the task for this week’s Riddler Classic.

Your goal is to devise a strategy to maximize your probability of winning Wordle in at most three guesses. After all, if Yvette Nicole Brown can do it, then so can your strategy! (No offense to Yvette Nicole Brown. She’s awesome.)

>>    I mean…
>>
>>    Alexa play the “It’s me season” quote by [@IssaRae](https://twitter.com/IssaRae?ref_src=twsrc%5Etfw) 😜
>>
>>    Some folks aim for the one and done kinda win. I’m not greedy. Anytime I get it at all, I’m running victory laps. To get it in three is CHRISTMAS![#Wordle](https://twitter.com/hashtag/Wordle?src=hash&ref_src=twsrc%5Etfw) 207 3/6
>>
>>    🟨🟨🟨⬛⬛
>>    🟨🟩⬛🟩🟩
>>    🟩🟩🟩🟩🟩
>>    — yvette nicole brown (@YNB) [January 12, 2022](https://twitter.com/YNB/status/1481221402447335425?ref_src=twsrc%5Etfw)

In particular, I want to know (1) your strategy, (2) the first word you would guess and (3) your probability of winning in three or fewer guesses.

To do this, you will need to access Wordle’s library of [2,315 mystery words](https://docs.google.com/spreadsheets/d/1-M0RIVVZqbeh0mZacdAsJyBrLuEmhKUhNaVAI-7pr2Y/edit#gid=0) as well as all [12,972 words you are allowed to guess](https://docs.google.com/spreadsheets/d/1KR5lsyI60J1Ek6YgJRU2hKsk4iAOWvlPLUWjAZ6m8sg/edit#gid=0). For the record, I pulled both of these word lists from Wordle’s source code and listed them alphabetically for your convenience.

Spoiler alert! If you enjoy playing Wordle daily and do not want to know the entire list of mystery words, then don’t look too closely at these lists. You have been warned!

[Submit your answer](https://docs.google.com/forms/d/e/1FAIpQLSd--4b02inswC5dfC3g94gv7PzWURh8-lqdZCYQHR8PHiY41Q/viewform?usp=sf_link)

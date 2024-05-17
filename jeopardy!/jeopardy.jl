### A Pluto.jl notebook ###
# v0.12.10

using Markdown
using InteractiveUtils

# ╔═╡ c2a644a0-25ef-11eb-2b15-e7a70de1da1e
md"""## Riddler Express

Last Sunday we lost [Alex Trebek](https://fivethirtyeight.com/features/remembering-alex-trebek-the-man-with-all-the-answers/), a giant in the world of game shows and trivia. The show he hosted over the course of four decades — Jeopardy! — has [previously appeared](https://fivethirtyeight.com/features/how-much-is-a-perfect-game-of-jeopardy-worth/) in this column. Today, it makes a return.

You’re playing the (single) Jeopardy! Round, and your opponents are simply no match for you. You choose first and never relinquish control, working your way horizontally across the board by first selecting all six $200 clues, then all six $400 clues, and so on, until you finally select all the $1,000 clues. You respond to each clue correctly before either of your opponents can.

One randomly selected clue is a Daily Double. Rather than award you the prize money associated with that clue, it instead allows you to double your current winnings or wager up to $1,000 should you have less than that. Being the aggressive player you are, you always bet the most you can. (In reality, the Daily Double is more likely to appear in [certain locations](https://www.vox.com/2015/3/3/8140405/jeopardy-daily-double-statistics) on the board than others, but for this problem assume it has an equal chance of appearing anywhere on the board.)

How much money do you expect to win during the Jeopardy! round?

*Extra credit*: Suppose you change your strategy. Instead of working your way horizontally across the board, you select random clues from anywhere on the board, one at a time. Now how much money do you expect to win during the Jeopardy! round?"""

# ╔═╡ 184cfed0-25f0-11eb-10d6-41be46d7e219
md"""---
If you are responding to the clues across the rows of the game board, that means after responding to the first clue you will have earned \\\$200; after the second, \\\$400; after the third, \\\$600; and so on. If the Daily Double is under clues one through five, then the value of the clue effectively becomes \\\$1000 because your total score up to that point on the board is less than \\\$1000. If the Daily Double is under any of the following clues, the value of that clue is effectively your total score at that point on the board. This means a Daily Double under clue six is worth \\\$1000; under clue seven, \\\$1200; under clue eight, \\\$1600 (because clue seven is worth \\\$400); and so on."""

# ╔═╡ c8787430-25f3-11eb-0c76-e1a6a570a690
"""
	cluevalue(x)

Return the value of clue `x`, where `x` is an integer in the range `0:29`.
"""
cluevalue(x) = 200 * (x ÷ 6 + 1)

# ╔═╡ b67d3940-25f4-11eb-2abf-2ffabbf4d40b
"""
	ddvalue(x)

Return the Daily Double value of a clue when the player has score `score`.
"""
function ddvalue(score)
	score ≤ 1000 && return 1000
	score
end

# ╔═╡ 1de8988e-25f5-11eb-37f3-2311c511a8a3
md"""Responding to all the clues correctly in absence of the Daily Double means you earn $$6 \times \left( \$200 + \$400 + \$600 + \$800 + \$1000 \right) = \$18000$$. That's no joke."""

# ╔═╡ 2a7ea450-25f5-11eb-3a64-bf53c91dfb29
sum(cluevalue.(0:29))

# ╔═╡ 67e31bb0-25fe-11eb-39a7-a3927959a2bf
md"""The worst thing that could happen to you earnings-wise is if the Daily Double appears somewhere in the first seven clues. In that case, you only earn an extra \\\$800."""

# ╔═╡ b3c0f3e0-25fe-11eb-3686-5dc6cb6b274a
ddvalue(0) + sum(cluevalue.(1:29))

# ╔═╡ d197cebe-25fe-11eb-2242-13dc718c615a
md"""The best thing that could happen is if the Daily Double is the very last clue, earning you a whopping \\\$16000 extra!"""

# ╔═╡ f8fd8fe0-25fe-11eb-1f45-074d7b6dbde8
sum(cluevalue.(0:28)) + ddvalue(sum(cluevalue.(0:28)))

# ╔═╡ b8e55fa0-25f4-11eb-3bc0-3d043441c7fd
md"""If the Daily Double has a $$\frac{1}{30}$$ chance of being under any of the 30 spaces on the board, then you have a $$\frac{1}{30}$$ chance of that space's value being replaced by either \\\$1000 or your current score, whichever is largest. That means the _expected value_ of the first clue is $$\$200 \times \frac{29}{30} + \$1000 \times \frac{1}{30}$$ because 29 out of 30 times the clue will be a regular old clue, but one out of 30 times it will be the Daily Double and net you \\\$1000 instead! The same goes for the expected value of other clues: for example, the expected value of the eighth clue is $$\$400 \times \frac{29}{30} + \$1600 \times \frac{1}{30}$$."""

# ╔═╡ faddcd20-25f4-11eb-10bc-6d5c51f3e9ed
md"""Your expected value from solving all of the clues is just the sum of the expected values of each clue, with the first clue extracted from the sum for programming convenience (note the limits of the inner sum):

$$\left[\frac{29}{30}\text{cluevalue}(0) + \frac{1}{30}\text{ddvalue}(0)\right] + \sum_{i=1}^{29}\left[\frac{29}{30} \text{cluevalue}(i) + \frac{1}{30} \text{ddvalue}\left( \sum_{j=0}^{i-1} \text{cluevalue}(j) \right) \right]$$
"""

# ╔═╡ 34180650-25ff-11eb-02e7-b9858330f10b
29/30 * cluevalue(0) + 1/30 * ddvalue(0) + sum(29/30 * cluevalue(i) + 1/30 * ddvalue(sum(cluevalue(j) for j ∈ 0:i-1)) for i ∈ 1:29)

# ╔═╡ 45144172-25f6-11eb-24f8-d7ac7cb07cf3
md"""To demonstrate that this is the answer, let's simulate 1,000,000 games by randomly drawing a Daily Double location for each game and playing the game out."""

# ╔═╡ 6193547e-25f6-11eb-0889-8f3603e89704
"""
	play(dd; order)

Play a perfect simulated game of Jeopardy! using `order` as the order of selected clues and a Daily Double at location `dd`, where values of `order` and `dd` are in the range `0:29`.
"""
function play(dd; order=0:29)
	score = 0
	for i ∈ order
		if i==dd
			score += ddvalue(score)
		else
			score += cluevalue(i)
		end
	end
	score
end

# ╔═╡ ae3b6510-25f7-11eb-238e-1b8e74b1c690
begin
	using Statistics
	N = 1_000_000
	games = play.(rand(0:29, N))
	mean(games)
end

# ╔═╡ e69c8b0e-25fb-11eb-237d-015f7e325b6f
begin
	using StatsPlots
	p = histogram(
		games,
		bins=18000:1000:36000,
		xlabel="Score",
		ylabel="Number of games",
		label="By row",
	)
end

# ╔═╡ fccaddd0-25f9-11eb-0d0e-3516abbef53d
begin
	using Random
	randomgames = [play(x[1]; order=x[2]) for x ∈ Iterators.zip(rand(0:29, N), (shuffle(0:29) for _ ∈ 1:N))]
	mean(randomgames)
end

# ╔═╡ 7303b7c0-25f9-11eb-3334-7385695b931a
md"""---
For extra credit, we can just update the simulation tool to allow for random order play and simulate another 10,000,000 games to figure out how different the expected value of those games are compared to row-by-row play."""

# ╔═╡ 10f61b50-25fd-11eb-1e23-ff75c9f4473d
stephist!(p,
	randomgames,
	bins=18000:1000:36000,
	label="Random order",
	linewidth=3,
)

# ╔═╡ Cell order:
# ╟─c2a644a0-25ef-11eb-2b15-e7a70de1da1e
# ╟─184cfed0-25f0-11eb-10d6-41be46d7e219
# ╠═c8787430-25f3-11eb-0c76-e1a6a570a690
# ╠═b67d3940-25f4-11eb-2abf-2ffabbf4d40b
# ╟─1de8988e-25f5-11eb-37f3-2311c511a8a3
# ╠═2a7ea450-25f5-11eb-3a64-bf53c91dfb29
# ╟─67e31bb0-25fe-11eb-39a7-a3927959a2bf
# ╠═b3c0f3e0-25fe-11eb-3686-5dc6cb6b274a
# ╟─d197cebe-25fe-11eb-2242-13dc718c615a
# ╠═f8fd8fe0-25fe-11eb-1f45-074d7b6dbde8
# ╟─b8e55fa0-25f4-11eb-3bc0-3d043441c7fd
# ╟─faddcd20-25f4-11eb-10bc-6d5c51f3e9ed
# ╠═34180650-25ff-11eb-02e7-b9858330f10b
# ╟─45144172-25f6-11eb-24f8-d7ac7cb07cf3
# ╠═6193547e-25f6-11eb-0889-8f3603e89704
# ╠═ae3b6510-25f7-11eb-238e-1b8e74b1c690
# ╠═e69c8b0e-25fb-11eb-237d-015f7e325b6f
# ╟─7303b7c0-25f9-11eb-3334-7385695b931a
# ╠═fccaddd0-25f9-11eb-0d0e-3516abbef53d
# ╠═10f61b50-25fd-11eb-1e23-ff75c9f4473d

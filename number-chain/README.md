# [July 28, 2017](https://fivethirtyeight.com/features/pick-a-number-any-number/)
## Pick A Number, Any Number

### Riddler Classic

From Itay Bavly, a chain-link number problem:

You start with the integers from one to 100, inclusive, and you want to organize them into a chain. The only rules for building this chain are that you can only use each number once and that each number must be adjacent in the chain to one of its factors or multiples. For example, you might build the chain:

4, 12, 24, 6, 60, 30, 10, 100, 25, 5, 1, 97

You have no numbers left to place after 97, leaving you with a finished chain of length 12.

What is the longest chain you can build?

Extra credit: What if you started with more numbers, e.g., one through 1,000?

[Submit your answer](https://docs.google.com/forms/d/e/1FAIpQLSfFKcfBsrIN690ooWIo2AoprIUHt-K0IUGmFaih6HWzDvXLVA/viewform?usp=sf_link)

### Note

The code treats the problem as a Hamiltonian path-finding problem, known to be NP-complete.  Fortunately, we can stop searching once we find a Hamiltonian path, and if we find that the longest path starting at N ends at M, we can be sure that the longest path starting at M ends at N because the graph is undirected.  This cuts down our search space, but not by much.

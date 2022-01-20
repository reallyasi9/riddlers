using ArgParse
using ProgressMeter
using Combinatorics
using Random

function parse_arguments(args=ARGS)
    s = ArgParseSettings()
    @add_arg_table! s begin
        "solutions"
            help = "file containing solutions to Wordle"
            required = true
        "guessables"
            help = "file containing guessable words"
            required = true
        "--guesses,-g"
            help = "(exact) number of guesses to optimize"
            arg_type = Int
            default = 2
    end

    parsed_args = parse_args(args, s)

    return parsed_args
end

struct Word
    letters::NTuple{5, UInt8}
end

function Word(word::String)
    @inbounds Word((word[1], word[2], word[3], word[4], word[5]) .- '`')
end

letters(w::Word) = w.letters

function Base.show(io::IO, x::Word)
    show(io, String(collect(x.letters .+ '`')))
end

function Base.print(io::IO, x::Word)
    print(io, String(collect(x.letters .+ '`')))
end

import Base.getindex
function Base.getindex(w::Word, i)
    return getindex(w.letters, i)
end

function Base.getindex(w::Word, cs::NTuple{5, Bool})
    return getindex(w.letters, collect(cs))
end

function Base.isdisjoint(x::Word, y::Word)
    return isdisjoint(x.letters, y.letters)
end

function Base.isdisjoint(x::Vector{Word})
    return isempty(intersect(letters.(x)))
end

struct Wordle
    posmap::NTuple{5,NTuple{26,BitSet}}
    lettermap::NTuple{26,BitSet}
    missingmap::NTuple{26,BitSet}
    words::Vector{Word}
end

function nwords(w::Wordle)
    return length(w.words)
end

function getword(w::Wordle, i)
    return w.words[i]
end

function wordswithletter(w::Wordle, c::UInt8)
    return @inbounds w.lettermap[c]
end

function wordsmissingletter(w::Wordle, c::UInt8)
    return @inbounds w.missingmap[c]
end

function wordsmatching(w::Wordle, c::Pair{Int,UInt8})
    return @inbounds w.posmap[c[1]][c[2]]
end

const ALPHABET = "abcdefghijklmnopqrstuvwxyz"

function Wordle(words)
    lettermap = ntuple(x->BitSet(), 26)
    missingmap = ntuple(x->BitSet(), 26)
    posmap = ntuple(y->ntuple(x -> BitSet(), 26), 5)
    for i in eachindex(words)
        for c in ALPHABET
            char = c - 'a' + 1
            if c in words[i]
                push!(lettermap[char], i)
                for pos in findall(c, words[i])
                    push!(posmap[pos][char], i)
                end
            else
                push!(missingmap[char], i)
            end
        end
    end

    return Wordle(posmap, lettermap, missingmap, Word.(words))
end

function Base.show(io::IO, w::Wordle)
    show(io, w.words)
end

function guess(target::Word, g::Word)
    matches = g.letters .== target.letters
    contains = in.(g.letters, (target.letters,))
    return [i=>g[i] for i in findall(matches)], [g[i] for i in findall(contains)], [g[i] for i in findall(.!contains)]
end

function emptyintersect(s, itrs...)
    isempty(s) && return intersect(itrs...)
    return intersect(s, itrs...)
end

function possibilities(w::Wordle, matches::Vector{Pair{Int,UInt8}}, contains::Vector{UInt8}, misses::Vector{UInt8})
    ma = mapreduce((c)->wordsmatching(w, c), emptyintersect, matches; init=BitSet())
    co = mapreduce((c)->wordswithletter(w, c), emptyintersect, contains; init=ma)
    return mapreduce((c)->wordsmissingletter(w, c), emptyintersect, misses; init=co)
end

function ambiguities(w::Wordle, target::Word, gs::Vector{Word})
    possible = mapreduce(x->possibilities(w, guess(target, x)...), intersect, gs)
    1 / length(possible)
end

function main(args=ARGS)
    a = parse_arguments(args)

    solutions = readlines(a["solutions"])
    guessables = Word.(readlines(a["guessables"]))

    wordle = Wordle(solutions)
    combs = collect(combinations(guessables, a["guesses"]))
    filter!(isdisjoint, combs)
    shuffle!(combs)
    ncombs = length(combs)
    # ncombs = 100
    best_combo = first(combs)
    best_prob = 0.
    p = Progress(ncombs; showspeed=true)
    solutions = Word.(solutions)
    Threads.@threads for combo in combs
        prob = mapreduce(x -> ambiguities(wordle, x, combo), +, solutions)
        if prob > best_prob
            best_prob = prob
            best_combo = combo
        end
        next!(p; showvalues=[(:best_combo,join(best_combo, " + ")), (:prob, best_prob/length(solutions))])
    end
    return best_combo, best_prob
end

# combs, ambicount = main();
main()

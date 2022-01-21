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
        "--guesses","-g"
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
    isempty(x) && return true
    return length(unique(mapreduce(y->collect(letters(y)), vcat, x))) == length(x)
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

struct ComboProb
    combo::Vector{Word}
    prob::Float64
end

combination(c::ComboProb) = c.combo
probability(c::ComboProb) = c.prob

const BEST_COMBOPROB_SO_FAR = ComboProb(Word[], 0.)
const WORDS = Channel{Vector{Word}}(100)
const UNFILTERED_COMBOPROBS = Channel{ComboProb}(100)
const FILTERED_COMBOPROBS = Channel{ComboProb}(Inf)

function comboprob_filter()
    for combo in UNFILTERED_COMBOPROBS
        if probability(combo) > probability(best)
            BEST_COMBOPROB_SO_FAR = combo
            put!(FILTERED_COMBOPROBS, combo)
        end
    end
    close(FILTERED_COMBOPROBS)
end

function calculate_comboprob(wordle::Wordle, solutions::Vector{Word})
    for combo in WORDS
        if !isdisjoint(combo)
            continue
        end
        prob = mapreduce(x -> ambiguities(wordle, x, combo), +, solutions)
        comboprob = ComboProb(combo, prob)
        put!(UNFILTERED_COMBOPROBS, comboprob)
    end
    close(UNFILTERED_COMBOPROBS)
end

function main(args=ARGS)
    a = parse_arguments(args)
    @info "starting" args=args parsed_args=a

    solutions = readlines(a["solutions"])
    guessables = Word.(readlines(a["guessables"]))

    wordle = Wordle(solutions)
    combs = combinations(guessables, a["guesses"])
    ncombs = length(combs)
    # ncombs = 100
    p = Progress(ncombs; showspeed=true)
    solutions = Word.(solutions)

    for i in 1:Threads.nthreads()
        @info "We starting this?" i
        @async errormonitor(calculate_comboprob($wordle, $solutions))
    end
    @info "Is this starting?"
    @async errormonitor(comboprob_filter())

    for combo in combs
        println("trying $combo")
        put!(WORDS, combo)
        next!(p; showvalues=[(:best_combo,join(combination(BEST_COMBOPROB_SO_FAR), " + ")), (:prob, probability(BEST_COMBOPROB_SO_FAR)/length(solutions))])
    end

    close(WORDS)
    # deal with the history of the process to make sure we have the correct values
    for result in FILTERED_COMBOPROBS
        println("$(join(combination(result), " + ")) => $(probability(result))")
    end

    return
end

# combs, ambicount = main();
main()

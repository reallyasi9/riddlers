using ArgParse
using Combinatorics
using ProgressMeter

function parse_arguments(args=ARGS)
    s = ArgParseSettings()
    @add_arg_table! s begin
        "solutions"
            help = "file containing solutions to Wordle"
            required = true
        "guessables"
            help = "file containing guessable words"
            required = true
    end

    parsed_args = parse_args(args, s)

    return parsed_args
end

struct Wordle
    lettermap::Dict{Char,Vector{String}}
    invmap::Dict{Char,Vector{String}}
    posmap::Dict{Tuple{Char,Int},Vector{String}}
    words::Vector{String}
end

const ALPHABET = "abcdefghijklmnopqrstuvwxyz"

function Wordle(words)
    lettermap = Dict{Char,Vector{String}}(
        c => String[] for c in ALPHABET
    )
    invmap = Dict{Char,Vector{String}}(
        c => String[] for c in ALPHABET
    )
    posmap = Dict{Tuple{Char,Int},Vector{String}}()
    for c in ALPHABET
        for i in 1:5
            posmap[(c,i)] = String[]
        end
    end
    for word in words
        for letter in ALPHABET
            if letter in word
                push!(lettermap[letter], word)
                for pos in findall(letter, word)
                    push!(posmap[(letter, pos)], word)
                end
            else
                push!(invmap[letter], word)
            end
        end
    end
    return Wordle(lettermap, invmap, posmap, words)
end

function guess(target::String, guesses)
    combined = reduce(*, guesses)
    letters = intersect(target, combined)
    misses = setdiff(combined, target)

    matches = Int[]
    for guess in guesses
        append!(matches, findall(collect(target) .== collect(guess)))
    end
    matchtupples = [(target[i],i) for i in unique(matches)]
    
    return letters, misses, matchtupples
end

function guess(target::String, guess::String)
    return guess(target, (guess,))
end

function ambiguities(w::Wordle, target::String, guesses)
    letters, misses, matchtupples = guess(target, guesses)
    wlm = isempty(letters) ? String[] : mapreduce(x->get(w.lettermap, x, ""), intersect, letters)
    winv = isempty(misses) ? String[] :  mapreduce(x->get(w.invmap, x, ""), intersect, misses)
    wmatch = isempty(matchtupples) ? String[] : mapreduce(x->get(w.posmap, x, ""), intersect, matchtupples)
    return intersect(wlm, winv, wmatch)
end

function main(args=ARGS)
    a = parse_arguments(args)

    solutions = readlines(a["solutions"])
    guessables = readlines(a["guessables"])

    wordle = Wordle(solutions)
    combs = collect(combinations(guessables, 2))
    # ncombs = length(combs)
    ncombs = 100
    ambicount = Vector{Tuple{Int,String,String}}(undef, ncombs)
    p = Progress(ncombs)
    Threads.@threads for i in 1:ncombs
        guesses = combs[i]
        n = mapreduce(x->length(ambiguities(wordle, x, guesses)), +, solutions; init=0)
        ambicount[i] = (n, guesses...)
        next!(p)
    end

    sort!(ambicount; lt = (x,y) -> x[1] < y[1])
    println(first(ambicount, 100))
    return ambicount
end

ambicount = main();

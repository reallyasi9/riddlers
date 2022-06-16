function parsewordle(word)
    [i => word[i] for i in 1:length(word) if word[i] ∈ 'a':'z']
end

is(letters) = function(word)
    all([word[pos] == letter for (pos, letter) in letters])
end

has(letters) = function(word)
    all([(letter in word) && (word[pos] != letter) for (pos, letter) in letters])
end

hasnot(letters) = function(word)
    all([letter ∉ word for letter in letters])
end

function wordle(iswords::Vector{String}, haswords::Vector{String}, hasnotletters, words)
    isletters = mapreduce(parsewordle, vcat, iswords)
    hasletters = mapreduce(parsewordle, vcat, haswords)
    return wordle(isletters, hasletters, hasnotletters, words)
end

wordle(isletters, hasletters, hasnotletters, words) = filter(hasnot(hasnotletters), filter(has(hasletters), filter(is(isletters), words)))
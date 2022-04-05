is(letters) = function(word)
    all([word[pos] == letter for (pos, letter) in letters])
end

has(letters) = function(word)
    all([(letter in word) && (word[pos] != letter) for (pos, letter) in letters])
end

hasnot(letters) = function(word)
    all([letter ∉ word for letter in letters])
end

wordle(isletters, hasletters, hasnotletters, words) = filter(hasnot(hasnotletters), filter(has(hasletters), filter(is(isletters), words)))
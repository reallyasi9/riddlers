include("FlagGuesser.jl")

import .FlagGuesser
using Random
using FileIO

f1 = load("flag_1.png.webp")
f2 = load("flag_2.png.webp")
f3 = load("flag_3.png.webp")

for guess ∈ sort([(v => k) for (k, v) in FlagGuesser.distances(f1)])[1:10]
    println(guess)
end
println()

for guess ∈ sort([(v => k) for (k, v) in FlagGuesser.distances(f2)])[1:10]
    println(guess)
end
println()

for guess ∈ sort([(v => k) for (k, v) in FlagGuesser.distances(f3)])[1:10]
    println(guess)
end
println()

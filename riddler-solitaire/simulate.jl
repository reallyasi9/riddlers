using AlgebraOfGraphics, CairoMakie, TypedTables
using Random

simulate(; cards=10, rng=Random.GLOBAL_RNG) = randperm(rng, cards + 1) .- 1
whentostop(v) = findfirst(==(0), v)
function earnings(v, stop)
    isnothing(stop) && return 0
    stop <= 0 && return 0
    sum(@view v[1:stop])
end

function mc(trials; cards=10, rng=Random.GLOBAL_RNG)
    stops = zeros(Int, trials)
    winnings = zeros(Int, trials)
    deck = collect(1:cards+1) .- 1
    for i = 1:trials
        shuffle!(rng, deck)
        stop = whentostop(deck)
        win = earnings(deck, stop)

        stops[i] = stop
        winnings[i] = win
    end
    stops, winnings
end

function plotmc(stops, winnings)
    t = Table(value = vcat(stops .- 1, winnings), variable = vcat(fill("stop card", length(stops)), fill("winnings", length(winnings))))
    axis = (width = 500, height = 250)
    specs = data(t) * mapping(:value, col=:variable) * histogram(bins=0:55)
    draw(specs; axis)
end

plotmc(mc(10000)...)
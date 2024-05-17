using Random
using Statistics
using StatsBase
using CairoMakie

struct Grasshopper
    history::Vector{Float64}

    function Grasshopper(x::Float64 = 0.)
        h = [x]
        new(h)
    end
end

function history(gh::Grasshopper)
    copy(gh.history)
end

function jump!(rng::AbstractRNG, gh::Grasshopper)
    pos = last(gh.history)
    left = max(0., pos - 0.2)
    right = min(1., pos + 0.2)
    pos = left + rand(rng) * (right - left)
    push!(gh.history, pos)
    pos
end

let gh = Grasshopper(), rng = Random.Xoshiro(0x42)
    for n in 1:1_000_000
        jump!(rng, gh)
    end
    h = history(gh)

    fig = Figure()
    density(fig[1,1], h)

    fig
end

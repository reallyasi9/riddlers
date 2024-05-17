using Plots

function findfailure(a)
    ã = reverse(a)
    s = cumsum(ã) ./ collect(1:length(a))
    return findfirst(abs.(s .- ã) .> 0.5)
end

plt = histogram(
    xlim=(0,100),
    ylim=(0,1),
    title="Stacks to failure",
    bins=1:100,
    normalize=:density,
)
f = Vector{Int}()
@gif for i ∈ 1:1000
    x = findfailure(cumsum(rand(100) .- 0.5))
    isnothing(x) && continue
    push!(f, x)
    histogram!(plt, f)
end every 1
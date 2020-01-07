using Colors
using Interpolations
using KernelDensity
using Distances

function distance(metric::PreMetric, i1::InterpKDE, i2::InterpKDE, range::AbstractRange)
    v1 = pdf(i1, range)
    v2 = pdf(i2, range)
    gzs = (v1 .> 0) .& (v2 .> 0)
    metric(v1[gzs], v2[gzs])
end

aspectratio(a::AbstractMatrix) = size(a, 1) / size(a, 2)

function distances(f::AbstractArray{<:Colorant})
    ar1 = aspectratio(f)
    f = vec(convert(Array{Lab}, f))
    l1 = InterpKDE(kde(comp1.(f)))
    a1 = InterpKDE(kde(comp2.(f)))
    b1 = InterpKDE(kde(comp3.(f)))
    metric = BhattacharyyaDist()
    rl = 0:1:100
    rab = -100:1:100
    dists = Dict{AbstractString, Real}()
    for (key, val) in flagsMap
        f2 = flag(key)
        ar2 = aspectratio(f2)
        f2 = vec(convert(Array{Lab}, f2))
        l2 = InterpKDE(kde(comp1.(f2)))
        a2 = InterpKDE(kde(comp2.(f2)))
        b2 = InterpKDE(kde(comp3.(f2)))

        d = sqrt(distance(metric, l1, l2, rl) + distance(metric, a1, a2, rab) + distance(metric, b1, b2, rab))

        d *= (1 + abs(ar1 - ar2))
        dists[val] = d
    end
    dists
end

using Colors
using HypothesisTests

function distance(flag1::AbstractArray{<:Colorant}, flag2::AbstractArray{<:Colorant})
    lab1 = vec(convert(Array{Lab}, flag1))
    lab2 = vec(convert(Array{Lab}, flag2))
    p1 = pvalue(KSampleADTest(comp1.(lab1), comp1.(lab2)), nsim=20)
    p2 = pvalue(KSampleADTest(comp2.(lab1), comp2.(lab2)), nsim=20)
    p3 = pvalue(KSampleADTest(comp3.(lab1), comp3.(lab2)), nsim=20)
    return 1 / (p1 + p2 + p3)
end

aspectratio(a::AbstractMatrix) = size(a, 1) / size(a, 2)

function distances(f::AbstractArray{<:Colorant})
    ar = aspectratio(f)
    f = vec(convert(Array{Lab}, f))
    d = Dict{AbstractString, Real}()
    for (key, val) in flagsMap
        f2 = flag(key)
        v = distance(f, f2)
        v *= (1 + abs(ar - aspectratio(f2)))
        d[val] = v
    end
    d
end

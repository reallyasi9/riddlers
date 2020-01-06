using Colors, ColorSchemeTools
using StatsBase: countmap

function pixelcount(flag::Array{T}) where T <: Colorant
    a = hex.(convert(Array{RGB}, flag))
    d = Dict(parse(XYZ, "#"*k) => v for (k,v) in countmap(a))
    reverse!(sort(collect(d), by=x->x[2]))
end

function dot(a::Colorant, b::Colorant)
    axyz = convert(XYZ, a)
    bxyz = convert(XYZ, b)
    comp1(axyz) * comp1(bxyz) + comp2(axyz) * comp2(bxyz) + comp3(axyz) + comp3(bxyz)
end

norm(a::Colorant) = sqrt(dot(a, a))

function dist(pc1, pc2; num_colors::Integer=4)
    # Look at the first num_colors numbers and check the differences
    d = 0
    for i in 1:min(num_colors, length(pc1), length(pc2))
        d += sqrt(dot(pc1[i][1], pc2[i][1])) * (abs(pc1[i][2] - pc2[i][2]) + 1)
    end
    d
end

aspectratio(a::Matrix) = size(a, 1) / size(a, 2)

function distances(f; num_colors::Integer=4)
    pc = pixelcount(f)
    ar = aspectratio(f)
    d = Dict{AbstractString, Real}()
    for (key, val) in flagsMap
        m = extract(flagfile(key), num_colors, 10)

        f2 = convert_to_scheme(m, flag(key))

        v = dist(pc, pixelcount(f2))
        v *= (1 + abs(ar - aspectratio(f2)))
        d[val] = v
    end
    d
end

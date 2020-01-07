using FileIO, ImageMagick

function flagfile(country::Symbol)
    if country ∉ keys(flagsMap)
        @error "country not recognized" country
        return
    end
    joinpath(@__DIR__, "flags", string(country) * ".gif")
end

flag(country::Symbol) = load(flagfile(country))

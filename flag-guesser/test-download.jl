include("FlagGuesser.jl")

import .FlagGuesser

abbrs = FlagGuesser.downloadflags()
@show abbrs

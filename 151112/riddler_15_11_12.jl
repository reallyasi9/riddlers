function binary(firstbreak::Integer, maxf::Integer, startf::Integer = 0, verbose::Bool = false)
  #println("Going from $startf to $maxf to find $firstbreak")
  if firstbreak == 1
    if verbose; print("1*"); end
    return (1, 1, true)
  end

  # Go until I break the first time
  f = nextpow2(firstbreak-startf)
  f = clamp(f, 0, maxf-startf)
  #println("Next pow $f")
  tbreak = Int(exponent(float(f)) + 1)
  #println("$tbreak so far")
  if verbose; print(join(startf + 2 .^ (0:tbreak-1), ", ")); end

  # If we hit the maximum, do another binary search
  if f == maxf
    sf = startf + prevpow2(firstbreak-startf)
    #println("New start for binary: $sf")
    if verbose; print(", "); end
    return (tbreak, sf, false)
  end

  if verbose; print("*"); end
  fstart = startf + prevpow2(firstbreak - 1)
  #println("New start: $fstart")
  return (tbreak, fstart, true)
end

function trials(firstbreak::Integer, maxf::Integer, verbose::Bool = false)

  if verbose; print("$firstbreak: "); end

  (tbreak, fstart, broken) = binary(firstbreak, maxf, 0, verbose)
  while !broken
    (tb, fstart, broken) = binary(firstbreak, maxf, fstart, verbose)
    tbreak += tb
  end

  # Linear search to the floor, but don't repeat work
  tlinear = firstbreak - fstart
  if ispow2(firstbreak)
    tlinear -= 1
  end
  tlinear = clamp(tlinear, 0, tlinear)

  if verbose
    if tlinear > 0
      print(", ")
      if ispow2(firstbreak)
        print(join(fstart+1:firstbreak-1, ", "))
      else
        print(join(fstart+1:firstbreak, ", "), "*")
      end
    end
    println(" = ", tbreak + tlinear)
  end

  return tbreak + tlinear
end

a100 = 1:100
t100 = map(a100) do f
  trials(f, 100)
end
println("max: ", findmax(t100))
trials(a100[indmax(t100)], 100, true)
println("median: ", median(t100), " mean: ", mean(t100))

t1000 = map(1:1000) do f
  trials(f, 1000)
end
println("max: ", findmax(t1000))
trials(indmax(t1000), 1000, true)
println("median: ", median(t1000), " mean: ", mean(t1000))

#=
1,2,4,8,16,32,64,65,67,69,73,81,82,84,88,96,97,99

=#

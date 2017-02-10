using DataFrames

const N = parse(ARGS[1])
const M = parse(ARGS[2])

const pickMatrix = Array{Int}(N, M)
const uberMatrix = Array{Float64}(N, M)

function binarySearch(imid::Int, target::Int, imax::Int)
  #println("Looking for $target from $start")
  s = Int(0)

  imin = Int(1)

  while imin <= imax
    s += 1
    #println("Step $s")

    if imid == target
      #println("=$imid")
      return s
    elseif imid < target
      #println(">$imid")
      imin = imid + 1
    else
      #println("<$imid")
      imax = imid - 1
    end

    imid = imin + div(imax - imin, 2)

  end
end

@time for n in 1:N
  const resultsMatrix = Array{Int}(n, n)
  for i in 1:n
    for j in 1:n
      resultsMatrix[i, j] = binarySearch(i, j, n)
    end
  end

  #println(resultsMatrix)

  # With 9 guesses, this becomes kind of fun...
  for m in 1:M
    const found = resultsMatrix .<= m
    valueMatrix = repmat(collect(1:n), 1, n)
    #println(valueMatrix)

    valueMatrix .*= found
    #println(valueMatrix)

    const means = mean(valueMatrix, 1)
    const (best, ibest) = findmax(means)

    #println("Best: $(findin(means, best)) -> $best")

    pickMatrix[n, m] = ibest
    uberMatrix[n, m] = best
  end
end

pickDf = convert(DataFrame, pickMatrix)
stackCols = keys(pickDf.colindex)
pickDf[:N] = 1:size(pickDf, 1)
stackedPickDf = stack(pickDf, stackCols)
stackedPickDf[:variable] = map(x -> pickDf.colindex[x], stackedPickDf[:variable])
rename!(stackedPickDf, [:variable], [:M])

uberDf = convert(DataFrame, uberMatrix)
uberDf[:N] = 1:size(uberDf, 1)
stackedUberDf = stack(uberDf)
stackedUberDf[:variable] = map(x -> uberDf.colindex[x], stackedUberDf[:variable])
rename!(stackedUberDf, [:variable], [:M])

writetable("picks.csv", stackedPickDf)
writetable("uber.csv", stackedUberDf)

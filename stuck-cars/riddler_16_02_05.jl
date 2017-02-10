addprocs()

using StatsBase
import LsqFit
import Gadfly

@everywhere const nTrials = 100::Int
@everywhere const N = Vector{Int}(round(logspace(0, 5, 50)))

@everywhere function naiveClumps(a::Vector{Int})
  slowestSoFar = typemax(Int)
  clumps = 0::Int
  for i in a
    if i < slowestSoFar
      clumps += 1
      slowestSoFar = i
    end
  end
  return clumps
end

@everywhere function simulate(n::Int)
  c = zeros(Int, nTrials)
  cars = collect(1:n)
  for iTrial in 1:nTrials
    c[iTrial] = naiveClumps(shuffle!(cars))
  end
  return mean(c), StatsBase.percentile(c, [25, 50, 75])
end

clumps = SharedArray(Float64, size(N))
clumps25 = SharedArray(Float64, size(N))
clumps50 = SharedArray(Float64, size(N))
clumps75 = SharedArray(Float64, size(N))
# Shuffle N to avoid work on the high end
@everywhere const Ns = shuffle(N)
@sync @parallel for i in 1:length(Ns)
  clumps[i], pcts = simulate(Ns[i])
  clumps25[i] = pcts[1]
  clumps50[i] = pcts[2]
  clumps75[i] = pcts[3]
end

Gadfly.plot(Gadfly.Geom.point, x=Ns, y=clumps,
  Gadfly.Geom.ribbon, ymin=clumps25, ymax=clumps75,
  Gadfly.Scale.x_log)
model(x, p) = p[1] + log(x)/log(p[2])
fit = LsqFit.curve_fit(model, Ns, clumps, [1., 2.])

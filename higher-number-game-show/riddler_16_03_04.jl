using Gadfly

const NTRIALS = 10000

staystay = zeros(NTRIALS)
stayhit = zeros(NTRIALS)
hitstay = zeros(NTRIALS)
hithit = zeros(NTRIALS)

staybest = zeros(NTRIALS)
hitbest = zeros(NTRIALS)

x = rand(Float64, NTRIALS, 4)

staystay = x[x[:, 1] .> x[:, 3], 1]
stayhit = x[x[:, 1] .> x[:, 4], 1]
hitstay = x[x[:, 2] .> x[:, 3], 1]
hithit = x[x[:, 2] .> x[:, 4], 1]

staybest = x[(x[:, 1] .> x[:, 3]) & (x[:, 1] .> x[:, 4]), 1]
hitbest = x[(x[:, 2] .> x[:, 3]) & (x[:, 2] .> x[:, 4]), 1]

draw(SVG("staybest.svg", 21cm, (21/golden)cm),
  plot(x=staybest, Geom.histogram(density=true)))
draw(SVG("hitbest.svg", 21cm, (21/golden)cm),
  plot(x=hitbest, Geom.histogram(density=true)))

e, counts = hist(staybest, 100)
cumcounts = cumsum(counts ./ length(staybest))
draw(SVG("staycum.svg", 21cm, (21/golden)cm),
  plot(x=midpoints(e), y=cumcounts, Geom.step))

draw(SVG("staystay.svg", 21cm, (21/golden)cm),
  plot(x=staystay, Geom.histogram(density=true)))
draw(SVG("stayhit.svg", 21cm, (21/golden)cm),
  plot(x=stayhit, Geom.histogram(density=true)))
draw(SVG("hitstay.svg", 21cm, (21/golden)cm),
  plot(x=hitstay, Geom.histogram(density=true)))
draw(SVG("hithit.svg", 21cm, (21/golden)cm),
  plot(x=hithit, Geom.histogram(density=true)))

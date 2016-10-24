
function gem_me()
  found = [0, 0, 0]
  while any(found .== 0)
    monster = rand(1:6)
    if monster <= 3
      found[1] += 3
    elseif monster <= 5
      found[2] += 2
    else
      found[3] += 1
    end
  end
  return found
end

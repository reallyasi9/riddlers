import HTTP

const factbookurl = "https://www.cia.gov/library/publications/the-world-factbook/docs/flagsoftheworld.html"
const flagsurl = "https://www.cia.gov/library/publications/the-world-factbook/attachments/flags"

function getpage(url::AbstractString)
    r = HTTP.get(url; status_exception=true)
    String(r.body)
end

function downloadflag(abbr::AbstractString)
    io = open(joinpath(@__DIR__, "flags", abbr * ".gif"), "w")
    HTTP.get("$flagsurl/$abbr-flag.gif", response_stream=io)
    close(io)
end

function downloadflags()
    body = getpage(factbookurl)
    abbrs = Dict{Symbol, String}()
    for m ∈ eachmatch(r"<img alt=\"(.*?) Flag\" title=\".*? Flag\" src=\"../attachments/flags/(.*?)-flag.gif\" />", body)
        country, abbr = m.captures
        downloadflag(abbr)
        abbrs[Symbol(abbr)] = country
    end
    abbrs
end

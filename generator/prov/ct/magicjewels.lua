local scripts = arg[0]:match("^(.*generator[/%\\])")
dofile(scripts.."lib/makereel.lua")

local symset = {
	1, -- 1 wild (2, 3, 4 reels only)
	8, -- 2 crown    5000
	8, -- 3 ruby     1500
	4, -- 4 diamond  500
	4, -- 5 emerald  200
	3, -- 6 amber    50
	3, -- 7 sapphire 50
	3, -- 8 amethyst 40
}

local neighbours = {
	--1, 2, 3, 4, 5, 6, 7, 8,
	{ 2, 0, 0, 0, 0, 0, 0, 0,}, -- 1 wild (2, 3, 4 reels only)
	{ 0, 2, 0, 0, 0, 0, 0, 0,}, -- 2 crown
	{ 0, 0, 2, 0, 0, 0, 0, 0,}, -- 3 ruby
	{ 0, 0, 0, 2, 0, 0, 0, 0,}, -- 4 diamond
	{ 0, 0, 0, 0, 2, 0, 0, 0,}, -- 5 emerald
	{ 0, 0, 0, 0, 0, 2, 0, 0,}, -- 6 amber
	{ 0, 0, 0, 0, 0, 0, 2, 0,}, -- 7 sapphire
	{ 0, 0, 0, 0, 0, 0, 0, 2,}, -- 8 amethyst
}

local function reelgen(n)
	local ss = tcopy(symset)
	if n == 1 or n == 5 then
		ss[1] = 0
	end
	return makereel(ss, neighbours)
end

if autoscan then
	return reelgen
end

math.randomseed(os.time())
printreel(reelgen(1))
printreel(reelgen(2))
printreel(reelgen(3))
printreel(reelgen(4))
printreel(reelgen(5))

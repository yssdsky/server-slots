
local mw = 3674.547370673521
local mr = 67.344781034542
local rtp = 95 -- place here value
local Pw = (rtp-mr)/(mw - mr)
print(string.format("Pw = %s, HRw = 1/%s", Pw, 1/Pw))

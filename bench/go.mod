module github.com/unilibs/uniwidth/bench

go 1.25.1

replace github.com/unilibs/uniwidth => ../

require (
	github.com/mattn/go-runewidth v0.0.19
	github.com/unilibs/uniwidth v0.0.0-00010101000000-000000000000
)

require github.com/clipperhouse/uax29/v2 v2.2.0 // indirect

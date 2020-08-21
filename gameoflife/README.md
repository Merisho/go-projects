# John Conway's Game of Life
This is the Golang implementation of [Conway's cellular automata](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life).
It outputs everything to CLI, so the field isn't displayed as an infinite grid: cells that have greater coordinates than allowed size or less than 0 won't be displayed.
However, those cells will still be live in memory, so this limitation is only the matter of output mechanism implementation.

## How to set a configuration
`configuration.go` contains necessary guidance and means to set the initial configuration.

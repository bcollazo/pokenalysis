# pokenalysis
Go scripts to analyze Gen 1 Pokemon data.
Built to learn Golang.

## Design
We will use [PokeApi's](https://pokeapi.co/) data.

Program will consist of several graphs / data summaries.
These will be accessible via different options to a command line tool.
The commands will download required data to ~/tmp if needed, subsequent executions
will attempt to data from ~/tmp directly.
The commands will compute data on the fly.

We now present the different options / graphs.

### Type Histogram
`poke histo` yields a histogram of the different types.

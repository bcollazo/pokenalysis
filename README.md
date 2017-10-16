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

### TODO:
1. Sql CLI for Pokemon.  Full with query parser, executor, ext...
2. Concurrently-executed (using go routines) pokemon battle simulation.
3. Define better GoodRatio.
4. Filtering by generation, or by any 2 numbers.  Maybe list of generations?
5. Find best solo pokemon.
6. Make simulation be able to run in multiple machines.  Use network to implement a PAXOS-kind of protocol
and make distributed system self-healing.

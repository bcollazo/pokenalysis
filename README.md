# pokenalysis
![Pokenalysis Logo](https://raw.githubusercontent.com/bcollazo/pokenalysis/master/images/logo.png)
Command-line tool that computes statistics about Pokemon video game series data.
Things computed:
- Ocurrance of Types
- How many pokemon is a Type super-effective against?
- Ratio of pokemons a Type is good against vs bad against.
- Best Move Set for any given pokemon.

All these statistics can be computed on any combination of pokemon generations and
results can be sorted in any direction.

I built this project to learn Go.

## Design
We use [PokeApi's](https://pokeapi.co/) data.

The first time the tool is run, it will download the data concurrently to a temporary folder.
Subsequent runs will read from local folder.

The commands compute data on the fly, some use as many processors as possible.

## Usage
`pokenalysis -command=<command> -gens=<gens> -sort=<sort>`
- command: (string) one of either 'histo', 'superhisto', 'goodratio', 'bestpoke'
- gens: (string) comma-separated list of ints from 1 - 7. e.g. 1,2,5,6
- sort: (int) -1, 0 or 1.

### Future Work:
- Sql CLI for Pokemon.  Full with query parser, executor, ext...
- Make simulation be able to run in multiple machines.
Use network to implement a PAXOS-kind of protocol and make distributed system self-healing.

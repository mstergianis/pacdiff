# WIP: pacdiff

pacdiff is a golang package differ. The intent is to provide utility beyond
`diff` when comparing two similar packages; think a db model and an api model.
These are two packages that represent the state of your application, that must
be kept in sync.

## Roadmap

- [ ] diff packages by parsing the code
- [ ] provide snapshot testing by storing diffs in a readable format that can be checked into version control
- [ ] code generation that converts models between the two packages via //go:generate tags

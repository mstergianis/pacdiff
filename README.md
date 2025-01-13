# WIP: pacdiff

pacdiff is a golang package differ. The intent is to provide utility beyond
`diff` when comparing two similar packages; think a db model and an api model.
These are two packages that represent the state of your application, that must
be kept in sync.

## Roadmap

- [ ] diff packages by parsing the code
  - [x] basic diffing functionality, with JSON encoded serialization format
  - [ ] unified diff output
  - [ ] new serialization format?
    - I think the issue I'm having is that the current JSON encoding, while it's
      easy to parse, isn't super intuitive for a person to read. That's not
      necessarily an issue if the serialization is only for snapshot testing.
      But if I want the snapshot tests to be readable in their own right I want
      to keep evaluating this.
- [ ] provide snapshot testing by storing diffs in a readable format that can be checked into version control
- [ ] code generation that converts models between the two packages via //go:generate tags

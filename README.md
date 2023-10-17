# pastebin

This is a Go web service that acts as a "pastebin", similar to [the OG pastebin.com](https://pastebin.com) or [paste.debian.net](https://paste.debian.net). It saves and displays snippets of text files.

## Internals

- Badger as the embedded key-value database
- Protocol Buffers for serialization of values stored in the embedded kv database
- Go standard library for the rest, including HTTP

## License

Apache-2.0

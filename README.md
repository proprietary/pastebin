# pastebin

This is a Go web service that acts as a "pastebin", similar to [the OG pastebin.com](https://pastebin.com) or [paste.debian.net](https://paste.debian.net). It saves and displays snippets of text files.

What makes this different?

Not only can you use the web interface, but you can also create pastes from the command line! No need to copy-paste or click through to uploading files. You can create pastes from the ``bash`` or ``zsh`` command line like so (``/usr/include/assert.h`` is just an example):

```bash
cat /usr/include/assert.h | curl --data-binary @- https://paste.libhack.so
https://paste.libhack.so/4tyx4jl9
```

You can share that link, and it will display that ``/usr/include/assert.h`` file you just uploaded.

In fact, you can make a bash alias, so that you may pipe files from the command line to create a convenient pastebin

```bash
echo "alias pastebin='curl -L --data-binary @- https://paste.libhack.so'" >> ~/.bashrc
source ~/.bashrc
# And now you can pipe output like so...
echo "hello world" | pastebin
https://paste.libhack.so/lktsssj5
```

## Internals

- Badger as the embedded key-value database
- Protocol Buffers for serialization of values stored in the embedded kv database
- Go standard library for the rest, including HTTP

## License

Apache-2.0

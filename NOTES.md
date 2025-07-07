# going further

- Subdirectories! I'd intended to support this on this run-through, but turns out fsnotify is architected different than I'd expected:
    I expected creating `bar/baz.txt` to make a single event `CREATE bar/baz.txt` but instead got `CREATE bar` as a dir and, nested, `CREATE baz.txt`. This brings me to:
- Recursive implementation! should refactor into `init` function accepting a base directory, so we can recursively call it on startup with existing subdirectories, and also on `CREATE` events.
- Deletion events! This would be super simple, just create another request/response types for file deletion with just a path field
- Diffs and compression! sending large text changes and/or images is expensive, so we should implement a text diff engine on the client and patch those diffs into the existing tree on the server. This also requires sending a checksum hash with each diff, to prevent getting out of sync, and a server->client message, to request the full file when we do. Compression is much more straightforward, just use an off the shelf lib like zlib
- Versioning for text files! Diffs make this easy, we already have content hashes and efficiently storable diffs, so just build a basic git-like thing on top of diffs.
- Ratelimits / raw TCP stream! Websockets don't allow for backpressure, so if our single server is meant to handle multiple clients (which we should, since on average this won't happen) and one client sends way too much data, (bottleneck with decoding json + base64) all other clients will stall/die. either change protocol to chunk on the client and listen for server ratelimit messages, or change protocol to raw tcp stream and implement backpressure

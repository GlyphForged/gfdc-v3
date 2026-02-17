# TODO.md

Using this file to keep track of where I am and what next steps are, because I
suck at stepping away and coming back. I guess it could also be seen as a bit of
a devlog, so gonna treat it as a combo.

For this session, my goal is to get the Go server running locally on my dev
machine. Next steps as best as I can figure them are to containerize and then
test locally. GHCR looks good for container hosting, so next would be that
workflow, then add an SSH deploy step to the gfdc-v3 repo on successful pull.
Finally, we'll want to get the droplet setup, including adding GHCR creds to the
droplet and performing a PR to trigger automation. If everything is working,
should be able to point dev.glyphforged.com to the droplet and get a test page
going.

## Initial Dev Session - Feb 16

Got the basic Go server up and running. Learned a few things about how the
server runs at a fundamental level.

The key items seem to be:
- Handler
- ServeMux
- http.Server
- ListenAndServe

ServeMux is the request multiplexer (router)
`http.Dir("./static")` creates a virtual filesystem rooted in static dir
  - Server is sandboxed to this directory
`http.FileServer(...)` gives us an `http.Handler` back. This handler is what's
  saved as `fileServer`.

FileServer:
  - Reads files
  - Parses MIME types
  - Writes them to the response
  - Handles directories

`mux.Handle` seems to initialize our virtual filesystem as `/`.
  - This is why test home page is served on 8080. Essentially `URI/index.html`

`http.Server` object saved as `server`.
  - Sets 'listen on' address
  - Sets router to use (in this case our `mux` router)

Finally we `ListenAndServe()` ad infinitum unless we hit an error, in which case
we panic and log the error.

`ListenAndServe()`
- Opens a TCP listener
- Accepts connections

Best Guess Call Stack:
```
Client makes request
  |
  V
TCP Connection
  |
  V
http.Server
  |
  V
ServeMux (router)
  |
  V
Handler
  |
  V
ResponseWriter crafts HTTP response
  |
  V
Client receives response
```

That's it for today. More tomorrow.

---

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

## Day 2... 6 days later... - Feb 22

Alright, so, turns out the Flu fucking sucks, and gives one zero brain power for
dev work. Awesome. But I'm feeling just well enough now to poke around a bit at
templates, something I completely forgot about.

I've got a basic system in mind, but need to do some more digging to figure out
how to connect a PostgreSQL db to the pipeline sooner than I thought, a lot of
this page data is going to come from that db, and a lot of the router/renderer
design is on hold until I get that built.

So for now, I've got 3 template folders: layouts/, partials/, and pages/, goal
is that they do what they say on the tin. As the pages grow and I piece together
a semblance of a theme, those will be where I build each part. The renderer
should stitch them together nicely as needed, pulling data from the db where it
fits. Need to get data for some basic test pages built soon.

Additionally, I've updated the main.go topull in a separate server module.
While the binary is monolithic, I don't think the code-base should be, so I'm
starting to think about what goes where. For now, my thought is routing and
rendering are really the server's domain. Auth and such can be separate modules
built and pulled into `main.go`.

Took some time to sit down and add some clear commentary on the sections that I
found myself having to go back and re-learn. Hopefully this makes spin-up next
time faster, as half the code-base was forgotten by the time I came back. As a
working dad of 2 I really shouldn't have expected to hit this code base daily.

Server is where I will build out basic pathing, though the slug handling I
believe will need some additional logic in a helper function or two. Slowly
getting the hang of Golang methods on a struct. The syntax was a bit odd until I
got a feel for what they're going for. Still kind of prefer Rust's impl. I swear
I'm going to wind up talking myself into doing v4 in Rust...

Finally, added some future-proofing with a quick update to how we decide on
where we are serving. Should come in useful when we containerize, or could just
be one more line to refactor. Time will tell.

Doesn't feel like a ton done today, but the concepts are becoming more concrete,
and I'm only feeling more confident about my decision to break away from the "JS
all the things" stack. Seriously, this feels much simpler than the Node
ecosystem.

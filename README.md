# GlyphForged.com v3

Written by Aaron Chavez

I've decided to make my design doc my top-level readme for a few reasons. First,
it gives me a place to store my thoughts as I build the strategy/design. Second,
if anyone is weird enough to want to know why I built the site the way I did,
they get an admittedly unadvisable peek into how my gray matter works. Finally,
since this isn't a user-facing application, it doesn't make too much sense for
me to keep "usage instructions" or similar in the README. This is admittedly
equal parts manifesto and design doc. I would apologize, but if you keep reading
you've brought it on yourself.

## 1. Hashtag Goals

I'm becoming quickly disenchanted with NodeJS and the whole "JavaScript all the
things" ecosystem. This is not a moral or judgemental stance, purely a
subjective one. I understand the purpose of Node, React, Angular, etc., and I
have built a handful of sites using a handful of variations on the Node
ecosystem. I have no major complaints beyond the typical security concerns
and minor gripes. At the end of the day, this site doesn't need to be too much.

So I'm kind of going back to basics, while still scratching the "learn a new
technology" itch. I've never written anything in Go before, so why not learn
with a simple little website server? I've never written a site with HTMX before,
so why not my personal site?

LLM/AI disclosure: I am not a militant anti-AI person, but I do see some
problems with some utilizations of LLMs. I don't like the idea of passing off
my learning and/or thinking to the computer, that is not, in my humble opinion,
its purpose. To quote Frank Herbert: "What do such machines really do? They
increase the number of things we can do without thinking. Things we do without
thinking-there's the real danger." But neither do I think LLMs are an
inherently bad thing. I do see real value in their ability to lower the barrier
to information.

I am using LLMs in this project, but in a VERY limited capacity. I have never
written my own design doc from start to finish. I've asked an LLM to provide
me a general strategy for how to write this doc, and I'm taking extensive
artistic liberty with what it's given me. Furthermore, I will likely be
leveraging the LLM to help me understand my options in building the site
using this stack. As a Go neophyte, but not a brand new developer, I find value
in using LLMs to help me understand the general skeleton of what I will be
building, but I refuse to put a single line of LLM generated code in this
project. Again, that's not a moral choice, but a subjective one. I want to
learn. I don't learn by copying and pasting, I learn by writing code that almost
works, updating it until it works, breaking it when I try to change it, then
learning all the ways I messed it up and re-writing it better, but not quite
right. Lather, rinse, repeat until you have a code base you aren't completely
ashamed by.

### Explicit Goals (TL;DR)

- Wat it is
  - A clean personal CV/showcase/personal musings site
- Wat it do
  - Personal portfolio and showoff spot
  - New long form ramblings without writing more code
  - Simple server-rendered HTML
- Why it do
  - Learnin' Go and HTMX
  - Not Node
- Wen it done
  - Projects, Games, Musings sections
  - Rants updated via CRUD post-login

## 2. Anti-Goals & Restraints

- Wat it not do
  - Client-heavy JS (No frontend framework, no build)
  - Comments or social media (No mod burden)
  - Non-admin accounts (Single trusted editor)
  - Dynamic Content/SPA (No Client-side routing or state)
  - Real-time shenanigans (No websockets, polling, live updates)
- How it tied down
  - Minimal JS
  - One server binary running on one droplet
  - Smol dependency list (External imports must have a reason to exist)
- How hard it be
  - Admin team: Me, myself, y yo.
  - Smol server
  - EZ Backup

## 3. Usage

We really only need two "usage" modes for this site, but it behooves me to stop
and think about this. I should keep in mind and put in writing somewhere a clear
line in the sand when it comes to the usage of the site.

So I've boiled it down to two states: Public & Admin.

- Public
  - Unauthenticated, read-only access. No data manipulation of any kind.
- Admin (Yours truly)
  - Authenticated. Control over the ravings. CRUD access to blog.

## 4. Architecture

### System Overview

- Single monolithic Go server.
- Server renders all HTML.
- HTMX is used to request partial HTML from the server.
- Server handles two kinds of requests:
  - Full page requests
  - Fragment requests

### Major Components

- Request handling
  - Routes incoming HTTP requests and decides what logic runs.
  - Owns every request until it is fully resolved.
- Rendering
  - Owns templates, layouts, and HTML fragments.
  - Provides the rendered page to be shown to the client.
- Content managements
  - Owns blog posts, projects, and read/write rules.
  - Effectively the database which the renderer calls from.
  - SQLite - Embedded DB, no need for a container.
- Authentication
  - Owns admin login state and access checks.
  - Simple "is this request legit?"" check.
- Static assets
  - Serves CSS, images, and WASM game files.

### Request Boundaries

- Public routes are strictly read-only, no auth.
- State mutation is restricted to authenticated requests via admin pages.
- Authentication is enforced before any write logic runs.
- WASM treated as static assets embedded in server-rendered pages.

## 5. Tech Stack & Y Tho

### Backend: Go

A Go server will compile to a single executable binary, minimizing complexity
and dependency. Allows for my goal of one server on one droplet with easy
backups/rollbacks/etc. Additionally, this will help enforce my server rendered
requirement. Go has a strong standard library for HTTP which allows for minimal
dependencies. Finally, using Go will force me to think about things that have
been historically abstracted away in prior web dev projects.

### Frontend: HTMX

The language of the web is HTML. To that end, where partial page changes are
required, I will be using HTMX. This ensures I do not need to build a dedicated
front-end application, I can simply query the server for the updated sections
of the page on an as-needed basis. No jQuery fragility, no clunky reloads, and
no overly complicated SPAs. HTMX is not a core dependency, but a nice-to-have
extension to my architecture. Core site functionality will not depend on it, and
if for any reason HTMX is unavailable or removed, the site will still function
albeit with more full-page reloads. HTMX is currently planned primarily during
admin workflows, though it may be used to reduce unnecessary navigation on the
user side, eventually. HTMX will NOT introduce user-side state.

### The Rejects

Note: None of this is objective, purely a best-fit analysis. Most of my web-dev
experience is in SPAs and purely static sites.

- Node.js + SPA (React/Angular/etc)
  - Unnecessary complexity
  - Too many sec vulns for a simple CV site
- Purely Static Site
  - Not enough flexibility
  - Not future-proof

## 6. Breadcrumbs (Site Structure)

### Top-Level

- Landing Page
- Projects
- Games
- Blog/Musings
- About/CV

### URL Shape

- Projects, games, and musing are listed on dedicated index pages.
- Individual blog posts/games/projects are accessed via stable slugs.
  - i.e. `/projects/{slug}`, `/games/{slug}`, `/musings/{slug}`.
- Admin access via subdomain.
  - subdomain.glyphforged.com
  - Any unauthorized access to this will lead to the login page.
  - Admin subdomain not linked from public pages and intentionally unindexed.
- Large projects may eventually also live in subdomains.
  - i.e. f1.glyphforged.com, pendejos.glyphforged.com, etc..
  - These are independent deployments
  - Not part of the main site's routing or content.
  - If it doesn't "fit" on the main site, it gets a subdomain.

## 7. Data Structures

- Only `Published` items are visible to public requests.
- `Draft` and `Archived` items return 404 to public users.
- All fields are mutable unless otherwise noted.

### Musings

Wat it is: A blog post with one of my ramblings.

Wat it need:

- Title
- Slug(s)
- Content
- Creation Date - Immutable
- Published At - Immutable
- Updated Date (Optional) (most recent change only)
- Tags
- Status (Published, Draft, Archived)
  - Visibility:
    - Published - Public
    - Draft & Archived - Admin-only. Public view renders 404.

### Project/Game (Showcase)

Wat it is: A page showcasing a project or group of projects.

Wat it need:

- Type (Enum - Game, Project, Etc)
- Title
- Slug(s)
- Description/Blurb
- Creation Date - Immutable
- Updated Date
- Tags
- HostingType (Enum)
- URI (WASMlocation, itch.io URL, etc.)
- Status (Published, Draft, Archived)
  - Visibility:
    - Published - Public
    - Draft & Archived - Admin-only. Public view renders 404.

### AdminUser - Private

Wat it is: An admin user. Should just be me.

Wat it need:

- Username - Immutable
- Name - Immutable
- PasswordHash
- Created date - Immutable
- LastLogin

## 8. Bouncer (Auth Model)

### Auth Model

Multi-admin capable, but flat. Any authorized user will be an admin, if you
aren't an admin, you're public. Login is a simple User + PW Hash (salted, of
course). No UI reset flow is needed, recovery will be handled via a CLI
subcommand built into the main binary. No HTTP-based recovery.

### Session Model

Server-side session store with sessions stored in memory. This allows for
easily revoked sessions with explicit control. Sliding expiration (inactivity)
with manual invalidation on logout. Expiration is derived from inactivity
duration on interaction. Expiration happens @ 1-hour since last activity. Call
me paranoid, but if passwords are reset, all sessions should immediately be
invalidated by clearing the session store.

All admin routes require a valid session before handler logic executes. This
check will be performed for each interaction that may lead to state mutation.
Templates do not perform authorization checks. (Auth before render).

Each session stores:

- userID
- createdAt
- lastActivityAt

On every authed request:

- Look up session by ID
- If not found, reject
- If expred, delete session, reject.
- If valid:
  - Update lastActivityAt
  - Proceed

### Da Cookies

- SessionID stored in HttpOnly cookie.
- Cookie marked secure (HTTPS only).
- SameSite policy = Lax to mitigate CSRF
- Cookie scoped to admin subdomain only.
  - (If admin lives on admin.gf.com, cookie is not valid for gf.com)

## 9. HTMXWTFBBQ (HTMX Patterns)

### Where it's used

- Admin CRUD flows (create, edit, delete musings)
- Preview functionality (blog at first, projects/games later?)
- That's really it for now.
- Fragment endpoints must return meaningful HTML when accessed directly.
- Fragment enpoints must not rely on client-side state.

### Responses

- Full page responses for:
  - Public navigation
  - Initial Loads
  - Non-HTMX clients
- Fragments for:
  - HTMX-triggered admin actions/state mutations.
- Note: Fragment endpoints never return JSON.

### Epic Fails

- When an HTMX request fails auth:
  - Treat as unauthenticated, no matter the reason
  - Unauthorized Admin route > Redirect to Login
  - Public access to draft/archive > 404 page
- Validation failures > Return HTML fragment with errors.
- Log all unauthorized admin access attempts

### The State of State

- The client is never an authority.
- Server is the sole source of truth.
- HTMX does NOT introduce application state.
- All mutations occur via normal HTTP semantics.

## 10. The Soapbox (Blog CRUD Flow)

CUD is all admin-only. Public cannot take CUD actions.

### Create

- New musings created with:
  - Status = Draft
  - CreationDate set (immutable after set)
  - NULL PublishedAt value
- After creation, redirect to edit page
- Draft may have initially empty fields
- Title & Content required to publish

### Read

- Public Read
  - Only published
  - Archived/Draft -> 404
- Admin Read
  - All stauses accessible

### Update

#### Edit

- Edit returns HTML fragments
- Validation failures return form with errors
- Successful updates return updated form fragment w/ success indicator
- If Status == Published
  - Updates affect live version immediately
  - UpdatedDate is updated
- If Status == Draft || Archived
  - Status remains unchanged
- Slug changes
  - Old slug added to alias list
  - 301 redirect maintained

#### Publish

- Requires Status == Draft || Archived
- Successful publish returns updated form fragment w/ success indicator
- On Publish:
  - Status == Published
  - PublishedAt == dateTime of first publish (Immutable)
- If Archive being re-published, update UpdatedAt

### Delete (Archive, really)

- Status = Archived
- Public requests now return 404
- Slug aliases resolve to canonical slug
- Canonical slug returns 404 if archived

### Preview

- Preview will never change state, despite being Admin only
- Preview will render draft content via a fragment

## 11. The Fun Zone (WASM Games and Project Strats)

Both projects and games are unified under a single identity: Showcase

### Rendering Model

Projects and games share a base layout template.
Every project/game should:

- Be server-rendered HTML
- Include
  - Title
  - Description/Blurb
  - Tags
  - Embedded project/game section

### Hosting Model

- The server is only responsible for delivering static files and HTML.
- Server does NOT handle game runtime behavior.

- HostingType = Enum [External, SelfHosted]
- If HostingType == External
  - Render embedded iframe/external link
  - Server does not proxy content
  - Server does not fetch remote content
- If HostingType == SelfHosted
  - WASM + JS glue served as static assets
  - Embedded directly in page
  - Server does not execute game logic

### Static Asset Strategy

- Static directory holds all static content, from images to WASM
- WASM files must be served with correct MIME type: `application/wasm`
- Static assets are read-only and not modified at runtime
- Each project/game has its own subfolder within the `/static/` directory
  - i.e. `/static/showcase/{slug}/`

### Isolation Boundaries

- Subdomains (i.e. larger projects) are separate. If we link to a subdomained
  project we are sending the client to a different server
- Self-hosted WASM runs in the browser
- The server should NEVER
  - Trust anything from the WASM runtime
  - Accept POSTS from game runtime to admin endpoints
- Game runtimes have no privileged communication with the server beyond public
  HTTP endpoints

## 12. Where It Lives (Deployment & ENV)

- "Cattle over pets"
  - Everything is containerized
  - Nothing critical lives exclusively on the host
- Fully containerized via Docker
  - Caddy (Reverse proxy + Auto-TLS)
  - Go Server App
  - SQLite DB
  - App and DB are internal-only
    - Communication via Docker network
    - DB is never publicly exposed
- Docker Volumes for storage
  - `static-assets` -> self-hosted static content
  - `app-data` - SQLite file
- Disposable Digital Ocean droplet (unless I cave and toss this on AWS)
- Docker compose as single source of truth

### Static Content

- Static content lives in a separate GH repo (for now)
- App container mounts the `static-assets` Docker volume on boot
- Static files served directly from the mounted volume
- CI pipeline
  - Build static artifacts where applicable
  - SSH into droplet
  - Sync files into `static-assets` volume

### App CI/CD

On push to main:
1. Build Docker image
2. Run tests
3. Push image to container registry
4. SSH into droplet
5. Pull latest image
6. Run `docker compuse up -d`

### Data Management

<!-- TODO: New Data Management Section for SQLIte -->

### MI Process

- In the event of an unrecoverable failure
  - Provision new droplet
  - Install Docker
  - Clone main site repo
  - Run `docker compose up -d`
  - Restore db from backup
  - Re-sync static content repo

## 13. Spying on myself (Observability)

I tend to over-engineer this kind of shit, so in the effort of reminding myself
to KISS, gonna log observability goals here.

The system needs to answer:
- Is the site up?
- Is the db reachable?
- Are requests succeeding?
- Any weirdness going on? (4xx/5xx spikes, etc.)

### Logging Strat

Application Logs

- Go app outputs structured JSON logs to stdout including:
  - Timestamp
  - Request ID
  - Method
  - Path
  - Response status
  - Latency
  - Error details

- Caddy logs include:
  - Enable access logs
  - Enable error logs

`docker logs` to access.

### Checkups

App should expose `/health`

Page verifies:
- App is running
- db connection is reachable
- 200 if healthy, 5xx if dependencies unavailable
- Exposed by Caddy

Use DigitalOcean tooling for system metrics monitoring and alerting

Use host tooling for updtime monitoring and alerting

## 14. OpSec (Security Considerations)

The goal here is to reduce the potential attack surface and enforce some basic
trust boundaries. It's just a CV site, but script kiddies gonna script kiddie.

### Transport Security

- All traffic over HTTPS.
- Caddy handles TLS provisioning
- HTTP redirects to HTTPS

### Network Isolation

- Caddy responsible for exposing public ports
- Go app and DB containers remain internal only
- DB is NEVER publicly accessible
- Inter-container comm happens via Docker network

### Auth and Session Security

- Server-side session store with in-memory sessions
- Session IDs SHA generated
- Session IDs stores in:
  - HttpOnly cookie
  - Secure cookie
  - SameSite=Lax
- Sliding expiration (1-hour)
- Invalid session treated as unauthorized
- Logout explicitly invalidates the session
- Password reset clears all active sessions
- Admin routes demand authenticated session 
- Auth checks BEFORE handler logic
- Public access to admin areas returns 404

### Defensive lines

- Basic rate limiting enforced at the proxy layer
  - Login endpoint rate-limited to mitigate brute-force attempts
  - General request rate limits to prevent abuse and resource exhaustion
  - IP throttling on excessive failed attempts

### Bobby Tables

- State-mutating input is validated server-side
- Client-side validation is purely UX enhancement
- Slug changes are sanitized and validated before persistence

### Failure & Breach Response

In case of suspected compromise:

1. Rotate admin credentials.
2. Clear all active sessions.
3. Rebuild and redeploy containers.
4. Restore database from known-good backup if necessary.
5. Review logs for suspicious activity.

## 15. Wanna Do These Later

- CRUD updates for projects and games
- Cloud db backup
- Log aggregation?
- Enable HSTS

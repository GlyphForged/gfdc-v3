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
with a simple little website server? I've never written a site in HTMX before,
so why not my personal site?

LLM/AI disclosure: I am not a militant anti-AI person, but I do see some
problems with some utilizations of LLMs. I don't like the idea of passing off
my learning and/or thinking to the computer, that is not, in my humble opinion,
its purpose. To quote Frank Herbert: “What do such machines really do? They
increase the number of things we can do without thinking. Things we do without
thinking-there’s the real danger.” But neither do I think LLMs are an
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
- Published At
- Updated Date (Optional) (most recent change only)
- Tags
- Status (Published, Draft, Archived)
  - Visibility:
    - Published - Public
    - Draft & Archived - Admin-only. Public view renders 404.

### Project/Game

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
- PasswordHash
- Created date - Immutable
- LastLogin

## Bouncer (Auth Model)

WIP

## HTMXWTFBBQ (HTMX Patterns)

WIP

## The Soapbox (Blog CRUD Flow)

WIP

## The Fun Zone (WASM Games and Project Strats)

WIP

## Where It Lives (Deployment & ENV)

WIP

## OpSec (Security Considerations)

WIP

## Wanna Do These Later

- CRUD updates for projects and games

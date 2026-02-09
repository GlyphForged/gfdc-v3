# GlyphForged.com v3

## Hashtag Goals

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

## Explicit Goals (TL;DR)

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

## Anti-Goals & Restraints

- Wat it not do
  - Client-heavy JS (No frontend framework, no build)
  - Comments or social media (No mod burden)
  - Non-admin accounts (Single trusted editor)
  - Dynamic Content/SPA (No Client-side routing or state)
  - Real-time shenanigans (No websockets, polling, live updates)
- How it tied up
  - Minimal JS
  - One server binary running on one droplet
  - Smol dependency list (External imports must have a reason to exist)
- How hard it be
  - Admin team: Me, myself, y yo.
  - Smol server
  - EZ Backup

## Usage

We really only need two "usage" modes for this site, but it behooves me to stop
and think about this. I should keep in mind and put in writing somewhere a clear
line in the sand when it comes to the usage of the site.

So I've boiled it down to two states: Public & Admin.

- Public
  - Unauthenticated, read-only access. No data manipulation of any kind.
- Admin (Yours truly)
  - Authenticated. Control over the ravings. CRUD access to blog.

## Architecture

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

## Tech Stack & Y Tho

WIP

## Breadcrumbs (Site Structure)

WIP

## Data Structures

WIP

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

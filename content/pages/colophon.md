---
title: Colophon
slug: colophon
subtitle: How this website is put together.
---

This is a small personal website built to stay simple, fast, and easy to maintain.
I'm hoping to keep it that way.

## Stack

The site is generated with a custom static site generator written in Go.
Pages are written in Markdown, rendered to static HTML, and served by a tiny Go HTTP server.

Production runs on a VPS with Docker Compose.
Caddy sits in front of the Go server and handles HTTPS automatically.

## Design

The design is intentionally minimal:

- plain HTML and CSS
- no JavaScript by default
- no analytics
- no cookies
- no tracking

The font is JetBrains Mono.
The color palette is dark green because I like it.

## Hosting

The site is hosted on a small VPS.
DNS points directly to the server, and Caddy manages TLS certificates through Let's Encrypt.

## Source

The source code for this site lives on GitHub:

[github.com/shinjustin/dotcom](https://github.com/shinjustin/dotcom)

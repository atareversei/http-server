# http-server: A Minimal HTTP/1.1 Server in Go

`http-server` is a custom-built HTTP/1.1 server written in Go, implemented from the ground up using raw TCP sockets. Inspired by Go’s standard `net/http` package, this project explores the internals of networking, HTTP protocol parsing, and secure server architecture.

Originally developed as part of a university networking course (IAUT, September 2024), it has since evolved into a deeper systems project aimed at understanding and replicating core web server behaviors in a clean, modular, and idiomatic Go structure.

## Features

- Raw TCP server with connection management
- Full HTTP/1.1 request parsing (methods, headers, body)
- Custom response generation with `ResponseWriter`
- Handler interface and minimal router (like `http.ServeMux`)
- Concurrent request handling via goroutines

## Package Structure

```md
http-server/
├── net/ # Low-level TCP server: listening, connections
├── http/ # HTTP protocol: request, response, handlers
├── internal/ # Internal utilities: CLI utilities
└── examples/ # Example servers using the http-server API
```

## Learning Objectives

This project is intended to:

- Understand how web servers work beneath abstractions
- Explore the full HTTP/1.1 specification (RFC 7230+)
- Study security risks like header injection, DoS, and slowloris
- Mimic the structure and ergonomics of Go's standard library
- Build reusable network layers from scratch

## Author

Ata — [github.com/atareversei](https://github.com/atareversei)

Originally developed at Islamic Azad University of Tabriz (IAUT) for the Computer Networks course, Fall 2024.

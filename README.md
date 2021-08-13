---
marp: true
theme: gaia
_class: invert
title: go flexible
---

# Flex

A collection of packages for building Go services.

---
# Goals

- Provide `flex`ible primitives to build upon.
- Solve common issues for 12 factor apps.
- Be a good citizen in the Docker and Kubernetes ecosystem.

---

# In Practice

- The 12 factors help define some standard behaviours and expectations.
- As the principles are broad, each language may implement differently.

---
# Primitive: Runner

A `Runner` represents anything which can "run" itself.
For example, an `HTTP` server.

```go
// Runner represents the behaviour for running a service worker.
type Runner interface {
        // Run should run start processing the worker and be a blocking operation.
        Run(context.Context) error
}
```

---

# Primitive: Halter

A `Halter` represents anything which can "halt" itself.
For example a `Kafka` broker.

```go
// Halter represents the behaviour for stopping a service worker.
type Halter interface {
        // Halt should tell the worker to stop doing work.
        Halt(context.Context) error
}
```

---

# Primitive: Worker

A `Worker` represents anything that can both "run" and "halt" itself.
For example, a `cron` job.

```go
// Worker represents the behaviour for a service worker.
type Worker interface {
        Runner
        Halter
}
```

---

# Optics

- Consider `io.Reader` and `io.Writer`.
- Both single method interfaces.
- Both result in incredibly varried uses.

The goal of `flex` is to provide a similar experience, for building `k8s` native applications.

<!-- Creating strong yet simple primitives, like io.Reader & io.Writer means 
opening yourself to an incredible variety of uses-cases, this in turn means
having the flexibility you need to build complex applications. -->
---
# The `flexhttp` plugin

```go
type Server struct{ *http.Server }

func NewHTTPServer(s *http.Server) *Server {
        return &Server{Server: s}
}

func (s *Server) Run(_ context.Context) error {
        log.Printf("serving on: http://localhost%s\n", s.Addr)
        return s.ListenAndServe()
}

func (s *Server) Halt(ctx context.Context) error {
        return s.Shutdown(ctx)
}
```

--- 
# The `flexmetrics` plugin

`flexmetrics` exposes prometheus and pprof metrics on an http server.
It's `Worker` is implemented by the concrete `Server` type.

```go
type Server struct {
        Path   string
        Server *http.Server
}
```
---

# Simple yet complete

```go
// Run will start the metrics server.
func (s *Server) Run(_ context.Context) error {
        /* abbreviated... */
        mux := http.NewServeMux()
        mux.Handle(s.Path, promhttp.Handler())
        mux.HandleFunc("/debug/pprof/", pprof.Index)
        mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
        mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
        mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
        mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
        /* abbreviated... */
}
```

---
# More than meets the eye

- Environment variables for configuration
- Configurable port-bindings
- Configurable http server
- Sane defaults if you choose to use them
---

# Live demo

## Goals:

- Build a non-trivial order system (think e-commerce).
- Demonstrate how `flex` eliminates a lot of complicated setup.
- Provide working example usecases for future reference.
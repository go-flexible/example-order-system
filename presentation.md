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

- Two of the most powerful interfaces in Go are `io.Reader` and `io.Writer`.
- Both are single method interfaces, which enable incredibly varried uses.
- `ReadWriter` is their union.

The goal of `flex` is to provide a similar experience, for building `k8s` native applications.
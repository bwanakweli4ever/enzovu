Here’s a better-formatted version of your `enzovu` framework documentation:

---

# Enzovu - Exotic Go Framework

**Enzovu** is a robust and elegant Go framework inspired by the strength and wisdom of elephants. It is designed to build powerful web applications with grace and efficiency.

This framework follows a simplified MVC (Model-View-Controller) pattern, making it easy to handle HTTP requests, manage data models, and render views. It is memory-safe, scalable, and built with Go's powerful features.

---

## Features

- **Memory-Safe Architecture**: Built with Go’s strong type system and memory safety guarantees.
- **Smart Routing**: Intelligent request handling with minimal overhead and maximum flexibility.
- **Scalable by Design**: From small projects to enterprise applications, scale with confidence.
- **MVC Architecture**: Supports the MVC pattern to separate concerns in your application.

---

## Installation

### Prerequisites

- Install **Go** (version 1.16 or later) from [Go's official website](https://golang.org/dl/).
- Set up your **Go workspace** (ensure `GOPATH` is properly configured).

### Setup

1. Clone the repository:

```bash
git clone https://github.com/bwanakweli4ever/enzovu.git
cd enzovu
```

---

## Using the `go-craft` CLI Tool

The `go-craft` CLI tool helps you create resources like models, controllers, and more for your application.

### 1. Create a Model

Use the following command to create a model (e.g., `User`):

```bash
go run cmd/go-craft.go create model User
```

### 2. Create a Controller

Use the following command to create a controller (e.g., `User`):

```bash
go run cmd/go-craft.go create controller User
```

### 3. Create a Migration File

Use the following command to create a migration file (e.g., `create_users_table`):

```bash
go run cmd/go-craft.go create migration create_users_table
```

---


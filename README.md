# Gate API Emulator

Backend service written in Go that emulates an external Gate API for managing WhatsApp channel subscriptions. This service is intended for development, integration testing, and load testing scenarios where the real external API is unavailable or unsuitable.

---

# Overview

Gate API Emulator simulates an external subscription gateway, allowing client systems to create and manage subscribers and subscriptions without connecting to the real provider.
Scenario:
* Create a new subscriber
* Update an existing subscriber
* Delete a subscriber
* Create a new subscription
* Client starts to get qr code a couple of times.
* User "scan" qr code - channel become "active"
* User does not "scan" qr code - channel become "qridle"
* Update an existing subscription
* Delete a subscription
* Client read chat history, he gets 3..5 dialogs with messages
* Amulationg when the messenger sends chat message
* Client posts message to channel
* Client reads the message

The service persists data in SQLite and can run locally or inside Docker.

This project is useful for:

* Integration testing
* Backend development
* Load testing
* API contract validation
* Offline development

---

# Features

* REST API for subscriber and subscription management
* SQLite database storage
* Docker support
* Environment-based configuration
* Persistent storage option
* Load testing utilities included
* Lightweight and fast

---

# Tech Stack

* Go
* SQLite
* Docker
* Makefile
* REST API

Project structure follows standard Go layout:

```
cmd/
internal/
config/
```

---

# Getting Started

## Run locally

Clone repository:

```
git clone https://github.com/antnzr/gatemulator.git
cd gatemulator
```

Run:

```
go run cmd/gatemulator/main.go
```

---

## Run with Docker

Build image:

```
docker build -t gatemulator .
```

Run container:

```
docker run -p 8080:8080 gatemulator
```

---

# Database persistence (IMPORTANT)

To persist data when running with Docker, create an empty database file in project root:

```
touch gatemulator.db
```

If this file does not exist, the application will use in-memory database and data will be lost after shutdown.

Currently stored entities:

* subscribers
* subscriptions

---

# Configuration

Environment variables are defined in:

```
env.example
```

Example:

```
PORT=8080
DB_PATH=gatemulator.db
```

---

# Project Structure

```
cmd/gatemulator      → application entry point
internal/            → business logic
config/              → configuration
test/load/           → load testing tools
```

---

# Development

Build:

```
make build
```

Run:

```
make run
```

---

# Use Case

This emulator is designed to replace external Gate API during development and testing of systems that depend on WhatsApp subscription services.

---

# License

Private / Internal use

---

# Author

Anton Nazarenko
GitHub: https://github.com/antnzr

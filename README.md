# Dideban 👁️

**Private Monitoring Guardian**

Dideban is a lightweight, fast, and self‑hosted monitoring system built for private infrastructures, VPCs, and production‑grade web applications.

Inspired by tools like Uptime Kuma, but designed to be:

* **More lightweight**
* **More extensible**
* **Production‑first**
* **Private by default**

---

## ✨ Features (MVP)

* 🔍 HTTP / HTTPS monitoring
* 📡 Ping (ICMP) checks
* 🖥️ Resource monitoring via lightweight agent
  * CPU usage & load
  * Memory usage
  * Disk usage
  * ⏱️ Metric collection latency tracking
* ⏱️ Fast and low‑overhead scheduler
* 🟢 Real‑time status dashboard
* 🚨 Alerting (Telegram – Bale – MVP)
* 🗄️ SQLite storage
* 📦 Single Go binary deployment
* 🔐 Fully self‑hosted (no external SaaS dependency)

---

## 🧠 Philosophy

Dideban is built with these principles in mind:

* **Minimal resource usage** (CPU & RAM)
* **Fast startup and execution**
* **Clear separation between core engine and UI**
* **Configurable & extensible by design**
* **No unnecessary abstractions**

---

## 🏗️ Architecture Overview

```
+-------------------------------+
|        Dideban Process        |
|                               |
|  +-------------------------+  |
|  |   Svelte Web UI (static)|  |
|  +-----------▲-------------+  |
|              |                |
|  +-----------+-------------+  |
|  |     Gin HTTP API        |  |
|  +-----------▲-------------+  |
|              |                |
|  +-----------+-------------+  |
|  |   Core Engine           |  |
|  |  - Scheduler            |  |
|  |  - Checks               |  |
|  |  - Alerts               |  |
|  +-----------▲-------------+  |
|              |                |
|  +-----------+-------------+  |
|  |   SQLite Storage        |  |
|  +-------------------------+  |
+-------------------------------+
```

> A reverse proxy or Docker is optional and not required for normal operation.

```
------------------+       +------------------+
|   Web Dashboard  | <-->  |   Go HTTP API    |
|   (Svelte)       |       |   (Core Engine)  |
+------------------+       +------------------+
|
v
+------------------+
|   SQLite Storage |
+------------------+
```

> Future versions may include an optional **Agent** for system‑level metrics.

---

## 🚀 Getting Started

Dideban is designed to run as a **single lightweight Go binary** by default.

No external web server (Nginx, Caddy, etc.) is required.
Docker is **optional** and provided only for convenience.

---

### Default Run Mode (Recommended)

Download the pre-built binary from GitHub Releases and run:

```bash
./dideban --config /etc/dideban/config.yaml
```

This will start:

* Core monitoring engine
* Scheduler
* Embedded HTTP API
* Embedded Web UI

All in a single process.

---

### Configuration

Dideban uses a **YAML configuration file** as the primary configuration source.
Environment variables can be used to override values.

Example:

```yaml
server:
  addr: ":8080"

storage:
  path: /var/lib/dideban/db.sqlite

alert:
  telegram:
    enabled: true
    token: "BOT_TOKEN"
    chat_id: "CHAT_ID"
```

---

### Development Mode (Backend)

For development purposes:

```bash
git clone https://github.com/MrYazdan/dideban.git
cd dideban/backend

go mod tidy
go run ./cmd/dideban
```

API will be available at:

```
http://localhost:8080
```

---

### Development Mode (Frontend)

Frontend is built with Svelte.

```bash
cd frontend
npm install
npm run dev
```

Dev UI:

```
http://localhost:5173
```

> In production, the frontend is built and **served directly by the Go binary**.

---

## 📦 Project Structure (Initial)

```
dideban/
├── backend/
│   ├── cmd/
│   │   └── dideban/
│   │       └── main.go
│   ├── internal/
│   │   ├── core/        # Engine & scheduler
│   │   ├── checks/      # HTTP, Ping, etc.
│   │   ├── alert/       # Alert dispatchers
│   │   ├── storage/     # SQLite implementation
│   │   └── api/         # Gin HTTP API
│   └── go.mod
│
├── frontend/
│   ├── src/
│   └── package.json
│
├── docs/
├── .gitignore
└── README.md
```

---

## 🔔 Alerting (MVP)

Supported in MVP:

* Telegram Bot notifications
* Bale.ai Bot notifications

Planned:

* Email
* Webhook
* Script execution
* Alert grouping & throttling

---

## 🐳 Docker (Optional)

Docker is provided as a convenience for CI/CD and container-based environments.

```bash
docker run -p 8080:8080 dideban/dideban:latest
```

> Docker is **not required** and is not the default execution method.

---

## 🛣️ Roadmap

### v0.1 – MVP

* [ ] HTTP checks
* [ ] Ping checks
* [ ] SQLite storage
* [ ] Simple UI
* [ ] Telegram alerts
* [ ] Bale.ai alerts

### v0.2

* [ ] Authentication
* [ ] Multi‑user support
* [ ] Configurable retention
* [ ] Status page

### v0.3

* [ ] Agent (system metrics)
* [ ] Docker monitoring
* [ ] Plugin system

---

## 🤝 Contributing

Contributions are welcome!

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/awesome`)
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

---

## 📄 License

MIT License

---

## ❤️ Name Origin

**Dideban (دیدبان)** means *Watcher / Guardian* in Persian — a silent observer that keeps your systems safe.

---

## ⭐ Star the Project

If you like the idea, consider giving the repo a star ⭐

It helps the project grow and stay motivated.

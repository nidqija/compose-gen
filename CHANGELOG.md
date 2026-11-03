# Compose Gen - Services Changelog

This document tracks all services added or updated in **Compose Gen** based on the Git commit history, categorized by release version milestones.

---

## Version 1.5.0 (`feature/cgen-6/add-grafana-k6`)

- **Release Date:** November 3, 2026

### Grafana-K6 Load Testing & Interactive Bubble Tea TUI

#### Grafana-K6 Load Testing Service

- **Service Name:** Grafana-K6 Load Testing (`grafana-k6`)
- **Template File:** `templates/grafana-k6.yml`
- **Config Template:** `templates/configs/load-test.js`
- **Image:** `grafana/k6:latest`
- **Container Name:** `k6_load_test`
- **Exposed Ports:**
  - `5665:5665` (K6 Built-in Real-Time Web Dashboard)
- **Environment Variables:**
  - `API_URL=https://test.k6.io`
  - `K6_WEB_DASHBOARD=true`
  - `K6_WEB_DASHBOARD_PORT=5665`
  - `K6_WEB_DASHBOARD_HOST=0.0.0.0`
- **Volume Mappings:**
  - `./tests/k6:/scripts`
- **Entrypoint:** `["k6", "run", "/scripts/load-test.js"]`
- **Description:**
  Introduced Grafana-K6 load and performance testing service with built-in real-time web dashboard visualization. Includes an automated scaffolding flow: when selected, `compose-gen` automatically scaffolds a starter test script at `tests/k6/load-test.js` configured with virtual users, test durations, and threshold validations against `https://test.k6.io`. Enables the container to run successfully out-of-the-box on the very first try without crashing due to missing script mounts.

#### Interactive Terminal UI (Charm Bubble Tea & Lip Gloss)

- **Bubble Tea TUI (`tui/`):** Re-architected the interactive CLI workflow with a modern, cyberpunk-inspired terminal user interface powered by Charm's `bubbletea`, `lipgloss`, and `bubbles`.
- **Big ASCII Art Title:** Added a dynamic cyber slant `COMPOSE-GEN` ASCII title banner featuring a vertical neon gradient (Electric Cyan to Neon Violet), with automatic responsive fallback to a compact 52-column banner on narrower terminal windows.
- **Interactive Multi-Select Checklist:** Keyboard navigation (`↑/↓` or `j/k`, `Space`/`x` to toggle, `a` for select all, `Enter` to generate) with selected item counters.
- **Live Service Inspector Card:** Real-time side panel previewing image tags, port mappings, volumes, environment variables, and auto-scaffolded files for the currently hovered service.
- **Merge & Overwrite Modal:** Safe conflict resolution dialog when `docker-compose.yml` already exists, allowing users to choose between `Merge Services`, `Overwrite Completely`, or `Cancel`.
- **Post-Generation Celebration Screen:** Detailed summary displaying activated services, scaffolded config paths, preserved files, and next-step commands (`compose-gen up`, `ps`, `logs`, `down`).
- **Backward Compatibility:** Preserved legacy plain prompt via `compose-gen --plain` or `--legacy` flags and updated command help output.

---

## Version 1.4.0 (`feature/cgen-5/added-cloudflare-tunnel`)

- **Release Date:** October 24, 2026
- **Relevant Commits:** `d6d09b3`, `c6b9c27`, `f75f592`, `c8330bb`

### Cloudflare Tunnel

- **Service Name:** Cloudflare Tunnel (`cloudflared`)
- **Template File:** `templates/cloudflared.yml`
- **Image:** `cloudflare/cloudflared:latest`
- **Container Name:** `standalone_cloudflared`
- **Command:** `tunnel --no-autoupdate --logfile /app/logs/tunnel.log --url http://host.docker.internal:8025`
- **Volume Mappings:**
  - `./logs:/app/logs`
- **Extra Hosts:**
  - `host.docker.internal:host-gateway`
- **Restart Policy:** `unless-stopped`
- **Description:**
  Added Cloudflare Tunnel (`cloudflared`) service to expose local development services securely to the internet without opening inbound firewall ports or configuring router port forwarding. Configured with `host.docker.internal:host-gateway` for host network resolution to forward traffic directly to local development services (such as the Mailpit web UI on port 8025). Provides instant HTTPS and public domain connectivity for external clients, mobile app testing, and third-party webhooks.
- **Dynamic URL Extraction & Synchronization:**
  - **Log Output:** Configured `templates/cloudflared.yml` with `--logfile /app/logs/tunnel.log` mounted to `./logs:/app/logs` so runtime tunnel logs are accessible to the host.
  - **Automatic URL Parsing (`FetchTunnelURL`):** Implemented `FetchTunnelURL()` in `features/commands.go` to scan `logs/tunnel.log` with a 10-second timeout, extracting the generated `https://[a-zA-Z0-9-]+\.trycloudflare\.com` URL via regex.
  - **Environment Configuration:** Automatically creates or updates `APP_BASE_URL` in `configs/.env` with the active tunnel URL.
  - **CLI Lifecycle Integration:** Integrated `FetchTunnelURL()` into `compose-gen up` in `main.go` to automatically discover and write the tunnel URL upon service startup.

---

## Version 1.3.0 (`feat/cgen-3/add-otel-feature`)

- **Release Date:** October 13, 2026
- **Relevant Commits:** `bf5fb5d`, `64bc6b4`

### OpenTelemetry Collector

- **Service Name:** OpenTelemetry Collector (`otel`)
- **Template File:** `templates/otel.yml`
- **Config Template:** `templates/configs/otel-collector-config.yaml`
- **Image:** `otel/opentelemetry-collector-contrib:latest`
- **Command:** `["--config=/etc/otelcol-contrib/config.yaml"]`
- **Exposed Ports:**
  - `127.0.0.1:4317:4317` (OTLP gRPC Receiver)
  - `127.0.0.1:4318:4318` (OTLP HTTP Receiver)
  - `127.0.0.1:8888:8888` (Internal Collector Prometheus Metrics)
  - `127.0.0.1:13133:13133` (Health Check Extension)
- **Description:**
  Introduced OpenTelemetry Collector Contrib service for unified collection, processing, and exporting of telemetry data (traces, metrics, and logs). Configured with dual OTLP receivers over gRPC and HTTP, an active health check extension, and internal metrics. When selected, the CLI automatically generates the `configs/otel-collector-config.yaml` file with overwrite prompt protection.
- **Volume Mappings:**
  - `./configs/otel-collector-config.yaml:/etc/otelcol-contrib/config.yaml:ro`

---

## Version 1.2.0 (`feature/cgen-1/add-nginx-feature`)

- **Release Date:** August 22, 2026
- **Relevant Commits:** `7229145`, `4b8e7e7`

### Nginx Web Server & Reverse Proxy

- **Service Name:** Nginx Web Server (`nginx`)
- **Template File:** `templates/nginx.yml`
- **Config Template:** `templates/configs/nginx.conf`
- **Image:** `nginx:alpine`
- **Exposed Ports:** `80:80`
- **Description:**
  Enhanced the Nginx service integration with full scaffolding support. When selected, the CLI automatically generates the `configs/nginx.conf` file (with an overwrite confirmation prompt if one already exists) and scaffolds the local `./html` directory for serving static web assets.
- **Volume Mappings:**
  - `./configs/nginx.conf:/etc/nginx/nginx.conf:ro`
  - `./html:/usr/share/nginx/html:ro`

---

## Version 1.1.0 (`development`)

- **Release Date:** August 21, 2026
- **Relevant Commits:** `f5a00cc`

### MinIO S3 Object Storage

- **Service Name:** MinIO S3 Object Storage (`minio`)
- **Template File:** `templates/minio.yml`
- **Image:** `minio/minio:latest`
- **Command:** `server /data --console-address ":9001"`
- **Exposed Ports:**
  - `127.0.0.1:9000:9000` (S3 API Endpoint)
  - `127.0.0.1:9001:9001` (Web Admin Console)
- **Default Credentials:**
  - **User:** `minioadmin`
  - **Password:** `minioadminpassword`
- **Description:**
  Introduced MinIO high-performance S3-compatible distributed object storage. Configured with dual ports for both direct S3 API access and the interactive Web Administration Console, backed by a persistent named volume (`minio_data`).
- **Volume Mappings:**
  - `minio_data:/data`

---

## Version 1.0.0 (`master` / Initial Release)

- **Release Date:** August 20, 2026
- **Relevant Commits:** `3563ebb`

### Apache ActiveMQ Artemis

- **Service Name:** Apache ActiveMQ Artemis (`activemq`)
- **Template File:** `templates/activemq.yml`
- **Image:** `apache/activemq-artemis:latest-alpine`
- **Exposed Ports:**
  - `8161` (Web Management Console)
  - `61616` (Core / OpenWire / JMS broker port)
  - `5672` (AMQP)
  - `1883` (MQTT)
  - `61613` (STOMP)
- **Default Credentials:**
  - **User:** `admin`
  - **Password:** `adminpassword`
  - **Anonymous Login:** `true`
- **Description:**
  Added Apache ActiveMQ Artemis multi-protocol message broker support with full management console and persistent message store.

### Mailpit

- **Service Name:** Mailpit SMTP Server & Web UI (`mailpit`)
- **Template File:** `templates/mailpit.yml`
- **Image:** `axllent/mailpit:latest`
- **Exposed Ports:**
  - `1025` (SMTP server)
  - `8025` (Web UI inbox)
- **Description:**
  Added Mailpit email testing and capture tool designed for local application development. Intercepts outgoing SMTP mail and provides a responsive webmail inbox.

### Redis

- **Service Name:** Redis Cache (`redis`)
- **Template File:** `templates/redis.yml`
- **Image:** `redis:alpine`
- **Exposed Ports:**
  - `6379:6379`
- **Description:**
  Added Redis in-memory key-value data structure store for caching and session management with persistent volume storage (`redis_data`).

---

## Summary of Services

| Version Added | Service Name | Key / Image | Primary Purpose |
| :--- | :--- | :--- | :--- |
| `v1.4.0` | **Cloudflare Tunnel** | `cloudflare/cloudflared:latest` | Secure tunnel to expose local development services to the internet |
| `v1.3.0` | **OpenTelemetry Collector** | `otel/opentelemetry-collector-contrib:latest` | Telemetry collector (OTLP gRPC/HTTP) with config scaffolding |
| `v1.2.0` | **Nginx Web Server** | `nginx:alpine` | Reverse proxy and static web server with config scaffolding |
| `v1.1.0` | **MinIO** | `minio/minio:latest` | S3-compatible local object storage & web console |
| `v1.0.0` | **Apache ActiveMQ Artemis** | `apache/activemq-artemis:latest-alpine` | Multi-protocol message broker (JMS, AMQP, MQTT, STOMP) |
| `v1.0.0` | **Mailpit** | `axllent/mailpit:latest` | Local SMTP capture server and web inbox |
| `v1.0.0` | **Redis** | `redis:alpine` | In-memory key-value caching and data store |

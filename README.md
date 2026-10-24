# Compose Gen (`compose-gen`)

An interactive CLI tool written in Go to quickly generate and compose `docker-compose.yml` files for local development environments. Select the services you need, and `compose-gen` will generate or merge the configuration along with necessary scaffolding (such as Nginx configuration files and directories).

---

## Features

- **Interactive Service Selection**: Choose services via numbers (e.g., `1, 3`), single choices, or type `all` to select all available services.
- **Embedded Templates**: Template definitions are embedded directly into the Go binary (`go:embed`), requiring no external dependencies at runtime.
- **Smart Merging & Conflict Handling**: Merges into existing `docker-compose.yml` configurations without accidentally overwriting unselected services, with confirmation prompts.
- **Automatic Config Scaffolding**: Automatically scaffolds service-specific configuration files (such as `configs/nginx.conf` and `configs/otel-collector-config.yaml`).
- **Cloudflare Tunnel URL Auto-Sync**: When `cloudflared` is used, `compose-gen up` automatically parses the generated `https://*.trycloudflare.com` URL from container logs and updates `APP_BASE_URL` in `configs/.env`.
- **Integrated Docker Compose Facade**: Convenience subcommands (`up`, `down`, `ps`, `logs`, `build`, `stop`) to manage the container lifecycle directly.
- **Lightweight & Portable**: Single standalone binary without runtime dependencies.

---

## Supported Services

| Service | Category | Ports | Credentials / Notes |
| :--- | :--- | :--- | :--- |
| **Apache ActiveMQ Artemis** | Message Broker | `8161` (Web Console), `61616` (JMS/Core), `5672` (AMQP), `1883` (MQTT), `61613` (STOMP) | User: `admin`<br>Password: `adminpassword` |
| **Cloudflare Tunnel** | Ingress / Secure Tunnel | Outbound Tunnel (No exposed host ports) | Exposes local services via HTTPS (`*.trycloudflare.com`); auto-syncs `APP_BASE_URL` to `configs/.env` |
| **Mailpit** | Email Testing / SMTP | `1025` (SMTP), `8025` (Web UI) | Local webmail inbox for testing |
| **MinIO** | S3 Object Storage | `9000` (S3 API), `9001` (Web Console) | User: `minioadmin`<br>Password: `minioadminpassword` |
| **Nginx** | Reverse Proxy / Web Server | `80` (HTTP) | Includes auto-generated `configs/nginx.conf` & `./html` volume |
| **OpenTelemetry Collector** | Observability / Telemetry | `4317` (OTLP gRPC), `4318` (OTLP HTTP), `8888` (Metrics), `13133` (Health) | Includes auto-generated `configs/otel-collector-config.yaml` |
| **PostgreSQL Database** | Relational Database | `5432` | User: `postgres`<br>Password: `postgrespassword`<br>Database: `app_db` |
| **Redis** | In-Memory Cache / KV Store | `6379` | Persistent volume enabled |

---

## Prerequisites

- [Go](https://go.dev/dl/) `1.24` or later (for building from source)
- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/)

---

## Installation & Build

### Build from Source

Clone the repository and build the binary:

```bash
git clone https://github.com/nidqija/compose-gen.git
cd compose-gen
go build -o compose-gen .
```

On Windows (PowerShell / CMD):
```powershell
go build -o compose-gen.exe .
```

### Install to Go Bin Path

```bash
go install .
```

---

## Usage

### Interactive Compose Generation

Run the tool without arguments in the directory where you want your `docker-compose.yml` generated:

```bash
compose-gen
```

#### Selection Syntax
- **Single item**: `2`
- **Multiple items**: `1, 3, 5` (comma-separated)
- **All items**: `all`

### Managing Services (CLI Commands)

`compose-gen` includes built-in wrappers around Docker Compose commands:

```bash
compose-gen up      # Start services in background (and auto-fetch Cloudflare URL if enabled)
compose-gen down    # Stop and remove containers
compose-gen logs    # Follow container logs
compose-gen ps      # View container status
compose-gen build   # Rebuild services
compose-gen stop    # Stop running containers
compose-gen help    # View available commands
```

---

## Project Structure

```text
.
├── main.go               # CLI entry point, menu logic & subcommand routing
├── features/
│   └── commands.go       # Docker Compose wrappers & Cloudflare URL syncing
├── go.mod                # Go module definition
├── go.sum                # Go checksum file
├── templates/            # YAML service definitions (embedded)
│   ├── activemq.yml
│   ├── cloudflared.yml
│   ├── mailpit.yml
│   ├── minio.yml
│   ├── nginx.yml
│   ├── otel.yml
│   ├── postgres.yml
│   ├── redis.yml
│   └── configs/
│       ├── nginx.conf                  # Nginx configuration template
│       └── otel-collector-config.yaml  # OpenTelemetry collector configuration template
├── LICENSE               # MIT License
└── README.md             # Project documentation
```

---

## Adding New Service Templates

To add a new service template:

1. Create a new YAML file in the `templates/` directory (e.g., `templates/postgres.yml`):
   ```yaml
   name: PostgreSQL Database (Port 5432)
   services:
     image: postgres:16-alpine
     environment:
       POSTGRES_USER: devuser
       POSTGRES_PASSWORD: devpassword
       POSTGRES_DB: appdb
     ports:
       - "5432:5432"
     volumes:
       - postgres_data:/var/lib/postgresql/data
     restart: unless-stopped
   volumes:
     postgres_data: {}
   ```
2. Recompile the binary with `go build -o compose-gen .`. The new service will automatically appear in the interactive list!

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

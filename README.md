# Compose Gen (`compose-gen`)

An interactive CLI tool written in Go to quickly generate and compose `docker-compose.yml` files for local development environments. Select the services you need, and `compose-gen` will generate or merge the configuration along with necessary scaffolding (such as Nginx configuration files and directories).

---

## Features

- **Interactive Service Selection**: Choose services via numbers (e.g., `1, 3`), single choices, or type `all` to select all available services.
- **Embedded Templates**: Template definitions are embedded directly into the Go binary (`go:embed`), requiring no external dependencies at runtime.
- **Smart Merging & Conflict Handling**: Merges into existing `docker-compose.yml` configurations without accidentally overwriting unselected services, with confirmation prompts.
- **Automatic Config Scaffolding**: Automatically scaffolds service-specific configuration files (such as `configs/nginx.conf` and `html/` directory for Nginx).
- **Lightweight & Portable**: Single standalone binary without runtime dependencies.

---

## Supported Services

| Service | Category | Ports | Credentials / Notes |
| :--- | :--- | :--- | :--- |
| **Apache ActiveMQ Artemis** | Message Broker | `8161` (Web Console), `61616` (JMS/Core), `5672` (AMQP), `1883` (MQTT), `61613` (STOMP) | User: `admin`<br>Password: `adminpassword` |
| **Mailpit** | Email Testing / SMTP | `1025` (SMTP), `8025` (Web UI) | Local webmail inbox for testing |
| **MinIO** | S3 Object Storage | `9000` (S3 API), `9001` (Web Console) | User: `minioadmin`<br>Password: `minioadminpassword` |
| **Nginx** | Reverse Proxy / Web Server | `80` (HTTP) | Includes auto-generated `configs/nginx.conf` & `./html` volume |
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

Run the tool in the directory where you want your `docker-compose.yml` generated:

```bash
./compose-gen
```

### Example Interactive Flow

```text
Available templates:
1. Apache ActiveMQ Artemis (JMS: 61616, Web Console: 8161)
2. Mailpit SMTP & Web UI (SMTP: 1025, Web: 8025)
3. MinIO S3 Object Storage (API: 9000, Web Console: 9001)
4. Nginx Web Server (Port 80)
5. Redis Cache (Port 6379)

Enter choices: 3, 5
docker-compose.yml generated successfully in docker-compose.yml
```

### Options for Selection
- **Single item**: `2`
- **Multiple items**: `1, 3, 5` (comma-separated)
- **All items**: `all`

### Start the Services

After generating your `docker-compose.yml`:

```bash
docker compose up -d
```

---

## Project Structure

```text
.
├── main.go               # CLI application entry point and logic
├── go.mod                # Go module definition
├── go.sum                # Go checksum file
├── templates/            # YAML service definitions (embedded)
│   ├── activemq.yml
│   ├── mailpit.yml
│   ├── minio.yml
│   ├── nginx.yml
│   ├── redis.yml
│   └── configs/
│       └── nginx.conf    # Default Nginx configuration template
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

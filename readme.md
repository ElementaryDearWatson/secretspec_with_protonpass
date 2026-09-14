# Go Podman App with SecretSpec & Proton Pass

A lightweight Go web application containerized with Podman, demonstrating secure secret management using [SecretSpec](https://github.com/cachix/secretspec) and Proton Pass Secure Notes.

---

## Features

- **Secure Secret Injection:** SecretSpec dynamically resolves secrets at runtime from Proton Pass without hardcoding credentials or storing sensitive files in source control.
- **Proton Pass Integration:** Uses `pass-cli` to securely fetch secrets from Proton Pass Secure Notes.
- **Podman & Alpine Containerization:** Multi-stage build producing a lightweight, minimal container image.
- **Local Fallback Logging:** Visual terminal indicators to demonstrate application behavior when secrets are present versus missing.

---

## Prerequisites

- [Go](https://golang.org/doc/install) **1.27+**
- [Podman](https://podman.io/) (or Docker)
- [SecretSpec CLI](https://github.com/cachix/secretspec)
- [Proton Pass CLI (`pass-cli`)](https://github.com/protonpass/pass-cli) authenticated with your Proton account

---

## Project Structure

```text
.
├── main.go           # Go application entrypoint
├── go.mod            # Go module definition
├── go.sum            # Go module checksums
├── secretspec.toml   # SecretSpec configuration mapping secrets to Proton Pass
├── Podmanfile        # Container build file
├── deploy.sh         # Deployment script
└── README.md         # Project documentation
```

---

## Configuration (`secretspec.toml`)

Ensure `secretspec.toml` is configured with schema `revision = 1` and points to your target Proton Pass vault and item:

```toml
[project]
name = "secret spec"
revision = 1

[providers]
default = "protonpass://secretspec"

[profiles.default]
API_KEY = { description = "Proton Pass API Key", required = true }
```

> **Note:** Store your API key as a Secure Note inside the `secretspec` vault in Proton Pass.

---

## Local Development & Running

### 1. Verification

Validate your SecretSpec configuration against Proton Pass:

```bash
secretspec check
```

### 2. Run without Secrets (Demonstration / Fallback Mode)

Running the binary directly without secret injection demonstrates missing key handling:

```bash
go run main.go
```

**Output:**
```text
🔴 WARNING: API_KEY is missing! Run with 'secretspec run --' to inject secrets.
Server starting on port 8080...
```

Test the endpoint:

```bash
curl http://localhost:8080/api-data
# Output: ❌ Error 500: API Key missing from environment
```

---

### 3. Run with SecretSpec Secret Injection

Use `secretspec run` to inject secrets resolved from Proton Pass into the process environment:

```bash
secretspec run -- go run main.go
```

**Output:**
```text
🟢 SUCCESS: API_KEY successfully resolved! (Key Length: 32)
Server starting on port 8080...
```

Test the endpoint:

```bash
curl http://localhost:8080/api-data
# Output: ✅ SecretSpec active: API Key loaded successfully (Length: 32)
```

---

## Containerization with Podman

### Build the Image

```bash
podman build -f Podmanfile -t go-podman-app .
```

### Run the Container

Run the application container using SecretSpec to supply the environment variables:

```bash
secretspec run -- podman run --rm -p 8080:8080 -e API_KEY go-podman-app
```

Alternatively, run using the provided deployment script:

```bash
chmod +x deploy.sh
./deploy.sh
```

---

## Endpoints

| Endpoint | Method | Description |
| :--- | :--- | :--- |
| `/` | `GET` | Hello World landing endpoint |
| `/api-data` | `GET` | Validates API key resolution status and returns key metadata |
| `/health` | `GET` | Container healthcheck endpoint (`200 OK`) |

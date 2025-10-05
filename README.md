# Kubernetes Homelab

[![GitOps](https://img.shields.io/badge/GitOps-FluxCD-blue)](https://fluxcd.io/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.32-326CE5?logo=kubernetes&logoColor=white)](https://kubernetes.io/)
[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)](https://golang.org/)
[![Service Mesh](https://img.shields.io/badge/Service_Mesh-Istio-466BB0?logo=istio&logoColor=white)](https://istio.io/)

> A GitOps-driven Kubernetes platform demonstrating enterprise-grade DevOps practices, infrastructure automation, and cloud-native architecture patterns.

## Key DevOps Competencies Demonstrated

### 1. **Kubernetes Orchestration & Management**
- Multi-namespace architecture with proper resource isolation
- PersistentVolume management with NFS storage classes
- **KEDA**-based autoscaling for workload management

### 2. **GitOps & Continuous Deployment**
- **FluxCD** implementation for declarative continuous delivery
- Automated synchronization from Git repositories
- Kustomize based configuration
- Dependency management between application deployments

### 3. **Infrastructure as Code (IaC)**
- **KCL (Kubernetes Configuration Language)** for type-safe infrastructure definitions
- Modular, reusable configuration libraries
- Version-controlled infrastructure components
- Automated manifest generation from code
- Schema validation and type checking

### 4. **Service Mesh & Advanced Networking**
- **Istio** service mesh implementation (v1.27.0)
- Gateway API configurations for ingress traffic
- VirtualService routing for advanced traffic management
- mTLS for secure service-to-service communication

### 5. **Custom Tooling & Automation**
- **Go**-based Terminal UI (TUI) tool for infrastructure automation
- Generator discovery and scaffolding system
- Interactive TUI using **tview** library
- Automated YAML/manifest generation
- Version management across Helm charts and Docker images
- Reduces manual configuration overhead by 70%+

### 6. **Observability & Monitoring**
- **VictoriaMetrics** for metrics storage and querying
- **Kiali** dashboard for service mesh visualization
- **Prometheus** integration for kiali metrics collection
- **Perses** for custom dashboards

### 7. **Security & Secrets Management**
- **HashiCorp Vault** for secrets management
- **Vault Secrets Operator** for Kubernetes integration
- OIDC-based RBAC for Kubernetes-API authentication and authorization

### 8. **Database Management**
- **CloudNativePG** operator for PostgreSQL cluster management
- **MariaDB** Deployment
- **Redis/Valkey** for caching

## Applications

### **Self-Hosted Applications**
| Application | Purpose | Technology |
|-------------|---------|------------|
| **Jellyfin** | Media streaming server | Helm Chart (v2.3.0) |
| **Forgejo** | Self-hosted Git service | Helm Chart (v14.0.0) |
| **Mealie** | Recipe management | Docker (v3.0.2) |
| **Booklore** | Book library management | Custom deployment |
| **Matrix Synapse** | Federated chat server | Custom + bridges |
| **Searxng** | Privacy-respecting search | Docker (2025.8.3) |
| **Glance** | Dashboard application | Docker (v0.8.4) |
| **Redlib** | Reddit alternative frontend | Docker |

### **Infrastructure Components**
- **Caddy** - Reverse proxy with automatic HTTPS
- **Gluetun** - VPN client for selective routing
- **CNPG** - PostgreSQL operator
- **KEDA** - Kubernetes-based event-driven autoscaling
- **CSI Driver NFS** - Network storage provisioning
- **Vault** - Secrets management platform
- **Vault-Secrets-Operator** - Manage Kubernetes Secrets from Vault
- **Reflector** - ConfigMap/Secret replication

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        FluxCD GitOps                        │
│                  (Continuous Reconciliation)                │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Istio Service Mesh                     │
│              (Gateway, VirtualServices, mTLS)               │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌──────────────┐    ┌──────────────────┐    ┌──────────────┐
│     Apps     │    │  Infrastructure  │    │  Monitoring  │
│              │    │                  │    │              │
│ • Jellyfin   │    │ • Vault          │    │ • Victoria   │
│ • Forgejo    │    │ • PostgreSQL     │    │   Metrics    │
│ • Mealie     │    │ • Redis/Valkey   │    │ • Kiali      │
│ • Matrix     │    │ • Caddy          │    │ • Perses     │
│ • Booklore   │    │ • KEDA           │    │ • ...        │
│ • ...        │    │ • ...            │    │              │
└──────────────┘    └──────────────────┘    └──────────────┘
        │                     │                     │
        └─────────────────────┴─────────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  NFS Storage     │
                    │  (PVC/PV)        │
                    └──────────────────┘
```

---

## Custom Terminal UI Tool
> **Todo**: Standalone CLI planned but not implemented yet

### **Features**
The custom Go Terminal UI tool (`cmd/cli-tool/`) provides:

1. **Generator System**
   - Automated discovery of infrastructure generators
   - Category-based organization (Apps, Infrastructure, Istio, Monitoring)
   - One-command deployment of complete application stacks

2. **Interactive TUI**
   - Real-time command output visualization
   - Tree-based generator navigation
   - Keyboard shortcuts for rapid operations

3. **Scaffolding Automation**
   - Template-based generator creation
   - Automatic boilerplate generation
   - Consistent project structure enforcement

4. **Version Management**
   - Centralized version tracking for Helm charts and Docker images
   - Automated version bumps across configurations
   - JSON-based version manifests

### **Usage Example**
```bash
# Start the interactive Terminal UI
go run cmd/cli-tool/main.go

# Commands available:
# [r] Run - Execute a generator
# [v] Version - Upgrade Helm/Docker versions
# [d] Discover - Auto-discover new generators
# [s] Scaffolding - Create new generator templates
# [q] Quit
```

---

## Project Structure

```
kubernetes-main/
├── cluster/                      # Kubernetes manifests (GitOps source)
│   ├── apps/                     # Application deployments
│   ├── infrastructure/           # Core infrastructure components
│   ├── istio/                    # Service mesh configuration
│   ├── monitoring/               # Observability stack
│   └── flux/                     # FluxCD bootstrap and automation
│
├── config_as_code/               # KCL configuration files (legacy)
│   ├── lib/                      # Reusable KCL libraries
│   │   ├── k8s_wrapper/          # Kubernetes resource wrappers
│   │   ├── helm_flux/            # Helm + Flux integration
│   │   ├── istio/                # Istio resource definitions
│   │   └── vault/                # Vault integration helpers
│   └── src/                      # Application-specific configs
│       ├── apps/
│       ├── infrastructure/
│       ├── istio/
│       └── monitoring/
│
├── cmd/                          # Go CLI tool
│   └── cli-tool/                 # Main CLI application
│
├── internal/                     # Internal Go packages
│   ├── cli/                      # CLI logic (TUI, scaffolding)
│   ├── generators/               # Generator implementations
│   └── pkg/utils/                # Utility functions
│
├── pkg/schema/                   # Go schemas for config validation
│   ├── cluster/                  # Cluster-level schemas
│   ├── generator/                # Generator metadata schemas
│   └── k8s/                      # Kubernetes resource schemas
│
├── setups/                       # Infrastructure setup scripts
│   ├── vps-setup/                # VPS hardening (UFW, Fail2ban)
│   ├── raspberry-pi-setup/       # ARM-based setup scripts
│   └── minecraft-server/         # Game server configuration
│
└── versions/                     # Centralized version tracking
    ├── apps.json
    ├── infrastructure.json
    ├── istio.json
    └── monitoring.json
```

## Quick Start

### **Prerequisites**
- Kubernetes cluster (1.28+)
- kubectl configured
- flux CLI installed
- Go 1.25+ (for CLI tool)

### **Bootstrap FluxCD**
```bash
flux bootstrap git \
  --components-extra=image-reflector-controller,image-automation-controller \
  --url=ssh://git@your-git-server.com/your-repo/kubernetes.git \
  --branch=main \
  --private-key-file=./id_ed25519 \
  --path=cluster/flux
```

### **Run the CLI Tool**
```bash
# Build and run
cd cmd/cli-tool
go run main.go

# Or build binary
go build -o k8s-cli
./k8s-cli
```

---

## Infrastructure Setup Scripts

Pre-configured hardening scripts for various environments:

### **VPS Setup** (`setups/vps-setup/`)
- UFW firewall configuration
- Fail2ban intrusion prevention
- SSH key-based authentication
- Root user access restriction
- NFS server setup

---

## Future Enhancements / Todo

Planned improvements include:
- [ ] Terraform/OpenTofu configuration of Vault 
- [ ] Backup strategies
- [ ] Chaos engineering integration (Chaos Mesh)
- [ ] Network Policies
---

## Operational Commands Reference

<details>
<summary>Click to expand useful kubectl commands</summary>

### Delete terminating namespace
```sh
kubectl get namespace "netbird" -o json \
  | tr -d "\n" | sed "s/\"finalizers\": \[[^]]\+\]/\"finalizers\": []/" \
  | kubectl replace --raw /api/v1/namespaces/netbird/finalize -f -
```

### Delete stuck CRD
```sh
kubectl patch crd nbgroups.netbird.io -p '{"metadata":{"finalizers":[]}}' --type=merge
```

### Remove stuck pod
```sh
kubectl delete pod <podname> -n <namespace> --force --grace-period=0
```

### Trigger CronJob manually
```sh
kubectl create job --from=cronjob/postgres-postgresql-pgdumpall manual-postgres -n postgres
```

### Creating a token for Kiali
```sh
kubectl create token kiali-service-account -n istio-system
```

### Copy cert from Caddy
```sh
kubectl cp <pod-name>:/data/caddy/pki/authorities/local/root.crt ./root.crt -n caddy
kubectl cp <pod-name>:/data/caddy/pki/authorities/local/root.key ./root.key -n caddy
openssl pkcs12 -export -in root.crt -inkey root.key -out caddy-server.p12
```

### Port-forward for database access
```sh
kubectl port-forward -n postgres svc/forgejo-postgresql-rw 5432:5432
```

### Create PostgreSQL dump
```sh
pg_dump --host 127.0.0.1 --username forgejo --format=custom --blobs --verbose --file "/path/to/backup.dump" forgejo
```

### Download files from pod
```sh
kubectl cp -n synapse <pod-name>:/data/media_store ./media_store
```

### Copying files from/to Vault
```sh
# copy from
kubectl cp <pod-name>:/vault/data ./vault/ -n vault
# copy to
kubectl cp ./vault-backup/ <pod-name>:/vault/data -n vault
```

### Mounting NFS
```sh
sudo mount -t nfs <server name>:<location on server> <local location>
```

### Generating Synapse config
```sh
podman run -ti --rm \
    --mount type=volume,src=synapse-data,dst=/data \
    -e SYNAPSE_SERVER_NAME=matrix.cluster.netbird.selfhosted \
    -e SYNAPSE_REPORT_STATS=yes \
    -e SYNAPSE_DATA_DIR=/media \
    ghcr.io/element-hq/synapse:latest generate
```

### Generate a new user in Synapse
```sh
register_new_matrix_user -u <username>
```

### Delete PostgreSQL backups
```sh
kubectl delete backup -A --all
```

### Changing registry mirrors
```sh
sudo nano /etc/rancher/k3s/registries.yaml
sudo systemctl restart k3s
```

</details>

---
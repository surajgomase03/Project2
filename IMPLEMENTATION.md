# Project2 Implementation Guide

Complete step-by-step guide for setting up Project2 from scratch.

---

## System Architecture

### Overall Architecture Diagram

```mermaid
graph TB
    subgraph "Development"
        A["Developer Pushes Code<br/>to GitHub"]
    end
    
    subgraph "CI/CD Pipeline"
        B["GitHub Actions<br/>CI Workflow"]
        C1["1. Build & Test"]
        C2["2. Code Quality<br/>golangci-lint"]
        C3["3. Build Docker<br/>Image"]
        C4["4. Push to<br/>Docker Hub"]
        C5["5. Update<br/>Helm Chart"]
    end
    
    subgraph "GitOps"
        D["Git Repository<br/>with Helm Charts"]
        E["ArgoCD<br/>CD Operator"]
    end
    
    subgraph "AWS EKS Cluster"
        F["Kubernetes Cluster"]
        G1["Deployment"]
        G2["Service"]
        G3["Ingress"]
        H["Running Pod<br/>surajgomase/project2:v1"]
    end
    
    subgraph "External"
        I["Docker Hub<br/>Registry"]
        J["User Access<br/>http://localhost:8080"]
    end
    
    A -->|push| B
    B --> C1 --> C2 --> C3 --> C4 --> I
    C3 --> C5 --> D
    D -->|watches| E
    E -->|deploys| F
    F --> G1 --> H
    F --> G2
    F --> G3
    H --> J
```

### CI/CD Pipeline Flow

```mermaid
graph LR
    A["Developer<br/>Commits Code"] -->|git push| B["GitHub<br/>Webhook"]
    B -->|triggers| C["GitHub<br/>Actions"]
    
    subgraph CI["CI Pipeline"]
        C1["Build<br/>Job"]
        C2["Code Quality<br/>Job"]
        C3["Docker Build<br/>Job"]
        C4["Helm Update<br/>Job"]
        C1 -->|success| C2
        C2 -->|success| C3
        C3 -->|success| C4
    end
    
    C --> CI
    
    C4 -->|commits| D["Git Repo<br/>Updated"]
    D -->|change event| E["ArgoCD<br/>CD"]
    E -->|syncs| F["EKS Cluster"]
    F --> G["✅ New Version<br/>Running"]
```

### Docker Build Process

```mermaid
graph TB
    A["Dockerfile<br/>Multi-stage Build"] -->|Stage 1| B["Go Builder<br/>golang:1.21"]
    B --> C["Copy go.mod<br/>& go.sum"]
    C --> D["RUN go mod<br/>download"]
    D --> E["Copy main.go<br/>+ templates<br/>+ static"]
    E --> F["RUN go build<br/>-o main"]
    
    F -->|Stage 2| G["Alpine<br/>Base Image"]
    G --> H["RUN apk add<br/>ca-certificates"]
    H --> I["COPY FROM<br/>builder: main<br/>templates<br/>static"]
    I --> J["EXPOSE 8080<br/>CMD ./main"]
    
    J -->|final| K["Docker Image<br/>~18.7MB"]
    K -->|docker push| L["Docker Hub<br/>surajgomase/project2:v1"]
```

### Kubernetes Deployment Architecture

```mermaid
graph TB
    A["AWS EKS Cluster"] -->|contains| B["Kubernetes Objects"]
    
    B --> C1["Deployment<br/>go-deployment"]
    B --> C2["Service<br/>go-service"]
    B --> C3["Ingress<br/>go-ingress"]
    
    C1 --> D1["ReplicaSet<br/>2 replicas"]
    D1 --> E1["Pod 1<br/>port 8080"]
    D1 --> E2["Pod 2<br/>port 8080"]
    
    C2 -->|Load Balancer| E1
    C2 -->|Load Balancer| E2
    
    E1 -->|running| F1["Container<br/>surajgomase/project2:v1<br/>port 8080"]
    E2 -->|running| F2["Container<br/>surajgomase/project2:v1<br/>port 8080"]
    
    C3 -->|routes traffic| C2
    
    F1 --> G["Go App<br/>main.go"]
    F2 --> G
    G -->|serves| H["templates/<br/>static/"]
```

### Helm Chart Structure

```mermaid
graph TB
    A["Helm Chart<br/>go-web-chart"] -->|contains| B["Chart.yaml<br/>Chart metadata"]
    A -->|contains| C["values.yaml<br/>Default values"]
    A -->|contains| D["templates/<br/>K8s manifests"]
    
    D --> D1["deployment.yaml<br/>Uses: image, tag<br/>replicaCount"]
    D --> D2["service.yaml<br/>Uses: service.type<br/>port"]
    D --> D3["ingress.yaml<br/>K8s ingress rules"]
    
    B --> B1["name: go-web-chart"]
    B --> B2["version: 0.1.0"]
    
    C --> C1["image.repository:<br/>surajgomase/project2"]
    C --> C2["image.tag: v1"]
    C --> C3["replicaCount: 2"]
    
    E["helm install<br/>go-web-app<br/>./go-web-chart"] -->|reads| A
    E -->|deploys| F["Kubernetes<br/>Resources"]
```

### ArgoCD GitOps Flow

```mermaid
graph LR
    A["Git Repository<br/>Helm/go-web-chart/"] -->|ArgoCD<br/>Watches| B["ArgoCD<br/>Application"]
    
    B -->|Desired State<br/>from Git| C["Git Commit<br/>updates values.yaml<br/>new image tag"]
    
    B -->|Current State<br/>in Cluster| D["EKS Cluster<br/>Running Pods"]
    
    B -->|Detects<br/>Drift| E{"Desired !=<br/>Current?"}
    
    E -->|Yes| F["Auto Sync<br/>or Manual Sync"]
    F -->|Apply Changes| D
    E -->|No| G["✅ In Sync"]
    
    D -->|Update| H["New Pods<br/>with Latest Image"]
```

### Complete End-to-End Flow

```mermaid
graph TB
    A["1. Developer<br/>Commits Code<br/>git push origin main"] -->|triggers| B["2. GitHub Actions<br/>Webhook"]
    
    B -->|starts| C["3. Build Job<br/>- Checkout<br/>- Setup Go 1.21<br/>- go build<br/>- go test"]
    
    C -->|success| D["4. Code Quality<br/>- golangci-lint<br/>- Check errors"]
    
    D -->|success| E["5. Docker Build<br/>- docker build<br/>- docker push"]
    
    E -->|pushed| F["Docker Hub<br/>surajgomase/project2:v1"]
    
    E -->|success| G["6. Update Helm<br/>- Update values.yaml<br/>- image.tag = v1<br/>- git push"]
    
    G -->|commit| H["7. Git Repo<br/>Updated<br/>Helm/values.yaml"]
    
    H -->|ArgoCD<br/>Detects Change| I["8. ArgoCD<br/>Syncs<br/>Kubernetes"]
    
    I -->|deploys| J["9. EKS Cluster<br/>- Pulls image<br/>- Updates pods<br/>- New version running"]
    
    J -->|result| K["✅ Application<br/>Running Latest<br/>Version"]
```

---

## Step 1: Code Analysis

**Objective:** Understand the codebase

- Get code from dev team
- Analyze all code thoroughly
- Check:
  - **Ports** - What port does the app run on? (e.g., 8080)
  - **Language** - What programming language? (e.g., Go 1.21)
  - **Dependencies** - What packages/libraries are used?
  - **Configuration** - Environment variables, configs needed?

**Commands:**
```bash
go mod tidy        # Check dependencies
go build           # Test local build
go test ./...      # Run tests
```

---

## Step 2: Docker Setup

**Objective:** Containerize the application

### 2.1 Write Dockerfile

Create `Dockerfile` in project root:

```dockerfile
# Build stage
FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY main.go .
COPY ./static ./static
COPY ./templates ./templates

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

EXPOSE 8080
CMD ["./main"]
```

### 2.2 Build Docker Image

```bash
docker build -t surajgomase/project2:v1 .
```

### 2.3 Test Locally

```bash
docker run -p 8080:8080 surajgomase/project2:v1
# Visit http://localhost:8080
```

---

## Step 3: EKS Cluster Provisioning

**Objective:** Set up Kubernetes cluster on AWS EKS

### 3.1 Prerequisites

```bash
# Install AWS CLI
aws --version

# Install kubectl
kubectl version --client

# Install eksctl
eksctl version
```

### 3.2 Create EKS Cluster

```bash
eksctl create cluster \
  --name project2-cluster \
  --region us-east-1 \
  --nodegroup-name standard-nodes \
  --node-type t3.medium \
  --nodes 2
```

### 3.3 Verify Cluster

```bash
kubectl get nodes
kubectl get pods -A
```

---

## Step 4: Kubernetes Manifest & Deployment

**Objective:** Deploy application on Kubernetes

### 4.1 Create K8s Manifests

Create `K8s/` folder with three files:

**K8s/Deployment.yml:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-deployment
spec:
  replicas: 2
  selector:
    matchLabels:
      app: go-web
  template:
    metadata:
      labels:
        app: go-web
    spec:
      containers:
      - name: go-app
        image: surajgomase/project2:v1
        ports:
        - containerPort: 8080
```

**K8s/Service.yml:**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: go-service
spec:
  selector:
    app: go-web
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: 8080
```

**K8s/Ingress.yml:**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: go-ingress
spec:
  rules:
  - host: example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: go-service
            port:
              number: 80
```

### 4.2 Deploy on Kubernetes

```bash
kubectl apply -f K8s/Deployment.yml
kubectl apply -f K8s/Service.yml
kubectl apply -f K8s/Ingress.yml
```

### 4.3 Verify Deployment

```bash
kubectl get deployments
kubectl get pods
kubectl get services
kubectl describe pod <pod-name>
```

---

## Step 5: Helm Setup

**Objective:** Package application with Helm for easier management

### 5.1 Install Helm

```bash
# macOS
brew install helm

# Linux
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

### 5.2 Create Helm Chart

```bash
helm create go-web-chart
```

This creates:
```
go-web-chart/
├── Chart.yaml
├── values.yaml
├── templates/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── ...
└── charts/
```

---

## Step 6: Helm Chart Customization

**Objective:** Use K8s manifests in Helm chart

### 6.1 Delete Unnecessary Files

```bash
cd go-web-chart
rm -rf charts/              # Delete charts folder
rm -rf templates/*          # Delete all templates
```

### 6.2 Copy K8s Manifests to Helm Templates

```bash
# Copy your K8s manifests to templates/
cp ../K8s/Deployment.yml templates/deployment.yaml
cp ../K8s/Service.yml templates/service.yaml
cp ../K8s/Ingress.yml templates/ingress.yaml
```

### 6.3 Update Chart.yaml

```yaml
apiVersion: v2
name: go-web-chart
description: A Helm chart for Project2 Go application
type: application
version: 0.1.0
appVersion: "1.0"
```

### 6.4 Update values.yaml

```yaml
# values.yaml
image:
  repository: surajgomase/project2
  tag: "v1"
  pullPolicy: IfNotPresent

replicaCount: 2

service:
  type: LoadBalancer
  port: 80
  targetPort: 8080
```

### 6.5 Modify Templates to Use Variables

Update `templates/deployment.yaml`:

```yaml
image: {{ .Values.image.repository }}:{{ .Values.image.tag }}
replicas: {{ .Values.replicaCount }}
```

---

## Step 7: Deploy with Helm

**Objective:** Install application using Helm chart

```bash
# Install Helm chart
helm install go-web-app ./go-web-chart

# Verify
helm list
helm status go-web-app

# Get pods
kubectl get pods

# Access application
kubectl get svc
```

---

## Step 8: CI/CD Setup with GitHub Actions

**Objective:** Automate build, test, and deployment

### 8.1 Create GitHub Actions Workflow

Create `.github/workflows/ci.yaml`:

```yaml
name: CI

on:
  push:
    branches:
      - main
    paths-ignore:
      - 'docs/**'
      - 'tests/**'
      - 'Helm/**'
      - 'error.md'
      - 'IMPLEMENTATION.md'
      - 'README.md'
      - 'EKS/**'
      - 'GITOPS/**'
      - 'ingress-controller/**'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Application build
        run: CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

      - name: Run tests
        run: go test ./...

  code-quality:
    runs-on: ubuntu-latest
    needs: build
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: v1.59.1

  build-docker-image:
    runs-on: ubuntu-latest
    needs: code-quality
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      - name: Docker login
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Build and push Docker image
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: ${{ secrets.DOCKER_USERNAME }}/project2:${{ github.run_id }}

  Update_helm_charts:
    runs-on: ubuntu-latest
    needs: build-docker-image
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          token: ${{ secrets.TOKEN }}

      - name: Update Helm Chart values.yaml
        run: |
          IMAGE_TAG=${{ secrets.DOCKER_USERNAME }}/project2:${{ github.run_id }}
          sed -i "s|tag: .*|tag: '${IMAGE_TAG##*:}'|g" Helm/go-web-chart/values.yaml

      - name: Commit and Push changes
        run: |
          git config --global user.name "github-actions"
          git add Helm/go-web-chart/values.yaml
          git commit -m "Update Helm chart with new image tag ${{ github.run_id }}"
          git push origin main
```

---

## Step 9: CI/CD Pipeline Explanation

### CI Steps (Continuous Integration)

1. **Build & Unit Test**
   - Checkout code
   - Set up Go 1.21
   - Build application: `go build -o main .`
   - Run tests: `go test ./...`

2. **Static Code Analysis**
   - Run golangci-lint with v1.59.1
   - Check for code quality issues
   - Fail if issues found

3. **Create & Push Docker Image**
   - Build Docker image with Dockerfile
   - Push to Docker Hub with run ID tag
   - Example: `surajgomase/project2:12345`

4. **Update Helm Chart**
   - Pull latest repo
   - Update `Helm/go-web-chart/values.yaml` with new image tag
   - Push changes back

### CD Steps (Continuous Deployment)

1. **Install ArgoCD**
   ```bash
   kubectl create namespace argocd
   kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
   ```

2. **Connect GitOps Repository**
   - Point ArgoCD to your Git repository
   - ArgoCD watches for changes in `Helm/` folder

3. **Deploy on Cluster**
   - ArgoCD detects new Helm chart changes
   - Automatically deploys to EKS cluster
   - Updates running pods with new image

---

## CI/CD Workflow Summary

```
Developer Pushes Code
         ↓
GitHub Actions Triggered
         ↓
Build & Test (Step 1) ✅
         ↓
Code Quality Check (Step 2) ✅
         ↓
Build & Push Docker Image (Step 3) ✅
         ↓
Update Helm Chart (Step 4) ✅
         ↓
Push Changes to Git
         ↓
ArgoCD Detects Changes
         ↓
Deploy to EKS Cluster (CD)
         ↓
✅ Application Running with Latest Code
```

---

## Setup GitHub Secrets

Required secrets for CI/CD to work:

1. **DOCKER_USERNAME** - Your Docker Hub username
2. **DOCKER_PASSWORD** - Your Docker Hub password/token
3. **TOKEN** - GitHub Personal Access Token (for pushing changes)

**How to Add Secrets:**
```
GitHub Repo → Settings → Secrets and variables → Actions → New repository secret
```

---

## Complete Deployment Checklist

- [ ] Code analyzed and understood
- [ ] Dockerfile created and tested locally
- [ ] Docker image built and pushed to Docker Hub
- [ ] EKS cluster provisioned and accessible
- [ ] K8s manifests created (Deployment, Service, Ingress)
- [ ] Helm chart created with custom templates
- [ ] values.yaml configured for image tags
- [ ] Helm chart deployed successfully
- [ ] GitHub Actions workflow file created
- [ ] GitHub secrets configured
- [ ] ArgoCD installed and connected
- [ ] First CI/CD pipeline run successful
- [ ] Application deployed on EKS cluster

---

## Troubleshooting

### Docker Image Not Found on Docker Hub
```bash
docker push surajgomase/project2:v1
```

### Kubernetes Pod Not Starting
```bash
kubectl describe pod <pod-name>
kubectl logs <pod-name>
```

### Helm Chart Update Not Deploying
```bash
helm upgrade go-web-app ./go-web-chart
kubectl rollout restart deployment go-deployment
```

### CI Pipeline Failing
Check `.github/workflows/ci.yaml` syntax and GitHub Actions logs

---

## References

- [Helm Documentation](https://helm.sh/docs/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [ArgoCD Documentation](https://argo-cd.readthedocs.io/)
- [EKS User Guide](https://docs.aws.amazon.com/eks/)

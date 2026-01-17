# Docker and Kubernetes Errors - Explanation

## Error 1: Docker Template and Static Folder Error

### The Problem
When running the Docker container, you got an error:
```
template not found
```

The Go website was trying to load HTML files from the `templates` folder and CSS files from the `static` folder, but these folders were **not inside the Docker container**.

### Why It Happened

#### Before Fix (Dockerfile)
```dockerfile
COPY --from=builder /app/main .
```

This only copied the **main binary** (the executable program) to the container. 

The `templates` and `static` folders were left behind on your computer, NOT in the container.

```
Your Computer:
├── main (copied ✓)
├── templates/ (NOT copied ✗)
├── static/ (NOT copied ✗)
└── go.mod

Inside Container:
└── main (only this exists ✗)
```

### The Solution

#### After Fix (Updated Dockerfile)
```dockerfile
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
```

Now we copy **everything needed**:
- The `main` executable program
- The `templates` folder with HTML files
- The `static` folder with CSS files

```
Inside Container (Now):
├── main ✓
├── templates/ ✓
│   ├── home.html
│   ├── about.html
│   └── contact.html
└── static/ ✓
    └── style.css
```

### In Simple Terms

Think of it like packing a suitcase for a trip:
- **Before**: You only packed your passport (main binary)
- **After**: You packed your passport + clothes + shoes (main + templates + static)

When the Go website runs inside Docker, it needs ALL these files to work properly. If any are missing, it crashes!

### How to Rebuild

```bash
docker build -t surajgomase/project2:v1 .
docker push surajgomase/project2:v1
docker run -p 8080:8080 -it surajgomase/project2:v1
```

Now visit `http://localhost:8080` and everything should work! ✓

---

## Error 2: Kubernetes ImagePullBackOff Error

### Error Message
```
Warning  Failed     4s (x4 over 107s)  kubelet  Failed to pull image "surajgomase/project2:v1": 
Error response from daemon: manifest for surajgomase/project2:v1 not found: manifest unknown
Error: ErrImagePull
Error: ImagePullBackOff
```

### Root Cause
The Docker image `surajgomase/project2:v1` exists **locally** on your machine but is **not on Docker Hub**. Kubernetes (Minikube) cannot pull the image from the registry because it hasn't been pushed there.

### Local Image Status
```
REPOSITORY              TAG       IMAGE ID       CREATED        SIZE
surajgomase/project2    latest    2394c8ee58ac   10 hours ago   18.7MB
surajgomase/project2    v1        2394c8ee58ac   10 hours ago   18.7MB
```

The image runs successfully locally with Docker:
```bash
docker run -p8080:8080 -it surajgomase/project2:v1
2026/01/17 18:42:57 Server starting on http://localhost:8080
```

But Kubernetes can't access it!

### Why This Happens

**Local Docker** = Image only on your computer
**Docker Hub** = Remote registry where Docker stores public images
**Kubernetes** = Needs to pull from Docker Hub (or another registry)

### Solutions

#### Option 1: Push Image to Docker Hub (Recommended for Production)
```bash
docker push surajgomase/project2:v1
```

#### Option 2: Load Local Image into Minikube (Best for Development)
```bash
minikube image load surajgomase/project2:v1
```

Then update your Helm values with `imagePullPolicy: Never`

#### Option 3: Use Minikube's Docker Daemon
```bash
eval $(minikube docker-env)
docker build -t surajgomase/project2:v1 .
```

---

## Error 3: Helm Release Not Found Error

### Error Message
```
Error: uninstall: Release not loaded: go-web-app: release: not found
Error: uninstall: Release not loaded: go-web-chart: release: not found
```

### Root Cause
The Helm releases `go-web-app` and `go-web-chart` don't exist or haven't been installed yet. You're trying to uninstall something that was never installed.

### Check What's Deployed
```bash
helm list -a
kubectl get deployments -A
kubectl get pods -A
```

### Reinstall Properly

First, fix the image pull issue (use Option 1 or 2 above), then:

```bash
helm install go-web-app ./go-web-chart
```

Or with local image:
```bash
helm install go-web-app ./go-web-chart --set imagePullPolicy=Never
```

### Verify Installation
```bash
helm list
kubectl get pods
kubectl describe pod <pod-name>
```

---

## Error 4: GitHub Actions golangci-lint Version Error

### Error Message
```
Error: Failed to run: Error: requested golangci-lint version 'v1.21' isn't supported: 
we support only v1.28.3 and later versions
```

### Root Cause
Your GitHub Actions CI workflow (`.github/workflows/ci.yaml`) is requesting golangci-lint version `v1.21`, which is no longer supported by the golangci-lint-action. The action only supports versions **v1.28.3 and later**.

### Location
**File:** `.github/workflows/ci.yaml` (Line 42)

### Current Code (Wrong)
```yaml
- name: Run golangci-lint
  uses: golangci/golangci-lint-action@v6
  with:
    version: v1.21
```

### Solution

Update the version to a supported version:

```yaml
- name: Run golangci-lint
  uses: golangci/golangci-lint-action@v6
  with:
    version: v1.28.3
```

Or use the latest version:

```yaml
- name: Run golangci-lint
  uses: golangci/golangci-lint-action@v6
  with:
    version: v1.59.1
```

### How to Fix

1. Open `.github/workflows/ci.yaml` in VS Code
2. Find line 42: `version: v1.21`
3. Change it to: `version: v1.28.3` (or latest)
4. Save the file
5. Commit and push:

```bash
git add .github/workflows/ci.yaml
git commit -m "fix: update golangci-lint version to v1.28.3"
git push origin main
```

Your GitHub Actions pipeline will now run successfully! ✓

---

## Summary Table

| Error | Cause | Solution |
|-------|-------|----------|
| Template not found | Missing files in Docker image | Copy templates & static folders in Dockerfile |
| ImagePullBackOff | Image not on Docker Hub | Push to Docker Hub or load into Minikube |
| Release not found | Helm not installed | Install Helm chart after fixing image issue |
| golangci-lint version | v1.21 no longer supported | Update to v1.28.3 or later in ci.yaml |

---

## Quick Fix Checklist

- [ ] Update Dockerfile to copy templates and static folders
- [ ] Rebuild Docker image: `docker build -t surajgomase/project2:v1 .`
- [ ] Push to Docker Hub: `docker push surajgomase/project2:v1`
- [ ] Update golangci-lint version in `.github/workflows/ci.yaml`
- [ ] Install Helm chart: `helm install go-web-app ./go-web-chart`
- [ ] Verify: `kubectl get pods`
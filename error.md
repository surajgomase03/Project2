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

## Summary Table

| Error | Cause | Solution |
|-------|-------|----------|
| Template not found | Missing files in Docker image | Copy templates & static folders in Dockerfile |
| ImagePullBackOff | Image not on Docker Hub | Push to Docker Hub or load into Minikube |
| Release not found | Helm not installed | Install Helm chart after fixing image issue |

---

## Quick Fix Checklist

- [ ] Update Dockerfile to copy templates and static folders
- [ ] Rebuild Docker image: `docker build -t surajgomase/project2:v1 .`
- [ ] Push to Docker Hub: `docker push surajgomase/project2:v1`
- [ ] Install Helm chart: `helm install go-web-app ./go-web-chart`
- [ ] Verify: `kubectl get pods`


# Docker Template and Static Folder Error - Explanation

## The Problem
When running the Docker container, you got an error:
```
template not found
```

The Go website was trying to load HTML files from the `templates` folder and CSS files from the `static` folder, but these folders were **not inside the Docker container**.

## Why It Happened

### Before Fix (Dockerfile)
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

## The Solution

### After Fix (Updated Dockerfile)
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

## In Simple Terms

Think of it like packing a suitcase for a trip:
- **Before**: You only packed your passport (main binary)
- **After**: You packed your passport + clothes + shoes (main + templates + static)

When the Go website runs inside Docker, it needs ALL these files to work properly. If any are missing, it crashes!

## How to Rebuild

```bash
docker build -t surajgomase/project1:latest .
docker push surajgomase/project1:latest
docker run -p 8080:8080 -it surajgomase/project1:latest
```

Now visit `http://localhost:8080` and everything should work! ✓

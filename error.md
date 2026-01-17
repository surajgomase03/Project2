# Project Errors & Solutions Guide

Quick reference for all errors and how to fix them.

---

## Error 1: Template Not Found (Docker)

**Problem:** When running Docker container: `template not found`

**Why:** Missing `templates/` and `static/` folders in Docker image

**Fix:** Update Dockerfile

```dockerfile
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates    ✅ Add this
COPY --from=builder /app/static ./static           ✅ Add this
```

**Test:**
```bash
docker build -t surajgomase/project2:v1 .
docker run -p 8080:8080 surajgomase/project2:v1
```

---

## Error 2: ImagePullBackOff (Kubernetes)

**Problem:** Kubernetes can't pull image from Docker Hub
```
Failed to pull image "surajgomase/project2:v1": manifest not found
Error: ImagePullBackOff
```

**Why:** Image exists locally but NOT on Docker Hub

**Fix - Option A (Recommended for Production):**
```bash
docker push surajgomase/project2:v1
```

**Fix - Option B (For Local Development):**
```bash
minikube image load surajgomase/project2:v1
```

---

## Error 3: Helm Release Not Found

**Problem:**
```
Error: uninstall: Release not loaded: go-web-app: release: not found
```

**Why:** Helm chart was never installed

**Fix:**
```bash
helm install go-web-app ./go-web-chart
```

**Verify:**
```bash
helm list
kubectl get pods
```

---

## Error 4: golangci-lint Version v1.21 Not Supported

**Problem (GitHub Actions):**
```
Error: requested golangci-lint version 'v1.21' isn't supported: 
we support only v1.28.3 and later versions
```

**Why:** v1.21 is too old

**Fix:** Update `.github/workflows/ci.yaml`

```yaml
- name: Run golangci-lint
  uses: golangci/golangci-lint-action@v6
  with:
    version: v1.59.1  ✅ Use 1.59.1 or higher
```

---

## Error 5: Go Version Mismatch (Linter Can't Parse Packages)

**Problem (GitHub Actions):**
```
Error: buildir: failed to load package goarch: could not load export data: 
cannot import "internal/goarch" (unknown bexport format version -1)
```

**Why:** Go 1.24.11 on GitHub Actions vs Go 1.21 in your Dockerfile - package format mismatch

**Fix:** Add Go 1.21 setup to ALL jobs in `.github/workflows/ci.yaml`

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'  ✅ Add this to build job

      - name: Application build
        run: CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

  code-quality:
    runs-on: ubuntu-latest
    needs: build
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'  ✅ Add this to code-quality job

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: v1.59.1
```

---

## Error 6: Invalid Parameter 'go-version' for golangci-lint-action

**Problem (GitHub Actions):**
```
Warning: Unexpected input(s) 'go-version', valid inputs are 
['version', 'install-mode', 'working-directory', 'github-token', ...]
```

**Why:** The `go-version` parameter doesn't exist in golangci-lint-action@v6

**Fix:** Remove `go-version` parameter, use the `Set up Go` step instead (see Error 5)

❌ **Wrong:**
```yaml
- name: Run golangci-lint
  uses: golangci/golangci-lint-action@v6
  with:
    version: v1.59.1
    go-version: '1.21'  ❌ This parameter doesn't exist
```

✅ **Right:**
```yaml
- name: Set up Go
  uses: actions/setup-go@v4
  with:
    go-version: '1.21'  ✅ Correct place for Go version

- name: Run golangci-lint
  uses: golangci/golangci-lint-action@v6
  with:
    version: v1.59.1
```

---

## Error 7: Unchecked Error from tmpl.Execute() (errcheck)

**Problem (GitHub Actions):**
```
Error: main.go:36:14: Error return value of `tmpl.Execute` is not checked (errcheck)
	tmpl.Execute(w, data)
	            ^
```

**Why:** You're not checking if template execution failed

**Fix:** Update `main.go` - check error from `tmpl.Execute()`

❌ **Wrong:**
```go
tmpl, err := template.ParseFiles("templates/home.html")
if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}
tmpl.Execute(w, data)  ❌ Not checking error!
```

✅ **Right:**
```go
tmpl, err := template.ParseFiles("templates/home.html")
if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}
if err := tmpl.Execute(w, data); err != nil {  ✅ Check error
    log.Printf("Error executing template: %v", err)
}
```

**Do this for ALL 3 handlers:**
- homeHandler (line 36)
- aboutHandler (line 50)
- contactHandler (line 71)

---

## Error 8: Git Push Rejected (Remote Has Changes)

**Problem:**
```
To https://github.com/surajgomase03/Project2.git
 ! [rejected]        main -> main (fetch first)
error: failed to push some refs to 'https://github.com/surajgomase03/Project2.git'
hint: Updates were rejected because the remote contains work that you do not
hint: have locally.
```

**Why:** Someone else pushed changes to GitHub that you don't have locally

**Fix:** Pull changes first, then push

```bash
git pull origin main
git push origin main
```

**What This Does:**
1. **git pull** - Fetches remote changes and merges them into your local branch
2. **git push** - Pushes your local commits to remote

**If There Are Merge Conflicts:**
```bash
# Git will mark conflicts in files
# Edit the files to resolve conflicts
# Then commit and push
git add .
git commit -m "Merge remote changes"
git push origin main
```

---

## Error 9: CI Infinite Loop (Helm Updates Trigger CI Again)

**Problem:**
```
CI Pipeline runs → Updates Helm chart → Pushes changes → CI triggers again → Updates Helm chart again → Pushes → CI triggers again → LOOP!
```

**Why This Happens:**

1. CI workflow updates `Helm/go-web-chart/values.yaml` with new Docker image tag
2. Helm changes are pushed back to GitHub
3. **Helm/ was NOT in paths-ignore** → CI gets triggered again
4. CI runs again and updates Helm
5. Pushes changes → CI triggers again
6. **Infinite loop! ♻️**

**Visual:**
```
Push code → CI starts → Updates Helm → Push to GitHub
                                            ↓
                     ← ← ← ← ← ← ← ← ← ← ← ←
```

**Solution:** Add unnecessary paths to `paths-ignore` in `.github/workflows/ci.yaml`

```yaml
on:
  push:
    branches:
      - main
    paths-ignore:
      - 'docs/**'
      - 'tests/**'
      - 'Helm/**'                    ✅ Ignore Helm updates
      - 'error.md'                   ✅ Ignore docs
      - 'IMPLEMENTATION.md'
      - 'README.md'
      - 'EKS/**'                     ✅ Ignore config folders
      - 'GITOPS/**'
      - 'ingress-controller/**'
```

**How It Works Now:**

1. Push code changes → CI starts ✅
2. CI updates Helm → Pushes changes
3. Helm/ is in paths-ignore → **CI does NOT trigger** ✅
4. **No loop!**

**Files to Ignore:**
- `Helm/**` - Helm charts (auto-updated by CI)
- `error.md` - Documentation
- `IMPLEMENTATION.md` - Documentation
- `README.md` - Documentation
- `EKS/**` - EKS config (manual only)
- `GITOPS/**` - GitOps config (manual only)
- `ingress-controller/**` - Ingress config (manual only)

**Test It:**
```bash
# Make a change in Helm/ and push
git add Helm/go-web-chart/values.yaml
git commit -m "test: update Helm values manually"
git push origin main

# Check GitHub Actions - CI should NOT trigger ✅
```

---

## Quick Fix Checklist

Copy & paste this to remember what needs to be fixed:

- [ ] **Dockerfile** - Add template/static copies
  ```bash
  git add Dockerfile
  git commit -m "fix: copy templates and static folders in Dockerfile"
  ```

- [ ] **main.go** - Fix tmpl.Execute() error handling
  ```bash
  git add main.go
  git commit -m "fix: add error handling for tmpl.Execute() calls"
  ```

- [ ] **.github/workflows/ci.yaml** - Add Go 1.21 setup to all jobs
  ```bash
  git add .github/workflows/ci.yaml
  git commit -m "fix: add Go 1.21 setup to all CI jobs"
  ```

- [ ] **Push Docker image** - For Kubernetes to pull it
  ```bash
  docker build -t surajgomase/project2:v1 .
  docker push surajgomase/project2:v1
  ```

- [ ] **Deploy to Kubernetes**
  ```bash
  helm install go-web-app ./go-web-chart
  kubectl get pods
  ```

---

## Error Summary Table

| # | Error | Cause | Solution |
|---|-------|-------|----------|
| 1 | Template not found | Missing files in Docker | Copy templates & static in Dockerfile |
| 2 | ImagePullBackOff | Image not on Docker Hub | Push image or load to Minikube |
| 3 | Release not found | Helm not installed | Run `helm install go-web-app ./go-web-chart` |
| 4 | golangci-lint v1.21 | Old version unsupported | Update to v1.59.1 |
| 5 | Go package format error | Go version mismatch | Add `Set up Go 1.21` step to all CI jobs |
| 6 | Invalid go-version param | Wrong action parameter | Remove from golangci-lint, use setup-go instead |
| 7 | tmpl.Execute unchecked | Missing error check | Add `if err := tmpl.Execute(...)`  |

---

## Testing Each Fix

### Test Docker Build
```bash
docker build -t surajgomase/project2:v1 .
docker run -p 8080:8080 surajgomase/project2:v1
# Visit http://localhost:8080
```

### Test Go Build
```bash
go build -o main .
./main
```

### Test Linting (Local)
```bash
golangci-lint run
```

### Test Kubernetes
```bash
helm install go-web-app ./go-web-chart
kubectl get pods
kubectl logs <pod-name>
```

---

## All Changes Summary

**Files Modified:**
1. `Dockerfile` - Added COPY for templates & static
2. `main.go` - Added error handling to tmpl.Execute()
3. `.github/workflows/ci.yaml` - Added Go 1.21 setup, updated golangci-lint to v1.59.1

**All changes pushed to:** `git push origin main`

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
Your GitHub Actions CI workflow (`.github/workflows/ci.yaml`) was requesting golangci-lint version `v1.21`, which is no longer supported by the golangci-lint-action. The action only supports versions **v1.28.3 and later**.

### Solution

Update the version to a supported version in `.github/workflows/ci.yaml`:

```yaml
- name: Run golangci-lint
  uses: golangci/golangci-lint-action@v6
  with:
    version: v1.28.3
```

---

## Error 5: Go Version Compatibility Error (golangci-lint & Go 1.24)

### Error Message
```
level=error msg="Running error: buildir: failed to load package goarch: could not load export data: 
cannot import \"internal/goarch\" (unknown bexport format version...)"
Error: golangci-lint exit with code 3
```

### Root Cause
GitHub Actions runner is using **Go 1.24.11**, but **golangci-lint 1.28.3** has compatibility problems with this newer Go version. There's a version skew between the Go compiler and the lint tool.

**Version Mismatch:**
- Your code: Built with Go 1.21 (see Dockerfile)
- GitHub Actions: Running Go 1.24.11 (latest)
- golangci-lint 1.28.3: Doesn't fully support Go 1.24

### Why This Happens

When Go versions are too different, the compiled packages have different formats. golangci-lint can't read the export data from the newer Go version.

### Solution

Fix your `.github/workflows/ci.yaml` to explicitly set Go 1.21 and use a compatible golangci-lint version:

```yaml
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
        version: v1.55.2
        go-version: '1.21'
```

### Key Changes

1. **Add Go setup** to both `build` and `code-quality` jobs
2. **Lock Go version to 1.21** (matches your Dockerfile)
3. **Upgrade golangci-lint to v1.55.2** (compatible with Go 1.21)

### Commit Changes

```bash
git add .github/workflows/ci.yaml
git commit -m "fix: set Go version to 1.21 and upgrade golangci-lint to v1.55.2"
git push origin main
```

Your CI pipeline will now pass! ✓

---

## Error 6: golangci-lint Invalid Input Parameter Error

### Error Message
```
Warning: Unexpected input(s) 'go-version', valid inputs are ['version', 'install-mode', 'working-directory', 'github-token', 'verify', 'only-new-issues', 'skip-cache', 'skip-save-cache', 'problem-matchers', 'args', 'cache-invalidation-interval']
level=error msg="Running error: buildir: failed to load package goarch: could not load export data: cannot import \"internal/goarch\" (unknown bexport format version -1..."
Error: golangci-lint exit with code 3
```

### Root Cause
The golangci-lint-action@v6 **doesn't support the `go-version` parameter**. It was trying to use an invalid parameter, and golangci-lint v1.28.3 still cannot parse Go 1.24.11 export format.

### The Issue

Your CI workflow had:
```yaml
- name: Run golangci-lint
  uses: golangci/golangci-lint-action@v6
  with:
    version: v1.28.3
    go-version: '1.21'  # ❌ Invalid parameter!
```

The action ignored the `go-version` parameter and still ran with Go 1.24.11, causing the same export format error.

### Solution

Remove the invalid `go-version` parameter and **upgrade golangci-lint to v1.60.0** which natively supports Go 1.24:

```yaml
- name: Run golangci-lint
  uses: golangci/golangci-lint-action@v6
  with:
    version: v1.60.0
```

### Why v1.60.0 Works

- ✅ Supports Go 1.24.11 export format
- ✅ Only uses valid parameters
- ✅ No version skew issues

### Commit Changes

```bash
git add .github/workflows/ci.yaml
git commit -m "fix: upgrade golangci-lint to v1.60.0 for Go 1.24 support"
git push origin main
```

---

## Summary Table

| Error | Cause | Solution |
|-------|-------|----------|
| Template not found | Missing files in Docker image | Copy templates & static folders in Dockerfile |
| ImagePullBackOff | Image not on Docker Hub | Push to Docker Hub or load into Minikube |
| Release not found | Helm not installed | Install Helm chart after fixing image issue |
| golangci-lint v1.21 unsupported | Old version no longer supported | Update to v1.28.3 or later |
| Go 1.24 compatibility (v1.28.3) | Version skew between Go and golangci-lint | Upgrade golangci-lint to v1.60.0 |
| Invalid go-version parameter | `go-version` not supported by golangci-lint-action@v6 | Remove parameter, use v1.60.0 |

---

## Quick Fix Checklist

- [ ] Update Dockerfile to copy templates and static folders
- [ ] Rebuild Docker image: `docker build -t surajgomase/project2:v1 .`
- [ ] Push to Docker Hub: `docker push surajgomase/project2:v1`
- [ ] Update `.github/workflows/ci.yaml` with Go 1.21 setup
- [ ] Fix unchecked errors in main.go (tmpl.Execute)
- [ ] Pull before push: `git pull origin main`
- [ ] Install Helm chart: `helm install go-web-app ./go-web-chart`
- [ ] Verify: `kubectl get pods`
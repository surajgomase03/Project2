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

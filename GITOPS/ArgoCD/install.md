# Argo CD Installation Guide

This guide provides step-by-step instructions to install and configure Argo CD using manifests.

## Prerequisites

- Kubernetes cluster (EKS, AKS, GKE, or local cluster)
- `kubectl` command-line tool configured to access your cluster

## Installation Steps

### 1. Create Argo CD Namespace

```bash
kubectl create namespace argocd
```

This creates a dedicated namespace for Argo CD components.

### 2. Install Argo CD Using Manifests

```bash
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

This command applies the official Argo CD installation manifests from the stable branch to your cluster.

**Wait for the installation to complete:**
```bash
kubectl wait --for=condition=Ready pod -l app.kubernetes.io/name=argocd-server -n argocd --timeout=300s
```

## Access Argo CD UI

### Using NodePort Service

#### For Linux/Mac:
```bash
kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "NodePort"}}'
```

#### For Windows (PowerShell):
```powershell
kubectl patch svc argocd-server -n argocd -p '{\"spec\": {\"type\": \"NodePort\"}}'
```

### Get the NodePort Service Details

```bash
kubectl get svc argocd-server -n argocd
```

This command will display the NodePort assigned to the Argo CD UI.

**Expected output:**
```
NAME            TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)
argocd-server   NodePort   10.x.x.x        <none>        80:30XXX/TCP, 443:30XXX/TCP
```

### Access the UI

Get your Node IP:
```bash
kubectl get nodes -o wide
```

Then access Argo CD UI at:
```
https://<NODE-IP>:<NODE-PORT>
```

Example: `https://192.168.1.100:30443`

**Note:** You may see a self-signed certificate warning. Accept it to proceed.

## Get Initial Admin Password

```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo
```

**For Windows PowerShell:**
```powershell
$secret = kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}"
[System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String($secret))
```

## Default Credentials

- **Username:** `admin`
- **Password:** Use the command above to retrieve the initial password

**Important:** Change the default password after first login.

## Verify Installation

```bash
kubectl get all -n argocd
```

This should show all Argo CD components are running:
- argocd-server (UI and API)
- argocd-repo-server (Git repository access)
- argocd-controller-manager (Reconciliation engine)
- argocd-notifications-controller
- argocd-application-controller

## Troubleshooting

### Check Pod Status
```bash
kubectl get pods -n argocd
```

### View Pod Logs
```bash
kubectl logs -n argocd -l app.kubernetes.io/name=argocd-server
```

### Delete and Reinstall
```bash
kubectl delete namespace argocd
```

Then repeat the installation steps above.

## Next Steps

1. Log in to the Argo CD UI
2. Add your Git repositories
3. Create applications for GitOps deployment
4. Configure notifications and webhooks

## References

- [Argo CD Documentation](https://argo-cd.readthedocs.io/)
- [Argo CD GitHub Repository](https://github.com/argoproj/argo-cd)

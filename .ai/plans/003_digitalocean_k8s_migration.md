# Migration to DigitalOcean Kubernetes (DOKS)

**Goal:** Migrate from Raspberry Pi K8s cluster to DigitalOcean Kubernetes while keeping GitOps workflow and adding monitoring

**Estimated Cost:** $40-55/month
**Timeline:** 1-2 days
**Difficulty:** Medium (you already know K8s!)

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     Cloudflare CDN + Tunnel                  │
│                   summitseekers.co.za                        │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              DigitalOcean Kubernetes Cluster                 │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Frontend   │  │   Backend    │  │ Cloudflared  │      │
│  │   (nginx)    │  │   (Go API)   │  │   Tunnel     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │         Flux CD (GitOps Controller)                  │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │   Monitoring Stack (Prometheus + Grafana)            │  │
│  └──────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│        DigitalOcean Managed PostgreSQL Database              │
│              (Automated backups, HA optional)                │
└─────────────────────────────────────────────────────────────┘
```

---

## Phase 1: DigitalOcean Setup (30 mins)

### Step 1.1: Create DigitalOcean Account
1. Sign up at https://www.digitalocean.com
2. Add payment method
3. (Optional) Use referral link for $200 credit: https://m.do.co/c/your-referral

### Step 1.2: Install doctl CLI
```bash
# macOS
brew install doctl

# Authenticate
doctl auth init
# Paste your DigitalOcean API token when prompted
```

### Step 1.3: Create DOKS Cluster
```bash
# Create cluster with 2 basic nodes (recommended for HA)
doctl kubernetes cluster create run-goals-cluster \
  --region sfo3 \
  --version 1.31.1-do.4 \
  --node-pool "name=worker-pool;size=s-2vcpu-4gb;count=2;auto-scale=true;min-nodes=1;max-nodes=3"

# Or create via web UI:
# - Go to Kubernetes in DigitalOcean dashboard
# - Click "Create Cluster"
# - Choose region closest to your users (sfo3, nyc3, or fra1)
# - Node pool: Basic nodes, 2GB RAM ($12/month each)
# - Start with 2 nodes for HA
```

**Node Size Options:**
- `s-2vcpu-2gb` - $18/month (2GB RAM) - minimum viable
- `s-2vcpu-4gb` - $24/month (4GB RAM) - recommended
- Start with 1 node to save money, scale to 2+ later

### Step 1.4: Get Cluster Credentials
```bash
# Download kubeconfig
doctl kubernetes cluster kubeconfig save run-goals-cluster

# Verify connection
kubectl cluster-info
kubectl get nodes
```

---

## Phase 2: Database Setup (20 mins)

### Option A: Managed PostgreSQL (Recommended - $15/month)

**Pros:**
- Automated backups (daily, 7-day retention)
- Automated updates
- Easy to scale
- Connection pooling built-in
- Monitoring included

**Setup:**
```bash
# Create managed PostgreSQL via web UI:
# 1. Go to Databases → Create Database
# 2. Choose PostgreSQL 16
# 3. Select same region as your cluster
# 4. Choose "Basic" plan - 1GB RAM ($15/month)
# 5. Name: summitseekers-db
# 6. Create database: run_goals
# 7. Wait 3-5 minutes for provisioning

# Or via CLI:
doctl databases create summitseekers-db \
  --engine pg \
  --region sfo3 \
  --size db-s-1vcpu-1gb \
  --version 16

# Get connection details:
doctl databases get summitseekers-db
```

**Important:** Note down these from the dashboard:
- Host (private network address)
- Port (usually 25060)
- Username (usually doadmin)
- Password
- Database name

### Option B: In-Cluster PostgreSQL (Cheaper - ~$0, uses existing node resources)

**Pros:**
- No extra cost
- Full control
- Works with existing manifests

**Cons:**
- You manage backups
- You manage updates
- Shares resources with apps

**Setup:** Keep your existing [k8s/app/db.yaml](../k8s/app/db.yaml), but update storage class (see Phase 3)

---

## Phase 3: Update Kubernetes Manifests

### Changes Needed:

#### 3.1: Update Database Configuration

**If using Managed PostgreSQL:**

Create new secret for database connection:
```yaml
# k8s/app/db-secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: postgres-credentials
  namespace: default
type: Opaque
stringData:
  host: "summitseekers-db-do-user-XXXXX-0.b.db.ondigitalocean.com"
  port: "25060"
  username: "doadmin"
  password: "YOUR_PASSWORD_HERE"  # Get from DO dashboard
  database: "run_goals"
  sslmode: "require"
```

**If using in-cluster PostgreSQL:**

Update [k8s/app/db.yaml](../k8s/app/db.yaml):
```yaml
# Change storage class from default to do-block-storage
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: do-block-storage  # ADD THIS LINE
  resources:
    requests:
      storage: 1Gi
```

#### 3.2: Update Backend Deployment

Update [k8s/app/backend.yaml](../k8s/app/backend.yaml):

```yaml
# If using managed PostgreSQL, update env vars:
env:
  - name: DATABASE_HOST
    valueFrom:
      secretKeyRef:
        name: postgres-credentials
        key: host
  - name: DATABASE_PORT
    valueFrom:
      secretKeyRef:
        name: postgres-credentials
        key: port
  - name: DATABASE_USER
    valueFrom:
      secretKeyRef:
        name: postgres-credentials
        key: username
  - name: DATABASE_PASSWORD
    valueFrom:
      secretKeyRef:
        name: postgres-credentials
        key: password
  - name: DATABASE_DBNAME
    valueFrom:
      secretKeyRef:
        name: postgres-credentials
        key: database
  - name: DATABASE_SSLMODE
    valueFrom:
      secretKeyRef:
        name: postgres-credentials
        key: sslmode
  # ... rest of env vars (Strava credentials, etc.)
```

#### 3.3: Update Ingress

Update [k8s/app/ingress.yaml](../k8s/app/ingress.yaml):

**No changes needed!** Traefik and Cloudflare Tunnel will work the same way.

#### 3.4: Remove Raspberry Pi Specific Configs

Check for and remove:
- ARM-specific node selectors
- Raspberry Pi node affinity rules
- Any `flannel-lease-gc` specific to your Pi cluster

Review [k8s/flux-system/flannel-lease-gc.yaml](../k8s/flux-system/flannel-lease-gc.yaml) - you may not need this on DOKS.

---

## Phase 4: Secrets Migration (15 mins)

### 4.1: Recreate Strava Credentials Secret
```bash
# Get your Strava credentials from current cluster or .env file
kubectl create secret generic strava-credentials \
  --from-literal=STRAVA_CLIENT_ID='your_client_id' \
  --from-literal=STRAVA_CLIENT_SECRET='your_client_secret' \
  --namespace=default
```

### 4.2: Recreate Cloudflare Tunnel Credentials
```bash
# Get credentials.json from your Raspberry Pi cluster
kubectl get secret cloudflared-credentials -n default -o yaml > cloudflared-secret-backup.yaml

# Apply to new cluster
kubectl apply -f cloudflared-secret-backup.yaml
```

### 4.3: Create Database Secret (if using managed PostgreSQL)
```bash
kubectl create secret generic postgres-credentials \
  --from-literal=host='your-db-host.db.ondigitalocean.com' \
  --from-literal=port='25060' \
  --from-literal=username='doadmin' \
  --from-literal=password='your_password' \
  --from-literal=database='run_goals' \
  --from-literal=sslmode='require' \
  --namespace=default
```

---

## Phase 5: Install Flux CD (20 mins)

### 5.1: Install Flux CLI
```bash
# macOS
brew install fluxcd/tap/flux

# Verify
flux --version
```

### 5.2: Check Prerequisites
```bash
flux check --pre
```

### 5.3: Bootstrap Flux on DOKS

You have two options:

**Option A: Use existing GitHub repo (recommended)**
```bash
# Export GitHub token
export GITHUB_TOKEN=<your-github-personal-access-token>

# Bootstrap Flux using your existing repo
flux bootstrap github \
  --owner=craeyefish \
  --repository=run-goals \
  --branch=main \
  --path=./k8s \
  --personal \
  --private=false
```

**Option B: Fresh Flux install**
If you want to redo the Flux setup:
```bash
flux bootstrap github \
  --owner=craeyefish \
  --repository=run-goals \
  --branch=main \
  --path=./k8s/flux-system \
  --personal
```

### 5.4: Verify Flux Installation
```bash
# Check Flux components
kubectl get pods -n flux-system

# Check GitRepository
flux get sources git

# Check Kustomizations
flux get kustomizations
```

### 5.5: Update Flux Image Automation

Your existing image automation should work! But verify these files:
- [k8s/flux-system/backend-image.yaml](../k8s/flux-system/backend-image.yaml)
- [k8s/flux-system/frontend-image.yaml](../k8s/flux-system/frontend-image.yaml)

Flux will continue watching GHCR for new images and auto-update your deployments.

---

## Phase 6: Deploy Applications (10 mins)

### 6.1: Let Flux Reconcile
```bash
# Force reconciliation
flux reconcile source git flux-system

# Watch deployments
kubectl get deployments -n default -w
```

### 6.2: Verify Pods
```bash
# Check all pods are running
kubectl get pods -n default

# Check logs
kubectl logs -n default deployment/run-goals-backend
kubectl logs -n default deployment/run-goals-frontend
kubectl logs -n default deployment/cloudflared
```

### 6.3: Initialize Database (if fresh start)
```bash
# Get backend pod name
kubectl get pods -n default -l app=run-goals-backend

# Exec into backend pod to run migrations (if you have them)
kubectl exec -it <backend-pod-name> -n default -- /bin/sh

# Or just let the backend auto-migrate on first run (check your Go code)
```

---

## Phase 7: Monitoring Stack (30 mins)

### 7.1: Install kube-prometheus-stack

Create monitoring namespace and install Prometheus + Grafana:

```bash
# Add Prometheus community Helm repo
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Create monitoring namespace
kubectl create namespace monitoring

# Install kube-prometheus-stack
helm install kube-prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --set prometheus.prometheusSpec.retention=7d \
  --set grafana.adminPassword='admin' \
  --set prometheus.prometheusSpec.storageSpec.volumeClaimTemplate.spec.resources.requests.storage=5Gi \
  --set prometheus.prometheusSpec.storageSpec.volumeClaimTemplate.spec.storageClassName=do-block-storage
```

### 7.2: Access Grafana Dashboard
```bash
# Port-forward Grafana
kubectl port-forward -n monitoring svc/kube-prometheus-grafana 3000:80

# Open http://localhost:3000
# Username: admin
# Password: admin (or what you set above)
```

### 7.3: Pre-configured Dashboards

The kube-prometheus-stack comes with these dashboards:
- **Kubernetes / Compute Resources / Cluster** - Overall cluster health
- **Kubernetes / Compute Resources / Namespace (Pods)** - Per-namespace metrics
- **Kubernetes / Compute Resources / Pod** - Individual pod metrics
- **Node Exporter / Nodes** - Node-level metrics

### 7.4: Create Custom Dashboard for Your App

In Grafana, create a dashboard to monitor:
- Backend API request rate
- Backend response times
- Database connection pool usage
- Frontend page load times
- Pod CPU/Memory usage

### 7.5: Set Up Alerts (Optional)

Create Prometheus alert rules for:
- Pod restart loops
- High memory usage (>80%)
- High CPU usage (>80%)
- Database connection failures

---

## Phase 8: DNS and Cloudflare Tunnel (15 mins)

### 8.1: Verify Cloudflared Deployment
```bash
# Check cloudflared pod
kubectl get pods -n default -l app=cloudflared

# Check logs
kubectl logs -n default -l app=cloudflared

# Verify tunnel is connected
# You should see: "Connection established" in logs
```

### 8.2: Update Cloudflare Tunnel (if needed)

Your existing [k8s/app/cloudflared.yaml](../k8s/app/cloudflared.yaml) should work as-is! It's already configured to route traffic to Traefik.

**No DNS changes needed** - Cloudflare Tunnel will automatically route to your new cluster.

### 8.3: Verify Traefik Ingress

```bash
# Check Traefik is running (should be installed by DOKS by default)
kubectl get pods -n kube-system -l app.kubernetes.io/name=traefik

# If not installed, DOKS uses DigitalOcean Load Balancer by default
# Your ingress resources will create a DO Load Balancer ($12/month)
```

**Cost-saving tip:** Use Cloudflare Tunnel (FREE) instead of DO Load Balancer ($12/month). Your current setup already does this!

---

## Phase 9: Testing (20 mins)

### 9.1: Health Checks
```bash
# Test backend health endpoint
kubectl port-forward -n default svc/run-goals-backend 8080:8080
curl http://localhost:8080/api/health  # or whatever health endpoint you have

# Test frontend
kubectl port-forward -n default svc/run-goals-frontend 4200:80
open http://localhost:4200
```

### 9.2: Test Through Cloudflare Tunnel
```bash
# Visit your domain
open https://summitseekers.co.za

# Test API endpoints
curl https://summitseekers.co.za/api/health
```

### 9.3: Test Strava OAuth Flow
1. Go to https://summitseekers.co.za
2. Click "Connect Strava"
3. Authorize with Strava
4. Verify redirect back to your app
5. Check that activities sync

### 9.4: Test Webhook
```bash
# Trigger a Strava activity sync
# Check backend logs for webhook processing
kubectl logs -n default -l app=run-goals-backend --tail=100 -f
```

---

## Phase 10: Monitoring and Optimization (Ongoing)

### 10.1: Monitor Costs
```bash
# Check DigitalOcean usage:
doctl invoice list
doctl invoice summary get <invoice-id>

# Watch your monthly spend in DO dashboard
```

### 10.2: Right-size Your Cluster

After running for a week, check actual resource usage:
```bash
# Check node resource usage
kubectl top nodes

# Check pod resource usage
kubectl top pods -n default

# If usage is low, you can:
# - Scale down to 1 node
# - Use smaller node size (s-2vcpu-2gb instead of s-2vcpu-4gb)
```

### 10.3: Enable Cluster Autoscaling

```bash
# Already configured in cluster creation above
# Cluster will auto-scale from 1 to 3 nodes based on load
```

### 10.4: Set Up Backup Strategy

**For Managed PostgreSQL:**
- Automated daily backups (included)
- 7-day retention (included)
- Point-in-time recovery available

**For In-Cluster PostgreSQL:**
```bash
# Set up CronJob for daily pg_dump backups to DigitalOcean Spaces (S3-compatible)
# Example: https://github.com/benjamin-maynard/kubernetes-cloud-mysql-backup
```

---

## Cost Breakdown

### Minimum Setup (~$27/month):
- 1x s-2vcpu-2gb node: $12/month
- Managed PostgreSQL (1GB): $15/month
- Cloudflare Tunnel: FREE
- **Total: $27/month**

### Recommended Setup (~$42/month):
- 2x s-2vcpu-2gb nodes (HA): $24/month
- Managed PostgreSQL (1GB): $15/month
- Block Storage (5GB for monitoring): $1/month
- Cloudflare Tunnel: FREE
- **Total: $40/month**

### With DigitalOcean Load Balancer (~$54/month):
- 2x s-2vcpu-2gb nodes: $24/month
- Managed PostgreSQL: $15/month
- DO Load Balancer: $12/month
- Block Storage: $1/month
- **Total: $52/month**

**Recommendation:** Stick with Cloudflare Tunnel (your current setup) to save $12/month!

---

## Rollback Plan

If something goes wrong:

1. **Keep Raspberry Pi running** during migration (parallel run)
2. **DNS is unchanged** - Cloudflare Tunnel can route to either cluster
3. **To rollback:** Just update Cloudflare Tunnel to point back to Pi cluster
4. **Database:** Since starting fresh, no data loss concern

---

## Post-Migration Enhancements

Once stable, consider adding:

### 1. Sealed Secrets
Stop storing secrets in plain YAML, use encrypted secrets in Git:
```bash
# Install sealed-secrets controller
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/controller.yaml

# Encrypt secrets before committing to Git
kubeseal --format yaml < secret.yaml > sealed-secret.yaml
```

### 2. Better Monitoring
- Set up Slack/Discord alerts from Prometheus
- Add application-level metrics (Prometheus client in Go backend)
- Set up uptime monitoring (UptimeRobot, Better Uptime)

### 3. Cost Alerts
- Set up billing alerts in DigitalOcean
- Get notified if spending exceeds $50/month

### 4. CI/CD Improvements
- Add automated testing in GitHub Actions
- Add staging environment (small cluster for testing)
- Use Flux notifications to alert on deployment success/failure

### 5. Horizontal Pod Autoscaling
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: backend-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: run-goals-backend
  minReplicas: 1
  maxReplicas: 5
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

---

## Helpful Commands Reference

```bash
# Switch kubectl context
kubectl config get-contexts
kubectl config use-context do-sfo3-run-goals-cluster

# Watch all resources
kubectl get all -n default
watch kubectl get pods -n default

# Force Flux reconciliation
flux reconcile source git flux-system
flux reconcile kustomization flux-system

# View Flux logs
flux logs --all-namespaces --follow

# Describe pod for troubleshooting
kubectl describe pod <pod-name> -n default

# Get events (useful for debugging)
kubectl get events -n default --sort-by='.lastTimestamp'

# Port-forward services locally
kubectl port-forward -n default svc/run-goals-backend 8080:8080
kubectl port-forward -n default svc/run-goals-frontend 4200:80
kubectl port-forward -n monitoring svc/kube-prometheus-grafana 3000:80

# Scale deployments manually
kubectl scale deployment run-goals-backend --replicas=2 -n default

# Update image manually (Flux will sync it back eventually)
kubectl set image deployment/run-goals-backend backend=ghcr.io/craeyefish/run-goals-backend:backend-v0.0.13 -n default

# Shell into pod
kubectl exec -it <pod-name> -n default -- /bin/sh

# Copy files from pod
kubectl cp default/<pod-name>:/path/to/file ./local-file
```

---

## Troubleshooting

### Pods stuck in ImagePullBackOff
```bash
# Check image exists in GHCR
# Verify you're using correct image tags
kubectl describe pod <pod-name> -n default

# Check if image is public or needs imagePullSecrets
```

### Database connection errors
```bash
# Verify secret exists
kubectl get secret postgres-credentials -n default -o yaml

# Check backend env vars
kubectl exec -it <backend-pod> -n default -- env | grep DATABASE

# Test connection from backend pod
kubectl exec -it <backend-pod> -n default -- /bin/sh
# Then try: nc -zv <db-host> <db-port>
```

### Flux not reconciling
```bash
# Check Flux components are running
kubectl get pods -n flux-system

# Check GitRepository status
flux get sources git -A

# Check for errors
flux logs --all-namespaces

# Force reconciliation
flux reconcile source git flux-system --with-source
```

### Cloudflare Tunnel not connecting
```bash
# Check cloudflared logs
kubectl logs -n default -l app=cloudflared

# Verify credentials secret exists
kubectl get secret cloudflared-credentials -n default

# Verify tunnel ID in ConfigMap
kubectl get configmap cloudflared-config -n default -o yaml
```

---

## Next Steps

1. **Review** this plan and ask questions
2. **Create** DigitalOcean account
3. **Provision** DOKS cluster
4. **Set up** managed PostgreSQL
5. **Migrate** secrets
6. **Bootstrap** Flux CD
7. **Deploy** applications
8. **Install** monitoring
9. **Test** everything
10. **Enjoy** your cloud-native K8s setup! 🚀

Want me to help you execute any of these steps? I can create the updated manifest files, help with commands, or troubleshoot issues as they come up!

# Archived Raspberry Pi Configurations

This directory contains Kubernetes configurations that were specific to the Raspberry Pi cluster and are not needed for DigitalOcean Kubernetes (DOKS).

## Files Archived

### flannel-lease-gc.yaml
**Reason for archival:** This CronJob was used to clean up stale Flannel CNI network leases on the Raspberry Pi cluster. DigitalOcean Kubernetes uses Cilium as the default CNI, so this is not needed.

**Original purpose:** Flannel on Raspberry Pi clusters can sometimes leave stale IP lease files. This CronJob ran every 15 minutes to clean them up and reload Flannel.

**Safe to delete:** Yes, after confirming DOKS migration is successful.

---

## Migration Date
Archived during migration from Raspberry Pi to DigitalOcean Kubernetes on 2025-12-08.

## Restoration
If you ever need to restore these configs for a Raspberry Pi cluster, they are preserved here. Simply copy them back to their original locations.

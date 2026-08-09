# Overview

CloudFlareDDNS is a small Go utility that checks the current public IP address and updates one or more Cloudflare DNS records to match it. It is designed to be run periodically so dynamic DNS entries stay in sync without manual intervention.

The application reads its configuration from environment variables and stores the last-known IPs in a local state file so it can avoid unnecessary updates.

# Usage

This project can be run either locally or in a container. Containers and Kubernetes are the recommended deployment model because they make it easy to run the updater on a schedule and keep the persistent state mounted reliably.

## Local usage

1. Build the binary:
   ```bash
   make build
   ```
2. Set the required environment variables:
   ```bash
   export CLOUDFLAREDDNS_API_TOKEN="..."
   export CLOUDFLAREDDNS_ZONE_ID="..."
   export CLOUDFLAREDDNS_DNS_RECORD_IDS="record-id-1,record-id-2"
   export CLOUDFLAREDDNS_PERSISTENCE_ROOT="$PWD/.state"
   ```
3. Run it:
   ```bash
   ./bin/cloudflareddns
   ```

The program compares the current public IP with the last known IP for each target record and updates Cloudflare only when a change is needed.

## Container usage

Build the image:

```bash
docker build -t cloudflareddns .
```

Run it with the same environment variables and a mounted persistence directory:

```bash
docker run --rm \
  -e CLOUDFLAREDDNS_API_TOKEN="..." \
  -e CLOUDFLAREDDNS_ZONE_ID="..." \
  -e CLOUDFLAREDDNS_DNS_RECORD_IDS="record-id-1,record-id-2" \
  -e CLOUDFLAREDDNS_PERSISTENCE_ROOT=/data \
  -v "$PWD/.state:/data" \
  cloudflareddns
```

For Kubernetes, run this as a CronJob or scheduled workload and mount a persistent volume for the state directory so the last-known IPs are preserved between runs.

# Planned Improvements

1. Support IPv6 DNS records
2. Support specifying DNS records by name instead of ID.
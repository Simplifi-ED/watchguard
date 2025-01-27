# Kubernetes Resource Watcher

This project is a Go-based application that dynamically watches Kubernetes resources in a cluster using informers. It automatically discovers available API resources and sets up watchers for them, logging add, update, and delete events.

## Features

- Dynamically discovers and watches all available Kubernetes API resources.
- Skips deprecated API versions in favor of preferred versions (configurable in the code).
- Logs `Add`, `Update`, and `Delete` events for each watched resource.
- Uses the Kubernetes dynamic client and informer factory.

## Prerequisites

- A running Kubernetes cluster.
- The application should be deployed inside the cluster (uses in-cluster configuration).
- Go 1.19+ installed for local development.

## How It Works

1. **Cluster Configuration**:
   - Uses Kubernetes' in-cluster configuration to authenticate with the API server.

2. **API Resource Discovery**:
   - Retrieves all available API resources using the discovery client.

3. **Dynamic Informers**:
   - Sets up informers for all watchable resources, skipping deprecated API versions if a preferred version is specified.

4. **Event Handling**:
   - Logs `Add`, `Update`, and `Delete` events for each resource.

## Usage

### Build and Deploy

1. **Build the Application**:

  ```bash
    go build -o k8s-resource-watcher
  ```

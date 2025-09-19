# Copilot Instructions for JellyNotifier

## Project Overview
JellyNotifier is a webhook receiver microservice built in Go that processes notifications from media management tools like Jellyfin, Overseerr, or similar systems. It's a simple HTTP server designed for containerized deployment.

## Architecture
- **Single binary**: All logic in `main.go` with a simple HTTP server
- **Webhook receiver**: Accepts POST requests at `/webhook` endpoint
- **Health monitoring**: `/health` endpoint for Kubernetes probes
- **Testing**: `/test` endpoint for development/debugging
- **Containerized**: Multi-stage Docker build with Alpine Linux base
- **Kubernetes-ready**: Complete K8s manifests in `k8s/` directory

## Key Data Structures
The core `Notification` struct handles complex nested JSON payloads with template-style field names:
- Fields use `{{media}}`, `{{request}}`, `{{issue}}`, `{{comment}}` as JSON keys
- Supports multiple notification types: media status updates, user requests, issue reports, comments
- Includes user identity data with Discord/Telegram integration fields

## Development Workflow
```bash
# Local development
go run main.go

# Build binary
go build -o jellynotifier .

# Build Docker image
docker build -t jellynotifier .

# Deploy to Kubernetes (in order)
kubectl apply -f k8s/01-namespace.yaml
kubectl apply -f k8s/02-deployment.yaml  
kubectl apply -f k8s/03-service.yaml
```

## Container Registry
Uses private registry: `registry.germainleignel.com/personal/jellynotifier:latest`

## Project Conventions
- **Logging**: Extensive structured logging for all notification fields
- **Error handling**: Simple HTTP status codes with logged details
- **Security**: Non-root container user, minimal Alpine base image
- **Resource limits**: Conservative CPU/memory limits for microservice deployment
- **Health checks**: Multiple probe types (liveness, readiness, startup) for robust K8s integration

## Common Modifications
- **Add new notification fields**: Update the nested structs in `Notification` type
- **Change logging format**: Modify the conditional logging blocks in `webhookHandler`
- **Update endpoints**: Add new HTTP handlers and register them in `main()`
- **Modify container resources**: Edit limits/requests in `k8s/02-deployment.yaml`
- **Change image tag**: Update image reference in deployment manifest

## Testing
- Use `/test` endpoint for quick connectivity verification
- Send POST requests to `/webhook` with sample JSON payloads
- Monitor logs for structured notification parsing output
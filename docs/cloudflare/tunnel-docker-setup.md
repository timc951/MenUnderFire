# Cloudflare Tunnel with Docker

Run Cloudflare Tunnel as a Docker container to connect to other containerized services.

## Prerequisites

- Docker and Docker Compose installed
- A Cloudflare account with a domain
- A tunnel created in the Cloudflare dashboard

## Creating a Tunnel

1. Go to **Cloudflare Zero Trust Dashboard**
2. Navigate to **Networks > Tunnels**
3. Click **Create a tunnel**
4. Choose **Cloudflared** as the connector
5. Name your tunnel and save
6. Copy the **tunnel token** for use in Docker

## Docker Compose Setup

### Basic Setup with Token

```yaml
# docker-compose.yml
services:
  cloudflared:
    image: cloudflare/cloudflared:latest
    command: tunnel run
    environment:
      - TUNNEL_TOKEN=<your-tunnel-token>
    networks:
      - web
    restart: unless-stopped

  mysite:
    image: nginx
    networks:
      - web
    # No ports needed - not exposed to host

networks:
  web:
```

### Configure Public Hostnames

In the Cloudflare dashboard, add public hostnames for your tunnel:

| Public Hostname | Service |
|----------------|---------|
| `mysite.yourdomain.com` | `http://mysite:80` |
| `api.yourdomain.com` | `http://api:3000` |

Use the **Docker service name** as the hostname (e.g., `mysite`, `api`), not `localhost`.

## Config File Approach

If you prefer managing configuration in files rather than the dashboard:

### Directory Structure

```
project/
├── docker-compose.yml
└── cloudflared/
    ├── config.yml
    └── credentials.json
```

### docker-compose.yml

```yaml
services:
  cloudflared:
    image: cloudflare/cloudflared:latest
    command: tunnel --config /etc/cloudflared/config.yml run
    volumes:
      - ./cloudflared:/etc/cloudflared:ro
    networks:
      - web
    restart: unless-stopped

  site1:
    image: nginx
    networks:
      - web

  site2:
    build: ./site2
    networks:
      - web

networks:
  web:
```

### cloudflared/config.yml

```yaml
tunnel: <tunnel-id>
credentials-file: /etc/cloudflared/credentials.json

ingress:
  - hostname: site1.yourdomain.com
    service: http://site1:80
  - hostname: site2.yourdomain.com
    service: http://site2:3000
  - hostname: api.yourdomain.com
    service: http://api:8080
    originRequest:
      noTLSVerify: true  # If your service uses self-signed certs
  - service: http_status:404  # Catch-all (required)
```

### Getting credentials.json

```bash
# On a machine with cloudflared installed
cloudflared tunnel login
cloudflared tunnel create my-tunnel

# Copy the credentials file
cp ~/.cloudflared/<tunnel-id>.json ./cloudflared/credentials.json
```

## Key Points

- **Use Docker service names** as hostnames (e.g., `http://mysite:80`), not `localhost`
- **All containers must be on the same Docker network**
- **No need to expose ports** on your services - the tunnel reaches them internally through the Docker network
- **One tunnel can serve multiple sites** by adding more hostnames in the ingress rules

## Multiple Sites Example

```yaml
services:
  cloudflared:
    image: cloudflare/cloudflared:latest
    command: tunnel run
    environment:
      - TUNNEL_TOKEN=${CLOUDFLARE_TUNNEL_TOKEN}
    networks:
      - web
    restart: unless-stopped

  portfolio:
    image: nginx
    volumes:
      - ./portfolio:/usr/share/nginx/html:ro
    networks:
      - web

  blog:
    image: ghost:5
    environment:
      url: https://blog.yourdomain.com
    networks:
      - web

  api:
    build: ./api
    environment:
      - DATABASE_URL=postgres://db:5432/myapp
    networks:
      - web
      - internal

  db:
    image: postgres:15
    volumes:
      - pgdata:/var/lib/postgresql/data
    networks:
      - internal  # Not on web network - not accessible via tunnel

networks:
  web:       # Services exposed through tunnel
  internal:  # Internal services only

volumes:
  pgdata:
```

## Environment Variables

Store your tunnel token in a `.env` file:

```bash
# .env
CLOUDFLARE_TUNNEL_TOKEN=eyJhIjoiYWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXoxMjM0NTY...
```

Then reference it in docker-compose.yml:
```yaml
environment:
  - TUNNEL_TOKEN=${CLOUDFLARE_TUNNEL_TOKEN}
```

## Troubleshooting

### Check tunnel status
```bash
docker compose logs cloudflared
```

### Verify container networking
```bash
# From inside the cloudflared container
docker compose exec cloudflared wget -qO- http://mysite:80
```

### Common issues

| Issue | Solution |
|-------|----------|
| "No ingress rules match" | Ensure the catch-all rule (`http_status:404`) is last |
| "Connection refused" | Check that services are on the same Docker network |
| "502 Bad Gateway" | Verify the service name and port are correct |
| Container keeps restarting | Check token is valid and not expired |

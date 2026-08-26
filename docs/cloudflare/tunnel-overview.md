# Cloudflare Tunnel Overview

Cloudflare Tunnel (formerly Argo Tunnel) allows you to expose local services to the internet without opening inbound ports on your firewall.

## Why It's Safe

- **No open inbound ports** - The tunnel creates an outbound-only connection to Cloudflare's edge, so you don't need to open ports on your firewall/router
- **No public IP exposure** - Your origin server's IP stays hidden behind Cloudflare
- **Built-in DDoS protection** - Traffic goes through Cloudflare's network
- **TLS encryption** - All traffic between Cloudflare and your origin is encrypted

## Cost

Cloudflare Tunnel is **free** on all plans, including the free tier.

You just need:
- A Cloudflare account (free)
- A domain on Cloudflare (free to add, you just need to own a domain)

### What Costs Money

| Feature | Cost |
|---------|------|
| Cloudflare Tunnel | Free |
| Cloudflare Access (first 50 users) | Free |
| Cloudflare Access (50+ users) | Paid (Zero Trust plans) |
| Load balancing across tunnels | Paid add-on |
| Advanced Zero Trust features | Paid tiers |

For a personal project or small site, you can run tunnels indefinitely at no cost.

## Basic Setup

### Install cloudflared

```bash
# Windows
winget install Cloudflare.cloudflared

# Mac
brew install cloudflared

# Linux (Debian/Ubuntu)
curl -L https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb -o cloudflared.deb
sudo dpkg -i cloudflared.deb
```

### Create and Run a Tunnel

```bash
# Authenticate with Cloudflare
cloudflared tunnel login

# Create a tunnel
cloudflared tunnel create my-tunnel

# Route traffic to your domain
cloudflared tunnel route dns my-tunnel mysite.yourdomain.com

# Run the tunnel (points to your local service)
cloudflared tunnel run --url http://localhost:3000 my-tunnel
```

## Configuration File

For more control, use a config file at `~/.cloudflared/config.yml`:

```yaml
tunnel: <your-tunnel-id>
credentials-file: ~/.cloudflared/<tunnel-id>.json

ingress:
  - hostname: mysite.yourdomain.com
    service: http://localhost:3000
  - service: http_status:404  # Catch-all for unmatched requests
```

Then run with:
```bash
cloudflared tunnel run
```

## Security Best Practices

1. **Use Cloudflare Access** - Add authentication in front of your tunnel for admin panels or sensitive apps
2. **Limit tunnel scope** - Only expose specific paths/services, not your entire network
3. **Run as a service** - Use `cloudflared service install` for production rather than manual runs
4. **Use a config file** - Define ingress rules explicitly rather than using command-line flags
5. **Enable logging** - Monitor tunnel traffic for suspicious activity
6. **Keep cloudflared updated** - Regularly update to get security patches

## Running as a System Service

For production use, install as a service:

```bash
# Install the service
sudo cloudflared service install

# Start the service
sudo systemctl start cloudflared

# Enable on boot
sudo systemctl enable cloudflared
```

The service will use the config file at `/etc/cloudflared/config.yml`.

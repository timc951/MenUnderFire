# Resend Email Setup

## Overview

This project uses [Resend](https://resend.com) for transactional emails. Resend provides:
- Simple REST API for sending emails
- React Email integration for beautiful templates
- Domain authentication (SPF, DKIM)
- 3,000 emails/month free tier
- Email delivery tracking and analytics

## Quick Start

1. Create a Resend account at https://resend.com
2. Add and verify your domain
3. Create an API key
4. Add the API key to your backend environment
5. Start sending emails

## Account Setup

### 1. Create Account

1. Go to https://resend.com
2. Click **Sign Up**
3. Create account with email or GitHub
4. Verify your email address

### 2. Add Your Domain

To send emails from your own domain (recommended for production):

1. Go to **Domains** in the Resend dashboard
2. Click **Add Domain**
3. Enter your domain (e.g., `menunderfire.com`)
4. Choose your region (US or EU based on user location)

### 3. Configure DNS Records

Resend will provide DNS records to add to your domain registrar:

| Type | Name | Value | Purpose |
|------|------|-------|---------|
| TXT | `_resend` | Verification code | Domain verification |
| TXT | Host varies | SPF record | Email authentication |
| CNAME | `resend._domainkey` | DKIM key | Email signing |

#### Adding DNS Records by Provider

**Cloudflare:**
1. Go to your domain → DNS → Records
2. Click **Add record** for each entry
3. Disable proxy (gray cloud) for CNAME records

**Namecheap:**
1. Go to Domain List → Manage → Advanced DNS
2. Add each record under Host Records

**GoDaddy:**
1. Go to DNS Management
2. Add each record type

**Route 53 (AWS):**
1. Go to Hosted Zones → Your domain
2. Create record for each entry

### 4. Wait for Verification

- DNS propagation takes 5 minutes to 48 hours
- Resend dashboard shows verification status
- You'll receive an email when verified

## API Key Setup

### 1. Create API Key

1. Go to **API Keys** in the Resend dashboard
2. Click **Create API Key**
3. Name it (e.g., `MenUnderFire Production`)
4. Select permission: **Full access** or **Sending access**
5. Optionally restrict to specific domain
6. Copy the key (shown only once!)

### 2. API Key Types

| Permission | Use Case |
|------------|----------|
| Full access | Development, CI/CD pipelines |
| Sending access | Production (more secure) |

## Environment Variables

### Backend

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `RESEND_API_KEY` | Yes | Resend API key | `re_123abc...` |
| `EMAIL_FROM_ADDRESS` | Yes | Sender email address | `noreply@menunderfire.com` |
| `EMAIL_FROM_NAME` | No | Sender display name | `Men Under Fire` |

Example backend `.env`:
```bash
RESEND_API_KEY=re_123abc456def...
EMAIL_FROM_ADDRESS=noreply@menunderfire.com
EMAIL_FROM_NAME=Men Under Fire
```

## Development Mode

For development without a verified domain:

1. Use Resend's test domain: `onboarding@resend.dev`
2. Emails can only be sent to the account owner's email
3. Create a separate API key for development

Example development `.env`:
```bash
RESEND_API_KEY=re_test_123abc...
EMAIL_FROM_ADDRESS=onboarding@resend.dev
EMAIL_FROM_NAME=Men Under Fire (Dev)
```

## Email Types

Configure these transactional emails:

| Email | Trigger | Priority |
|-------|---------|----------|
| Invitation | User invited to organization | High |
| Welcome | User completes signup | Medium |
| Password Reset | User requests reset | High |
| Group Notification | New message in group | Medium |
| Report Summary | Weekly/monthly reports | Low |

## Go Integration

Install the Resend Go SDK:

```bash
go get github.com/resend/resend-go/v2
```

Example usage:

```go
import "github.com/resend/resend-go/v2"

client := resend.NewClient("re_123abc...")

params := &resend.SendEmailRequest{
    From:    "Men Under Fire <noreply@menunderfire.com>",
    To:      []string{"user@example.com"},
    Subject: "Welcome to Men Under Fire",
    Html:    "<h1>Welcome!</h1><p>Thanks for joining.</p>",
}

sent, err := client.Emails.Send(params)
```

## Testing Emails

### Test Mode

1. Use development API key
2. Emails only go to account owner
3. No charges for test emails

### Resend Testing Tools

- **Logs**: View all sent emails in dashboard
- **Preview**: Test email rendering before sending
- **Webhooks**: Track delivery, opens, clicks

## Webhooks (Optional)

Configure webhooks to track email events:

1. Go to **Webhooks** in Resend dashboard
2. Add endpoint URL (e.g., `https://api.menunderfire.com/webhooks/email`)
3. Select events:
   - `email.sent`
   - `email.delivered`
   - `email.bounced`
   - `email.complained`
4. Copy signing secret for verification

## Rate Limits

| Tier | Emails/Month | Emails/Second |
|------|--------------|---------------|
| Free | 3,000 | 1 |
| Pro ($20/mo) | 50,000 | 10 |
| Enterprise | Unlimited | Custom |

## Troubleshooting

### "Domain not verified" error
- Check DNS records are correctly added
- Wait up to 48 hours for propagation
- Use `dig` or online DNS checker to verify records

### Emails going to spam
1. Ensure SPF and DKIM records are set
2. Use a professional from address (not personal email)
3. Include unsubscribe link in marketing emails
4. Avoid spam trigger words in subject

### "Invalid API key" error
- Verify key starts with `re_`
- Check for extra spaces when copying
- Ensure key has correct permissions

### Emails not sending in development
- Test emails only send to account owner's email
- Verify API key is for correct environment

## Security Best Practices

1. **Never commit API keys** - Use environment variables
2. **Use sending-only keys in production** - Limits damage if leaked
3. **Rotate keys periodically** - Create new key, update env, delete old
4. **Restrict keys to domains** - Prevents misuse

## Resources

- [Resend Documentation](https://resend.com/docs)
- [Resend Go SDK](https://github.com/resend/resend-go)
- [React Email](https://react.email) - Build email templates
- [Email Testing Tools](https://resend.com/docs/dashboard/emails/logs)

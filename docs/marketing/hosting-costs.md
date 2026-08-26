# Hosting & Infrastructure Costs

## Overview

This document breaks down the expected hosting and operational costs for Men Under Fire at various scales.

---

## Cost Breakdown by Service

### Starting Phase (0-500 users)

| Service | Provider Options | Free Tier | Paid Tier |
|---------|-----------------|-----------|-----------|
| **Frontend Hosting** | Vercel, Netlify, Cloudflare Pages | Free | $20/month |
| **Database** | Supabase, PlanetScale, Neon | Free (generous) | $25-50/month |
| **Authentication** | Clerk, Auth0, Supabase Auth | Free to 50,000 MAU | $20-100/month |
| **Email** | Resend, SendGrid, Postmark | Free tier | $20/month |
| **Domain** | Namecheap, Cloudflare | - | $15/year |
| **SSL** | Included with hosting | Free | Free |

**Total Starting Cost: $0-50/month**

### Growth Phase (500-5,000 users)

| Service | Estimated Cost |
|---------|---------------|
| Frontend Hosting | $20/month |
| Database | $50-100/month |
| Authentication | $50-150/month |
| Email | $30-50/month |
| CDN/Assets | $20/month |
| Monitoring | $20/month |
| Backups | $20/month |

**Total Growth Cost: $150-400/month**

### Scale Phase (5,000+ users)

| Service | Estimated Cost |
|---------|---------------|
| Frontend Hosting | $50/month |
| Database | $200-500/month |
| Authentication | $200-400/month |
| Email | $100/month |
| CDN/Assets | $50/month |
| Monitoring/Logging | $100/month |
| Backups | $50/month |
| Support tools | $50/month |

**Total Scale Cost: $500-1,500/month**

---

## Recommended Stack (Cost-Optimized)

### Frontend
- **Vercel** (free tier generous, scales well)
- Alternative: Cloudflare Pages (very generous free tier)

### Database
- **Supabase** (PostgreSQL, free tier includes 500MB, auth built-in)
- Alternative: PlanetScale (MySQL, generous free tier)

### Authentication
- **Clerk** (50,000 MAU free, excellent DX)
- Alternative: Supabase Auth (included with Supabase)
- Alternative: Auth0 (25,000 MAU free)

### Email
- **Resend** (3,000 emails/month free)
- Alternative: SendGrid (100 emails/day free)

### File Storage
- **Supabase Storage** (1GB free)
- Alternative: Cloudflare R2 (10GB free, no egress fees)

---

## Monthly Cost by User Count

| Active Users | Estimated Monthly Cost |
|--------------|----------------------|
| 0-100 | $0-20 |
| 100-500 | $20-50 |
| 500-1,000 | $50-100 |
| 1,000-2,500 | $100-200 |
| 2,500-5,000 | $200-400 |
| 5,000-10,000 | $400-800 |

---

## Break-Even Analysis

### Donation Model

| Monthly Cost | Donors Needed @ $5/mo | Donors Needed @ $10/mo |
|--------------|----------------------|------------------------|
| $50 | 10 donors | 5 donors |
| $100 | 20 donors | 10 donors |
| $200 | 40 donors | 20 donors |
| $500 | 100 donors | 50 donors |

### Organization Model

| Monthly Cost | Churches Needed @ $79/mo | Churches Needed @ $149/mo |
|--------------|-------------------------|--------------------------|
| $100 | 2 churches | 1 church |
| $500 | 7 churches | 4 churches |
| $1,000 | 13 churches | 7 churches |

---

## Cost Optimization Tips

1. **Stay on free tiers as long as possible**
   - Most services have generous free tiers
   - Don't over-provision early

2. **Use Supabase for multiple services**
   - Database + Auth + Storage in one
   - Reduces vendor complexity

3. **Defer premium features**
   - Don't pay for monitoring until you need it
   - Use free Sentry/LogRocket tiers initially

4. **Annual billing discounts**
   - Most services offer 10-20% off for annual
   - Only commit once stable

5. **Cloudflare for everything possible**
   - Free CDN, DDoS protection, DNS
   - R2 storage has no egress fees

---

## Hidden Costs to Plan For

| Cost | Estimate | Notes |
|------|----------|-------|
| Domain renewal | $15/year | Annual |
| SSL certificates | $0 | Usually free with hosting |
| Error monitoring | $0-30/month | Sentry free tier is good |
| Analytics | $0 | Plausible/Umami self-hosted or free tiers |
| Backup storage | $5-20/month | Important for data safety |
| Development tools | $0-50/month | GitHub free, etc. |

---

## Authentication Provider Comparison

**MAU = Monthly Active Users** - unique users who authenticate at least once per billing month.

### Hosted Solutions (Recommended)

| Provider | Free Tier | Paid Pricing | Best For |
|----------|-----------|--------------|----------|
| **Clerk** | 50,000 MAU | Pro: $20/mo + usage | Best DX, React/Next.js apps |
| **Auth0** | 25,000 MAU | Essentials: $35/mo+ | Enterprise features needed |
| **SuperTokens** | 5,000 MAU (cloud) | $0.02/MAU | Self-host option, predictable costs |
| **Supabase Auth** | Included with Supabase | Included | Already using Supabase |
| **Logto** | 50,000 MAU | Variable | Open-source alternative |

### Self-Hosted Options

| Provider | Cost | Effort | Notes |
|----------|------|--------|-------|
| **SuperTokens** | Free (unlimited) | Medium | Open-source, good docs |
| **Keycloak** | Free (but hosting costs) | High | Enterprise SSO, needs 2GB+ RAM |
| **Stack Auth** | Free (MIT/AGPL) | Medium | Modern, developer-friendly |
| **Hanko** | Free | Medium | Passkey-first approach |

### Self-Hosted Keycloak (Not Recommended)

Why self-hosting Keycloak is typically overkill:

| Factor | Details |
|--------|---------|
| **Memory Required** | Minimum 750MB, recommended 2GB for production |
| **Hosting Cost** | ~$25-50/mo (container + database) on Render/Railway |
| **Operational Burden** | Upgrades, security patches, backups, SSL, monitoring |
| **Best Use Case** | Enterprise SSO, SAML federation, complex LDAP |

Unless you need specific Keycloak features (SAML, LDAP, complex federation), hosted solutions like Clerk provide better value with zero infrastructure management.

### Recommendation

For Men Under Fire: **Clerk** or **Supabase Auth**

- **Clerk**: Best if auth is separate from database. 50,000 MAU free covers growth phase.
- **Supabase Auth**: Best if already using Supabase. Simplifies stack.

---

## Recommendation

**Start lean, scale smart:**

1. Use Supabase for database + auth + storage (free tier)
2. Deploy frontend to Vercel (free tier)
3. Use Resend for email (free tier)
4. Total starting cost: **~$15/year** (domain only)

As you grow past 500 active users, budget for $100-200/month and plan for 10+ paying organizations or 30+ regular donors to cover costs.

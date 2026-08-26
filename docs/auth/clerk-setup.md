# Clerk.com Authentication Setup

## Overview

This project uses [Clerk.com](https://clerk.com) for authentication. Clerk provides:
- User authentication (email/password, social providers, passkeys)
- JWT token generation for API access
- User management dashboard
- 50,000 MAU free tier

## Quick Start

1. Create a Clerk application at https://dashboard.clerk.com
2. Copy your publishable key to `frontend/.env`
3. Configure your backend with the JWKS URL
4. Run the application

## Environment Variables

### Frontend

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `VITE_AUTH_PUBLISHABLE_KEY` | Yes | Clerk publishable key | `pk_test_abc123...` |
| `VITE_API_URL` | Yes | Backend API URL | `http://localhost:7001/api` |

Example `frontend/.env`:
```bash
VITE_API_URL=http://localhost:7001/api
VITE_AUTH_PUBLISHABLE_KEY=pk_test_abc123...
```

### Backend

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `AUTH_ISSUER` | Yes | JWT issuer URL | `https://your-app.clerk.accounts.dev` |
| `AUTH_JWKS_URL` | Yes | JWKS endpoint for token validation | `https://your-app.clerk.accounts.dev/.well-known/jwks.json` |
| `AUTH_AUDIENCE` | No | JWT audience claim (optional) | `your-api-identifier` |

Example backend `.env`:
```bash
AUTH_ISSUER=https://your-app.clerk.accounts.dev
AUTH_JWKS_URL=https://your-app.clerk.accounts.dev/.well-known/jwks.json
AUTH_AUDIENCE=
```

## Clerk Dashboard Setup

### 1. Create Application

1. Go to https://dashboard.clerk.com
2. Click "Create application"
3. Choose authentication methods:
   - Email (recommended)
   - Google (optional)
   - Other social providers as needed

### 2. Get Your Keys

1. In your Clerk dashboard, go to **API Keys**
2. Copy the **Publishable key** (starts with `pk_`)
3. Note the **Frontend API URL** (e.g., `https://your-app.clerk.accounts.dev`)

### 3. Configure Redirect URLs

1. Go to **Paths** in the dashboard
2. Set allowed redirect URLs:
   - Development: `http://localhost:8001`
   - Production: `https://your-domain.com`

### 4. JWT Template (Required)

The backend needs `email` and `name` claims in the JWT. Clerk's default JWT doesn't include these.

1. Go to **JWT Templates**
2. Click **New template** → **Blank**
3. Name it exactly: `menunderfire`
4. Set the claims:
```json
{
  "email": "{{user.primary_email_address}}",
  "name": "{{user.full_name}}"
}
```
5. Click **Save**

The frontend automatically uses this template via `getToken({ template: 'menunderfire' })`.

## Finding Your JWKS URL

Your JWKS URL is derived from your Clerk Frontend API:
- Frontend API: `https://your-app.clerk.accounts.dev`
- JWKS URL: `https://your-app.clerk.accounts.dev/.well-known/jwks.json`

## Switching Auth Providers

The environment variables are designed to be generic. To switch to a different provider:

1. Update `VITE_AUTH_PUBLISHABLE_KEY` or equivalent
2. Update `AUTH_ISSUER` and `AUTH_JWKS_URL` for backend
3. Modify `AuthProvider.tsx` to use the new provider's SDK
4. Update `useAuth.ts` to use the new provider's hooks

## Troubleshooting

### "Missing publishable key" error
Ensure `VITE_AUTH_PUBLISHABLE_KEY` is set in your `.env` file and restart the dev server.

### Token validation fails on backend
1. Verify `AUTH_ISSUER` matches your Clerk Frontend API URL
2. Verify `AUTH_JWKS_URL` is accessible (test with curl)
3. Check that the JWT is being passed in the `Authorization: Bearer <token>` header

### User sync fails
The frontend syncs users to the backend via `GET /users/me`. Ensure:
1. Backend is running and accessible
2. `VITE_API_URL` is correct
3. Backend JWT validation is configured with Clerk's JWKS

## Resources

- [Clerk Documentation](https://clerk.com/docs)
- [Clerk React SDK](https://clerk.com/docs/quickstarts/react)
- [JWT Templates](https://clerk.com/docs/backend-requests/making/jwt-templates)

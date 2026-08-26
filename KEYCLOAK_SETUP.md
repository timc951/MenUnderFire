# Keycloak Setup for Men Under Fire

This guide walks through setting up Keycloak as the identity provider for the Men Under Fire application.

## Prerequisites



- Keycloak running at http://localhost:9001
- Admin credentials: `admin` / `admin`

## Step 1: Create a New Realm

1. Log in to Keycloak Admin Console at http://localhost:9001
2. Click the dropdown in the top-left corner (shows "master")
3. Click **Create realm**
4. Enter the following:
   - **Realm name**: `menunderfire`
5. Click **Create**

## Step 2: Create the Frontend Client

1. In the left sidebar, click **Clients**
2. Click **Create client**
3. **General Settings**:
   - **Client type**: OpenID Connect
   - **Client ID**: `menunderfire-frontend`
   - Click **Next**
4. **Capability config**:
   - **Client authentication**: OFF (this is a public SPA client)
   - **Authorization**: OFF
   - **Authentication flow**: Check only **Standard flow** and **Direct access grants**
   - Click **Next**
5. **Login settings**:
   - **Root URL**: `http://localhost:8001`
   - **Home URL**: `http://localhost:8001`
   - **Valid redirect URIs**: `http://localhost:8001/*`
   - **Valid post logout redirect URIs**: `http://localhost:8001/*`
   - **Web origins**: `http://localhost:8001`
   - Click **Save**

## Step 3: Configure Client Scopes (Optional)

The default scopes include `openid`, `profile`, and `email` which match what the app currently uses. No changes needed unless you want custom claims.

## Step 4: Create Roles

1. In the left sidebar, click **Realm roles**
2. Click **Create role**
3. Create the following roles (click **Create role** for each):

| Role Name | Description |
|-----------|-------------|
| `user` | Standard authenticated user |
| `org_admin` | Organization administrator |
| `site_admin` | Site-wide administrator |

## Step 5: Create a Test User

1. In the left sidebar, click **Users**
2. Click **Add user**
3. Fill in the details:
   - **Username**: `testuser`
   - **Email**: `testuser@example.com`
   - **Email verified**: ON
   - **First name**: `Test`
   - **Last name**: `User`
4. Click **Create**
5. Go to the **Credentials** tab
6. Click **Set password**
   - **Password**: `testpass123`
   - **Temporary**: OFF
   - Click **Save**
7. Go to the **Role mapping** tab
8. Click **Assign role**
9. Select the `user` role and click **Assign**

## Step 6: Get Configuration Values

After setup, you'll need these values for your frontend:

1. Go to **Realm settings** in the left sidebar
2. Click the **General** tab
3. At the bottom, click **OpenID Endpoint Configuration**

Key values for your app:

| Setting | Value |
|---------|-------|
| Keycloak URL | `http://localhost:9001` |
| Realm | `menunderfire` |
| Client ID | `menunderfire-frontend` |
| Authorization Endpoint | `http://localhost:9001/realms/menunderfire/protocol/openid-connect/auth` |
| Token Endpoint | `http://localhost:9001/realms/menunderfire/protocol/openid-connect/token` |
| UserInfo Endpoint | `http://localhost:9001/realms/menunderfire/protocol/openid-connect/userinfo` |

## Step 7: Update Frontend Environment Variables

Update your `.env` or `.env.local` file in the frontend directory:

```env
# Keycloak Configuration
VITE_KEYCLOAK_URL=http://localhost:9001
VITE_KEYCLOAK_REALM=menunderfire
VITE_KEYCLOAK_CLIENT_ID=menunderfire-frontend
```

## Step 8: Update AuthProvider (Code Change Required)

The current frontend uses `@auth0/auth0-react`. To use Keycloak, you'll need to switch to a Keycloak adapter such as:

- `@react-keycloak/web` - React bindings for keycloak-js
- `react-oidc-context` - Generic OIDC library that works with Keycloak

See the separate migration guide for code changes.

---

## Additional Configuration

### Enable User Registration (Optional)

1. Go to **Realm settings**
2. Click the **Login** tab
3. Enable **User registration**
4. Click **Save**

### Configure Email (Optional)

1. Go to **Realm settings**
2. Click the **Email** tab
3. Configure your SMTP server settings
4. Click **Save**

### Add Social Login (Optional)

1. Go to **Identity providers** in the left sidebar
2. Click **Add provider**
3. Select Google, GitHub, etc.
4. Configure with your OAuth credentials

---

## Production Considerations

Before deploying to production:

1. **Change admin password** - Update the default admin credentials
2. **Enable HTTPS** - Configure SSL/TLS for Keycloak
3. **Use strong passwords** - Update database and admin passwords
4. **Configure hostname** - Set `KC_HOSTNAME` environment variable
5. **Switch to production mode** - Change `start-dev` to `start` in docker-compose
6. **Set up backups** - Configure PostgreSQL backups for the keycloak database

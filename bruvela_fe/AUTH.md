# Authentication - Bruvela Bakehouse

## Overview

Frontend menggunakan JWT-based authentication yang terintegrasi dengan backend API.

## Features

- ✅ Login dengan email & password
- ✅ JWT token storage di cookies
- ✅ Auto-redirect ke login jika belum authenticated
- ✅ Auto-redirect ke dashboard jika sudah login
- ✅ Logout functionality
- ✅ User info display di header
- ✅ Protected routes (semua halaman kecuali login)

## Test Account

```
Email: admin@bruvela.com
Password: admin123
```

## API Endpoints

### Login
```
POST /api/v1/auth/login
Body: {
  "email": "admin@bruvela.com",
  "password": "admin123"
}

Response: {
  "token": "eyJhbGc...",
  "user": {
    "id": "uuid",
    "email": "admin@bruvela.com",
    "name": "Admin",
    "role": "admin"
  }
}
```

### Get Current User
```
GET /api/v1/auth/me
Headers: {
  "Authorization": "Bearer <token>"
}

Response: {
  "id": "uuid",
  "email": "admin@bruvela.com",
  "name": "Admin",
  "role": "admin"
}
```

## Usage

### Login Page

Halaman login tersedia di `/login` dengan UI yang clean dan modern.

Features:
- Email & password input
- Remember me checkbox
- Auto-fill untuk testing (admin@bruvela.com / admin123)
- Loading state saat login
- Error handling dengan toast notification

### useAuth Composable

```typescript
const { 
  login,      // Login function
  logout,     // Logout function
  getMe,      // Get current user
  isAuthenticated, // Computed boolean
  user,       // Current user state
  token       // JWT token
} = useAuth()
```

### Login Example

```typescript
const { login } = useAuth()

const handleLogin = async () => {
  const { data, error } = await login(email, password)
  
  if (error) {
    // Handle error
    toast.add({ title: 'Error', description: error.error })
    return
  }
  
  // Success - auto redirect to dashboard
  router.push('/')
}
```

### Logout Example

```typescript
const { logout } = useAuth()

const handleLogout = () => {
  logout() // Clears token & user, redirects to login
}
```

### Check Authentication

```typescript
const { isAuthenticated, user } = useAuth()

// In template
<div v-if="isAuthenticated">
  Welcome, {{ user?.name }}
</div>
```

## Protected Routes

Semua routes kecuali `/login` dilindungi oleh auth middleware.

Middleware akan:
1. Check apakah user sudah login (ada token)
2. Jika belum login → redirect ke `/login`
3. Jika sudah login di `/login` → redirect ke `/`

## Token Storage

Token disimpan di cookies dengan:
- Name: `auth_token`
- Max Age: 7 days
- HttpOnly: false (untuk client-side access)
- Secure: true (production only)

## Auto-load User Data

Saat aplikasi dimuat dan user sudah authenticated:
1. Plugin auth akan auto-call `getMe()`
2. User data akan di-load ke state
3. User info ditampilkan di header

## Error Handling

### 401 Unauthorized

Jika API return 401:
1. Token invalid/expired
2. User akan di-logout otomatis
3. Redirect ke login page

### Login Failed

Jika login gagal:
- Toast notification dengan error message
- Form tetap bisa digunakan
- Password field di-clear

## Security

- Token disimpan di cookies (lebih aman dari localStorage)
- Password tidak pernah disimpan
- Token dikirim via Authorization header
- Auto-logout saat token expired

## Files

### Core Files
- `composables/useAuth.ts` - Auth composable
- `pages/login.vue` - Login page
- `plugins/auth.ts` - Global auth middleware & auto-load user

### Integration
- `composables/useApi.ts` - Auto-attach token ke API calls
- `layouts/default.vue` - Display user info & logout button

## Testing

1. **Start Backend**
```bash
cd /Users/yogatama.egiantoro/Documents/Development/bruvela-be
make run
```

2. **Start Frontend**
```bash
cd bruvela_fe
npm run dev
```

3. **Test Login**
- Buka http://localhost:3000
- Auto-redirect ke `/login`
- Login dengan `admin@bruvela.com` / `admin123`
- Redirect ke dashboard
- User info muncul di header

4. **Test Logout**
- Klik avatar di header
- Klik "Logout"
- Redirect ke login page

5. **Test Protected Routes**
- Logout dulu
- Coba akses http://localhost:3000/orders
- Auto-redirect ke login

## Troubleshooting

### Token tidak tersimpan
- Check browser cookies
- Pastikan backend return token dengan benar

### Auto-redirect loop
- Clear cookies
- Restart dev server

### 401 di semua API calls
- Check token di cookies
- Pastikan backend JWT secret sama
- Check token format di Authorization header

## Next Steps

- [ ] Add "Forgot Password" functionality
- [ ] Add "Remember Me" persistent login
- [ ] Add refresh token mechanism
- [ ] Add role-based access control
- [ ] Add session timeout warning

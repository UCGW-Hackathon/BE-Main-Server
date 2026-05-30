---

## 📑 Table of Contents

### Part A — General & Authentication
- [0. General Information](#0-general-information)
- [1. Authentication](#1-authentication)

### Part B — User (Customer) Side
- [2. User Profile](#2-user-profile)
- [3. Home (User)](#3-home-user)
- [4. Service Categories](#4-service-categories)
- [5. Nearby Workers (Tukang Dekat)](#5-nearby-workers-tukang-dekat)
- [6. Worker Detail (Detail Tukang)](#6-worker-detail-detail-tukang)
- [7. Orders — User Side](#7-orders--user-side)
- [8. Tracking — User Side](#8-tracking--user-side)
- [9. Purchase Tracking — User Side](#9-purchase-tracking--user-side)
- [10. Chat — User Side](#10-chat--user-side)
- [11. Rating & Review — User Side](#11-rating--review--user-side)
- [12. Invoice & Payment — User Side](#12-invoice--payment--user-side)
- [13. Knowledge / FAQ / Guide](#13-knowledge--faq--guide)
- [14. Notifications](#14-notifications)

### Part C — Worker (Tukang) Side
- [15. Worker Profile & Verification](#15-worker-profile--verification)
- [16. Home (Worker)](#16-home-worker)
- [17. Incoming Orders — Worker Side](#17-incoming-orders--worker-side)
- [18. Order Management — Worker Side](#18-order-management--worker-side)
- [19. Purchase Management — Worker Side (AI-Assisted)](#19-purchase-management--worker-side-ai-assisted)
- [20. Chat — Worker Side](#20-chat--worker-side)
- [21. Rating — Worker Side (Rate Customer)](#21-rating--worker-side-rate-customer)
- [22. History & Statistics — Worker Side](#22-history--statistics--worker-side)
- [23. Wallet — Worker Side](#23-wallet--worker-side)
- [24. Worker Location Updates](#24-worker-location-updates)

### Part D — Shared
- [25. WebSocket Events](#25-websocket-events)
- [26. Enums & Shared Schemas](#26-enums--shared-schemas)
- [27. Error Responses](#27-error-responses)
- [28. Endpoint Summary Table](#28-endpoint-summary-table)

---

## 0. General Information

### 0.1 Base URL

| Environment | URL |
|---|---|
| Production | `https://api.situkang.id/v1` |
| Staging | `https://staging-api.situkang.id/v1` |
| Development | `http://localhost:3000/v1` |

### 0.2 Content Types

| Content-Type | Usage |
|---|---|
| `application/json` | Default request/response body |
| `multipart/form-data` | File upload (avatar, photos, receipts) |

### 0.3 Authentication — RBAC via JWT

SiTukang uses a **single backend** for both User and Worker. Role differentiation is handled via the `role` field embedded in the JWT access token. Middleware checks this field to authorize or deny access to role-specific endpoints.

```yaml
securitySchemes:
  BearerAuth:
    type: http
    scheme: bearer
    bearerFormat: JWT
    description: |
      JWT token containing user role.
      Include in header: Authorization: Bearer <token>
```

**JWT Payload Structure:**

```json
{
  "sub": "550e8400-e29b-41d4-a716-446655440000",
  "role": "user",
  "email": "budi@email.com",
  "full_name": "Budi Santoso",
  "iat": 1698230400,
  "exp": 1698234000,
  "jti": "unique-token-id"
}
```

**Middleware RBAC Flow:**

```
Request
  → Extract Bearer token from Authorization header
  → Verify JWT signature & expiry
  → Decode payload → extract `role`
  → Match role against endpoint requirement:
      ├── role == "user"   → User endpoints ✅ | Worker endpoints ❌
      ├── role == "worker" → Worker endpoints ✅ | User endpoints ❌
      ├── role == "admin"  → All endpoints ✅
      └── Invalid/expired  → 401 Unauthorized
```

### 0.4 Rate Limiting

All responses include rate limiting headers:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1698230400
```

### 0.5 Pagination Convention

All list endpoints use consistent pagination:

```json
{
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 150,
    "total_pages": 15
  }
}
```

### 0.6 Common HTTP Status Codes

| Code | Usage |
|---|---|
| `200` | Success (GET, PUT, PATCH) |
| `201` | Resource created (POST) |
| `204` | Success, no content (DELETE) |
| `400` | Bad Request — invalid parameters |
| `401` | Unauthorized — invalid/expired token |
| `403` | Forbidden — role mismatch |
| `404` | Not Found |
| `409` | Conflict — duplicate resource |
| `422` | Unprocessable Entity — business validation failed |
| `429` | Rate Limited |
| `500` | Internal Server Error |

---

## 1. Authentication

> Auth endpoints are **role-agnostic** — both user and worker use the same endpoints. The `role` is set during registration and embedded in the JWT.

---

### 1.1 `POST /auth/register`

> Register a new account (user or worker).

**Request Body:**

```json
{
  "full_name": "Budi Santoso",
  "email": "budi@email.com",
  "phone": "+6281234567890",
  "password": "SecureP@ss123",
  "password_confirmation": "SecureP@ss123",
  "role": "user",
  "latitude": -7.257500,
  "longitude": 112.752100,
  "address": "Jl. Merdeka No. 12, Surabaya"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `full_name` | string | ✅ | Full name (max 255 chars) |
| `email` | string | ✅ | Unique email address |
| `phone` | string | ✅ | Phone number with country code |
| `password` | string | ✅ | Min 8 chars, must include uppercase, lowercase, number |
| `password_confirmation` | string | ✅ | Must match password |
| `role` | enum | ✅ | `user` or `worker` |
| `latitude` | number | ❌ | Current latitude |
| `longitude` | number | ❌ | Current longitude |
| `address` | string | ❌ | Full address |

**Response `201 Created`:**

```json
{
  "status": "success",
  "message": "Registrasi berhasil",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "full_name": "Budi Santoso",
    "email": "budi@email.com",
    "phone": "+6281234567890",
    "role": "user",
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

**Error `409 Conflict`:**

```json
{
  "status": "error",
  "message": "Email sudah terdaftar",
  "error_code": "EMAIL_ALREADY_EXISTS"
}
```

---

### 1.2 `POST /auth/login`

> Login for any role.

**Request Body:**

```json
{
  "email": "budi@email.com",
  "password": "SecureP@ss123"
}
```

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "full_name": "Budi Santoso",
    "role": "user",
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

**Error `401 Unauthorized`:**

```json
{
  "status": "error",
  "message": "Email atau password salah",
  "error_code": "INVALID_CREDENTIALS"
}
```

---

### 1.3 `POST /auth/refresh`

> Refresh access token using refresh token. Old refresh token is invalidated (rotation).

**Request Body:**

```json
{
  "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4..."
}
```

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "bmV3IHJlZnJlc2ggdG9rZW4...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

---

### 1.4 `POST /auth/logout`

> Logout and invalidate refresh token.

🔒 **Auth Required** — Any role

**Request Headers:**

```
Authorization: Bearer <access_token>
```

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Berhasil logout"
}
```

---

### 1.5 `POST /auth/forgot-password`

> Send password reset link to email.

**Request Body:**

```json
{
  "email": "budi@email.com"
}
```

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Link reset password telah dikirim ke email Anda"
}
```

---

### 1.6 `POST /auth/reset-password`

> Reset password using token from email.

**Request Body:**

```json
{
  "token": "reset-token-from-email",
  "password": "NewSecureP@ss456",
  "password_confirmation": "NewSecureP@ss456"
}
```

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Password berhasil direset"
}
```

---

## 2. User Profile

> Endpoints for user profile management. Used in Profile tab and Home screen greeting.

---

### 2.1 `GET /users/me`

> Get current authenticated user's profile.

🔒 **Auth Required** — Role: `user`, `worker`, `admin`

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "full_name": "Budi Santoso",
    "email": "budi@email.com",
    "phone": "+6281234567890",
    "role": "user",
    "avatar_url": "https://cdn.handydirect.id/avatars/budi_001.jpg",
    "address": "Jl. Merdeka No. 12, Surabaya",
    "latitude": -7.257500,
    "longitude": 112.752100,
    "is_active": true,
    "email_verified_at": "2023-01-15T10:00:00Z",
    "phone_verified_at": "2023-01-15T10:05:00Z",
    "last_login_at": "2026-05-30T08:00:00Z",
    "created_at": "2023-01-15T10:00:00Z",
    "active_orders_count": 1
  }
}
```

---

### 2.2 `PUT /users/me`

> Update user profile info.

🔒 **Auth Required** — Role: `user`, `worker`, `admin`

**Request Body:**

```json
{
  "full_name": "Budi Santoso",
  "phone": "+6281234567890",
  "address": "Jl. Merdeka No. 14, Surabaya",
  "latitude": -7.257600,
  "longitude": 112.752200
}
```

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Profil berhasil diperbarui",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "full_name": "Budi Santoso",
    "phone": "+6281234567890",
    "address": "Jl. Merdeka No. 14, Surabaya",
    "latitude": -7.257600,
    "longitude": 112.752200,
    "updated_at": "2026-05-30T08:30:00Z"
  }
}
```

---

### 2.3 `PUT /users/me/avatar`

> Upload or update profile photo.

🔒 **Auth Required** — Role: `user`, `worker`, `admin`
📎 **Content-Type:** `multipart/form-data`

| Field | Type | Required | Description |
|---|---|---|---|
| `avatar` | file | ✅ | Image file (jpg/png, max 5MB) |

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Foto profil berhasil diperbarui",
  "data": {
    "avatar_url": "https://cdn.handydirect.id/avatars/budi_001_v2.jpg"
  }
}
```

---

### 2.4 `PUT /users/me/location`

> Update current location. Used for "Lokasi saat ini" display on Home screen.

🔒 **Auth Required** — Role: `user`, `worker`, `admin`

**Request Body:**

```json
{
  "latitude": -7.257500,
  "longitude": 112.752100,
  "address": "Jl. Merdeka No. 12, Surabaya"
}
```

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Lokasi berhasil diperbarui"
}
```



---

## 3. Home (User)

> Aggregated data for the User Home screen: greeting, current location, active order banner, promotions, articles, category grid, and featured workers.

---

### 3.1 `GET /home`

> Fetch all data for user home screen. Displayed as: "Selamat datang, Halo Budi!" with location, promo banners, category grid, and "Tukang Unggulan Terdekat".

🔒 **Auth Required** — Role: `user`

**Query Parameters:**

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `latitude` | number | ✅ | — | User's current latitude |
| `longitude` | number | ✅ | — | User's current longitude |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "user_summary": {
      "full_name": "Budi",
      "avatar_url": "https://cdn.handydirect.id/avatars/budi_001.jpg",
      "current_address": "Jl. Merdeka No. 12"
    },
    "active_order": {
      "order_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "status": "in_progress",
      "worker_name": "Ahmad Jaelani",
      "service_name": "Perbaikan Pipa",
      "eta_minutes": 8
    },
    "promotions": [
      {
        "promo_id": "promo-001",
        "title": "Diskon 20% Jasa AC",
        "description": "Berlaku hingga 31 Oktober 2023",
        "image_url": "https://cdn.handydirect.id/promos/diskon_ac.jpg",
        "cta_label": "Klaim Sekarang",
        "deep_link": "/promo/diskon-ac-20",
        "valid_until": "2023-10-31T23:59:59Z"
      }
    ],
    "articles": [
      {
        "article_id": "art-001",
        "title": "Tips Merawat Atap",
        "thumbnail_url": "https://cdn.handydirect.id/articles/tips_atap.jpg",
        "cta_label": "Baca →",
        "slug": "tips-merawat-atap"
      }
    ],
    "service_categories": [
      { "category_id": "cat-01", "name": "AC", "icon_url": "https://cdn.handydirect.id/icons/ac.png", "slug": "ac" },
      { "category_id": "cat-02", "name": "Pipa", "icon_url": "https://cdn.handydirect.id/icons/pipa.png", "slug": "pipa" },
      { "category_id": "cat-03", "name": "Atap", "icon_url": "https://cdn.handydirect.id/icons/atap.png", "slug": "atap" },
      { "category_id": "cat-04", "name": "Listrik", "icon_url": "https://cdn.handydirect.id/icons/listrik.png", "slug": "listrik" },
      { "category_id": "cat-05", "name": "Kunci", "icon_url": "https://cdn.handydirect.id/icons/kunci.png", "slug": "kunci" },
      { "category_id": "cat-06", "name": "Kayu", "icon_url": "https://cdn.handydirect.id/icons/kayu.png", "slug": "kayu" },
      { "category_id": "cat-07", "name": "Cat", "icon_url": "https://cdn.handydirect.id/icons/cat.png", "slug": "cat" },
      { "category_id": "cat-08", "name": "Kebun", "icon_url": "https://cdn.handydirect.id/icons/kebun.png", "slug": "kebun" }
    ],
    "featured_workers": [
      {
        "worker_id": "w-001",
        "full_name": "Ahmad Jaelani",
        "specialization": "Spesialis AC & Listrik",
        "avatar_url": "https://cdn.handydirect.id/workers/ahmad.jpg",
        "rating": 4.9,
        "distance_km": 1.2,
        "completed_jobs": 150,
        "is_verified": true
      },
      {
        "worker_id": "w-002",
        "full_name": "Budi Santoso",
        "specialization": "Pipa & Sanitasi",
        "avatar_url": "https://cdn.handydirect.id/workers/budis.jpg",
        "rating": 4.8,
        "distance_km": 2.5,
        "completed_jobs": 89,
        "is_verified": true
      }
    ]
  }
}
```

---

## 4. Service Categories

> Category grid and service listings. Maps to "Kategori Jasa" section with icons (AC, Pipa, Atap, Listrik, Kunci, Kayu, Cat, Kebun) and "Lihat Semua".

---

### 4.1 `GET /categories`

> Get all service categories.

**Query Parameters:**

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `page` | integer | ❌ | `1` | Page number |
| `per_page` | integer | ❌ | `20` | Items per page |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "category_id": "cat-01",
      "name": "AC",
      "slug": "ac",
      "icon_url": "https://cdn.handydirect.id/icons/ac.png",
      "description": "Servis, pemasangan, dan perbaikan AC",
      "display_order": 1,
      "is_active": true
    },
    {
      "category_id": "cat-02",
      "name": "Pipa",
      "slug": "pipa",
      "icon_url": "https://cdn.handydirect.id/icons/pipa.png",
      "description": "Perbaikan dan pemasangan pipa air",
      "display_order": 2,
      "is_active": true
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total": 8,
    "total_pages": 1
  }
}
```

---

### 4.2 `GET /categories/{category_id}/services`

> Get services within a specific category.

**Path Parameters:**

| Param | Type | Description |
|---|---|---|
| `category_id` | uuid | Category ID |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "service_id": "svc-001",
      "category_id": "cat-04",
      "name": "Instalasi Listrik",
      "slug": "instalasi-listrik",
      "description": "Pemasangan dan perbaikan instalasi listrik",
      "icon_url": "https://cdn.handydirect.id/icons/instalasi_listrik.png",
      "base_price": 150000,
      "price_unit": "per kunjungan",
      "estimated_duration": "1-3 jam",
      "is_active": true
    },
    {
      "service_id": "svc-002",
      "category_id": "cat-04",
      "name": "Perbaikan Ledeng",
      "slug": "perbaikan-ledeng",
      "description": "Perbaikan pipa bocor dan ledeng",
      "icon_url": "https://cdn.handydirect.id/icons/perbaikan_ledeng.png",
      "base_price": 120000,
      "price_unit": "per kunjungan",
      "estimated_duration": "1-2 jam",
      "is_active": true
    }
  ]
}
```

---

## 5. Nearby Workers (Tukang Dekat)

> Search and filter nearby workers by location, category, rating, etc. Maps to "Tukang Unggulan Terdekat" and "Lihat Semua" page.

---

### 5.1 `GET /workers/nearby`

> Get nearby workers sorted by distance, rating, or price.

🔒 **Auth Required** — Role: `user`

**Query Parameters:**

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `latitude` | number | ✅ | — | User's latitude |
| `longitude` | number | ✅ | — | User's longitude |
| `radius_km` | number | ❌ | `10` | Search radius in km |
| `category_id` | uuid | ❌ | — | Filter by service category |
| `service_id` | uuid | ❌ | — | Filter by specific service |
| `min_rating` | number | ❌ | — | Minimum rating (1.0–5.0) |
| `sort_by` | string | ❌ | `distance` | `distance` \| `rating` \| `price` \| `completed_jobs` |
| `sort_order` | string | ❌ | `asc` | `asc` \| `desc` |
| `page` | integer | ❌ | `1` | Page number |
| `per_page` | integer | ❌ | `10` | Items per page |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "worker_id": "w-001",
      "full_name": "Ahmad Jaelani",
      "avatar_url": "https://cdn.handydirect.id/workers/ahmad.jpg",
      "specialization": "Spesialis AC & Listrik",
      "rating": 4.9,
      "total_reviews": 120,
      "completed_jobs": 150,
      "distance_km": 1.2,
      "is_verified": true,
      "is_available": true,
      "base_price": 150000,
      "price_unit": "per kunjungan",
      "services": ["Instalasi Listrik", "Perbaikan Ledeng", "Servis AC"],
      "latitude": -7.260000,
      "longitude": 112.750000
    },
    {
      "worker_id": "w-002",
      "full_name": "Budi Santoso",
      "avatar_url": "https://cdn.handydirect.id/workers/budis.jpg",
      "specialization": "Pipa & Sanitasi",
      "rating": 4.8,
      "total_reviews": 89,
      "completed_jobs": 89,
      "distance_km": 2.5,
      "is_verified": true,
      "is_available": true,
      "base_price": 120000,
      "price_unit": "per kunjungan",
      "services": ["Perbaikan Pipa", "Pipe Repair & Cleaning"],
      "latitude": -7.262000,
      "longitude": 112.748000
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

---

### 5.2 `GET /workers/search`

> Search workers by keyword. Used by search bar: "Cari tukang atau jenis kerusakan..."

🔒 **Auth Required** — Role: `user`

**Query Parameters:**

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `q` | string | ✅ | — | Search keyword (name, specialization, damage type) |
| `latitude` | number | ✅ | — | User's latitude |
| `longitude` | number | ✅ | — | User's longitude |
| `radius_km` | number | ❌ | `10` | Search radius |
| `page` | integer | ❌ | `1` | — |
| `per_page` | integer | ❌ | `10` | — |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "worker_id": "w-001",
      "full_name": "Ahmad Jaelani",
      "specialization": "Spesialis AC & Listrik",
      "avatar_url": "https://cdn.handydirect.id/workers/ahmad.jpg",
      "rating": 4.9,
      "total_reviews": 120,
      "completed_jobs": 150,
      "distance_km": 1.2,
      "is_verified": true,
      "is_available": true,
      "base_price": 150000,
      "price_unit": "per kunjungan",
      "services": ["Instalasi Listrik", "Perbaikan Ledeng", "Servis AC"]
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 5,
    "total_pages": 1
  }
}
```

---

## 6. Worker Detail (Detail Tukang)

> Detail view of a worker. Maps to the mockup screen showing Ahmad Jaelani's profile with "Terverifikasi" badge, rating 4.9, 120 Ulasan, 150+ Pesanan Selesai, services (Instalasi Listrik, Perbaikan Ledeng, Servis AC), review quotes, and Booking Fee Rp2.000.

---

### 6.1 `GET /workers/{worker_id}`

> Get full worker profile detail.

🔒 **Auth Required** — Role: `user`

**Path Parameters:**

| Param | Type | Description |
|---|---|---|
| `worker_id` | uuid | Worker's user ID |

**Query Parameters:**

| Param | Type | Required | Description |
|---|---|---|---|
| `latitude` | number | ❌ | For distance calculation |
| `longitude` | number | ❌ | For distance calculation |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "worker_id": "w-001",
    "full_name": "Ahmad Jaelani",
    "avatar_url": "https://cdn.handydirect.id/workers/ahmad.jpg",
    "cover_photo_url": "https://cdn.handydirect.id/workers/ahmad_cover.jpg",
    "specialization": "Spesialis Kelistrikan & Ledeng",
    "bio": "Berpengalaman lebih dari 8 tahun dalam menangani instalasi listrik rumah tangga dan perbaikan pipa bocor dengan standar keamanan tinggi.",
    "is_verified": true,
    "verification_status": "verified",
    "rating": 4.9,
    "total_reviews": 120,
    "completed_jobs": 150,
    "distance_km": 1.2,
    "is_available": true,
    "member_since": "2021-03-10T00:00:00Z",
    "booking_fee": 2000,
    "services": [
      {
        "service_id": "svc-001",
        "name": "Instalasi Listrik",
        "icon_url": "https://cdn.handydirect.id/icons/listrik.png",
        "base_price": 150000,
        "price_unit": "per kunjungan"
      },
      {
        "service_id": "svc-002",
        "name": "Perbaikan Ledeng",
        "icon_url": "https://cdn.handydirect.id/icons/ledeng.png",
        "base_price": 120000,
        "price_unit": "per kunjungan"
      },
      {
        "service_id": "svc-003",
        "name": "Servis AC",
        "icon_url": "https://cdn.handydirect.id/icons/ac.png",
        "base_price": 175000,
        "price_unit": "per kunjungan"
      }
    ],
    "top_reviews": [
      {
        "review_id": "rev-001",
        "user_name": "Budi S.",
        "user_location": "Jakarta Selatan",
        "rating": 5,
        "comment": "Tukang Ahmad kerjanya cepat, rapi, dan sangat profesional. Pipa bocor langsung beres dalam sejam.",
        "created_at": "2023-10-15T08:30:00Z"
      }
    ]
  }
}
```

---

### 6.2 `GET /workers/{worker_id}/reviews`

> Get paginated reviews for a worker. Includes rating distribution breakdown.

🔒 **Auth Required** — Role: `user`

**Path Parameters:**

| Param | Type | Description |
|---|---|---|
| `worker_id` | uuid | Worker's user ID |

**Query Parameters:**

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `rating` | integer | ❌ | — | Filter by rating star (1–5) |
| `sort_by` | string | ❌ | `created_at` | `created_at` \| `rating` |
| `sort_order` | string | ❌ | `desc` | `asc` \| `desc` |
| `page` | integer | ❌ | `1` | — |
| `per_page` | integer | ❌ | `10` | — |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "review_id": "rev-001",
      "order_id": "ord-001",
      "user_name": "Budi S.",
      "user_avatar_url": "https://cdn.handydirect.id/avatars/budi_001.jpg",
      "user_location": "Jakarta Selatan",
      "rating": 5,
      "comment": "Tukang Ahmad kerjanya cepat, rapi, dan sangat profesional. Pipa bocor langsung beres dalam sejam.",
      "tags": ["cepat", "rapi", "profesional"],
      "created_at": "2023-10-15T08:30:00Z"
    }
  ],
  "meta": {
    "average_rating": 4.9,
    "rating_distribution": {
      "5": 95,
      "4": 18,
      "3": 5,
      "2": 1,
      "1": 1
    },
    "current_page": 1,
    "per_page": 10,
    "total": 120,
    "total_pages": 12
  }
}
```

---

### 6.3 `GET /workers/{worker_id}/services`

> Get list of services offered by a worker.

🔒 **Auth Required** — Role: `user`

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "service_id": "svc-001",
      "name": "Instalasi Listrik",
      "category": "Listrik",
      "description": "Pemasangan dan perbaikan instalasi listrik rumah tangga",
      "base_price": 150000,
      "price_unit": "per kunjungan",
      "estimated_duration": "1-3 jam",
      "icon_url": "https://cdn.handydirect.id/icons/listrik.png"
    },
    {
      "service_id": "svc-002",
      "name": "Perbaikan Ledeng",
      "category": "Pipa",
      "description": "Perbaikan pipa bocor, saluran air, dan ledeng",
      "base_price": 120000,
      "price_unit": "per kunjungan",
      "estimated_duration": "1-2 jam",
      "icon_url": "https://cdn.handydirect.id/icons/ledeng.png"
    },
    {
      "service_id": "svc-003",
      "name": "Servis AC",
      "category": "AC",
      "description": "Pembersihan, isi freon, dan perbaikan AC",
      "base_price": 175000,
      "price_unit": "per kunjungan",
      "estimated_duration": "1-2 jam",
      "icon_url": "https://cdn.handydirect.id/icons/ac.png"
    }
  ]
}
```



---

## 7. Orders — User Side

> Order/report creation and management from the user perspective. User selects a worker, creates an order, and can view or cancel orders.

---

### 7.1 `POST /orders`

> Create a new order/report. User selects a worker and service, describes the problem, and provides location.

🔒 **Auth Required** — Role: `user`

**Request Body:**

```json
{
  "worker_id": "w-001",
  "service_id": "svc-002",
  "title": "Perbaikan Pipa Bocor",
  "description": "Pipa wastafel bocor di bagian sambungan bawah. Air menetes cukup deras saat keran dinyalakan.",
  "location": {
    "latitude": -7.257500,
    "longitude": 112.752100,
    "address": "Jl. Merdeka No. 12, Surabaya",
    "address_detail": "Lantai 1, dekat dapur, rumah pagar hitam cat putih"
  },
  "preferred_date": "2023-10-25",
  "preferred_time_start": "09:00",
  "preferred_time_end": "12:00",
  "urgency": "normal",
  "photos": ["base64_encoded_photo_string_1"],
  "notes": "Tolong bawa seal tape dan sambungan pipa cadangan"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `worker_id` | uuid | ✅ | Selected worker's ID |
| `service_id` | uuid | ✅ | Selected service ID |
| `title` | string | ✅ | Order title (max 255 chars) |
| `description` | string | ✅ | Problem description |
| `location` | object | ✅ | Work location details |
| `location.latitude` | number | ✅ | Latitude |
| `location.longitude` | number | ✅ | Longitude |
| `location.address` | string | ✅ | Full address |
| `location.address_detail` | string | ❌ | Additional location details |
| `preferred_date` | date | ❌ | Preferred date (YYYY-MM-DD) |
| `preferred_time_start` | time | ❌ | Preferred start time (HH:MM) |
| `preferred_time_end` | time | ❌ | Preferred end time (HH:MM) |
| `urgency` | enum | ❌ | `normal` \| `urgent` (default: `normal`) |
| `photos` | array | ❌ | Problem photos (max 5, base64 or presigned URL) |
| `notes` | string | ❌ | Additional notes for worker |

**Response `201 Created`:**

```json
{
  "status": "success",
  "message": "Pesanan berhasil dibuat",
  "data": {
    "order_id": "ord-001",
    "order_number": "HD-20231025-001",
    "status": "pending",
    "worker_id": "w-001",
    "worker_name": "Ahmad Jaelani",
    "service_name": "Perbaikan Ledeng",
    "booking_fee": 2000,
    "estimated_base_price": 120000,
    "created_at": "2023-10-25T08:00:00Z"
  }
}
```

---

### 7.2 `GET /orders`

> Get user's order list. Used in the "Orders" tab in bottom navigation.

🔒 **Auth Required** — Role: `user`

**Query Parameters:**

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `status` | string | ❌ | — | Filter: `pending` \| `accepted` \| `in_progress` \| `completed` \| `cancelled` \| `rejected` |
| `sort_by` | string | ❌ | `created_at` | `created_at` \| `status` |
| `sort_order` | string | ❌ | `desc` | `asc` \| `desc` |
| `page` | integer | ❌ | `1` | — |
| `per_page` | integer | ❌ | `10` | — |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "order_id": "ord-005",
      "order_number": "HD-20231024-005",
      "title": "Pipe Repair & Cleaning",
      "status": "completed",
      "worker": {
        "worker_id": "w-002",
        "full_name": "Budi Santoso",
        "avatar_url": "https://cdn.handydirect.id/workers/budis.jpg"
      },
      "service_name": "Pipe Repair & Cleaning",
      "total_price": 350000,
      "created_at": "2023-10-24T14:00:00Z",
      "completed_at": "2023-10-24T15:30:00Z"
    },
    {
      "order_id": "ord-004",
      "order_number": "HD-20231023-004",
      "title": "AC Maintenance",
      "status": "completed",
      "worker": {
        "worker_id": "w-003",
        "full_name": "Siti Aminah",
        "avatar_url": "https://cdn.handydirect.id/workers/siti.jpg"
      },
      "service_name": "AC Maintenance",
      "total_price": 175000,
      "created_at": "2023-10-23T09:00:00Z",
      "completed_at": "2023-10-23T10:15:00Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 12,
    "total_pages": 2
  }
}
```

---

### 7.3 `GET /orders/{order_id}`

> Get full detail of a specific order.

🔒 **Auth Required** — Role: `user`

**Path Parameters:**

| Param | Type | Description |
|---|---|---|
| `order_id` | uuid | Order ID |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "order_id": "ord-001",
    "order_number": "HD-20231025-001",
    "title": "Perbaikan Pipa Bocor",
    "description": "Pipa wastafel bocor di bagian sambungan bawah. Air menetes cukup deras saat keran dinyalakan.",
    "status": "in_progress",
    "urgency": "normal",
    "worker": {
      "worker_id": "w-001",
      "full_name": "Ahmad Jaelani",
      "avatar_url": "https://cdn.handydirect.id/workers/ahmad.jpg",
      "specialization": "Spesialis Listrik & Air",
      "phone": "+62812xxxxxxx",
      "rating": 4.9,
      "total_reviews": 120,
      "is_verified": true
    },
    "service": {
      "service_id": "svc-002",
      "name": "Perbaikan Ledeng",
      "category": "Pipa"
    },
    "location": {
      "latitude": -7.257500,
      "longitude": 112.752100,
      "address": "Jl. Merdeka No. 12, Surabaya",
      "address_detail": "Lantai 1, dekat dapur"
    },
    "schedule": {
      "preferred_date": "2023-10-25",
      "preferred_time_start": "09:00",
      "preferred_time_end": "12:00"
    },
    "pricing": {
      "booking_fee": 2000,
      "base_service_fee": 150000,
      "total_material_cost": 75000,
      "total_additional_cost": 0,
      "grand_total": 225000
    },
    "photos": [
      "https://cdn.handydirect.id/orders/ord-001/photo_1.jpg"
    ],
    "notes": "Tolong bawa seal tape dan sambungan pipa cadangan",
    "timeline": [
      {
        "event": "order_created",
        "label": "Pesanan Dibuat",
        "timestamp": "2023-10-25T08:00:00Z",
        "is_completed": true
      },
      {
        "event": "order_accepted",
        "label": "Pesanan Diterima Worker",
        "timestamp": "2023-10-25T08:05:00Z",
        "is_completed": true
      },
      {
        "event": "worker_on_the_way",
        "label": "Tukang Menuju Lokasi",
        "timestamp": "2023-10-25T09:00:00Z",
        "is_completed": true
      },
      {
        "event": "worker_arrived",
        "label": "Tukang Tiba di Lokasi",
        "timestamp": "2023-10-25T09:15:00Z",
        "is_completed": true
      },
      {
        "event": "work_in_progress",
        "label": "Pengerjaan Berlangsung",
        "timestamp": "2023-10-25T09:20:00Z",
        "is_completed": true
      }
    ],
    "purchase_summary": {
      "total_items": 3,
      "total_cost": 75000,
      "pending_approval": 1,
      "approved": 2,
      "rejected": 0
    },
    "has_unread_chat": true,
    "can_cancel": true,
    "cancellation_policy": "Pembatalan gratis sebelum worker tiba di lokasi",
    "created_at": "2023-10-25T08:00:00Z",
    "updated_at": "2023-10-25T09:20:00Z"
  }
}
```

---

### 7.4 `POST /orders/{order_id}/cancel`

> Cancel an order (subject to cancellation policy).

🔒 **Auth Required** — Role: `user`

**Request Body:**

```json
{
  "reason": "Sudah menemukan tukang lain",
  "reason_category": "found_other_worker"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `reason` | string | ✅ | Cancellation reason text |
| `reason_category` | enum | ❌ | `changed_mind` \| `found_other_worker` \| `too_long` \| `other` |

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Pesanan berhasil dibatalkan",
  "data": {
    "order_id": "ord-001",
    "status": "cancelled",
    "cancellation_fee": 0,
    "refund_amount": 2000,
    "cancelled_at": "2023-10-25T08:30:00Z"
  }
}
```

**Error `422 Unprocessable Entity`:**

```json
{
  "status": "error",
  "message": "Pesanan tidak dapat dibatalkan karena worker sudah tiba di lokasi",
  "error_code": "CANCEL_NOT_ALLOWED"
}
```

---

## 8. Tracking — User Side

> Real-time tracking of worker location, job progress, and status. Maps to the tracking map mockup showing worker location on map with ETA, worker info card, Chat and Telepon buttons, and "Batalkan Pesanan".

---

### 8.1 `GET /orders/{order_id}/tracking`

> Get full tracking state: worker location, ETA, activity status, timeline, purchase summary.

🔒 **Auth Required** — Role: `user`

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "order_id": "ord-001",
    "order_status": "on_the_way",
    "worker": {
      "worker_id": "w-001",
      "full_name": "Ahmad Jaelani",
      "avatar_url": "https://cdn.handydirect.id/workers/ahmad.jpg",
      "specialization": "Spesialis Listrik & Air",
      "rating": 4.9,
      "total_reviews": 120
    },
    "worker_location": {
      "latitude": -7.259000,
      "longitude": 112.751000,
      "heading": 45.0,
      "updated_at": "2023-10-25T09:12:00Z"
    },
    "user_location": {
      "latitude": -7.257500,
      "longitude": 112.752100,
      "address": "Jl. Merdeka No. 12"
    },
    "eta_minutes": 8,
    "status_label": "Tukang menuju lokasi Anda",
    "activity_status": "on_the_way",
    "timeline": [
      {
        "event": "order_accepted",
        "label": "Pesanan Diterima",
        "timestamp": "2023-10-25T08:05:00Z",
        "is_completed": true
      },
      {
        "event": "worker_on_the_way",
        "label": "Tukang Menuju Lokasi",
        "timestamp": "2023-10-25T09:00:00Z",
        "is_completed": true
      },
      {
        "event": "worker_arrived",
        "label": "Tukang Tiba",
        "timestamp": null,
        "is_completed": false
      },
      {
        "event": "work_in_progress",
        "label": "Pengerjaan",
        "timestamp": null,
        "is_completed": false
      },
      {
        "event": "completed",
        "label": "Selesai",
        "timestamp": null,
        "is_completed": false
      }
    ],
    "purchase_summary": {
      "total_items": 0,
      "total_cost": 0,
      "pending_approval": 0
    },
    "can_cancel": true,
    "can_chat": true,
    "can_call": true
  }
}
```

---

### 8.2 `GET /orders/{order_id}/tracking/location`

> Poll worker's current location. Alternative to WebSocket for simple polling.

🔒 **Auth Required** — Role: `user`

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "latitude": -7.259000,
    "longitude": 112.751000,
    "heading": 45.0,
    "speed_kmh": 25,
    "eta_minutes": 8,
    "activity_status": "on_the_way",
    "updated_at": "2023-10-25T09:12:00Z"
  }
}
```

---

### 8.3 WebSocket — Real-time Tracking (User)

> **Endpoint:** `wss://api.handydirect.id/v1/ws/tracking/{order_id}`

**Authentication:** JWT token as query param `?token=xxx` or in initial connection header.

**Server → Client Events:**

**Location Update:**

```json
{
  "event": "location_update",
  "data": {
    "latitude": -7.258500,
    "longitude": 112.751500,
    "heading": 45.0,
    "eta_minutes": 5,
    "updated_at": "2023-10-25T09:13:00Z"
  }
}
```

**Status Change:**

```json
{
  "event": "status_change",
  "data": {
    "previous_status": "on_the_way",
    "new_status": "arrived",
    "label": "Tukang telah tiba di lokasi Anda",
    "timestamp": "2023-10-25T09:15:00Z"
  }
}
```

**New Purchase Added by Worker:**

```json
{
  "event": "new_purchase",
  "data": {
    "purchase_id": "pur-001",
    "item_name": "Pipa PVC 1/2 Inch",
    "total_price": 50000,
    "status": "pending_approval",
    "needs_user_approval": true
  }
}
```



---

## 9. Purchase Tracking — User Side

> User views, approves, rejects, or requests clarification on material/tool purchases made by the worker during the job. Core of the AI-assisted purchase tracking feature.

---

### 9.1 `GET /orders/{order_id}/purchases`

> Get all purchase items for an order, with summary and AI-generated summary text.

🔒 **Auth Required** — Role: `user`

**Query Parameters:**

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `status` | string | ❌ | — | Filter: `draft` \| `pending_approval` \| `approved` \| `rejected` \| `needs_clarification` |
| `page` | integer | ❌ | `1` | — |
| `per_page` | integer | ❌ | `20` | — |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "order_id": "ord-001",
    "purchases": [
      {
        "purchase_id": "pur-001",
        "item_name": "Pipa PVC 1/2 Inch",
        "category": "material",
        "quantity": 2,
        "unit": "meter",
        "unit_price": 25000,
        "total_price": 50000,
        "reason": "Untuk mengganti bagian pipa yang rusak di bawah wastafel",
        "receipt_photo_url": "https://cdn.handydirect.id/receipts/pur-001.jpg",
        "status": "approved",
        "confidence": 0.95,
        "needs_clarification": false,
        "clarification_question": null,
        "risk_flags": [],
        "approved_at": "2023-10-25T10:00:00Z",
        "created_at": "2023-10-25T09:45:00Z"
      },
      {
        "purchase_id": "pur-002",
        "item_name": "Lem Pipa",
        "category": "material",
        "quantity": 1,
        "unit": "kaleng",
        "unit_price": 15000,
        "total_price": 15000,
        "reason": "Untuk merekatkan sambungan pipa baru",
        "receipt_photo_url": null,
        "status": "pending_approval",
        "confidence": 0.92,
        "needs_clarification": false,
        "clarification_question": null,
        "risk_flags": [],
        "approved_at": null,
        "created_at": "2023-10-25T10:15:00Z"
      },
      {
        "purchase_id": "pur-003",
        "item_name": "Alat tidak diketahui",
        "category": "alat",
        "quantity": 1,
        "unit": "pcs",
        "unit_price": 200000,
        "total_price": 200000,
        "reason": "",
        "receipt_photo_url": null,
        "status": "needs_clarification",
        "confidence": 0.30,
        "needs_clarification": true,
        "clarification_question": "Nama alat belum jelas. Mohon worker mengisi nama item, jumlah, dan alasan pembelian.",
        "risk_flags": [
          {
            "type": "data_tidak_lengkap",
            "message": "Nama item dan alasan pembelian belum diisi"
          },
          {
            "type": "harga_tidak_wajar",
            "message": "Harga Rp200.000 terlihat tinggi untuk item tanpa deskripsi"
          }
        ],
        "approved_at": null,
        "created_at": "2023-10-25T10:30:00Z"
      }
    ],
    "summary": {
      "total_items": 3,
      "total_cost": 265000,
      "approved_cost": 50000,
      "pending_cost": 15000,
      "rejected_cost": 0,
      "needs_clarification_cost": 200000,
      "ai_summary": "Total pembelian tambahan saat ini adalah Rp265.000, terdiri dari Pipa PVC 1/2 Inch (disetujui), Lem Pipa (menunggu persetujuan), dan 1 item yang memerlukan klarifikasi."
    }
  },
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total": 3,
    "total_pages": 1
  }
}
```

---

### 9.2 `GET /orders/{order_id}/purchases/{purchase_id}`

> Get detail of a single purchase item including audit log and AI explanation.

🔒 **Auth Required** — Role: `user`

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "purchase_id": "pur-001",
    "order_id": "ord-001",
    "item_name": "Pipa PVC 1/2 Inch",
    "category": "material",
    "quantity": 2,
    "unit": "meter",
    "unit_price": 25000,
    "total_price": 50000,
    "reason": "Untuk mengganti bagian pipa yang rusak di bawah wastafel",
    "receipt_photo_url": "https://cdn.handydirect.id/receipts/pur-001.jpg",
    "status": "approved",
    "confidence": 0.95,
    "needs_clarification": false,
    "clarification_question": null,
    "clarification_response": null,
    "ai_explanation": "Worker menambahkan pembelian Pipa PVC 1/2 Inch sebanyak 2 meter seharga Rp25.000/meter karena diperlukan untuk mengganti bagian pipa yang rusak di bawah wastafel.",
    "risk_flags": [],
    "audit_log": [
      {
        "action": "created",
        "actor_type": "worker",
        "actor_name": "Ahmad Jaelani",
        "timestamp": "2023-10-25T09:45:00Z",
        "note": "Pembelian dimasukkan oleh worker"
      },
      {
        "action": "ai_processed",
        "actor_type": "system",
        "actor_name": "AI Assistant",
        "timestamp": "2023-10-25T09:45:05Z",
        "note": "Input dirapikan oleh AI, confidence 95%"
      },
      {
        "action": "submitted",
        "actor_type": "worker",
        "actor_name": "Ahmad Jaelani",
        "timestamp": "2023-10-25T09:45:10Z",
        "note": "Dikirim untuk persetujuan user"
      },
      {
        "action": "approved",
        "actor_type": "user",
        "actor_name": "Budi Santoso",
        "timestamp": "2023-10-25T10:00:00Z",
        "note": null
      }
    ],
    "created_at": "2023-10-25T09:45:00Z",
    "updated_at": "2023-10-25T10:00:00Z"
  }
}
```

---

### 9.3 `PATCH /orders/{order_id}/purchases/{purchase_id}/approve`

> Approve a purchase item. Once approved, it is included in the final invoice.

🔒 **Auth Required** — Role: `user`

**Request Body (optional):**

```json
{
  "note": "OK, setuju"
}
```

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Pembelian berhasil disetujui",
  "data": {
    "purchase_id": "pur-002",
    "status": "approved",
    "approved_at": "2023-10-25T10:20:00Z"
  }
}
```

---

### 9.4 `PATCH /orders/{order_id}/purchases/{purchase_id}/reject`

> Reject a purchase item.

🔒 **Auth Required** — Role: `user`

**Request Body:**

```json
{
  "reason": "Item tidak relevan dengan pekerjaan yang diminta"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `reason` | string | ✅ | Rejection reason |

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Pembelian ditolak",
  "data": {
    "purchase_id": "pur-003",
    "status": "rejected",
    "rejected_at": "2023-10-25T10:25:00Z"
  }
}
```

---

### 9.5 `PATCH /orders/{order_id}/purchases/{purchase_id}/clarify`

> Request clarification on a purchase. Worker will receive a notification.

🔒 **Auth Required** — Role: `user`

**Request Body:**

```json
{
  "question": "Untuk apa alat ini digunakan? Mohon jelaskan nama alat dan kebutuhan spesifiknya."
}
```

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Permintaan klarifikasi berhasil dikirim ke worker",
  "data": {
    "purchase_id": "pur-003",
    "status": "needs_clarification",
    "clarification_question": "Untuk apa alat ini digunakan? Mohon jelaskan nama alat dan kebutuhan spesifiknya."
  }
}
```

---

### 9.6 `PATCH /orders/{order_id}/purchases/bulk-approve`

> Approve multiple purchase items at once.

🔒 **Auth Required** — Role: `user`

**Request Body:**

```json
{
  "purchase_ids": ["pur-001", "pur-002", "pur-004"],
  "note": "Semua sudah dicek, setuju"
}
```

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "3 pembelian berhasil disetujui",
  "data": {
    "approved_count": 3,
    "approved_ids": ["pur-001", "pur-002", "pur-004"]
  }
}
```

---

## 10. Chat — User Side

> Real-time chat between user and worker per order. Maps to the chat mockup showing conversation between Budi and Ahmad Jaelani with text messages, image of broken pipe, and status "Sedang menuju lokasi".

---

### 10.1 `GET /orders/{order_id}/chat/messages`

> Get chat message history (cursor-based pagination, newest first).

🔒 **Auth Required** — Role: `user`

**Query Parameters:**

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `before` | string (ISO 8601) | ❌ | — | Cursor: fetch messages before this timestamp |
| `limit` | integer | ❌ | `50` | Messages per request (max 100) |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "order_id": "ord-001",
    "worker": {
      "worker_id": "w-001",
      "full_name": "Ahmad Jaelani",
      "avatar_url": "https://cdn.handydirect.id/workers/ahmad.jpg",
      "status_label": "Sedang menuju lokasi"
    },
    "messages": [
      {
        "message_id": "msg-001",
        "sender_type": "user",
        "sender_name": "Budi",
        "content": "Halo Pak Ahmad, apakah sudah di jalan?",
        "message_type": "text",
        "media_url": null,
        "is_read": true,
        "read_at": "2023-10-25T10:15:30Z",
        "created_at": "2023-10-25T10:15:00Z"
      },
      {
        "message_id": "msg-002",
        "sender_type": "worker",
        "sender_name": "Ahmad Jaelani",
        "content": "Halo Pak, iya saya sudah di jalan. Mungkin sekitar 15 menit lagi sampai. Boleh dikirimkan foto pipa yang bocornya lagi supaya saya siapkan alat yang tepat?",
        "message_type": "text",
        "media_url": null,
        "is_read": true,
        "read_at": "2023-10-25T10:17:15Z",
        "created_at": "2023-10-25T10:17:00Z"
      },
      {
        "message_id": "msg-003",
        "sender_type": "user",
        "sender_name": "Budi",
        "content": "Ini pipa yang di bawah wastafel dapur Pak. Bocornya lumayan deras kalau keran dinyalakan.",
        "message_type": "image",
        "media_url": "https://cdn.handydirect.id/chat/img_pipe_leak.jpg",
        "is_read": true,
        "read_at": "2023-10-25T10:20:30Z",
        "created_at": "2023-10-25T10:20:00Z"
      },
      {
        "message_id": "msg-004",
        "sender_type": "worker",
        "sender_name": "Ahmad Jaelani",
        "content": "Baik Pak, saya mengerti. Saya bawa seal tape dan sambungan pipa cadangan. Tunggu sebentar ya.",
        "message_type": "text",
        "media_url": null,
        "is_read": true,
        "read_at": "2023-10-25T10:22:15Z",
        "created_at": "2023-10-25T10:22:00Z"
      }
    ],
    "has_more": false
  }
}
```

---

### 10.2 `POST /orders/{order_id}/chat/messages`

> Send a new chat message.

🔒 **Auth Required** — Role: `user`

**Request Body (text):**

```json
{
  "content": "Halo Pak Ahmad, apakah sudah di jalan?",
  "message_type": "text"
}
```

**Request Body (image — `multipart/form-data`):**

| Field | Type | Required | Description |
|---|---|---|---|
| `message_type` | string | ✅ | `image` |
| `content` | string | ❌ | Image caption |
| `media` | file | ✅ | Image file (jpg/png, max 10MB) |

**Response `201 Created`:**

```json
{
  "status": "success",
  "data": {
    "message_id": "msg-005",
    "sender_type": "user",
    "sender_name": "Budi",
    "content": "Halo Pak Ahmad, apakah sudah di jalan?",
    "message_type": "text",
    "media_url": null,
    "created_at": "2023-10-25T10:15:00Z"
  }
}
```

---

### 10.3 `PATCH /orders/{order_id}/chat/read`

> Mark all messages in this chat as read.

🔒 **Auth Required** — Role: `user`

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Semua pesan ditandai sudah dibaca"
}
```

---

### 10.4 `GET /chats`

> Get list of all active chat conversations. Used in "Chat" tab on bottom navigation.

🔒 **Auth Required** — Role: `user`

**Query Parameters:**

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `page` | integer | ❌ | `1` | — |
| `per_page` | integer | ❌ | `20` | — |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "order_id": "ord-001",
      "worker": {
        "worker_id": "w-001",
        "full_name": "Ahmad Jaelani",
        "avatar_url": "https://cdn.handydirect.id/workers/ahmad.jpg",
        "is_online": true
      },
      "order_title": "Perbaikan Pipa Bocor",
      "order_status": "in_progress",
      "last_message": {
        "content": "Baik Pak, saya mengerti. Saya bawa seal tape...",
        "sender_type": "worker",
        "message_type": "text",
        "created_at": "2023-10-25T10:22:00Z"
      },
      "unread_count": 0
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total": 3,
    "total_pages": 1
  }
}
```

---

### 10.5 WebSocket — Real-time Chat

> **Endpoint:** `wss://api.handydirect.id/v1/ws/chat/{order_id}`

**Server → Client Events:**

**New Message:**

```json
{
  "event": "new_message",
  "data": {
    "message_id": "msg-006",
    "sender_type": "worker",
    "sender_name": "Ahmad Jaelani",
    "content": "Saya sudah sampai di depan rumah Pak.",
    "message_type": "text",
    "media_url": null,
    "created_at": "2023-10-25T10:30:00Z"
  }
}
```

**Typing Indicator:**

```json
{
  "event": "typing",
  "data": {
    "sender_type": "worker",
    "is_typing": true
  }
}
```

**Message Read:**

```json
{
  "event": "message_read",
  "data": {
    "reader_type": "worker",
    "read_until": "2023-10-25T10:20:00Z"
  }
}
```

**Client → Server Events:**

```json
{
  "event": "send_message",
  "data": {
    "content": "Oke Pak, saya bukakan pintu",
    "message_type": "text"
  }
}
```

```json
{
  "event": "typing",
  "data": { "is_typing": true }
}
```

---

## 11. Rating & Review — User Side

> User rates the worker after job completion. Maps to the "Bagaimana layanan Budi?" mockup with star rating and comment box.

---

### 11.1 `POST /orders/{order_id}/rating`

> Submit rating and review after job is completed.

🔒 **Auth Required** — Role: `user`

**Request Body:**

```json
{
  "rating": 4,
  "comment": "Kerjanya bagus dan rapi, pipa sudah tidak bocor lagi. Terima kasih Pak Ahmad!",
  "tags": ["cepat", "rapi", "profesional"]
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `rating` | integer | ✅ | 1–5 stars |
| `comment` | string | ❌ | Review text (max 1000 chars) |
| `tags` | array | ❌ | Quick tags: `cepat`, `rapi`, `profesional`, `ramah`, etc. |

**Response `201 Created`:**

```json
{
  "status": "success",
  "message": "Rating berhasil dikirim. Terima kasih!",
  "data": {
    "review_id": "rev-010",
    "order_id": "ord-001",
    "worker_id": "w-001",
    "rating": 4,
    "comment": "Kerjanya bagus dan rapi, pipa sudah tidak bocor lagi. Terima kasih Pak Ahmad!",
    "tags": ["cepat", "rapi", "profesional"],
    "created_at": "2023-10-25T12:00:00Z"
  }
}
```

**Error `409 Conflict`:**

```json
{
  "status": "error",
  "message": "Anda sudah memberikan rating untuk pesanan ini",
  "error_code": "RATING_ALREADY_EXISTS"
}
```

---

### 11.2 `GET /orders/{order_id}/rating`

> Get the rating already submitted by user for this order.

🔒 **Auth Required** — Role: `user`

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "review_id": "rev-010",
    "order_id": "ord-001",
    "worker_id": "w-001",
    "rating": 4,
    "comment": "Kerjanya bagus dan rapi, pipa sudah tidak bocor lagi. Terima kasih Pak Ahmad!",
    "tags": ["cepat", "rapi", "profesional"],
    "created_at": "2023-10-25T12:00:00Z"
  }
}
```

**Error `404 Not Found`:**

```json
{
  "status": "error",
  "message": "Belum ada rating untuk pesanan ini",
  "error_code": "RATING_NOT_FOUND"
}
```



---

## 12. Invoice & Payment — User Side

> Invoice detail and payment confirmation. Maps to the "Kerja Selesai!" mockup showing "RINCIAN INVOICE" with Biaya Jasa Dasar, Biaya Material, Total Pembayaran Cash, and payment instruction.

---

### 12.1 `GET /orders/{order_id}/invoice`

> Get final invoice for a completed order.

🔒 **Auth Required** — Role: `user`

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "invoice_id": "inv-001",
    "order_id": "ord-001",
    "order_number": "HD-20231025-001",
    "invoice_number": "INV-20231025-001",
    "worker": {
      "worker_id": "w-001",
      "full_name": "Ahmad Jaelani",
      "avatar_url": "https://cdn.handydirect.id/workers/ahmad.jpg"
    },
    "service_name": "Perbaikan Ledeng",
    "line_items": {
      "base_service_fee": {
        "label": "Biaya Jasa Dasar",
        "amount": 150000
      },
      "material_cost": {
        "label": "Biaya Material",
        "amount": 75000,
        "breakdown": [
          {
            "item_name": "Pipa PVC 1/2 Inch",
            "quantity": 2,
            "unit": "meter",
            "unit_price": 25000,
            "total_price": 50000
          },
          {
            "item_name": "Lem Pipa",
            "quantity": 1,
            "unit": "kaleng",
            "unit_price": 15000,
            "total_price": 15000
          },
          {
            "item_name": "Seal Tape",
            "quantity": 1,
            "unit": "pcs",
            "unit_price": 10000,
            "total_price": 10000
          }
        ]
      },
      "additional_cost": {
        "label": "Biaya Tambahan",
        "amount": 0
      },
      "booking_fee": {
        "label": "Booking Fee",
        "amount": 2000
      },
      "platform_fee": {
        "label": "Biaya Layanan",
        "amount": 0
      },
      "discount": {
        "label": "Diskon",
        "amount": 0,
        "promo_code": null
      }
    },
    "grand_total": 225000,
    "currency": "IDR",
    "payment_method": "cash",
    "payment_status": "unpaid",
    "payment_instruction": "Silakan lakukan pembayaran secara tunai langsung kepada teknisi kami di lokasi.",
    "ai_report": {
      "work_summary": "Perbaikan pipa bocor di bawah wastafel dapur. Pipa PVC yang rusak diganti dengan yang baru, sambungan diperbaiki dengan lem pipa, dan area sekitar dibersihkan.",
      "materials_used": "Pipa PVC 1/2 Inch 2 meter, Lem Pipa 1 kaleng, Seal Tape 1 pcs",
      "total_service_fee": 150000,
      "total_material_fee": 75000,
      "worker_notes": "Pekerjaan selesai. Pipa sudah tidak bocor. Disarankan pengecekan berkala setiap 6 bulan.",
      "all_purchases_approved": true
    },
    "created_at": "2023-10-25T11:30:00Z"
  }
}
```

---

### 12.2 `POST /orders/{order_id}/payment`

> Confirm or initiate payment.

🔒 **Auth Required** — Role: `user`

**Request Body:**

```json
{
  "payment_method": "cash",
  "payment_proof_url": null,
  "promo_code": null
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `payment_method` | enum | ✅ | `cash` \| `bank_transfer` \| `ewallet` |
| `payment_proof_url` | string | ❌ | URL of payment proof (for bank_transfer) |
| `promo_code` | string | ❌ | Promo code to apply |

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Pembayaran berhasil dikonfirmasi",
  "data": {
    "payment_id": "pay-001",
    "order_id": "ord-001",
    "amount": 225000,
    "payment_method": "cash",
    "payment_status": "paid",
    "paid_at": "2023-10-25T11:45:00Z",
    "receipt_url": null
  }
}
```

---

### 12.3 `GET /orders/{order_id}/invoice/pdf`

> Download invoice as PDF file.

🔒 **Auth Required** — Role: `user`

**Response `200 OK`:**

```
Content-Type: application/pdf
Content-Disposition: attachment; filename="INV-20231025-001.pdf"
```

---

## 13. Knowledge / FAQ / Guide

> Information center, tutorials, and FAQ. Maps to "Pusat Bantuan" and article sections on Home.

---

### 13.1 `GET /knowledge/articles`

> Get list of knowledge base articles.

**Query Parameters:**

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `category` | string | ❌ | — | `faq` \| `guide` \| `tips` \| `safety` \| `payment` |
| `q` | string | ❌ | — | Search keyword |
| `page` | integer | ❌ | `1` | — |
| `per_page` | integer | ❌ | `10` | — |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "article_id": "art-001",
      "title": "Tips Merawat Atap Rumah",
      "slug": "tips-merawat-atap",
      "category": "tips",
      "thumbnail_url": "https://cdn.handydirect.id/articles/tips_atap.jpg",
      "excerpt": "Pelajari cara merawat atap rumah agar tetap awet dan tidak bocor...",
      "read_time_minutes": 5,
      "published_at": "2023-10-20T08:00:00Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 15,
    "total_pages": 2
  }
}
```

---

### 13.2 `GET /knowledge/articles/{article_id}`

> Get full article detail.

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "article_id": "art-001",
    "title": "Tips Merawat Atap Rumah",
    "slug": "tips-merawat-atap",
    "category": "tips",
    "thumbnail_url": "https://cdn.handydirect.id/articles/tips_atap.jpg",
    "content_html": "<h2>Perawatan Rutin Atap</h2><p>Atap rumah merupakan bagian penting yang melindungi seluruh bangunan...</p>",
    "read_time_minutes": 5,
    "author": "Tim HandyDirect",
    "tags": ["atap", "perawatan", "rumah"],
    "published_at": "2023-10-20T08:00:00Z"
  }
}
```

---

### 13.3 `GET /knowledge/faq`

> Get FAQ list.

**Query Parameters:**

| Param | Type | Required | Description |
|---|---|---|---|
| `category` | string | ❌ | `general` \| `payment` \| `tracking` \| `security` \| `cancellation` |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "faq_id": "faq-001",
      "question": "Bagaimana cara membatalkan pesanan?",
      "answer": "Anda dapat membatalkan pesanan sebelum worker tiba di lokasi melalui halaman tracking. Pembatalan setelah worker tiba dapat dikenakan biaya pembatalan.",
      "category": "cancellation",
      "display_order": 1
    },
    {
      "faq_id": "faq-002",
      "question": "Apakah biaya material otomatis masuk tagihan?",
      "answer": "Tidak. Setiap pembelian material oleh worker harus disetujui terlebih dahulu oleh Anda sebelum masuk ke tagihan akhir. Anda dapat menyetujui, menolak, atau meminta klarifikasi.",
      "category": "payment",
      "display_order": 2
    },
    {
      "faq_id": "faq-003",
      "question": "Bagaimana cara kerja AI purchase tracking?",
      "answer": "AI membantu merapikan input pembelian dari worker menjadi data terstruktur, memberikan ringkasan biaya, dan mendeteksi pembelian yang tidak wajar. Semua keputusan tetap di tangan Anda.",
      "category": "tracking",
      "display_order": 3
    }
  ]
}
```

---

## 14. Notifications

> Push notifications and in-app notifications for users and workers.

---

### 14.1 `GET /notifications`

> Get notification list for current user.

🔒 **Auth Required** — Role: `user`, `worker`

**Query Parameters:**

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `is_read` | boolean | ❌ | — | Filter read/unread |
| `type` | string | ❌ | — | `order` \| `purchase` \| `chat` \| `promo` \| `system` \| `payment` |
| `page` | integer | ❌ | `1` | — |
| `per_page` | integer | ❌ | `20` | — |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "notification_id": "notif-001",
      "type": "purchase",
      "title": "Pembelian Baru Perlu Persetujuan",
      "body": "Ahmad Jaelani menambahkan pembelian Lem Pipa (1 kaleng, Rp15.000). Setujui atau tolak pembelian ini.",
      "is_read": false,
      "deep_link": "/orders/ord-001/purchases/pur-002",
      "metadata": {
        "order_id": "ord-001",
        "purchase_id": "pur-002"
      },
      "created_at": "2023-10-25T10:15:00Z"
    },
    {
      "notification_id": "notif-002",
      "type": "order",
      "title": "Pesanan Diterima",
      "body": "Ahmad Jaelani telah menerima pesanan Anda. Estimasi tiba: 15 menit.",
      "is_read": true,
      "deep_link": "/orders/ord-001/tracking",
      "metadata": {
        "order_id": "ord-001"
      },
      "created_at": "2023-10-25T08:05:00Z"
    }
  ],
  "meta": {
    "unread_count": 3,
    "current_page": 1,
    "per_page": 20,
    "total": 15,
    "total_pages": 1
  }
}
```

---

### 14.2 `PATCH /notifications/{notification_id}/read`

> Mark a single notification as read.

🔒 **Auth Required** — Role: `user`, `worker`

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Notifikasi ditandai sudah dibaca"
}
```

---

### 14.3 `PATCH /notifications/read-all`

> Mark all notifications as read.

🔒 **Auth Required** — Role: `user`, `worker`

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Semua notifikasi ditandai sudah dibaca"
}
```



---

# Part C — Worker (Tukang) Side

> All endpoints in this section require `role: worker` in the JWT token. Middleware will return `403 Forbidden` if a user with role `user` attempts to access these endpoints.

---

## 15. Worker Profile & Verification

> Worker-specific profile management and identity verification. Maps to "Profil & Keahlian" menu on worker home screen.

---

### 15.1 `GET /worker/profile`

> Get worker's full profile including worker_profiles data.

🔒 **Auth Required** — Role: `worker`

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "user_id": "w-001",
    "full_name": "Ahmad Jaelani",
    "email": "ahmad@email.com",
    "phone": "+6281234567891",
    "avatar_url": "https://cdn.handydirect.id/workers/ahmad.jpg",
    "address": "Jl. Kebon Jeruk No. 5, Jakarta",
    "latitude": -6.200000,
    "longitude": 106.816666,
    "worker_profile": {
      "specialization": "Spesialis Kelistrikan & Ledeng",
      "bio": "Berpengalaman lebih dari 8 tahun dalam menangani instalasi listrik rumah tangga dan perbaikan pipa bocor dengan standar keamanan tinggi.",
      "cover_photo_url": "https://cdn.handydirect.id/workers/ahmad_cover.jpg",
      "verification_status": "verified",
      "id_card_url": "https://cdn.handydirect.id/verification/ahmad_ktp.jpg",
      "certificate_urls": [
        "https://cdn.handydirect.id/certs/ahmad_cert1.jpg",
        "https://cdn.handydirect.id/certs/ahmad_cert2.jpg"
      ],
      "base_price": 150000,
      "price_unit": "per kunjungan",
      "booking_fee": 2000,
      "rating_avg": 4.9,
      "total_reviews": 120,
      "completed_jobs": 150,
      "is_available": true,
      "verified_at": "2021-04-15T10:00:00Z",
      "member_since": "2021-03-10T00:00:00Z"
    },
    "services": [
      {
        "service_id": "svc-001",
        "name": "Instalasi Listrik",
        "category": "Listrik",
        "custom_price": 160000,
        "is_active": true
      },
      {
        "service_id": "svc-002",
        "name": "Perbaikan Ledeng",
        "category": "Pipa",
        "custom_price": null,
        "is_active": true
      },
      {
        "service_id": "svc-003",
        "name": "Servis AC",
        "category": "AC",
        "custom_price": 175000,
        "is_active": true
      }
    ]
  }
}
```

---

### 15.2 `PUT /worker/profile`

> Update worker profile details.

🔒 **Auth Required** — Role: `worker`

**Request Body:**

```json
{
  "full_name": "Ahmad Jaelani",
  "phone": "+6281234567891",
  "address": "Jl. Kebon Jeruk No. 5, Jakarta",
  "specialization": "Spesialis Kelistrikan & Ledeng",
  "bio": "Berpengalaman lebih dari 8 tahun dalam menangani instalasi listrik rumah tangga dan perbaikan pipa bocor.",
  "base_price": 150000,
  "price_unit": "per kunjungan",
  "services": [
    { "service_id": "svc-001", "custom_price": 160000, "is_active": true },
    { "service_id": "svc-002", "custom_price": null, "is_active": true },
    { "service_id": "svc-003", "custom_price": 175000, "is_active": true }
  ]
}
```

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Profil worker berhasil diperbarui",
  "data": {
    "user_id": "w-001",
    "full_name": "Ahmad Jaelani",
    "specialization": "Spesialis Kelistrikan & Ledeng",
    "updated_at": "2026-05-30T09:00:00Z"
  }
}
```

---

### 15.3 `PUT /worker/profile/cover-photo`

> Upload or update worker cover photo.

🔒 **Auth Required** — Role: `worker`
📎 **Content-Type:** `multipart/form-data`

| Field | Type | Required | Description |
|---|---|---|---|
| `cover_photo` | file | ✅ | Cover image (jpg/png, max 10MB) |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "cover_photo_url": "https://cdn.handydirect.id/workers/ahmad_cover_v2.jpg"
  }
}
```

---

### 15.4 `POST /worker/profile/verification`

> Submit identity verification documents (KTP, certificates).

🔒 **Auth Required** — Role: `worker`
📎 **Content-Type:** `multipart/form-data`

| Field | Type | Required | Description |
|---|---|---|---|
| `id_card` | file | ✅ | KTP photo (jpg/png, max 5MB) |
| `certificates` | file[] | ❌ | Skill certificates (max 5 files, max 5MB each) |
| `selfie_with_id` | file | ❌ | Selfie holding KTP |

**Response `201 Created`:**

```json
{
  "status": "success",
  "message": "Dokumen verifikasi berhasil dikirim. Proses verifikasi membutuhkan 1-3 hari kerja.",
  "data": {
    "verification_status": "pending",
    "submitted_at": "2026-05-30T09:00:00Z"
  }
}
```

---

### 15.5 `GET /worker/profile/verification`

> Get current verification status.

🔒 **Auth Required** — Role: `worker`

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "verification_status": "verified",
    "id_card_url": "https://cdn.handydirect.id/verification/ahmad_ktp.jpg",
    "certificate_urls": [
      "https://cdn.handydirect.id/certs/ahmad_cert1.jpg"
    ],
    "submitted_at": "2021-04-10T10:00:00Z",
    "verified_at": "2021-04-15T10:00:00Z",
    "rejection_reason": null
  }
}
```

---

## 16. Home (Worker)

> Worker home dashboard. Maps to the worker home mockup showing: "Bpk. Ahmad" with Terverifikasi badge, "Pendapatan Hari Ini: Rp150.000", availability toggle "Siap Menerima Order", quick menu grid (Riwayat, Penarikan, Profil & Keahlian, Pusat Bantuan), and "Ringkasan Minggu Ini" with Tingkat Penerimaan 95% and Rating Rata-rata 4.9.

---

### 16.1 `GET /worker/home`

> Aggregated dashboard data for worker home screen.

🔒 **Auth Required** — Role: `worker`

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "worker_summary": {
      "full_name": "Bpk. Ahmad",
      "avatar_url": "https://cdn.handydirect.id/workers/ahmad.jpg",
      "is_verified": true,
      "verification_status": "verified",
      "is_available": true
    },
    "earnings": {
      "today": 150000,
      "this_week": 1250000,
      "this_month": 8450000
    },
    "wallet_balance": 150000,
    "weekly_summary": {
      "acceptance_rate": 95,
      "average_rating": 4.9,
      "jobs_completed_this_week": 12,
      "jobs_completed_this_month": 42
    },
    "incoming_orders_count": 1,
    "active_order": {
      "order_id": "ord-012",
      "title": "Perbaikan Pipa Bocor",
      "status": "in_progress",
      "user_name": "Budi Santoso",
      "started_at": "2023-10-25T10:00:00Z"
    },
    "quick_menu": [
      {
        "id": "history",
        "label": "Riwayat",
        "icon": "clock",
        "subtitle": "12 Selesai bulan ini",
        "deep_link": "/worker/history"
      },
      {
        "id": "withdrawal",
        "label": "Penarikan",
        "icon": "wallet",
        "subtitle": "Tarik dana ke bank",
        "deep_link": "/worker/wallet"
      },
      {
        "id": "profile",
        "label": "Profil & Keahlian",
        "icon": "user-check",
        "subtitle": "Update sertifikat",
        "deep_link": "/worker/profile"
      },
      {
        "id": "help",
        "label": "Pusat Bantuan",
        "icon": "headset",
        "subtitle": "Hubungi CS",
        "deep_link": "/knowledge"
      }
    ]
  }
}
```

---

### 16.2 `PATCH /worker/availability`

> Toggle worker's availability status ("Siap Menerima Order" toggle). When off, worker won't appear in search results and won't receive new orders.

🔒 **Auth Required** — Role: `worker`

**Request Body:**

```json
{
  "is_available": true
}
```

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Status ketersediaan berhasil diperbarui",
  "data": {
    "is_available": true,
    "status_label": "Siap Menerima Order",
    "status_description": "Anda terlihat online oleh pelanggan di sekitar."
  }
}
```

---

## 17. Incoming Orders — Worker Side

> Incoming order management. Maps to the "Order Baru Masuk!" popup/sheet showing order title, category, description, distance, countdown timer, and "Slide untuk Terima Order".

---

### 17.1 `GET /worker/orders/incoming`

> Get list of pending orders assigned to this worker. Each order has a countdown timer for auto-rejection.

🔒 **Auth Required** — Role: `worker`

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "order_id": "ord-020",
      "order_number": "HD-20231025-020",
      "title": "Pipa Bocor di Dapur",
      "service_name": "Perbaikan Pipa & Saluran Air",
      "category_name": "Pipa",
      "category_icon_url": "https://cdn.handydirect.id/icons/pipa.png",
      "description_preview": "Pipa wastafel bocor di bagian sambungan bawah. Air menetes cukup deras saat keran dinyalakan.",
      "user": {
        "user_id": "u-005",
        "full_name": "Budi Santoso",
        "avatar_url": "https://cdn.handydirect.id/avatars/budi_001.jpg"
      },
      "location": {
        "address": "Jl. Merdeka Selatan No. 45",
        "distance_km": 1.5
      },
      "urgency": "normal",
      "estimated_base_price": 150000,
      "countdown_seconds": 30,
      "expires_at": "2023-10-25T08:00:30Z",
      "created_at": "2023-10-25T08:00:00Z"
    }
  ]
}
```

---

### 17.2 `GET /worker/orders/incoming/{order_id}`

> Get full detail of an incoming order. Maps to "Detail Order" screen showing map, user info (Budi Santoso, Diterima, 2.4 km), address (Jl. Merdeka Selatan No. 45, Kec. Menteng, Jakarta Pusat), and "Detail Keluhan" section.

🔒 **Auth Required** — Role: `worker`

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "order_id": "ord-020",
    "order_number": "HD-20231025-020",
    "title": "Pipa Bocor di Dapur",
    "service": {
      "service_id": "svc-002",
      "name": "Perbaikan Pipa & Saluran Air",
      "category": "Pipa",
      "icon_url": "https://cdn.handydirect.id/icons/pipa.png"
    },
    "description": "Pipa wastafel bocor di bagian sambungan bawah. Air merembes cukup deras dari bawah wastafel cuci piring. Sepertinya sambungan pipa P-trap kendor atau pecah. Butuh perbaikan segera karena lantai dapur mulai basah.",
    "status": "pending",
    "urgency": "normal",
    "user": {
      "user_id": "u-005",
      "full_name": "Budi Santoso",
      "avatar_url": "https://cdn.handydirect.id/avatars/budi_001.jpg",
      "phone": "+62812xxxxxxx",
      "total_orders": 5,
      "member_label": "Pelanggan Setia"
    },
    "location": {
      "latitude": -6.186486,
      "longitude": 106.834091,
      "address": "Jl. Merdeka Selatan No. 45",
      "address_detail": "Kec. Menteng, Jakarta Pusat, 10110. (Rumah pagar hitam, cat putih)",
      "distance_km": 2.4
    },
    "schedule": {
      "preferred_date": "2023-10-25",
      "preferred_time_start": "09:00",
      "preferred_time_end": "12:00"
    },
    "photos": [
      "https://cdn.handydirect.id/orders/ord-020/photo_1.jpg"
    ],
    "notes": "Tolong bawa seal tape dan sambungan pipa cadangan",
    "estimated_base_price": 150000,
    "booking_fee": 2000,
    "countdown_seconds": 25,
    "expires_at": "2023-10-25T08:00:30Z",
    "created_at": "2023-10-25T08:00:00Z"
  }
}
```

---

### 17.3 `POST /worker/orders/{order_id}/accept`

> Accept an incoming order. Status changes to `accepted`. Triggers notification to user. Maps to "Slide untuk Terima Order".

🔒 **Auth Required** — Role: `worker`

**Request Body (optional):**

```json
{
  "estimated_arrival_minutes": 15,
  "note": "Saya akan tiba dalam 15 menit"
}
```

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Order berhasil diterima",
  "data": {
    "order_id": "ord-020",
    "status": "accepted",
    "accepted_at": "2023-10-25T08:05:00Z",
    "estimated_arrival_minutes": 15
  }
}
```

**Error `422 Unprocessable Entity`:**

```json
{
  "status": "error",
  "message": "Order sudah expired atau diambil worker lain",
  "error_code": "ORDER_EXPIRED"
}
```

---

### 17.4 `POST /worker/orders/{order_id}/reject`

> Reject an incoming order.

🔒 **Auth Required** — Role: `worker`

**Request Body:**

```json
{
  "reason": "Sedang dalam perjalanan ke order lain",
  "reason_category": "busy"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `reason` | string | ❌ | Rejection reason text |
| `reason_category` | enum | ✅ | `busy` \| `too_far` \| `not_my_expertise` \| `personal` \| `other` |

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Order ditolak",
  "data": {
    "order_id": "ord-020",
    "status": "rejected",
    "rejected_at": "2023-10-25T08:02:00Z"
  }
}
```



---

## 18. Order Management — Worker Side

> Order lifecycle management from the worker perspective. Maps to "Detail Order" (map + user info + Detail Keluhan), "Perbaikan #012 — Sedang Diproses" (work detail, material list, cost estimate), "Menunggu Pembayaran" (invoice summary), and action slides ("Geser untuk Mulai Kerja", "Geser untuk Membuat Tagihan", "Geser untuk Selesai").

---

### 18.1 `GET /worker/orders`

> Get all orders assigned to worker (active, completed, cancelled).

🔒 **Auth Required** — Role: `worker`

**Query Parameters:**

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `status` | string | ❌ | — | `pending` \| `accepted` \| `on_the_way` \| `arrived` \| `in_progress` \| `completed` \| `cancelled` \| `rejected` |
| `date_from` | date | ❌ | — | Filter from date (YYYY-MM-DD) |
| `date_to` | date | ❌ | — | Filter to date (YYYY-MM-DD) |
| `sort_by` | string | ❌ | `created_at` | `created_at` \| `status` \| `total_price` |
| `sort_order` | string | ❌ | `desc` | `asc` \| `desc` |
| `page` | integer | ❌ | `1` | — |
| `per_page` | integer | ❌ | `10` | — |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "order_id": "ord-012",
      "order_number": "HD-20231025-012",
      "title": "Perbaikan Pipa Bocor",
      "status": "in_progress",
      "user": {
        "user_id": "u-005",
        "full_name": "Budi Santoso",
        "avatar_url": "https://cdn.handydirect.id/avatars/budi_001.jpg"
      },
      "service_name": "Perbaikan Pipa & Saluran Air",
      "category_name": "Pipa",
      "location_address": "Jl. Merdeka No.45, Jakarta",
      "total_price": 215000,
      "created_at": "2023-10-25T08:00:00Z",
      "started_at": "2023-10-25T10:00:00Z"
    },
    {
      "order_id": "ord-005",
      "order_number": "HD-20231024-005",
      "title": "Pipe Repair & Cleaning",
      "status": "completed",
      "user": {
        "user_id": "u-003",
        "full_name": "Budi Santoso",
        "avatar_url": "https://cdn.handydirect.id/avatars/budi_003.jpg"
      },
      "service_name": "Pipe Repair & Cleaning",
      "category_name": "Pipa",
      "location_address": "Jl. Sudirman No. 10, Jakarta",
      "total_price": 350000,
      "created_at": "2023-10-24T14:00:00Z",
      "completed_at": "2023-10-24T15:30:00Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 150,
    "total_pages": 15
  }
}
```

---

### 18.2 `GET /worker/orders/{order_id}`

> Get full order detail from worker perspective. Maps to "Perbaikan #012 — Sedang Diproses" screen: DETAIL PEKERJAAN, lokasi, waktu mulai, Biaya Material list (Pipa PVC, Lem Pipa), estimasi total biaya, informasi pelanggan (Ibu Sarah Wijaya), and action buttons.

🔒 **Auth Required** — Role: `worker`

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "order_id": "ord-012",
    "order_number": "#012",
    "title": "Perbaikan Pipa Bocor",
    "status": "in_progress",
    "status_label": "Sedang Diproses",
    "description": "Pipa wastafel bocor di bagian sambungan bawah. Air merembes cukup deras dari bawah wastafel cuci piring. Sepertinya sambungan pipa P-trap kendor atau pecah. Butuh perbaikan segera karena lantai dapur mulai basah.",
    "service": {
      "service_id": "svc-002",
      "name": "Perbaikan Pipa & Saluran Air",
      "category": "Pipa",
      "icon_url": "https://cdn.handydirect.id/icons/pipa.png"
    },
    "user": {
      "user_id": "u-005",
      "full_name": "Ibu Sarah Wijaya",
      "avatar_url": "https://cdn.handydirect.id/avatars/sarah.jpg",
      "phone": "+62812xxxxxxx",
      "total_orders": 8,
      "member_label": "Pelanggan Setia"
    },
    "location": {
      "latitude": -6.186486,
      "longitude": 106.834091,
      "address": "Jl. Merdeka No.45, Jakarta",
      "address_detail": "Kec. Menteng, Jakarta Pusat, 10110. (Rumah pagar hitam, cat putih)",
      "distance_km": 2.4
    },
    "schedule": {
      "started_at": "2023-10-25T10:00:00Z",
      "preferred_date": "2023-10-25",
      "preferred_time_start": "09:00",
      "preferred_time_end": "12:00"
    },
    "photos": [
      "https://cdn.handydirect.id/orders/ord-012/photo_1.jpg"
    ],
    "notes": "Tolong bawa seal tape dan sambungan pipa cadangan",
    "pricing": {
      "base_service_fee": 150000,
      "total_material_cost": 65000,
      "total_additional_cost": 0,
      "grand_total": 215000
    },
    "purchases": [
      {
        "purchase_id": "pur-101",
        "item_name": "Pipa PVC 1/2 Inch",
        "category": "material",
        "quantity": 2,
        "unit": "meter",
        "unit_price": 25000,
        "total_price": 50000,
        "status": "approved",
        "icon_url": "https://cdn.handydirect.id/icons/pipa_material.png"
      },
      {
        "purchase_id": "pur-102",
        "item_name": "Lem Pipa",
        "category": "material",
        "quantity": 1,
        "unit": "kaleng",
        "unit_price": 15000,
        "total_price": 15000,
        "status": "approved",
        "icon_url": "https://cdn.handydirect.id/icons/lem_material.png"
      }
    ],
    "timeline": [
      {
        "event": "order_created",
        "label": "Pesanan Dibuat",
        "timestamp": "2023-10-25T08:00:00Z",
        "is_completed": true
      },
      {
        "event": "order_accepted",
        "label": "Order Diterima",
        "timestamp": "2023-10-25T08:05:00Z",
        "is_completed": true
      },
      {
        "event": "worker_on_the_way",
        "label": "Menuju Lokasi",
        "timestamp": "2023-10-25T09:00:00Z",
        "is_completed": true
      },
      {
        "event": "worker_arrived",
        "label": "Tiba di Lokasi",
        "timestamp": "2023-10-25T09:15:00Z",
        "is_completed": true
      },
      {
        "event": "work_in_progress",
        "label": "Pengerjaan Dimulai",
        "timestamp": "2023-10-25T10:00:00Z",
        "is_completed": true
      }
    ],
    "available_actions": [
      {
        "action": "add_purchase",
        "label": "+ Tambah",
        "enabled": true
      },
      {
        "action": "generate_invoice",
        "label": "Geser untuk Membuat Tagihan",
        "enabled": true
      },
      {
        "action": "chat",
        "label": "Hubungi Budi (Chat)",
        "enabled": true
      }
    ],
    "can_add_purchase": true,
    "can_generate_invoice": true,
    "can_chat": true,
    "created_at": "2023-10-25T08:00:00Z",
    "updated_at": "2023-10-25T10:00:00Z"
  }
}
```

---

### 18.3 `PATCH /worker/orders/{order_id}/status`

> Update order status. Maps to action slides:
> - `on_the_way` → Worker starts heading to location
> - `arrived` → Worker arrives at location
> - `in_progress` → "Geser untuk Mulai Kerja"
> - `work_paused` → Pause work (e.g., to buy materials)
> - `completed` → "Geser untuk Selesai" (after invoice generated)

🔒 **Auth Required** — Role: `worker`

**Request Body:**

```json
{
  "status": "in_progress",
  "note": "Mulai pengerjaan",
  "latitude": -6.186486,
  "longitude": 106.834091
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `status` | enum | ✅ | New status: `on_the_way` \| `arrived` \| `in_progress` \| `work_paused` \| `completed` |
| `note` | string | ❌ | Status change note |
| `latitude` | number | ❌ | Current latitude (for location verification) |
| `longitude` | number | ❌ | Current longitude |

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Status order berhasil diperbarui",
  "data": {
    "order_id": "ord-012",
    "previous_status": "arrived",
    "new_status": "in_progress",
    "status_label": "Sedang Diproses",
    "updated_at": "2023-10-25T10:00:00Z"
  }
}
```

**Error `422 Unprocessable Entity`:**

```json
{
  "status": "error",
  "message": "Tidak dapat mengubah status ke 'completed' sebelum invoice dibuat",
  "error_code": "INVALID_STATUS_TRANSITION"
}
```

**Valid State Transitions:**

```
pending → accepted (via /accept endpoint)
accepted → on_the_way
on_the_way → arrived
arrived → in_progress
in_progress → work_paused
work_paused → in_progress
in_progress → completed (requires invoice to be generated first)
```

---

### 18.4 `POST /worker/orders/{order_id}/generate-invoice`

> Generate final invoice. Maps to "Geser untuk Membuat Tagihan". Calculates total from base service fee + approved material costs. After invoiced generated, worker can proceed to "Geser untuk Selesai". Shows "Menunggu Pembayaran" screen with breakdown (Biaya Jasa + Biaya Material = Total Tagihan) and Informasi Pelanggan.

🔒 **Auth Required** — Role: `worker`

**Request Body:**

```json
{
  "base_service_fee": 150000,
  "worker_notes": "Pekerjaan selesai. Pipa sudah tidak bocor. Disarankan pengecekan berkala setiap 6 bulan."
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `base_service_fee` | integer | ✅ | Final service fee (Rupiah) |
| `worker_notes` | string | ❌ | Worker's final notes |

**Response `201 Created`:**

```json
{
  "status": "success",
  "message": "Invoice berhasil dibuat",
  "data": {
    "invoice_id": "inv-012",
    "invoice_number": "INV-20231025-012",
    "order_id": "ord-012",
    "base_service_fee": 150000,
    "total_material_cost": 65000,
    "total_additional_cost": 0,
    "booking_fee": 2000,
    "platform_fee": 0,
    "discount_amount": 0,
    "grand_total": 215000,
    "payment_status": "unpaid",
    "payment_instruction": "Silakan kumpulkan pembayaran tunai dari pelanggan sesuai total di bawah.",
    "user": {
      "user_id": "u-005",
      "full_name": "Ibu Sarah Wijaya",
      "avatar_url": "https://cdn.handydirect.id/avatars/sarah.jpg",
      "member_label": "Pelanggan Setia"
    },
    "summary": {
      "biaya_jasa": 150000,
      "biaya_material": 65000,
      "total_tagihan": 215000
    },
    "created_at": "2023-10-25T11:30:00Z"
  }
}
```

---

## 19. Purchase Management — Worker Side (AI-Assisted)

> Worker adds, edits, deletes, and submits material/tool purchases. Includes AI-assisted input processing ("Bantu Rapikan dengan AI") and receipt scanning. Maps to the "+ Tambah" button on the work detail screen, the material list (Pipa PVC, Lem Pipa with delete icons), and AI processing flow.

---

### 19.1 `POST /worker/orders/{order_id}/purchases`

> Add a new purchase item manually. Each field is entered by the worker directly.

🔒 **Auth Required** — Role: `worker`

**Request Body:**

```json
{
  "item_name": "Pipa PVC 1/2 Inch",
  "category": "material",
  "quantity": 2,
  "unit": "meter",
  "unit_price": 25000,
  "total_price": 50000,
  "reason": "Untuk mengganti pipa yang rusak di bawah wastafel",
  "receipt_photo_url": null
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `item_name` | string | ✅ | Item name |
| `category` | enum | ✅ | `material` \| `alat` \| `sparepart` \| `bahan_bangunan` \| `biaya_tambahan` \| `lainnya` |
| `quantity` | number | ✅ | Quantity |
| `unit` | string | ✅ | Unit (pcs, meter, kaleng, sak, liter, etc.) |
| `unit_price` | integer | ✅ | Unit price (Rupiah) |
| `total_price` | integer | ✅ | Total price (Rupiah) |
| `reason` | string | ❌ | Reason for purchase |
| `receipt_photo_url` | string | ❌ | Receipt photo URL |

**Response `201 Created`:**

```json
{
  "status": "success",
  "message": "Pembelian berhasil ditambahkan",
  "data": {
    "purchase_id": "pur-201",
    "item_name": "Pipa PVC 1/2 Inch",
    "category": "material",
    "quantity": 2,
    "unit": "meter",
    "unit_price": 25000,
    "total_price": 50000,
    "reason": "Untuk mengganti pipa yang rusak di bawah wastafel",
    "status": "draft",
    "created_at": "2023-10-25T10:30:00Z"
  }
}
```

---

### 19.2 `POST /worker/orders/{order_id}/purchases/ai-process`

> Send raw text input to AI for structuring — "Bantu Rapikan dengan AI" button. The backend processes the input via OpenAI API and returns structured purchase data.

🔒 **Auth Required** — Role: `worker`

**Request Body:**

```json
{
  "raw_input": "Beli pipa PVC 2 meter 2 pcs 50 ribu, semen 1 sak 65.000, baut lem pipa dan sambungan total 85 ribu",
  "order_context": {
    "service_type": "Perbaikan Pipa",
    "description": "Pipa bocor di bawah wastafel"
  }
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `raw_input` | string | ✅ | Raw text input from worker |
| `order_context` | object | ❌ | Context for AI (service type, description) to improve accuracy |

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Input berhasil diproses oleh AI",
  "data": {
    "order_id": "ord-012",
    "items": [
      {
        "item_name": "Pipa PVC",
        "category": "material",
        "quantity": 2,
        "unit": "pcs",
        "unit_price": 25000,
        "total_price": 50000,
        "reason": "Untuk mengganti bagian pipa yang rusak",
        "confidence": 0.95,
        "needs_clarification": false,
        "clarification_question": null
      },
      {
        "item_name": "Semen",
        "category": "bahan_bangunan",
        "quantity": 1,
        "unit": "sak",
        "unit_price": 65000,
        "total_price": 65000,
        "reason": "Untuk menambal area sekitar pipa",
        "confidence": 0.92,
        "needs_clarification": false,
        "clarification_question": null
      },
      {
        "item_name": "Baut, Lem Pipa, Sambungan",
        "category": "material",
        "quantity": 1,
        "unit": "paket",
        "unit_price": 85000,
        "total_price": 85000,
        "reason": "Perlengkapan perbaikan pipa",
        "confidence": 0.78,
        "needs_clarification": true,
        "clarification_question": "Input berisi beberapa item dalam satu harga total. Mohon pisahkan masing-masing item: baut (jumlah? harga?), lem pipa (jumlah? harga?), sambungan (jumlah? harga?)."
      }
    ],
    "summary": "Total 3 entri pembelian senilai Rp200.000. 1 entri perlu dipisahkan agar lebih detail.",
    "risk_flags": [
      {
        "type": "data_tidak_lengkap",
        "message": "Item 'baut, lem pipa, dan sambungan' perlu dipisahkan per item untuk transparansi."
      }
    ],
    "approval_status": "draft"
  }
}
```

---

### 19.3 `POST /worker/orders/{order_id}/purchases/receipt-scan`

> Upload receipt/nota photo for OCR + AI processing. Backend extracts text via OCR, then processes with LLM.

🔒 **Auth Required** — Role: `worker`
📎 **Content-Type:** `multipart/form-data`

| Field | Type | Required | Description |
|---|---|---|---|
| `receipt_photo` | file | ✅ | Receipt/nota photo (jpg/png, max 10MB) |
| `notes` | string | ❌ | Worker's notes about this receipt |

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Nota berhasil dipindai dan diproses oleh AI",
  "data": {
    "receipt_photo_url": "https://cdn.handydirect.id/receipts/scan_ord012_001.jpg",
    "ocr_raw_text": "Toko Bangunan Jaya\nPipa PVC 1/2\" 2m x Rp25.000 = Rp50.000\nLem Pipa 1 klg = Rp15.000\nTotal: Rp65.000",
    "items": [
      {
        "item_name": "Pipa PVC 1/2 Inch",
        "category": "material",
        "quantity": 2,
        "unit": "meter",
        "unit_price": 25000,
        "total_price": 50000,
        "reason": "Pembelian dari Toko Bangunan Jaya",
        "confidence": 0.93,
        "needs_clarification": false,
        "clarification_question": null
      },
      {
        "item_name": "Lem Pipa",
        "category": "material",
        "quantity": 1,
        "unit": "kaleng",
        "unit_price": 15000,
        "total_price": 15000,
        "reason": "Pembelian dari Toko Bangunan Jaya",
        "confidence": 0.91,
        "needs_clarification": false,
        "clarification_question": null
      }
    ],
    "summary": "Nota dari Toko Bangunan Jaya dengan total Rp65.000 untuk 2 item material.",
    "risk_flags": [],
    "approval_status": "draft"
  }
}
```

---

### 19.4 `PUT /worker/orders/{order_id}/purchases/{purchase_id}`

> Edit an existing purchase item (e.g., after AI processing, worker can correct data).

🔒 **Auth Required** — Role: `worker`

**Request Body:**

```json
{
  "item_name": "Pipa PVC 1/2 Inch",
  "category": "material",
  "quantity": 2,
  "unit": "meter",
  "unit_price": 25000,
  "total_price": 50000,
  "reason": "Untuk mengganti pipa P-trap yang pecah di bawah wastafel"
}
```

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Pembelian berhasil diperbarui",
  "data": {
    "purchase_id": "pur-201",
    "item_name": "Pipa PVC 1/2 Inch",
    "category": "material",
    "quantity": 2,
    "unit": "meter",
    "unit_price": 25000,
    "total_price": 50000,
    "reason": "Untuk mengganti pipa P-trap yang pecah di bawah wastafel",
    "status": "draft",
    "updated_at": "2023-10-25T10:35:00Z"
  }
}
```

---

### 19.5 `DELETE /worker/orders/{order_id}/purchases/{purchase_id}`

> Delete a purchase item. Maps to the trash/delete icon (🗑️) next to each material item in the mockup.

🔒 **Auth Required** — Role: `worker`

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Pembelian berhasil dihapus"
}
```

**Error `422 Unprocessable Entity`:**

```json
{
  "status": "error",
  "message": "Pembelian yang sudah disetujui user tidak dapat dihapus",
  "error_code": "PURCHASE_ALREADY_APPROVED"
}
```

---

### 19.6 `POST /worker/orders/{order_id}/purchases/{purchase_id}/submit`

> Submit a draft purchase for user approval. Changes status from `draft` to `pending_approval`.

🔒 **Auth Required** — Role: `worker`

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Pembelian berhasil dikirim untuk persetujuan user",
  "data": {
    "purchase_id": "pur-201",
    "status": "pending_approval",
    "submitted_at": "2023-10-25T10:40:00Z"
  }
}
```

---

### 19.7 `POST /worker/orders/{order_id}/purchases/bulk-submit`

> Submit all draft purchases at once for user approval.

🔒 **Auth Required** — Role: `worker`

**Request Body:**

```json
{
  "purchase_ids": ["pur-201", "pur-202", "pur-203"]
}
```

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "3 pembelian berhasil dikirim untuk persetujuan user",
  "data": {
    "submitted_count": 3,
    "submitted_ids": ["pur-201", "pur-202", "pur-203"]
  }
}
```

---

### 19.8 `PATCH /worker/orders/{order_id}/purchases/{purchase_id}/clarify-response`

> Respond to a user's clarification request on a purchase item.

🔒 **Auth Required** — Role: `worker`

**Request Body:**

```json
{
  "response": "Alat yang dibeli adalah kunci pipa ukuran 14 inch, diperlukan untuk membuka mur sambungan pipa yang sudah karatan.",
  "updated_item_name": "Kunci Pipa 14 Inch",
  "updated_reason": "Diperlukan untuk membuka mur sambungan pipa yang sudah karatan"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `response` | string | ✅ | Clarification response text |
| `updated_item_name` | string | ❌ | Updated item name (if changed) |
| `updated_reason` | string | ❌ | Updated reason (if changed) |

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Klarifikasi berhasil dikirim. Menunggu persetujuan user.",
  "data": {
    "purchase_id": "pur-203",
    "status": "pending_approval",
    "clarification_response": "Alat yang dibeli adalah kunci pipa ukuran 14 inch...",
    "updated_at": "2023-10-25T10:50:00Z"
  }
}
```



---

## 20. Chat — Worker Side

> Chat functionality from the worker perspective. Same data structure as user side; the `sender_type` field differentiates messages. Uses the same WebSocket endpoint as Section 10.5.

---

### 20.1 `GET /worker/orders/{order_id}/chat/messages`

> Get chat history for an order from worker perspective.

🔒 **Auth Required** — Role: `worker`

**Query Parameters:**

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `before` | string (ISO 8601) | ❌ | — | Cursor: fetch messages before this timestamp |
| `limit` | integer | ❌ | `50` | Messages per request (max 100) |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "order_id": "ord-012",
    "user": {
      "user_id": "u-005",
      "full_name": "Budi Santoso",
      "avatar_url": "https://cdn.handydirect.id/avatars/budi_001.jpg",
      "status_label": "Pelanggan"
    },
    "messages": [
      {
        "message_id": "msg-001",
        "sender_type": "user",
        "sender_name": "Budi Santoso",
        "content": "Halo Pak Ahmad, apakah sudah di jalan?",
        "message_type": "text",
        "media_url": null,
        "is_read": true,
        "read_at": "2023-10-25T10:15:30Z",
        "created_at": "2023-10-25T10:15:00Z"
      },
      {
        "message_id": "msg-002",
        "sender_type": "worker",
        "sender_name": "Ahmad Jaelani",
        "content": "Halo Pak, iya saya sudah di jalan. Mungkin sekitar 15 menit lagi sampai.",
        "message_type": "text",
        "media_url": null,
        "is_read": true,
        "read_at": "2023-10-25T10:17:15Z",
        "created_at": "2023-10-25T10:17:00Z"
      }
    ],
    "has_more": false
  }
}
```

---

### 20.2 `POST /worker/orders/{order_id}/chat/messages`

> Send a chat message from worker.

🔒 **Auth Required** — Role: `worker`

**Request Body (text):**

```json
{
  "content": "Saya sudah sampai di depan rumah. Rumah pagar hitam cat putih kan Pak?",
  "message_type": "text"
}
```

**Request Body (image — `multipart/form-data`):**

| Field | Type | Required | Description |
|---|---|---|---|
| `message_type` | string | ✅ | `image` |
| `content` | string | ❌ | Image caption |
| `media` | file | ✅ | Image file (jpg/png, max 10MB) |

**Response `201 Created`:**

```json
{
  "status": "success",
  "data": {
    "message_id": "msg-010",
    "sender_type": "worker",
    "sender_name": "Ahmad Jaelani",
    "content": "Saya sudah sampai di depan rumah. Rumah pagar hitam cat putih kan Pak?",
    "message_type": "text",
    "media_url": null,
    "created_at": "2023-10-25T10:30:00Z"
  }
}
```

---

### 20.3 `PATCH /worker/orders/{order_id}/chat/read`

> Mark all messages as read by worker.

🔒 **Auth Required** — Role: `worker`

**Response `200 OK`:**

```json
{
  "status": "success",
  "message": "Semua pesan ditandai sudah dibaca"
}
```

---

### 20.4 `GET /worker/chats`

> Get all active chat conversations for worker. Used in "Chat" tab on worker bottom navigation.

🔒 **Auth Required** — Role: `worker`

**Query Parameters:**

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `page` | integer | ❌ | `1` | — |
| `per_page` | integer | ❌ | `20` | — |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "order_id": "ord-012",
      "user": {
        "user_id": "u-005",
        "full_name": "Budi Santoso",
        "avatar_url": "https://cdn.handydirect.id/avatars/budi_001.jpg",
        "is_online": true
      },
      "order_title": "Perbaikan Pipa Bocor",
      "order_status": "in_progress",
      "last_message": {
        "content": "Saya sudah sampai di depan rumah...",
        "sender_type": "worker",
        "message_type": "text",
        "created_at": "2023-10-25T10:30:00Z"
      },
      "unread_count": 1
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total": 2,
    "total_pages": 1
  }
}
```

> **Note:** WebSocket chat events use the same endpoint and format as [Section 10.5](#105-websocket--real-time-chat). The `sender_type` field is `worker` for outgoing and `user` for incoming messages.

---

## 21. Rating — Worker Side (Rate Customer)

> Worker rates the customer after job completion. Maps to the "Beri Rating Konsumen" mockup: showing customer avatar (Ibu Sarah Wijaya), star rating, text input "Catatan tentang Konsumen", and tag chips (Ramah, Lokasi Akurat, Pembayaran Cepat).

---

### 21.1 `POST /worker/orders/{order_id}/customer-rating`

> Submit rating and review for the customer.

🔒 **Auth Required** — Role: `worker`

**Request Body:**

```json
{
  "rating": 5,
  "comment": "Lokasi mudah ditemukan, sangat ramah dan pembayaran langsung tunai tanpa masalah.",
  "tags": ["Ramah", "Lokasi Akurat", "Pembayaran Cepat"]
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `rating` | integer | ✅ | 1–5 stars |
| `comment` | string | ❌ | Comment about customer (max 500 chars) |
| `tags` | array | ❌ | Quick tags: `Ramah`, `Lokasi Akurat`, `Pembayaran Cepat`, `Responsif`, `Kooperatif` |

**Response `201 Created`:**

```json
{
  "status": "success",
  "message": "Rating konsumen berhasil dikirim",
  "data": {
    "customer_review_id": "crev-001",
    "order_id": "ord-012",
    "user_id": "u-005",
    "rating": 5,
    "comment": "Lokasi mudah ditemukan, sangat ramah dan pembayaran langsung tunai tanpa masalah.",
    "tags": ["Ramah", "Lokasi Akurat", "Pembayaran Cepat"],
    "created_at": "2023-10-25T12:10:00Z"
  }
}
```

**Error `409 Conflict`:**

```json
{
  "status": "error",
  "message": "Anda sudah memberikan rating untuk konsumen pada pesanan ini",
  "error_code": "CUSTOMER_RATING_ALREADY_EXISTS"
}
```

---

### 21.2 `GET /worker/orders/{order_id}/customer-rating`

> Get the customer rating already submitted by worker for this order.

🔒 **Auth Required** — Role: `worker`

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "customer_review_id": "crev-001",
    "order_id": "ord-012",
    "user_id": "u-005",
    "user_name": "Ibu Sarah Wijaya",
    "user_avatar_url": "https://cdn.handydirect.id/avatars/sarah.jpg",
    "rating": 5,
    "comment": "Lokasi mudah ditemukan, sangat ramah dan pembayaran langsung tunai tanpa masalah.",
    "tags": ["Ramah", "Lokasi Akurat", "Pembayaran Cepat"],
    "created_at": "2023-10-25T12:10:00Z"
  }
}
```

---

## 22. History & Statistics — Worker Side

> Worker's job history and performance statistics. Maps to the "Orders" tab mockup showing "Pendapatan Bulanan: Rp 8.450.000", filter tabs (Bulanan/Mingguan/Harian), and Job History list with each job's date, time, price, and COMPLETED status.

---

### 22.1 `GET /worker/history`

> Get worker's order history with earnings summary and period filter.

🔒 **Auth Required** — Role: `worker`

**Query Parameters:**

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `period` | string | ❌ | `monthly` | `daily` \| `weekly` \| `monthly` |
| `status` | string | ❌ | — | `completed` \| `cancelled` \| `in_progress` |
| `date_from` | date | ❌ | — | Start date (YYYY-MM-DD) |
| `date_to` | date | ❌ | — | End date (YYYY-MM-DD) |
| `sort_by` | string | ❌ | `created_at` | `created_at` \| `total_price` |
| `sort_order` | string | ❌ | `desc` | `asc` \| `desc` |
| `page` | integer | ❌ | `1` | — |
| `per_page` | integer | ❌ | `10` | — |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "earnings_summary": {
      "period": "monthly",
      "period_label": "Pendapatan Bulanan",
      "total_earnings": 8450000,
      "currency": "IDR"
    },
    "orders": [
      {
        "order_id": "ord-005",
        "order_number": "HD-20231024-005",
        "title": "Pipe Repair & Cleaning",
        "status": "completed",
        "user": {
          "full_name": "Budi Santoso",
          "avatar_url": "https://cdn.handydirect.id/avatars/budi_001.jpg"
        },
        "service_name": "Pipe Repair & Cleaning",
        "category_name": "Pipa",
        "category_icon_url": "https://cdn.handydirect.id/icons/pipa.png",
        "total_price": 350000,
        "date": "2023-10-24",
        "time_range": "14:00 - 15:30",
        "duration_label": "1.5 jam",
        "created_at": "2023-10-24T14:00:00Z",
        "completed_at": "2023-10-24T15:30:00Z"
      },
      {
        "order_id": "ord-004",
        "order_number": "HD-20231023-004",
        "title": "AC Maintenance",
        "status": "completed",
        "user": {
          "full_name": "Siti Aminah",
          "avatar_url": "https://cdn.handydirect.id/avatars/siti.jpg"
        },
        "service_name": "AC Maintenance",
        "category_name": "AC",
        "category_icon_url": "https://cdn.handydirect.id/icons/ac.png",
        "total_price": 175000,
        "date": "2023-10-23",
        "time_range": "09:00 - 10:15",
        "duration_label": "1.25 jam",
        "created_at": "2023-10-23T09:00:00Z",
        "completed_at": "2023-10-23T10:15:00Z"
      },
      {
        "order_id": "ord-003",
        "order_number": "HD-20231021-003",
        "title": "Interior Wall Painting",
        "status": "completed",
        "user": {
          "full_name": "Andi Wijaya",
          "avatar_url": "https://cdn.handydirect.id/avatars/andi.jpg"
        },
        "service_name": "Interior Wall Painting",
        "category_name": "Cat",
        "category_icon_url": "https://cdn.handydirect.id/icons/cat.png",
        "total_price": 1200000,
        "date": "2023-10-21",
        "time_range": "Full Day",
        "duration_label": "Full Day",
        "created_at": "2023-10-21T08:00:00Z",
        "completed_at": "2023-10-21T17:00:00Z"
      },
      {
        "order_id": "ord-002",
        "order_number": "HD-20231020-002",
        "title": "Door Hinge Replacement",
        "status": "completed",
        "user": {
          "full_name": "Linda Hartono",
          "avatar_url": "https://cdn.handydirect.id/avatars/linda.jpg"
        },
        "service_name": "Door Hinge Replacement",
        "category_name": "Kayu",
        "category_icon_url": "https://cdn.handydirect.id/icons/kayu.png",
        "total_price": 120000,
        "date": "2023-10-20",
        "time_range": "16:30 - 17:15",
        "duration_label": "45 menit",
        "created_at": "2023-10-20T16:30:00Z",
        "completed_at": "2023-10-20T17:15:00Z"
      }
    ]
  },
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 42,
    "total_pages": 5
  }
}
```

---

### 22.2 `GET /worker/statistics`

> Get worker performance statistics.

🔒 **Auth Required** — Role: `worker`

**Query Parameters:**

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `period` | string | ❌ | `monthly` | `daily` \| `weekly` \| `monthly` \| `yearly` \| `all_time` |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "period": "monthly",
    "period_label": "Bulan Ini",
    "earnings": {
      "total": 8450000,
      "service_fees": 6300000,
      "material_reimbursements": 2150000,
      "bonuses": 0,
      "currency": "IDR"
    },
    "jobs": {
      "total": 42,
      "completed": 40,
      "cancelled": 1,
      "in_progress": 1,
      "rejected": 3
    },
    "performance": {
      "acceptance_rate": 95,
      "completion_rate": 98,
      "average_rating": 4.9,
      "total_reviews": 38,
      "on_time_arrival_rate": 92,
      "average_response_time_seconds": 25
    },
    "top_services": [
      { "service_name": "Perbaikan Pipa", "count": 15, "earnings": 3200000 },
      { "service_name": "Servis AC", "count": 12, "earnings": 2800000 },
      { "service_name": "Instalasi Listrik", "count": 8, "earnings": 1650000 }
    ],
    "daily_breakdown": [
      { "date": "2023-10-25", "jobs": 2, "earnings": 350000 },
      { "date": "2023-10-24", "jobs": 1, "earnings": 350000 },
      { "date": "2023-10-23", "jobs": 2, "earnings": 295000 }
    ]
  }
}
```

---

## 23. Wallet — Worker Side

> Worker's digital wallet for earnings and withdrawals. Maps to "Penarikan" menu on worker home, showing saldo wallet balance and transaction history.

---

### 23.1 `GET /worker/wallet`

> Get wallet balance and summary.

🔒 **Auth Required** — Role: `worker`

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "wallet_id": "wal-001",
    "balance": 150000,
    "total_earnings": 8450000,
    "total_withdrawn": 8300000,
    "pending_earnings": 215000,
    "is_active": true,
    "currency": "IDR",
    "updated_at": "2023-10-25T11:45:00Z"
  }
}
```

---

### 23.2 `GET /worker/wallet/transactions`

> Get wallet transaction history.

🔒 **Auth Required** — Role: `worker`

**Query Parameters:**

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `type` | string | ❌ | — | `earning` \| `withdrawal` \| `refund` \| `bonus` \| `fee` |
| `status` | string | ❌ | — | `pending` \| `completed` \| `failed` \| `cancelled` |
| `date_from` | date | ❌ | — | Start date |
| `date_to` | date | ❌ | — | End date |
| `page` | integer | ❌ | `1` | — |
| `per_page` | integer | ❌ | `20` | — |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": [
    {
      "transaction_id": "tx-001",
      "type": "earning",
      "amount": 350000,
      "balance_before": 8100000,
      "balance_after": 8450000,
      "description": "Pendapatan dari order #HD-20231024-005 (Pipe Repair & Cleaning)",
      "order_id": "ord-005",
      "order_number": "HD-20231024-005",
      "status": "completed",
      "completed_at": "2023-10-24T16:00:00Z",
      "created_at": "2023-10-24T15:30:00Z"
    },
    {
      "transaction_id": "tx-002",
      "type": "withdrawal",
      "amount": 5000000,
      "balance_before": 5150000,
      "balance_after": 150000,
      "description": "Penarikan ke rekening BCA ****1234",
      "order_id": null,
      "order_number": null,
      "status": "completed",
      "completed_at": "2023-10-23T14:00:00Z",
      "created_at": "2023-10-23T12:00:00Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total": 85,
    "total_pages": 5
  }
}
```

---

### 23.3 `POST /worker/wallet/withdraw`

> Create a withdrawal request. Maps to "Penarikan — Tarik dana ke bank".

🔒 **Auth Required** — Role: `worker`

**Request Body:**

```json
{
  "amount": 5000000,
  "bank_name": "BCA",
  "account_number": "1234567890",
  "account_holder_name": "Ahmad Jaelani",
  "notes": "Penarikan bulanan"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `amount` | integer | ✅ | Withdrawal amount (Rupiah, min 50000) |
| `bank_name` | string | ✅ | Bank name (BCA, BNI, BRI, Mandiri, etc.) |
| `account_number` | string | ✅ | Bank account number |
| `account_holder_name` | string | ✅ | Account holder name |
| `notes` | string | ❌ | Notes |

**Response `201 Created`:**

```json
{
  "status": "success",
  "message": "Permintaan penarikan berhasil dibuat. Dana akan diproses dalam 1-3 hari kerja.",
  "data": {
    "transaction_id": "tx-100",
    "type": "withdrawal",
    "amount": 5000000,
    "balance_before": 5150000,
    "balance_after": 150000,
    "status": "pending",
    "estimated_completion": "2023-10-28T00:00:00Z",
    "created_at": "2023-10-25T12:00:00Z"
  }
}
```

**Error `422 Unprocessable Entity`:**

```json
{
  "status": "error",
  "message": "Saldo tidak mencukupi untuk penarikan. Saldo saat ini: Rp150.000",
  "error_code": "INSUFFICIENT_BALANCE"
}
```

---

## 24. Worker Location Updates

> Worker sends periodic location updates while en route or during a job. These updates are pushed to the user via WebSocket.

---

### 24.1 `PUT /worker/location`

> Update worker's current location. Called periodically (every 5-10 seconds) by the mobile app when worker is on the way to a job.

🔒 **Auth Required** — Role: `worker`

**Request Body:**

```json
{
  "latitude": -6.188000,
  "longitude": 106.832000,
  "heading": 180.0,
  "speed_kmh": 30,
  "accuracy_meters": 5.0,
  "order_id": "ord-012"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `latitude` | number | ✅ | Current latitude |
| `longitude` | number | ✅ | Current longitude |
| `heading` | number | ❌ | Compass heading in degrees (0-360) |
| `speed_kmh` | number | ❌ | Current speed in km/h |
| `accuracy_meters` | number | ❌ | GPS accuracy in meters |
| `order_id` | uuid | ❌ | Active order ID (if en route to order) |

**Response `200 OK`:**

```json
{
  "status": "success",
  "data": {
    "latitude": -6.188000,
    "longitude": 106.832000,
    "eta_minutes": 5,
    "distance_remaining_km": 1.2,
    "updated_at": "2023-10-25T09:12:00Z"
  }
}
```



---

# Part D — Shared

---

## 25. WebSocket Events

> Summary of all real-time WebSocket channels and their events.

### 25.1 WebSocket Channels

| Channel | Endpoint | Auth | Description |
|---|---|---|---|
| Tracking | `wss://api.handydirect.id/v1/ws/tracking/{order_id}` | JWT (query param) | Real-time worker location, status changes, purchase notifications |
| Chat | `wss://api.handydirect.id/v1/ws/chat/{order_id}` | JWT (query param) | Real-time messaging between user and worker |

### 25.2 Tracking Events

| Direction | Event | Description |
|---|---|---|
| Server → Client | `location_update` | Worker location updated (lat, lng, heading, eta) |
| Server → Client | `status_change` | Order status changed (e.g., arrived, in_progress) |
| Server → Client | `new_purchase` | Worker added a new purchase item |
| Server → Client | `purchase_status_change` | Purchase approved/rejected by user |
| Server → Client | `order_completed` | Job completed, invoice ready |

### 25.3 Chat Events

| Direction | Event | Description |
|---|---|---|
| Server → Client | `new_message` | New chat message received |
| Server → Client | `typing` | Counterpart is typing |
| Server → Client | `message_read` | Messages have been read |
| Client → Server | `send_message` | Send a text message |
| Client → Server | `typing` | Send typing indicator |

### 25.4 Connection Management

```json
// Connection
{
  "url": "wss://api.handydirect.id/v1/ws/tracking/{order_id}?token=<jwt_token>",
  "protocols": ["v1.handydirect"]
}

// Heartbeat (every 30 seconds)
{ "event": "ping" }
{ "event": "pong" }

// Connection error
{
  "event": "error",
  "data": {
    "code": "TOKEN_EXPIRED",
    "message": "Token expired, please refresh"
  }
}
```

---

## 26. Enums & Shared Schemas

### 26.1 Enums

```yaml
OrderStatus:
  enum: [pending, accepted, on_the_way, arrived, in_progress, work_paused, completed, cancelled, rejected]

PurchaseStatus:
  enum: [draft, pending_approval, approved, rejected, needs_clarification]

PurchaseCategory:
  enum: [material, alat, sparepart, bahan_bangunan, biaya_tambahan, lainnya]

RiskFlagType:
  enum: [harga_tidak_wajar, item_tidak_relevan, data_tidak_lengkap, nota_tidak_jelas, duplikat, alasan_tidak_lengkap]

MessageType:
  enum: [text, image, system]

PaymentMethod:
  enum: [cash, bank_transfer, ewallet]

PaymentStatus:
  enum: [unpaid, pending, paid, refunded]

Urgency:
  enum: [normal, urgent]

UserRole:
  enum: [user, worker, admin]

VerificationStatus:
  enum: [unverified, pending, verified, rejected]

NotificationType:
  enum: [order, purchase, chat, promo, system, payment]

WalletTransactionType:
  enum: [earning, withdrawal, refund, bonus, fee]

WalletTransactionStatus:
  enum: [pending, completed, failed, cancelled]

OrderCancelReason:
  enum: [changed_mind, found_other_worker, too_long, other]

OrderRejectReason:
  enum: [busy, too_far, not_my_expertise, personal, other]

AuditAction:
  enum: [created, ai_processed, submitted, approved, rejected, clarification_requested, clarification_responded, edited, deleted]
```

### 26.2 Standard Response Schema

**Success Response:**

```json
{
  "status": "success",
  "message": "Optional success message",
  "data": {},
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

**Error Response:**

```json
{
  "status": "error",
  "message": "Human-readable error message",
  "error_code": "MACHINE_READABLE_CODE",
  "errors": [
    {
      "field": "email",
      "message": "Email sudah terdaftar"
    }
  ]
}
```

### 26.3 AI Purchase Output Schema

> Structured JSON schema returned by the LLM for purchase tracking processing.

```json
{
  "order_id": "string (UUID)",
  "items": [
    {
      "item_name": "string",
      "category": "material | alat | sparepart | bahan_bangunan | biaya_tambahan | lainnya",
      "quantity": "number",
      "unit": "string",
      "unit_price": "number (integer, Rupiah)",
      "total_price": "number (integer, Rupiah)",
      "reason": "string",
      "confidence": "number (0.00 – 1.00)",
      "needs_clarification": "boolean",
      "clarification_question": "string | null"
    }
  ],
  "summary": "string (human-readable summary of all purchases)",
  "risk_flags": [
    {
      "type": "harga_tidak_wajar | item_tidak_relevan | data_tidak_lengkap | nota_tidak_jelas | duplikat | alasan_tidak_lengkap",
      "message": "string"
    }
  ],
  "approval_status": "draft"
}
```

---

## 27. Error Responses

### 27.1 Standard Error Format

All error responses follow the same structure:

```json
{
  "status": "error",
  "message": "Deskripsi error yang mudah dipahami",
  "error_code": "KODE_ERROR_MESIN",
  "errors": []
}
```

### 27.2 Error Examples by Status Code

**`400 Bad Request` — Invalid Parameters:**

```json
{
  "status": "error",
  "message": "Validasi gagal",
  "error_code": "VALIDATION_ERROR",
  "errors": [
    { "field": "email", "message": "Format email tidak valid" },
    { "field": "password", "message": "Password minimal 8 karakter" }
  ]
}
```

**`401 Unauthorized` — Invalid/Expired Token:**

```json
{
  "status": "error",
  "message": "Token tidak valid atau sudah expired",
  "error_code": "UNAUTHORIZED"
}
```

**`403 Forbidden` — Role Mismatch:**

```json
{
  "status": "error",
  "message": "Anda tidak memiliki akses ke endpoint ini. Dibutuhkan role: worker",
  "error_code": "FORBIDDEN"
}
```

**`404 Not Found`:**

```json
{
  "status": "error",
  "message": "Order dengan ID ord-999 tidak ditemukan",
  "error_code": "NOT_FOUND"
}
```

**`409 Conflict` — Duplicate Resource:**

```json
{
  "status": "error",
  "message": "Anda sudah memberikan rating untuk pesanan ini",
  "error_code": "RATING_ALREADY_EXISTS"
}
```

**`422 Unprocessable Entity` — Business Logic Error:**

```json
{
  "status": "error",
  "message": "Pesanan tidak dapat dibatalkan karena worker sudah tiba di lokasi",
  "error_code": "CANCEL_NOT_ALLOWED"
}
```

**`429 Too Many Requests` — Rate Limited:**

```json
{
  "status": "error",
  "message": "Terlalu banyak request. Coba lagi dalam 60 detik.",
  "error_code": "RATE_LIMITED"
}
```

**`500 Internal Server Error`:**

```json
{
  "status": "error",
  "message": "Terjadi kesalahan internal. Silakan coba lagi nanti.",
  "error_code": "INTERNAL_ERROR"
}
```

### 27.3 Error Code Reference

| Error Code | HTTP | Description |
|---|---|---|
| `VALIDATION_ERROR` | 400 | One or more fields failed validation |
| `UNAUTHORIZED` | 401 | Missing, invalid, or expired JWT token |
| `INVALID_CREDENTIALS` | 401 | Wrong email/password combination |
| `FORBIDDEN` | 403 | Authenticated but not authorized (role mismatch) |
| `NOT_FOUND` | 404 | Requested resource does not exist |
| `EMAIL_ALREADY_EXISTS` | 409 | Email already registered |
| `PHONE_ALREADY_EXISTS` | 409 | Phone number already registered |
| `RATING_ALREADY_EXISTS` | 409 | Rating already submitted for this order |
| `CUSTOMER_RATING_ALREADY_EXISTS` | 409 | Worker already rated customer for this order |
| `CANCEL_NOT_ALLOWED` | 422 | Order cannot be cancelled at current status |
| `ORDER_EXPIRED` | 422 | Order acceptance window has expired |
| `INVALID_STATUS_TRANSITION` | 422 | Invalid order status transition |
| `PURCHASE_ALREADY_APPROVED` | 422 | Cannot delete/edit approved purchase |
| `INSUFFICIENT_BALANCE` | 422 | Wallet balance too low for withdrawal |
| `INVOICE_NOT_GENERATED` | 422 | Cannot complete order before invoice is generated |
| `RATE_LIMITED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Unexpected server error |

---

## 28. Endpoint Summary Table

### Part A — Authentication (Any Role)

| # | Method | Endpoint | Role | Description |
|---|---|---|---|---|
| 1 | POST | `/auth/register` | any | Register new account |
| 2 | POST | `/auth/login` | any | Login |
| 3 | POST | `/auth/refresh` | any | Refresh access token |
| 4 | POST | `/auth/logout` | auth | Logout |
| 5 | POST | `/auth/forgot-password` | any | Request password reset |
| 6 | POST | `/auth/reset-password` | any | Reset password |

### Part B — User Profile (Any Authenticated Role)

| # | Method | Endpoint | Role | Description |
|---|---|---|---|---|
| 7 | GET | `/users/me` | auth | Get own profile |
| 8 | PUT | `/users/me` | auth | Update profile |
| 9 | PUT | `/users/me/avatar` | auth | Upload avatar |
| 10 | PUT | `/users/me/location` | auth | Update location |

### Part C — User (Customer) Endpoints

| # | Method | Endpoint | Role | Description |
|---|---|---|---|---|
| 11 | GET | `/home` | user | Home screen data |
| 12 | GET | `/categories` | any | List categories |
| 13 | GET | `/categories/{id}/services` | any | Services in category |
| 14 | GET | `/workers/nearby` | user | Nearby workers |
| 15 | GET | `/workers/search` | user | Search workers |
| 16 | GET | `/workers/{id}` | user | Worker detail |
| 17 | GET | `/workers/{id}/reviews` | user | Worker reviews |
| 18 | GET | `/workers/{id}/services` | user | Worker services |
| 19 | POST | `/orders` | user | Create order |
| 20 | GET | `/orders` | user | List orders |
| 21 | GET | `/orders/{id}` | user | Order detail |
| 22 | POST | `/orders/{id}/cancel` | user | Cancel order |
| 23 | GET | `/orders/{id}/tracking` | user | Tracking data |
| 24 | GET | `/orders/{id}/tracking/location` | user | Worker location (poll) |
| 25 | GET | `/orders/{id}/purchases` | user | List purchases |
| 26 | GET | `/orders/{id}/purchases/{pid}` | user | Purchase detail |
| 27 | PATCH | `/orders/{id}/purchases/{pid}/approve` | user | Approve purchase |
| 28 | PATCH | `/orders/{id}/purchases/{pid}/reject` | user | Reject purchase |
| 29 | PATCH | `/orders/{id}/purchases/{pid}/clarify` | user | Request clarification |
| 30 | PATCH | `/orders/{id}/purchases/bulk-approve` | user | Bulk approve purchases |
| 31 | GET | `/orders/{id}/chat/messages` | user | Chat history |
| 32 | POST | `/orders/{id}/chat/messages` | user | Send message |
| 33 | PATCH | `/orders/{id}/chat/read` | user | Mark chat read |
| 34 | GET | `/chats` | user | List conversations |
| 35 | POST | `/orders/{id}/rating` | user | Rate worker |
| 36 | GET | `/orders/{id}/rating` | user | Get rating |
| 37 | GET | `/orders/{id}/invoice` | user | Get invoice |
| 38 | POST | `/orders/{id}/payment` | user | Confirm payment |
| 39 | GET | `/orders/{id}/invoice/pdf` | user | Download invoice PDF |
| 40 | GET | `/knowledge/articles` | any | List articles |
| 41 | GET | `/knowledge/articles/{id}` | any | Article detail |
| 42 | GET | `/knowledge/faq` | any | List FAQ |
| 43 | GET | `/notifications` | auth | List notifications |
| 44 | PATCH | `/notifications/{id}/read` | auth | Mark notification read |
| 45 | PATCH | `/notifications/read-all` | auth | Mark all read |

### Part D — Worker Endpoints

| # | Method | Endpoint | Role | Description |
|---|---|---|---|---|
| 46 | GET | `/worker/profile` | worker | Get worker profile |
| 47 | PUT | `/worker/profile` | worker | Update worker profile |
| 48 | PUT | `/worker/profile/cover-photo` | worker | Upload cover photo |
| 49 | POST | `/worker/profile/verification` | worker | Submit verification docs |
| 50 | GET | `/worker/profile/verification` | worker | Get verification status |
| 51 | GET | `/worker/home` | worker | Worker home dashboard |
| 52 | PATCH | `/worker/availability` | worker | Toggle availability |
| 53 | GET | `/worker/orders/incoming` | worker | List incoming orders |
| 54 | GET | `/worker/orders/incoming/{id}` | worker | Incoming order detail |
| 55 | POST | `/worker/orders/{id}/accept` | worker | Accept order |
| 56 | POST | `/worker/orders/{id}/reject` | worker | Reject order |
| 57 | GET | `/worker/orders` | worker | List all worker orders |
| 58 | GET | `/worker/orders/{id}` | worker | Order detail (worker view) |
| 59 | PATCH | `/worker/orders/{id}/status` | worker | Update order status |
| 60 | POST | `/worker/orders/{id}/generate-invoice` | worker | Generate invoice |
| 61 | POST | `/worker/orders/{id}/purchases` | worker | Add purchase |
| 62 | POST | `/worker/orders/{id}/purchases/ai-process` | worker | AI process raw input |
| 63 | POST | `/worker/orders/{id}/purchases/receipt-scan` | worker | Scan receipt (OCR + AI) |
| 64 | PUT | `/worker/orders/{id}/purchases/{pid}` | worker | Edit purchase |
| 65 | DELETE | `/worker/orders/{id}/purchases/{pid}` | worker | Delete purchase |
| 66 | POST | `/worker/orders/{id}/purchases/{pid}/submit` | worker | Submit for approval |
| 67 | POST | `/worker/orders/{id}/purchases/bulk-submit` | worker | Bulk submit purchases |
| 68 | PATCH | `/worker/orders/{id}/purchases/{pid}/clarify-response` | worker | Respond to clarification |
| 69 | GET | `/worker/orders/{id}/chat/messages` | worker | Chat history |
| 70 | POST | `/worker/orders/{id}/chat/messages` | worker | Send message |
| 71 | PATCH | `/worker/orders/{id}/chat/read` | worker | Mark chat read |
| 72 | GET | `/worker/chats` | worker | List conversations |
| 73 | POST | `/worker/orders/{id}/customer-rating` | worker | Rate customer |
| 74 | GET | `/worker/orders/{id}/customer-rating` | worker | Get customer rating |
| 75 | GET | `/worker/history` | worker | Order history |
| 76 | GET | `/worker/statistics` | worker | Performance statistics |
| 77 | GET | `/worker/wallet` | worker | Wallet balance |
| 78 | GET | `/worker/wallet/transactions` | worker | Transaction history |
| 79 | POST | `/worker/wallet/withdraw` | worker | Request withdrawal |
| 80 | PUT | `/worker/location` | worker | Update location |

### Part E — WebSocket Channels

| # | Protocol | Endpoint | Role | Description |
|---|---|---|---|---|
| 81 | WSS | `/ws/tracking/{order_id}` | auth | Real-time tracking |
| 82 | WSS | `/ws/chat/{order_id}` | auth | Real-time chat |

---

### 📊 Totals

| Category | Count |
|---|---|
| Authentication | 6 |
| User Profile | 4 |
| User (Customer) Endpoints | 35 |
| Worker Endpoints | 35 |
| WebSocket Channels | 2 |
| **Grand Total** | **82** |

---

> **📌 Document Version:** 1.0.0 | **Last Updated:** 2026-05-30 | **Total REST Endpoints:** 80 | **Total WebSocket Channels:** 2 | **Database Tables:** 24 | **Architecture:** Single Backend with RBAC via JWT


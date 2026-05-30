---

## 📑 Table of Contents

- [1. Architecture Overview](#1-architecture-overview)
  - [Key Principles](#key-principles)
- [2. RBAC \& Authentication Strategy](#2-rbac--authentication-strategy)
  - [JWT Token Payload](#jwt-token-payload)
  - [Middleware Flow](#middleware-flow)
  - [Refresh Token Strategy](#refresh-token-strategy)
- [3. Enum Types](#3-enum-types)
- [4. Entity Relationship Diagram (ERD)](#4-entity-relationship-diagram-erd)
- [5. Table Definitions](#5-table-definitions)
  - [5.1 `users`](#51-users)
  - [5.2 `refresh_tokens`](#52-refresh_tokens)
  - [5.3 `worker_profiles`](#53-worker_profiles)
  - [5.4 `categories`](#54-categories)
  - [5.5 `services`](#55-services)
  - [5.6 `worker_services`](#56-worker_services)
  - [5.7 `orders`](#57-orders)
  - [5.8 `order_photos`](#58-order_photos)
  - [5.9 `order_timeline`](#59-order_timeline)
  - [5.10 `purchases`](#510-purchases)
  - [5.11 `purchase_risk_flags`](#511-purchase_risk_flags)
  - [5.12 `purchase_audit_logs`](#512-purchase_audit_logs)
  - [5.13 `chat_messages`](#513-chat_messages)
  - [5.14 `reviews`](#514-reviews)
  - [5.15 `review_tags`](#515-review_tags)
  - [5.16 `invoices`](#516-invoices)
  - [5.17 `invoice_line_items`](#517-invoice_line_items)
  - [5.18 `payments`](#518-payments)
  - [5.19 `notifications`](#519-notifications)
  - [5.20 `articles`](#520-articles)
  - [5.21 `faqs`](#521-faqs)
  - [5.22 `promotions`](#522-promotions)
  - [5.23 `worker_wallets`](#523-worker_wallets)
  - [5.24 `wallet_transactions`](#524-wallet_transactions)
- [6. Migration Order](#6-migration-order)
- [7. Indexes Summary](#7-indexes-summary)
- [8. Notes \& Conventions](#8-notes--conventions)
  - [8.1 Naming Conventions](#81-naming-conventions)
  - [8.2 Monetary Values](#82-monetary-values)
  - [8.3 Soft Delete](#83-soft-delete)
  - [8.4 Timestamps](#84-timestamps)
  - [8.5 JSONB Columns](#85-jsonb-columns)
  - [8.6 Geospatial](#86-geospatial)
  - [8.7 Trigger: Auto-update `updated_at`](#87-trigger-auto-update-updated_at)
  - [8.8 Required PostgreSQL Extensions](#88-required-postgresql-extensions)

---

## 1. Architecture Overview

HandyDirect menggunakan **satu tabel `users`** untuk semua aktor (User, Worker, Admin). Diferensiasi dilakukan melalui kolom `role` yang disimpan sebagai enum. JWT access token yang dihasilkan saat login akan mengandung informasi `role` sehingga middleware di backend dapat melakukan **Role-Based Access Control (RBAC)** tanpa perlu query tambahan ke database untuk setiap request.

### Key Principles

- **Single `users` table** — semua aktor (user, worker, admin) berada di satu tabel.
- **`worker_profiles` table** — data spesifik worker (spesialisasi, bio, verifikasi) disimpan terpisah dan di-link via `user_id`.
- **UUID v4** — semua primary key menggunakan UUID untuk keamanan dan distribusi.
- **Soft delete** — menggunakan kolom `deleted_at` (nullable TIMESTAMPTZ) pada tabel yang membutuhkan.
- **Audit timestamps** — setiap tabel memiliki `created_at` dan `updated_at`.
- **PostgreSQL enums** — digunakan untuk kolom dengan nilai terbatas dan tetap.

---

## 2. RBAC & Authentication Strategy

### JWT Token Payload

```json
{
  "sub": "user_id (UUID)",
  "role": "user | worker | admin",
  "email": "user@email.com",
  "iat": 1698230400,
  "exp": 1698234000,
  "jti": "unique-token-id"
}
```

### Middleware Flow

```
Request → Extract JWT → Verify Signature → Decode Payload → Check role
  ├── role == "user"   → Allow access to user endpoints
  ├── role == "worker" → Allow access to worker endpoints
  ├── role == "admin"  → Allow access to admin endpoints
  └── role mismatch    → Return 403 Forbidden
```

### Refresh Token Strategy

- Access token: short-lived (1 hour)
- Refresh token: long-lived (30 days), stored in `refresh_tokens` table
- Refresh tokens are rotated on every refresh (old token invalidated)

---

## 3. Enum Types

Berikut adalah semua PostgreSQL enum types yang digunakan:

```sql
-- =============================================
-- ENUM TYPES
-- =============================================

CREATE TYPE user_role AS ENUM ('user', 'worker', 'admin');

CREATE TYPE order_status AS ENUM (
  'pending',
  'accepted',
  'on_the_way',
  'arrived',
  'in_progress',
  'work_paused',
  'completed',
  'cancelled',
  'rejected'
);

CREATE TYPE order_urgency AS ENUM ('normal', 'urgent');

CREATE TYPE purchase_status AS ENUM (
  'draft',
  'pending_approval',
  'approved',
  'rejected',
  'needs_clarification'
);

CREATE TYPE purchase_category AS ENUM (
  'material',
  'alat',
  'sparepart',
  'bahan_bangunan',
  'biaya_tambahan',
  'lainnya'
);

CREATE TYPE risk_flag_type AS ENUM (
  'harga_tidak_wajar',
  'item_tidak_relevan',
  'data_tidak_lengkap',
  'nota_tidak_jelas',
  'duplikat',
  'alasan_tidak_lengkap'
);

CREATE TYPE message_type AS ENUM ('text', 'image', 'system');

CREATE TYPE payment_method AS ENUM ('cash', 'bank_transfer', 'ewallet');

CREATE TYPE payment_status AS ENUM ('unpaid', 'pending', 'paid', 'refunded');

CREATE TYPE verification_status AS ENUM ('unverified', 'pending', 'verified', 'rejected');

CREATE TYPE notification_type AS ENUM ('order', 'purchase', 'chat', 'promo', 'system', 'payment');

CREATE TYPE article_category AS ENUM ('faq', 'guide', 'tips', 'safety', 'payment');

CREATE TYPE faq_category AS ENUM ('general', 'payment', 'tracking', 'security', 'cancellation');

CREATE TYPE wallet_tx_type AS ENUM ('earning', 'withdrawal', 'refund', 'bonus', 'fee');

CREATE TYPE wallet_tx_status AS ENUM ('pending', 'completed', 'failed', 'cancelled');

CREATE TYPE audit_action AS ENUM (
  'created', 'ai_processed', 'submitted',
  'approved', 'rejected', 'clarification_requested',
  'clarification_responded', 'edited', 'deleted'
);
```

---

## 4. Entity Relationship Diagram (ERD)

```
┌──────────────┐       ┌──────────────────┐       ┌─────────────┐
│    users     │1─────1│  worker_profiles  │       │  categories │
│  (all roles) │       │  (worker only)    │       │             │
└──────┬───────┘       └────────┬─────────┘       └──────┬──────┘
       │                        │                         │
       │                        │ M:N                     │ 1:N
       │                 ┌──────┴──────┐           ┌──────┴──────┐
       │                 │worker_services│          │  services   │
       │                 └──────┬──────┘           └──────┬──────┘
       │                        │                         │
       │ 1:N (as user)          │                         │
       │                        │                         │
  ┌────┴────────────────────────┴─────────────────────────┴───┐
  │                          orders                            │
  │  (user_id, worker_id, service_id, category_id)            │
  └───┬───────┬──────────┬──────────┬──────────┬──────────┬───┘
      │       │          │          │          │          │
      │1:N    │1:N       │1:N       │1:N       │1:1       │1:1
      │       │          │          │          │          │
 ┌────┴──┐ ┌──┴────┐ ┌───┴────┐ ┌──┴─────┐ ┌─┴──────┐ ┌┴───────┐
 │order_  │ │order_ │ │purch-  │ │chat_   │ │invoices│ │reviews │
 │photos  │ │time-  │ │ases    │ │messages│ │        │ │        │
 │        │ │line   │ │        │ │        │ │        │ │        │
 └────────┘ └───────┘ └───┬────┘ └────────┘ └───┬────┘ └───┬────┘
                           │                      │          │
                     ┌─────┴──────┐          ┌────┴────┐ ┌───┴─────┐
                     │purchase_   │          │invoice_ │ │review_  │
                     │risk_flags  │          │line_    │ │tags     │
                     │            │          │items    │ │         │
                     ├────────────┤          └─────────┘ └─────────┘
                     │purchase_   │
                     │audit_logs  │
                     └────────────┘

  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
  │notifications │    │  articles    │    │   faqs       │
  │ (per user)   │    │ (knowledge)  │    │              │
  └──────────────┘    └──────────────┘    └──────────────┘

  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
  │ promotions   │    │worker_wallets│───1│wallet_trans- │
  │              │    │ (per worker) │ :N │actions       │
  └──────────────┘    └──────────────┘    └──────────────┘

  ┌──────────────┐
  │refresh_tokens│
  │ (per user)   │
  └──────────────┘
```

---

## 5. Table Definitions

### 5.1 `users`

> Tabel utama untuk semua aktor. Kolom `role` digunakan sebagai dasar RBAC di JWT middleware.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `full_name` | VARCHAR(255) | NO | — | Nama lengkap |
| `email` | VARCHAR(255) | NO | — | Email (unique) |
| `phone` | VARCHAR(20) | NO | — | Nomor telepon (unique) |
| `password_hash` | VARCHAR(255) | NO | — | Hashed password (bcrypt) |
| `role` | user_role | NO | `'user'` | Role: user, worker, admin |
| `avatar_url` | TEXT | YES | NULL | URL foto profil |
| `address` | TEXT | YES | NULL | Alamat lengkap |
| `latitude` | DECIMAL(10,7) | YES | NULL | Latitude lokasi |
| `longitude` | DECIMAL(10,7) | YES | NULL | Longitude lokasi |
| `is_active` | BOOLEAN | NO | `TRUE` | Status aktif akun |
| `email_verified_at` | TIMESTAMPTZ | YES | NULL | Waktu verifikasi email |
| `phone_verified_at` | TIMESTAMPTZ | YES | NULL | Waktu verifikasi phone |
| `last_login_at` | TIMESTAMPTZ | YES | NULL | Login terakhir |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu diupdate |
| `deleted_at` | TIMESTAMPTZ | YES | NULL | Soft delete |

**Constraints & Indexes:**
- PK: `id`
- UNIQUE: `email`, `phone`
- INDEX: `role`, `is_active`, `latitude + longitude` (GIST), `deleted_at`

```sql
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name       VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL UNIQUE,
    phone           VARCHAR(20) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    role            user_role NOT NULL DEFAULT 'user',
    avatar_url      TEXT,
    address         TEXT,
    latitude        DECIMAL(10, 7),
    longitude       DECIMAL(10, 7),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    email_verified_at TIMESTAMPTZ,
    phone_verified_at TIMESTAMPTZ,
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_active ON users(is_active);
CREATE INDEX idx_users_location ON users USING GIST (
    ST_MakePoint(longitude, latitude)
) WHERE latitude IS NOT NULL AND longitude IS NOT NULL;
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

---

### 5.2 `refresh_tokens`

> Menyimpan refresh token untuk JWT rotation strategy.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `user_id` | UUID | NO | — | FK → users.id |
| `token_hash` | VARCHAR(255) | NO | — | Hashed refresh token |
| `device_info` | VARCHAR(500) | YES | NULL | Info device/browser |
| `ip_address` | INET | YES | NULL | IP address saat login |
| `expires_at` | TIMESTAMPTZ | NO | — | Waktu expired |
| `revoked_at` | TIMESTAMPTZ | YES | NULL | Waktu revoke (jika di-revoke) |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |

```sql
CREATE TABLE refresh_tokens (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash      VARCHAR(255) NOT NULL,
    device_info     VARCHAR(500),
    ip_address      INET,
    expires_at      TIMESTAMPTZ NOT NULL,
    revoked_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens(expires_at);
```

---

### 5.3 `worker_profiles`

> Data spesifik worker. Satu user dengan role='worker' memiliki tepat satu worker_profile.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `user_id` | UUID | NO | — | FK → users.id (UNIQUE) |
| `specialization` | VARCHAR(255) | YES | NULL | Spesialisasi utama |
| `bio` | TEXT | YES | NULL | Deskripsi pengalaman |
| `cover_photo_url` | TEXT | YES | NULL | Foto cover profil |
| `verification_status` | verification_status | NO | `'unverified'` | Status verifikasi |
| `id_card_url` | TEXT | YES | NULL | Foto KTP (untuk verifikasi) |
| `certificate_urls` | JSONB | YES | `'[]'` | Array URL sertifikat |
| `base_price` | INTEGER | YES | NULL | Harga dasar (Rupiah) |
| `price_unit` | VARCHAR(50) | YES | `'per kunjungan'` | Satuan harga |
| `booking_fee` | INTEGER | NO | `2000` | Booking fee (Rupiah) |
| `rating_avg` | DECIMAL(2,1) | NO | `0.0` | Rata-rata rating |
| `total_reviews` | INTEGER | NO | `0` | Total ulasan |
| `completed_jobs` | INTEGER | NO | `0` | Total pekerjaan selesai |
| `is_available` | BOOLEAN | NO | `TRUE` | Sedang tersedia/tidak |
| `verified_at` | TIMESTAMPTZ | YES | NULL | Waktu diverifikasi |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu diupdate |

```sql
CREATE TABLE worker_profiles (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    specialization      VARCHAR(255),
    bio                 TEXT,
    cover_photo_url     TEXT,
    verification_status verification_status NOT NULL DEFAULT 'unverified',
    id_card_url         TEXT,
    certificate_urls    JSONB DEFAULT '[]',
    base_price          INTEGER,
    price_unit          VARCHAR(50) DEFAULT 'per kunjungan',
    booking_fee         INTEGER NOT NULL DEFAULT 2000,
    rating_avg          DECIMAL(2, 1) NOT NULL DEFAULT 0.0,
    total_reviews       INTEGER NOT NULL DEFAULT 0,
    completed_jobs      INTEGER NOT NULL DEFAULT 0,
    is_available        BOOLEAN NOT NULL DEFAULT TRUE,
    verified_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_worker_profiles_user ON worker_profiles(user_id);
CREATE INDEX idx_worker_profiles_available ON worker_profiles(is_available);
CREATE INDEX idx_worker_profiles_rating ON worker_profiles(rating_avg DESC);
CREATE INDEX idx_worker_profiles_verification ON worker_profiles(verification_status);
```

---

### 5.4 `categories`

> Kategori jasa: AC, Pipa, Atap, Listrik, Kunci, Kayu, Cat, Kebun, dll.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `name` | VARCHAR(100) | NO | — | Nama kategori |
| `slug` | VARCHAR(100) | NO | — | Slug URL-friendly (unique) |
| `icon_url` | TEXT | YES | NULL | URL ikon kategori |
| `description` | TEXT | YES | NULL | Deskripsi kategori |
| `display_order` | INTEGER | NO | `0` | Urutan tampilan |
| `is_active` | BOOLEAN | NO | `TRUE` | Status aktif |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu diupdate |

```sql
CREATE TABLE categories (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    slug            VARCHAR(100) NOT NULL UNIQUE,
    icon_url        TEXT,
    description     TEXT,
    display_order   INTEGER NOT NULL DEFAULT 0,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_active ON categories(is_active, display_order);
```

---

### 5.5 `services`

> Layanan spesifik dalam kategori. Contoh: kategori "Listrik" → "Instalasi Listrik", "Perbaikan Listrik".

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `category_id` | UUID | NO | — | FK → categories.id |
| `name` | VARCHAR(255) | NO | — | Nama layanan |
| `slug` | VARCHAR(255) | NO | — | Slug (unique) |
| `description` | TEXT | YES | NULL | Deskripsi layanan |
| `icon_url` | TEXT | YES | NULL | URL ikon |
| `base_price` | INTEGER | YES | NULL | Harga dasar (Rupiah) |
| `price_unit` | VARCHAR(50) | YES | `'per kunjungan'` | Satuan harga |
| `estimated_duration` | VARCHAR(50) | YES | NULL | Estimasi durasi (e.g. "1-3 jam") |
| `is_active` | BOOLEAN | NO | `TRUE` | Status aktif |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu diupdate |

```sql
CREATE TABLE services (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id         UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    name                VARCHAR(255) NOT NULL,
    slug                VARCHAR(255) NOT NULL UNIQUE,
    description         TEXT,
    icon_url            TEXT,
    base_price          INTEGER,
    price_unit          VARCHAR(50) DEFAULT 'per kunjungan',
    estimated_duration  VARCHAR(50),
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_services_category ON services(category_id);
CREATE INDEX idx_services_slug ON services(slug);
CREATE INDEX idx_services_active ON services(is_active);
```

---

### 5.6 `worker_services`

> Junction table: relasi many-to-many antara worker dan services yang mereka tawarkan.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `worker_id` | UUID | NO | — | FK → users.id (role=worker) |
| `service_id` | UUID | NO | — | FK → services.id |
| `custom_price` | INTEGER | YES | NULL | Harga kustom worker (override) |
| `is_active` | BOOLEAN | NO | `TRUE` | Masih ditawarkan? |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |

```sql
CREATE TABLE worker_services (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    worker_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    service_id      UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    custom_price    INTEGER,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(worker_id, service_id)
);

CREATE INDEX idx_worker_services_worker ON worker_services(worker_id);
CREATE INDEX idx_worker_services_service ON worker_services(service_id);
```

---

### 5.7 `orders`

> Tabel utama pesanan/laporan dari user ke worker.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `order_number` | VARCHAR(50) | NO | — | Nomor order (unique, e.g. HD-20231025-001) |
| `user_id` | UUID | NO | — | FK → users.id (pemesan) |
| `worker_id` | UUID | NO | — | FK → users.id (pekerja) |
| `service_id` | UUID | NO | — | FK → services.id |
| `category_id` | UUID | NO | — | FK → categories.id |
| `title` | VARCHAR(255) | NO | — | Judul laporan/pesanan |
| `description` | TEXT | NO | — | Deskripsi masalah |
| `status` | order_status | NO | `'pending'` | Status pesanan |
| `urgency` | order_urgency | NO | `'normal'` | Tingkat urgensi |
| `location_address` | TEXT | NO | — | Alamat lokasi pekerjaan |
| `location_detail` | VARCHAR(500) | YES | NULL | Detail tambahan lokasi |
| `location_lat` | DECIMAL(10,7) | NO | — | Latitude lokasi |
| `location_lng` | DECIMAL(10,7) | NO | — | Longitude lokasi |
| `preferred_date` | DATE | YES | NULL | Tanggal yang diinginkan |
| `preferred_time_start` | TIME | YES | NULL | Jam mulai yang diinginkan |
| `preferred_time_end` | TIME | YES | NULL | Jam selesai yang diinginkan |
| `notes` | TEXT | YES | NULL | Catatan tambahan |
| `booking_fee` | INTEGER | NO | `2000` | Booking fee (Rupiah) |
| `base_service_fee` | INTEGER | YES | NULL | Biaya jasa dasar |
| `total_material_cost` | INTEGER | NO | `0` | Total biaya material (approved) |
| `total_additional_cost` | INTEGER | NO | `0` | Total biaya tambahan (approved) |
| `grand_total` | INTEGER | YES | NULL | Grand total (dihitung saat selesai) |
| `cancellation_reason` | TEXT | YES | NULL | Alasan pembatalan |
| `cancellation_category` | VARCHAR(50) | YES | NULL | Kategori alasan batal |
| `cancelled_by` | UUID | YES | NULL | Siapa yang membatalkan |
| `accepted_at` | TIMESTAMPTZ | YES | NULL | Waktu order diterima |
| `started_at` | TIMESTAMPTZ | YES | NULL | Waktu pengerjaan dimulai |
| `completed_at` | TIMESTAMPTZ | YES | NULL | Waktu selesai |
| `cancelled_at` | TIMESTAMPTZ | YES | NULL | Waktu dibatalkan |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu diupdate |

```sql
CREATE TABLE orders (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number          VARCHAR(50) NOT NULL UNIQUE,
    user_id               UUID NOT NULL REFERENCES users(id),
    worker_id             UUID NOT NULL REFERENCES users(id),
    service_id            UUID NOT NULL REFERENCES services(id),
    category_id           UUID NOT NULL REFERENCES categories(id),
    title                 VARCHAR(255) NOT NULL,
    description           TEXT NOT NULL,
    status                order_status NOT NULL DEFAULT 'pending',
    urgency               order_urgency NOT NULL DEFAULT 'normal',
    location_address      TEXT NOT NULL,
    location_detail       VARCHAR(500),
    location_lat          DECIMAL(10, 7) NOT NULL,
    location_lng          DECIMAL(10, 7) NOT NULL,
    preferred_date        DATE,
    preferred_time_start  TIME,
    preferred_time_end    TIME,
    notes                 TEXT,
    booking_fee           INTEGER NOT NULL DEFAULT 2000,
    base_service_fee      INTEGER,
    total_material_cost   INTEGER NOT NULL DEFAULT 0,
    total_additional_cost INTEGER NOT NULL DEFAULT 0,
    grand_total           INTEGER,
    cancellation_reason   TEXT,
    cancellation_category VARCHAR(50),
    cancelled_by          UUID REFERENCES users(id),
    accepted_at           TIMESTAMPTZ,
    started_at            TIMESTAMPTZ,
    completed_at          TIMESTAMPTZ,
    cancelled_at          TIMESTAMPTZ,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_orders_worker ON orders(worker_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_service ON orders(service_id);
CREATE INDEX idx_orders_category ON orders(category_id);
CREATE INDEX idx_orders_created ON orders(created_at DESC);
CREATE INDEX idx_orders_number ON orders(order_number);
CREATE INDEX idx_orders_user_status ON orders(user_id, status);
CREATE INDEX idx_orders_worker_status ON orders(worker_id, status);
```

---

### 5.8 `order_photos`

> Foto yang dilampirkan user saat membuat order.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `order_id` | UUID | NO | — | FK → orders.id |
| `photo_url` | TEXT | NO | — | URL foto |
| `caption` | VARCHAR(500) | YES | NULL | Keterangan foto |
| `display_order` | INTEGER | NO | `0` | Urutan tampilan |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |

```sql
CREATE TABLE order_photos (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id        UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    photo_url       TEXT NOT NULL,
    caption         VARCHAR(500),
    display_order   INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_order_photos_order ON order_photos(order_id);
```


---

### 5.9 `order_timeline`

> Timeline/progress pekerjaan. Setiap perubahan status dicatat di sini untuk tracking halaman user.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `order_id` | UUID | NO | — | FK → orders.id |
| `event` | VARCHAR(50) | NO | — | Kode event (e.g. order_created, worker_arrived) |
| `label` | VARCHAR(255) | NO | — | Label tampilan (e.g. "Pesanan Dibuat") |
| `description` | TEXT | YES | NULL | Deskripsi tambahan |
| `actor_id` | UUID | YES | NULL | FK → users.id (siapa yang trigger) |
| `actor_type` | VARCHAR(20) | YES | NULL | 'user', 'worker', 'system' |
| `metadata` | JSONB | YES | NULL | Data tambahan (latitude, longitude, notes, dll) |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu event terjadi |

```sql
CREATE TABLE order_timeline (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id        UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    event           VARCHAR(50) NOT NULL,
    label           VARCHAR(255) NOT NULL,
    description     TEXT,
    actor_id        UUID REFERENCES users(id),
    actor_type      VARCHAR(20),
    metadata        JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_order_timeline_order ON order_timeline(order_id);
CREATE INDEX idx_order_timeline_event ON order_timeline(event);
CREATE INDEX idx_order_timeline_created ON order_timeline(order_id, created_at);
```

---

### 5.10 `purchases`

> Item pembelian material/alat/biaya tambahan oleh worker. Inti dari fitur AI-assisted purchase tracking.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `order_id` | UUID | NO | — | FK → orders.id |
| `worker_id` | UUID | NO | — | FK → users.id (worker yang memasukkan) |
| `item_name` | VARCHAR(255) | NO | — | Nama item/barang |
| `category` | purchase_category | NO | `'material'` | Kategori pembelian |
| `quantity` | DECIMAL(10,2) | NO | `1` | Jumlah |
| `unit` | VARCHAR(50) | NO | `'pcs'` | Satuan (pcs, meter, sak, liter, dll) |
| `unit_price` | INTEGER | NO | `0` | Harga satuan (Rupiah) |
| `total_price` | INTEGER | NO | `0` | Harga total (Rupiah) |
| `reason` | TEXT | YES | NULL | Alasan pembelian |
| `receipt_photo_url` | TEXT | YES | NULL | URL foto nota/struk |
| `status` | purchase_status | NO | `'draft'` | Status approval |
| `confidence` | DECIMAL(3,2) | YES | NULL | AI confidence score (0.00 – 1.00) |
| `needs_clarification` | BOOLEAN | NO | `FALSE` | Perlu klarifikasi dari user? |
| `clarification_question` | TEXT | YES | NULL | Pertanyaan klarifikasi dari AI/user |
| `clarification_response` | TEXT | YES | NULL | Respons klarifikasi dari worker |
| `ai_explanation` | TEXT | YES | NULL | Penjelasan AI tentang pembelian ini |
| `raw_input` | TEXT | YES | NULL | Input mentah dari worker (sebelum AI proses) |
| `ai_processed_at` | TIMESTAMPTZ | YES | NULL | Waktu diproses AI |
| `approved_by` | UUID | YES | NULL | FK → users.id (user yang approve) |
| `approved_at` | TIMESTAMPTZ | YES | NULL | Waktu diapprove |
| `rejected_by` | UUID | YES | NULL | FK → users.id (user yang reject) |
| `rejected_at` | TIMESTAMPTZ | YES | NULL | Waktu ditolak |
| `rejection_reason` | TEXT | YES | NULL | Alasan penolakan |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu diupdate |

```sql
CREATE TABLE purchases (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id                UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    worker_id               UUID NOT NULL REFERENCES users(id),
    item_name               VARCHAR(255) NOT NULL,
    category                purchase_category NOT NULL DEFAULT 'material',
    quantity                DECIMAL(10, 2) NOT NULL DEFAULT 1,
    unit                    VARCHAR(50) NOT NULL DEFAULT 'pcs',
    unit_price              INTEGER NOT NULL DEFAULT 0,
    total_price             INTEGER NOT NULL DEFAULT 0,
    reason                  TEXT,
    receipt_photo_url       TEXT,
    status                  purchase_status NOT NULL DEFAULT 'draft',
    confidence              DECIMAL(3, 2),
    needs_clarification     BOOLEAN NOT NULL DEFAULT FALSE,
    clarification_question  TEXT,
    clarification_response  TEXT,
    ai_explanation          TEXT,
    raw_input               TEXT,
    ai_processed_at         TIMESTAMPTZ,
    approved_by             UUID REFERENCES users(id),
    approved_at             TIMESTAMPTZ,
    rejected_by             UUID REFERENCES users(id),
    rejected_at             TIMESTAMPTZ,
    rejection_reason        TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_purchases_order ON purchases(order_id);
CREATE INDEX idx_purchases_worker ON purchases(worker_id);
CREATE INDEX idx_purchases_status ON purchases(status);
CREATE INDEX idx_purchases_order_status ON purchases(order_id, status);
CREATE INDEX idx_purchases_created ON purchases(created_at DESC);
```

---

### 5.11 `purchase_risk_flags`

> Risk flags yang dihasilkan AI untuk setiap item pembelian.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `purchase_id` | UUID | NO | — | FK → purchases.id |
| `type` | risk_flag_type | NO | — | Jenis risk flag |
| `message` | TEXT | NO | — | Pesan/penjelasan risk flag |
| `is_resolved` | BOOLEAN | NO | `FALSE` | Sudah ditindaklanjuti? |
| `resolved_by` | UUID | YES | NULL | FK → users.id |
| `resolved_at` | TIMESTAMPTZ | YES | NULL | Waktu resolved |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |

```sql
CREATE TABLE purchase_risk_flags (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    purchase_id     UUID NOT NULL REFERENCES purchases(id) ON DELETE CASCADE,
    type            risk_flag_type NOT NULL,
    message         TEXT NOT NULL,
    is_resolved     BOOLEAN NOT NULL DEFAULT FALSE,
    resolved_by     UUID REFERENCES users(id),
    resolved_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_purchase_risk_flags_purchase ON purchase_risk_flags(purchase_id);
CREATE INDEX idx_purchase_risk_flags_type ON purchase_risk_flags(type);
CREATE INDEX idx_purchase_risk_flags_unresolved ON purchase_risk_flags(is_resolved) WHERE is_resolved = FALSE;
```

---

### 5.12 `purchase_audit_logs`

> Audit trail untuk setiap aksi pada item pembelian. Menjamin transparansi proses approval.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `purchase_id` | UUID | NO | — | FK → purchases.id |
| `action` | audit_action | NO | — | Jenis aksi |
| `actor_id` | UUID | YES | NULL | FK → users.id (pelaku aksi) |
| `actor_name` | VARCHAR(255) | YES | NULL | Nama pelaku (untuk display) |
| `actor_type` | VARCHAR(20) | NO | — | 'user', 'worker', 'system', 'ai' |
| `note` | TEXT | YES | NULL | Catatan tambahan |
| `old_data` | JSONB | YES | NULL | Data sebelum perubahan |
| `new_data` | JSONB | YES | NULL | Data setelah perubahan |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu aksi |

```sql
CREATE TABLE purchase_audit_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    purchase_id     UUID NOT NULL REFERENCES purchases(id) ON DELETE CASCADE,
    action          audit_action NOT NULL,
    actor_id        UUID REFERENCES users(id),
    actor_name      VARCHAR(255),
    actor_type      VARCHAR(20) NOT NULL,
    note            TEXT,
    old_data        JSONB,
    new_data        JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_purchase_audit_purchase ON purchase_audit_logs(purchase_id);
CREATE INDEX idx_purchase_audit_action ON purchase_audit_logs(action);
CREATE INDEX idx_purchase_audit_actor ON purchase_audit_logs(actor_id);
CREATE INDEX idx_purchase_audit_created ON purchase_audit_logs(created_at DESC);
```

---

### 5.13 `chat_messages`

> Pesan chat antara user dan worker per order.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `order_id` | UUID | NO | — | FK → orders.id |
| `sender_id` | UUID | NO | — | FK → users.id (pengirim) |
| `sender_type` | VARCHAR(20) | NO | — | 'user', 'worker', 'system' |
| `content` | TEXT | YES | NULL | Isi pesan (nullable jika image-only) |
| `message_type` | message_type | NO | `'text'` | Tipe pesan |
| `media_url` | TEXT | YES | NULL | URL media (foto, dll) |
| `is_read` | BOOLEAN | NO | `FALSE` | Sudah dibaca penerima? |
| `read_at` | TIMESTAMPTZ | YES | NULL | Waktu dibaca |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dikirim |

```sql
CREATE TABLE chat_messages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id        UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    sender_id       UUID NOT NULL REFERENCES users(id),
    sender_type     VARCHAR(20) NOT NULL,
    content         TEXT,
    message_type    message_type NOT NULL DEFAULT 'text',
    media_url       TEXT,
    is_read         BOOLEAN NOT NULL DEFAULT FALSE,
    read_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_chat_messages_order ON chat_messages(order_id);
CREATE INDEX idx_chat_messages_sender ON chat_messages(sender_id);
CREATE INDEX idx_chat_messages_order_created ON chat_messages(order_id, created_at DESC);
CREATE INDEX idx_chat_messages_unread ON chat_messages(order_id, is_read) WHERE is_read = FALSE;
```

---

### 5.14 `reviews`

> Rating dan review yang diberikan user setelah pekerjaan selesai.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `order_id` | UUID | NO | — | FK → orders.id (UNIQUE — 1 review per order) |
| `user_id` | UUID | NO | — | FK → users.id (reviewer) |
| `worker_id` | UUID | NO | — | FK → users.id (yang di-review) |
| `rating` | SMALLINT | NO | — | Rating 1–5 |
| `comment` | TEXT | YES | NULL | Ulasan teks |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu diupdate |

```sql
CREATE TABLE reviews (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id        UUID NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id),
    worker_id       UUID NOT NULL REFERENCES users(id),
    rating          SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reviews_order ON reviews(order_id);
CREATE INDEX idx_reviews_user ON reviews(user_id);
CREATE INDEX idx_reviews_worker ON reviews(worker_id);
CREATE INDEX idx_reviews_worker_rating ON reviews(worker_id, rating);
CREATE INDEX idx_reviews_created ON reviews(created_at DESC);
```

---

### 5.15 `review_tags`

> Tag yang dipilih user pada review (e.g. "cepat", "rapi", "profesional").

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `review_id` | UUID | NO | — | FK → reviews.id |
| `tag` | VARCHAR(50) | NO | — | Nama tag |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |

```sql
CREATE TABLE review_tags (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    review_id       UUID NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    tag             VARCHAR(50) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_review_tags_review ON review_tags(review_id);
CREATE INDEX idx_review_tags_tag ON review_tags(tag);
```

---

### 5.16 `invoices`

> Invoice/laporan akhir untuk setiap order yang selesai.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `order_id` | UUID | NO | — | FK → orders.id (UNIQUE) |
| `invoice_number` | VARCHAR(50) | NO | — | Nomor invoice (unique) |
| `base_service_fee` | INTEGER | NO | `0` | Biaya jasa dasar |
| `total_material_cost` | INTEGER | NO | `0` | Total biaya material (approved) |
| `total_additional_cost` | INTEGER | NO | `0` | Biaya tambahan |
| `booking_fee` | INTEGER | NO | `2000` | Booking fee |
| `platform_fee` | INTEGER | NO | `0` | Biaya layanan platform |
| `discount_amount` | INTEGER | NO | `0` | Jumlah diskon |
| `promo_code` | VARCHAR(50) | YES | NULL | Kode promo yang digunakan |
| `grand_total` | INTEGER | NO | `0` | Grand total pembayaran |
| `currency` | VARCHAR(3) | NO | `'IDR'` | Mata uang |
| `payment_instruction` | TEXT | YES | NULL | Instruksi pembayaran |
| `ai_work_summary` | TEXT | YES | NULL | Ringkasan pekerjaan oleh AI |
| `ai_materials_summary` | TEXT | YES | NULL | Ringkasan material oleh AI |
| `worker_notes` | TEXT | YES | NULL | Catatan dari worker |
| `all_purchases_approved` | BOOLEAN | NO | `TRUE` | Semua pembelian sudah diapprove? |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu diupdate |

```sql
CREATE TABLE invoices (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id                UUID NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
    invoice_number          VARCHAR(50) NOT NULL UNIQUE,
    base_service_fee        INTEGER NOT NULL DEFAULT 0,
    total_material_cost     INTEGER NOT NULL DEFAULT 0,
    total_additional_cost   INTEGER NOT NULL DEFAULT 0,
    booking_fee             INTEGER NOT NULL DEFAULT 2000,
    platform_fee            INTEGER NOT NULL DEFAULT 0,
    discount_amount         INTEGER NOT NULL DEFAULT 0,
    promo_code              VARCHAR(50),
    grand_total             INTEGER NOT NULL DEFAULT 0,
    currency                VARCHAR(3) NOT NULL DEFAULT 'IDR',
    payment_instruction     TEXT,
    ai_work_summary         TEXT,
    ai_materials_summary    TEXT,
    worker_notes            TEXT,
    all_purchases_approved  BOOLEAN NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_invoices_order ON invoices(order_id);
CREATE INDEX idx_invoices_number ON invoices(invoice_number);
```


---

### 5.17 `invoice_line_items`

> Rincian baris-baris pada invoice (breakdown biaya).

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `invoice_id` | UUID | NO | — | FK → invoices.id |
| `label` | VARCHAR(255) | NO | — | Label item (e.g. "Pipa PVC 2 meter") |
| `description` | TEXT | YES | NULL | Deskripsi tambahan |
| `category` | VARCHAR(50) | NO | — | Kategori: 'service', 'material', 'additional', 'fee', 'discount' |
| `quantity` | DECIMAL(10,2) | YES | `1` | Jumlah |
| `unit` | VARCHAR(50) | YES | NULL | Satuan |
| `unit_price` | INTEGER | YES | NULL | Harga satuan |
| `amount` | INTEGER | NO | `0` | Total amount untuk line item ini (Rupiah) |
| `purchase_id` | UUID | YES | NULL | FK → purchases.id (jika dari purchase) |
| `display_order` | INTEGER | NO | `0` | Urutan tampilan |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |

```sql
CREATE TABLE invoice_line_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id      UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    label           VARCHAR(255) NOT NULL,
    description     TEXT,
    category        VARCHAR(50) NOT NULL,
    quantity        DECIMAL(10, 2) DEFAULT 1,
    unit            VARCHAR(50),
    unit_price      INTEGER,
    amount          INTEGER NOT NULL DEFAULT 0,
    purchase_id     UUID REFERENCES purchases(id),
    display_order   INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_invoice_line_items_invoice ON invoice_line_items(invoice_id);
CREATE INDEX idx_invoice_line_items_purchase ON invoice_line_items(purchase_id);
```

---

### 5.18 `payments`

> Catatan pembayaran untuk setiap order.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `order_id` | UUID | NO | — | FK → orders.id |
| `invoice_id` | UUID | NO | — | FK → invoices.id |
| `user_id` | UUID | NO | — | FK → users.id (pembayar) |
| `amount` | INTEGER | NO | — | Jumlah pembayaran (Rupiah) |
| `currency` | VARCHAR(3) | NO | `'IDR'` | Mata uang |
| `payment_method` | payment_method | NO | — | Metode pembayaran |
| `payment_status` | payment_status | NO | `'unpaid'` | Status pembayaran |
| `payment_proof_url` | TEXT | YES | NULL | URL bukti pembayaran |
| `transaction_ref` | VARCHAR(255) | YES | NULL | Referensi transaksi (external) |
| `paid_at` | TIMESTAMPTZ | YES | NULL | Waktu bayar |
| `refunded_at` | TIMESTAMPTZ | YES | NULL | Waktu refund |
| `refund_amount` | INTEGER | YES | NULL | Jumlah refund |
| `refund_reason` | TEXT | YES | NULL | Alasan refund |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu diupdate |

```sql
CREATE TABLE payments (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id            UUID NOT NULL REFERENCES orders(id),
    invoice_id          UUID NOT NULL REFERENCES invoices(id),
    user_id             UUID NOT NULL REFERENCES users(id),
    amount              INTEGER NOT NULL,
    currency            VARCHAR(3) NOT NULL DEFAULT 'IDR',
    payment_method      payment_method NOT NULL,
    payment_status      payment_status NOT NULL DEFAULT 'unpaid',
    payment_proof_url   TEXT,
    transaction_ref     VARCHAR(255),
    paid_at             TIMESTAMPTZ,
    refunded_at         TIMESTAMPTZ,
    refund_amount       INTEGER,
    refund_reason       TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_order ON payments(order_id);
CREATE INDEX idx_payments_invoice ON payments(invoice_id);
CREATE INDEX idx_payments_user ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(payment_status);
CREATE INDEX idx_payments_method ON payments(payment_method);
```

---

### 5.19 `notifications`

> Notifikasi push/in-app untuk user.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `user_id` | UUID | NO | — | FK → users.id (penerima) |
| `type` | notification_type | NO | — | Tipe notifikasi |
| `title` | VARCHAR(255) | NO | — | Judul notifikasi |
| `body` | TEXT | NO | — | Isi notifikasi |
| `deep_link` | VARCHAR(500) | YES | NULL | Deep link ke halaman terkait |
| `metadata` | JSONB | YES | NULL | Data tambahan (order_id, purchase_id, dll) |
| `is_read` | BOOLEAN | NO | `FALSE` | Sudah dibaca? |
| `read_at` | TIMESTAMPTZ | YES | NULL | Waktu dibaca |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |

```sql
CREATE TABLE notifications (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type            notification_type NOT NULL,
    title           VARCHAR(255) NOT NULL,
    body            TEXT NOT NULL,
    deep_link       VARCHAR(500),
    metadata        JSONB,
    is_read         BOOLEAN NOT NULL DEFAULT FALSE,
    read_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, is_read) WHERE is_read = FALSE;
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_created ON notifications(user_id, created_at DESC);
```

---

### 5.20 `articles`

> Artikel knowledge base: tips, panduan, edukasi.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `title` | VARCHAR(500) | NO | — | Judul artikel |
| `slug` | VARCHAR(500) | NO | — | Slug URL-friendly (unique) |
| `category` | article_category | NO | — | Kategori artikel |
| `thumbnail_url` | TEXT | YES | NULL | URL thumbnail |
| `excerpt` | TEXT | YES | NULL | Ringkasan singkat |
| `content_html` | TEXT | NO | — | Konten HTML |
| `read_time_minutes` | INTEGER | YES | NULL | Estimasi waktu baca (menit) |
| `author` | VARCHAR(255) | YES | `'Tim HandyDirect'` | Penulis |
| `tags` | JSONB | YES | `'[]'` | Array tag |
| `is_published` | BOOLEAN | NO | `FALSE` | Sudah dipublikasi? |
| `published_at` | TIMESTAMPTZ | YES | NULL | Waktu publikasi |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu diupdate |

```sql
CREATE TABLE articles (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title               VARCHAR(500) NOT NULL,
    slug                VARCHAR(500) NOT NULL UNIQUE,
    category            article_category NOT NULL,
    thumbnail_url       TEXT,
    excerpt             TEXT,
    content_html        TEXT NOT NULL,
    read_time_minutes   INTEGER,
    author              VARCHAR(255) DEFAULT 'Tim HandyDirect',
    tags                JSONB DEFAULT '[]',
    is_published        BOOLEAN NOT NULL DEFAULT FALSE,
    published_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_articles_slug ON articles(slug);
CREATE INDEX idx_articles_category ON articles(category);
CREATE INDEX idx_articles_published ON articles(is_published, published_at DESC);
CREATE INDEX idx_articles_tags ON articles USING GIN(tags);
```

---

### 5.21 `faqs`

> Frequently Asked Questions.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `question` | TEXT | NO | — | Pertanyaan |
| `answer` | TEXT | NO | — | Jawaban |
| `category` | faq_category | NO | — | Kategori FAQ |
| `display_order` | INTEGER | NO | `0` | Urutan tampilan |
| `is_active` | BOOLEAN | NO | `TRUE` | Status aktif |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu diupdate |

```sql
CREATE TABLE faqs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question        TEXT NOT NULL,
    answer          TEXT NOT NULL,
    category        faq_category NOT NULL,
    display_order   INTEGER NOT NULL DEFAULT 0,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_faqs_category ON faqs(category);
CREATE INDEX idx_faqs_active ON faqs(is_active, display_order);
```

---

### 5.22 `promotions`

> Promo banner yang ditampilkan di Home user.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `title` | VARCHAR(255) | NO | — | Judul promo |
| `description` | TEXT | YES | NULL | Deskripsi promo |
| `image_url` | TEXT | NO | — | URL gambar promo |
| `cta_label` | VARCHAR(100) | YES | NULL | Label CTA button |
| `deep_link` | VARCHAR(500) | YES | NULL | Deep link ke halaman terkait |
| `promo_code` | VARCHAR(50) | YES | NULL | Kode promo (jika ada) |
| `discount_percent` | DECIMAL(5,2) | YES | NULL | Persentase diskon |
| `discount_amount` | INTEGER | YES | NULL | Jumlah diskon tetap (Rupiah) |
| `min_order_amount` | INTEGER | YES | NULL | Minimum order untuk promo |
| `display_order` | INTEGER | NO | `0` | Urutan tampilan |
| `is_active` | BOOLEAN | NO | `TRUE` | Status aktif |
| `valid_from` | TIMESTAMPTZ | YES | NULL | Mulai berlaku |
| `valid_until` | TIMESTAMPTZ | YES | NULL | Berakhir |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu diupdate |

```sql
CREATE TABLE promotions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title               VARCHAR(255) NOT NULL,
    description         TEXT,
    image_url           TEXT NOT NULL,
    cta_label           VARCHAR(100),
    deep_link           VARCHAR(500),
    promo_code          VARCHAR(50),
    discount_percent    DECIMAL(5, 2),
    discount_amount     INTEGER,
    min_order_amount    INTEGER,
    display_order       INTEGER NOT NULL DEFAULT 0,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    valid_from          TIMESTAMPTZ,
    valid_until         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_promotions_active ON promotions(is_active, display_order);
CREATE INDEX idx_promotions_valid ON promotions(valid_from, valid_until);
CREATE INDEX idx_promotions_code ON promotions(promo_code) WHERE promo_code IS NOT NULL;
```

---

### 5.23 `worker_wallets`

> Wallet/dompet digital worker. Satu worker memiliki satu wallet.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `worker_id` | UUID | NO | — | FK → users.id (UNIQUE) |
| `balance` | BIGINT | NO | `0` | Saldo saat ini (Rupiah) |
| `total_earnings` | BIGINT | NO | `0` | Total pendapatan keseluruhan |
| `total_withdrawn` | BIGINT | NO | `0` | Total yang sudah ditarik |
| `pending_earnings` | BIGINT | NO | `0` | Pendapatan yang belum masuk |
| `is_active` | BOOLEAN | NO | `TRUE` | Status aktif |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu diupdate |

```sql
CREATE TABLE worker_wallets (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    worker_id           UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    balance             BIGINT NOT NULL DEFAULT 0,
    total_earnings      BIGINT NOT NULL DEFAULT 0,
    total_withdrawn     BIGINT NOT NULL DEFAULT 0,
    pending_earnings    BIGINT NOT NULL DEFAULT 0,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_worker_wallets_worker ON worker_wallets(worker_id);
```

---

### 5.24 `wallet_transactions`

> Riwayat transaksi wallet worker.

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | UUID | NO | `gen_random_uuid()` | Primary key |
| `wallet_id` | UUID | NO | — | FK → worker_wallets.id |
| `order_id` | UUID | YES | NULL | FK → orders.id (jika dari order) |
| `type` | wallet_tx_type | NO | — | Tipe transaksi |
| `amount` | INTEGER | NO | — | Jumlah (Rupiah, positif) |
| `balance_before` | BIGINT | NO | — | Saldo sebelum transaksi |
| `balance_after` | BIGINT | NO | — | Saldo setelah transaksi |
| `description` | TEXT | YES | NULL | Deskripsi transaksi |
| `reference_id` | VARCHAR(255) | YES | NULL | ID referensi eksternal |
| `status` | wallet_tx_status | NO | `'pending'` | Status transaksi |
| `completed_at` | TIMESTAMPTZ | YES | NULL | Waktu selesai |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu dibuat |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | Waktu diupdate |

```sql
CREATE TABLE wallet_transactions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id           UUID NOT NULL REFERENCES worker_wallets(id) ON DELETE CASCADE,
    order_id            UUID REFERENCES orders(id),
    type                wallet_tx_type NOT NULL,
    amount              INTEGER NOT NULL,
    balance_before      BIGINT NOT NULL,
    balance_after       BIGINT NOT NULL,
    description         TEXT,
    reference_id        VARCHAR(255),
    status              wallet_tx_status NOT NULL DEFAULT 'pending',
    completed_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_wallet_tx_wallet ON wallet_transactions(wallet_id);
CREATE INDEX idx_wallet_tx_order ON wallet_transactions(order_id);
CREATE INDEX idx_wallet_tx_type ON wallet_transactions(type);
CREATE INDEX idx_wallet_tx_status ON wallet_transactions(status);
CREATE INDEX idx_wallet_tx_created ON wallet_transactions(wallet_id, created_at DESC);
```

---

## 6. Migration Order

Urutan migrasi database harus mengikuti dependency graph (tabel yang di-referensi harus dibuat lebih dulu):

```
Phase 1 — Enum Types (no dependencies)
  └── All CREATE TYPE statements

Phase 2 — Independent Tables (no FK dependencies)
  ├── users
  ├── categories
  ├── articles
  ├── faqs
  └── promotions

Phase 3 — First-level Dependencies
  ├── refresh_tokens        (depends on: users)
  ├── worker_profiles       (depends on: users)
  ├── services              (depends on: categories)
  ├── worker_wallets        (depends on: users)
  └── notifications         (depends on: users)

Phase 4 — Second-level Dependencies
  ├── worker_services       (depends on: users, services)
  └── orders                (depends on: users, services, categories)

Phase 5 — Order-dependent Tables
  ├── order_photos          (depends on: orders)
  ├── order_timeline        (depends on: orders, users)
  ├── purchases             (depends on: orders, users)
  ├── chat_messages         (depends on: orders, users)
  ├── reviews               (depends on: orders, users)
  ├── invoices              (depends on: orders)
  └── wallet_transactions   (depends on: worker_wallets, orders)

Phase 6 — Third-level Dependencies
  ├── purchase_risk_flags   (depends on: purchases, users)
  ├── purchase_audit_logs   (depends on: purchases, users)
  ├── review_tags           (depends on: reviews)
  ├── invoice_line_items    (depends on: invoices, purchases)
  └── payments              (depends on: orders, invoices, users)
```

---

## 7. Indexes Summary

| Table | Total Indexes | Key Index Strategy |
|---|---|---|
| `users` | 4 | Role, active status, geospatial (GIST) |
| `refresh_tokens` | 3 | User lookup, token hash, expiry |
| `worker_profiles` | 4 | User, availability, rating, verification |
| `categories` | 2 | Slug, active + order |
| `services` | 3 | Category, slug, active |
| `worker_services` | 2 + UNIQUE | Worker, service |
| `orders` | 9 | User, worker, status, service, composite |
| `order_photos` | 1 | Order |
| `order_timeline` | 3 | Order, event, composite time |
| `purchases` | 5 | Order, worker, status, composite, time |
| `purchase_risk_flags` | 3 | Purchase, type, partial (unresolved) |
| `purchase_audit_logs` | 4 | Purchase, action, actor, time |
| `chat_messages` | 4 | Order, sender, composite time, partial (unread) |
| `reviews` | 5 | Order, user, worker, composite rating, time |
| `review_tags` | 2 | Review, tag |
| `invoices` | 2 | Order, invoice number |
| `invoice_line_items` | 2 | Invoice, purchase |
| `payments` | 5 | Order, invoice, user, status, method |
| `notifications` | 4 | User, partial (unread), type, composite time |
| `articles` | 4 | Slug, category, published, GIN (tags) |
| `faqs` | 2 | Category, active + order |
| `promotions` | 3 | Active + order, validity, partial (promo code) |
| `worker_wallets` | 1 + UNIQUE | Worker |
| `wallet_transactions` | 5 | Wallet, order, type, status, composite time |
| **TOTAL** | **~81 indexes** | |

---

## 8. Notes & Conventions

### 8.1 Naming Conventions

| Convention | Rule | Example |
|---|---|---|
| Table names | snake_case, plural | `worker_profiles`, `chat_messages` |
| Column names | snake_case | `full_name`, `created_at` |
| Primary keys | Always `id` (UUID) | `id UUID PRIMARY KEY` |
| Foreign keys | `{referenced_table_singular}_id` | `user_id`, `order_id` |
| Timestamps | Always TIMESTAMPTZ | `created_at`, `updated_at` |
| Boolean columns | Prefix `is_` or `has_` | `is_active`, `is_read`, `needs_clarification` |
| Enum types | snake_case | `user_role`, `order_status` |
| Indexes | `idx_{table}_{column(s)}` | `idx_orders_user_status` |

### 8.2 Monetary Values

- Semua nilai uang disimpan dalam **satuan terkecil (Rupiah, tanpa desimal)** sebagai `INTEGER` atau `BIGINT`.
- Contoh: Rp150.000 disimpan sebagai `150000`.
- Ini menghindari masalah floating-point arithmetic.

### 8.3 Soft Delete

- Hanya tabel `users` yang menggunakan soft delete (`deleted_at`).
- Tabel lain menggunakan `ON DELETE CASCADE` dari parent.

### 8.4 Timestamps

- Seluruh timestamp menggunakan `TIMESTAMPTZ` (timestamp with time zone).
- Semua waktu disimpan dalam UTC dan dikonversi di client/frontend.

### 8.5 JSONB Columns

- `metadata` (order_timeline, notifications) — untuk data fleksibel yang tidak perlu kolom dedicated.
- `certificate_urls` (worker_profiles) — array URL sertifikat.
- `tags` (articles) — array tag artikel.
- `old_data` / `new_data` (purchase_audit_logs) — snapshot data sebelum/sesudah perubahan.

### 8.6 Geospatial

- Menggunakan PostGIS extension untuk geographic queries.
- Index GIST pada `users` untuk pencarian nearby workers.
- Formula jarak menggunakan `ST_DistanceSphere()` atau `ST_DWithin()`.

```sql
-- Contoh query: mencari worker dalam radius 10km
SELECT u.*, wp.*, 
       ST_DistanceSphere(
           ST_MakePoint(u.longitude, u.latitude),
           ST_MakePoint(:user_lng, :user_lat)
       ) / 1000.0 AS distance_km
FROM users u
JOIN worker_profiles wp ON wp.user_id = u.id
WHERE u.role = 'worker'
  AND u.is_active = TRUE
  AND wp.is_available = TRUE
  AND ST_DWithin(
      ST_MakePoint(u.longitude, u.latitude)::geography,
      ST_MakePoint(:user_lng, :user_lat)::geography,
      10000  -- 10km in meters
  )
ORDER BY distance_km ASC;
```

### 8.7 Trigger: Auto-update `updated_at`

```sql
CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply to all tables with updated_at column
CREATE TRIGGER set_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at BEFORE UPDATE ON worker_profiles
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at BEFORE UPDATE ON services
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at BEFORE UPDATE ON purchases
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at BEFORE UPDATE ON reviews
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at BEFORE UPDATE ON invoices
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at BEFORE UPDATE ON articles
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at BEFORE UPDATE ON faqs
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at BEFORE UPDATE ON promotions
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at BEFORE UPDATE ON worker_wallets
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at BEFORE UPDATE ON wallet_transactions
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
```

### 8.8 Required PostgreSQL Extensions

```sql
-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";      -- For gen_random_uuid() fallback
CREATE EXTENSION IF NOT EXISTS "pgcrypto";        -- For gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS "postgis";         -- For geospatial queries
```

---

> **📌 Document Version:** 1.0.0 | **Last Updated:** 2026-05-30 | **Total Tables:** 24 | **Total Indexes:** ~81

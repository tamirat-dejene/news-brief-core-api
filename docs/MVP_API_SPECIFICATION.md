# NewsBrief MVP API Specification v1.0

## 2-Week MVP Implementation - Enhanced API Documentation

**Audience:** Frontend developers building the MVP  
**Scope:** Essential endpoints for 2-week implementation with core features  
**Architecture:** Monolithic Core API + Scraper Service

---

## 🌐 MVP Service Architecture (2-Service Model)

| Service                  | URL                         | Responsibility                                                   |
| ------------------------ | --------------------------- | ---------------------------------------------------------------- |
| **Core API** (Go Gin)    | `http://localhost:8080`     | User Auth, Feed, Stories, Summaries, Chat, Briefs, Notifications |
| **Scraper** (FastAPI)    | `http://localhost:8001`     | Content extraction from URLs (HTTP)                              |
| **Vector DB** (Pinecone) | `https://api.pinecone.io`   | Stores vector embeddings for semantic search                     |
| **MongoDB**              | `mongodb://localhost:27017` | Primary data store                                               |

### **Architectural Philosophy:**

- **Core API as Monolith**: Handles all business logic including auth, user management, content processing, AI chat, brief generation, and notifications.
- **Dedicated Scraper**: Separate service for web scraping to isolate Python dependencies and external web requests.
- **Direct HTTP Communication**: Core API calls Scraper synchronously over HTTP (no message queue for MVP).
- **Embedded AI**: Core API includes Gemini integration for summarization and chat responses.

---

## 🔐 Authentication & Security

### **JWT Token Format (MongoDB ObjectId)**

```json
{
  "sub": "507f1f77bcf86cd799439011",
  "iss": "newsbrief.et",
  "exp": 1693123200,
  "iat": 1693119600,
  "aud": ["mobile", "web"],
  "scope": ["read:stories", "write:preferences"],
  "email_verified": true,
  "token_id": "66c1234567890abcdef12345"
}
```

### **Unified Token System**

All authentication tokens managed in single MongoDB collection:

```javascript
// MongoDB tokens collection structure
{
  _id: ObjectId("66c1234567890abcdef12345"),
  user_id: ObjectId("507f1f77bcf86cd799439011"),
  token_type: "refresh_token", // "email_verify" | "password_reset"
  token_hash: "sha256_hash_value",
  expires_at: ISODate("2025-08-26T10:30:00Z"),
  device_info: { // only for refresh_token
    user_agent: "NewsBrief/1.0.0 (iOS 16.0; iPhone14,2)",
    ip_address: "192.168.1.100"
  },
  used_at: null, // only for one-time tokens
  created_at: ISODate("2025-08-25T10:30:00Z")
}
```

### **Rate Limiting (MVP)**

| Endpoint Pattern      | Limit  | Window | Scope    |
| --------------------- | ------ | ------ | -------- |
| `GET /v1/feed`        | 60 req | 1 min  | Per IP   |
| `GET /v1/search`      | 30 req | 1 min  | Per IP   |
| `POST /v1/chat/query` | 5 req  | 1 min  | Per User |
| `POST /v1/auth/login` | 5 req  | 1 min  | Per IP   |
| `PATCH /v1/me/*`      | 30 req | 1 min  | Per User |

---

## 📊 Standardized Error Response

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Field validation failed",
    "details": { "field": "email", "reason": "invalid_format" },
    "request_id": "req_550e8400-e29b-41d4-a716-446655440000",
    "timestamp": "2025-08-25T10:30:00Z"
  }
}
```

**Common Status Codes:**

- `200` - Success
- `201` - Created
- `400` - Validation Error
- `401` - Authentication Required
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `429` - Rate Limited
- `500` - Internal Error

---

## 🔑 Authentication Endpoints

### **POST /v1/auth/register**

Create new account with email/password.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "SecureP@ssw0rd123!",
  "name": "John Doe",
  "lang": "am"
}
```

**Success Response (201):**

```json
{
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "name": "John Doe",
    "email_verified": false,
    "preferred_lang": "am",
    "created_at": "2025-08-25T10:30:00Z",
    "preferences": {
      "lang": "am",
      "topics": [],
      "data_saver": true,
      "brief_type": "short",
      "notifications": {
        "daily_brief": { "am": true, "pm": false },
        "wifi_only": true
      }
    }
  },
  "verification_token_id": "66c1234567890abcdef1234a",
  "verification_sent": true
}
```

### **POST /v1/auth/login**

Authenticate with email/password.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "SecureP@ssw0rd123!"
}
```

**Success Response (200):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "66c1234567890abcdef12349",
  "expires_in": 3600,
  "token_type": "Bearer",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "name": "John Doe",
    "email_verified": true,
    "last_login": "2025-08-25T10:30:00Z",
    "preferences": {
      "lang": "am",
      "topics": ["economy", "agriculture", "politics"],
      "data_saver": true,
      "brief_type": "short",
      "notifications": {
        "daily_brief": { "am": true, "pm": false },
        "wifi_only": true
      }
    }
  }
}
```

### **POST /v1/auth/refresh**

Refresh access token.

**Request Body:**

```json
{
  "refresh_token": "66c1234567890abcdef12349"
}
```

**Success Response (200):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "token_type": "Bearer",
  "refresh_token": "66c1234567890abcdef1234b",
  "token_rotated": true
}
```

### **POST /v1/auth/verify-email**

Verify email using unified token system.

**Request Body:**

```json
{
  "token": "66c1234567890abcdef1234a"
}
```

**Success Response (200):**

```json
{
  "message": "Email verified successfully",
  "user_id": "507f1f77bcf86cd799439011",
  "email_verified": true
}
```

### **POST /v1/auth/forgot-password**

Request password reset.

**Request Body:**

```json
{
  "email": "user@example.com"
}
```

**Success Response (200):**

```json
{
  "message": "Password reset instructions sent",
  "reset_token_id": "66c1234567890abcdef1234c",
  "expires_in_minutes": 30
}
```

### **POST /v1/auth/reset-password**

Reset password using reset token.

**Request Body:**

```json
{
  "token": "66c1234567890abcdef1234c",
  "new_password": "NewSecureP@ssw0rd123!"
}
```

**Success Response (200):**

```json
{
  "message": "Password reset successfully",
  "user_id": "507f1f77bcf86cd799439011",
  "all_sessions_invalidated": true
}
```

---

## 📱 Core Content Endpoints

### **GET /v1/feed**

Get paginated story feed with enhanced filtering.

**Query Parameters:**
| Field | Type | Required | Description | Example |
|-------|------|----------|-------------|---------|
| `lang` | `"am" \| "en"` | No | Summary language preference | `am` |
| `topic` | `string` | No | Topic filter from `/v1/topics` | `economy` |
| `source` | `string` | No | News outlet filter | `addisstandard` |
| `brief_type` | `"short" \| "medium"` | No | Summary length preference | `short` |
| `since` | `string` | No | ISO8601 timestamp | `2025-08-25T00:00:00Z` |
| `limit` | `number` | No | Page size (1-50, default: 20) | `20` |
| `cursor` | `string` | No | Pagination cursor | `eyJ0aW1lc3RhbXAiOjE2OTMxMjMyMDB9` |

**Success Response (200):**

```json
{
  "items": [
    {
      "id": "507f1f77bcf86cd799439012",
      "title": "Ethiopia launches new agricultural initiative",
      "source": {
        "key": "addisstandard",
        "name": "Addis Standard",
        "logo_url": "https://cdn.newsbrief.et/sources/addisstandard.png"
      },
      "url": "https://addisstandard.com/news/ethiopia-agri-2025",
      "published_at": "2025-08-25T09:10:00Z",
      "summary_short": "Government announces $50M investment in rural farming targeting 100,000 farmers.",
      "summary_bullets": [
        "Government announces $50M investment in rural farming",
        "Program targets 100,000 smallholder farmers nationwide",
        "Focus on drought-resistant crop varieties and irrigation",
        "Initiative includes training programs and equipment subsidies",
        "Expected to increase agricultural productivity by 35%"
      ],
      "summary_lang": "am",
      "topic_tags": ["agriculture", "economy"],
      "topic_image": "https://cdn.newsbrief.et/topics/agriculture.jpg",
      "processing_status": "completed",
      "content_hash": "sha256:a1b2c3d4e5f6...",
      "reading_time_minutes": {
        "short": 1,
        "medium": 3
      },
      "engagement_score": 0.87
    }
  ],
  "next_cursor": "eyJfaWQiOiI2NmMxMjM0NTY3ODkwYWJjZGVmMTIzNDUifQ==",
  "total_available": 156,
  "server_time": "2025-08-25T10:30:00Z"
}
```

### **GET /v1/story/:id**

Get detailed story by MongoDB ObjectId.

**Path Parameters:**
| Field | Type | Description |
|-------|------|-------------|
| `id` | `string` | MongoDB ObjectId (24-char hex string) |

**Success Response (200):**

```json
{
  "id": "507f1f77bcf86cd799439012",
  "title": "Ethiopia launches new agricultural initiative",
  "source": {
    "key": "addisstandard",
    "name": "Addis Standard",
    "logo_url": "https://cdn.newsbrief.et/sources/addisstandard.png"
  },
  "url": "https://addisstandard.com/news/ethiopia-agri-2025",
  "published_at": "2025-08-25T09:10:00Z",
  "summary_short": "Government announces $50M farming investment targeting 100,000 farmers.",
  "summary_bullets": [
    "Government announces $50M investment in rural farming",
    "Program targets 100,000 smallholder farmers nationwide",
    "Focus on drought-resistant crop varieties and irrigation"
  ],
  "summary_lang": "am",
  "topic_tags": ["agriculture", "economy"],
  "topic_image": "https://cdn.newsbrief.et/topics/agriculture.jpg",
  "processing_status": "completed",
  "content_hash": "sha256:a1b2c3d4e5f6...",
  "reading_time_minutes": 3,
  "word_count": 450,
  "scraped_at": "2025-08-25T09:05:00Z",
  "summarized_at": "2025-08-25T09:08:00Z"
}
```

### **GET /v1/search**

Enhanced search with MongoDB text index + Pinecone semantic search.

**Query Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `q` | `string` | Yes | Search query (min 2 chars) |
| `type` | `"text" \| "semantic"` | No | Search method (default: `text`) |
| `lang` | `"am" \| "en"` | No | Preferred result language |
| `topic` | `string` | No | Topic filter |
| `source` | `string` | No | Source filter |
| `since` | `string` | No | ISO8601 timestamp |
| `limit` | `number` | No | Results per page (1-50) |
| `cursor` | `string` | No | Pagination cursor |

**Success Response (200):**

```json
{
  "query": "agriculture investment",
  "items": [
    {
      "id": "507f1f77bcf86cd799439012",
      "title": "Ethiopia launches new agricultural initiative",
      "source": {
        "key": "addisstandard",
        "name": "Addis Standard",
        "logo_url": "https://cdn.newsbrief.et/sources/addisstandard.png"
      },
      "url": "https://addisstandard.com/news/ethiopia-agri-2025",
      "published_at": "2025-08-25T09:10:00Z",
      "summary_short": "Government announces $50M farming investment.",
      "summary_bullets": ["Government announces...", "..."],
      "summary_lang": "am",
      "topic_tags": ["agriculture", "economy"],
      "processing_status": "completed",
      "relevance_score": 0.95,
      "matched_terms": ["agriculture", "investment"],
      "text_score": 1.2,
      "highlight_snippets": [
        "Government announces $50M <mark>investment</mark> in rural <mark>agriculture</mark>",
        "New <mark>agricultural</mark> <mark>investment</mark> program launches nationwide"
      ]
    }
  ],
  "next_cursor": "eyJfaWQiOiI2NmMxMjM0NTY3ODkwYWJjZGVmMTIzNDYiLCJzY29yZSI6MC45NX0=",
  "total_matches": 12,
  "search_time_ms": 45,
  "search_method": "mongodb_text_index"
}
```

### **GET /v1/topics**

Get topic categories with localized labels and images.

**Success Response (200):**

```json
{
  "topics": [
    {
      "key": "economy",
      "label": { "en": "Economy", "am": "ኢኮኖሚ" },
      "description": {
        "en": "Business, finance, and economic news",
        "am": "የንግድ፣ የፋይናንስ እና የኢኮኖሚ ዜናዎች"
      },
      "image_url": "https://cdn.newsbrief.et/topics/economy.jpg",
      "story_count": 45,
      "subscribed_sources": 8,
      "last_updated": "2025-08-25T10:00:00Z"
    },
    {
      "key": "agriculture",
      "label": { "en": "Agriculture", "am": "ግብርና" },
      "description": {
        "en": "Farming, livestock, and agricultural development",
        "am": "የግብርና፣ የእንስሳት ሀብት እና የግብርና ልማት"
      },
      "image_url": "https://cdn.newsbrief.et/topics/agriculture.jpg",
      "story_count": 23,
      "subscribed_sources": 6,
      "last_updated": "2025-08-25T09:30:00Z"
    }
  ],
  "total_topics": 8,
  "last_updated": "2025-08-20T00:00:00Z"
}
```

### **GET /v1/sources**

Get available news outlets with subscription support.

**Query Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `search` | `string` | No | Search outlets by name |
| `country` | `string` | No | Filter by country code |
| `lang` | `"am" \| "en"` | No | Filter by language |
| `subscribed_only` | `boolean` | No | Show only subscribed outlets |

**Success Response (200):**

```json
{
  "sources": [
    {
      "key": "addisstandard",
      "name": "Addis Standard",
      "description": "Independent news outlet covering Ethiopian politics and society",
      "url": "https://addisstandard.com",
      "logo_url": "https://cdn.newsbrief.et/sources/addisstandard.png",
      "country": "ET",
      "languages": ["en", "am"],
      "topics": ["politics", "economy", "society"],
      "rss_feeds": [
        {
          "url": "https://addisstandard.com/feed/",
          "topic": "general",
          "active": true
        }
      ],
      "reliability_score": 0.92,
      "update_frequency": "hourly",
      "avg_articles_per_day": 15,
      "last_updated": "2025-08-25T10:15:00Z"
    }
  ],
  "total_sources": 25,
  "available_countries": ["ET", "KE", "UG"],
  "available_topics": [
    "politics",
    "economy",
    "agriculture",
    "society",
    "sports",
    "international"
  ]
}
```

---

## 🤖 AI Chat Endpoints (NEW)

### **POST /v1/chat/query** 🔐

AI-powered news query with web search and summarization.

**Request Body:**

```json
{
  "query": "What is the latest news about Ethiopia's renewable energy projects?",
  "lang": "am",
  "max_sources": 5,
  "brief_type": "medium"
}
```

**Success Response (200):**

```json
{
  "query_id": "507f1f77bcf86cd799439015",
  "query": "What is the latest news about Ethiopia's renewable energy projects?",
  "response": {
    "summary_bullets": [
      "Ethiopia launches 300MW solar power project in Afar region",
      "Government signs $2B renewable energy agreement with international partners",
      "New wind farm construction begins in Tigray with 200MW capacity",
      "Ethiopia targets 30GW renewable energy capacity by 2030",
      "Hydroelectric projects face environmental concerns from local communities"
    ],
    "summary_lang": "am",
    "confidence_score": 0.92,
    "sources_used": [
      {
        "url": "https://ethiopianherald.com/renewable-energy-2025",
        "title": "Ethiopia's Renewable Energy Expansion",
        "source": "Ethiopian Herald",
        "scraped_at": "2025-08-25T10:30:00Z",
        "relevance_score": 0.95
      }
    ],
    "vector_matches": 3,
    "new_scrapes": 2,
    "processing_time_ms": 2340
  },
  "created_at": "2025-08-25T10:30:00Z"
}
```

### **GET /v1/chat/history** 🔐

Get user's chat query history.

**Query Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `limit` | `number` | No | Results per page (1-50) |
| `cursor` | `string` | No | Pagination cursor |

**Success Response (200):**

```json
{
  "items": [
    {
      "query_id": "507f1f77bcf86cd799439015",
      "query": "What is the latest news about Ethiopia's renewable energy projects?",
      "created_at": "2025-08-25T10:30:00Z",
      "lang": "am",
      "sources_count": 5,
      "processing_time_ms": 2340
    }
  ],
  "next_cursor": "eyJfaWQiOiI2NmMxMjM0NTY3ODkwYWJjZGVmMTIzNDYifQ==",
  "total_queries": 12
}
```

### **GET /v1/chat/query/:id** 🔐

Get detailed results for a specific chat query.

**Path Parameters:**
| Field | Type | Description |
|-------|------|-------------|
| `id` | `string` | Query ID (MongoDB ObjectId) |

**Success Response (200):**
Same format as POST `/v1/chat/query` response.

---

## 📊 Daily Brief Endpoints (ENHANCED)

### **GET /v1/daily-brief**

Get curated morning/evening story collection.

**Query Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `slot` | `"am" \| "pm"` | Yes | Brief timing |
| `lang` | `"am" \| "en"` | No | Content language |
| `date` | `string` | No | ISO date (YYYY-MM-DD) |

**Success Response (200):**

```json
{
  "id": "66c1234567890abcdef12347",
  "slot": "am",
  "lang": "am",
  "date": "2025-08-25",
  "title": "የጠዋት ዜና ማጠቃለያ",
  "description": "Today's 5 most important stories for Ethiopia",
  "created_at": "2025-08-25T06:00:00Z",
  "story_count": 5,
  "estimated_read_time_minutes": 8,
  "curation_algorithm": "ai_scored_v1.2",
  "stories": [
    {
      "id": "507f1f77bcf86cd799439012",
      "title": "Ethiopia launches new agricultural initiative",
      "source": {
        "key": "addisstandard",
        "name": "Addis Standard",
        "logo_url": "https://cdn.newsbrief.et/sources/addisstandard.png"
      },
      "published_at": "2025-08-25T05:30:00Z",
      "summary_bullets": ["Government announces...", "..."],
      "summary_lang": "am",
      "topic_tags": ["agriculture"],
      "processing_status": "completed",
      "brief_position": 1,
      "curation_score": 0.95,
      "reading_time_minutes": 2
    }
  ],
  "content_hash": "sha256:f1g2h3i4j5k6...",
  "last_updated": "2025-08-25T06:00:00Z"
}
```

### **GET /v1/me/briefs** 🔐

Get personalized daily briefs (alias for backward compatibility).

**Query Parameters:**
Same as `/v1/daily-brief`

---

## 👤 User Profile & Preferences (ENHANCED)

### **GET /v1/me** 🔐

Get authenticated user profile and preferences.

**Success Response (200):**

```json
{
  "id": "507f1f77bcf86cd799439011",
  "email": "user@example.com",
  "name": "John Doe",
  "email_verified": true,
  "created_at": "2025-08-25T10:30:00Z",
  "last_lo-gin": "2025-08-25T10:30:00Z",
  "preferences": {
    "topics": ["economy", "agriculture", "politics"],
    "subscribed_sources": ["addisstandard", "ethiopianherald"],
    "send_notifications": false,
    "notifications": {
      "daily_brief": { "am": true, "pm": false }
    }
  },
  "stats": {
    "chat_queries": 8,
    "last_sync": "2025-08-25T09:15:00Z"
  }
}
```

### **PATCH /v1/me/password** 🔐

Change user password.

**Request Body:**

```json
{
  "current_password": "OldP@ssw0rd123!",
  "new_password": "NewSecureP@ssw0rd456!"
}
```

**Success Response (200):**

```json
{
  "message": "Password updated successfully",
  "all_sessions_invalidated": true,
  "new_access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "new_refresh_token": "66c1234567890abcdef1234d"
}
```

### **PATCH /v1/me/preferences** 🔐

Update user preferences (partial update).

**Request Body:**

```json
{
  "lang": "en",
  "topics": ["economy", "agriculture", "technology"],
  "brief_type": "medium",
  "data_saver": false,
  "notifications": {
    "daily_brief": { "am": true, "pm": true }
  }
}
```

**Success Response (200):**

```json
{
  "lang": "en",
  "topics": ["economy", "agriculture", "technology"],
  "brief_type": "medium",
  "data_saver": false,
  "notifications": {
    "daily_brief": { "am": true, "pm": true },
    "wifi_only": true
  },
  "updated_at": "2025-08-25T10:30:00Z"
}
```

### **GET /v1/me/subscriptions** 🔐

Get user's news outlet subscriptions.

**Success Response (200):**

```json
{
  "subscriptions": [
    {
      "source_key": "addisstandard",
      "source_name": "Addis Standard",
      "subscribed_at": "2025-08-01T10:00:00Z",
      "topics": ["politics", "economy"],
      "notification_enabled": true,
      "last_article": "2025-08-25T09:30:00Z"
    }
  ],
  "total_subscriptions": 3,
  "subscription_limit": 10,
  "premium_features": false
}
```

### **POST /v1/me/subscriptions** 🔐

Subscribe to a news outlet.

**Request Body:**

```json
{
  "source_key": "ethiopianherald",
  "topics": ["economy", "politics"],
  "notifications": true
}
```

**Success Response (201):**

```json
{
  "subscription": {
    "source_key": "ethiopianherald",
    "source_name": "Ethiopian Herald",
    "topics": ["economy", "politics"],
    "notification_enabled": true,
    "subscribed_at": "2025-08-25T10:30:00Z"
  },
  "total_subscriptions": 4
}
```

### **DELETE /v1/me/subscriptions/:source_key** 🔐

Unsubscribe from news outlet.

**Path Parameters:**
| Field | Type | Description |
|-------|------|-------------|
| `source_key` | `string` | News outlet identifier |

**Success Response (200):**

```json
{
  "message": "Successfully unsubscribed from Addis Standard",
  "source_key": "addisstandard",
  "unsubscribed_at": "2025-08-25T10:30:00Z",
  "remaining_subscriptions": 2
}
```

### **PATCH /v1/me/subscriptions/:source_key** 🔐

Update subscription preferences for a news outlet.

**Path Parameters:**
| Field | Type | Description |
|-------|------|-------------|
| `source_key` | `string` | News outlet identifier |

**Request Body:**

```json
{
  "topics": ["politics", "economy", "agriculture"],
  "notifications": false
}
```

**Success Response (200):**

```json
{
  "subscription": {
    "source_key": "addisstandard",
    "source_name": "Addis Standard",
    "subscribed_at": "2025-08-01T10:00:00Z",
    "topics": ["politics", "economy", "agriculture"],
    "notification_enabled": false,
    "updated_at": "2025-08-25T10:30:00Z"
  }
}
```

---

## 🔔 Notifications & Analytics (NEW)

### **GET /v1/me/notifications** 🔐

Get user's notifications.

**Query Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `limit` | `number` | No | Results per page (1-50) |
| `unread_only` | `boolean` | No | Show only unread notifications |

**Success Response (200):**

```json
{
  "notifications": [
    {
      "id": "notif_66c1234567890abcdef1234f",
      "type": "daily_brief_ready",
      "title": "Your Morning Brief is ready!",
      "message": "Catch up on the latest news for August 26.",
      "data": { "brief_id": "66c1234567890abcdef12347" },
      "is_read": false,
      "created_at": "2025-08-26T06:01:00Z"
    }
  ],
  "unread_count": 1,
  "total_notifications": 24
}
```

### **POST /v1/me/notifications/mark-read** 🔐

Mark notifications as read.

**Request Body:**

```json
{
  "notification_ids": ["notif_66c1234567890abcdef1234f"]
}
```

**Success Response (200):**

```json
{
  "message": "Notifications marked as read",
  "updated_count": 1,
  "remaining_unread": 0
}
```

### **GET /v1/me/analytics** 🔐

Get user reading analytics.

**Query Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `period` | `"week" \| "month"` | No | Analytics timeframe (default: `week`) |

**Success Response (200):**

```json
{
  "period": "month",
  "date_range": {
    "from": "2025-07-25T00:00:00Z",
    "to": "2025-08-25T00:00:00Z"
  },
  "reading_stats": {
    "articles_read": 89,
    "time_spent_minutes": 267,
    "daily_average": 8.9,
    "completion_rate": 0.76
  },
  "topic_breakdown": [
    { "topic": "economy", "articles": 34, "percentage": 38.2 },
    { "topic": "agriculture", "articles": 28, "percentage": 31.5 },
    { "topic": "politics", "articles": 27, "percentage": 30.3 }
  ],
  "source_breakdown": [
    { "source": "addisstandard", "articles": 45, "percentage": 50.6 },
    { "source": "ethiopianherald", "articles": 32, "percentage": 36.0 }
  ]
}
```

---

## 🔄 Sync & Health Endpoints (NEW)

### **GET /v1/sync/manifest**

Provide sync metadata for offline-first apps.

**Query Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `since` | `string` | No | Client's last sync timestamp |
| `lang` | `"am" \| "en"` | No | Preferred content language |

**Success Response (200):**

```json
{
  "server_time": "2025-08-25T10:30:00Z",
  "content_version": "v1.0.0",
  "stories": {
    "total_count": 234,
    "updated_since": "2025-08-25T09:00:00Z",
    "new_count": 12,
    "updated_count": 3,
    "content_hash": "sha256:a1b2c3d4...",
    "last_object_id": "66c1234567890abcdef12348"
  },
  "daily_briefs": {
    "am": {
      "available": true,
      "brief_id": "66c1234567890abcdef12347",
      "date": "2025-08-25",
      "updated": "2025-08-25T06:00:00Z"
    },
    "pm": {
      "available": false,
      "date": "2025-08-25",
      "next_expected": "2025-08-25T18:00:00Z"
    }
  },
  "topics": {
    "updated": "2025-08-20T00:00:00Z",
    "content_hash": "sha256:i9j0k1l2..."
  }
}
```

### **GET /health**

Basic service health check.

**Success Response (200):**

```json
{
  "status": "ok",
  "timestamp": "2025-08-25T10:30:00Z",
  "version": "1.0.0"
}
```

---

## 🔧 Internal Service Communication

### **Core API -> Scraper Communication**

**POST {SCRAPER_BASE_URL}/scrape**

**Request:**

```json
{
  "url": "https://addisstandard.com/news/ethiopia-agri-2025",
  "source_key": "addisstandard",
  "timeout_seconds": 30
}
```

**Response:**

```json
{
  "url": "https://addisstandard.com/news/ethiopia-agri-2025",
  "title": "Ethiopia launches new agricultural initiative",
  "text": "The Ethiopian government today announced...",
  "published_at": "2025-08-25T09:10:00Z",
  "author": "Reporter Name",
  "lang": "en",
  "word_count": 450,
  "scraped_at": "2025-08-25T10:30:00Z",
  "content_hash": "sha256:a1b2c3d4e5f6..."
}
```

### **Core API Processing Flow**

1. **Content Ingestion**: Receive scraped content from Scraper service
2. **Text Processing**: Clean and validate content
3. **AI Summarization**: Generate summaries using Gemini API
4. **Vector Embedding**: Create and store embeddings in Pinecone
5. **Data Persistence**: Store story in MongoDB
6. **Background Jobs**: Brief generation and notifications

---

## 🐳 MVP Deployment

### **Docker Compose (Local Development)**

```yaml
version: "3.8"
services:
  mongodb:
    image: mongo:7.0
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_DATABASE: newsbrief

  core-api:
    build: ./core-api
    ports:
      - "8080:8080"
    depends_on:
      - mongodb
    environment:
      MONGODB_URI: "mongodb://mongodb:27017/newsbrief"
      SCRAPER_BASE_URL: "http://scraper:8001"
      PINECONE_API_KEY: ${PINECONE_API_KEY}
      GEMINI_API_KEY: ${GEMINI_API_KEY}
      JWT_SECRET: ${JWT_SECRET}

  scraper:
    build: ./scraper
    ports:
      - "8001:8001"
    environment:
      LOG_LEVEL: info
      TIMEOUT_SECONDS: 30
```

### **Environment Variables**

```bash
# Core API
MONGODB_URI=mongodb://localhost:27017/newsbrief
SCRAPER_BASE_URL=http://localhost:8001
PINECONE_API_KEY=your_pinecone_api_key
GEMINI_API_KEY=your_gemini_api_key
JWT_SECRET=your_jwt_secret_key

# Scraper
LOG_LEVEL=info
TIMEOUT_SECONDS=30
```

---

## 📈 Implementation Priority

### **Phase 1 (Week 1)**

1. ✅ Authentication endpoints (`/v1/auth/*`)
2. ✅ Basic user management (`/v1/me`)
3. ✅ Story feed (`/v1/feed`)
4. ✅ Story details (`/v1/story/:id`)
5. ✅ Topics and sources (`/v1/topics`, `/v1/sources`)

### **Phase 2 (Week 2)**

1. ✅ Enhanced search (`/v1/search`)
2. ✅ AI chat functionality (`/v1/chat/*`)
3. ✅ Daily briefs (`/v1/daily-brief`)
4. ✅ User subscriptions (`/v1/me/subscriptions`)
5. ✅ Notifications (`/v1/me/notifications`)

### **Phase 3 (Post-MVP)**

1. Analytics dashboard
2. Advanced user preferences
3. Push notifications
4. Offline sync optimization

---

This enhanced MVP specification includes all core features from the professional API while maintaining the simple 2-service architecture suitable for rapid development and deployment.

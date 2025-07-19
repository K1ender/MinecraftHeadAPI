# Minecraft Skin Head Renderer

This service fetches and renders the **head image** of a Minecraft player from their Mojang skin, with optional Redis caching.

## 🧠 Features

- Fetches UUID by nickname via Mojang API
- Retrieves and decodes skin textures
- Crops and resizes the player’s head from the full skin
- Optional overlay support (second layer)
- Redis-based caching (24h) to reduce API usage
- Exposes simple HTTP API

## 🚀 Usage

### 1. Environment Variables

| Variable     | Description                                | Default         |
|--------------|--------------------------------------------|-----------------|
| `REDIS_URL`  | Redis connection string (e.g. `redis://localhost:6379/0`) | **Required**     |
| `REAL_IP`    | If not set to `false`, resolves actual IP for logging | `localhost`      |

### 2. Start Server

```bash
go run ./cmd/server
````

Server will start at: `http://<your-ip>:8080`

## 🖼️ API Endpoint

```
GET /head/{nickname}
```

* **Returns**: PNG image of the user's Minecraft head (base64-decoded PNG)
* **Example**:

  ```
  curl http://localhost:8080/head/Notch --output notch.png
  ```

## 🗂 Project Structure

```text
.
├── cmd/                   # Entry point
│   └── server/            # Main app logic
├── internal/
│   ├── api/               # HTTP handlers
│   ├── service/           # Head rendering and caching logic
│   ├── mojang/            # Mojang API client
│   └── util/              # Image and network utilities
```

## 💾 Redis Caching

* Redis is used to cache rendered head images using UUID as key.
* TTL: 24 hours
* If Redis is down or the cache misses, the image is fetched directly from Mojang.

## 🛠 Dependencies

* [Go](https://golang.org/) 1.20+
* [disintegration/imaging](https://github.com/disintegration/imaging) – image manipulation
* [go-redis/redis/v9](https://github.com/redis/go-redis) – Redis client
* [uuid](https://github.com/google/uuid) – UUID parsing

## 📄 License

MIT – free to use and modify.
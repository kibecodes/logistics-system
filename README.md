# FastaBiz 

FastaBiz helps businesses streamline **inventory, deliveries, and customer orders** — all in one platform. Whether you’re a farmer, shop owner, or everyday buyer, we provide a simple way to manage operations and connect with customers.

---

## ✨ Why FastaBiz?  

Many small and medium businesses struggle with:  

- ❌ Manual tracking of stock, orders, and deliveries  
- ❌ Poor visibility into customer feedback & delivery status  
- ❌ No central tool for both **business owners and buyers**  

**FastaBiz solves this with an integrated platform**:  
- 📦 Smart inventory management  
- 🚚 Delivery & driver coordination  
- 🛒 Customer ordering & tracking  
- 📊 Dashboards for insights 

---

## 🚀 Features

- **Role-based access control**: Admin, Driver, Customer  
- **Full CRUD APIs** for orders, deliveries, payments, feedback, notifications  
- **Business storefronts**: Businesses get unique links to share their inventory  
- **Customer tools**: Browse, order, track deliveries, review  
- **Dockerized microservices**: Backend, DB, Kong API Gateway  
- **Security**: JWT authentication, rate limiting via Kong  
- **CI/CD ready**: GitHub Actions with API tests 

---

## 🛠️ Tech Stack

| Layer       | Technologies                                 |
|-------------|----------------------------------------------|
| Frontend    | Blazor (C#), TailwindCSS                     |
| Backend     | Go (Chi, Clean Architecture, Swagger)        |
| Gateway     | Kong (JWT auth + rate limiting)              |
| Database    | PostgreSQL                                   |
| CI/CD       | GitHub Actions + Docker + Postman/Newman     |
| Containerization | Docker, Docker Compose                  |

---

## 📁 Repository Structure

```
logistics-system/
├── apps/
│   ├── logistics-backend/        # Backend (Go APIs, infra & configs)
│   │   ├── kong/                 # Kong declarative config
│   │   │   └── kong.yml
│   │   ├── postman/              # API test collections
│   │   │   ├── collection.json
│   │   │   └── environment.json
│   │   ├── .github/              # CI workflows
│   │   │   └── workflows/
│   │   │       └── api-tests.yml
│   │   ├── .env.docker           # Docker environment variables
│   │   ├── Dockerfile            # Backend Dockerfile
│   │   └── docker-compose.yml    # Compose services
│   └── logistics-frontend/       # Frontend (Blazor app)
├── .gitignore
├── README.md
└── LICENSE

```

---

## 🖼️ Project Flow (Business Use Case)  

1️⃣ **Business Owners**: Register → Upload inventory → Share store link  
2️⃣ **Customers**: Browse via link → Place orders → Track deliveries  
3️⃣ **Admins/Drivers**: Manage deliveries, drivers, and feedback  

*(Illustrations & screenshots will be added here — AI-generated concept images for now, real dashboard shots later.)*  

---

## ⚙️ Getting Started

### Prerequisites

- Docker & Docker Compose
- Git

---

### 🚀 Running Locally with Docker

```bash
git clone https://github.com/kibecodes/logistics-system.git
cd logistics-system

# Start all services: DB, backend, Kong
docker compose up --build
```

- **API & Swagger:** `http://localhost:8000/api/swagger/index.html`
- **Backend logs:** `docker compose logs -f backend`
- **Kong Admin:** `http://localhost:8001`

---

### 🧪 Testing APIs

```bash
docker run --rm \
  -v "${PWD}/postman:/etc/newman" \
  postman/newman:alpine run collection.json \
  --environment=environment.json --reporters cli

```

---

## 🧩 Environment Configuration

**.env.docker**

```env
PUBLIC_API_BASE_URL=http://localhost:8000/api
INTERNAL_API_BASE_URL=http://backend:8080
DATABASE_URL=postgres://admin:secret@db:5432/logistics_db?sslmode=disable
PORT=8080
JWT_SECRET=<your-secret>
```

Kong connects to the backend on `http://backend:8080` internally, while clients use `localhost:8000`.

---

## 📈 Roadmap

✅ Proof of Concept APIs

🚧 Business logic for orders, drivers, routes

🚧 Frontend dashboards for Admin, Driver, Customer

🚧 gRPC/Kafka integration for async flows

🚧 Production CI/CD & monitoring

---

## 🤝 Contributing

Your contributions are welcome! Suggested areas:

- Completing business logic and clean architecture layers
- Adding frontend user interfaces or dashboards
- Production-grade logging, monitoring, and gateway enhancements
- Message bus integrations (Kafka / RabbitMQ)

---

## 📝 License

MIT License – see [LICENSE](LICENSE)

---
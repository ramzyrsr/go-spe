# FinTech Payment API 🏦💳
This is a Golang-based FinTech Payment API that supports secure authentication, transaction notifications, and transaction status checks. It follows Domain-Driven Design (DDD) and implements Redis caching, RabbitMQ for asynchronous processing, and PostgreSQL for database management.

---

## 🚀 Features
- ✅ Secure Authentication using HMAC
- ✅ Asynchronous Transaction Processing via RabbitMQ
- ✅ Optimized Status Checks using Redis Caching

---

## 🛠️ Tech Stack
| **Component**         | **Technology**   |
|-----------------------|------------------|
| Language              | Golang           |
| Web Framework         | Gin              |
| Database              | PostgreSQL       |
| Cache                 | Redis            |
| Async Processing      | RabbitMQ         |

---

## 📦 Installation
1️⃣ **Clone the Repository:**
```sh
git clone git@github.com:ramzyrsr/go-spe.git
cd go-spe
```
2️⃣ **Install Dependencies**
```sh
go mod tidy
```
3️⃣ **Setup Environment Variables**

Create a .env file in the root directory and add:
```sh
POSTGRES_DSN=
SECRET_KEY=
REDIS_URL=
RABBITMQ_URL=
```
4️⃣ **Run the API**
```sh
go run cmd/main.go
```
---

## 📡 API Endpoints

### 💳 Transaction Notification
| **Method**  | **Endpoint**                     | **Description**                 |
|-------------|-----------------------------------|---------------------------------|
| POST        | /transaction-notification         | Notify transaction status      |

### 🔍 Check Transaction Status
| **Method**  | **Endpoint**                     | **Description**                          |
|-------------|-----------------------------------|------------------------------------------|
| POST        | /check-status                     | Check status based on `bill_number`      |


---

## 👨‍💻 Contributors
- Ramzy Syahrul Ramadhan - [GitHub](https://github.com/ramzyrsr)
- Feel free to contribute by opening issues & pull requests!
---
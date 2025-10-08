# 🧩 Fullstack Developer Test Challenge

## 🚀 How to Run the Stack Locally

Clone repository dan pastikan Docker sudah ter-install:

```bash
git clone https://github.com/username/fullstack-test.git
cd fullstack-test

Jalankan semua container:
docker-compose up -d --build

🔍 Default Ports
Service	URL / Host	Port
Product Service	http://localhost:3000	3000
Order Service	http://localhost:3001	3001

⚙️ Environment Variables
product-service/.env
DATABASE_URL="mysql://root:root@mysql:3306/fullstack-test?parseTime=true"
REDIS_HOST=redis
REDIS_PORT=6379
KAFKA_BROKER=kafka:9092

order-service/.env
DB_HOST=mysql
DB_PORT=3306
DB_NAME=fullstack-test
DB_USER=root
DB_PASSWORD=root

REDIS_HOST=redis
REDIS_PORT=6379

PRODUCT_SERVICE_URL=http://product-service:3000
# KAFKA_BROKER=kafka:9092

🏗️ Architecture Overview
Aplikasi ini dibangun dengan arsitektur microservices yang terdiri dari dua service utama:

Product Service (NestJS)

Order Service (Golang)

Keduanya berkomunikasi secara event-driven menggunakan Apache Kafka, dengan MySQL sebagai database dan Redis untuk caching.

🔹 Product Service (NestJS)
Database: MySQL (products table)

Cache: Redis

Endpoints:

POST /products → menambahkan produk baru

GET /products/:id → mengambil detail produk (menggunakan cache Redis)

Event Flow:

Mengirim event product.created ke Kafka saat produk dibuat

Mendengarkan event order.created untuk mengurangi stok produk

🔹 Order Service (Golang)
Database: MySQL (orders table)

Cache: Redis

Endpoints:

POST /orders → membuat order baru (validasi productId ke product-service)

GET /orders/product/:productId → menampilkan daftar order (cached)

Event Flow:

Mengirim event order.created ke Kafka agar product-service tahu stok perlu dikurangi

🔹 Event & Communication (Kafka)
Topik utama: order.created, product.created

Kelebihan: Sistem loose-coupled, scalable, dan mampu menangani ribuan request per detik

🔹 Caching (Redis)
Cache digunakan untuk data yang sering diakses:

Produk → GET /products/:id

Order berdasarkan produk → GET /orders/product/:productId

🧪 API Testing (Postman)
 Product Service Endpoints
 Create Product
POST http://localhost:3000/products
Content-Type: application/json

{
  "name": "Nasi Goreng",
  "price": 25000,
  "qty": 100
}

Get Product by ID (Cached)
GET http://localhost:3000/products/1

 Order Service Endpoints
 Create Order
POST http://localhost:3001/orders
Content-Type: application/json

{
  "productId": 1,
  "qty": 2
}

Get Orders by Product ID (Cached)
GET http://localhost:3001/orders/product/1

🧩 Fullstack Developer Test Challenge
🚀 How to Run the Stack Locally

Clone repository dan pastikan Docker sudah ter-install:

git clone https://github.com/username/fullstack-test.git
cd fullstack-test


Jalankan semua container:

docker-compose up -d --build

🔍 Default Ports
Service	URL / Host	Port
Product Service	http://localhost:3000
	3000
Order Service	http://localhost:3001
	3001
MySQL	localhost:3306	3306
Redis	localhost:6379	6379
Kafka (optional)	localhost:9092	9092
⚙️ Environment Variables
🧱 product-service/.env
DATABASE_URL="mysql://root:root@mysql:3306/fullstack-test?parseTime=true"
REDIS_HOST=redis
REDIS_PORT=6379
KAFKA_BROKER=kafka:9092

📦 order-service/.env
DB_HOST=mysql
DB_PORT=3306
DB_NAME=fullstack-test
DB_USER=root
DB_PASSWORD=root

REDIS_HOST=redis
REDIS_PORT=6379

PRODUCT_SERVICE_URL=http://product-service:3000
# KAFKA_BROKER=kafka:9092

🏗️ Brief Explanation of Architecture

Aplikasi ini dibangun dengan arsitektur microservices yang terdiri dari dua service utama:
Product Service (NestJS) dan Order Service (Golang).
Keduanya berkomunikasi secara event-driven menggunakan Apache Kafka, dengan MySQL sebagai database dan Redis untuk caching data yang sering diakses.

🔹 Product Service (NestJS)

Database: MySQL (products table)

Cache: Redis

Endpoints:

POST /products → menambahkan produk baru

GET /products/:id → mengambil detail produk (menggunakan cache Redis)

Event Flow:

Mengirim event product.created ke Kafka saat produk dibuat

Mendengarkan event order.created untuk mengurangi stok produk

Alur Singkat:
Produk dibuat → kirim product.created → saat ada order.created → update stok produk.

🔹 Order Service (Golang)

Database: MySQL (orders table)

Cache: Redis

Endpoints:

POST /orders → membuat order baru (validasi productId ke product-service)

GET /orders/product/:productId → menampilkan daftar order (cached)

Event Flow:

Mengirim event order.created ke Kafka agar product-service tahu stok perlu dikurangi

Memproses event order.created di background melalui consumer

Alur Singkat:
Order dibuat → validasi produk → simpan ke DB → kirim event order.created → product-service update stok.

🔹 Event & Communication (Kafka)

Topik utama:

order.created → dikirim oleh order-service, diterima oleh product-service

product.created → dikirim oleh product-service (opsional untuk logging)

Kelebihan:
Sistem menjadi loose-coupled, scalable, dan mampu menangani ribuan request per detik.

🔹 Caching (Redis)

Cache digunakan untuk data yang sering diakses:

Produk → GET /products/:id

Order berdasarkan produk → GET /orders/product/:productId

Cache diperbarui otomatis saat data berubah di database.

🔹 Database (MySQL)

Product Service: tabel products (id, name, price, qty, createdAt)

Order Service: tabel orders (id, productId, totalPrice, status, createdAt)

Setiap service memiliki database sendiri (Database per Service).

🔹 Infrastructure (Docker Compose)

Seluruh komponen dijalankan menggunakan Docker Compose untuk memudahkan setup dan isolasi environment.

Container yang digunakan:

product-service (NestJS)

order-service (Golang)

mysql

redis

zookeeper

kafka

Command untuk menjalankan:

docker-compose up --build

🧪 API Testing (Postman)

Koleksi Postman tersedia di folder:

/postman_collection.json


Import file tersebut ke Postman untuk mencoba semua endpoint di bawah ini.

🔹 Product Service
Create Product
POST http://localhost:3000/products
Content-Type: application/json
Body:
{
  "name": "Nasi Goreng",
  "price": 25000,
  "qty": 100
}

Get Product by ID (Cached)
GET http://localhost:3000/products/1

🔹 Order Service
Create Order
POST http://localhost:3001/orders
Content-Type: application/json
Body:
{
  "productId": 1,
  "qty": 2
}

Get Orders by Product ID (Cached)
GET http://localhost:3001/orders/product/1

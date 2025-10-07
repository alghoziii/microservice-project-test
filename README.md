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
DATABASE_URL="mysql://root:root@mysql:3306/fullstack-test?parseTime=true"
REDIS_HOST=redis
REDIS_PORT=6379
KAFKA_BROKER=kafka:9092

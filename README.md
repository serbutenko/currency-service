# **Сurrency-service**

gRPC-сервис для конвертации валют с кэшированием в Redis.  
Курсы берутся из [exchangerate.host](https://exchangerate.host/).

---

## **🛠 Используемые технологии**
| Технология  | Описание  |
|------------|-----------|
| **Go** | Основной язык разработки |
| **gRPC** | Коммуникация между сервисами |
| **Protocol Buffers (Protobuf)** | Описание API и сериализация данных |
| **Redis** | Хранилище для кэширования валютных курсов и списка валют |
| **Docker** | Контейнеризация сервиса |
| **Docker Compose** | Оркестрация контейнеров |
| **Exchangerate.host API** | Внешний источник данных о курсах валют |

---

## **⚙️ Установка и запуск**
### **1️⃣ Клонирование репозитория**
```bash
git clone https://github.com/serbutenko/currency-service.git
cd currency-service
```

### **2️⃣ Настройка окружения**
```bash
cp .env.example .env
```
и указать свой API_KEY от [exchangerate.host](https://exchangerate.host/).

### **3️⃣ Запуск через Docker Compose**
```bash
docker-compose up --build
```
Поднимутся два контейнера:
- `currency-service` (gRPC-сервис, порт `50051`)
- `currency-redis` (Redis, порт `6379`)

---

## **📡 gRPC API Методы**
| Метод  | Описание |
|--------|----------|
| **GetRate** | Получить курса валюты A по отношению к валюте B. |
| **Convert** | Конвертация определённой суммы из одной валюты в другую с учетом курса. |
| **ListCurrencies** | Предоставить список всех валют, которые поддерживаются. |

#### GetRate:
```bash
grpcurl -plaintext -d '{"from":"USD","to":"EUR"}' localhost:50051 currency.v1.CurrencyService/GetRate
```
#### Convert:
```bash
grpcurl -plaintext -d '{"from":"USD","to":"EUR", "amount":100}' localhost:50051 currency.v1.CurrencyService/Convert
```
#### ListCurrencies:
```bash
grpcurl -plaintext localhost:50051 currency.v1.CurrencyService/ListCurrencies
```

---

## **📂 Структура проекта**
```
currency-service/
├── cmd/
│ └── server/ # main.go — точка входа
├── internal/
│ ├── app/ # бизнес-логика
│ ├── ports/grpc/ # gRPC-хендлеры
│ ├── adapters/ # интеграция с внешними API
│ ├── cache/ # кэш (Redis, in-memory)
│ └── config/ # конфигурация
├── proto/ # gRPC контракты (.proto)
├── Dockerfile
├── docker-compose.yml
├── .env.example
└── README.md
```

---
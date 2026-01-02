# EEG Analyzer Backend API

REST API для анализа EEG (электроэнцефалографических) данных. Поддерживает анализ одного файла по нескольким ритмам и анализ нескольких файлов по одному ритму.

## 🚀 Быстрый старт

### Вариант 1: Запуск без Docker

#### Требования
- Go 1.21 или выше

#### Установка и запуск

```bash
cd backend

# Установить зависимости
go mod download

# Запустить сервер
go run main.go
```

Сервер будет доступен на `http://localhost:3000`

### Вариант 2: Запуск с Docker

#### Требования
- Docker
- Docker Compose

#### Запуск

```bash
cd backend

# Собрать и запустить контейнер
docker-compose up --build

# Или в фоновом режиме
docker-compose up -d --build

# Остановить
docker-compose down
```

### Вариант 3: Сборка бинарника

```bash
cd backend

# Для Windows
go build -o eeg-analyzer.exe

# Для Linux/Mac
go build -o eeg-analyzer

# Запустить
./eeg-analyzer
```

## 📡 API Endpoints

### 1. Health Check

Проверка статуса API.

```
GET /health
```

**Ответ:**
```json
{
  "status": "ok",
  "message": "EEG Analyzer API is running",
  "version": "1.0.0"
}
```

### 2. Analyze EEG

Анализ EEG данных.

```
POST /analyze
```

#### Single Mode (один файл, несколько ритмов)

**Запрос:**
```json
{
  "analysisId": "550e8400-e29b-41d4-a716-446655440000",
  "analysisMode": "SINGLE",
  "file": {
    "id": "file-1",
    "filename": "eeg_data.csv",
    "experimentName": "Alpha Test",
    "timeColumn": "Time",
    "amplitudeColumn": "Amplitude",
    "rawFile": "VGltZSxBbXBsaXR1ZGUKMC4wMDAsMS4yMwo..."
  },
  "rhythms": ["ALPHA", "BETA", "THETA"],
  "brainZone": "FRONTAL"
}
```

**Ответ:**
```json
{
  "analysisId": "550e8400-e29b-41d4-a716-446655440000",
  "analysisMode": "SINGLE",
  "experimentName": "Alpha Test",
  "rhythms": ["ALPHA", "BETA", "THETA"],
  "absolutePowers": [
    ["ALPHA", 12.45],
    ["BETA", 8.32],
    ["THETA", 5.67]
  ],
  "relativePowers": [
    ["ALPHA", 45.2],
    ["BETA", 30.1],
    ["THETA", 24.7]
  ],
  "dataByRhythm": {
    "ALPHA": {
      "psdPlot": {
        "seriesMetadata": [{"dataKey": "psd", "legend": "PSD"}],
        "data": [{"x": 0, "psd": 1.5}, {"x": 1, "psd": 2.3}],
        "yLogarithmic": true
      },
      "signalPlot": {
        "seriesMetadata": [
          {"dataKey": "raw", "legend": "Raw Signal"},
          {"dataKey": "filtered", "legend": "Filtered Signal"}
        ],
        "data": [{"x": 0, "raw": 12.5, "filtered": 11.2}]
      }
    }
  }
}
```

#### Group Mode (несколько файлов, один ритм)

**Запрос:**
```json
{
  "analysisId": "550e8400-e29b-41d4-a716-446655440001",
  "analysisMode": "GROUP",
  "files": [
    {
      "id": "file-1",
      "filename": "subject1.csv",
      "experimentName": "Subject 1",
      "timeColumn": "Time",
      "amplitudeColumn": "Amplitude",
      "rawFile": "VGltZSxBbXBsaXR1ZGUKMC4wMDAsMS4yMwo..."
    },
    {
      "id": "file-2",
      "filename": "subject2.csv",
      "experimentName": "Subject 2",
      "timeColumn": "Time",
      "amplitudeColumn": "Signal",
      "rawFile": "VGltZSxTaWduYWwKMC4wMDAsMS41Ngo..."
    }
  ],
  "rhythm": "ALPHA",
  "brainZone": "OCCIPITAL"
}
```

## 🧠 Поддерживаемые ритмы

| Ритм | Частота (Hz) | Описание |
|------|-------------|----------|
| DELTA | 0.5-4 | Глубокий сон |
| THETA | 4-8 | Сонливость, медитация |
| ALPHA | 8-13 | Расслабленное бодрствование |
| BETA | 13-30 | Активность, концентрация |
| GAMMA | 30-100 | Высшие когнитивные функции |
| MU | 8-13 | Моторное торможение |
| LAMBDA | 4-8 | Визуальное исследование |
| KAPPA | 8-13 | Вариант альфа |

## 📊 Формат CSV

CSV файл должен содержать:
- Колонку времени (название указывается в `timeColumn`)
- Колонку амплитуды (название указывается в `amplitudeColumn`)
- Заголовок (первая строка)

**Пример:**
```csv
Time,Amplitude,Channel
0.000,12.45,Fp1
0.004,15.23,Fp1
0.008,18.67,Fp1
```

Файл передаётся в base64 кодировке в поле `rawFile`.

## 🧪 Тестирование

### С помощью curl

```bash
# Health check
curl http://localhost:3000/health

# Анализ (пример с тестовыми данными)
curl -X POST http://localhost:3000/analyze \
  -H "Content-Type: application/json" \
  -d @test_request.json
```

### С помощью Swagger UI

Откройте в браузере:
```
http://localhost:3000/swagger/index.html
```

Swagger UI позволяет:
- Просмотреть документацию API
- Интерактивно тестировать эндпоинты
- Посмотреть примеры запросов и ответов

### Генерация Swagger документации

Если вы изменили код и нужно обновить Swagger документацию:

```bash
# Установить swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Сгенерировать документацию
swag init
```

## 🔧 Конфигурация

### Переменные окружения

- `PORT` - порт сервера (по умолчанию: 3000)
- `GIN_MODE` - режим Gin (release/debug, по умолчанию: release)

### CORS

API настроен для работы с фронтендом на:
- `http://localhost:5173` (локальная разработка)
- `https://vad1mchk.github.io` (продакшн)

Для добавления других доменов отредактируйте [main.go](main.go):
```go
AllowOrigins: []string{"http://localhost:5173", "https://example.com"},
```

## 📁 Структура проекта

```
backend/
├── main.go                 # Entry point
├── go.mod                  # Go dependencies
├── Dockerfile              # Docker configuration
├── docker-compose.yml      # Docker Compose setup
├── models/                 # Data models
│   ├── types.go           # Constants and enums
│   ├── request.go         # Request structures
│   ├── response.go        # Response structures
│   └── errors.go          # Error definitions
├── handlers/              # HTTP handlers
│   ├── health.go          # Health check
│   └── analyze.go         # Analysis endpoint
├── analysis/              # Signal processing
│   ├── csv_parser.go      # CSV parsing
│   ├── filter.go          # Signal filtering
│   ├── fft.go             # FFT and PSD
│   ├── rhythms.go         # Rhythm analysis
│   └── downsampler.go     # Data decimation
└── testdata/              # Test CSV files
    ├── sample_eeg_alpha.csv
    ├── sample_eeg_beta.csv
    └── sample_eeg_theta.csv
```

## 🔬 Алгоритмы обработки сигналов

### 1. Предобработка
- Удаление DC смещения (вычитание среднего)
- Нормализация амплитуды

### 2. Фильтрация
- Butterworth bandpass фильтр для каждого ритма
- Удаление шума вне диапазона частот

### 3. FFT и PSD
- Применение окна Хэмминга для уменьшения спектральной утечки
- БПФ (Fast Fourier Transform)
- Расчёт спектральной плотности мощности (PSD)

### 4. Извлечение мощности ритмов
- Абсолютная мощность: средняя PSD в частотном диапазоне
- Относительная мощность: процент от общей мощности

### 5. Децимация данных (Downsampling)
- **LTTB** (Largest-Triangle-Three-Buckets) - сохраняет визуальные характеристики
- Уменьшает количество точек до 1000-2000 для эффективной визуализации
- Применяется перед отправкой данных фронтенду

## 🐛 Отладка

### Логи

```bash
# В режиме debug
GIN_MODE=debug go run main.go

# В Docker с логами
docker-compose up
```

### Проверка зависимостей

```bash
go mod verify
go mod tidy
```

### Типичные проблемы

#### Порт занят
```
Error: listen tcp :3000: bind: address already in use
```
Решение: Измените порт через переменную окружения
```bash
PORT=3001 go run main.go
```

#### CORS ошибки
Убедитесь, что домен фронтенда добавлен в `AllowOrigins` в [main.go](main.go).

## 🚀 Деплой

### VPS/Cloud

1. Скомпилируйте бинарник:
```bash
GOOS=linux GOARCH=amd64 go build -o eeg-analyzer
```

2. Загрузите на сервер:
```bash
scp eeg-analyzer user@server:/path/to/app/
```

3. Запустите:
```bash
./eeg-analyzer
```

### Docker

```bash
# Собрать образ
docker build -t eeg-analyzer .

# Запустить
docker run -p 3000:3000 eeg-analyzer
```

### systemd service (Linux)

Создайте `/etc/systemd/system/eeg-analyzer.service`:
```ini
[Unit]
Description=EEG Analyzer API
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/eeg-analyzer
ExecStart=/opt/eeg-analyzer/eeg-analyzer
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable eeg-analyzer
sudo systemctl start eeg-analyzer
```

## 📝 Лицензия

MIT

## 🤝 Вклад

Фронтенд разработан отдельной командой и доступен на: https://vad1mchk.github.io/CourseworkEEG/

## 📧 Контакты

При возникновении вопросов или проблем создайте issue в репозитории.

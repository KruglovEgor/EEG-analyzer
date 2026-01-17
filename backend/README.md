# EEG Analyzer Backend

REST API для анализа ЭЭГ (электроэнцефалография) сигналов с поддержкой множественных ритмов и режимов анализа.

## Запуск

### Локально
```bash
go mod download
go run .
```

### Docker
```bash
docker-compose up --build
```

Сервер запускается на порту `3000`. Swagger документация доступна на `/swagger/index.html`.

## API Endpoints

### POST /analyze
Мультичастотный анализ ЭЭГ сигнала (multipart/form-data)

**Режим SINGLE**: один файл, несколько ритмов
```bash
curl -X POST http://localhost:3000/analyze \
  -F "file=@data.csv" \
  -F "analysisMode=SINGLE" \
  -F "analysisId=test-123" \
  -F "experimentName=Test" \
  -F "timeColumn=time" \
  -F "amplitudeColumn=amplitude" \
  -F "rhythms=ALPHA,BETA,THETA"
```

**Режим GROUP**: несколько файлов, один ритм
```bash
curl -X POST http://localhost:3000/analyze \
  -F "file=@subject1.csv" \
  -F "file=@subject2.csv" \
  -F "analysisMode=GROUP" \
  -F "analysisId=test-456" \
  -F "experimentName=Subject 1" \
  -F "experimentName=Subject 2" \
  -F "rhythm=ALPHA" \
  -F "filterMin=8" \
  -F "filterMax=13"
```

### POST /preview
Предпросмотр эффекта фильтров (multipart/form-data)

```bash
curl -X POST http://localhost:3000/preview \
  -F "file=@data.csv" \
  -F "previewId=preview-1" \
  -F "experimentName=Test" \
  -F "rhythm=ALPHA" \
  -F 'filterParams={"filterMin":8,"filterMax":13,"filterOrder":2,"nPerSeg":1024,"nOverlap":512}'
```

### GET /health
Проверка состояния сервера

```bash
curl http://localhost:3000/health
```

## Поддерживаемые ритмы

| Ритм | Частота (Hz) | Описание |
|------|-------------|----------|
| DELTA | 0.5-4 | Глубокий сон |
| THETA | 4-8 | Дремота, медитация |
| ALPHA | 8-13 | Расслабленное бодрствование |
| BETA | 13-30 | Активное мышление |
| GAMMA | 30-100 | Высшие когнитивные функции |
| MU | 8-13 | Моторное торможение |
| LAMBDA | 4-8 | Визуальное сканирование |
| KAPPA | 8-13 | Вариант альфа |

## Параметры фильтров

```json
{
  "filterMin": 8.0,        // Нижняя граница частоты (Hz)
  "filterMax": 13.0,       // Верхняя граница частоты (Hz)
  "filterOrder": 1,        // Порядок Баттерворта (1-4)
  "nPerSeg": 1024,         // Размер сегмента для Welch
  "nOverlap": 512          // Перекрытие сегментов
}
```

### Ограничения
- **filterOrder**: минимум 1 (рекомендуется 1-4)
- **filterMin/Max**: должны быть >0, max > min
- **nPerSeg**: должен быть >0 и <длины сигнала
- **nOverlap**: должен быть ≥0 и <nPerSeg
- **Размер файла**: максимум 32 МБ

## Алгоритм обработки

1. **Удаление DC смещения**: вычитание среднего значения
2. **FFT пре-фильтр (0.5-40 Hz)**: только для PSD расчёта
3. **Welch PSD**: спектральная плотность мощности
4. **Трапециевидная интеграция**: расчёт абсолютной/относительной мощности
5. **Butterworth фильтр**: для визуализации сигнала
6. **LTTB даунсемплинг**: сжатие до 2000 точек

## Технологии

- Go 1.21+
- Gin Web Framework
- Swagger/OpenAPI
- Docker & Docker Compose

# Extension Backend

Go-сервис для браузерного расширения. Принимает текст, отправляет в Gemini AI и возвращает объяснение, summary или глоссарий.

## Требования

- Go 1.24+
- API-ключ Gemini — получить на [aistudio.google.com/app/apikey](https://aistudio.google.com/app/apikey)

## Запуск

**1. Клонировать репозиторий**

```bash
git clone https://github.com/Alpha-x-prog/extension_backend.git
cd extension_backend
```

**2. Создать файл `.env`**

```bash
cp .env.example .env
```

Открыть `.env` и вставить свой ключ:

```
GEMINI_API_KEY=твой_ключ_сюда
PORT=8080
```

**3. Установить зависимости**

```bash
go mod download
```

**4. Запустить сервер**

```bash
# Linux / macOS
export $(cat .env | xargs) && go run .

# Windows PowerShell
$env:GEMINI_API_KEY="твой_ключ"; $env:PORT="8080"; go run .
```

Сервер запустится на `http://localhost:8080`.

---

## API

### POST `/explain`

Объясняет выделенный фрагмент текста простым языком.

**Запрос:**
```json
{ "text": "Текст для объяснения" }
```

**Ответ:**
```json
{ "answer": "Развёрнутое объяснение с глоссарием и выводами..." }
```

| Ограничение | Значение |
|---|---|
| Максимальная длина text | 5000 символов |

---

### POST `/summary`

Структурированный пересказ большого текста в формате Markdown.

**Запрос:**
```json
{ "text": "Большой текст для пересказа..." }
```

**Ответ:**
```json
{ "summary_md": "# Summary\n\n..." }
```

| Ограничение | Значение |
|---|---|
| Минимальная длина text | 300 символов |
| Максимальная длина text | 20 000 символов |

---

### POST `/glossary`

Выделяет термины и аббревиатуры из текста и объясняет их.

**Запрос:**
```json
{ "text": "Текст с терминами..." }
```

**Ответ:**
```json
{ "glossary_md": "# Глоссарий\n\n## REST API\n..." }
```

| Ограничение | Значение |
|---|---|
| Минимальная длина text | 100 символов |
| Максимальная длина text | 5000 символов |

---

## Коды ответов

| Код | Причина |
|---|---|
| 200 | Успех |
| 400 | Невалидный JSON, пустой текст или текст слишком короткий |
| 405 | Неправильный HTTP-метод (ожидается POST) |
| 413 | Текст превышает максимальную длину |
| 502 | AI-сервис недоступен |
| 504 | Превышен таймаут запроса к AI (30 сек) |

**Формат ошибки:**
```json
{ "error": "описание ошибки" }
```

---

## Структура проекта

```
extension_backend/
├── main.go                      # Точка входа
├── go.mod / go.sum              # Зависимости
├── .env.example                 # Пример переменных окружения
└── internal/
    ├── service/
    │   └── gemini.go            # Gemini AI клиент
    └── handler/
        ├── common.go            # Утилиты (writeJSON, writeAIError)
        ├── explain.go           # POST /explain
        ├── summary.go           # POST /summary
        └── glossary.go         # POST /glossary
```

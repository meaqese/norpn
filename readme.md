

# norpn (HTTP Calc)

## Описание

Этот проект представляет собой сервер, реализующий http для вычисления арифметических выражений. Он принимает строковое арифметическое выражение в формате json через http post-запрос (и не только :), выполняет его вычисление и возвращает результат в формате json.

Калькулятор поддерживает стандартные арифметические операции:
- Сложение (+)
- Вычитание (-)
- Умножение (*)
- Деление (/)

Поддерживаются круглые скобки для задания порядка выполнения операций, а так же скобки внутри скобок.

---

## Установка и запуск

### Шаг 1: Клонировать репозиторий

Клонируем репу с гитхаба и переходим на склонированную папку:

```bash
git clone https://github.com/meaqese/norpn.git
cd norpn
```

### Шаг 2: Запуск приложения

Для запуска сервера выполните команду:

```bash
go run ./cmd/main.go
```

После запуска сервер будет доступен на порту `8080` по адресу: `http://localhost:8080`. 
Порт можно менять в переменных окружения - `PORT`.

---

## Описание API

### 1. `/api/v1/calculate`

**Метод:** POST

**Описание:** Принимает арифметическое выражение в JSON формате, выполняет вычисления и возвращает результат.

#### Пример запроса:

```bash
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-d '{"expression": "(2+2)*2"}'
```

Если курл не сработал, можете попробовать экранировать кавычки, - `\"`, или лучше используйте Postman

**Пример ответа (успешное выполнение):**

```json
{"result": "8"}
```

### Возможные ошибки

1. **Ошибка 422 (Unprocessable Entity):**

Это ошибка возникает, если передано недопустимое выражение (например, если в нем присутствуют посторонние символы или есть синтаксическая ошибка).

Пример запроса:

```bash
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-d '{"expression": "2+x"}'
```

**Пример ответа (ошибка 422):**

```json
{"error":"Expression is not valid","description":"no one value to return"}
```

2. **Ошибка 500 (Internal Server Error):**

Эта ошибка возникает, если произошло что-то непредвиденное при выполнении вычислений.

Ошибку возпроизвести не удалось, поскольку обрабатываются известные ошибки, но можете попытаться.

**Пример ответа (ошибка 500):**

```json
{"error":"Internal server error","description":""}
```
---

### Тесты

**http api**
```bash
go test ./internal/transport/rest
```

**калькулятор**
```bash
go test ./pkg/norpn
```


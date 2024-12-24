# Финальный проект 1 семестра

REST API сервис для загрузки и выгрузки данных о ценах.

## Требования к системе

Для успешной работы сервиса необходимы следующие компоненты:
- **Go**: версия 1.23.3 или выше. Убедитесь, что язык программирования установлен корректно и доступен в `PATH`.
- **PostgreSQL**: версия 13.3 или выше.
  Настройки базы данных должны соответствовать следующим параметрам:
  - Хост: `localhost`.
  - Порт: `5432`.
  - Имя базы данных: `project-sem-1`.
  - Пользователь: `validator`.
  - Пароль: `val1dat0r`.

## Установка и запуск

### Шаг 1. Подготовка окружения

1. **Клонирование репозитория**  
   Скачайте исходный код проекта из удалённого репозитория:
   ```bash
   git clone https://github.com/dmitryfedoroff/itmo-devops-sem1-project-template.git
   cd itmo-devops-sem1-project-template
   ```

2. **Установка зависимостей**  
   Проверьте установленную версию Go:
   ```bash
   go version
   ```
   Затем выполните установку необходимых модулей:
   ```bash
   go mod tidy
   ```

3. **Настройка базы данных**  
   Убедитесь, что PostgreSQL запущен. Выполните подключение к серверу и создайте необходимые сущности:
   ```bash
   psql -h localhost -p 5432 -U postgres
   ```
   Внутри консоли PostgreSQL выполните следующие команды:
   ```sql
   CREATE DATABASE "project-sem-1";
   CREATE USER validator WITH PASSWORD 'val1dat0r';
   GRANT ALL PRIVILEGES ON DATABASE "project-sem-1" TO validator;
   ```

4. **Применение миграций**  
   Для создания таблицы `prices` выполните:
   ```bash
   chmod +x scripts/prepare.sh
   ./scripts/prepare.sh
   ```

### Шаг 2. Запуск сервиса

После завершения подготовки запустите приложение следующей командой:
```bash
chmod +x scripts/run.sh
./scripts/run.sh
```

## Тестирование

Директория `sample_data` - это пример директории, которая является разархивированной версией файла `sample_data.zip`

### Подготовка тестов

Для проверки функциональности выполните:
1. Убедитесь, что тестовый скрипт обладает правами на выполнение:
   ```bash
   chmod +x scripts/tests.sh
   ```
2. Запустите тесты всех уровней:
   ```bash
   for idx in {1..3}; do ./scripts/tests.sh $idx; done
   ```

### Уровни тестирования

<details>
<summary><b>Уровень 1: базовая проверка корректности работы POST и GET запросов</b></summary>

   На этом уровне тестирования выполняются базовые проверки:
   - POST-запросы на загрузку данных в формате CSV: проверка успешной обработки корректных данных.  
   - GET-запросы для выгрузки всех данных: проверяется, что файл выгрузки содержит корректные записи.  
   - Работа PostgreSQL с минимальными запросами: подсчёт общего числа записей в таблице.  
   Результат: Все тесты успешно завершены.

![](/screenshots/github_test_level_1.png)

![](/screenshots/local_test_level_1.png)

</details>

<details>
<summary><b>Уровень 2: тестирование загрузки данных через архивы (ZIP/TAR) и фильтров</b></summary>

   На этом уровне тестирования проверяются:  
   - POST-запросы с загрузкой данных в форматах ZIP и TAR: убедиться, что архивы корректно обрабатываются.  
   - Базовые GET-запросы: проверка успешной выгрузки данных в виде архива ZIP.  
   - Тесты API первого уровня: POST и GET запросы с корректными данными.  
   - Работа PostgreSQL с простыми запросами для подсчёта количества записей, уникальных категорий и общей стоимости.  
   Результат: Все проверки завершены успешно.

![](/screenshots/github_test_level_2.png)

![](/screenshots/local_test_level_2.png)

</details>

<details>
<summary><b>Уровень 3: сложные сценарии, включая работу с некорректными данными, дубликатами и продвинутую статистику</b></summary>

   На этом уровне тестирования проверяются:
   - Корректная обработка POST-запросов с архивами (ZIP) и некорректными данными: убедиться, что сервис фильтрует неверные записи и корректно подсчитывает статистику (например, общее количество записей, количество дубликатов).
   - Успешное обнаружение дубликатов в загружаемых данных.
   - Функциональность GET-запросов с фильтрами: проверяется, что выгружаемые данные не содержат некорректных записей.
   - Работа PostgreSQL с запросами на сложные выборки, включая использование фильтров по дате, цене и подсчёт статистики.  
   Результат: Все проверки прошли успешно.

![](/screenshots/github_test_level_3.png)

![](/screenshots/local_test_level_3.png)

</details>

<details>
<summary><b>Финальная проверка</b></summary>

Финальная проверка подтверждает, что хотя бы один уровень тестов успешно пройден.

![](/screenshots/github_check_test_results.png)

</details>

## Пример использования API

### Загрузка данных (POST)
Загрузите данные о ценах с помощью следующей команды:
```bash
curl -F "file=@sample_data.zip" http://localhost:8080/api/v0/prices
```

### Выгрузка данных с фильтрами (GET)
Получите данные за конкретный период и в определённом ценовом диапазоне:
```bash
curl "http://localhost:8080/api/v0/prices?start=2024-01-01&end=2024-01-31&min=100&max=1000" -o response.zip
```

## Возможные проблемы и их решение

1. **Используется устаревшая версия Go**

   Если команда `go version` возвращает версию Go ниже `1.23.3`, это может привести к несовместимости с проектом.

   Найдите текущую версию Go:
   ```bash
   go version
   ```
   Если версия ниже `1.23.3`, удалите устаревшую версию:
   ```bash
   sudo apt remove --purge <устаревшая версия>
   sudo rm -rf /usr/local/go
   ```
   Скачайте архив с версией `1.23.3`:
   ```bash
   wget https://go.dev/dl/go1.23.3.linux-amd64.tar.gz
   ```
   Установите Go:
   ```bash
   sudo tar -C /usr/local -xzf go1.23.3.linux-amd64.tar.gz
   ```
   Добавьте Go в `PATH`:
   ```bash
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   source ~/.bashrc
   ```
   Проверьте, что версия Go обновлена:
   ```bash
   go version
   ```
   Ожидаемый результат:
   ```
   go version go1.23.3 linux/amd64
   ```

2. **Ошибка подключения к базе данных**  
   Убедитесь, что PostgreSQL работает, а параметры подключения совпадают с указанными в `config.yml`.

3. **Конфликт портов**  
   Если порт `8080` занят, завершите процесс, который использует порт, или измените порт в конфигурационном файле `config.yml`.

4. **Ошибка миграции базы данных**  
   Проверьте правильность данных для подключения к базе данных, а также доступность PostgreSQL.

## Контакты

- **Автор**: Дмитрий Федоров
- **Эл. почта**: [fedoroffx@gmail.com](mailto:fedoroffx@gmail.com)
- **Telegram**: [https://t.me/dmitryfedoroff](https://t.me/dmitryfedoroff) 
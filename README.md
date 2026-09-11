# Posts Service

Сервис для добавления и чтения постов и комментариев с использованием GraphQL.

Проект реализован на Go с разделением на слои в соответствии с принципами чистой архитектуры. Основные пакеты:

- `cmd` - точка входа в приложение;
- `config` - конфигурация приложения;
- `internal/app` - инициализация зависимостей и запуск сервера;
- `internal/delivery` - транспортный слой, GraphQL API;
- `internal/usecase` - бизнес-логика приложения;
- `internal/repository` - работа с хранилищем данных;
- `internal/entity` - основные сущности приложения;
- `internal/mocks` - сгенерированные моки для unit-тестов;
- `db` - SQL-схема базы данных.

Основные слои приложения покрыты unit-тестами: отдельно тестируются usecase, repository и GraphQL delivery. Для изоляции зависимостей используются сгенерированные моки

Проект поддерживает два варианта хранения данных:

- in-memory
- PostgreSQL

Тип хранилища выбирается при запуске приложения с помощью переменной окружения `STORAGE`.

In-memory:

```bash
STORAGE=inmemory go run ./cmd/main
```

Запуск с PostgreSQL:

```bash
STORAGE=postgres go run ./cmd/main
```

Параметры подключения PostgreSQL берутся из `config/config.yaml` и при необходимости могут быть переопределены через `PG_HOST`, `PG_PORT`, `PG_USER`, `PG_PASSWORD`, `PG_DB`, `PG_SSL_MODE`.

## Docker

Для распространения сервиса используется Docker.

Сборка образа:

```bash
docker build -t posts-service .
```

Запуск сервиса с in-memory хранилищем:

```bash
docker run --rm \
  -p 8000:8000 \
  -e STORAGE=inmemory \
  -e SRV_HOST=0.0.0.0 \
  posts-service
```

Для запуска сервиса вместе с PostgreSQL используется Docker Compose:

```bash
docker compose up --build
```

При запуске через `docker compose`:

- поднимается контейнер с приложением;
- поднимается контейнер PostgreSQL;
- приложение запускается с `STORAGE=postgres`;
- SQL-схема из `db/ddl.sql` применяется при инициализации базы данных.

---

## GraphQL API

После запуска сервиса GraphQL Playground доступен по адресу:

http://localhost:8000/

GraphQL endpoint:

http://localhost:8000/graphql

---

## Возможности

### Посты

#### Создание поста

При создании поста можно указать, разрешено ли оставлять под ним комментарии.

Пост с разрешёнными комментариями:

```graphql
mutation {
  createPost(
    params: {
      userId: 1
      title: "First post"
      description: "Post description"
      commentsAllowed: true
    }
  ) {
    id
    userId
    title
    description
    commentsAllowed
  }
}
```

Пост с запрещёнными комментариями:

```graphql
mutation {
  createPost(
    params: {
      userId: 1
      title: "Closed post"
      description: "Comments are disabled"
      commentsAllowed: false
    }
  ) {
    id
    userId
    title
    description
    commentsAllowed
  }
}
```

#### Получение списка постов

Список постов возвращается с пагинацией через `limit` и `offset`.

```graphql
query {
  posts(
    params: {
      limit: 10
      offset: 0
    }
  ) {
    len
    posts {
      id
      userId
      title
      description
      commentsAllowed
    }
  }
}
```

#### Получение поста по ID

```graphql
query {
  post(id: 1) {
    id
    userId
    title
    description
    commentsAllowed
  }
}
```


### Комментарии

#### Создание комментария

Комментарий верхнего уровня создаётся без `replyToCommentId`.

```graphql
mutation {
  createComment(
    params: {
      postId: 1
      userId: 2
      text: "First comment"
    }
  ) {
    id
    postId
    replyToCommentId
    userId
    text
  }
}
```

#### Создание ответа на комментарий

Для ответа на другой комментарий указывается `replyToCommentId`.

```graphql
mutation {
  createComment(
    params: {
      postId: 1
      replyToCommentId: 1
      userId: 3
      text: "Reply to comment"
    }
  ) {
    id
    postId
    replyToCommentId
    userId
    text
  }
}
```

#### Получение комментариев поста

Возвращаются комментарии верхнего уровня. Для списка используется пагинация через `limit` и `offset`.

```graphql
query {
  comments(
    params: {
      postId: 1
      limit: 10
      offset: 0
    }
  ) {
    len
    comments {
      id
      postId
      replyToCommentId
      userId
      text
    }
  }
}
```

#### Получение ответов на комментарий

Ответы на конкретный комментарий также возвращаются с пагинацией.

```graphql
query {
  replies(
    params: {
      postId: 1
      replyToCommentId: 1
      limit: 10
      offset: 0
    }
  ) {
    len
    comments {
      id
      postId
      replyToCommentId
      userId
      text
    }
  }
}
```

### GraphQL Subscriptions

Для новых комментариев реализована GraphQL Subscription.

Клиент может подписаться на определённый пост и получать новые комментарии без повторного выполнения query.

```graphql
subscription {
  commentAdded(postId: 1) {
    id
    postId
    replyToCommentId
    userId
    text
  }
}
```

После создания нового комментария к посту `1` подписанный клиент автоматически получит его через открытое WebSocket-соединение.

---

## Производительность

GraphQL API спроектирован таким образом, чтобы избежать классической проблемы N+1.

Комментарии не загружаются через вложенные резолверы `Post.comments` или `Comment.replies`. Вместо этого для получения комментариев и ответов используются отдельные запросы `comments` и `replies` с указанием `postId` или `replyToCommentId`.

Таким образом, получение списка постов не приводит к выполнению отдельного запроса к хранилищу для комментариев каждого поста.

Каждый комментарий содержит `postId` и `replyToCommentId`. Для комментариев верхнего уровня `replyToCommentId` равен `null`, а для ответов содержит ID комментария, на который был создан ответ. Такая модель позволяет строить дерево комментариев с неограниченной глубиной вложенности.

Комментарии загружаются по уровням: верхний уровень через `comments`, ответы на конкретный комментарий через `replies`. Это позволяет не загружать всё дерево комментариев целиком.

Для списков комментариев и ответов используется пагинация через `limit` и `offset`.

Для безопасного одновременного доступа к данным in-memory хранилища используется `sync.RWMutex`, позволяющий параллельно выполнять операции чтения.

Доступ к списку подписчиков GraphQL Subscriptions также синхронизирован с помощью `sync.RWMutex`.
---

## Обработка ошибок и граничных случаев

В сервисе предусмотрена обработка основных некорректных сценариев при работе с постами и комментариями.

Для постов:

- заголовок поста не может быть пустым;
- длина заголовка ограничена 200 символами;
- длина описания ограничена 10000 символами;
- `id` поста должен быть больше 0;
- `limit` должен быть больше 0;
- `offset` не может быть отрицательным;
- если `limit` больше 50, он автоматически ограничивается значением 50;
- при запросе несуществующего поста возвращается ошибка.

Для комментариев:

- текст комментария не может быть пустым;
- длина текста комментария ограничена 2000 символами;
- `postId` должен быть больше 0;
- нельзя создать комментарий к несуществующему посту;
- нельзя создать комментарий, если у поста запрещены комментарии;
- при создании ответа проверяется существование комментария, на который отправляется ответ;
- нельзя создать ответ на комментарий, принадлежащий другому посту;
- при запросе ответов проверяется существование родительского комментария;
- `replyToCommentId` при запросе ответов должен быть указан и быть больше 0;
- `limit` должен быть больше 0;
- `offset` не может быть отрицательным;
- если `limit` больше 50, он автоматически ограничивается значением 50.
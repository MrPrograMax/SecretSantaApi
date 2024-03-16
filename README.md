# Тайный Санта RESTFul API

Этот проект представляет собой RESTFul API сервис для игры в Тайного Санту. Сервис позволяет управлять группами и участниками игры, проводить жеребьевку и получать информацию о том, кому нужно подарить подарок.

## Типы данных
#### Модель 1: Group
```go
type Group struct {
  Id          int    json:"id" db:"id"
  Name        string json:"name" db:"name" binding:"required"
  Description string json:"description" db:"description"
}
```
Описание:
- **Id**: уникальный идентификатор группы.
- **Name**: название группы (обязательное поле).
- **Description**: описание группы.

#### Модель 2: Participant
```go
type Participant struct {
  Id          int    json:"id" db:"id"
  Name        string json:"name" db:"name" binding:"required"
  Wish        string json:"wish" db:"wish"
  RecipientId int    json:"recipientId" db:"recipient_id"
}
```
Описание:
- **Id**: уникальный идентификатор участника.
- **Name**: имя участника (обязательное поле).
- **Wish**: желание участника.
- **RecipientId**: идентификатор получателя (до жеребьевки null).

## Доступные действия

### Группы

| Метод       | Описание               | Комментарий|
| ------------- |:------------------|:------------------|
|`POST`   - **/api/group**     | Добавление группы с возможностью указания названия (name) и описания (description) |  |
|`GET`    - **/api/groups**    | Получение краткой информации о всех группах (без информации об участниках)| |
|`GET`    - **/api/group/:id** | Получение полной информации о группе по идентификатору (с информацией об участниках)| До проведения жеребьевки recipient у участников не заполнен (null)|
|`PUT`    - **/api/group/:id** | Редактирование группы по идентификатору. Можно редактировать только свойства name и description|
|`DELETE` - **/api/group/:id** | Удаление группы по идентификатору |



### Участники

| Метод       | Описание               | Комментарий|
| ------------- |:------------------|:------------------|
|`POST`   - **/api/group/:id/participant**                          | Добавление участника в группу по идентификатору группы |  |
|`POST`   - **/api/group/:id/toss**                                 | Проведение жеребьевки в группе по идентификатору группы| <li> Проведение жеребьевки возможно только в том случае, когда количество участников группы >= 3 <li> Участнику в качестве подопечного нельзя выдать самого себя <li> Участник не может быть подопечным одновременно у двух и более участников|
|`GET`    - **/api/group/:id/participant/:participantId/recipient** | Получение информации для конкретного участника группы о том, кому он дарит подарок| |
|`DELETE` - **/api/group/:id/participant/:participantId**           | Удаление участника из группы по идентификаторам группы и участника|


## Техническая информация

- Язык программирования: GoLang
- База данных: Postgres
- Порт запуска: 8080

## Принципы проектирования

Этот проект соответствует принципам SOLID.
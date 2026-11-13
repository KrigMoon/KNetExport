# KNetExport

Модуль на Go, собираемый как C-разделяемая библиотека через cgo. Экспортирует API библиотек [KJournal](https://github.com/KrigMoon/KJournal) и [KNet](https://github.com/KrigMoon/KNet) для вызова из C. Сам по себе модуль не самостоятелен и не имеет собственной бизнес-логики: он принимает данные из C в виде JSON-строк, передаёт их в KJournal и KNet и возвращает результат обратно в C также в виде JSON-строк. KNet использует общий экземпляр KJournal, поэтому журнал запускается до сети. Все возвращённые C-строки освобождаются через FreeString.

## API

### Общее

Каждая экспортируемая функция возвращает C-строку с JSON-объектом FuncResponse. Поле Level принимает значения "Success" или "Error", поле Response содержит полезную нагрузку либо текст ошибки. При успехе Response может быть пустой строкой, если функция не возвращает данных. Все возвращённые указатели выделены в C-куче и освобождаются через FreeString.

| Функция | Принимает | Возвращает |
|---|---|---|
| `FreeString` | `*C.char` | — |

### KNet

Функции API библиотеки KNet. Принимают JSON-строки, возвращают FuncResponse. Для работы сети журнал KJournal должен быть запущен.

| Функция | Принимает | Возвращает |
|---|---|---|
| `RunKNet` | JSON `KNSettings` | FuncResponse |
| `EndworkKNet` | — | FuncResponse |
| `KNetMetrics` | — | FuncResponse, `Response` = JSON метрик |
| `KNExecute` | строку команды | FuncResponse |
| `KNResponse` | — | FuncResponse, `Response` = JSON `Response` или `""` |
| `KNWriteResponse` | JSON `Response` | FuncResponse |
| `KNGetCursor` | — | FuncResponse, `Response` = JSON курсора |
| `KNExternalPrograms` | — | FuncResponse, `Response` = JSON списка программ |
| `KNPeerList` | — | FuncResponse, `Response` = JSON списка пиров |
| `KNAddReaderPackages` | `id *C.char` (nil допустим), `buff C.int` | FuncResponse, `Response` = идентификатор читателя |
| `KNDelReaderPackages` | `id *C.char` | FuncResponse |
| `KNPackage` | `id *C.char` | FuncResponse, `Response` = protobuf-JSON пакета или `""` |

### KJournal

Функции API библиотеки KJournal. Принимают JSON-строки, возвращают FuncResponse.

| Функция | Принимает | Возвращает |
|---|---|---|
| `RunKJournal` | JSON `KJSettings` | FuncResponse |
| `EndworkKJournal` | — | FuncResponse |
| `KJournalMetrics` | — | FuncResponse, `Response` = JSON метрик |
| `KJNewWriter` | — | FuncResponse |
| `KJDelWriter` | — | FuncResponse |
| `KJWrite` | JSON `Record` | FuncResponse |
| `KJNewReader` | `id *C.char` (nil допустим), `buff C.int` | FuncResponse, `Response` = идентификатор читателя |
| `KJDelReader` | `id *C.char` | FuncResponse |
| `KJRead` | `id *C.char` | FuncResponse, `Response` = JSON `Record` или `""` |
| `KJReadGroup` | `id *C.char` | FuncResponse, `Response` = JSON-массив `Record` |

## Экспортируемые структуры

Структуры, передаваемые через JSON между C и Go.

### DataTime

Представление момента времени. Используется в настройках, записях журнала и ответах команд. Используется для однозначной передачи времени или промежутка через границу FFI

| Название | Тип | Описание |
|---|---|---|
| `Year` | `int` | Год |
| `Month` | `int` | Месяц |
| `Day` | `int` | День |
| `Hour` | `int` | Час |
| `Minute` | `int` | Минута |
| `Second` | `int` | Секунда |
| `Miliseconds` | `int` | Миллисекунды |

### FuncResponse

Обёртка ответа, возвращаемая каждой экспортируемой функцией.

| Название | Тип | Описание |
|---|---|---|
| `Level` | `string` | `"Success"` или `"Error"` |
| `Response` | `string` | Полезная нагрузка или текст ошибки |

### KJSettings

Настройки запуска KJournal. Передаются в `RunKJournal`. Поля `AlarmRWGroup` и `AlarmWrite` имеют тип `DataTime` и преобразуются в `time.Duration`.

| Название | Тип | Описание |
|---|---|---|
| `DefaultReaderBuf` | `int` | Размер буфера читателя по умолчанию |
| `MaxInGroup` | `int` | Максимальное число записей в группе |
| `FileDumperDir` | `string` | Каталог дампа журнала в файлы |
| `AlarmRWGroup` | `DataTime` | Порог тревоги чтения/записи группы |
| `AlarmWrite` | `DataTime` | Порог тревоги записи |
| `CrashReport` | `string` | Каталог отчётов о сбоях |

### KNSettings

Настройки запуска KNet. Передаются в `RunKNet`. Поле `MyNode` содержит protobuf-JSON объекта `Node`.

| Название | Тип | Описание |
|---|---|---|
| `ExternalCommandsMap` | `[]string` | Список внешних команд |
| `AutoRun` | `[]string` | Команды, выполняемые при старте |
| `FileStorageDir` | `string` | Каталог файлового хранилища |
| `RemoteControl` | `bool` | Признак удалённого управления |
| `MyNode` | `json.RawMessage` | Описание узла в protobuf-JSON |
| `Password` | `string` | Пароль ядра сети |

### Record

Запись журнала. Передаётся в `KJWrite`, возвращается из `KJRead` и `KJReadGroup`.

| Название | Тип | Описание |
|---|---|---|
| `Author` | `string` | Автор записи |
| `Level` | `string` | Уровень записи |
| `Code` | `string` | Код записи |
| `Time` | `DataTime` | Время создания |
| `Message` | `string` | Текст сообщения |

### Response

Ответ на команду. Передаётся в `KNWriteResponse`, возвращается из `KNResponse`.

| Название | Тип | Описание |
|---|---|---|
| `CommandID` | `string` | Идентификатор команды |
| `Level` | `string` | Уровень ответа |
| `Time` | `DataTime` | Время формирования |
| `ResponseMessage` | `string` | Текст ответа |
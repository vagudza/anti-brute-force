package bucket

// Limiter представляет собой интерфейс для работы с bucket
type Limiter interface {
	// Allow проверяет, можно ли пропустить запрос
	// возвращает true, если запрос допустим, false если лимит превышен
	Allow(key string) (bool, error)

	// Reset сбрасывает bucket для указанного ключа
	Reset(key string) error

	// Close освобождает ресурсы
	Close()
}

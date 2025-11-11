package service

import (
	"time"
)

// SignatureService предоставляет методы для работы с лицензиями:
// генерация лицензий, проверка их подлинности и получение публичного ключа.
type SignatureService interface {
	// Issue создаёт новую лицензию для пользователя с указанным временем действия.
	// Возвращает закодированную строку лицензии и ошибку, если операция не удалась.
	Issue(userID string, duration time.Duration) (string, error)

	// Verify проверяет корректность переданного секретного payload лицензии.
	// Возвращает подпись лицензии или ошибку.
	Verify(secretPayload string) (string, error)

	// GetPublicKey возвращает публичный ключ в формате PEM для проверки лицензий.
	GetPublicKey() ([]byte, error)
}

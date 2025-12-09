package util

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

const UserCookieName = "user_uuid"

// GetUserID извлекает UUID из куки пользователя или генерирует новый,
// если кука не найдена или недействительна. Новый UUID устанавливается в HTTP-ответ.
func GetUserID(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie(UserCookieName)

	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	newUUID := uuid.New().String()

	// Создаем новую куку
	newCookie := &http.Cookie{
		Name:     UserCookieName,
		Value:    newUUID,
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour), // Срок годности: 1 год
		HttpOnly: true,                                 // Защита от XSS (недоступна из JavaScript)
		SameSite: http.SameSiteLaxMode,                 // Защита от CSRF
	}

	http.SetCookie(w, newCookie)

	return newUUID
}

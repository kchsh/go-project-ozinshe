package handlers

func NewApiError(msg string) map[string]string {
	return map[string]string{"error": msg}
}

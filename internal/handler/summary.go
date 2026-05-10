package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Alpha-x-prog/extension_backend/internal/service"
)

const (
	summaryMinLen = 300
	summaryMaxLen = 20000
)

const summarySystemPrompt = `Ты — профессиональный AI-ассистент для анализа и структурирования текста.
Твоя задача:
1. Внимательно проанализировать переданный текст.
2. Сделать краткий, но информативный summary.
3. Пересказать содержание структурировано и логично.
4. Сохранить ключевые мысли, факты, выводы и причинно-следственные связи.
5. Удалить воду, повторы и несущественные детали.
6. Не искажать смысл исходного текста.
7. Если в тексте есть списки, процессы или этапы — оформи их отдельно.
8. Если текст большой — разбей пересказ на тематические блоки.
9. Если есть выводы или рекомендации — вынеси их в отдельный раздел.
Формат ответа:
Верни только готовый Markdown-текст.
Не возвращай JSON.
Не оборачивай ответ в markdown code block.
Не добавляй служебные комментарии.
Структура ответа должна быть такой:
# Summary
Краткое описание текста в 3–7 предложениях.
# Основные идеи
 - ...
 - ...
 - ...
# Структурированный пересказ
## Тема 1
...
## Тема 2
...
# Ключевые выводы
 - ...
 - ...
Правила:
 - Пиши ясно, профессионально и лаконично.
 - Не добавляй информацию от себя.
 - Не используй вводные фразы вроде «в тексте говорится».
 - Сохраняй терминологию оригинала.
 - Если встречаются даты, числа, имена или технические термины — обязательно сохрани их.
 - Если текст написан на русском — отвечай на русском.
 - Если текст написан на английском — отвечай на английском.
 - Если текст смешанный — используй основной язык текста.
 - Если исходный текст недостаточно связный, содержит мусор от страницы, меню, рекламу или обрывки интерфейса — игнорируй нерелевантные элементы и пересказывай только содержательную часть.
 - Не упоминай, что ты являешься AI-ассистентом.
 - Не добавляй разделы, если для них нет информации в исходном тексте.
 - Если в тексте нет рекомендаций или выводов, сформулируй только выводы, которые прямо следуют из текста.`

type summaryRequest struct {
	Text string `json:"text"`
}

type summaryResponse struct {
	SummaryMD string `json:"summary_md"`
}

// NewSummaryHandler handles POST /summary.
// Accepts {"text": "..."}, returns {"summary_md": "..."}.
// text: min 300 chars, max 20000 chars.
func NewSummaryHandler(gc *service.GeminiClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
			return
		}

		if r.Body == nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
			return
		}

		var req summaryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
			return
		}

		text := strings.TrimSpace(req.Text)
		if text == "" {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "text is required"})
			return
		}

		if len(text) > summaryMaxLen {
			writeJSON(w, http.StatusRequestEntityTooLarge, ErrorResponse{Error: "text is too large"})
			return
		}

		if len(text) < summaryMinLen {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "message is too short"})
			return
		}

		result, err := gc.Generate(summarySystemPrompt, text)
		if err != nil {
			writeAIError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, summaryResponse{SummaryMD: result})
	}
}

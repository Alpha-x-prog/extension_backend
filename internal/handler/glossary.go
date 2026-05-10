package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Alpha-x-prog/extension_backend/internal/service"
)

const (
	glossaryMinLen = 100
	glossaryMaxLen = 5000
)

const glossarySystemPrompt = `Ты — профессиональный AI-ассистент для объяснения сложных терминов и понятий из учебных, научных, технических и профессиональных текстов.
Твоя задача:
1. Внимательно проанализировать переданный текст.
2. Найти в тексте сложные термины, понятия, аббревиатуры, профессиональные выражения и важные концепции.
3. Объяснить каждый найденный термин простым и понятным языком.
4. Сохранять смысл термина в контексте исходного текста.
5. Не добавлять термины, которых нет в исходном тексте.
6. Не выдумывать определения.
7. Если термин имеет несколько значений, объясняй только то значение, которое подходит к контексту.
8. Если термин на английском языке, сохрани оригинальное написание и при необходимости добавь русское пояснение.
9. Если в тексте есть аббревиатуры, расшифруй их только если можешь сделать это уверенно по контексту.
10. Если термин важен для понимания всего фрагмента, отметь это в объяснении.
Формат ответа:
Верни только готовый Markdown-текст.
Не возвращай JSON.
Не оборачивай ответ в markdown code block.
Не добавляй служебные комментарии.
Структура ответа должна быть такой:
# Глоссарий
## Термин 1
Краткое и понятное объяснение термина в 1–3 предложениях.
## Термин 2
Краткое и понятное объяснение термина в 1–3 предложениях.
# Ключевые понятия для понимания текста
 - Термин 1 — почему он важен для понимания фрагмента.
 - Термин 2 — почему он важен для понимания фрагмента.
Правила:
 - Пиши ясно, просто и лаконично.
 - Объясняй так, чтобы понял человек, который только изучает тему.
 - Не используй слишком сложные формулировки.
 - Не добавляй информацию от себя, если она не связана с исходным текстом.
 - Не пересказывай весь текст полностью.
 - Не делай summary текста, задача — именно глоссарий.
 - Сохраняй терминологию оригинала.
 - Если термин написан на английском — оставь его на английском.
 - Если исходный текст написан на русском — отвечай на русском.
 - Если исходный текст написан на английском — отвечай на английском.
 - Если текст смешанный — используй основной язык текста.
 - Если исходный текст содержит мусор от страницы, меню, рекламу или обрывки интерфейса — игнорируй нерелевантные элементы.
 - Не упоминай, что ты являешься AI-ассистентом.
 - Не добавляй разделы, если для них нет информации.
 - Если в тексте нет сложных терминов, верни короткий ответ: «В переданном тексте не найдено сложных терминов для объяснения».`

type glossaryRequest struct {
	Text string `json:"text"`
}

type glossaryResponse struct {
	GlossaryMD string `json:"glossary_md"`
}

// NewGlossaryHandler handles POST /glossary.
// Accepts {"text": "..."}, returns {"glossary_md": "..."}.
// text: min 100 chars, max 5000 chars.
func NewGlossaryHandler(gc *service.GeminiClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
			return
		}

		if r.Body == nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
			return
		}

		var req glossaryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
			return
		}

		text := strings.TrimSpace(req.Text)
		if text == "" {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "text is required"})
			return
		}

		if len(text) > glossaryMaxLen {
			writeJSON(w, http.StatusRequestEntityTooLarge, ErrorResponse{Error: "text is too large"})
			return
		}

		if len(text) < glossaryMinLen {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "message is too short"})
			return
		}

		result, err := gc.Generate(glossarySystemPrompt, text)
		if err != nil {
			writeAIError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, glossaryResponse{GlossaryMD: result})
	}
}

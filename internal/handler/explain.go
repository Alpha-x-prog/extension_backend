package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const maxTextLength = 5000

const geminiModel = "gemini-2.0-flash"

const systemPrompt = `**Роль ассистента**
Ты — AI-помощник для изучения сложных материалов в браузере.
Твоя задача — объяснять выделенный пользователем фрагмент текста простым, понятным и точным языком.

**Контекст продукта**
Пользователь читает статью, документацию, учебный материал или профессиональный текст и выделяет непонятный фрагмент.
Твоя задача — быстро объяснить именно этот фрагмент, не заставляя пользователя копировать текст, формулировать длинный запрос или открывать отдельный чат.

**Главная цель**
Помоги пользователю понять смысл выделенного текста, ключевые термины и практическую пользу фрагмента.

**Формат ответа**

1. **Коротко**
    - 2–4 предложения.
    - Объясни главную мысль выделенного фрагмента простым языком.

2. **Подробное объяснение**
    - 3–5 коротких абзацев.
    - Каждый абзац не длиннее 100–120 слов.
    - Начни с сути, затем объясни детали, причины и связи.
    - Не уходи далеко от выделенного текста.

3. **Пример или аналогия**
    - Если тема сложная, приведи простой пример из жизни, обучения или IT.
    - Если пример может исказить смысл, лучше не используй аналогию.

4. **Глоссарий**
    - Выпиши ключевые термины из выделенного текста.
    - Для каждого термина дай короткое объяснение простыми словами.
    - Не добавляй случайные термины, которых нет в тексте, если они не нужны для понимания.

5. **Что важно запомнить**
    - 2–4 пункта с главными выводами.

6. **Источники и проверка**
    - Если тема техническая, научная, юридическая, медицинская или быстро меняющаяся, укажи, что информацию желательно проверить по первоисточнику.
    - Если возможно, предложи тип источника: официальная документация, научная статья, стандарт, учебник.
    - Не выдумывай ссылки и источники.

**Язык и стиль**
- Пиши по-русски.
- Используй дружелюбный, но деловой тон.
- Избегай сложных слов без объяснения.
- Не используй чрезмерный жаргон.
- Если термин важен, сначала объясни его простыми словами.

**Ограничения**
- Объясняй именно выделенный фрагмент.
- Не пересказывай всю тему целиком, если это не нужно.
- Не добавляй неподтверждённые факты.
- Если добавляешь внешний контекст, явно пометь его как «дополнительный контекст».
- Если исходный текст неоднозначный, скажи: «Это можно понять так...» и предложи наиболее вероятную интерпретацию.
- Если информации недостаточно, честно напиши: «По выделенному фрагменту нельзя точно сказать...»

**Самокоррекция**
- Если в ответе обнаружена ошибка, исправь её и кратко поясни, что изменилось.

**Критерии успешного ответа**
- Пользователь понял смысл выделенного текста.
- Ответ короткий, но не поверхностный.
- Термины объяснены простым языком.
- Нет лишней информации и выдуманных фактов.
- Ответ помогает продолжить чтение без потери фокуса.`

// Gemini REST API request/response types

type geminiPart struct {
	Text string `json:"text"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiRequest struct {
	SystemInstruction geminiContent   `json:"system_instruction"`
	Contents          []geminiContent `json:"contents"`
}

type geminiCandidate struct {
	Content geminiContent `json:"content"`
}

type geminiResponse struct {
	Candidates []geminiCandidate `json:"candidates"`
}

// HTTP handler types

type ExplainRequest struct {
	Text string `json:"text"`
}

type ExplainResponse struct {
	Answer string `json:"answer"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func NewExplainHandler(apiKey string) http.HandlerFunc {
	apiURL := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		geminiModel, apiKey,
	)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
			return
		}

		if r.Body == nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
			return
		}

		var req ExplainRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
			return
		}

		if len(req.Text) == 0 {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "text is required"})
			return
		}

		if len(req.Text) > maxTextLength {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "text exceeds maximum length of 5000 characters"})
			return
		}

		answer, err := callGemini(apiURL, req.Text)
		if err != nil {
			log.Printf("Gemini API error: %v", err)
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to get explanation"})
			return
		}

		writeJSON(w, http.StatusOK, ExplainResponse{Answer: answer})
	}
}

func callGemini(apiURL, text string) (string, error) {
	body := geminiRequest{
		SystemInstruction: geminiContent{
			Parts: []geminiPart{{Text: systemPrompt}},
		},
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: text}}},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini returned %d: %s", resp.StatusCode, string(raw))
	}

	var gemResp geminiResponse
	if err := json.Unmarshal(raw, &gemResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(gemResp.Candidates) == 0 || len(gemResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from gemini")
	}

	return gemResp.Candidates[0].Content.Parts[0].Text, nil
}

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"wordbook/models"
	"wordbook/utils"

	"github.com/gin-gonic/gin"
)

type AIResponse struct {
	Definition string   `json:"definition"`
	Sentences  []string `json:"sentences"`
}

func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}

	var exist models.User
	if models.DB.Where("username = ?", req.Username).First(&exist).Error == nil {
		c.JSON(400, gin.H{"error": "用户名已存在"})
		return
	}

	user := models.User{Username: req.Username}
	user.SetPassword(req.Password)
	models.DB.Create(&user)

	token, _ := utils.GenerateToken(user.ID)
	c.JSON(200, gin.H{"token": token, "user": user})
}

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	c.BindJSON(&req)

	var user models.User
	if models.DB.Where("username = ?", req.Username).First(&user).Error != nil {
		c.JSON(401, gin.H{"error": "用户不存在"})
		return
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(401, gin.H{"error": "密码错误"})
		return
	}

	token, _ := utils.GenerateToken(user.ID)
	c.JSON(200, gin.H{"token": token, "user": user})
}

func QueryWord(c *gin.Context) {
	word := c.Query("word")
	provider := c.Query("ai_provider")
	userID := c.GetUint("userID")

	var exist models.Word
	if models.DB.Where("user_id = ? AND word = ?", userID, word).First(&exist).Error == nil {
		var sentences []string
		json.Unmarshal([]byte(exist.Sentences), &sentences)
		c.JSON(200, gin.H{
			"from_cache":  true,
			"word":        exist.Word,
			"definition":  exist.Definition,
			"sentences":   sentences,
			"ai_provider": exist.AIProvider,
		})
		return
	}

	result, err := callAI(word, provider)
	if err != nil {
		c.JSON(500, gin.H{"error": "AI服务错误: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"from_cache":  false,
		"word":        word,
		"definition":  result.Definition,
		"sentences":   result.Sentences,
		"ai_provider": provider,
	})
}

func callAI(word, provider string) (*AIResponse, error) {
	if provider == "deepseek" {
		return callDeepSeek(word)
	}
	return callQwen(word)
}

func callQwen(word string) (*AIResponse, error) {
	apiKey := os.Getenv("QWEN_API_KEY")
	if apiKey == "" {
		return &AIResponse{
			Definition: "请配置 QWEN_API_KEY 环境变量",
			Sentences:  []string{"在 .env 文件中添加 QWEN_API_KEY=你的密钥", "然后运行 docker-compose down && docker-compose up -d --build"},
		}, nil
	}

	prompt := fmt.Sprintf(`请为单词 "%s" 生成英文学习内容。必须返回严格的JSON格式，不要包含任何其他文字：
{
    "definition": "中文释义（简洁准确）",
    "sentences": ["英文例句1", "英文例句2", "英文例句3"]
}`, word)

	reqBody := map[string]interface{}{
		"model": "qwen-turbo",
		"input": map[string]interface{}{
			"messages": []map[string]string{
				{"role": "user", "content": prompt},
			},
		},
		"parameters": map[string]interface{}{
			"result_format": "message",
		},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Output struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		} `json:"output"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if len(result.Output.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	var aiResp AIResponse
	content := result.Output.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &aiResp); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %v", err)
	}

	return &aiResp, nil
}

func callDeepSeek(word string) (*AIResponse, error) {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return &AIResponse{
			Definition: "请配置 DEEPSEEK_API_KEY 环境变量",
			Sentences:  []string{"在 .env 文件中添加 DEEPSEEK_API_KEY=你的密钥", "然后重启服务"},
		}, nil
	}

	prompt := fmt.Sprintf(`请为单词 "%s" 生成英文学习内容。返回JSON格式：{"definition": "中文释义", "sentences": ["例句1", "例句2", "例句3"]}`, word)

	reqBody := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	var aiResp AIResponse
	if err := json.Unmarshal([]byte(result.Choices[0].Message.Content), &aiResp); err != nil {
		return nil, err
	}

	return &aiResp, nil
}

func SaveWord(c *gin.Context) {
	var req struct {
		Word       string   `json:"word"`
		Definition string   `json:"definition"`
		Sentences  []string `json:"sentences"`
		AIProvider string   `json:"ai_provider"`
	}
	c.BindJSON(&req)

	userID := c.GetUint("userID")

	var exist models.Word
	if models.DB.Where("user_id = ? AND word = ?", userID, req.Word).First(&exist).Error == nil {
		c.JSON(409, gin.H{"error": "单词已保存"})
		return
	}

	sentencesJSON, _ := json.Marshal(req.Sentences)
	word := models.Word{
		UserID:     userID,
		Word:       req.Word,
		Definition: req.Definition,
		Sentences:  string(sentencesJSON),
		AIProvider: req.AIProvider,
	}
	models.DB.Create(&word)

	c.JSON(200, gin.H{"message": "保存成功", "id": word.ID})
}

func GetWords(c *gin.Context) {
	userID := c.GetUint("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var words []models.Word
	var total int64

	query := models.DB.Model(&models.Word{}).Where("user_id = ?", userID)
	query.Count(&total)
	query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&words)

	var result []gin.H
	for _, w := range words {
		var sentences []string
		json.Unmarshal([]byte(w.Sentences), &sentences)
		result = append(result, gin.H{
			"id":          w.ID,
			"word":        w.Word,
			"definition":  w.Definition,
			"sentences":   sentences,
			"ai_provider": w.AIProvider,
			"created_at":  w.CreatedAt.Format("2006-01-02"),
		})
	}

	c.JSON(200, gin.H{
		"total":     total,
		"page":      page,
		"page_size": size,
		"words":     result,
	})
}

func DeleteWord(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("userID")

	result := models.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Word{})
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "单词不存在"})
		return
	}
	c.JSON(200, gin.H{"message": "删除成功"})
}

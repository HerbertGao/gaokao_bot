package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/herbertgao/gaokao_bot/internal/model"
	"github.com/herbertgao/gaokao_bot/internal/service"
	"github.com/herbertgao/gaokao_bot/internal/util"
)

const (
	// MaxTemplatesPerUser 每个用户最多可创建的模板数量
	// 限制为 10 个以防止滥用和控制数据库存储
	MaxTemplatesPerUser = 10

	// MaxTemplateContentLength 模板内容最大长度（字符数）
	// 限制为 140 字符以确保：
	// 1. 适合 Telegram 消息显示
	// 2. 防止过长文本影响用户体验
	// 3. 控制数据库字段大小
	MaxTemplateContentLength = 140

	// MaxTemplateNameLength 模板名称最大长度（字符数）
	// 限制为 20 字符以保持名称简洁易读
	MaxTemplateNameLength = 20
)

// TemplateHandler 模板处理器
type TemplateHandler struct {
	templateService *service.UserTemplateService
}

// NewTemplateHandler 创建模板处理器
func NewTemplateHandler(templateService *service.UserTemplateService) *TemplateHandler {
	return &TemplateHandler{
		templateService: templateService,
	}
}

// CreateTemplateRequest 创建模板请求
type CreateTemplateRequest struct {
	TemplateName    string `json:"template_name"`
	TemplateContent string `json:"template_content" binding:"required"`
}

// UpdateTemplateRequest 更新模板请求
type UpdateTemplateRequest struct {
	TemplateName    string `json:"template_name"`
	TemplateContent string `json:"template_content" binding:"required"`
}

// GetTemplates 获取模板列表
func (h *TemplateHandler) GetTemplates(c *gin.Context) {
	userID := c.GetInt64("user_id")

	templates, err := h.templateService.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to get templates: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    templates,
	})
}

// CreateTemplate 创建模板
func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	// 检查用户模板数量限制
	existingTemplates, err := h.templateService.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to check template count: %v", err),
		})
		return
	}

	if len(existingTemplates) >= MaxTemplatesPerUser {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   fmt.Sprintf("模板数量已达上限（最多 %d 个）", MaxTemplatesPerUser),
		})
		return
	}

	// 验证模板内容
	if err := validateTemplateContent(req.TemplateContent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// 验证模板名称
	if req.TemplateName != "" {
		if err := validateTemplateName(req.TemplateName); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}

	// 生成 ID
	id := util.GenerateID()

	template := &model.UserTemplate{
		ID:              id,
		UserID:          userID,
		TemplateName:    req.TemplateName,
		TemplateContent: req.TemplateContent,
	}

	if err := h.templateService.Create(template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to create template: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    template,
	})
}

// UpdateTemplate 更新模板
func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid template ID",
		})
		return
	}

	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	// 验证模板内容
	if err := validateTemplateContent(req.TemplateContent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// 验证模板名称
	if req.TemplateName != "" {
		if err := validateTemplateName(req.TemplateName); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}

	// 检查模板是否存在且属于当前用户
	existingTemplate, err := h.templateService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to get template: %v", err),
		})
		return
	}

	if existingTemplate == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Template not found",
		})
		return
	}

	if existingTemplate.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Permission denied",
		})
		return
	}

	// 更新模板
	existingTemplate.TemplateName = req.TemplateName
	existingTemplate.TemplateContent = req.TemplateContent

	if err := h.templateService.Update(existingTemplate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to update template: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    existingTemplate,
	})
}

// DeleteTemplate 删除模板
func (h *TemplateHandler) DeleteTemplate(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid template ID",
		})
		return
	}

	// 检查模板是否存在且属于当前用户
	template, err := h.templateService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to get template: %v", err),
		})
		return
	}

	if template == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Template not found",
		})
		return
	}

	if template.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Permission denied",
		})
		return
	}

	if err := h.templateService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to delete template: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// validateTemplateContent 验证模板内容
func validateTemplateContent(content string) error {
	if content == "" {
		return fmt.Errorf("模板内容不能为空")
	}

	if len(content) > MaxTemplateContentLength {
		return fmt.Errorf("模板内容不能超过 %d 字符", MaxTemplateContentLength)
	}

	// 必须包含 {exam} 和 {time}
	if !strings.Contains(content, "{exam}") {
		return fmt.Errorf("模板必须包含 {exam} 变量")
	}

	if !strings.Contains(content, "{time}") {
		return fmt.Errorf("模板必须包含 {time} 变量")
	}

	return nil
}

// validateTemplateName 验证模板名称
func validateTemplateName(name string) error {
	if len(name) > MaxTemplateNameLength {
		return fmt.Errorf("模板标题不能超过 %d 字符", MaxTemplateNameLength)
	}
	return nil
}

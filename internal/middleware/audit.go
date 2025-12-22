package middleware

import (
	"context"

	"github.com/Mag1cFall/magtrade/internal/model"
	"github.com/Mag1cFall/magtrade/internal/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuditLogger struct {
	repo *repository.AuditLogRepository
	log  *zap.Logger
}

func NewAuditLogger(log *zap.Logger) *AuditLogger {
	return &AuditLogger{
		repo: repository.NewAuditLogRepository(),
		log:  log,
	}
}

func (a *AuditLogger) Log(ctx context.Context, userID int64, action, resource, resourceID, ip, userAgent, details string) {
	auditLog := &model.AuditLog{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		IP:         ip,
		UserAgent:  userAgent,
		Details:    details,
	}

	if err := a.repo.Create(ctx, auditLog); err != nil {
		a.log.Error("failed to create audit log", zap.Error(err))
	}
}

func (a *AuditLogger) LogFromGin(c *gin.Context, userID int64, action, resource, resourceID, details string) {
	a.Log(c.Request.Context(), userID, action, resource, resourceID, c.ClientIP(), c.Request.UserAgent(), details)
}

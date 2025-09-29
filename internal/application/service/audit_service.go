package service

import (
	"fms_audit/internal/domain"
)

// AuditService định nghĩa interface cho business logic
type AuditService interface {
	Save(log domain.AuditLog) error
}

// auditServiceImpl là implement mặc định
type auditServiceImpl struct{}

// NewAuditService khởi tạo service - THÊM THAM SỐ NHƯNG KHÔNG DÙNG
func NewAuditService(repo interface{}) AuditService {
	// Tạm thời bỏ qua repo, không dùng đến
	return &auditServiceImpl{}
}

// Save thực hiện lưu log (hiện demo in ra console)
func (s *auditServiceImpl) Save(log domain.AuditLog) error {
	// TODO: Sau này sẽ lưu vào database
	// Tạm thời in ra console để test
	println("✅ Saved Audit Log:", log.ID, log.Action, log.UserID)
	return nil
}

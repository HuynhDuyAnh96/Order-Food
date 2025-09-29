package repository

import (
	"fms_audit/internal/domain"
	"fmt"
)

type AuditRepository struct{}

func NewAuditRepository() *AuditRepository {
	return &AuditRepository{}
}

func (r *AuditRepository) Save(log domain.AuditLog) error {
	// TODO: connect DB (Mongo, Postgres,...)
	fmt.Printf("Saving audit log: %+v\n", log)
	return nil
}

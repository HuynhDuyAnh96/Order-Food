package grpcinfra

import (
	"context"
	"time"

	"fms_audit/internal/application/service"
	"fms_audit/internal/domain"
	"fms_audit/internal/infrastructure/grpc/proto"
)

type AuditServer struct {
	proto.UnimplementedAuditServiceServer
	svc service.AuditService
}

func NewAuditServer(svc service.AuditService) *AuditServer {
	return &AuditServer{svc: svc}
}

func (s *AuditServer) SaveAudit(ctx context.Context, req *proto.SaveAuditRequest) (*proto.SaveAuditResponse, error) {
	if req == nil {
		return &proto.SaveAuditResponse{Success: false}, nil
	}

	var auditLog domain.AuditLog

	// Cách 1: Nếu có log wrapper
	if req.GetLog() != nil {
		logProto := req.GetLog()
		auditLog = domain.AuditLog{
			ID:        logProto.GetId(),
			Action:    logProto.GetAction(),
			UserID:    logProto.GetUserId(),
			CreatedAt: time.Unix(logProto.GetCreatedAt(), 0),
			Metadata:  logProto.GetMetadata(),
		}
	} else {
		// Cách 2: Fallback - in ra để debug
		println("⚠️  Received request without log wrapper")
		return &proto.SaveAuditResponse{Success: false}, nil
	}

	err := s.svc.Save(auditLog)
	return &proto.SaveAuditResponse{Success: err == nil}, nil
}

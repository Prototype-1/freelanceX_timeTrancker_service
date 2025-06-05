package service

import (
	"context"
	"time"
"github.com/Prototype-1/freelanceX_timeTrancker_service/internal/model"
"github.com/Prototype-1/freelanceX_timeTrancker_service/internal/repository"
pb "github.com/Prototype-1/freelanceX_timeTrancker_service/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type TimeLogService interface {
	CreateTimeLog(ctx context.Context, req *pb.CreateTimeLogRequest) (*pb.CreateTimeLogResponse, error)
	GetTimeLogsByUser(ctx context.Context, req *pb.GetTimeLogsByUserRequest) (*pb.TimeLogsResponse, error)
	GetTimeLogsByProject(ctx context.Context, req *pb.GetTimeLogsByProjectRequest) (*pb.TimeLogsResponse, error)
	UpdateTimeLog(ctx context.Context, req *pb.UpdateTimeLogRequest) (*pb.UpdateTimeLogResponse, error)
	DeleteTimeLog(ctx context.Context, req *pb.DeleteTimeLogRequest) (*pb.DeleteTimeLogResponse, error)
}

type timeLogService struct {
	repo repository.TimeLogRepository
	pb.UnimplementedTimeLogServiceServer  
}

func NewTimeLogService(repo repository.TimeLogRepository) pb.TimeLogServiceServer {
	return &timeLogService{repo: repo}
}

func extractRole(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	roles := md.Get("role")
	if len(roles) == 0 {
		return ""
	}
	return roles[0]
}

func extractUserID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	ids := md.Get("user_id")
	if len(ids) == 0 {
		return ""
	}
	return ids[0]
}

func (s *timeLogService) CreateTimeLog(ctx context.Context, req *pb.CreateTimeLogRequest) (*pb.CreateTimeLogResponse, error) {
	if extractRole(ctx) != "freelancer" {
		return nil, status.Error(codes.PermissionDenied, "only freelancers can create time logs")
	}
	userID := extractUserID(ctx)

	start := req.StartTime.AsTime()
	end := req.EndTime.AsTime()
	if end.Before(start) {
		return nil, status.Error(codes.InvalidArgument, "end time cannot be before start time")
	}

	durationMinutes := int(end.Sub(start).Minutes())

	log := &model.TimeLog{
		ID:        uuid.New(),
		UserID:    uuid.MustParse(userID),
		ProjectID: uuid.MustParse(req.ProjectId),
		TaskName:  req.TaskName,
		StartTime: req.StartTime.AsTime(),
		EndTime:   req.EndTime.AsTime(),
		Duration:  durationMinutes,
		Source:    req.Source.String(),
	}

	createdLog, err := s.repo.CreateTimeLog(log)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create time log: %v", err)
	}

	return &pb.CreateTimeLogResponse{
		LogId:         createdLog.ID.String(),
		DurationHours: float64(createdLog.Duration) / 60,
	}, nil
}

func (s *timeLogService) GetTimeLogsByUser(ctx context.Context, req *pb.GetTimeLogsByUserRequest) (*pb.TimeLogsResponse, error) {
	role := extractRole(ctx)
	requesterID := extractUserID(ctx)

	if role != "freelancer" && role != "admin" {
		return nil, status.Error(codes.PermissionDenied, "access denied")
	}

	if role == "freelancer" && req.UserId != requesterID {
		return nil, status.Error(codes.PermissionDenied, "cannot view others' logs")
	}

	var fromTime, toTime *time.Time
if req.DateFrom != nil {
    t := req.DateFrom.AsTime()
    fromTime = &t
}
if req.DateTo != nil {
    t := req.DateTo.AsTime()
    toTime = &t
}

userID, err := uuid.Parse(req.UserId)
if err != nil {
    return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
}

projectID, err := uuid.Parse(req.ProjectId)
if err != nil {
    return nil, status.Errorf(codes.InvalidArgument, "invalid project ID: %v", err)
}

logs, err := s.repo.GetTimeLogsByUser(userID, projectID, fromTime, toTime)
if err != nil {
    return nil, status.Errorf(codes.Internal, "error retrieving logs: %v", err)
}
	return &pb.TimeLogsResponse{Logs: toResponseLogs(logs)}, nil
}

func (s *timeLogService) GetTimeLogsByProject(ctx context.Context, req *pb.GetTimeLogsByProjectRequest) (*pb.TimeLogsResponse, error) {
	role := extractRole(ctx)
	if role != "freelancer" && role != "admin" {
		return nil, status.Error(codes.PermissionDenied, "access denied")
	}

	logs, err := s.repo.GetTimeLogsByProject(uuid.MustParse(req.ProjectId), nil, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error retrieving logs: %v", err)
	}

return &pb.TimeLogsResponse{Logs: toResponseLogs(logs)}, nil

}

func (s *timeLogService) UpdateTimeLog(ctx context.Context, req *pb.UpdateTimeLogRequest) (*pb.UpdateTimeLogResponse, error) {
	requesterID := extractUserID(ctx)
	role := extractRole(ctx)

	logID, err := uuid.Parse(req.LogId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid log ID")
	}

	log, err := s.repo.GetTimeLogByID(logID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "time log not found")
	}

	if role != "admin" && log.UserID.String() != requesterID {
		return nil, status.Error(codes.PermissionDenied, "not allowed to update this time log")
	}
	if log.Source != model.MANUAL {
		return nil, status.Error(codes.FailedPrecondition, "only manual time logs can be updated")
	}

	log.StartTime = req.StartTime.AsTime()
	log.EndTime = req.EndTime.AsTime()

	updated, err := s.repo.UpdateTimeLog(log)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update failed: %v", err)
	}

	return &pb.UpdateTimeLogResponse{
		LogId:         updated.ID.String(),
		DurationHours: float64(updated.Duration) / 60,
	}, nil
}

func (s *timeLogService) DeleteTimeLog(ctx context.Context, req *pb.DeleteTimeLogRequest) (*pb.DeleteTimeLogResponse, error) {
	requesterID := extractUserID(ctx)
	role := extractRole(ctx)

	logID, err := uuid.Parse(req.LogId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid log ID")
	}

	log, err := s.repo.GetTimeLogByID(logID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "log not found")
	}

	if role != "admin" && log.UserID.String() != requesterID {
		return nil, status.Error(codes.PermissionDenied, "not allowed to delete this log")
	}

	err = s.repo.DeleteTimeLog(logID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "deletion failed: %v", err)
	}

	return &pb.DeleteTimeLogResponse{
    LogId:  log.ID.String(),
    Status: "Time log deleted successfully",
}, nil

}

func toResponseLogs(logs []model.TimeLog) []*pb.TimeLog {
	var res []*pb.TimeLog
	for _, log := range logs {
		res = append(res, &pb.TimeLog{
			LogId:     log.ID.String(),
			UserId:    log.UserID.String(),
			ProjectId: log.ProjectID.String(),
			TaskName:  log.TaskName,
			StartTime: timestamppb.New(log.StartTime),
			EndTime:   timestamppb.New(log.EndTime),
			Duration:  int32(log.Duration), 
			Source:    pb.TimeLogSource(pb.TimeLogSource_value[log.Source]),
		})
	}
	return res
}


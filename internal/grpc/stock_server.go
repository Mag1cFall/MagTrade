package grpc

import (
	"context"
	"net"

	"github.com/Mag1cFall/magtrade/internal/cache"
	pb "github.com/Mag1cFall/magtrade/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type StockServer struct {
	pb.UnimplementedStockServiceServer
	stockService *cache.StockService
	log          *zap.Logger
}

func NewStockServer(log *zap.Logger) *StockServer {
	return &StockServer{
		stockService: cache.NewStockService(),
		log:          log,
	}
}

func (s *StockServer) GetStock(ctx context.Context, req *pb.GetStockRequest) (*pb.GetStockResponse, error) {
	stock, err := s.stockService.GetStock(ctx, req.FlashSaleId)
	if err != nil {
		return &pb.GetStockResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.GetStockResponse{
		Stock:   int32(stock),
		Success: true,
		Message: "success",
	}, nil
}

func (s *StockServer) DeductStock(ctx context.Context, req *pb.DeductStockRequest) (*pb.DeductStockResponse, error) {
	result, err := s.stockService.DeductStock(ctx, req.FlashSaleId, req.UserId, int(req.Quantity), int(req.Limit))
	if err != nil {
		return &pb.DeductStockResponse{
			Success: false,
			Code:    -99,
			Message: err.Error(),
		}, nil
	}

	return &pb.DeductStockResponse{
		Success: result.Success,
		Code:    int32(result.Code),
		Message: result.Message,
	}, nil
}

func (s *StockServer) RestoreStock(ctx context.Context, req *pb.RestoreStockRequest) (*pb.RestoreStockResponse, error) {
	err := s.stockService.RestoreStock(ctx, req.FlashSaleId, req.UserId, int(req.Quantity))
	if err != nil {
		return &pb.RestoreStockResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.RestoreStockResponse{
		Success: true,
		Message: "success",
	}, nil
}

func (s *StockServer) InitStock(ctx context.Context, req *pb.InitStockRequest) (*pb.InitStockResponse, error) {
	err := s.stockService.InitStock(ctx, req.FlashSaleId, int(req.Stock))
	if err != nil {
		return &pb.InitStockResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.InitStockResponse{
		Success: true,
		Message: "success",
	}, nil
}

func StartGRPCServer(addr string, log *zap.Logger) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer()
	stockServer := NewStockServer(log)
	pb.RegisterStockServiceServer(server, stockServer)
	reflection.Register(server)

	log.Info("gRPC server starting", zap.String("addr", addr))

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Error("gRPC server failed", zap.Error(err))
		}
	}()

	return server, nil
}

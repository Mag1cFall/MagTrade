package grpc

import (
	"context"
	"time"

	pb "github.com/Mag1cFall/magtrade/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StockClient struct {
	conn   *grpc.ClientConn
	client pb.StockServiceClient
	log    *zap.Logger
}

func NewStockClient(addr string, log *zap.Logger) (*StockClient, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	log.Info("gRPC client connected", zap.String("addr", addr))

	return &StockClient{
		conn:   conn,
		client: pb.NewStockServiceClient(conn),
		log:    log,
	}, nil
}

func (c *StockClient) GetStock(ctx context.Context, flashSaleID int64) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.GetStock(ctx, &pb.GetStockRequest{
		FlashSaleId: flashSaleID,
	})
	if err != nil {
		return 0, err
	}

	return int(resp.Stock), nil
}

type DeductResult struct {
	Success bool
	Code    int
	Message string
}

func (c *StockClient) DeductStock(ctx context.Context, flashSaleID, userID int64, quantity, limit int) (*DeductResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.DeductStock(ctx, &pb.DeductStockRequest{
		FlashSaleId: flashSaleID,
		UserId:      userID,
		Quantity:    int32(quantity),
		Limit:       int32(limit),
	})
	if err != nil {
		return nil, err
	}

	return &DeductResult{
		Success: resp.Success,
		Code:    int(resp.Code),
		Message: resp.Message,
	}, nil
}

func (c *StockClient) RestoreStock(ctx context.Context, flashSaleID, userID int64, quantity int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := c.client.RestoreStock(ctx, &pb.RestoreStockRequest{
		FlashSaleId: flashSaleID,
		UserId:      userID,
		Quantity:    int32(quantity),
	})
	return err
}

func (c *StockClient) InitStock(ctx context.Context, flashSaleID int64, stock int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := c.client.InitStock(ctx, &pb.InitStockRequest{
		FlashSaleId: flashSaleID,
		Stock:       int32(stock),
	})
	return err
}

func (c *StockClient) Close() error {
	return c.conn.Close()
}

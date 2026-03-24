package payment

import (
    "context"
    "fmt"

    model "github.com/Loboo34/travel/models"
)

type PaymentResult struct {
    Reference string 
    Status    model.PaymentStatus
}

type Provider interface {
    Charge(ctx context.Context, req ChargeRequest) (*PaymentResult, error)
    Refund(ctx context.Context, reference string, amount int64) error
}

type ChargeRequest struct {
    Amount   int64
    Currency string
    Method   model.PaymentMethod
    UserID   string
    Metadata map[string]string 
}


type StubProvider struct{}

func NewStubProvider() *StubProvider {
    return &StubProvider{}
}

func (p *StubProvider) Charge(ctx context.Context, req ChargeRequest) (*PaymentResult, error) {
    
    return &PaymentResult{
        Reference: fmt.Sprintf("stub_%d", req.Amount),
        Status:    model.PaymentStatus("paid"),
    }, nil
}

func (p *StubProvider) Refund(ctx context.Context, reference string, amount int64) error {
    return nil
}
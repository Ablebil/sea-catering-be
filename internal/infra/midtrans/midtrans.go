package midtrans

import (
	"time"

	conf "github.com/Ablebil/sea-catering-be/config"
	"github.com/Ablebil/sea-catering-be/internal/domain/dto"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransItf interface {
	CreateTransaction(req *dto.MidtransRequest) (*dto.PaymentResponse, error)
}

type Midtrans struct {
	snapClient *snap.Client
	conf       *conf.Config
}

func NewMidtrans(conf *conf.Config) MidtransItf {
	var snapClient snap.Client
	snapClient.New(conf.MidtransServerKey, midtrans.Sandbox)

	return &Midtrans{
		snapClient: &snapClient,
		conf:       conf,
	}
}

func (m *Midtrans) CreateTransaction(req *dto.MidtransRequest) (*dto.PaymentResponse, error) {
	var items []midtrans.ItemDetails
	for _, item := range req.ItemDetails {
		items = append(items, midtrans.ItemDetails{
			ID:    item.ID,
			Price: item.Price,
			Qty:   item.Qty,
			Name:  item.Name,
		})
	}

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.OrderID,
			GrossAmt: req.Amount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: req.CustomerDetails.Name,
			Email: req.CustomerDetails.Email,
			Phone: req.CustomerDetails.Phone,
		},
		Items: &items,
		Expiry: &snap.ExpiryDetails{
			StartTime: time.Now().Format("2006-01-02 15:04:05 -0700"),
			Unit:      "minutes",
			Duration:  int64(m.conf.MidtransPaymentDuration.Minutes()),
		},
		CustomField1: req.SubscriptionID.String(),
	}

	snapResp, err := m.snapClient.CreateTransaction(snapReq)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		Token:       snapResp.Token,
		RedirectURL: snapResp.RedirectURL,
	}, nil
}

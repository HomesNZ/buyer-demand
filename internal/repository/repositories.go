package repository

import (
	"github.com/HomesNZ/buyer-demand/internal/repository/address"
	buyerDemand "github.com/HomesNZ/buyer-demand/internal/repository/buyer_demand"
)

type Repositories interface {
	Address() address.Repo
	BuyerDemand() buyerDemand.Repo
}

type repositories struct {
	address     address.Repo
	buyerDemand buyerDemand.Repo
}

func (r repositories) Address() address.Repo {
	return r.address
}

func (r repositories) BuyerDemand() buyerDemand.Repo {
	return r.buyerDemand
}

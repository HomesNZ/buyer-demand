package repository

import (
	"github.com/HomesNZ/buyer-demand/internal/repository/address"
	buyerDemand "github.com/HomesNZ/buyer-demand/internal/repository/buyer_demand"
	propertyClaim "github.com/HomesNZ/buyer-demand/internal/repository/property_claim"
)

type Repositories interface {
	Address() address.Repo
	BuyerDemand() buyerDemand.Repo
	PropertyClaim() propertyClaim.Repo
}

type repositories struct {
	address       address.Repo
	buyerDemand   buyerDemand.Repo
	propertyClaim propertyClaim.Repo
}

func (r repositories) Address() address.Repo {
	return r.address
}

func (r repositories) BuyerDemand() buyerDemand.Repo {
	return r.buyerDemand
}

func (r repositories) PropertyClaim() propertyClaim.Repo {
	return r.propertyClaim
}

package system_services

import (
	"github.com/oswaldom-code/api-template-gin/src/adapters/repository"
	"github.com/oswaldom-code/api-template-gin/src/application/system_services/ports"
)

type Health interface {
	TestDb() error
}
type healthImp struct {
	r ports.Store
}

func HealthService() (Health, error) {
	repo, err := repository.NewRepository()
	if err != nil {
		return nil, err
	}
	return &healthImp{r: repo}, nil
}

func (p *healthImp) TestDb() error {
	return p.r.TestDb()
}

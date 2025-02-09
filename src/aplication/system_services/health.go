package system_services

import (
	"github.com/oswaldom-code/api-template-gin/pkg/config"
	"github.com/oswaldom-code/api-template-gin/src/adapters/repository"
	"github.com/oswaldom-code/api-template-gin/src/aplication/system_services/ports"
)

type Health interface {
	TestDb() error
}
type healthImp struct {
	r ports.Repository
}

func HealthService() Health {
	repo, err := repository.NewRepository(config.GetDBConfig())
	if err != nil {
		panic(err) // TODO: handle error
	}
	return &healthImp{r: repo}
}

func (p *healthImp) TestDb() error {
	return p.r.TestDb()
}

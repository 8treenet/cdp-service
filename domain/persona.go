package domain

import (
	"time"

	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/domain/aggregate"
	"github.com/8treenet/cdp-service/domain/vo"
	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindService(func() *PersonaService {
			return &PersonaService{}
		})
		initiator.InjectController(func(ctx freedom.Context) (service *PersonaService) {
			initiator.FetchService(ctx, &service)
			return
		})
	})
}

// PersonaService .
type PersonaService struct {
	Worker              freedom.Worker
	PersonaFactory      *aggregate.PersonaFactory
	CustomerRepository  *repository.CustomerRepository
	BehaviourRepository *repository.BehaviourRepository
	PersonaRepository   *repository.PersonaRepository
}

// CreatePersona
func (service *PersonaService) CreatePersona(req vo.ReqCreatePersona) error {
	cmd := service.PersonaFactory.CreatePersonaCMD(req.Name, req.Title, req.DateRange, req.XmlData)
	err := cmd.Do()
	if err != nil {
		service.Worker.Logger().Errorf("CreatePersona req:%v err:%v", req, err)
		return err
	}
	return nil
}

// ExecuteDayJob 每日全量
func (service *PersonaService) ExecuteDayJob() {
	cmds, err := service.PersonaFactory.JobPersonaCmds()
	if err != nil {
		service.Worker.Logger().Errorf("PersonaService.ExecuteDayJob.JobPersonaCmds err:%v", err)
		return
	}

	var validCmds []*aggregate.PersonaJob
	for _, cmd := range cmds {
		if cmd.Deleted == 0 {
			validCmds = append(validCmds, cmd)
		}
	}

	if len(validCmds) == 0 {
		return
	}
	startId := 0
	size := 500
	for {
		userIds := []string{}
		var err error
		userIds, startId, err = service.CustomerRepository.GetRangeUserIds(startId, size)
		if err != nil {
			service.Worker.Logger().Errorf("PersonaService.ExecuteDayJob.GetRangeUserIds startId:%v, size:%v, err:%v", startId, size, err)
			break
		}
		if len(userIds) != 0 {
			service.job(userIds, validCmds)
		}
		if len(userIds) <= size {
			break
		}
	}
}

// ExecuteRefreshJob 小时增量
func (service *PersonaService) ExecuteRefreshJob() {
	cmds, err := service.PersonaFactory.JobPersonaCmds()
	if err != nil {
		service.Worker.Logger().Errorf("PersonaService.ExecuteDayJob.JobPersonaCmds err:%v", err)
		return
	}
	var validCmds []*aggregate.PersonaJob
	for _, cmd := range cmds {
		if cmd.Deleted == 0 {
			validCmds = append(validCmds, cmd)
			continue
		}

		hour := cmd.GetDeadHour()
		if hour > 0 && hour < 5*24 {
			service.PersonaRepository.DeleteCustomerProfileByPersona(&cmd.Persona)
		}
	}

	if len(validCmds) == 0 {
		return
	}

	size := 500
	for {
		userIds := service.BehaviourRepository.FetchActiveCustomer(size)
		if len(userIds) != 0 {
			service.job(userIds, validCmds)
			continue
		}
		break
	}
}

// job 处理画像
func (service *PersonaService) job(userIds []string, cmds []*aggregate.PersonaJob) {
	now := time.Now()
	for _, cmd := range cmds {
		err := cmd.Do(userIds, now)
		if err == nil {
			continue
		}
		service.Worker.Logger().Errorf("PersonaService.job.Do cmd:%v, err:%v", *cmd, err)
	}
}

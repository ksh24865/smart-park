package logicCoreUC

import (
	"github.com/ksh24865/smart-park/logic-core/adapter"
	"github.com/ksh24865/smart-park/logic-core/domain/repository"
	"github.com/ksh24865/smart-park/logic-core/domain/service"
)

type logicCoreUsecase struct {
	rr repository.RegistRepo
	ks service.KafkaConsumerGroup
	es service.ElasticClient
	ls service.LogicService
}

func NewLogicCoreUsecase(rr repository.RegistRepo,
	ks service.KafkaConsumerGroup,
	es service.ElasticClient,
	ls service.LogicService) *logicCoreUsecase {
	lcu := &logicCoreUsecase{
		rr: rr,
		ks: ks,
		es: es,
		ls: ls,
	}

	in := lcu.ks.GetOutput()
	out := lcu.es.GetInput()

	go func() {
		for rawData := range in {
			ld, err := lcu.ToLogicData(&rawData)
			// log.Println("add metadata :", ld)

			if err != nil {
				continue
			}

			lchs, err := lcu.ls.GetLogicChans(ld.SensorID)
			if err == nil {
				for _, ch := range lchs {
					if len(ch) != cap(ch) {
						ch <- ld
					}
				}
			}
			out <- lcu.ToDocument(&ld)
		}
	}()

	return lcu
}

func (lcu *logicCoreUsecase) AppendSinkAddr(sa *adapter.SinkAddr) error {
	lcu.rr.AppendSinkAddr(sa.Sid, &sa.Addr)

	return nil
}

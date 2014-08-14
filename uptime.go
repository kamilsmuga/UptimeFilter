package uptime

import (
	"github.com/mozilla-services/heka/pipeline"
)

type UptimeFilter struct {
	// epoch
	startHour uint32
	endHour   uint32
	// in seconds
	initialUptime uint32
	lastUptime    uint32
	totalUptime   uint32
	// true if uptime dropped
	dived  bool
	output string
}

func (f *UptimeFilter) Init(config interface{}) error {

}

func (f *UptimeFilter) Run(runner FilterRunner, helper PluginHelper) (err error) {
	var (
		pack    *PipelinePack
		output  OutputRunner
		payload string
		ok      bool
	)
	if output, ok = helper.Output(f.output); !ok {
		runner.LogError("No output: ")
		return
	}
	inChan := runner.InChan()
	for pack = range inChan {
		payload = pack.Message.GetPayload()
	}
	runner.LogMessage(payload)
	pack.Recycle()
	return
}

func init() {
	RegisterPlugin("UptimeFilter", func() interface{} {
		return new(UptimeFilter)
	})
}

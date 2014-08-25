package uptime

import (
	"github.com/mozilla-services/heka/pipeline"
	"log"
	"strconv"
	"strings"
	"time"
)

type UptimeFilter struct {
	// epoch
	startHour int64
	endHour   int64
	// in seconds
	initialUptime int64
	lastUptime    int64
	totalUptime   int64
	// true if uptime dropped
	dived  bool
	output string
	hours  map[int64]bool
}

func (f *UptimeFilter) Init(config interface{}) error {
	return nil
}

func (f *UptimeFilter) Run(runner pipeline.FilterRunner, helper pipeline.PluginHelper) (err error) {
	var (
		pack    *pipeline.PipelinePack
		payload string
	)

	inChan := runner.InChan()
	for pack = range inChan {
		payload = pack.Message.GetPayload()
		runner.LogMessage("Payload: " + payload)
		if f.hours == nil {
			f.hours = make(map[int64]bool)
		}
		var epoch int64 = f.GetEpoch(payload)
		f.startHour, f.endHour = f.FigureOutStartAndEndHour(epoch)
		if !f.hours[f.startHour] {
			f.InitFilterForStartHour(f.startHour, payload)
		} else {
			f.CalculateUptimeFor(f.startHour, f.endHour)
			// f.hours[&f.startHour] = false
			log.Printf("Length of map: &d", len(f.hours))
		}
		log.Printf("Start hour: %d", f.startHour)
		log.Printf("End hour: %d", f.endHour)
		log.Printf("EPOCH: %d", epoch)
		pack.Recycle()
	}
	return
}

func (f *UptimeFilter) CalculateUptimeFor(startHour, endHour int64) {
	// TODO
	log.Printf("This function will calculate total uptime for period between: %d and %d", startHour, endHour)
}

func (f *UptimeFilter) InitFilterForStartHour(startHour int64, payload string) {
	f.initialUptime = f.GetUptime(payload)
	f.lastUptime = f.initialUptime
	f.totalUptime = f.initialUptime
	f.dived = false
	f.hours[startHour] = true
	log.Printf("Stuff for hour: %d initialized", &startHour)
}

func (f *UptimeFilter) GetEpoch(payload string) (epoch int64) {
	split := strings.Split(payload, " ")
	if len(split) < 3 {
		return 1
	}
	if epochValue, err := strconv.Atoi(strings.TrimSpace(split[2])); err == nil {
		return int64(epochValue)
	} else {
		log.Printf("Error while trying to convert Epoch string to int64: %s", err)
	}
	return 0
}

func (f *UptimeFilter) GetUptime(payload string) (uptime int64) {
	split := strings.Split(payload, " ")
	if len(split) < 3 {
		return 1
	}
	if uptimeValue, err := strconv.ParseFloat(split[1], 64); err == nil {
		return int64(uptimeValue)
	} else {
		log.Printf("Error while trying to convert Uptime string to int64: %s", err)
	}
	return 0
}

// UDP payloads with uptime can be delivered anytime. There's only a little chance
// that this will occurr at e.g. 06:00
// Need to parse epoch to figure out to which time period (between start and
// end hour) epoch value belongs
// Return start and end hour
func (f *UptimeFilter) FigureOutStartAndEndHour(epoch int64) (startHour, endHour int64) {
	startEpoch := time.Unix(epoch, 0)
	start := startEpoch.Hour()
	endEpoch := startEpoch.Add(time.Hour)
	end := endEpoch.Hour()
	return int64(start), int64(end)
}

func init() {
	pipeline.RegisterPlugin("UptimeFilter", func() interface{} {
		return new(UptimeFilter)
	})
}

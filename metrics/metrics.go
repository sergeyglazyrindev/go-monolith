package metrics

import (
	"encoding/json"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"github.com/uadmin/rrd"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// SystemMetrics enables GoMonolith system metrics to be recorded
var SystemMetrics = false

// UserMetrics enables the user metrics to be recorded
var UserMetrics = false

func getRRDTemplate(name string) (*rrd.RRD, error) {
	tmpl := rrd.RRD{}
	buf, err := ioutil.ReadFile("templates/go-monolith/rrd/" + name + ".json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf, &tmpl)
	if err != nil {
		return nil, err
	}
	return &tmpl, nil
}

// NewMetric creates a new metric
func NewMetric(name string, template string) error {
	if strings.HasPrefix(name, "go-monolith/") && !SystemMetrics {
		return nil
	}
	if !strings.HasPrefix(name, "go-monolith/") && !UserMetrics {
		return nil
	}
	tmpl, err := getRRDTemplate(template)
	if err != nil {
		core.Trail(core.ERROR, "NewMetric.getRRDTemplate. %s", err.Error())
		return err
	}
	fName := "rrd/" + name + ".rrd"
	err = rrd.CreateRRD(fName, *tmpl)
	if err != nil {
		core.Trail(core.ERROR, "NetMetric.rrd.CreateRRD. %s", err.Error())
	}
	return err
}

// SetMetric sets the value of a gauge metric
func SetMetric(name string, value float64) {
	go setRRDValue(name, value, "gauge")
}

// IncrementMetric increments the value of a
func IncrementMetric(name string) {
	go setRRDValue(name, 1, "absolute")
}

// TimeMetric runs a function and times it as a metric
func TimeMetric(name string, div float64, f func()) {
	sTime := time.Now()
	f()
	SetMetric(name, float64(time.Since(sTime).Nanoseconds())/div)
}

func setRRDValue(name string, value float64, tmpl string) error {
	var err error
	if strings.HasPrefix(name, "go-monolith/") && !SystemMetrics {
		return nil
	}
	if !strings.HasPrefix(name, "go-monolith/") && !UserMetrics {
		return nil
	}
	fName := "rrd/" + name + ".rrd"
	if _, err = os.Stat(fName); err != nil {
		err = NewMetric(name, tmpl)
		if err != nil {
			return err
		}
	}
	err = rrd.UpdateRRD(fName, 1, value)
	if err != nil {
		core.Trail(core.ERROR, "setRRDValue.rrd.UpdateRRD. %s", err.Error())
	}
	return err
}

package autoconfiguratorcmd

import (
	"wnd/api"
	"wnd/modules"
	"wnd/base"
	"wnd/utils/log"

	"flag"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

const (
	_FLAG_FORMAT = "%s.%s"
)

type cmdConfigurator struct {
	Runtime base.Runtime `inject:""`
}

type flagVar struct {
	i string
}

func (f *flagVar) Set(value string) error {
	f.i = value
	return nil
}

func (f *flagVar) String() string {
	return f.i
}

func New() modules.AutoConfigurator {
	return &cmdConfigurator{}
}

func (this *cmdConfigurator) ID() string {
	return "autoConfiguratorCmd"
}

func (this *cmdConfigurator) Priority() int8 {
	return math.MinInt8
}

func (this *cmdConfigurator) Init() error {
	return this.Configure()
}

func (this *cmdConfigurator) Configure() error {
	modules := this.Runtime.GetByType(reflect.TypeOf((*api.ConfigurableModule)(nil)).Elem())

	for _, v := range modules {
		c,_ := v.(api.ConfigurableModule)
		m,_ := v.(api.GameModule)
		
		tkvs := c.Configuration()
		
		log.Debugf("configurable %s => %# v", m.ID(), tkvs)

		for _, tkv := range tkvs {
			v := new(flagVar)
			
			flagId := fmt.Sprintf(_FLAG_FORMAT, m.ID(), tkv.Key)
			
			flag.Var(v, flagId, string(tkv.Name))
			
			var e error
			
			if v.i == "" {
				continue
			} else {
				log.Debugf("%s got flag %v", flagId, v.i)
			}

			switch tkv.Type {
			case reflect.String:
				tkv.Value = v.i
			case reflect.Int:
				tkv.Value, e = strconv.Atoi(v.i)
			case reflect.Bool:
				tkv.Value, e = strconv.ParseBool(v.i)
			default:
				log.Warnf("unsupported config type: %v", tkv.Type.String())
			}

			if e != nil {
				log.Errorf("error parsing %v as %v", v.i, tkv.Type.String())
			}
			
			c.SetConfiguration(tkv)
		}
	}

	return nil
}

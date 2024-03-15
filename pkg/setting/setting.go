package setting

import (
	"bypctl/pkg/global"
	"github.com/gogf/gf/v2/util/gconv"
)

func Init() {

	global.Log.Infof("setting System init----> %v", gconv.String(global.Conf.System))
}

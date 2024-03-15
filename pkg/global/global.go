package global

import (
	"bypctl/pkg/config"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	Log  *logrus.Logger
	Conf config.Config
)

package xk6sftp

import (
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/sftp", new(Client))
}

var logger = logrus.New()

package main

import (
	"bytes"
	"flag"
	"os"

	"github.com/opensourceways/community-robot-lib/giteeclient"
	"github.com/opensourceways/community-robot-lib/logrusutil"
	liboptions "github.com/opensourceways/community-robot-lib/options"
	framework "github.com/opensourceways/community-robot-lib/robot-gitee-framework"
	"github.com/sirupsen/logrus"
)

type options struct {
	service liboptions.ServiceOptions
	gitee   liboptions.GiteeOptions
}

func (o *options) Validate() error {
	if err := o.service.Validate(); err != nil {
		return err
	}

	return o.gitee.Validate()
}

func GetSecret(secretPath string) []byte {
	token, err := os.ReadFile(secretPath)

	if err != nil {
		logrus.WithError(err).Fatal("read token file fail")
		return nil
	}

	token = bytes.TrimRight(token, "\r\n")
	return token
}

func GetTokenGenerator(secretPath string) func() []byte {
	return func() []byte {
		return GetSecret(secretPath)
	}
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.gitee.AddFlags(fs)
	o.service.AddFlags(fs)

	fs.Parse(args)
	return o
}

func main() {
	logrusutil.ComponentInit(botName)

	o := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err := o.Validate(); err != nil {
		logrus.WithError(err).Fatal("Invalid options")
	}

	c := giteeclient.NewClient(GetTokenGenerator(o.gitee.TokenPath))

	if err := os.Remove(o.gitee.TokenPath); err != nil {
		logrus.WithError(err).Error("fatal error occurred while deleting token")
	}

	v, err := c.GetBot()
	if err != nil {
		logrus.WithError(err).Error("Error get bot name")
	}

	r := newRobot(c, v.Login)

	framework.Run(r, o.service)
}

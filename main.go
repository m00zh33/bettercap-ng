package main

import (
	"runtime"
	"strings"

	"github.com/op/go-logging"

	"github.com/evilsocket/bettercap-ng/core"
	"github.com/evilsocket/bettercap-ng/session"
	"github.com/evilsocket/bettercap-ng/session/modules"
)

var sess *session.Session
var log = logging.MustGetLogger("mitm")
var err error

func main() {
	if sess, err = session.New(); err != nil {
		panic(err)
	}

	log.Infof("Starting %s v%s\n", core.Name, core.Version)
	log.Infof("Build: date=%s os=%s arch=%s\n", core.BuildDate, runtime.GOOS, runtime.GOARCH)

	sess.Register(session_modules.NewProber(sess))
	sess.Register(session_modules.NewDiscovery(sess))
	sess.Register(session_modules.NewArpSpoofer(sess))
	sess.Register(session_modules.NewSniffer(sess))
	sess.Register(session_modules.NewHttpProxy(sess))

	if err = sess.Start(); err != nil {
		log.Fatal(err)
	}

	defer sess.Close()

	if *sess.Options.Commands != "" {
		for _, cmd := range strings.Split(*sess.Options.Commands, ";") {
			cmd = strings.Trim(cmd, "\r\n\t ")
			if err = sess.Run(cmd); err != nil {
				log.Fatal(err)
			}
		}
	}

	if *sess.Options.Caplet != "" {
		if err = sess.RunCaplet(*sess.Options.Caplet); err != nil {
			log.Fatal(err)
		}
	}

	for sess.Active {
		line, err := sess.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		if line == "" || line[0] == '#' {
			continue
		}

		if err = sess.Run(line); err != nil {
			log.Error(err)
		}
	}
}

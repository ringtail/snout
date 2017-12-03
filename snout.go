package snout

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ringtail/snout/core"
)

const USAGE_DESC = `

=====================================================

     ███████╗███╗   ██╗ ██████╗ ██╗   ██╗████████╗
     ██╔════╝████╗  ██║██╔═══██╗██║   ██║╚══██╔══╝
     ███████╗██╔██╗ ██║██║   ██║██║   ██║   ██║
     ╚════██║██║╚██╗██║██║   ██║██║   ██║   ██║
     ███████║██║ ╚████║╚██████╔╝╚██████╔╝   ██║
     ╚══════╝╚═╝  ╚═══╝ ╚═════╝  ╚═════╝    ╚═╝

  snout is tool to improve your system performance
=====================================================


Usage:
	snout [commands|flags]

The commands & flags are:
	help 				print snout help
	version 			print the version to stdout
	--debug 			switch on debug mode

Examples:
	# start in normal mode
	snout

	# start in debug mode
	snout  --debug
`

var (
	h     bool
	v     bool
	debug bool
)

func init() {
	flag.BoolVar(&h, "help", false, "--help")
	flag.BoolVar(&v, "version", false, "--version")
	flag.BoolVar(&debug, "debug", false, "--debug")
}

func main() {
	flag.Parse()
	if h == true {
		fmt.Println(USAGE_DESC)
		return
	}
	if v == true {
		fmt.Println(core.VERSION)
		return
	}

	customFormatter := new(log.TextFormatter)
	customFormatter.DisableTimestamp = true
	log.SetFormatter(customFormatter)

	if debug == true {
		log.SetLevel(log.DebugLevel)
	}

	st := &core.Snout{}
	st.Run()
}

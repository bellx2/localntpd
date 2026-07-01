package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kardianos/service"
	"github.com/t7b/localntpd/internal/ntp"
)

const (
	serviceName        = "localntpd"
	serviceDisplayName = "Simple NTP Server"
	serviceDescription = "A simple NTP server that serves the local PC's clock"
)

var (
	addr    = flag.String("addr", ":123", "listen address (e.g. :123, 0.0.0.0:123)")
	stratum = flag.Uint("stratum", 2, "stratum (1-15)")
)

type program struct {
	cancel context.CancelFunc
}

func (p *program) Start(s service.Service) error {
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	go func() {
		srv := ntp.NewServer(*addr, byte(*stratum))
		if err := srv.Run(ctx); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()
	return nil
}

func (p *program) Stop(s service.Service) error {
	if p.cancel != nil {
		p.cancel()
	}
	return nil
}

func main() {
	flag.Usage = printUsage

	svcConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceDisplayName,
		Description: serviceDescription,
		Arguments:   []string{"run"},
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// 先頭引数がコマンド（フラグでない）ならサービス制御などを処理
	args := os.Args[1:]
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		switch cmd := args[0]; cmd {
		case "install", "uninstall", "start", "stop", "restart", "status":
			if err := service.Control(s, cmd); err != nil {
				log.Fatalf("%s failed: %v", cmd, err)
			}
			fmt.Printf("%s done\n", cmd)
			return
		case "help":
			printUsage()
			return
		case "run":
			args = args[1:] // "run" を除いて残りをフラグとして解析
		default:
			fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
			printUsage()
			os.Exit(1)
		}
	}

	// 残り（引数なし / フラグ直指定 / "run" 以降）を解析。-h/--help は flag.Usage を呼ぶ
	flag.CommandLine.Parse(args)

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `localntpd - a simple NTP server that serves the local PC's clock

Usage:
  localntpd [command] [options]

Commands:
  run        run in the foreground (default)
  install    register as a system service
  uninstall  remove the service
  start      start the service
  stop       stop the service
  restart    restart the service
  status     show the service status
  help       show this help

Options:
  -addr string     listen address (default: :123)
  -stratum uint    stratum (default: 2)

Examples:
  localntpd run -addr :12345          # run on a non-privileged port
  localntpd install                   # register as a service (requires administrator privileges)
  localntpd start

Notes:
  Binding to port 123 requires administrator/root privileges.
  For testing, use a non-privileged port such as -addr :12345.
`)
}

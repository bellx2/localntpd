package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bellx2/localntpd/internal/ntp"
	"github.com/kardianos/service"
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
	done   chan struct{}
}

func (p *program) Start(s service.Service) error {
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.done = make(chan struct{})

	go func() {
		defer close(p.done)
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
	if p.done != nil {
		<-p.done
	}
	return nil
}

func main() {
	flag.Usage = printUsage

	// 先頭引数がコマンド（フラグでない）ならサービス制御などを処理
	args := os.Args[1:]
	cmd := ""
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		cmd = args[0]
		args = args[1:]
	}

	switch cmd {
	case "help":
		printUsage()
		return
	case "install", "uninstall", "start", "stop", "restart", "status", "run", "":
		// 続行
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}

	// install / run ではフラグを解析する（サービス起動時の Arguments にも使う）
	if cmd == "install" || cmd == "run" || cmd == "" {
		if err := flag.CommandLine.Parse(args); err != nil {
			os.Exit(2)
		}
		if err := validateFlags(); err != nil {
			log.Fatal(err)
		}
	}

	svcConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceDisplayName,
		Description: serviceDescription,
		Arguments:   serviceArguments(),
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	switch cmd {
	case "install", "uninstall", "start", "stop", "restart", "status":
		if err := service.Control(s, cmd); err != nil {
			log.Fatalf("%s failed: %v", cmd, err)
		}
		fmt.Printf("%s done\n", cmd)
		return
	}

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

// serviceArguments はサービス登録時に渡す起動引数を組み立てる
func serviceArguments() []string {
	args := []string{"run"}
	if *addr != ":123" {
		args = append(args, "-addr", *addr)
	}
	if *stratum != 2 {
		args = append(args, "-stratum", fmt.Sprintf("%d", *stratum))
	}
	return args
}

func validateFlags() error {
	if *stratum < 1 || *stratum > 15 {
		return fmt.Errorf("stratum must be between 1 and 15, got %d", *stratum)
	}
	return nil
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
  localntpd install -addr :12345      # register as a service with custom addr
  localntpd start

Notes:
  Binding to port 123 requires administrator/root privileges.
  For testing, use a non-privileged port such as -addr :12345.
  Options passed to install are stored as service arguments.
`)
}

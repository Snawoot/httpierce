package main

import (
    "flag"
    "log"
    //"net"
    "time"
    "os"
)

type CLIArgs struct {
    Bind, Dst string
    Server bool
    Timeout time.Duration
    VPN bool
    flagset *flag.FlagSet
}

func NewCLIArgs() *CLIArgs {
    args := CLIArgs{flagset: flag.NewFlagSet(os.Args[0], flag.ContinueOnError)}
    args.flagset.StringVar(&args.Bind, "bind", "", "listen address")
    args.flagset.StringVar(&args.Dst, "dst", "", "target address")
    args.flagset.BoolVar(&args.Server, "server", false, "server-side mode")
    args.flagset.BoolVar(&args.VPN, "V", false, "VPN mode. Used by shadowsocks on Android")
    args.flagset.DurationVar(&args.Timeout, "timeout", 10 * time.Second, "connect timeout")
    return &args
}

func (args *CLIArgs) Update(values []string) error {
    return args.flagset.Parse(values)
}

func main() {
    args := NewCLIArgs()

    pluginArgs, err := NewPluginArgs()
    if err != nil {
        log.Printf("main: plugin interface error: %v", err)
    } else {
        log.Print("main: running in plugin mode")

        opts := pluginArgs.ExportOptions()

        if err := args.Update(opts); err != nil {
            log.Printf("main: WARNING: CLIArgs.Update: %v", err)
        }

        if args.Server {
            args.Dst = pluginArgs.GetLocalAddr()
            args.Bind = pluginArgs.GetRemoteAddr()
        } else {
            args.Bind = pluginArgs.GetLocalAddr()
            args.Dst = pluginArgs.GetRemoteAddr()
        }
    }
    err = args.Update(os.Args[1:])
    if err != nil {
        log.Fatalf("main: commandLine.Parse: %v", err)
    }

    if args.Bind == "" {
        log.Fatal("main: bind addr is required")
    }
    if args.Dst == "" {
        log.Fatal("main: destination addr is required")
    }

    /*if args.server {
        lc := net.ListenConfig{Control: core.GetControlFunc(&core.TcpConfig{EnableTFO: tfo})}
        l, err := lc.Listen(context.Background(), "tcp", args.bind)
        if err != nil {
            log.Fatalf("main: net.Listen: %v", err)
        }

        err = core.DoServer(l, certificates, args.dst, sendPaddingData, timeout)
        if err != nil {
            log.Fatalf("main: doServer: %v", err)
        }

    } else { // do client
        lc := net.ListenConfig{}
        l, err := lc.Listen(context.Background(), "tcp", args.bind)
        if err != nil {
            log.Fatalf("main: net.Listen: %v", err)
        }

        err = core.DoClient(l, args.dst, host, rootCAs, insecureSkipVerify, sendPaddingData, timeout, vpn, tfo)
        if err != nil {
            log.Fatalf("main: doServer: %v", err)
        }
    }
    */
}

package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Unknwon/com"
	"github.com/appscode/go/flags"
	"github.com/appscode/go/net"
	"github.com/go-macaron/auth"
	"github.com/go-macaron/toolbox"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	hostUtil "github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	flag "github.com/spf13/pflag"
	macaron "gopkg.in/macaron.v1"
)

func main() {
	selfIP := net.GetInternalIP()
	if selfIP == "" {
		// may be peers are running in HostNetwork mode and host only has public IP
		selfIP = net.GetExternalIPs()[0]
	}
	log.Println("Detected IP for hostfacts server:", selfIP)

	host := flag.String("host", selfIP, "Http server ip address")
	port := flag.Int("port", 56977, "Http server port")
	caCertFile := flag.String("caCertFile", "", "File containing CA certificate")
	certFile := flag.String("certFile", "", "File container server TLS certificate")
	keyFile := flag.String("keyFile", "", "File containing server TLS private key")

	flags.InitFlags()
	flags.DumpAll()

	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())

	// auth
	username := os.Getenv("AUTH_USERNAME")
	password := os.Getenv("AUTH_PASSWORD")
	token := os.Getenv("AUTH_TOKEN")
	if username != "" && password != "" {
		m.Use(auth.Basic(username, password))
	} else if token != "" {
		m.Use(auth.Bearer(token))
	}

	m.Use(toolbox.Toolboxer(m))
	m.Use(macaron.Renderer(macaron.RenderOptions{
		IndentJSON: true,
	}))

	m.Get("/cpu", func(ctx *macaron.Context) {
		r, _ := cpu.Info()
		ctx.JSON(200, r)
	})

	m.Get("/virt_mem", func(ctx *macaron.Context) {
		r, _ := mem.VirtualMemory()
		ctx.JSON(200, r)
	})
	m.Get("/swap_mem", func(ctx *macaron.Context) {
		r, _ := mem.SwapMemory()
		ctx.JSON(200, r)
	})

	m.Get("/host", func(ctx *macaron.Context) {
		r, _ := hostUtil.Info()
		ctx.JSON(200, r)
	})
	m.Get("/uptime", func(ctx *macaron.Context) {
		r, _ := hostUtil.Uptime()
		ctx.JSON(200, r)
	})

	m.Get("/disks", func(ctx *macaron.Context) {
		r, _ := disk.Partitions(true)
		ctx.JSON(200, r)
	})

	m.Get("/du", func(ctx *macaron.Context) {
		paths := ctx.QueryStrings("p")
		du := make([]*disk.UsageStat, len(paths))
		for i, p := range paths {
			du[i], _ = disk.Usage(p)
		}
		ctx.JSON(200, du)
	})

	m.Get("/load", func(ctx *macaron.Context) {
		l, _ := load.Avg()
		ctx.JSON(200, l)
	})

	addr := *host + ":" + com.ToStr(*port)
	log.Printf("listening on %s (%s)\n", addr, macaron.Env)

	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      m,
	}
	if *caCertFile == "" && *certFile == "" && *keyFile == "" {
		log.Fatalln(srv.ListenAndServe())
	} else {
		/*
			Ref:
			 - https://blog.cloudflare.com/exposing-go-on-the-internet/
			 - http://www.bite-code.com/2015/06/25/tls-mutual-auth-in-golang/
			 - http://www.hydrogen18.com/blog/your-own-pki-tls-golang.html
		*/
		caCert, err := ioutil.ReadFile(*caCertFile)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
			SessionTicketsDisabled:   true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				// tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
				// tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
			ClientCAs:  caCertPool,
			ClientAuth: tls.VerifyClientCertIfGiven,
		}
		tlsConfig.BuildNameToCertificate()

		srv.TLSConfig = tlsConfig
		log.Fatalln(srv.ListenAndServeTLS(*certFile, *keyFile))
	}
}

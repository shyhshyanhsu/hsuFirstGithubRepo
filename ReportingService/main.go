package main

import (
	"context"
	"flag"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"os"
)

func main() {
	var (
		listen= flag.String("listen", ":8090", "HTTP listen address")
		proxy= flag.String("proxy", ":8091", "Comma-separated list of URLs to proxy POS APIs requests")
	)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

	var svc ReportingService
	svc = reportingService{}
	svc = loggingMiddleware{logger, svc}
	svc = proxyingMiddleware(context.Background(), *proxy, logger)(svc)

	posHandler := httptransport.NewServer(
		makePOSEndpoint(svc),
		decodeBusinessesRequest,
		encodeResponse,
	)

	http.Handle("/reporting", posHandler)
	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}


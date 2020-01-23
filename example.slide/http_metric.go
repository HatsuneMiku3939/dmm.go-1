package main

import (
	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/HatsuneMiku3939/ocecho"
	"go.opencensus.io/stats/view"
	"net/http"
)

func InitPrometheus() error {
	pe, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		panic(err)
	}
	if err := view.Register(ocecho.DefaultServerViews...); err != nil {
		panic(err)
	}

	view.RegisterExporter(pe)
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", pe)
		if err := http.ListenAndServe(":8888", mux); err != nil {
			panic(err)
		}
	}()
	return nil
}

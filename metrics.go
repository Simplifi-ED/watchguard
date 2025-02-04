package main

import "github.com/prometheus/client_golang/prometheus"

var (
	resourceEvents = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "k8s_resource_events_total",
			Help: "Total number of Kubernetes resource events by type and resource",
		},
		[]string{"event_type", "group", "resource"},
	)

	watchedResourcesGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "k8s_watched_resources_total",
			Help: "Total number of Kubernetes resources being watched",
		},
	)

	watchEvents = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "k8s_watch_events_total",
			Help: "Total number of Kubernetes watch events",
		},
		[]string{"resource", "namespace", "event_type"},
	)
)

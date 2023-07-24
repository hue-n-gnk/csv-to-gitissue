package factory

import (
	"time"
)

const (
	// ConnMaxIdle must be lower than ConnMaxOpen, and strongly recommended to keep it low,
	// since PostgreSQL allocates 2-3MB memory per ConnMaxIdle.
	//
	// Estimate:
	//   ConnMaxIdle * 3MB Memory * GKE Pods = Allocated Server Memory
	//
	// Examples:
	//   Assume we have 4 APIs scaled to have 3 Pod replicas, plus Workers in single Pod.
	//   The total memory allocation to PostgreSQL instance will be 196MB.
	//
	//   a) 4 APIs:    ConnMaxIdle(5) * 3MB * Pods(12) --> 180MB
	//   b) 4 Workers: ConnMaxIdle(5) * 3MB * Pods(4)  -->  16MB
	//                                        TOTAL:   --> 196MB
	ConnMaxIdle = 5

	// ConnMaxOpen must be lower than PostgreSQL instance's max connections.
	//
	// Estimate:
	//   ConnMaxOpen * GKE Pods = Max Connections
	//
	// Examples:
	//   Assume we have 4 APIs scaled to have 3 Pod replicas, plus Workers in single Pod.
	//   The total max connections to PostgreSQL instance will be 320 max connections.
	//
	//   a) 4 APIs:    ConnMaxOpen(20) * Pods(12) --> 240
	//   b) 4 Workers: ConnMaxOpen(20) * Pods(4)  -->  80
	//                                   TOTAL:   --> 320 Max Connections
	ConnMaxOpen = 20

	// Setting this number too low will increase premature timeout between REST API and datastores,
	// and too high will make load balancers to go timeout.
	//
	TransactionTimeout = time.Second * 15
)

package common

import "os"

var (
	ClusterName = os.Getenv("NAIS_CLUSTER_NAME")
)

package texas

import (
	"fmt"
	"os"
)

type Target struct {
	Application string
	Namespace   string
	Cluster     string
}

func (f Target) String() string {
	cluster := f.Cluster
	if cluster == "" {
		cluster = os.Getenv("NAIS_CLUSTER_NAME")
	}
	namespace := f.Namespace
	if namespace == "" {
		namespace = os.Getenv("NAIS_NAMESPACE")
	}
	return fmt.Sprintf("api://%s.%s.%s/.default", cluster, namespace, f.Application)
}

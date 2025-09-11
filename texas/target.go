package texas

import (
	"fmt"
	"os"
)

var (
	cluster   = os.Getenv("NAIS_CLUSTER_NAME")
	namespace = os.Getenv("NAIS_NAMESPACE")
)

type Target struct {
	Application string
	Namespace   string
	Cluster     string
}

func (f Target) String() string {
	c := f.Cluster
	if c == "" {
		c = cluster
	}
	n := f.Namespace
	if n == "" {
		n = namespace
	}
	return fmt.Sprintf("api://%s.%s.%s/.default", c, n, f.Application)
}

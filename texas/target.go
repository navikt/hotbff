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

func (t Target) String() string {
	c := t.Cluster
	if c == "" {
		c = cluster
	}
	n := t.Namespace
	if n == "" {
		n = namespace
	}
	return fmt.Sprintf("api://%s.%s.%s/.default", c, n, t.Application)
}

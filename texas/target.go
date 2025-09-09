package texas

import (
	"fmt"

	"github.com/navikt/hotbff/common"
)

type Target struct {
	Application string
	Namespace   string
	Cluster     string
}

func (f Target) String() string {
	cluster := f.Cluster
	if cluster == "" {
		cluster = common.ClusterName
	}
	namespace := f.Namespace
	if namespace == "" {
		namespace = common.Namespace
	}
	return fmt.Sprintf("api://%s.%s.%s/.default", cluster, namespace, f.Application)
}

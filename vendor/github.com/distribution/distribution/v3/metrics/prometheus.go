package metrics

import "github.com/docker/go-metrics"

const (
	// NamespacePrefix is the namespace of prometheus metrics
	NamespacePrefix = "registry"
)

var (
	// StorageNamespace is the prometheus namespace of blob/cache related operations
	StorageNamespace = metrics.NewNamespace(NamespacePrefix, "storage", nil)

	// NotificationsNamespace is the prometheus namespace of notification related metrics
	NotificationsNamespace = metrics.NewNamespace(NamespacePrefix, "notifications", nil)

	// ProxyNamespace is the prometheus namespace of proxy related metrics
	ProxyNamespace = metrics.NewNamespace(NamespacePrefix, "proxy", nil)
)

package entity

type PrometheusInstanceParams struct {
	PromForCloudService *PromForCloudService `json:"prom_for_cloud_service"`
}

type PromForCloudService struct {
	CesMetricNamespaces []string `json:"ces_metric_namespaces"`
}

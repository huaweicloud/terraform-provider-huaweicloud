package aom

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"io"
	"strings"
	"time"
)

func ResourcePrometheusInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourcePrometheusInstanceRead,
		CreateContext: resourcePrometheusInstancePatch,
		DeleteContext: resourcePrometheusInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"prom_for_cloud_service": {
				Type:     schema.TypeList,
				Optional: true,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"ces_metric_namespaces": {
						Type:     schema.TypeList,
						Optional: true,
						Required: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				}},
			},
		},
	}
}

func resourcePrometheusInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, dErr := httpclient_go.NewHttpClientGo(config)
	if dErr != nil {
		return dErr
	}
	client.WithMethod(httpclient_go.MethodGet).WithUrlWithoutEndpoint(config, "aom",
		config.GetRegion(d), "v1/"+config.TenantID+"/prometheus-instances?action=prom_for_cloud_service")

	resp, err := client.Do()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AOM alarm rule")
	}

	mErr := &multierror.Error{}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		mErr = multierror.Append(mErr, err)
	}
	rlt := &entity.PrometheusInstanceParams{}

	err = json.Unmarshal(body, rlt)

	if err != nil {
		mErr = multierror.Append(mErr, err)
	}
	d.Set("prom_for_cloud_service", buildPrometheusInstanceMap(rlt))
	d.Set("action", "prom_for_cloud_service")
	d.Set("project_id", config.TenantID)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error getting AOM prometheus instance fields: %s", err)
	}

	return nil
}

func buildPrometheusInstanceMap(rlt *entity.PrometheusInstanceParams) []map[string]interface{} {
	var m = make(map[string]interface{})
	m["ces_metric_namespaces"] = rlt.PromForCloudService.CesMetricNamespaces
	return []map[string]interface{}{m}
}

func resourcePrometheusInstancePatch(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, dErr := httpclient_go.NewHttpClientGo(config)
	if dErr != nil {
		return dErr
	}
	var p = d.Get("prom_for_cloud_service").([]interface{})[0]
	service := p.(map[string]interface{})["ces_metric_namespaces"].([]interface{})
	namespace := make([]string, 0)
	for _, s := range service {
		namespace = append(namespace, s.(string))
	}
	patchOpts := &entity.PrometheusInstanceParams{
		PromForCloudService: &entity.PromForCloudService{CesMetricNamespaces: namespace},
	}
	client.WithMethod(httpclient_go.MethodPost).WithUrlWithoutEndpoint(config, "aom",
		config.GetRegion(d), "v1/"+config.TenantID+"/prometheus-instances?action=prom_for_cloud_service").WithBody(patchOpts)
	r, err := client.Do()

	if r.StatusCode != 204 || err != nil {
		return common.CheckDeletedDiag(d, err, "error add prometheus-instances")
	}
	d.SetId(strings.Join(namespace, ","))
	return resourcePrometheusInstanceRead(context.TODO(), d, meta)
}

func resourcePrometheusInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.Set("prom_for_cloud_service",
		entity.PrometheusInstanceParams{PromForCloudService: &entity.PromForCloudService{CesMetricNamespaces: []string{}}})
	return resourcePrometheusInstancePatch(ctx, d, meta)
}

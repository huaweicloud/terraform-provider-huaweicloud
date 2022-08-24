package aom

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"io"
	"regexp"
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
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 100),
					validation.StringMatch(regexp.MustCompile(
						"^[\u4e00-\u9fa5A-Za-z0-9]([\u4e00-\u9fa5-_A-Za-z0-9]*[\u4e00-\u9fa5A-Za-z0-9])?$"),
						"The name can only consist of letters, digits, underscores (_),"+
							" hyphens (-) and chinese characters, and it must start and end with letters,"+
							" digits or chinese characters."),
				),
			},
			"prom_for_cloud_service": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"ces_metric_namespaces": {
						Type:     schema.TypeList,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				}},
			},

			"action": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error getting AOM prometheus instance fields: %w", err)
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

func resourcePrometheusInstanceDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return nil
}

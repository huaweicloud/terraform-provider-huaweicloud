package aom

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	prometheusInstanceNotExistsCode = "AOM.11017014"
)

// @API AOM POST /v1/{project_id}/aom/prometheus
// @API AOM GET /v1/{project_id}/aom/prometheus
// @API AOM PUT /v1/{project_id}/aom/prometheus
// @API AOM DELETE /v1/{project_id}/aom/prometheus
func ResourcePromInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePromInstanceCreate,
		ReadContext:   resourcePromInstanceRead,
		UpdateContext: resourcePromInstanceUpdate,
		DeleteContext: resourcePromInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the prometheus instance is located.`,
			},

			// Required Parameters
			"prom_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the prometheus instance.`,
			},
			"prom_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The type of the prometheus instance.`,
			},

			// Optional Parameters
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The enterprise project ID to which the prometheus instance belongs.`,
			},
			"prom_version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The version of the prometheus instance.`,
			},
			"prom_limits": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compactor_blocks_retention_period": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The retention period for the compactor blocks.`,
						},
					},
				},
				MaxItems:    1,
				Description: `The limit configurations of the prometheus instance.`,
			},

			// Attributes
			"prom_http_api_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The HTTP URL for calling the prometheus instance.`,
			},
			"remote_read_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The remote read address of the prometheus instance.`,
			},
			"remote_write_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The remote write address of the prometheus instance.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the prometheus instance, in RFC3339 format.`,
			},
		},
	}
}

func buildCreatePrometheusInstanceBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"prom_name":             d.Get("prom_name"),
		"prom_type":             d.Get("prom_type"),
		"prom_version":          utils.ValueIgnoreEmpty(d.Get("prom_version")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}
}

func resourcePromInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("aom", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	httpUrl := "v1/{project_id}/aom/prometheus"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCreatePrometheusInstanceBodyParams(d, cfg)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating AOM prometheus instance: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId := utils.PathSearch(fmt.Sprintf("prometheus[?prom_name=='%s']|[0].prom_id", d.Get("prom_name")), respBody, "").(string)
	if resourceId == "" {
		return diag.Errorf("unable to find prometheus instance ID from the API response")
	}

	d.SetId(resourceId)

	if _, ok := d.GetOk("prom_limits"); ok {
		err = updatePrometheusInstance(cfg, client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourcePromInstanceRead(ctx, d, meta)
}

func GetPrometheusInstanceById(client *golangsdk.ServiceClient, isntanceId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/aom/prometheus?prom_id={instance_id}"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", isntanceId)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":          "application/json",
			"Enterprise-Project-Id": "all_granted_eps",
		},
	}
	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AOM prometheus instance: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening AOM prometheus instance: %s", err)
	}

	instance := utils.PathSearch("prometheus[0]", respBody, nil)
	if instance == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/aom/prometheus",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("prometheus instance (%s) not found", isntanceId)),
			},
		}
	}

	return instance, nil
}

func flattenPrometheusInstancePromLimits(promLimits interface{}) interface{} {
	if promLimits == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"compactor_blocks_retention_period": utils.PathSearch("compactor_blocks_retention_period", promLimits, nil),
		},
	}
}

func resourcePromInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	instance, err := GetPrometheusInstanceById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AOM prometheus instance")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("prom_name", utils.PathSearch("prom_name", instance, nil)),
		d.Set("prom_type", utils.PathSearch("prom_type", instance, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", instance, nil)),
		d.Set("prom_version", utils.PathSearch("prom_version", instance, nil)),
		d.Set("prom_limits", flattenPrometheusInstancePromLimits(utils.PathSearch("prom_limits", instance, nil))),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("prom_create_timestamp", instance, float64(0)).(float64))/1000, false)),
		d.Set("remote_write_url", utils.PathSearch("prom_spec_config.remote_write_url", instance, nil)),
		d.Set("remote_read_url", utils.PathSearch("prom_spec_config.remote_read_url", instance, nil)),
		d.Set("prom_http_api_endpoint", utils.PathSearch("prom_spec_config.prom_http_api_endpoint", instance, nil)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting AOM prometheus instance fields: %s", err)
	}

	return nil
}

func buildPrometheusInstancePromLimits(promLimits []interface{}) interface{} {
	if len(promLimits) < 1 {
		return nil
	}

	promLimit := promLimits[0]
	return map[string]interface{}{
		"compactor_blocks_retention_period": utils.PathSearch("compactor_blocks_retention_period", promLimit, nil),
	}
}

func buildUpdatePrometheusInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"prom_id":     d.Id(),
		"prom_name":   d.Get("prom_name"),
		"prom_limits": buildPrometheusInstancePromLimits(d.Get("prom_limits").([]interface{})),
	}
}

func updatePrometheusInstance(cfg *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/aom/prometheus"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeaders(cfg, d),
		JSONBody:         utils.RemoveNil(buildUpdatePrometheusInstanceBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return err
	}
	return nil
}

func resourcePromInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("aom", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	err = updatePrometheusInstance(cfg, client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePromInstanceRead(ctx, d, meta)
}

func resourcePromInstanceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("aom", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	httpUrl := "v1/{project_id}/aom/prometheus?prom_id={instance_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeaders(cfg, d),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", prometheusInstanceNotExistsCode),
			"error deleting AOM prometheus instance")
	}
	return nil
}

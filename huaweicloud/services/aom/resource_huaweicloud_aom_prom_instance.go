package aom

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	prometheusInstanceNotExistsCode = "AOM.11017014"
)

// @API AOM POST /v1/{project_id}/aom/prometheus
// @API AOM DELETE /v1/{project_id}/aom/prometheus
// @API AOM GET /v1/{project_id}/aom/prometheus
func ResourcePromInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePromInstanceCreate,
		ReadContext:   resourcePromInstanceRead,
		DeleteContext: resourcePromInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"prom_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"prom_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"prom_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			// attributes
			"remote_write_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_read_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prom_http_api_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePromInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "aom"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createPrometheusInstanceHttpUrl := "v1/{project_id}/aom/prometheus"
	createPrometheusInstanceHttpUrl = strings.ReplaceAll(createPrometheusInstanceHttpUrl, "{project_id}", client.ProjectID)
	createPrometheusInstancePath := client.Endpoint + createPrometheusInstanceHttpUrl

	createPrometheusInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createPrometheusInstanceOpt.JSONBody = utils.RemoveNil(buildCreatePrometheusInstanceBodyParams(d, cfg))
	createPrometheusInstanceResp, err := client.Request("POST", createPrometheusInstancePath, &createPrometheusInstanceOpt)
	if err != nil {
		return diag.Errorf("error creating AOM prometheus instance: %s", err)
	}
	createPrometheusInstanceRespBody, err := utils.FlattenResponse(createPrometheusInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}
	expression := fmt.Sprintf("prometheus[?prom_name== '%s'].prom_id | [0]", d.Get("prom_name"))
	id, err := jmespath.Search(expression, createPrometheusInstanceRespBody)
	if err != nil || id == nil {
		return diag.Errorf("error creating AOM prometheus instance: ID is not found in API response")
	}

	d.SetId(id.(string))

	return resourcePromInstanceRead(ctx, d, meta)
}

func buildCreatePrometheusInstanceBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"prom_name":             d.Get("prom_name"),
		"prom_type":             d.Get("prom_type"),
		"prom_version":          d.Get("prom_version"),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}
	return bodyParams
}

func resourcePromInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "aom"
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	getPrometheusInstanceHttpUrl := "v1/{project_id}/aom/prometheus"
	getPrometheusInstanceHttpUrl = strings.ReplaceAll(getPrometheusInstanceHttpUrl, "{project_id}", client.ProjectID)
	getPrometheusInstanceHttpUrl += fmt.Sprintf("?prom_id=%s", d.Id())
	getPrometheusInstancePath := client.Endpoint + getPrometheusInstanceHttpUrl

	getPrometheusInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getPrometheusInstanceResp, err := client.Request("GET", getPrometheusInstancePath, &getPrometheusInstanceOpt)
	if err != nil {
		diag.Errorf("error retrieving AOM prometheus instance: %s", err)
	}

	getPrometheusInstanceRespBody, err := utils.FlattenResponse(getPrometheusInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	curJson, err := jmespath.Search("prometheus[0]", getPrometheusInstanceRespBody)
	if err != nil {
		return diag.Errorf("error parsing AOM prometheus instance: %s", err)
	}

	if curJson == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving AOM prometheus instance")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("prom_name", utils.PathSearch("prom_name", curJson, nil)),
		d.Set("prom_type", utils.PathSearch("prom_type", curJson, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", curJson, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("prom_create_timestamp", curJson, float64(0)).(float64))/1000, false)),
		d.Set("prom_version", utils.PathSearch("prom_version", curJson, nil)),
		d.Set("remote_write_url", utils.PathSearch("prom_spec_config.remote_write_url", curJson, nil)),
		d.Set("remote_read_url", utils.PathSearch("prom_spec_config.remote_read_url", curJson, nil)),
		d.Set("prom_http_api_endpoint", utils.PathSearch("prom_spec_config.prom_http_api_endpoint", curJson, nil)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting AOM prometheus instance fields: %s", err)
	}

	return nil
}

func resourcePromInstanceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "aom"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	deletePrometheusInstanceHttpUrl := "v1/{project_id}/aom/prometheus"
	deletePrometheusInstanceHttpUrl += fmt.Sprintf("?prom_id=%s", d.Id())
	deletePrometheusInstancePath := client.Endpoint + deletePrometheusInstanceHttpUrl
	deletePrometheusInstancePath = strings.ReplaceAll(deletePrometheusInstancePath, "{project_id}", client.ProjectID)

	deletePrometheusInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deletePrometheusInstancePath, &deletePrometheusInstanceOpt)
	if err != nil && !hasErrorCode(err, prometheusInstanceNotExistsCode) {
		return diag.Errorf("error deleting AOM prometheus instance: %s", err)
	}

	return nil
}

package lb

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB POST /v2/{project_id}/elb/healthmonitors
// @API ELB GET /v2/{project_id}/elb/healthmonitors/{healthmonitor_id}
// @API ELB PUT /v2/{project_id}/elb/healthmonitors/{healthmonitor_id}
// @API ELB DELETE /v2/{project_id}/elb/healthmonitors/{healthmonitor_id}
func ResourceMonitorV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMonitorV2Create,
		ReadContext:   resourceMonitorV2Read,
		UpdateContext: resourceMonitorV2Update,
		DeleteContext: resourceMonitorV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delay": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_retries": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"url_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"http_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expected_codes": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"admin_state_up": {
				Type:       schema.TypeBool,
				Default:    true,
				Optional:   true,
				Deprecated: "tenant_id is deprecated",
			},
			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "tenant_id is deprecated",
			},
		},
	}
}

func resourceMonitorV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/elb/healthmonitors"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateMonitorBodyParams(d))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ELB monitor: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error retrieving ELB monitor: %s", err)
	}
	monitorId := utils.PathSearch("healthmonitor.id", createRespBody, "").(string)
	if monitorId == "" {
		return diag.Errorf("error creating ELB monitor: ID is not found in API response")
	}

	d.SetId(monitorId)

	return resourceMonitorV2Read(ctx, d, meta)
}

func buildCreateMonitorBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"pool_id":        d.Get("pool_id"),
		"type":           d.Get("type"),
		"delay":          d.Get("delay"),
		"timeout":        d.Get("timeout"),
		"max_retries":    d.Get("max_retries"),
		"domain_name":    utils.ValueIgnoreEmpty(d.Get("domain_name")),
		"tenant_id":      utils.ValueIgnoreEmpty(d.Get("tenant_id")),
		"url_path":       utils.ValueIgnoreEmpty(d.Get("url_path")),
		"http_method":    utils.ValueIgnoreEmpty(d.Get("http_method")),
		"expected_codes": utils.ValueIgnoreEmpty(d.Get("expected_codes")),
		"name":           utils.ValueIgnoreEmpty(d.Get("name")),
		"monitor_port":   utils.ValueIgnoreEmpty(d.Get("port")),
		"admin_state_up": utils.ValueIgnoreEmpty(d.Get("admin_state_up")),
	}

	return map[string]interface{}{"healthmonitor": bodyParams}
}

func resourceMonitorV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v2/{project_id}/elb/healthmonitors/{healthmonitor_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{healthmonitor_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ELB monitor")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("tenant_id", utils.PathSearch("healthmonitor.tenant_id", getRespBody, nil)),
		d.Set("type", utils.PathSearch("healthmonitor.type", getRespBody, nil)),
		d.Set("delay", utils.PathSearch("healthmonitor.delay", getRespBody, nil)),
		d.Set("timeout", utils.PathSearch("healthmonitor.timeout", getRespBody, nil)),
		d.Set("max_retries", utils.PathSearch("healthmonitor.max_retries", getRespBody, nil)),
		d.Set("url_path", utils.PathSearch("healthmonitor.url_path", getRespBody, nil)),
		d.Set("http_method", utils.PathSearch("healthmonitor.http_method", getRespBody, nil)),
		d.Set("expected_codes", utils.PathSearch("healthmonitor.expected_codes", getRespBody, nil)),
		d.Set("admin_state_up", utils.PathSearch("healthmonitor.admin_state_up", getRespBody, nil)),
		d.Set("name", utils.PathSearch("healthmonitor.name", getRespBody, nil)),
		d.Set("pool_id", utils.PathSearch("healthmonitor.pools[0].id", getRespBody, nil)),
		d.Set("port", utils.PathSearch("healthmonitor.monitor_port", getRespBody, nil)),
		d.Set("domain_name", utils.PathSearch("healthmonitor.domain_name", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceMonitorV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/elb/healthmonitors/{healthmonitor_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{healthmonitor_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateMonitorBodyParams(d))
	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating ELB monitor: %s", err)
	}

	return resourceMonitorV2Read(ctx, d, meta)
}

func buildUpdateMonitorBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type":           d.Get("type"),
		"delay":          d.Get("delay"),
		"timeout":        d.Get("timeout"),
		"max_retries":    d.Get("max_retries"),
		"domain_name":    utils.ValueIgnoreEmpty(d.Get("domain_name")),
		"url_path":       utils.ValueIgnoreEmpty(d.Get("url_path")),
		"http_method":    utils.ValueIgnoreEmpty(d.Get("http_method")),
		"expected_codes": utils.ValueIgnoreEmpty(d.Get("expected_codes")),
		"name":           utils.ValueIgnoreEmpty(d.Get("name")),
		"monitor_port":   utils.ValueIgnoreEmpty(d.Get("port")),
		"admin_state_up": utils.ValueIgnoreEmpty(d.Get("admin_state_up")),
	}
	return map[string]interface{}{"healthmonitor": bodyParams}
}

func resourceMonitorV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/elb/healthmonitors/{healthmonitor_id}"
		product = "elb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{healthmonitor_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting ELB monitor")
	}

	return nil
}

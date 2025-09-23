package dws

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS POST /v2/{project_id}/clusters/{cluster_id}/workload
// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/workload
func ResourceWorkLoadConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkLoadConfigurationCreate,
		ReadContext:   resourceWorkLoadConfigurationRead,
		UpdateContext: resourceWorkLoadConfigurationUpdate,
		DeleteContext: resourceWorkLoadConfigurationDelete,
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
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the DWS cluster ID.",
			},
			"workload_switch": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the workload management switch.",
			},
			"max_concurrency_num": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the maximum number of concurrent tasks on a single CN.",
			},
		},
	}
}

func resourceWorkLoadConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		clusterId = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient("dws", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	err = modifyWorkloadConfiguration(client, d, clusterId)
	if err != nil {
		return diag.Errorf("error setting workload configuration for DWS cluster(%s): %s", clusterId, err)
	}

	d.SetId(clusterId)
	return resourceWorkLoadConfigurationRead(ctx, d, meta)
}

func buildModifyConfigurationParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"workload_switch":     d.Get("workload_switch"),
		"max_concurrency_num": d.Get("max_concurrency_num"),
	}

	return map[string]interface{}{
		"workload_status": params,
	}
}

func modifyWorkloadConfiguration(client *golangsdk.ServiceClient, d *schema.ResourceData, clusterId string) error {
	httpUrl := "v2/{project_id}/clusters/{cluster_id}/workload"
	modifyPath := client.Endpoint + httpUrl
	modifyPath = strings.ReplaceAll(modifyPath, "{project_id}", client.ProjectID)
	modifyPath = strings.ReplaceAll(modifyPath, "{cluster_id}", clusterId)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildModifyConfigurationParams(d)),
	}

	_, err := client.Request("POST", modifyPath, &opts)
	return err
}

func GetWorkloadConfiguration(client *golangsdk.ServiceClient, clusterId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/clusters/{cluster_id}/workload"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func resourceWorkLoadConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	respBody, err := GetWorkloadConfiguration(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, parseWorkLoadPlanError(err), "DWS workload configuration")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("cluster_id", d.Id()),
		d.Set("workload_switch", utils.PathSearch("workload_status.workload_switch", respBody, "")),
		// The type returned by the query interface is int.
		d.Set("max_concurrency_num", strconv.Itoa(int(utils.PathSearch("workload_status.max_concurrency_num", respBody, float64(0)).(float64)))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceWorkLoadConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dws", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	if err = modifyWorkloadConfiguration(client, d, d.Id()); err != nil {
		return diag.Errorf("error updating workload configuration: %s", err)
	}

	return resourceWorkLoadConfigurationRead(ctx, d, meta)
}

func resourceWorkLoadConfigurationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only used to modify the workload configuration. Deleting this resource will
	not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

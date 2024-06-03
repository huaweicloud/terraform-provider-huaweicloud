package dws

import (
	"context"
	"encoding/json"
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

// @API DWS PUT /v2/{project_id}/clusters/{cluster_id}/workload/queues
// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/workload/queues
// @API DWS DELETE /v2/{project_id}/clusters/{cluster_id}/workload/queues
func ResourceWorkLoadQueue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkLoadQueueCreate,
		ReadContext:   resourceWorkLoadQueueRead,
		DeleteContext: resourceWorkLoadQueueDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceWorkloadQueueImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"configuration": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"resource_value": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"logical_cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceWorkLoadQueueCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/workload/queues"
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", d.Get("cluster_id").(string))
	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateWorkloadQueueBodyParams(d)),
	}

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DWS workload queue: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceWorkLoadQueueRead(ctx, d, meta)
}

func buildCreateWorkloadQueueBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"workload_queue": map[string]interface{}{
			"workload_queue_name":         d.Get("name"),
			"workload_resource_item_list": buildCreateConfigurationBodyParams(d),
			"logical_cluster_name":        utils.ValueIgnoreEmpty(d.Get("logical_cluster_name")),
		},
	}
}

func buildCreateConfigurationBodyParams(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("configuration").(*schema.Set)
	if rawParams.Len() == 0 {
		return nil
	}

	params := make([]map[string]interface{}, 0, rawParams.Len())
	for _, rawParam := range rawParams.List() {
		raw := rawParam.(map[string]interface{})
		param := map[string]interface{}{
			"resource_name":  raw["resource_name"],
			"resource_value": raw["resource_value"],
		}
		params = append(params, param)
	}

	return params
}

func resourceWorkLoadQueueRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/workload/queues"
		product = "dws"
	)

	getClient, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	getPath := getClient.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", getClient.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", d.Get("cluster_id").(string))
	getOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}

	getResp, err := getClient.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DWS workload queue")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	expression := fmt.Sprintf("workload_queue_name_list[?@=='%s']|[0]", d.Id())
	resp := utils.PathSearch(expression, getRespBody, nil)
	if resp == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "DWS workload queue")
	}

	// Due to API restrictions, only the queue name can be found in the API response.
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", d.Id()),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceWorkLoadQueueDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/workload/queues"
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl + "?workload_queue_name={name}"
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", d.Get("cluster_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{name}", d.Get("name").(string))
	// Due to API restrictions, the request body must pass in an empty JSON.
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
		JSONBody:         json.RawMessage("{}"),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting DWS workload queue: %s", err)
	}

	return nil
}

func resourceWorkloadQueueImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <cluster_id>/<name>")
	}

	d.Set("cluster_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

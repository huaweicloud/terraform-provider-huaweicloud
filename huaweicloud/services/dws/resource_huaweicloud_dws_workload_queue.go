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
// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/workload/queues{queue_name}
// @API DWS DELETE /v2/{project_id}/clusters/{cluster_id}/workload/queues
func ResourceWorkLoadQueue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkLoadQueueCreate,
		ReadContext:   resourceWorkLoadQueueRead,
		// The API only supports updating "configuration" parameter, but the update cannot be implemented due to
		// inconsistencies between the parameters on the API and the resource.
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

// GetWorkloadQueueByName is a method used to obtain resource pool information by resource pool name.
func GetWorkloadQueueByName(client *golangsdk.ServiceClient, clusterId, queueName, logicalClusterName string) (interface{}, error) {
	httpUrl := "v2/{project_id}/clusters/{cluster_id}/workload/queues/{queue_name}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterId)
	getPath = strings.ReplaceAll(getPath, "{queue_name}", queueName)

	if logicalClusterName != "" {
		getPath = fmt.Sprintf("%s?logical_cluster_name=%s", getPath, logicalClusterName)
	}

	getOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		// 1. "DWS.0047": The cluster ID is a standard UUID, the status code is 404.
		// 2. The API response includes these cases about resource not found:
		//   a. The queue name does not exist, e.g. { "workload_res_code": 1, "workload_res_str": "", "workload_queue": null }.
		//   b. Logical cluster name does not exist in logical cluster mode, e.g.
		//   { "workload_res_code": -1,"workload_res_str": "xxx", "workload_queue": null }.
		//   c. Unspecifies logical cluster name in logical cluster mode, e.g.
		//   { "workload_res_code": 1, "workload_res_str": "xxx", "workload_queue": null }.
		return nil, common.ConvertExpected500ErrInto404Err(err, "workload_res_code")
	}
	return utils.FlattenResponse(resp)
}

func resourceWorkLoadQueueRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	getRespBody, err := GetWorkloadQueueByName(client, d.Get("cluster_id").(string), d.Id(), d.Get("logical_cluster_name").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DWS workload queue")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("workload_queue.queue_name", getRespBody, nil)),
		d.Set("logical_cluster_name", utils.PathSearch("workload_queue.logical_cluster_name", getRespBody, nil)),
		// Since the API return value is inconsistent with the resource, the "configuration" parameter is ignored.
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

	if logicalClusterName, ok := d.GetOk("logical_cluster_name"); ok {
		deletePath = fmt.Sprintf("%s&logical_cluster_name=%s", deletePath, logicalClusterName)
	}
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
	importedId := d.Id()
	parts := strings.Split(d.Id(), "/")

	// In logical cluster mode, the "logical_cluster_name" parameter must not be set, otherwise the query interface will report an error.
	if len(parts) != 2 && len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import id, must be '<cluster_id>/<name>' or "+
			"'<cluster_id>/<name>/<logical_cluster_name>', but got '%s'", importedId)
	}

	mErr := multierror.Append(nil, d.Set("cluster_id", parts[0]))
	d.SetId(parts[1])

	if len(parts) == 3 {
		mErr = multierror.Append(mErr, d.Set("logical_cluster_name", parts[2]))
	}

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

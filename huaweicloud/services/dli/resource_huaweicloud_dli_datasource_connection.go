// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DLI
// ---------------------------------------------------------------

package dli

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

// @API DLI POST /v2.0/{project_id}/datasource/enhanced-connections
// @API DLI POST /v2.0/{project_id}/datasource/enhanced-connections/{id}/routes
// @API DLI GET /v2.0/{project_id}/datasource/enhanced-connections/{id}
// @API DLI PUT /v2.0/{project_id}/datasource/enhanced-connections/{id}
// @API DLI POST /v2.0/{project_id}/datasource/enhanced-connections/{id}/associate-queue
// @API DLI POST /v2.0/{project_id}/datasource/enhanced-connections/{id}/disassociate-queue
// @API DLI DELETE /v2.0/{project_id}/datasource/enhanced-connections/{id}/routes/{name}
// @API DLI DELETE /v2.0/{project_id}/datasource/enhanced-connections/{id}
func ResourceDatasourceConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatasourceConnectionCreate,
		UpdateContext: resourceDatasourceConnectionUpdate,
		ReadContext:   resourceDatasourceConnectionRead,
		DeleteContext: resourceDatasourceConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of a datasource connection.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The VPC ID of the service to be connected.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The subnet ID of the service to be connected.`,
			},
			"route_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The route table ID associated with the subnet of the service to be connected.`,
			},
			"queues": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `List of queue names that are available for datasource connections.`,
			},
			"hosts": {
				Type:        schema.TypeList,
				Elem:        datasourceConnectionHostSchema(),
				Optional:    true,
				Computed:    true,
				Description: `The user-defined host information. A maximum of 20,000 records are supported.`,
			},
			"routes": {
				Type:        schema.TypeSet,
				Elem:        datasourceConnectionRouteSchema(),
				Optional:    true,
				Computed:    true,
				Description: `List of routes.`,
			},
			"tags": common.TagsForceNewSchema(),
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The connection status.`,
			},
		},
	}
}

func datasourceConnectionHostSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The user-defined host name.`,
			},
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `IPv4 address of the host.`,
			},
		},
	}
	return &sc
}

func datasourceConnectionRouteSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The route Name`,
			},
			"cidr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The CIDR of the route.`,
			},
		},
	}
	return &sc
}

func resourceDatasourceConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createDatasourceConnection: create a DLI enhanced connection.
	var (
		createDatasourceConnectionHttpUrl = "v2.0/{project_id}/datasource/enhanced-connections"
		createDatasourceConnectionProduct = "dli"
	)
	createDatasourceConnectionClient, err := cfg.NewServiceClient(createDatasourceConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	createDatasourceConnectionPath := createDatasourceConnectionClient.Endpoint + createDatasourceConnectionHttpUrl
	createDatasourceConnectionPath = strings.ReplaceAll(createDatasourceConnectionPath, "{project_id}",
		createDatasourceConnectionClient.ProjectID)

	createDatasourceConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createDatasourceConnectionOpt.JSONBody = utils.RemoveNil(buildCreateDatasourceConnectionBodyParams(d))
	createDatasourceConnectionResp, err := createDatasourceConnectionClient.Request("POST",
		createDatasourceConnectionPath, &createDatasourceConnectionOpt)

	if err != nil {
		return diag.Errorf("error creating DatasourceConnection: %s", err)
	}

	createDatasourceConnectionRespBody, err := utils.FlattenResponse(createDatasourceConnectionResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if !utils.PathSearch("is_success", createDatasourceConnectionRespBody, true).(bool) {
		return diag.Errorf("unable to create the enhanced connection: %s",
			utils.PathSearch("message", createDatasourceConnectionRespBody, "Message Not Found"))
	}

	connectionId := utils.PathSearch("connection_id", createDatasourceConnectionRespBody, "").(string)
	if connectionId == "" {
		return diag.Errorf("unable to find the connection ID of the DLI Data Source from the API response")
	}
	d.SetId(connectionId)

	// add routes
	if v, ok := d.GetOk("routes"); ok {
		err = addRoutes(createDatasourceConnectionClient, d.Id(), v.(*schema.Set))
		if err != nil {
			return diag.Errorf("error adding routes to DatasourceConnection: %s", d.Id())
		}
	}

	return resourceDatasourceConnectionRead(ctx, d, meta)
}

func buildCreateDatasourceConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":            utils.ValueIgnoreEmpty(d.Get("name")),
		"dest_vpc_id":     utils.ValueIgnoreEmpty(d.Get("vpc_id")),
		"dest_network_id": utils.ValueIgnoreEmpty(d.Get("subnet_id")),
		"routetable_id":   utils.ValueIgnoreEmpty(d.Get("route_table_id")),
		"queues":          d.Get("queues").(*schema.Set).List(),
		"hosts":           buildCreateDatasourceConnectionRequestBodyHost(d.Get("hosts")),
		"tags":            utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}
	return bodyParams
}

func buildCreateDatasourceConnectionRequestBodyHost(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name": utils.ValueIgnoreEmpty(raw["name"]),
				"ip":   utils.ValueIgnoreEmpty(raw["ip"]),
			}
		}
		return rst
	}
	return nil
}

func getConnectionById(client *golangsdk.ServiceClient, connectionId string) (interface{}, error) {
	var (
		httpUrl = "v2.0/{project_id}/datasource/enhanced-connections/{connection_id}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{connection_id}", connectionId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceDatasourceConnectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	getDatasourceConnectionClient, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	getDatasourceConnectionRespBody, err := getConnectionById(getDatasourceConnectionClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DatasourceConnection")
	}
	if !utils.PathSearch("is_success", getDatasourceConnectionRespBody, true).(bool) {
		return diag.Errorf("unable to query the enhanced connection: %s",
			utils.PathSearch("message", getDatasourceConnectionRespBody, "Message Not Found"))
	}

	if utils.PathSearch("status", getDatasourceConnectionRespBody, "") == "DELETED" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "the datasource connection has been deleted")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getDatasourceConnectionRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("dest_vpc_id", getDatasourceConnectionRespBody, nil)),
		d.Set("subnet_id", utils.PathSearch("dest_network_id", getDatasourceConnectionRespBody, nil)),
		d.Set("queues", utils.PathSearch("available_queue_info[*].name", getDatasourceConnectionRespBody, nil)),
		d.Set("hosts", flattenGetDatasourceConnectionResponseBodyHost(getDatasourceConnectionRespBody)),
		d.Set("routes", flattenGetDatasourceConnectionResponseBodyRoute(getDatasourceConnectionRespBody)),
		d.Set("status", utils.PathSearch("status", getDatasourceConnectionRespBody, nil)),
		d.Set("tags", d.Get("tags")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetDatasourceConnectionResponseBodyHost(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("hosts", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
			"ip":   utils.PathSearch("ip", v, nil),
		})
	}
	return rst
}

func flattenGetDatasourceConnectionResponseBodyRoute(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("routes", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
			"cidr": utils.PathSearch("cidr", v, nil),
		})
	}
	return rst
}

func resourceDatasourceConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateDatasourceConnectionHostsChanges := []string{
		"hosts",
	}

	if d.HasChanges(updateDatasourceConnectionHostsChanges...) {
		// updateDatasourceConnectionHosts: update hosts
		var (
			updateDatasourceConnectionHostsHttpUrl = "v2.0/{project_id}/datasource/enhanced-connections/{id}"
			updateDatasourceConnectionHostsProduct = "dli"
		)
		updateDatasourceConnectionHostsClient, err := cfg.NewServiceClient(updateDatasourceConnectionHostsProduct, region)
		if err != nil {
			return diag.Errorf("error creating DLI Client: %s", err)
		}

		updateDatasourceConnectionHostsPath := updateDatasourceConnectionHostsClient.Endpoint + updateDatasourceConnectionHostsHttpUrl
		updateDatasourceConnectionHostsPath = strings.ReplaceAll(updateDatasourceConnectionHostsPath, "{project_id}",
			updateDatasourceConnectionHostsClient.ProjectID)
		updateDatasourceConnectionHostsPath = strings.ReplaceAll(updateDatasourceConnectionHostsPath, "{id}", d.Id())

		updateDatasourceConnectionHostsOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateDatasourceConnectionHostsOpt.JSONBody = utils.RemoveNil(buildUpdateDatasourceConnectionHostsBodyParams(d, cfg))
		resp, err := updateDatasourceConnectionHostsClient.Request("PUT", updateDatasourceConnectionHostsPath, &updateDatasourceConnectionHostsOpt)
		if err != nil {
			return diag.Errorf("error updating DatasourceConnection: %s", err)
		}
		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}
		if !utils.PathSearch("is_success", respBody, true).(bool) {
			return diag.Errorf("unable to update the hosts configuration: %s",
				utils.PathSearch("message", respBody, "Message Not Found"))
		}
	}

	// updateDatasourceConnectionQueues: update queues
	updateDatasourceConnectionQueuesChanges := []string{
		"queues",
	}

	if d.HasChanges(updateDatasourceConnectionQueuesChanges...) {
		o, n := d.GetChange("queues")

		addRaws := n.(*schema.Set).Difference(o.(*schema.Set))
		delRaws := o.(*schema.Set).Difference(n.(*schema.Set))
		if addRaws.Len() > 0 {
			var (
				updateDatasourceConnectionQueuesHttpUrl = "v2.0/{project_id}/datasource/enhanced-connections/{id}/associate-queue"
				updateDatasourceConnectionQueuesProduct = "dli"
			)
			updateDatasourceConnectionQueuesClient, err := cfg.NewServiceClient(updateDatasourceConnectionQueuesProduct, region)
			if err != nil {
				return diag.Errorf("error creating DLI Client: %s", err)
			}

			updateDatasourceConnectionQueuesPath := updateDatasourceConnectionQueuesClient.Endpoint + updateDatasourceConnectionQueuesHttpUrl
			updateDatasourceConnectionQueuesPath = strings.ReplaceAll(updateDatasourceConnectionQueuesPath, "{project_id}",
				updateDatasourceConnectionQueuesClient.ProjectID)
			updateDatasourceConnectionQueuesPath = strings.ReplaceAll(updateDatasourceConnectionQueuesPath, "{id}", d.Id())

			updateDatasourceConnectionQueuesOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			updateDatasourceConnectionQueuesOpt.JSONBody = buildDatasourceConnectionQueuesBodyParams(addRaws)
			requestResp, err := updateDatasourceConnectionQueuesClient.Request("POST", updateDatasourceConnectionQueuesPath,
				&updateDatasourceConnectionQueuesOpt)
			if err != nil {
				return diag.Errorf("error updating DatasourceConnection: %s", err)
			}
			respBody, err := utils.FlattenResponse(requestResp)
			if err != nil {
				return diag.FromErr(err)
			}
			if !utils.PathSearch("is_success", respBody, true).(bool) {
				return diag.Errorf("unable to associate the queues: %s",
					utils.PathSearch("message", respBody, "Message Not Found"))
			}
		}

		if delRaws.Len() > 0 {
			var (
				deleteDatasourceConnectionQueuesHttpUrl = "v2.0/{project_id}/datasource/enhanced-connections/{id}/disassociate-queue"
				deleteDatasourceConnectionQueuesProduct = "dli"
			)
			deleteDatasourceConnectionQueuesClient, err := cfg.NewServiceClient(deleteDatasourceConnectionQueuesProduct, region)
			if err != nil {
				return diag.Errorf("error creating DLI Client: %s", err)
			}

			deleteDatasourceConnectionQueuesPath := deleteDatasourceConnectionQueuesClient.Endpoint + deleteDatasourceConnectionQueuesHttpUrl
			deleteDatasourceConnectionQueuesPath = strings.ReplaceAll(deleteDatasourceConnectionQueuesPath, "{project_id}",
				deleteDatasourceConnectionQueuesClient.ProjectID)
			deleteDatasourceConnectionQueuesPath = strings.ReplaceAll(deleteDatasourceConnectionQueuesPath, "{id}", d.Id())

			deleteDatasourceConnectionQueuesOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			deleteDatasourceConnectionQueuesOpt.JSONBody = buildDatasourceConnectionQueuesBodyParams(delRaws)
			requestResp, err := deleteDatasourceConnectionQueuesClient.Request("POST", deleteDatasourceConnectionQueuesPath,
				&deleteDatasourceConnectionQueuesOpt)
			if err != nil {
				return diag.Errorf("error updating DatasourceConnection: %s", err)
			}
			respBody, err := utils.FlattenResponse(requestResp)
			if err != nil {
				return diag.FromErr(err)
			}
			if !utils.PathSearch("is_success", respBody, true).(bool) {
				return diag.Errorf("unable to disassociate the queues: %s",
					utils.PathSearch("message", respBody, "Message Not Found"))
			}
		}
	}

	// updateDatasourceConnectionRoutes: update routes
	updateDatasourceConnectionRoutesChanges := []string{
		"routes",
	}

	if d.HasChanges(updateDatasourceConnectionRoutesChanges...) {
		connectionRouteClient, err := cfg.NewServiceClient("dli", region)
		if err != nil {
			return diag.Errorf("error creating DLI Client: %s", err)
		}

		o, n := d.GetChange("routes")
		addRaws := n.(*schema.Set).Difference(o.(*schema.Set))
		delRaws := o.(*schema.Set).Difference(n.(*schema.Set))

		if addRaws.Len() > 0 {
			err := addRoutes(connectionRouteClient, d.Id(), addRaws)
			if err != nil {
				return diag.Errorf("error updating DatasourceConnection: %s", err)
			}
		}

		if delRaws.Len() > 0 {
			err := removeRoutes(connectionRouteClient, d.Id(), delRaws)
			if err != nil {
				return diag.Errorf("error updating DatasourceConnection: %s", err)
			}
		}
	}
	return resourceDatasourceConnectionRead(ctx, d, meta)
}

func addRoutes(connectionRouteClient *golangsdk.ServiceClient, id string, addRaws *schema.Set) error {
	var (
		addConnectionRouteHttpUrl = "v2.0/{project_id}/datasource/enhanced-connections/{id}/routes"
	)

	addConnectionRoutePath := connectionRouteClient.Endpoint + addConnectionRouteHttpUrl
	addConnectionRoutePath = strings.ReplaceAll(addConnectionRoutePath, "{project_id}", connectionRouteClient.ProjectID)
	addConnectionRoutePath = strings.ReplaceAll(addConnectionRoutePath, "{id}", id)

	addConnectionRouteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	for _, params := range addRaws.List() {
		addConnectionRouteOpt.JSONBody = params
		requestResp, err := connectionRouteClient.Request("POST", addConnectionRoutePath, &addConnectionRouteOpt)
		if err != nil {
			return err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return err
		}
		if !utils.PathSearch("is_success", respBody, true).(bool) {
			return fmt.Errorf("unable to add the routes: %s",
				utils.PathSearch("message", respBody, "Message Not Found"))
		}
	}
	return nil
}

func removeRoutes(connectionRouteClient *golangsdk.ServiceClient, id string, raws *schema.Set) error {
	for _, params := range raws.List() {
		var (
			removeDatasourceConnectionRoutesHttpUrl = "v2.0/{project_id}/datasource/enhanced-connections/{id}/routes/{name}"
		)
		removeDatasourceConnectionRoutesPath := connectionRouteClient.Endpoint + removeDatasourceConnectionRoutesHttpUrl
		removeDatasourceConnectionRoutesPath = strings.ReplaceAll(removeDatasourceConnectionRoutesPath, "{project_id}",
			connectionRouteClient.ProjectID)
		removeDatasourceConnectionRoutesPath = strings.ReplaceAll(removeDatasourceConnectionRoutesPath, "{id}", id)
		removeDatasourceConnectionRoutesPath = strings.ReplaceAll(removeDatasourceConnectionRoutesPath, "{name}",
			fmt.Sprintf("%v", utils.PathSearch("name", params, nil)))

		removeDatasourceConnectionRoutesOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		requestResp, err := connectionRouteClient.Request("DELETE", removeDatasourceConnectionRoutesPath,
			&removeDatasourceConnectionRoutesOpt)
		if err != nil {
			return err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return err
		}
		if !utils.PathSearch("is_success", respBody, true).(bool) {
			return fmt.Errorf("unable to remove the routes: %s",
				utils.PathSearch("message", respBody, "Message Not Found"))
		}
	}
	return nil
}

func buildDatasourceConnectionQueuesBodyParams(v *schema.Set) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"queues": v.List(),
	}
	return bodyParams
}

func buildUpdateDatasourceConnectionHostsBodyParams(d *schema.ResourceData, _ *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"hosts": buildUpdateDatasourceConnectionHostsRequestBodyHost(d.Get("hosts")),
	}
	return bodyParams
}

func buildUpdateDatasourceConnectionHostsRequestBodyHost(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name": utils.ValueIgnoreEmpty(raw["name"]),
				"ip":   utils.ValueIgnoreEmpty(raw["ip"]),
			}
		}
		return rst
	}
	return nil
}

func resourceDatasourceConnectionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteDatasourceConnectionHttpUrl = "v2.0/{project_id}/datasource/enhanced-connections/{id}"
		deleteDatasourceConnectionProduct = "dli"
	)
	deleteDatasourceConnectionClient, err := cfg.NewServiceClient(deleteDatasourceConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	deleteDatasourceConnectionPath := deleteDatasourceConnectionClient.Endpoint + deleteDatasourceConnectionHttpUrl
	deleteDatasourceConnectionPath = strings.ReplaceAll(deleteDatasourceConnectionPath, "{project_id}", deleteDatasourceConnectionClient.ProjectID)
	deleteDatasourceConnectionPath = strings.ReplaceAll(deleteDatasourceConnectionPath, "{id}", d.Id())

	deleteDatasourceConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	requestResp, err := deleteDatasourceConnectionClient.Request("DELETE", deleteDatasourceConnectionPath, &deleteDatasourceConnectionOpt)
	if err != nil {
		return diag.Errorf("error deleting DatasourceConnection: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if !utils.PathSearch("is_success", respBody, true).(bool) {
		return diag.Errorf("unable to delete the enhanced connection: %s",
			utils.PathSearch("message", respBody, "Message Not Found"))
	}

	return nil
}

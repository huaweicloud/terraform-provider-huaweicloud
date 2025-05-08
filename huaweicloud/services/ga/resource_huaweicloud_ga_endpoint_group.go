// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product GA
// ---------------------------------------------------------------

package ga

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GA POST /v1/endpoint-groups
// @API GA GET /v1/endpoint-groups/{endpoint_group_id}
// @API GA PUT /v1/endpoint-groups/{endpoint_group_id}
// @API GA DELETE /v1/endpoint-groups/{endpoint_group_id}
func ResourceEndpointGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointGroupCreate,
		UpdateContext: resourceEndpointGroupUpdate,
		ReadContext:   resourceEndpointGroupRead,
		DeleteContext: resourceEndpointGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the endpoint group name.`,
			},
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the region where the endpoint group belongs.`,
			},
			"listeners": {
				Type:        schema.TypeList,
				MaxItems:    1,
				MinItems:    1,
				Elem:        EndpointGroupIdSchema(),
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the listeners associated with the endpoint group.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the information about the endpoint group.`,
			},
			"traffic_dial_percentage": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the percentage of traffic distributed to the endpoint group.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the provisioning status.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies when the endpoint group was created.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies when the endpoint group was updated.`,
			},
			"frozen_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The frozen details of cloud services or resources.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The status of a cloud service or resource.`,
						},
						"effect": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The status of the resource after being forzen.`,
						},
						"scene": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The service scenario.`,
						},
					},
				},
			},
		},
	}
}

func EndpointGroupIdSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the associated listener.`,
			},
		},
	}
	return &sc
}

func resourceEndpointGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		httpUrl = "v1/endpoint-groups"
		product = "ga"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateEndpointGroupBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating GA endpoint group: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	groupId := utils.PathSearch("endpoint_group.id", respBody, "").(string)
	if groupId == "" {
		return diag.Errorf("error creating GA endpoint group: ID is not found in API response")
	}
	d.SetId(groupId)

	err = createEndpointGroupWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the GA endpoint group (%s) creation to complete: %s", d.Id(), err)
	}
	return resourceEndpointGroupRead(ctx, d, meta)
}

func buildCreateEndpointGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"endpoint_group": map[string]interface{}{
			"description":             utils.ValueIgnoreEmpty(d.Get("description")),
			"listeners":               buildCreateEndpointGroupRequestBodyId(d.Get("listeners")),
			"name":                    utils.ValueIgnoreEmpty(d.Get("name")),
			"region_id":               utils.ValueIgnoreEmpty(d.Get("region_id")),
			"traffic_dial_percentage": utils.ValueIgnoreEmpty(d.Get("traffic_dial_percentage")),
		},
	}
	return bodyParams
}

func buildCreateEndpointGroupRequestBodyId(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"id": utils.ValueIgnoreEmpty(raw["id"]),
			}
		}
		return rst
	}
	return nil
}

func createEndpointGroupWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		httpUrl          = "v1/endpoint-groups/{endpoint_group_id}"
		product          = "ga"
		targetStatus     = []string{"ACTIVE"}
		unexpectedStatus = []string{"ERROR"}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{endpoint_group_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Request("GET", requestPath, &requestOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			respBody, err := utils.FlattenResponse(resp)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch(`endpoint_group.status`, respBody, "").(string)
			if utils.StrSliceContains(targetStatus, status) {
				return respBody, "COMPLETED", nil
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return respBody, status, nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	return err
}

func resourceEndpointGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/endpoint-groups/{endpoint_group_id}"
		product = "ga"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{endpoint_group_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GA endpoint group")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("created_at", utils.PathSearch("endpoint_group.created_at", respBody, nil)),
		d.Set("description", utils.PathSearch("endpoint_group.description", respBody, nil)),
		d.Set("listeners", flattenGetEndpointGroupResponseBodyId(respBody)),
		d.Set("name", utils.PathSearch("endpoint_group.name", respBody, nil)),
		d.Set("region_id", utils.PathSearch("endpoint_group.region_id", respBody, nil)),
		d.Set("status", utils.PathSearch("endpoint_group.status", respBody, nil)),
		d.Set("traffic_dial_percentage", utils.PathSearch("endpoint_group.traffic_dial_percentage", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("endpoint_group.updated_at", respBody, nil)),
		d.Set("frozen_info", flattenEndpointGroupFrozenInfo(utils.PathSearch("endpoint_group.frozen_info", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetEndpointGroupResponseBodyId(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("endpoint_group.listeners", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id": utils.PathSearch("id", v, nil),
		})
	}
	return rst
}

func flattenEndpointGroupFrozenInfo(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	frozenInfo := map[string]interface{}{
		"status": utils.PathSearch("status", resp, nil),
		"effect": utils.PathSearch("effect", resp, nil),
		"scene":  utils.PathSearch("scene", resp, []string{}),
	}

	return []map[string]interface{}{frozenInfo}
}

func resourceEndpointGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	updateEndpointGroupHasChanges := []string{
		"description",
		"name",
		"traffic_dial_percentage",
	}

	if d.HasChanges(updateEndpointGroupHasChanges...) {
		var (
			httpUrl = "v1/endpoint-groups/{endpoint_group_id}"
			product = "ga"
		)
		client, err := conf.NewServiceClient(product, region)
		if err != nil {
			return diag.Errorf("error creating GA client: %s", err)
		}

		requestPath := client.Endpoint + httpUrl
		requestPath = strings.ReplaceAll(requestPath, "{endpoint_group_id}", d.Id())
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateEndpointGroupBodyParams(d)),
		}

		_, err = client.Request("PUT", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating GA endpoint group: %s", err)
		}
		err = updateEndpointGroupWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the GA endpoint group (%s) update to complete: %s", d.Id(), err)
		}
	}
	return resourceEndpointGroupRead(ctx, d, meta)
}

func buildUpdateEndpointGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"endpoint_group": map[string]interface{}{
			"description":             d.Get("description"),
			"name":                    utils.ValueIgnoreEmpty(d.Get("name")),
			"traffic_dial_percentage": utils.ValueIgnoreEmpty(d.Get("traffic_dial_percentage")),
		},
	}
	return bodyParams
}

func updateEndpointGroupWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		httpUrl          = "v1/endpoint-groups/{endpoint_group_id}"
		product          = "ga"
		targetStatus     = []string{"ACTIVE"}
		unexpectedStatus = []string{"ERROR"}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{endpoint_group_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Request("GET", requestPath, &requestOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			respBody, err := utils.FlattenResponse(resp)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch(`endpoint_group.status`, respBody, "").(string)
			if utils.StrSliceContains(targetStatus, status) {
				return respBody, "COMPLETED", nil
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return respBody, status, nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	return err
}

func resourceEndpointGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		httpUrl = "v1/endpoint-groups/{endpoint_group_id}"
		product = "ga"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{endpoint_group_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting GA endpoint group: %s", err)
	}

	err = deleteEndpointGroupWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the GA endpoint group (%s) delete to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteEndpointGroupWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		httpUrl          = "v1/endpoint-groups/{endpoint_group_id}"
		product          = "ga"
		unexpectedStatus = []string{"ERROR"}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{endpoint_group_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Request("GET", requestPath, &requestOpt)
			if err != nil {
				var errDefault404 golangsdk.ErrDefault404
				if errors.As(err, &errDefault404) {
					// When the error code is `404`, the value of respBody is nil, and a non-null value is returned to
					// avoid continuing the loop check.
					return "Resource Not Found", "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			respBody, err := utils.FlattenResponse(resp)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch(`endpoint_group.status`, respBody, "").(string)
			if utils.StrSliceContains(unexpectedStatus, status) {
				return respBody, status, nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	return err
}

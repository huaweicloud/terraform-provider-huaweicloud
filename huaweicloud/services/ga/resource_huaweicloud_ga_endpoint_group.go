// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product GA
// ---------------------------------------------------------------

package ga

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

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
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	// createEndpointGroup: Create a GA endpoint group.
	var (
		createEndpointGroupHttpUrl = "v1/endpoint-groups"
		createEndpointGroupProduct = "ga"
	)
	createEndpointGroupClient, err := conf.NewServiceClient(createEndpointGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating EndpointGroup Client: %s", err)
	}

	createEndpointGroupPath := createEndpointGroupClient.Endpoint + createEndpointGroupHttpUrl

	createEndpointGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createEndpointGroupOpt.JSONBody = utils.RemoveNil(buildCreateEndpointGroupBodyParams(d))
	createEndpointGroupResp, err := createEndpointGroupClient.Request("POST", createEndpointGroupPath, &createEndpointGroupOpt)
	if err != nil {
		return diag.Errorf("error creating EndpointGroup: %s", err)
	}

	createEndpointGroupRespBody, err := utils.FlattenResponse(createEndpointGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("endpoint_group.id", createEndpointGroupRespBody)
	if err != nil {
		return diag.Errorf("error creating EndpointGroup: ID is not found in API response")
	}
	d.SetId(id.(string))

	err = createEndpointGroupWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Create of EndpointGroup (%s) to complete: %s", d.Id(), err)
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
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// createEndpointGroupWaiting: missing operation notes
			var (
				createEndpointGroupWaitingHttpUrl = "v1/endpoint-groups/{endpoint_group_id}"
				createEndpointGroupWaitingProduct = "ga"
			)
			createEndpointGroupWaitingClient, err := config.NewServiceClient(createEndpointGroupWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating EndpointGroup Client: %s", err)
			}

			createEndpointGroupWaitingPath := createEndpointGroupWaitingClient.Endpoint + createEndpointGroupWaitingHttpUrl
			createEndpointGroupWaitingPath = strings.ReplaceAll(createEndpointGroupWaitingPath, "{endpoint_group_id}", d.Id())

			createEndpointGroupWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			createEndpointGroupWaitingResp, err := createEndpointGroupWaitingClient.Request("GET",
				createEndpointGroupWaitingPath, &createEndpointGroupWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createEndpointGroupWaitingRespBody, err := utils.FlattenResponse(createEndpointGroupWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`endpoint_group.status`, createEndpointGroupWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `endpoint_group.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"ACTIVE",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createEndpointGroupWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return createEndpointGroupWaitingRespBody, status, nil
			}

			return createEndpointGroupWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceEndpointGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	// getEndpointGroup: Query the GA Endpoint Group detail
	var (
		getEndpointGroupHttpUrl = "v1/endpoint-groups/{endpoint_group_id}"
		getEndpointGroupProduct = "ga"
	)
	getEndpointGroupClient, err := conf.NewServiceClient(getEndpointGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating EndpointGroup Client: %s", err)
	}

	getEndpointGroupPath := getEndpointGroupClient.Endpoint + getEndpointGroupHttpUrl
	getEndpointGroupPath = strings.ReplaceAll(getEndpointGroupPath, "{endpoint_group_id}", d.Id())

	getEndpointGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getEndpointGroupResp, err := getEndpointGroupClient.Request("GET", getEndpointGroupPath, &getEndpointGroupOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving EndpointGroup")
	}

	getEndpointGroupRespBody, err := utils.FlattenResponse(getEndpointGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("created_at", utils.PathSearch("endpoint_group.created_at", getEndpointGroupRespBody, nil)),
		d.Set("description", utils.PathSearch("endpoint_group.description", getEndpointGroupRespBody, nil)),
		d.Set("listeners", flattenGetEndpointGroupResponseBodyId(getEndpointGroupRespBody)),
		d.Set("name", utils.PathSearch("endpoint_group.name", getEndpointGroupRespBody, nil)),
		d.Set("region_id", utils.PathSearch("endpoint_group.region_id", getEndpointGroupRespBody, nil)),
		d.Set("status", utils.PathSearch("endpoint_group.status", getEndpointGroupRespBody, nil)),
		d.Set("traffic_dial_percentage", utils.PathSearch("endpoint_group.traffic_dial_percentage", getEndpointGroupRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("endpoint_group.updated_at", getEndpointGroupRespBody, nil)),
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

func resourceEndpointGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	updateEndpointGrouphasChanges := []string{
		"description",
		"name",
		"traffic_dial_percentage",
	}

	if d.HasChanges(updateEndpointGrouphasChanges...) {
		// updateEndpointGroup: Update the configuration of GA Endpoint Group
		var (
			updateEndpointGroupHttpUrl = "v1/endpoint-groups/{endpoint_group_id}"
			updateEndpointGroupProduct = "ga"
		)
		updateEndpointGroupClient, err := conf.NewServiceClient(updateEndpointGroupProduct, region)
		if err != nil {
			return diag.Errorf("error creating EndpointGroup Client: %s", err)
		}

		updateEndpointGroupPath := updateEndpointGroupClient.Endpoint + updateEndpointGroupHttpUrl
		updateEndpointGroupPath = strings.ReplaceAll(updateEndpointGroupPath, "{endpoint_group_id}", d.Id())

		updateEndpointGroupOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateEndpointGroupOpt.JSONBody = utils.RemoveNil(buildUpdateEndpointGroupBodyParams(d))
		_, err = updateEndpointGroupClient.Request("PUT", updateEndpointGroupPath, &updateEndpointGroupOpt)
		if err != nil {
			return diag.Errorf("error updating EndpointGroup: %s", err)
		}
		err = updateEndpointGroupWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the Update of EndpointGroup (%s) to complete: %s", d.Id(), err)
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
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// updateEndpointGroupWaiting: missing operation notes
			var (
				updateEndpointGroupWaitingHttpUrl = "v1/endpoint-groups/{endpoint_group_id}"
				updateEndpointGroupWaitingProduct = "ga"
			)
			updateEndpointGroupWaitingClient, err := config.NewServiceClient(updateEndpointGroupWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating EndpointGroup Client: %s", err)
			}

			updateEndpointGroupWaitingPath := updateEndpointGroupWaitingClient.Endpoint + updateEndpointGroupWaitingHttpUrl
			updateEndpointGroupWaitingPath = strings.ReplaceAll(updateEndpointGroupWaitingPath, "{endpoint_group_id}", d.Id())

			updateEndpointGroupWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			updateEndpointGroupWaitingResp, err := updateEndpointGroupWaitingClient.Request("GET",
				updateEndpointGroupWaitingPath, &updateEndpointGroupWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			updateEndpointGroupWaitingRespBody, err := utils.FlattenResponse(updateEndpointGroupWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`endpoint_group.status`, updateEndpointGroupWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `endpoint_group.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"ACTIVE",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return updateEndpointGroupWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return updateEndpointGroupWaitingRespBody, status, nil
			}

			return updateEndpointGroupWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceEndpointGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	// deleteEndpointGroup: Delete an existing GA Endpoint Group
	var (
		deleteEndpointGroupHttpUrl = "v1/endpoint-groups/{endpoint_group_id}"
		deleteEndpointGroupProduct = "ga"
	)
	deleteEndpointGroupClient, err := conf.NewServiceClient(deleteEndpointGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating EndpointGroup Client: %s", err)
	}

	deleteEndpointGroupPath := deleteEndpointGroupClient.Endpoint + deleteEndpointGroupHttpUrl
	deleteEndpointGroupPath = strings.ReplaceAll(deleteEndpointGroupPath, "{endpoint_group_id}", d.Id())

	deleteEndpointGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteEndpointGroupClient.Request("DELETE", deleteEndpointGroupPath, &deleteEndpointGroupOpt)
	if err != nil {
		return diag.Errorf("error deleting EndpointGroup: %s", err)
	}

	err = deleteEndpointGroupWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the Delete of EndpointGroup (%s) to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteEndpointGroupWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// deleteEndpointGroupWaiting: missing operation notes
			var (
				deleteEndpointGroupWaitingHttpUrl = "v1/endpoint-groups/{endpoint_group_id}"
				deleteEndpointGroupWaitingProduct = "ga"
			)
			deleteEndpointGroupWaitingClient, err := config.NewServiceClient(deleteEndpointGroupWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating EndpointGroup Client: %s", err)
			}

			deleteEndpointGroupWaitingPath := deleteEndpointGroupWaitingClient.Endpoint + deleteEndpointGroupWaitingHttpUrl
			deleteEndpointGroupWaitingPath = strings.ReplaceAll(deleteEndpointGroupWaitingPath, "{endpoint_group_id}", d.Id())

			deleteEndpointGroupWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			deleteEndpointGroupWaitingResp, err := deleteEndpointGroupWaitingClient.Request("GET",
				deleteEndpointGroupWaitingPath, &deleteEndpointGroupWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
					return "Resource Not Found", "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			deleteEndpointGroupWaitingRespBody, err := utils.FlattenResponse(deleteEndpointGroupWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`endpoint_group.status`, deleteEndpointGroupWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `endpoint_group.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return deleteEndpointGroupWaitingRespBody, status, nil
			}

			return deleteEndpointGroupWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

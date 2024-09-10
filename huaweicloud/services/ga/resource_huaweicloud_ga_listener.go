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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GA POST /v1/listeners
// @API GA DELETE /v1/listeners/{listener_id}
// @API GA GET /v1/listeners/{listener_id}
// @API GA PUT /v1/listeners/{listener_id}
// @API GA POST /v1/{resource_type}/{resource_id}/tags/create
// @API GA DELETE /v1/{resource_type}/{resource_id}/tags/delete
func ResourceListener() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceListenerCreate,
		UpdateContext: resourceListenerUpdate,
		ReadContext:   resourceListenerRead,
		DeleteContext: resourceListenerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the global accelerator associated with the listener.`,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `|-
                    Specifies the listener name. The name can contain 1 to 64 characters.
                    Only letters, digits, and hyphens (-) are allowed.`,
			},
			"port_ranges": {
				Type:        schema.TypeList,
				MaxItems:    10,
				MinItems:    1,
				Elem:        ListenerPortRangeSchema(),
				Required:    true,
				Description: `Specifies the port range used by the listener.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the protocol used by the listener to forward requests.`,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP",
				}, false),
			},
			"client_affinity": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `|-
                    Specifies the client affinity. The value can be one of the following:
                      - **Source IP address**: Requests from the same IP address are forwarded to the same endpoint.
                      - **NONE**: Requests are evenly forwarded across the endpoints.`,
				ValidateFunc: validation.StringInSlice([]string{
					"SOURCE_IP", "NONE",
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `|-
                    Specifies the information about the listener. The value can contain 0 to 255 characters.
                    The following characters are not allowed: <>`,
			},
			"tags": common.TagsSchema(),

			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `|-
                    Specifies the provisioning status. The value can be one of the following:
                      - **ACTIVE**: The resource is running.
                      - **PENDING**: The status is to be determined.
                      - **ERROR**: Failed to create the resource.
                      - **DELETING**: The resource is being deleted.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies when the listener was created.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies when the listener was updated.`,
			},
		},
	}
}

func ListenerPortRangeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"from_port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the start port number.`,
			},
			"to_port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the end port number.`,
			},
		},
	}
	return &sc
}

func resourceListenerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	// createListener: Create a GA Listener.
	var (
		createListenerHttpUrl = "v1/listeners"
		createListenerProduct = "ga"
	)
	createListenerClient, err := conf.NewServiceClient(createListenerProduct, region)
	if err != nil {
		return diag.Errorf("error creating Listener Client: %s", err)
	}

	createListenerPath := createListenerClient.Endpoint + createListenerHttpUrl

	createListenerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createListenerOpt.JSONBody = utils.RemoveNil(buildCreateListenerBodyParams(d))
	createListenerResp, err := createListenerClient.Request("POST", createListenerPath, &createListenerOpt)
	if err != nil {
		return diag.Errorf("error creating Listener: %s", err)
	}

	createListenerRespBody, err := utils.FlattenResponse(createListenerResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("listener.id", createListenerRespBody)
	if err != nil {
		return diag.Errorf("error creating Listener: ID is not found in API response")
	}
	d.SetId(id.(string))

	err = createListenerWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Create of Listener (%s) to complete: %s", d.Id(), err)
	}
	return resourceListenerRead(ctx, d, meta)
}

func buildCreateListenerBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"listener": map[string]interface{}{
			"accelerator_id":  utils.ValueIgnoreEmpty(d.Get("accelerator_id")),
			"client_affinity": utils.ValueIgnoreEmpty(d.Get("client_affinity")),
			"description":     utils.ValueIgnoreEmpty(d.Get("description")),
			"name":            utils.ValueIgnoreEmpty(d.Get("name")),
			"port_ranges":     buildCreateListenerRequestBodyPortRange(d.Get("port_ranges")),
			"protocol":        utils.ValueIgnoreEmpty(d.Get("protocol")),
			"tags":            utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
		},
	}
	return bodyParams
}

func buildCreateListenerRequestBodyPortRange(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"from_port": utils.ValueIgnoreEmpty(raw["from_port"]),
				"to_port":   utils.ValueIgnoreEmpty(raw["to_port"]),
			}
		}
		return rst
	}
	return nil
}

func createListenerWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// createListenerWaiting: missing operation notes
			var (
				createListenerWaitingHttpUrl = "v1/listeners/{listener_id}"
				createListenerWaitingProduct = "ga"
			)
			createListenerWaitingClient, err := config.NewServiceClient(createListenerWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating Listener Client: %s", err)
			}

			createListenerWaitingPath := createListenerWaitingClient.Endpoint + createListenerWaitingHttpUrl
			createListenerWaitingPath = strings.ReplaceAll(createListenerWaitingPath, "{listener_id}", d.Id())

			createListenerWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			createListenerWaitingResp, err := createListenerWaitingClient.Request("GET", createListenerWaitingPath, &createListenerWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createListenerWaitingRespBody, err := utils.FlattenResponse(createListenerWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`listener.status`, createListenerWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `listener.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"ACTIVE",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createListenerWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return createListenerWaitingRespBody, status, nil
			}

			return createListenerWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceListenerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	// getListener: Query the GA Listener detail
	var (
		getListenerHttpUrl = "v1/listeners/{listener_id}"
		getListenerProduct = "ga"
	)
	getListenerClient, err := conf.NewServiceClient(getListenerProduct, region)
	if err != nil {
		return diag.Errorf("error creating Listener Client: %s", err)
	}

	getListenerPath := getListenerClient.Endpoint + getListenerHttpUrl
	getListenerPath = strings.ReplaceAll(getListenerPath, "{listener_id}", d.Id())

	getListenerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getListenerResp, err := getListenerClient.Request("GET", getListenerPath, &getListenerOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Listener")
	}

	getListenerRespBody, err := utils.FlattenResponse(getListenerResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("accelerator_id", utils.PathSearch("listener.accelerator_id", getListenerRespBody, nil)),
		d.Set("client_affinity", utils.PathSearch("listener.client_affinity", getListenerRespBody, nil)),
		d.Set("created_at", utils.PathSearch("listener.created_at", getListenerRespBody, nil)),
		d.Set("description", utils.PathSearch("listener.description", getListenerRespBody, nil)),
		d.Set("name", utils.PathSearch("listener.name", getListenerRespBody, nil)),
		d.Set("port_ranges", flattenGetListenerResponseBodyPortRange(getListenerRespBody)),
		d.Set("protocol", utils.PathSearch("listener.protocol", getListenerRespBody, nil)),
		d.Set("status", utils.PathSearch("listener.status", getListenerRespBody, nil)),
		d.Set("tags", flattenGetListenerResponseBodyResourceTag(getListenerRespBody)),
		d.Set("updated_at", utils.PathSearch("listener.updated_at", getListenerRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetListenerResponseBodyPortRange(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("listener.port_ranges", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"from_port": utils.PathSearch("from_port", v, nil),
			"to_port":   utils.PathSearch("to_port", v, nil),
		})
	}
	return rst
}

func flattenGetListenerResponseBodyResourceTag(resp interface{}) map[string]interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("listener.tags", resp, make([]interface{}, 0))
	return utils.FlattenTagsToMap(curJson)
}

func resourceListenerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	updateListenerhasChanges := []string{
		"client_affinity",
		"description",
		"name",
		"port_ranges",
	}

	if d.HasChanges(updateListenerhasChanges...) {
		// updateListener: Update the configuration of GA Listener
		var (
			updateListenerHttpUrl = "v1/listeners/{listener_id}"
			updateListenerProduct = "ga"
		)
		updateListenerClient, err := conf.NewServiceClient(updateListenerProduct, region)
		if err != nil {
			return diag.Errorf("error creating Listener Client: %s", err)
		}

		updateListenerPath := updateListenerClient.Endpoint + updateListenerHttpUrl
		updateListenerPath = strings.ReplaceAll(updateListenerPath, "{listener_id}", d.Id())

		updateListenerOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateListenerOpt.JSONBody = utils.RemoveNil(buildUpdateListenerBodyParams(d))
		_, err = updateListenerClient.Request("PUT", updateListenerPath, &updateListenerOpt)
		if err != nil {
			return diag.Errorf("error updating Listener: %s", err)
		}
		err = updateListenerWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the Update of Listener (%s) to complete: %s", d.Id(), err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		client, err := conf.NewServiceClient("ga", region)
		if err != nil {
			return diag.Errorf("error creating GA Client: %s", err)
		}

		oldRaw, newRaw := d.GetChange("tags")
		oldMap := oldRaw.(map[string]interface{})
		newMap := newRaw.(map[string]interface{})

		// remove old tags
		if len(oldMap) > 0 {
			if err = deleteTags(client, "ga-listeners", d.Id(), oldMap); err != nil {
				return diag.FromErr(err)
			}
		}

		// set new tags
		if len(newMap) > 0 {
			if err := createTags(client, "ga-listeners", d.Id(), newMap); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceListenerRead(ctx, d, meta)
}

func buildUpdateListenerBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"listener": map[string]interface{}{
			"client_affinity": utils.ValueIgnoreEmpty(d.Get("client_affinity")),
			"description":     d.Get("description"),
			"name":            utils.ValueIgnoreEmpty(d.Get("name")),
			"port_ranges":     buildUpdateListenerRequestBodyPortRange(d.Get("port_ranges")),
		},
	}
	return bodyParams
}

func buildUpdateListenerRequestBodyPortRange(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"from_port": utils.ValueIgnoreEmpty(raw["from_port"]),
				"to_port":   utils.ValueIgnoreEmpty(raw["to_port"]),
			}
		}
		return rst
	}
	return nil
}

func updateListenerWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// updateListenerWaiting: missing operation notes
			var (
				updateListenerWaitingHttpUrl = "v1/listeners/{listener_id}"
				updateListenerWaitingProduct = "ga"
			)
			updateListenerWaitingClient, err := config.NewServiceClient(updateListenerWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating Listener Client: %s", err)
			}

			updateListenerWaitingPath := updateListenerWaitingClient.Endpoint + updateListenerWaitingHttpUrl
			updateListenerWaitingPath = strings.ReplaceAll(updateListenerWaitingPath, "{listener_id}", d.Id())

			updateListenerWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			updateListenerWaitingResp, err := updateListenerWaitingClient.Request("GET", updateListenerWaitingPath, &updateListenerWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			updateListenerWaitingRespBody, err := utils.FlattenResponse(updateListenerWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`listener.status`, updateListenerWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `listener.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"ACTIVE",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return updateListenerWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return updateListenerWaitingRespBody, status, nil
			}

			return updateListenerWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceListenerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	// deleteListener: Delete an existing GA Listener
	var (
		deleteListenerHttpUrl = "v1/listeners/{listener_id}"
		deleteListenerProduct = "ga"
	)
	deleteListenerClient, err := conf.NewServiceClient(deleteListenerProduct, region)
	if err != nil {
		return diag.Errorf("error creating Listener Client: %s", err)
	}

	deleteListenerPath := deleteListenerClient.Endpoint + deleteListenerHttpUrl
	deleteListenerPath = strings.ReplaceAll(deleteListenerPath, "{listener_id}", d.Id())

	deleteListenerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteListenerClient.Request("DELETE", deleteListenerPath, &deleteListenerOpt)
	if err != nil {
		return diag.Errorf("error deleting Listener: %s", err)
	}

	err = deleteListenerWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the Delete of Listener (%s) to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteListenerWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// deleteListenerWaiting: missing operation notes
			var (
				deleteListenerWaitingHttpUrl = "v1/listeners/{listener_id}"
				deleteListenerWaitingProduct = "ga"
			)
			deleteListenerWaitingClient, err := config.NewServiceClient(deleteListenerWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating Listener Client: %s", err)
			}

			deleteListenerWaitingPath := deleteListenerWaitingClient.Endpoint + deleteListenerWaitingHttpUrl
			deleteListenerWaitingPath = strings.ReplaceAll(deleteListenerWaitingPath, "{listener_id}", d.Id())

			deleteListenerWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			deleteListenerWaitingResp, err := deleteListenerWaitingClient.Request("GET", deleteListenerWaitingPath, &deleteListenerWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
					return "Resource Not Found", "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			deleteListenerWaitingRespBody, err := utils.FlattenResponse(deleteListenerWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`listener.status`, deleteListenerWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `listener.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return deleteListenerWaitingRespBody, status, nil
			}

			return deleteListenerWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

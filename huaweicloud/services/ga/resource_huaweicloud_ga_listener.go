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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

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

		CustomizeDiff: config.MergeDefaultTags(),

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
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		httpUrl = "v1/listeners"
		product = "ga"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateListenerBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating GA listener: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	listenerId := utils.PathSearch("listener.id", respBody, "").(string)
	if listenerId == "" {
		return diag.Errorf("error creating GA listener: ID is not found in API response")
	}
	d.SetId(listenerId)

	err = createListenerWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the GA listener (%s) creation to complete: %s", d.Id(), err)
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
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		httpUrl          = "v1/listeners/{listener_id}"
		product          = "ga"
		targetStatus     = []string{"ACTIVE"}
		unexpectedStatus = []string{"ERROR"}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{listener_id}", d.Id())
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Request("GET", requestPath, &createOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			respBody, err := utils.FlattenResponse(resp)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch(`listener.status`, respBody, "").(string)
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

func resourceListenerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/listeners/{listener_id}"
		product = "ga"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{listener_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GA listener")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("accelerator_id", utils.PathSearch("listener.accelerator_id", respBody, nil)),
		d.Set("client_affinity", utils.PathSearch("listener.client_affinity", respBody, nil)),
		d.Set("created_at", utils.PathSearch("listener.created_at", respBody, nil)),
		d.Set("description", utils.PathSearch("listener.description", respBody, nil)),
		d.Set("name", utils.PathSearch("listener.name", respBody, nil)),
		d.Set("port_ranges", flattenGetListenerResponseBodyPortRange(respBody)),
		d.Set("protocol", utils.PathSearch("listener.protocol", respBody, nil)),
		d.Set("status", utils.PathSearch("listener.status", respBody, nil)),
		d.Set("tags", flattenGetListenerResponseBodyResourceTag(respBody)),
		d.Set("updated_at", utils.PathSearch("listener.updated_at", respBody, nil)),
		d.Set("frozen_info", flattenListenerFrozenInfo(utils.PathSearch("listener.frozen_info", respBody, nil))),
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

func flattenListenerFrozenInfo(resp interface{}) []map[string]interface{} {
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

func resourceListenerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		product = "ga"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	updateListenerHasChanges := []string{
		"client_affinity",
		"description",
		"name",
		"port_ranges",
	}

	if d.HasChanges(updateListenerHasChanges...) {
		requestPath := client.Endpoint + "v1/listeners/{listener_id}"
		requestPath = strings.ReplaceAll(requestPath, "{listener_id}", d.Id())
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateListenerBodyParams(d)),
		}

		_, err = client.Request("PUT", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating GA listener: %s", err)
		}

		err = updateListenerWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the GA listener (%s) update to complete: %s", d.Id(), err)
		}
	}

	// update tags
	if d.HasChange("tags") {
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
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		httpUrl          = "v1/listeners/{listener_id}"
		product          = "ga"
		targetStatus     = []string{"ACTIVE"}
		unexpectedStatus = []string{"ERROR"}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{listener_id}", d.Id())
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

			status := utils.PathSearch(`listener.status`, respBody, "").(string)
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

func resourceListenerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		httpUrl = "v1/listeners/{listener_id}"
		product = "ga"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{listener_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting GA listener: %s", err)
	}

	err = deleteListenerWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the GA listener (%s) delete to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteListenerWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		httpUrl          = "v1/listeners/{listener_id}"
		product          = "ga"
		unexpectedStatus = []string{"ERROR"}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{listener_id}", d.Id())
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

			status := utils.PathSearch(`listener.status`, respBody, "").(string)
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

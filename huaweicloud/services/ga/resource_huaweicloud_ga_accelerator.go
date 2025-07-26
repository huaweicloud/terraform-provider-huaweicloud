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

// @API GA POST /v1/accelerators
// @API GA GET /v1/accelerators/{accelerator_id}
// @API GA PUT /v1/accelerators/{accelerator_id}
// @API GA DELETE /v1/accelerators/{accelerator_id}
// @API GA POST /v1/{resource_type}/{resource_id}/tags/create
// @API GA DELETE /v1/{resource_type}/{resource_id}/tags/delete
func ResourceAccelerator() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAcceleratorCreate,
		UpdateContext: resourceAcceleratorUpdate,
		ReadContext:   resourceAcceleratorRead,
		DeleteContext: resourceAcceleratorDelete,
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the global accelerator name.`,
			},
			"ip_sets": {
				Type:     schema.TypeList,
				MaxItems: 2,
				Elem:     AcceleratorAccelerateIpSchema(),
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the description about the global accelerator.`,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
				ForceNew: true,
				Description: `|-
					Specifies the enterprise project ID of the tenant. The value is **0** or a string that
					contains a maximum of 36 characters in UUID format with hyphens (-). **0** indicates the
					default enterprise project.`,
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
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the specification ID.`,
			},
			"frozen_info": {
				Type:        schema.TypeList,
				Elem:        AcceleratorFrozenInfoSchema(),
				Computed:    true,
				Description: `Indicates the frozen details of cloud services or resources.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies when the global accelerator was created.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies when the global accelerator was updated.`,
			},
		},
	}
}

func AcceleratorAccelerateIpSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"area": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `|-
					Specifies the acceleration area. The value can be one of the following:
					  - **OUTOFCM**: Outside the Chinese mainland
					  - **CM**: Chinese mainland`,
			},
			"ip_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "IPV4",
				Description: `Specifies the IP address version.`,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the IP address.`,
			},
		},
	}
	return &sc
}

func AcceleratorFrozenInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"effect": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: `|-
					Specifies the status of the resource after being frozen. The value can be one of the following:
					  - **1** (default): The resource is frozen and can be released.
					  - **2**: The resource is frozen and cannot be released.
					  - **3**: The resource is frozen and cannot be renewed.`,
			},
			"scene": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Description: `|-
					Specifies the service scenario. The value can be one of the following:
					  - **ARREAR**: The cloud service is in arrears, including expiration of yearly/monthly resources
					  	and fee deduction failure of pay-per-use resources.
					  - **POLICE**: The cloud service is frozen for public security.
					  - **ILLEGAL**: The cloud service is frozen due to violation of laws and regulations.
					  - **VERIFY**: The cloud service is frozen because the user fails to pass the real-name authentication.
					  - **PARTNER**: A partner freezes their customer's resources.`,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: `|-
					Specifies the status of a cloud service or resource. The value can be one of the following:
					  - **0**: unfrozen/normal (The cloud service will recover after being unfrozen.)
					  - **1**: frozen (Resources and data will be retained, but the cloud service cannot be used.)
					  - **2**: deleted/terminated (Both resources and data will be cleared.)`,
			},
		},
	}
	return &sc
}

func resourceAcceleratorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		httpUrl = "v1/accelerators"
		product = "ga"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateAcceleratorBodyParams(d, conf)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating GA accelerator: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	acceleratorId := utils.PathSearch("accelerator.id", respBody, "").(string)
	if acceleratorId == "" {
		return diag.Errorf("error creating GA accelerator: unable to find the accelerator ID from the API response")
	}
	d.SetId(acceleratorId)

	err = createAcceleratorWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the accelerator (%s) creation to complete: %s", acceleratorId, err)
	}
	return resourceAcceleratorRead(ctx, d, meta)
}

func buildCreateAcceleratorBodyParams(d *schema.ResourceData, conf *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"accelerator": buildCreateAcceleratorAcceleratorChildBody(d, conf),
	}
	return bodyParams
}

func buildCreateAcceleratorAcceleratorChildBody(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	params := map[string]interface{}{
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"ip_sets":               buildCreateAcceleratorIpSetsChildBody(d),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"tags":                  utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
	}
	return params
}

func buildCreateAcceleratorIpSetsChildBody(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("ip_sets").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	ipSets := make([]map[string]interface{}, 0, len(rawParams))
	for _, v := range rawParams {
		raw := v.(map[string]interface{})
		params := map[string]interface{}{
			"area":       utils.ValueIgnoreEmpty(raw["area"]),
			"ip_address": utils.ValueIgnoreEmpty(raw["ip_address"]),
			"ip_type":    utils.ValueIgnoreEmpty(raw["ip_type"]),
		}
		ipSets = append(ipSets, params)
	}

	return ipSets
}

func createAcceleratorWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		httpUrl          = "v1/accelerators/{accelerator_id}"
		product          = "ga"
		targetStatus     = []string{"ACTIVE"}
		unexpectedStatus = []string{"ERROR"}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{accelerator_id}", d.Id())
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

			status := utils.PathSearch(`accelerator.status`, respBody, "").(string)
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

func resourceAcceleratorRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/accelerators/{accelerator_id}"
		product = "ga"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{accelerator_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GA accelerator")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("name", utils.PathSearch("accelerator.name", respBody, nil)),
		d.Set("description", utils.PathSearch("accelerator.description", respBody, nil)),
		d.Set("ip_sets", flattenAccelerateIpSets(utils.PathSearch("accelerator.ip_sets", respBody, make([]interface{}, 0)))),
		d.Set("enterprise_project_id", utils.PathSearch("accelerator.enterprise_project_id", respBody, nil)),
		d.Set("tags", flattenGetAcceleratorResponseBodyResourceTag(respBody)),
		d.Set("status", utils.PathSearch("accelerator.status", respBody, nil)),
		d.Set("flavor_id", utils.PathSearch("accelerator.flavor_id", respBody, nil)),
		d.Set("frozen_info", flattenGetAcceleratorResponseBodyFrozenInfo(respBody)),
		d.Set("created_at", utils.PathSearch("accelerator.created_at", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("accelerator.updated_at", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAccelerateIpSets(resp interface{}) []map[string]interface{} {
	rawArray, _ := resp.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		params := map[string]interface{}{
			"area":       utils.PathSearch("area", v, nil),
			"ip_address": utils.PathSearch("ip_address", v, nil),
			"ip_type":    utils.PathSearch("ip_type", v, nil),
		}
		rst[i] = params
	}

	return rst
}

func flattenGetAcceleratorResponseBodyResourceTag(resp interface{}) map[string]interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("accelerator.tags", resp, make([]interface{}, 0))
	return utils.FlattenTagsToMap(curJson)
}

func flattenGetAcceleratorResponseBodyFrozenInfo(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("accelerator.frozen_info", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"effect": utils.PathSearch("effect", curJson, nil),
			"scene":  utils.PathSearch("scene", curJson, nil),
			"status": utils.PathSearch("status", curJson, nil),
		},
	}
	return rst
}

func resourceAcceleratorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		product = "ga"
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	if d.HasChanges("name", "description") {
		requestPath := client.Endpoint + "v1/accelerators/{accelerator_id}"
		requestPath = strings.ReplaceAll(requestPath, "{accelerator_id}", d.Id())
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateAcceleratorBodyParams(d)),
		}

		_, err = client.Request("PUT", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating GA accelerator: %s", err)
		}

		err = updateAcceleratorWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the GA accelerator (%s) update to complete: %s", d.Id(), err)
		}
	}

	// Editing tags will cause the status to change to pending. Instances in the pending status do not support editing tags.
	if d.HasChange("tags") {
		oldRaw, newRaw := d.GetChange("tags")
		oldMap := oldRaw.(map[string]interface{})
		newMap := newRaw.(map[string]interface{})

		// remove old tags
		if len(oldMap) > 0 {
			if err = deleteTags(client, "ga-accelerators", d.Id(), oldMap); err != nil {
				return diag.FromErr(err)
			}

			err = updateAcceleratorWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.Errorf("error waiting for the GA accelerator (%s) delete tags to complete: %s", d.Id(), err)
			}
		}

		// set new tags
		if len(newMap) > 0 {
			if err := createTags(client, "ga-accelerators", d.Id(), newMap); err != nil {
				return diag.FromErr(err)
			}

			err = updateAcceleratorWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.Errorf("error waiting for the GA accelerator (%s) create tags to complete: %s", d.Id(), err)
			}
		}
	}

	return resourceAcceleratorRead(ctx, d, meta)
}

func createTags(createTagsClient *golangsdk.ServiceClient, resourceType, resourceId string, tags map[string]interface{}) error {
	createTagsHttpUrl := "v1/{resource_type}/{resource_id}/tags/create"
	createTagsPath := createTagsClient.Endpoint + createTagsHttpUrl
	createTagsPath = strings.ReplaceAll(createTagsPath, "{resource_type}", resourceType)
	createTagsPath = strings.ReplaceAll(createTagsPath, "{resource_id}", resourceId)
	createTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	createTagsOpt.JSONBody = map[string]interface{}{
		"tags": utils.ExpandResourceTags(tags),
	}

	_, err := createTagsClient.Request("POST", createTagsPath, &createTagsOpt)
	if err != nil {
		return fmt.Errorf("error creating GA (%s) tags: %s", resourceType, err)
	}
	return nil
}

func deleteTags(deleteTagsClient *golangsdk.ServiceClient, resourceType, resourceId string, tags map[string]interface{}) error {
	deleteTagsHttpUrl := "v1/{resource_type}/{resource_id}/tags/delete"
	deleteTagsPath := deleteTagsClient.Endpoint + deleteTagsHttpUrl
	deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{resource_type}", resourceType)
	deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{resource_id}", resourceId)
	deleteTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	deleteTagsOpt.JSONBody = map[string]interface{}{
		"tags": utils.ExpandResourceTags(tags),
	}

	_, err := deleteTagsClient.Request("DELETE", deleteTagsPath, &deleteTagsOpt)
	if err != nil {
		return fmt.Errorf("error deleting GA (%s) tags: %s", resourceType, err)
	}
	return nil
}

func buildUpdateAcceleratorBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"accelerator": buildUpdateAcceleratorAcceleratorChildBody(d),
	}
	return bodyParams
}

func buildUpdateAcceleratorAcceleratorChildBody(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return params
}

func updateAcceleratorWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		httpUrl          = "v1/accelerators/{accelerator_id}"
		product          = "ga"
		targetStatus     = []string{"ACTIVE"}
		unexpectedStatus = []string{"ERROR"}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{accelerator_id}", d.Id())
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

			status := utils.PathSearch(`accelerator.status`, respBody, "").(string)
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

func resourceAcceleratorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		httpUrl = "v1/accelerators/{accelerator_id}"
		product = "ga"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{accelerator_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting GA accelerator: %s", err)
	}

	err = deleteAcceleratorWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the GA accelerator (%s) delete to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteAcceleratorWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		httpUrl          = "v1/accelerators/{accelerator_id}"
		product          = "ga"
		unexpectedStatus = []string{"ERROR"}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{accelerator_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Request("GET", requestPath, &deleteOpt)
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

			status := utils.PathSearch(`accelerator.status`, respBody, "").(string)
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

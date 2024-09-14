// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product GA
// ---------------------------------------------------------------

package ga

import (
	"context"
	"fmt"
	"log"
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

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the global accelerator name.`,
			},
			"ip_sets": {
				Type:     schema.TypeList,
				MaxItems: 1,
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
				ValidateFunc: validation.StringInSlice([]string{
					"OUTOFCM", "CM",
				}, false),
			},
			"ip_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "IPV4",
				Description: `Specifies the IP address version.`,
				ValidateFunc: validation.StringInSlice([]string{
					"IPV4",
				}, false),
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
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	// createAccelerator: Create a GA Accelerator.
	var (
		createAcceleratorHttpUrl = "v1/accelerators"
		createAcceleratorProduct = "ga"
	)
	createAcceleratorClient, err := conf.NewServiceClient(createAcceleratorProduct, region)
	if err != nil {
		return diag.Errorf("error creating Accelerator Client: %s", err)
	}

	createAcceleratorPath := createAcceleratorClient.Endpoint + createAcceleratorHttpUrl

	createAcceleratorOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createAcceleratorOpt.JSONBody = utils.RemoveNil(buildCreateAcceleratorBodyParams(d, conf))
	createAcceleratorResp, err := createAcceleratorClient.Request("POST", createAcceleratorPath, &createAcceleratorOpt)
	if err != nil {
		return diag.Errorf("error creating Accelerator: %s", err)
	}

	createAcceleratorRespBody, err := utils.FlattenResponse(createAcceleratorResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("accelerator.id", createAcceleratorRespBody)
	if err != nil {
		return diag.Errorf("error creating Accelerator: ID is not found in API response")
	}
	d.SetId(id.(string))

	err = createAcceleratorWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Create of Accelerator (%s) to complete: %s", d.Id(), err)
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

	raw := rawParams[0].(map[string]interface{})
	params := map[string]interface{}{
		"area":       utils.ValueIgnoreEmpty(raw["area"]),
		"ip_address": utils.ValueIgnoreEmpty(raw["ip_address"]),
		"ip_type":    utils.ValueIgnoreEmpty(raw["ip_type"]),
	}

	return []map[string]interface{}{params}
}

func createAcceleratorWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// createAcceleratorWaiting: missing operation notes
			var (
				createAcceleratorWaitingHttpUrl = "v1/accelerators/{accelerator_id}"
				createAcceleratorWaitingProduct = "ga"
			)
			createAcceleratorWaitingClient, err := config.NewServiceClient(createAcceleratorWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating Accelerator Client: %s", err)
			}

			createAcceleratorWaitingPath := createAcceleratorWaitingClient.Endpoint + createAcceleratorWaitingHttpUrl
			createAcceleratorWaitingPath = strings.ReplaceAll(createAcceleratorWaitingPath, "{accelerator_id}", d.Id())

			createAcceleratorWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			createAcceleratorWaitingResp, err := createAcceleratorWaitingClient.Request("GET",
				createAcceleratorWaitingPath, &createAcceleratorWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createAcceleratorWaitingRespBody, err := utils.FlattenResponse(createAcceleratorWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`accelerator.status`, createAcceleratorWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `accelerator.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"ACTIVE",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createAcceleratorWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return createAcceleratorWaitingRespBody, status, nil
			}

			return createAcceleratorWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceAcceleratorRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	// getAccelerator: Query the GA accelerator detail
	var (
		getAcceleratorHttpUrl = "v1/accelerators/{accelerator_id}"
		getAcceleratorProduct = "ga"
	)
	getAcceleratorClient, err := conf.NewServiceClient(getAcceleratorProduct, region)
	if err != nil {
		return diag.Errorf("error creating Accelerator Client: %s", err)
	}

	getAcceleratorPath := getAcceleratorClient.Endpoint + getAcceleratorHttpUrl
	getAcceleratorPath = strings.ReplaceAll(getAcceleratorPath, "{accelerator_id}", d.Id())

	getAcceleratorOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getAcceleratorResp, err := getAcceleratorClient.Request("GET", getAcceleratorPath, &getAcceleratorOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Accelerator")
	}

	getAcceleratorRespBody, err := utils.FlattenResponse(getAcceleratorResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("name", utils.PathSearch("accelerator.name", getAcceleratorRespBody, nil)),
		d.Set("description", utils.PathSearch("accelerator.description", getAcceleratorRespBody, nil)),
		d.Set("ip_sets", flattenGetAcceleratorResponseBodyAccelerateIp(getAcceleratorRespBody)),
		d.Set("enterprise_project_id", utils.PathSearch("accelerator.enterprise_project_id", getAcceleratorRespBody, nil)),
		d.Set("tags", flattenGetAcceleratorResponseBodyResourceTag(getAcceleratorRespBody)),
		d.Set("status", utils.PathSearch("accelerator.status", getAcceleratorRespBody, nil)),
		d.Set("flavor_id", utils.PathSearch("accelerator.flavor_id", getAcceleratorRespBody, nil)),
		d.Set("frozen_info", flattenGetAcceleratorResponseBodyFrozenInfo(getAcceleratorRespBody)),
		d.Set("created_at", utils.PathSearch("accelerator.created_at", getAcceleratorRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("accelerator.updated_at", getAcceleratorRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetAcceleratorResponseBodyAccelerateIp(resp interface{}) []interface{} {
	var rst []interface{}
	curJson, err := jmespath.Search("accelerator.ip_sets", resp)
	if err != nil {
		log.Printf("[ERROR] error parsing accelerator.ip_sets from response= %#v", resp)
		return rst
	}

	curArray := curJson.([]interface{})
	rst = []interface{}{
		map[string]interface{}{
			"area":       utils.PathSearch("area", curArray[0], nil),
			"ip_address": utils.PathSearch("ip_address", curArray[0], nil),
			"ip_type":    utils.PathSearch("ip_type", curArray[0], nil),
		},
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
	curJson, err := jmespath.Search("accelerator.frozen_info", resp)
	if err != nil {
		log.Printf("[ERROR] error parsing accelerator.frozen_info from response= %#v", resp)
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
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	updateAcceleratorhasChanges := []string{
		"name",
		"description",
	}

	if d.HasChanges(updateAcceleratorhasChanges...) {
		// updateAccelerator: Update the configuration of GA accelerator
		var (
			updateAcceleratorHttpUrl = "v1/accelerators/{accelerator_id}"
			updateAcceleratorProduct = "ga"
		)
		updateAcceleratorClient, err := conf.NewServiceClient(updateAcceleratorProduct, region)
		if err != nil {
			return diag.Errorf("error creating Accelerator Client: %s", err)
		}

		updateAcceleratorPath := updateAcceleratorClient.Endpoint + updateAcceleratorHttpUrl
		updateAcceleratorPath = strings.ReplaceAll(updateAcceleratorPath, "{accelerator_id}", d.Id())

		updateAcceleratorOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateAcceleratorOpt.JSONBody = utils.RemoveNil(buildUpdateAcceleratorBodyParams(d))
		_, err = updateAcceleratorClient.Request("PUT", updateAcceleratorPath, &updateAcceleratorOpt)
		if err != nil {
			return diag.Errorf("error updating Accelerator: %s", err)
		}
		err = updateAcceleratorWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the Update of Accelerator (%s) to complete: %s", d.Id(), err)
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
			if err = deleteTags(client, "ga-accelerators", d.Id(), oldMap); err != nil {
				return diag.FromErr(err)
			}
		}

		// set new tags
		if len(newMap) > 0 {
			if err := createTags(client, "ga-accelerators", d.Id(), newMap); err != nil {
				return diag.FromErr(err)
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
		return fmt.Errorf("error creating tags: %s", err)
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
		return fmt.Errorf("error deleting tags: %s", err)
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
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// updateAcceleratorWaiting: missing operation notes
			var (
				updateAcceleratorWaitingHttpUrl = "v1/accelerators/{accelerator_id}"
				updateAcceleratorWaitingProduct = "ga"
			)
			updateAcceleratorWaitingClient, err := config.NewServiceClient(updateAcceleratorWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating Accelerator Client: %s", err)
			}

			updateAcceleratorWaitingPath := updateAcceleratorWaitingClient.Endpoint + updateAcceleratorWaitingHttpUrl
			updateAcceleratorWaitingPath = strings.ReplaceAll(updateAcceleratorWaitingPath, "{accelerator_id}", d.Id())

			updateAcceleratorWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			updateAcceleratorWaitingResp, err := updateAcceleratorWaitingClient.Request("GET",
				updateAcceleratorWaitingPath, &updateAcceleratorWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			updateAcceleratorWaitingRespBody, err := utils.FlattenResponse(updateAcceleratorWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`accelerator.status`, updateAcceleratorWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `accelerator.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"ACTIVE",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return updateAcceleratorWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return updateAcceleratorWaitingRespBody, status, nil
			}

			return updateAcceleratorWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceAcceleratorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	// deleteAccelerator: Delete an existing GA Accelerator
	var (
		deleteAcceleratorHttpUrl = "v1/accelerators/{accelerator_id}"
		deleteAcceleratorProduct = "ga"
	)
	deleteAcceleratorClient, err := conf.NewServiceClient(deleteAcceleratorProduct, region)
	if err != nil {
		return diag.Errorf("error creating Accelerator Client: %s", err)
	}

	deleteAcceleratorPath := deleteAcceleratorClient.Endpoint + deleteAcceleratorHttpUrl
	deleteAcceleratorPath = strings.ReplaceAll(deleteAcceleratorPath, "{accelerator_id}", d.Id())

	deleteAcceleratorOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteAcceleratorClient.Request("DELETE", deleteAcceleratorPath, &deleteAcceleratorOpt)
	if err != nil {
		return diag.Errorf("error deleting Accelerator: %s", err)
	}

	err = deleteAcceleratorWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the Delete of Accelerator (%s) to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteAcceleratorWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// deleteAcceleratorWaiting: missing operation notes
			var (
				deleteAcceleratorWaitingHttpUrl = "v1/accelerators/{accelerator_id}"
				deleteAcceleratorWaitingProduct = "ga"
			)
			deleteAcceleratorWaitingClient, err := config.NewServiceClient(deleteAcceleratorWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating Accelerator Client: %s", err)
			}

			deleteAcceleratorWaitingPath := deleteAcceleratorWaitingClient.Endpoint + deleteAcceleratorWaitingHttpUrl
			deleteAcceleratorWaitingPath = strings.ReplaceAll(deleteAcceleratorWaitingPath, "{accelerator_id}", d.Id())

			deleteAcceleratorWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			deleteAcceleratorWaitingResp, err := deleteAcceleratorWaitingClient.Request("GET",
				deleteAcceleratorWaitingPath, &deleteAcceleratorWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
					return "Resource Not Found", "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			deleteAcceleratorWaitingRespBody, err := utils.FlattenResponse(deleteAcceleratorWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`accelerator.status`, deleteAcceleratorWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `accelerator.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return deleteAcceleratorWaitingRespBody, status, nil
			}

			return deleteAcceleratorWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

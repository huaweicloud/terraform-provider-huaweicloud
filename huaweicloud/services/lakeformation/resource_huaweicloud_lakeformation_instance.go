package lakeformation

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	instanceParamKeys          = []string{"specs"}
	instanceNonUpdatableParams = []string{"shared", "enterprise_project_id"}
)

// @API LakeFormation POST /v1/{project_id}/instances
// @API LakeFormation GET /v1/{project_id}/instances/{instance_id}
// @API LakeFormation PUT /v1/{project_id}/instances/{instance_id}
// @API LakeFormation POST /v1/{project_id}/instances/{instance_id}/scale
// @API LakeFormation PUT /v1/{project_id}/instances/{instance_id}/tags
// @API LakeFormation DELETE /v1/{project_id}/instances/{instance_id}
func ResourceInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		UpdateContext: resourceInstanceUpdate,
		DeleteContext: resourceInstanceDelete,

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(instanceNonUpdatableParams),
			config.MergeDefaultTags(),
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Hour),
			Update: schema.DefaultTimeout(2 * time.Hour),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the instance is located.`,
			},

			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the instance.`,
			},
			"shared": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Whether the instance is shared.`,
			},

			// Optional parameters.
			"specs": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spec_code": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Description: utils.SchemaDesc(`The specification code.`,
								utils.SchemaDescInput{
									Required: true,
								},
							),
						},
						// The `Required` property of numeric types may cause unexpected changes when additional
						// objects are returned from the remote source, thus deviating from the Diff logic.
						"stride_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `The stride number of the specification.`,
						},
						"product_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Description: utils.SchemaDesc(`The product ID of the specification.`,
								utils.SchemaDescInput{
									Computed: true,
								},
							),
						},
					},
				},
				DiffSuppressFunc: utils.SuppressObjectSliceDiffs(),
				Description:      `The list of specifications.`,
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the instance.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project.`,
			},
			"tags": common.TagsSchema(`The key/value pairs to associate with the instance.`),
			"to_recycle_bin": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Whether to put the instance into the recycle bin when delete postpaid instance.`,
			},

			// Attributes.
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the instance.`,
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the instance is the default instance.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the instance, in RFC3339 format.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the instance, in RFC3339 format.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},

			// Internal attributes.
			"specs_origin": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spec_code": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The specification code.`,
						},
						"stride_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `The stride number.`,
						},
						"product_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The product ID.`,
						},
					},
				},
				DiffSuppressFunc: utils.SuppressDiffAll,
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
the new value next time the change is made. The corresponding parameter name is 'specs'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildInstanceSpecs(specs []interface{}) []interface{} {
	if len(specs) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(specs))
	for _, spec := range specs {
		result = append(result, map[string]interface{}{
			// Required parameters.
			"spec_code":  utils.PathSearch("spec_code", spec, nil),
			"stride_num": utils.PathSearch("stride_num", spec, nil),
			// Optional parameters.
			"product_id": utils.PathSearch("product_id", spec, nil),
		})
	}

	return result
}

func buildCreateInstanceBodyParams(cfg *config.Config, d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		// Required parameters.
		"name": d.Get("name"),
		// The current POST API does not return any information (instance ID or order ID) when creating a prePaid
		// type instance, so, the instance status cannot be tracked, that's why it is forced to be set to postPaid.
		"charge_mode": "postPaid",
		"shared":      d.Get("shared"),
		// Optional parameters.
		"specs":                 buildInstanceSpecs(d.Get("specs").([]interface{})),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"tags":                  utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
	}

	return bodyParams
}

func GetInstanceById(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/instances/{instance_id}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	// Instances in the recycle bin will return a 'DELETING' status and the value of attribute 'in_recycle_bin' will be true.
	if utils.PathSearch("status", respBody, "").(string) == "DELETING" && utils.PathSearch("in_recycle_bin", respBody, false).(bool) {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/instances/{instance_id}",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the instance (%s) has been moved to the recycle bin", instanceId)),
			},
		}
	}
	return respBody, nil
}

func instanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetInstanceById(client, instanceId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return "not_found", "COMPLETED", nil
			}
			return respBody, "ERROR", err
		}

		statusResp := utils.PathSearch("status", respBody, "").(string)
		if utils.StrSliceContains([]string{"RESOURCE_PREPARATION_FAIL", "SCALE_FAIL"},
			statusResp) {
			return respBody, "ERROR", fmt.Errorf("unexpected status (%s)", statusResp)
		}

		if utils.StrSliceContains(targets, statusResp) {
			return respBody, "COMPLETED", nil
		}
		return "continue", "PENDING", nil
	}
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("lakeformation", region)
	if err != nil {
		return diag.Errorf("error creating LakeFormation client: %s", err)
	}

	httpUrl := "v1/{project_id}/instances"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildCreateInstanceBodyParams(cfg, d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating instance: %s", err)
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceId := utils.PathSearch("instance_id", respBody, "").(string)
	if instanceId == "" {
		return diag.Errorf("unable to find the instance ID from the API response")
	}
	d.SetId(instanceId)

	err = utils.RefreshSliceParamOriginValues(d, instanceParamKeys)
	if err != nil {
		// Don't report errors if origin values are failed to be refreshed.
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      instanceStateRefreshFunc(client, instanceId, []string{"RUNNING"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        45 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the status of instance (%s) to become running: %s", instanceId, err)
	}

	return resourceInstanceRead(ctx, d, meta)
}

func orderSpecsBySpecsOrderOrigin(specs, specsOrigin []interface{}) []interface{} {
	if len(specsOrigin) == 0 {
		return specs
	}

	sortedSpecs := make([]interface{}, 0, len(specs))
	specsCopy := specs
	for copyIndex, spec := range specsCopy {
		specCode := utils.PathSearch("spec_code", spec, "").(string)
		for originIndex, specOrigin := range specsOrigin {
			if utils.PathSearch("spec_code", specOrigin, "").(string) != specCode {
				continue
			}
			sortedSpecs = append(sortedSpecs, specsCopy[copyIndex])
			// Retain additional specification configurations from the remote service and these will not be sorted.
			specsCopy = append(specsCopy[:copyIndex], specsCopy[copyIndex+1:]...)
			specsOrigin = append(specsOrigin[:originIndex], specsOrigin[originIndex+1:]...)
		}
	}
	sortedSpecs = append(sortedSpecs, specsCopy...)
	sortedSpecs = append(sortedSpecs, specsOrigin...)
	return sortedSpecs
}

func flattenInstanceSpecs(specs, specsOrderOrigin []interface{}) []map[string]interface{} {
	if len(specs) < 1 {
		return nil
	}

	sortedSpecs := orderSpecsBySpecsOrderOrigin(specs, specsOrderOrigin)
	result := make([]map[string]interface{}, 0, len(sortedSpecs))
	for _, spec := range sortedSpecs {
		result = append(result, map[string]interface{}{
			"spec_code":  utils.PathSearch("spec_code", spec, nil),
			"stride_num": utils.PathSearch("stride_num", spec, nil),
			"product_id": utils.PathSearch("product_id", spec, nil),
		})
	}

	return result
}

func resourceInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("lakeformation", region)
	if err != nil {
		return diag.Errorf("error creating LakeFormation client: %s", err)
	}

	respBody, err := GetInstanceById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving LakeFormation instance")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		// Required parameters.
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("shared", utils.PathSearch("shared", respBody, nil)),
		// Optional parameters.
		d.Set("specs", flattenInstanceSpecs(utils.PathSearch("specs", respBody, make([]interface{}, 0)).([]interface{}),
			d.Get("specs_origin").([]interface{}))),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", respBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", respBody,
			make([]interface{}, 0)).([]interface{}))),
		// Attributes.
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("is_default", utils.PathSearch("default_instance", respBody, nil)),
		d.Set("create_time", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_time",
			respBody, "").(string))/1000, false)),
		d.Set("update_time", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("update_time",
			respBody, "").(string))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
	}
}

func updateInstance(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/instances/{instance_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildUpdateInstanceBodyParams(d),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func updateInstanceSpecs(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl    = "v1/{project_id}/instances/{instance_id}/scale"
		instanceId = d.Id()
		rawConfig  = d.GetRawConfig()
	)

	scalePath := client.Endpoint + httpUrl
	scalePath = strings.ReplaceAll(scalePath, "{project_id}", client.ProjectID)
	scalePath = strings.ReplaceAll(scalePath, "{instance_id}", d.Id())

	scaleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"specs": buildInstanceSpecs(utils.GetNestedObjectFromRawConfig(rawConfig, "specs").([]interface{})),
		},
	}

	_, err := client.Request("POST", scalePath, &scaleOpt)
	if err != nil {
		return fmt.Errorf("error updating instance (%s) specs: %s", instanceId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      instanceStateRefreshFunc(client, d.Id(), []string{"RUNNING"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        45 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the status of instance (%s) scale to complete: %s", d.Id(), err)
	}
	return nil
}

func updateInstanceTags(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/instances/{instance_id}/tags"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"tags": utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
		},
	}
	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("lakeformation", region)
	if err != nil {
		return diag.Errorf("error creating client: %s", err)
	}

	if d.HasChanges("name", "description") {
		err = updateInstance(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("specs") {
		err = updateInstanceSpecs(ctx, client, d)
		if err != nil {
			return diag.FromErr(err)
		}
		err = utils.RefreshSliceParamOriginValues(d, instanceParamKeys)
		if err != nil {
			// Don't report errors if origin values are failed to be refreshed.
			log.Printf("[WARN] Unable to refresh the origin values: %s", err)
		}
	}

	if d.HasChange("tags") {
		err = updateInstanceTags(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceInstanceRead(ctx, d, meta)
}

func deleteInstance(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/instances/{instance_id}?to_recycle_bin={to_recycle_bin}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Id())
	deletePath = strings.ReplaceAll(deletePath, "{to_recycle_bin}", strconv.FormatBool(d.Get("to_recycle_bin").(bool)))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Id()
	)

	client, err := cfg.NewServiceClient("lakeformation", region)
	if err != nil {
		return diag.Errorf("error creating client: %s", err)
	}

	err = deleteInstance(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting instance")
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: instanceStateRefreshFunc(client, instanceId, nil),
		Timeout: d.Timeout(schema.TimeoutDelete),
		// Some instances will not enter the recycle bin immediately after deletion, and there will be a period of time
		// during which data is lost, which may be a design issue on the service side.
		Delay:        10 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the deletion of instance (%s) to complete: %s", instanceId, err)
	}
	return nil
}

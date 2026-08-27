package cceautopilot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var autopilotAddonNonUpdatableParams = []string{
	"cluster_id", "addon_template_name",
}

// @API CCE POST /autopilot/v3/addons
// @API CCE GET /autopilot/v3/addons/{id}
// @API CCE PUT /autopilot/v3/addons/{id}
// @API CCE DELETE /autopilot/v3/addons/{id}
func ResourceAutopilotAddon() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAutopilotAddonCreate,
		ReadContext:   resourceAutopilotAddonRead,
		UpdateContext: resourceAutopilotAddonUpdate,
		DeleteContext: resourceAutopilotAddonDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(autopilotAddonNonUpdatableParams),

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
			},
			"addon_template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeMap,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateAddonBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	spec, err := buildAddonSpecBodyParams(d)
	if err != nil {
		return nil, err
	}
	bodyParams := map[string]interface{}{
		"kind":       "Addon",
		"apiVersion": "v3",
		"metadata":   buildCreateAddonMetadataBodyParams(d),
		"spec":       spec,
	}

	return bodyParams, nil
}

func buildCreateAddonMetadataBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":  utils.ValueIgnoreEmpty(d.Get("name")),
		"alias": utils.ValueIgnoreEmpty(d.Get("alias")),
		"annotations": map[string]interface{}{
			"addon.install/type": "install",
		},
	}

	return bodyParams
}

func buildAddonSpecBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	values, err := buildAddonValuesBodyParams(d)
	if err != nil {
		return nil, err
	}
	bodyParams := map[string]interface{}{
		"clusterID":         d.Get("cluster_id"),
		"addonTemplateName": d.Get("addon_template_name"),
		"version":           utils.ValueIgnoreEmpty(d.Get("version")),
		"values":            values,
	}

	return bodyParams, nil
}

func buildAddonValuesBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	valuesRaw := d.Get("values").(map[string]interface{})
	bodyParams := make(map[string]interface{}, len(valuesRaw))
	for k, v := range valuesRaw {
		var value map[string]interface{}
		err := json.Unmarshal([]byte(v.(string)), &value)
		if err != nil {
			return nil, err
		}
		bodyParams[k] = value
	}

	return bodyParams, nil
}

func resourceAutopilotAddonCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cce", region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	createPath := client.Endpoint + "autopilot/v3/addons"
	createParams, err := buildCreateAddonBodyParams(d)
	if err != nil {
		return diag.Errorf("error building create options of CCE autopilot addon: %s", err)
	}
	createAddonOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(createParams),
	}

	createAddonResp, err := client.Request("POST", createPath, &createAddonOpt)
	if err != nil {
		return diag.Errorf("error creating CCE autopilot add-on: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createAddonResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("metadata.uid", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CCE autopilot add-on: ID is not found in API response")
	}
	d.SetId(id)

	err = waitingForAddonJobCompleted(ctx, client, d, d.Timeout(schema.TimeoutCreate), []string{"running", "available"})
	if err != nil {
		return diag.Errorf("error waiting for creating CCE autopilot add-on (%s) to complete: %s", id, err)
	}

	return resourceAutopilotAddonRead(ctx, d, meta)
}

func resourceAutopilotAddonRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("cce", region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	getAddonRespBody, err := getAutopilotAddon(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCE autopilot add-on")
	}

	// values not set, because the response if different from the user input
	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("cluster_id", utils.PathSearch("spec.clusterID", getAddonRespBody, nil)),
		d.Set("addon_template_name", utils.PathSearch("spec.addonTemplateName", getAddonRespBody, nil)),
		d.Set("version", utils.PathSearch("spec.version", getAddonRespBody, nil)),
		d.Set("name", utils.PathSearch("metadata.name", getAddonRespBody, nil)),
		d.Set("alias", utils.PathSearch("metadata.alias", getAddonRespBody, nil)),
		d.Set("created_at", utils.PathSearch("metadata.creationTimestamp", getAddonRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("metadata.updateTimestamp", getAddonRespBody, nil)),
		d.Set("status", utils.PathSearch("status.status", getAddonRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getAutopilotAddon(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	getPath := client.Endpoint + "autopilot/v3/addons/{id}"
	getPath = strings.ReplaceAll(getPath, "{id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getAddonResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getAddonResp)
}

func buildUpdateAddonBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	spec, err := buildAddonSpecBodyParams(d)
	if err != nil {
		return nil, err
	}

	bodyParams := map[string]interface{}{
		"kind":       "Addon",
		"apiVersion": "v3",
		"metadata":   buildUpdateAddonMetadataBodyParams(d),
		"spec":       spec,
	}

	return bodyParams, nil
}

func buildUpdateAddonMetadataBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":  utils.ValueIgnoreEmpty(d.Get("name")),
		"alias": utils.ValueIgnoreEmpty(d.Get("alias")),
		"annotations": map[string]interface{}{
			"addon.upgrade/type": "upgrade",
		},
	}

	return bodyParams
}

func resourceAutopilotAddonUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateAddonProduct = "cce"
		updateAddonHttpUrl = "autopilot/v3/addons/{id}"
	)

	client, err := cfg.NewServiceClient(updateAddonProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	if d.HasChangeExcept("enable_force_new") {
		updatePath := client.Endpoint + updateAddonHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{id}", d.Id())

		updateOpts, err := buildUpdateAddonBodyParams(d)
		if err != nil {
			return diag.Errorf("error building update options of CCE autopilot add-on: %s", err)
		}
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			JSONBody: utils.RemoveNil(updateOpts),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating CCE autopilot add-on: %s", err)
		}

		err = waitingForAddonJobCompleted(ctx, client, d, d.Timeout(schema.TimeoutUpdate), []string{"running", "available"})
		if err != nil {
			return diag.Errorf("error waiting for updating CCE autopilot add-on (%s) to complete: %s", d.Id(), err)
		}
	}

	return resourceAutopilotAddonRead(ctx, d, meta)
}

func resourceAutopilotAddonDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cce", region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	deletePath := client.Endpoint + "autopilot/v3/addons/{id}"
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CCE autopilot add-on")
	}

	err = waitingForAddonJobCompleted(ctx, client, d, d.Timeout(schema.TimeoutDelete), nil)
	if err != nil {
		return diag.Errorf("error waiting for deleting CCE autopilot add-on (%s) to complete: %s", d.Id(), err)
	}

	return nil
}

func waitingForAddonJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	t time.Duration, targets []string) error {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshAddonStatus(client, d, targets),
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshAddonStatus(client *golangsdk.ServiceClient, d *schema.ResourceData, targets []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Expect the status of CCE autopilot add-on to be any one of the status list: %v.", targets)
		addonRespBody, err := getAutopilotAddon(client, d)
		if err != nil {
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				// On delete (targets is nil), the resource is already gone, which is the expected
				// outcome. Return a non-nil empty object so the SDK state machine matches "COMPLETED"
				// against Target=["COMPLETED"] and finishes normally.
				// On create/update (targets is non-nil), a 404 means the resource does not exist,
				// return nil so the SDK increments notfoundTick and eventually reports NotFoundError.
				if len(targets) == 0 {
					return map[string]interface{}{}, "COMPLETED", nil
				}
				return nil, "COMPLETED", nil
			}
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status.status", addonRespBody, nil)
		if status == nil {
			return nil, "ERROR", fmt.Errorf("error parsing status from response body")
		}

		statusStr := status.(string)
		invalidStatuses := []string{"installFailed", "upgradeFailed", "deleteFailed", "rollbackFailed", "abnormal", "unknown"}
		if utils.IsStrContainsSliceElement(statusStr, invalidStatuses, true, true) {
			reason := utils.PathSearch("status.reason", addonRespBody, "")
			message := utils.PathSearch("status.message", addonRespBody, "")
			return nil, "ERROR", fmt.Errorf("addon status is %s, reason: %s, message: %s",
				statusStr, reason, message)
		}

		if utils.StrSliceContains(targets, statusStr) {
			return addonRespBody, "COMPLETED", nil
		}
		return addonRespBody, "PENDING", nil
	}
}

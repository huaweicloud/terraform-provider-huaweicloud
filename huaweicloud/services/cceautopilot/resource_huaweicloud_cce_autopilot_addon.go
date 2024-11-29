package cceautopilot

import (
	"context"
	"encoding/json"
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

	var (
		createAddonHttpUrl = "autopilot/v3/addons"
		createAddonProduct = "cce"
	)
	createAddonClient, err := cfg.NewServiceClient(createAddonProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	createAddonPath := createAddonClient.Endpoint + createAddonHttpUrl

	createAddonOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createOpts, err := buildCreateAddonBodyParams(d)
	if err != nil {
		return diag.Errorf("error building create options of CCE autopolit addon: %s", err)
	}
	createAddonOpt.JSONBody = utils.RemoveNil(createOpts)
	createAddonResp, err := createAddonClient.Request("POST", createAddonPath, &createAddonOpt)
	if err != nil {
		return diag.Errorf("error creating CCE autopolit add-on: %s", err)
	}

	createAddonRespBody, err := utils.FlattenResponse(createAddonResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("metadata.uid", createAddonRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CCE autopilot add-on: ID is not found in API response")
	}
	d.SetId(id)

	err = addonWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for creating CCE autopilot add-on (%s) to complete: %s", id, err)
	}

	return resourceAutopilotAddonRead(ctx, d, meta)
}

func resourceAutopilotAddonRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getAddonHttpUrl = "autopilot/v3/addons/{id}"
		getAddonProduct = "cce"
	)
	getAddonClient, err := cfg.NewServiceClient(getAddonProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	getAddonPath := getAddonClient.Endpoint + getAddonHttpUrl
	getAddonPath = strings.ReplaceAll(getAddonPath, "{id}", d.Id())

	getAddonOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getAddonResp, err := getAddonClient.Request("GET", getAddonPath, &getAddonOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCE autopolit add-on")
	}

	getAddonRespBody, err := utils.FlattenResponse(getAddonResp)
	if err != nil {
		return diag.FromErr(err)
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

	updateAddonClient, err := cfg.NewServiceClient(updateAddonProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	updateAddonPath := updateAddonClient.Endpoint + updateAddonHttpUrl
	updateAddonPath = strings.ReplaceAll(updateAddonPath, "{id}", d.Id())

	updateOpts, err := buildUpdateAddonBodyParams(d)
	if err != nil {
		return diag.Errorf("error building update options of CCE autopolit add-on: %s", err)
	}
	updateAddonOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(updateOpts),
	}

	_, err = updateAddonClient.Request("PUT", updateAddonPath, &updateAddonOpt)
	if err != nil {
		return diag.Errorf("error updating CCE autopolit add-on: %s", err)
	}

	return resourceAutopilotAddonRead(ctx, d, meta)
}

func resourceAutopilotAddonDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteAddonHttpUrl = "autopilot/v3/addons/{id}"
		deleteAddonProduct = "cce"
	)
	deleteAddonClient, err := cfg.NewServiceClient(deleteAddonProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	deleteAddonPath := deleteAddonClient.Endpoint + deleteAddonHttpUrl
	deleteAddonPath = strings.ReplaceAll(deleteAddonPath, "{project_id}", deleteAddonClient.ProjectID)
	deleteAddonPath = strings.ReplaceAll(deleteAddonPath, "{id}", d.Id())

	deleteAddonOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = deleteAddonClient.Request("DELETE", deleteAddonPath, &deleteAddonOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CCE autopolit add-on")
	}

	err = addonWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for deleting CCE autopilot add-on (%s) to complete: %s", d.Id(), err)
	}

	return nil
}

func addonWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			var (
				adonWaitingHttpUrl = "autopilot/v3/addons/{id}"
				adonWaitingProduct = "cce"
			)
			adonWaitingClient, err := cfg.NewServiceClient(adonWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating CCE client: %s", err)
			}

			adonWaitingPath := adonWaitingClient.Endpoint + adonWaitingHttpUrl
			adonWaitingPath = strings.ReplaceAll(adonWaitingPath, "{id}", d.Id())

			adonWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}
			adonWaitingResp, err := adonWaitingClient.Request("GET", adonWaitingPath, &adonWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return adonWaitingResp, "COMPLETED", nil
				}
				return nil, "ERROR", err
			}

			adonWaitingRespBody, err := utils.FlattenResponse(adonWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			status := utils.PathSearch(`status.status`, adonWaitingRespBody, nil)
			if status == nil {
				return nil, "ERROR", fmt.Errorf("error parsing %s from response body", `status.phase`)
			}

			targetStatus := []string{
				"running",
			}
			if utils.StrSliceContains(targetStatus, status.(string)) {
				return adonWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"installFailed", "upgradeFailed", "deleteFailed", "rollbackFailed",
			}
			if utils.StrSliceContains(unexpectedStatus, status.(string)) {
				return adonWaitingRespBody, status.(string), nil
			}

			return adonWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

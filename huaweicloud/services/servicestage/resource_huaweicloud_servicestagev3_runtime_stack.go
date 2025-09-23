package servicestage

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	runtimeStackJsonObjectParamKeys = []string{
		"spec",
	}
	runtimeStackNonUpdatableParams = []string{
		"deploy_mode",
		"type",
	}
)

// @API ServiceStage POST /v3/{project_id}/cas/runtimestacks
// @API ServiceStage GET /v3/{project_id}/cas/runtimestacks/{runtimestack_id}
// @API ServiceStage PUT /v3/{project_id}/cas/runtimestacks/{runtimestack_id}
// @API ServiceStage DELETE /v3/{project_id}/cas/runtimestacks/{runtimestack_id}
// @API ServiceStage GET /v3/{project_id}/cas/runtimestacks
func ResourceV3RuntimeStack() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3RuntimeStackCreate,
		ReadContext:   resourceV3RuntimeStackRead,
		UpdateContext: resourceV3RuntimeStackUpdate,
		DeleteContext: resourceV3RuntimeStackDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV3RuntimeStackImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(runtimeStackNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the runtime stack is located.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the runtime stack.`,
			},
			"deploy_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The deploy mode of the runtime stack.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the runtime stack.`,
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The version of the runtime stack.`,
			},
			"spec": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: utils.SuppressObjectDiffs(),
				Description:      `The configuration of runtime stack, in JSON format.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the runtime stack.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the runtime stack.`,
			},
			"creator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator name of the runtime stack.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the runtime stack, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the runtime stack, in RFC3339 format.`,
			},
			"component_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of components associated with the runtime stack.`,
			},
			// Internal parameters/attributes.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"spec_origin": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
the new value next time the change is made. The corresponding parameter name is 'spec'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildV3RuntimeStackCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"deploy_mode": d.Get("deploy_mode"),
		"type":        d.Get("type"),
		"version":     d.Get("version"),
		"spec":        utils.StringToJson(d.Get("spec").(string)),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func resourceV3RuntimeStackCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/cas/runtimestacks"
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildV3RuntimeStackCreateBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating runtime stack: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	runtimeStackId := utils.PathSearch("id", respBody, "").(string)
	if runtimeStackId == "" {
		return diag.Errorf("failed to find the runtime stack ID from the API response")
	}
	d.SetId(runtimeStackId)

	// If the request is successful, obtain the values ​​of all JSON parameters first and save them to the
	// corresponding '_origin' attributes for subsequent determination and construction of the request body during
	// next updates.
	err = utils.RefreshObjectParamOriginValues(d, runtimeStackJsonObjectParamKeys)
	if err != nil {
		return diag.Errorf("unable to refresh the origin values: %s", err)
	}

	return resourceV3RuntimeStackRead(ctx, d, meta)
}

func GetV3RuntimeStackById(client *golangsdk.ServiceClient, runtimeStackId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/cas/runtimestacks/{runtimestack_id}"

	queryPath := client.Endpoint + httpUrl
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath = strings.ReplaceAll(queryPath, "{runtimestack_id}", runtimeStackId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	requestResp, err := client.Request("GET", queryPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceV3RuntimeStackRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		runtimeStackId = d.Id()
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	runtimeStack, err := GetV3RuntimeStackById(client, runtimeStackId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying runtime stack (%s)", runtimeStackId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", runtimeStack, nil)),
		d.Set("deploy_mode", utils.PathSearch("deploy_mode", runtimeStack, nil)),
		d.Set("type", utils.PathSearch("type", runtimeStack, nil)),
		d.Set("version", utils.PathSearch("version", runtimeStack, nil)),
		d.Set("spec", utils.JsonToString(utils.PathSearch("spec", runtimeStack, nil))),
		d.Set("description", utils.PathSearch("description", runtimeStack, nil)),
		d.Set("status", utils.PathSearch("status", runtimeStack, nil)),
		d.Set("creator", utils.PathSearch("creator", runtimeStack, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", runtimeStack, float64(0)).(float64))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", runtimeStack, float64(0)).(float64))/1000, false)),
		d.Set("component_count", utils.PathSearch("component_count", runtimeStack, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildV3RuntimeStackUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"deploy_mode": d.Get("deploy_mode"),
		"type":        d.Get("type"),
		"version":     d.Get("version"),
		"spec":        utils.StringToJson(d.Get("spec").(string)),
		"description": d.Get("description"),
	}
}

func resourceV3RuntimeStackUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		runtimeStackId = d.Id()
		httpUrl        = "v3/{project_id}/cas/runtimestacks/{runtimestack_id}"
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{runtimestack_id}", runtimeStackId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildV3RuntimeStackUpdateBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &opt)
	if err != nil {
		return diag.Errorf("error updating runtime stack (%s): %s", runtimeStackId, err)
	}

	// If the request is successful, obtain the values ​​of all JSON parameters first and save them to the
	// corresponding '_origin' attributes for subsequent determination and construction of the request body during
	// next updates.
	err = utils.RefreshObjectParamOriginValues(d, runtimeStackJsonObjectParamKeys)
	if err != nil {
		return diag.Errorf("unable to refresh the origin values: %s", err)
	}

	return resourceV3RuntimeStackRead(ctx, d, meta)
}

func resourceV3RuntimeStackDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		runtimeStackId = d.Id()
		httpUrl        = "v3/{project_id}/cas/runtimestacks/{runtimestack_id}"
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{runtimestack_id}", runtimeStackId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildV3RuntimeStackCreateBodyParams(d)),
	}

	_, err = client.Request("DELETE", deletePath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting runtime stack (%s)", runtimeStackId))
	}
	return nil
}

func resourceV3RuntimeStackImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()

	if !utils.IsUUID(importedId) {
		var (
			cfg    = meta.(*config.Config)
			region = cfg.GetRegion(d)
		)
		client, err := cfg.NewServiceClient("servicestage", region)
		if err != nil {
			return nil, fmt.Errorf("error creating ServiceStage client: %s", err)
		}
		runtimeStacks, err := listV3RuntimeStacks(client)
		if err != nil {
			return nil, err
		}
		runtimeStackId := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0].id", importedId), runtimeStacks, "").(string)
		if runtimeStackId == "" {
			return nil, fmt.Errorf("unable to find the runtime stack by its name (%s)", importedId)
		}
		d.SetId(runtimeStackId)
	}

	return []*schema.ResourceData{d}, nil
}

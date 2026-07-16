package gaussdb

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var parameterTemplateSaveNonUpdatableParams = []string{
	"config_id",
	"name",
}

// @API GaussDB POST /v3/{project_id}/configurations/{config_id}/save
// @API GaussDB GET /v3.1/{project_id}/configurations/{config_id}
// @API GaussDB PUT /v3/{project_id}/configurations/{config_id}
// @API GaussDB DELETE /v3/{project_id}/configurations/{config_id}
func ResourceParameterTemplateSave() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceParameterTemplateSaveCreate,
		ReadContext:   resourceParameterTemplateSaveRead,
		UpdateContext: resourceParameterTemplateSaveUpdate,
		DeleteContext: resourceParameterTemplateSaveDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(parameterTemplateSaveNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"config_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"configuration_parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     configTemplateParametersSchema(),
			},
		},
	}
}

func configTemplateParametersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"need_restart": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"readonly": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"value_range": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func buildCreateParameterTemplateSaveBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return bodyParams
}

func resourceParameterTemplateSaveCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/configurations/{config_id}/save"
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{config_id}", d.Get("config_id").(string))
	createOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateParameterTemplateSaveBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error saving GaussDB instance parameters as a parameter template: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	configId := utils.PathSearch("config_id", respBody, "").(string)
	if configId == "" {
		return diag.Errorf("error saving GaussDB instance parameters as a parameter template: unable to find config ID")
	}

	d.SetId(configId)

	parameterValue := d.Get("values").(map[string]interface{})
	if len(parameterValue) > 0 {
		_, err = updateParameterTemplateSaveParameters(client, d, d.Id())
		if err != nil {
			return diag.Errorf("error updating GaussDB template parameters: %s", err)
		}

		// This API is a is an asynchronous interface, but because the `job_id` not take effect,
		// so add the wait time.
		// lintignore:R018
		time.Sleep(10 * time.Second)
	}

	return resourceParameterTemplateSaveRead(ctx, d, meta)
}

func resourceParameterTemplateSaveRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	templateInfo, err := GetParameterTemplateSaveInfo(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB parameter template")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", templateInfo, nil)),
		d.Set("description", utils.PathSearch("description", templateInfo, nil)),
		d.Set("engine_version", utils.PathSearch("engine_version", templateInfo, nil)),
		d.Set("instance_mode", utils.PathSearch("instance_mode", templateInfo, nil)),
		d.Set("created_at", utils.PathSearch("created_at", templateInfo, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", templateInfo, nil)),
		d.Set("configuration_parameters", flattenTemplateParametersInfo(
			utils.PathSearch("configuration_parameters", templateInfo, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTemplateParametersInfo(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"name":         utils.PathSearch("name", v, nil),
			"value":        utils.PathSearch("value", v, nil),
			"need_restart": utils.PathSearch("need_restart", v, nil),
			"readonly":     utils.PathSearch("readonly", v, nil),
			"value_range":  utils.PathSearch("value_range", v, nil),
			"data_type":    utils.PathSearch("data_type", v, nil),
			"description":  utils.PathSearch("description", v, nil),
		})
	}

	return rst
}

func GetParameterTemplateSaveInfo(client *golangsdk.ServiceClient, templateId string) (interface{}, error) {
	httpUrl := "v3.1/{project_id}/configurations/{config_id}"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{config_id}", templateId)
	listOpts := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", listPath, &listOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceParameterTemplateSaveUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	if d.HasChanges("values", "description") {
		_, err = updateParameterTemplateSaveParameters(client, d, d.Id())
		if err != nil {
			return diag.Errorf("error updating GaussDB template parameters: %s", err)
		}

		// This API is a is an asynchronous interface, but because the `job_id` not take effect,
		// so add the wait time.
		// lintignore:R018
		time.Sleep(10 * time.Second)
	}

	return resourceParameterTemplateSaveRead(ctx, d, meta)
}

func updateParameterTemplateSaveParameters(client *golangsdk.ServiceClient, d *schema.ResourceData, templateId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/configurations/{config_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{config_id}", templateId)
	updateopt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateParameterTemplateSaveBodyParams(d)),
	}

	resp, err := client.Request("PUT", updatePath, &updateopt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func buildUpdateParameterTemplateSaveBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"values":      d.Get("values"),
		"description": d.Get("description"),
	}

	return bodyParams
}

func resourceParameterTemplateSaveDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/configurations/{config_id}"
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{config_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting GaussDB parameter template")
	}

	return nil
}

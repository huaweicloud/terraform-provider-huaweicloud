package workspace

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var appApplicationBatchPublishNonUpdatableParams = []string{
	"app_group_id",
	"applications",
	"applications.*.name",
	"applications.*.execute_path",
	"applications.*.source_type",
	"applications.*.version",
	"applications.*.command_param",
	"applications.*.work_path",
	"applications.*.icon_path",
	"applications.*.icon_index",
	"applications.*.description",
	"applications.*.publisher",
	"applications.*.source_image_ids",
	"applications.*.sandbox_enable",
	"applications.*.is_pre_boot",
	"applications.*.app_extended_info",
}

// @API Workspace POST /v1/{project_id}/app-groups/{app_group_id}/apps
func ResourceAppApplicationBatchPublish() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppApplicationBatchPublishCreate,
		ReadContext:   resourceAppApplicationBatchPublishRead,
		UpdateContext: resourceAppApplicationBatchPublishUpdate,
		DeleteContext: resourceAppApplicationBatchPublishDelete,

		CustomizeDiff: config.FlexibleForceNew(appApplicationBatchPublishNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the applications to be published are located.`,
			},
			"app_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the application group.`,
			},
			"applications": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the application.`,
						},
						"execute_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The execution path of the application.`,
						},
						"source_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The type of the application.`,
						},
						"version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The version of the application.`,
						},
						"command_param": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The command line parameters used to start the application.`,
						},
						"work_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The working directory of the application.`,
						},
						"icon_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The path where the application icon is located.`,
						},
						"icon_index": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `The icon index of the application.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The description of the application.`,
						},
						"publisher": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The publisher of the application.`,
						},
						"source_image_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of image IDs to which the application belongs.`,
						},
						"sandbox_enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether to run in sandbox mode.`,
						},
						"is_pre_boot": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether to enable application pre-boot.`,
						},
						"app_extended_info": {
							Type:        schema.TypeMap,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The extended information of the custom application.`,
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the published application.`,
						},
					},
				},
				Description: `The list of applications to be published.`,
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceAppApplicationBatchPublishCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                = meta.(*config.Config)
		applicationGroupId = d.Get("app_group_id").(string)
		httpUrl            = "v1/{project_id}/app-groups/{app_group_id}/apps"
		applications       = d.Get("applications").([]interface{})
	)

	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{app_group_id}", applicationGroupId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         utils.RemoveNil(buildApplicationBatchPublishBodyParams(applications)),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error batch publishing applications under application group (%s): %s", applicationGroupId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	result := make([]map[string]interface{}, 0, len(applications))
	for _, v := range applications {
		item := v.(map[string]interface{})
		// Set the ID of the published application.
		item["id"] = utils.PathSearch(fmt.Sprintf("items[?name=='%s']|[0].id", utils.PathSearch("name", v, "").(string)),
			respBody, nil)
		result = append(result, item)
	}

	return diag.FromErr(d.Set("applications", result))
}

func buildApplicationBatchPublishBodyParams(applications []interface{}) map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(applications))
	for _, item := range applications {
		result = append(result, map[string]interface{}{
			"name":              utils.PathSearch("name", item, nil),
			"execute_path":      utils.PathSearch("execute_path", item, nil),
			"source_type":       utils.PathSearch("source_type", item, nil),
			"version":           utils.ValueIgnoreEmpty(utils.PathSearch("version", item, nil)),
			"command_param":     utils.ValueIgnoreEmpty(utils.PathSearch("command_param", item, nil)),
			"icon_uri":          utils.ValueIgnoreEmpty(utils.PathSearch("icon_uri", item, nil)),
			"work_path":         utils.ValueIgnoreEmpty(utils.PathSearch("work_path", item, nil)),
			"icon_path":         utils.ValueIgnoreEmpty(utils.PathSearch("icon_path", item, nil)),
			"icon_index":        utils.ValueIgnoreEmpty(utils.PathSearch("icon_index", item, nil)),
			"description":       utils.ValueIgnoreEmpty(utils.PathSearch("description", item, nil)),
			"publisher":         utils.ValueIgnoreEmpty(utils.PathSearch("publisher", item, nil)),
			"source_image_ids":  utils.ValueIgnoreEmpty(utils.PathSearch("source_image_ids", item, nil)),
			"sandbox_enable":    utils.ValueIgnoreEmpty(utils.PathSearch("sandbox_enable", item, nil)),
			"is_pre_boot":       utils.ValueIgnoreEmpty(utils.PathSearch("is_pre_boot", item, nil)),
			"app_extended_info": utils.ValueIgnoreEmpty(utils.PathSearch("app_extended_info", item, nil)),
		})
	}

	return map[string]interface{}{
		"items": result,
	}
}

func resourceAppApplicationBatchPublishRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppApplicationBatchPublishUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppApplicationBatchPublishDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch publish applications. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

package codeartsbuild

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

var templateNonUpdatableParams = []string{
	"name", "description", "tool_type", "parameters",
	"steps", "steps.*.module_id", "steps.*.name", "steps.*.properties", "steps.*.version", "steps.*.enable",
}

// @API CodeArtsBuild POST /v1/template/create
// @API CodeArtsBuild GET /v1/template/custom
// @API CodeArtsBuild DELETE /v1/template/{uuid}/delete
func ResourceCodeArtsBuildTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBuildTemplateCreate,
		ReadContext:   resourceBuildTemplateRead,
		UpdateContext: resourceBuildTemplateUpdate,
		DeleteContext: resourceBuildTemplateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(templateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the build template.`,
			},
			"steps": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Specifies the build execution steps.`,
				Elem:        resourceSchemeTaskSteps(),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the template description.`,
			},
			"tool_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the tool type.`,
			},
			"parameters": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the build execution parameter list.`,
				Elem:        resourceSchemeTaskParameters(),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"favorite": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the template is favorite.`,
			},
			"nick_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the nick name.`,
			},
			"template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates ID in database.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the template type.`,
			},
			"public": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the template is public.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the template creation time.`,
			},
			"weight": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the weight of the template.`,
			},
			"scope": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the scope.`,
			},
		},
	}
}

func resourceBuildTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("codearts_build", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CodeArts Build client: %s", err)
	}

	httpUrl := "v1/template/create"
	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateBuildTemplateBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CodeArts Build template: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(createRespBody); err != nil {
		return diag.Errorf("error creating CodeArts Build template: %s", err)
	}

	id := utils.PathSearch("result.uuid", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the CodeArts Build template ID from the API response")
	}

	d.SetId(id)

	return resourceBuildTemplateRead(ctx, d, meta)
}

func buildCreateOrUpdateBuildTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
		"tool_type":   utils.ValueIgnoreEmpty(d.Get("tool_type")),
		"parameters":  buildBuildTaskParameters(d),
		"template": map[string]interface{}{
			"steps": buildBuildTaskSteps(d),
		},
	}

	return bodyParams
}

func GetBuildTemplate(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	httpUrl := "v1/template/custom"
	listPath := client.Endpoint + httpUrl
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	pageSize := 100
	page := 1
	for {
		currentPath := listPath + fmt.Sprintf("?page_size=%d&page=%d", pageSize, page)
		listResp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return nil, err
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return nil, err
		}

		templates := utils.PathSearch("result.items", listRespBody, make([]interface{}, 0)).([]interface{})
		if len(templates) == 0 {
			return nil, golangsdk.ErrDefault404{}
		}

		searchPath := fmt.Sprintf("result.items[?uuid=='%s']|[0]", id)
		result := utils.PathSearch(searchPath, listRespBody, nil)
		if result != nil {
			return result, nil
		}

		page++
	}
}

func resourceBuildTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_build", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Build client: %s", err)
	}

	template, err := GetBuildTemplate(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CodeArts Build template")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", template, nil)),
		d.Set("description", utils.PathSearch("description", template, nil)),
		d.Set("tool_type", utils.PathSearch("tool_type", template, nil)),
		d.Set("parameters", flattenBuildTaskParameters(template)),
		d.Set("steps", flattenBuildTaskSteps(d, utils.PathSearch("template", template, nil))),
		d.Set("nick_name", utils.PathSearch("nick_name", template, nil)),
		d.Set("template_id", utils.PathSearch("id", template, nil)),
		d.Set("type", utils.PathSearch("type", template, nil)),
		d.Set("public", utils.PathSearch("public", template, nil)),
		d.Set("create_time", utils.PathSearch("create_time", template, nil)),
		d.Set("weight", utils.PathSearch("weight", template, nil)),
		d.Set("scope", utils.PathSearch("scope", template, nil)),

		//nolint:misspell
		d.Set("favorite", utils.PathSearch("favourite", template, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBuildTemplateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBuildTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_build", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Build client: %s", err)
	}

	httpUrl := "v1/template/{uuid}/delete"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{uuid}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertUndefinedErrInto404Err(err, 422, "error_code", templateNotFoundErr),
			"error deleting CodeArts Build template")
	}

	return nil
}

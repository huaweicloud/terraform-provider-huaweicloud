// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DLI
// ---------------------------------------------------------------

package dli

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DLI POST /v1.0/{project_id}/sqls
// @API DLI GET /v1.0/{project_id}/sqls
// @API DLI PUT /v1.0/{project_id}/sqls/{id}
// @API DLI POST /v1.0/{project_id}/sqls-deletion
func ResourceSQLTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSQLTemplateCreate,
		UpdateContext: resourceSQLTemplateUpdate,
		ReadContext:   resourceSQLTemplateRead,
		DeleteContext: resourceSQLTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
				Description: `The name of the SQL template.`,
			},
			"sql": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The statement of the SQL template.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the SQL template.`,
			},
			"group": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The group of the SQL template.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user ID of owner.`,
			},
		},
	}
}

func resourceSQLTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createSQLTemplate: create a SQLTemplate.
	var (
		createSQLTemplateHttpUrl = "v1.0/{project_id}/sqls"
		createSQLTemplateProduct = "dli"
	)
	createSQLTemplateClient, err := cfg.NewServiceClient(createSQLTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	createSQLTemplatePath := createSQLTemplateClient.Endpoint + createSQLTemplateHttpUrl
	createSQLTemplatePath = strings.ReplaceAll(createSQLTemplatePath, "{project_id}", createSQLTemplateClient.ProjectID)

	createSQLTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createSQLTemplateOpt.JSONBody = utils.RemoveNil(buildSQLTemplateBodyParams(d))
	createSQLTemplateResp, err := createSQLTemplateClient.Request("POST", createSQLTemplatePath, &createSQLTemplateOpt)
	if err != nil {
		return diag.Errorf("error creating SQLTemplate: %s", err)
	}

	createSQLTemplateRespBody, err := utils.FlattenResponse(createSQLTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if !utils.PathSearch("is_success", createSQLTemplateRespBody, true).(bool) {
		return diag.Errorf("unable to create the SQL template: %s",
			utils.PathSearch("message", createSQLTemplateRespBody, "Message Not Found"))
	}

	templateId := utils.PathSearch("sql_id", createSQLTemplateRespBody, "").(string)
	if templateId == "" {
		return diag.Errorf("unable to find the SQL template ID the API response")
	}
	d.SetId(templateId)

	return resourceSQLTemplateRead(ctx, d, meta)
}

func buildSQLTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sql_name":    utils.ValueIgnoreEmpty(d.Get("name")),
		"sql":         utils.ValueIgnoreEmpty(d.Get("sql")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"group":       utils.ValueIgnoreEmpty(d.Get("group")),
	}
	return bodyParams
}

func resourceSQLTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getSQLTemplate: Query the SQLTemplate.
	var (
		getSQLTemplateHttpUrl = "v1.0/{project_id}/sqls"
		getSQLTemplateProduct = "dli"
	)
	getSQLTemplateClient, err := cfg.NewServiceClient(getSQLTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	getSQLTemplatePath := getSQLTemplateClient.Endpoint + getSQLTemplateHttpUrl
	getSQLTemplatePath = strings.ReplaceAll(getSQLTemplatePath, "{project_id}", getSQLTemplateClient.ProjectID)

	getSQLTemplatequeryParams := buildGetSQLTemplateQueryParams(d)
	getSQLTemplatePath += getSQLTemplatequeryParams

	getSQLTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSQLTemplateResp, err := getSQLTemplateClient.Request("GET", getSQLTemplatePath, &getSQLTemplateOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SQLTemplate")
	}

	getSQLTemplateRespBody, err := utils.FlattenResponse(getSQLTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if !utils.PathSearch("is_success", getSQLTemplateRespBody, true).(bool) {
		return diag.Errorf("unable to query the SQL templates: %s",
			utils.PathSearch("message", getSQLTemplateRespBody, "Message Not Found"))
	}

	jsonPath := fmt.Sprintf("sqls[?sql_id=='%s']|[0]", d.Id())
	sqlTemplate := utils.PathSearch(jsonPath, getSQLTemplateRespBody, nil)
	if sqlTemplate == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("sql_name", sqlTemplate, nil)),
		d.Set("sql", utils.PathSearch("sql", sqlTemplate, nil)),
		d.Set("description", utils.PathSearch("description", sqlTemplate, nil)),
		d.Set("group", utils.PathSearch("group", sqlTemplate, nil)),
		d.Set("owner", utils.PathSearch("owner", sqlTemplate, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetSQLTemplateQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&keyword=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func resourceSQLTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateSQLTemplateChanges := []string{
		"name",
		"sql",
		"description",
		"group",
	}

	if d.HasChanges(updateSQLTemplateChanges...) {
		// updateSQLTemplate: update SQLTemplate
		var (
			updateSQLTemplateHttpUrl = "v1.0/{project_id}/sqls/{id}"
			updateSQLTemplateProduct = "dli"
		)
		updateSQLTemplateClient, err := cfg.NewServiceClient(updateSQLTemplateProduct, region)
		if err != nil {
			return diag.Errorf("error creating DLI Client: %s", err)
		}

		updateSQLTemplatePath := updateSQLTemplateClient.Endpoint + updateSQLTemplateHttpUrl
		updateSQLTemplatePath = strings.ReplaceAll(updateSQLTemplatePath, "{project_id}", updateSQLTemplateClient.ProjectID)
		updateSQLTemplatePath = strings.ReplaceAll(updateSQLTemplatePath, "{id}", d.Id())

		updateSQLTemplateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateSQLTemplateOpt.JSONBody = utils.RemoveNil(buildSQLTemplateBodyParams(d))
		requestResp, err := updateSQLTemplateClient.Request("PUT", updateSQLTemplatePath, &updateSQLTemplateOpt)
		if err != nil {
			return diag.Errorf("error updating SQLTemplate: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}
		if !utils.PathSearch("is_success", respBody, true).(bool) {
			return diag.Errorf("unable to update the SQL template: %s",
				utils.PathSearch("message", respBody, "Message Not Found"))
		}
	}
	return resourceSQLTemplateRead(ctx, d, meta)
}

func resourceSQLTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteSQLTemplateHttpUrl = "v1.0/{project_id}/sqls-deletion"
		deleteSQLTemplateProduct = "dli"
	)
	deleteSQLTemplateClient, err := cfg.NewServiceClient(deleteSQLTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	deleteSQLTemplatePath := deleteSQLTemplateClient.Endpoint + deleteSQLTemplateHttpUrl
	deleteSQLTemplatePath = strings.ReplaceAll(deleteSQLTemplatePath, "{project_id}", deleteSQLTemplateClient.ProjectID)

	deleteSQLTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	deleteSQLTemplateOpt.JSONBody = utils.RemoveNil(buildDeleteSQLTemplateBodyParams(d))
	requestResp, err := deleteSQLTemplateClient.Request("POST", deleteSQLTemplatePath, &deleteSQLTemplateOpt)
	if err != nil {
		return diag.Errorf("error deleting SQLTemplate: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if !utils.PathSearch("is_success", respBody, true).(bool) {
		return diag.Errorf("unable to delete the SQL template: %s",
			utils.PathSearch("message", respBody, "Message Not Found"))
	}

	return nil
}

func buildDeleteSQLTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sql_ids": []string{d.Id()},
	}
	return bodyParams
}

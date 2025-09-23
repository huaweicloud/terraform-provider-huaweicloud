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

// @API DLI POST /v1.0/{project_id}/streaming/job-templates
// @API DLI GET /v1.0/{project_id}/streaming/job-templates
// @API DLI PUT /v1.0/{project_id}/streaming/job-templates/{id}
// @API DLI DELETE /v1.0/{project_id}/streaming/job-templates/{id}
func ResourceFlinkTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFlinkTemplateCreate,
		UpdateContext: resourceFlinkTemplateUpdate,
		ReadContext:   resourceFlinkTemplateRead,
		DeleteContext: resourceFlinkTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

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
				Description: `The name of the flink template.`,
			},
			"sql": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The statement of the flink template.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the flink template.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The type of the flink template.`,
			},
			"tags": common.TagsForceNewSchema(),
		},
	}
}

func resourceFlinkTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createFlinkTemplate: create a Flink template.
	var (
		createFlinkTemplateHttpUrl = "v1.0/{project_id}/streaming/job-templates"
		createFlinkTemplateProduct = "dli"
	)
	createFlinkTemplateClient, err := cfg.NewServiceClient(createFlinkTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	createFlinkTemplatePath := createFlinkTemplateClient.Endpoint + createFlinkTemplateHttpUrl
	createFlinkTemplatePath = strings.ReplaceAll(createFlinkTemplatePath, "{project_id}", createFlinkTemplateClient.ProjectID)

	createFlinkTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createFlinkTemplateOpt.JSONBody = utils.RemoveNil(buildCreateFlinkTemplateBodyParams(d))
	createFlinkTemplateResp, err := createFlinkTemplateClient.Request("POST", createFlinkTemplatePath, &createFlinkTemplateOpt)
	if err != nil {
		return diag.Errorf("error creating FlinkTemplate: %s", err)
	}

	createFlinkTemplateRespBody, err := utils.FlattenResponse(createFlinkTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if !utils.PathSearch("is_success", createFlinkTemplateRespBody, true).(bool) {
		return diag.Errorf("unable to create the flink template: %s",
			utils.PathSearch("message", createFlinkTemplateRespBody, "Message Not Found"))
	}

	templateId := utils.PathSearch("template.template_id", createFlinkTemplateRespBody, float64(0)).(float64)
	if templateId == 0 {
		return diag.Errorf("unable to find the Flink template ID from the API response")
	}
	d.SetId(fmt.Sprint(templateId))

	return resourceFlinkTemplateRead(ctx, d, meta)
}

func buildCreateFlinkTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":     utils.ValueIgnoreEmpty(d.Get("name")),
		"sql_body": utils.ValueIgnoreEmpty(d.Get("sql")),
		"desc":     utils.ValueIgnoreEmpty(d.Get("description")),
		"job_type": utils.ValueIgnoreEmpty(d.Get("type")),
		"tags":     utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}
	return bodyParams
}

func resourceFlinkTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getFlinkTemplate: Query the Flink template.
	var (
		getFlinkTemplateHttpUrl = "v1.0/{project_id}/streaming/job-templates"
		getFlinkTemplateProduct = "dli"
	)
	getFlinkTemplateClient, err := cfg.NewServiceClient(getFlinkTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	getFlinkTemplatePath := getFlinkTemplateClient.Endpoint + getFlinkTemplateHttpUrl
	getFlinkTemplatePath = strings.ReplaceAll(getFlinkTemplatePath, "{project_id}", getFlinkTemplateClient.ProjectID)

	getFlinkTemplatequeryParams := buildGetFlinkTemplateQueryParams(d)
	getFlinkTemplatePath += getFlinkTemplatequeryParams

	getFlinkTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getFlinkTemplateResp, err := getFlinkTemplateClient.Request("GET", getFlinkTemplatePath, &getFlinkTemplateOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving FlinkTemplate")
	}

	getFlinkTemplateRespBody, err := utils.FlattenResponse(getFlinkTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if !utils.PathSearch("is_success", getFlinkTemplateRespBody, true).(bool) {
		return diag.Errorf("unable to query the flink templates: %s",
			utils.PathSearch("message", getFlinkTemplateRespBody, "Message Not Found"))
	}

	jsonPath := fmt.Sprintf("template_list.templates[?template_id==`%s`]|[0]", d.Id())
	flinkTemplate := utils.PathSearch(jsonPath, getFlinkTemplateRespBody, nil)
	if flinkTemplate == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", flinkTemplate, nil)),
		d.Set("sql", utils.PathSearch("sql_body", flinkTemplate, nil)),
		d.Set("description", utils.PathSearch("desc", flinkTemplate, nil)),
		d.Set("type", utils.PathSearch("job_type", flinkTemplate, nil)),
		d.Set("tags", d.Get("tags")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetFlinkTemplateQueryParams(d *schema.ResourceData) string {
	res := "?&limit=100"
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}

	return res
}

func resourceFlinkTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateFlinkTemplateChanges := []string{
		"name",
		"sql",
		"description",
	}

	if d.HasChanges(updateFlinkTemplateChanges...) {
		// updateFlinkTemplate: update Flink template
		var (
			updateFlinkTemplateHttpUrl = "v1.0/{project_id}/streaming/job-templates/{id}"
			updateFlinkTemplateProduct = "dli"
		)
		updateFlinkTemplateClient, err := cfg.NewServiceClient(updateFlinkTemplateProduct, region)
		if err != nil {
			return diag.Errorf("error creating DLI Client: %s", err)
		}

		updateFlinkTemplatePath := updateFlinkTemplateClient.Endpoint + updateFlinkTemplateHttpUrl
		updateFlinkTemplatePath = strings.ReplaceAll(updateFlinkTemplatePath, "{project_id}", updateFlinkTemplateClient.ProjectID)
		updateFlinkTemplatePath = strings.ReplaceAll(updateFlinkTemplatePath, "{id}", d.Id())

		updateFlinkTemplateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateFlinkTemplateOpt.JSONBody = utils.RemoveNil(buildUpdateFlinkTemplateBodyParams(d))
		requestResp, err := updateFlinkTemplateClient.Request("PUT", updateFlinkTemplatePath, &updateFlinkTemplateOpt)
		if err != nil {
			return diag.Errorf("error updating FlinkTemplate: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}
		if !utils.PathSearch("is_success", respBody, true).(bool) {
			return diag.Errorf("unable to update the flink template: %s",
				utils.PathSearch("message", respBody, "Message Not Found"))
		}
	}
	return resourceFlinkTemplateRead(ctx, d, meta)
}

func buildUpdateFlinkTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":     utils.ValueIgnoreEmpty(d.Get("name")),
		"sql_body": utils.ValueIgnoreEmpty(d.Get("sql")),
		"desc":     utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceFlinkTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteFlinkTemplate: delete Flink template
	var (
		deleteFlinkTemplateHttpUrl = "v1.0/{project_id}/streaming/job-templates/{id}"
		deleteFlinkTemplateProduct = "dli"
	)
	deleteFlinkTemplateClient, err := cfg.NewServiceClient(deleteFlinkTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	deleteFlinkTemplatePath := deleteFlinkTemplateClient.Endpoint + deleteFlinkTemplateHttpUrl
	deleteFlinkTemplatePath = strings.ReplaceAll(deleteFlinkTemplatePath, "{project_id}", deleteFlinkTemplateClient.ProjectID)
	deleteFlinkTemplatePath = strings.ReplaceAll(deleteFlinkTemplatePath, "{id}", d.Id())

	deleteFlinkTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	requestResp, err := deleteFlinkTemplateClient.Request("DELETE", deleteFlinkTemplatePath, &deleteFlinkTemplateOpt)
	if err != nil {
		return diag.Errorf("error deleting FlinkTemplate: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if !utils.PathSearch("is_success", respBody, true).(bool) {
		return diag.Errorf("unable to delete the flink template: %s",
			utils.PathSearch("message", respBody, "Message Not Found"))
	}

	return nil
}

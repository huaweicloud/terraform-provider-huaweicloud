package rgc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var templateNonUpdatableParams = []string{"template_name", "template_type", "template_description", "template_body"}

// @API RGC POST /v1/rgc/templates
// @API RGC DELETE /v1/rgc/templates/{template_name}
// @API RFS GET /v1/{project_id}/templates
func ResourceTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTemplateCreate,
		UpdateContext: resourceTemplateUpdate,
		ReadContext:   resourceTemplateRead,
		DeleteContext: resourceTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(templateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_body": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createTemplateHttpUrl = "v1/rgc/templates"
		createTemplateProduct = "rgc"
	)

	createTemplateClient, err := cfg.NewServiceClient(createTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	createTemplatePath := createTemplateClient.Endpoint + createTemplateHttpUrl

	createTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildTemplateBodyParams(d)),
	}
	createTemplateResp, err := createTemplateClient.Request("POST", createTemplatePath, &createTemplateOpt)
	if err != nil {
		return diag.Errorf("error creating template: %s", err)
	}

	createTemplateRespBody, err := utils.FlattenResponse(createTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	templateId := utils.PathSearch("template_id", createTemplateRespBody, "").(string)
	if templateId == "" {
		return diag.Errorf("error get template id from API response.")
	}
	d.SetId(templateId)
	return resourceTemplateRead(ctx, d, meta)
}

func resourceTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listTemplateHttpUrl = "v1/{project_id}/templates"
		listTemplateProduct = "rfs"
	)

	listTemplateClient, err := cfg.NewServiceClient(listTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	projectId := listTemplateClient.ProjectID

	var templateMeta interface{}
	var marker string

	for {
		listTemplatePath := listTemplateClient.Endpoint + listTemplateHttpUrl + buildTemplateQueryParams(marker)
		listTemplatePath = strings.ReplaceAll(listTemplatePath, "{project_id}", projectId)

		listTemplateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Client-Request-Id": randUUID,
				"Content-Type":      "application/json",
				"X-Language":        "en-us",
			},
		}

		listTemplateResp, err := listTemplateClient.Request("GET", listTemplatePath, &listTemplateOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving template metadata")
		}

		listTemplateRespBody, err := utils.FlattenResponse(listTemplateResp)
		if err != nil {
			return diag.FromErr(err)
		}

		templateMeta = utils.PathSearch(fmt.Sprintf("templates[?template_id=='%s']|[0]", d.Id()), listTemplateRespBody, nil)

		if templateMeta != nil {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", listTemplateRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	if templateMeta == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("create_time", utils.PathSearch("create_time", templateMeta, nil)),
		d.Set("update_time", utils.PathSearch("create_time", templateMeta, nil)),
		d.Set("latest_version_id", utils.PathSearch("create_time", templateMeta, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceTemplateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteTemplateHttpUrl = "v1/rgc/templates/{template_name}"
		deleteTemplateProduct = "rgc"
	)

	deleteTemplateClient, err := cfg.NewServiceClient(deleteTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	deleteTemplatePath := deleteTemplateClient.Endpoint + deleteTemplateHttpUrl
	deleteTemplatePath = strings.ReplaceAll(deleteTemplatePath, "{template_name}", d.Get("template_name").(string))

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	deleteTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Client-Request-Id": randUUID,
			"Content-Type":      "application/json",
			"X-Language":        "en-us",
		},
		OkCodes: []int{
			204,
		},
	}

	_, err = deleteTemplateClient.Request("DELETE", deleteTemplatePath, &deleteTemplateOpt)
	if err != nil {
		return diag.Errorf("error deleting template: %s", err)
	}

	return nil
}

func buildTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"template_name":        d.Get("template_name").(string),
		"template_type":        d.Get("template_type").(string),
		"template_body":        utils.ValueIgnoreEmpty(d.Get("template_body")),
		"template_description": utils.ValueIgnoreEmpty(d.Get("template_description")),
	}

	return bodyParams
}

func buildTemplateQueryParams(marker string) string {
	// the default value of limit is 200
	res := "?limit=200"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

package rfs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RFS POST /v1/{project_id}/templates
// @API RFS GET /v1/{project_id}/templates
// @API RFS PATCH /v1/{project_id}/templates/{template_name}/metadata
// @API RFS DELETE /v1/{project_id}/templates/{template_name}
func ResourceRfsTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRfsTemplateCreate,
		UpdateContext: resourceRfsTemplateUpdate,
		ReadContext:   resourceRfsTemplateRead,
		DeleteContext: resourceRfsTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"template_name",
			"template_body",
			"template_uri",
			"version_description",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_body": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_description": {
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
			"template_id": {
				Type:     schema.TypeString,
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
			"latest_version_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRfsTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/templates"
		product = "rfs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request ID: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId,
			"Content-Type":      "application/json",
		},
		JSONBody: utils.RemoveNil(buildCreateRfsTemplateBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating RFS template: %s", err)
	}
	_, err = utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	templateName := d.Get("template_name").(string)
	d.SetId(templateName)

	return resourceRfsTemplateRead(ctx, d, meta)
}

func buildCreateRfsTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"template_name":        d.Get("template_name"),
		"template_description": utils.ValueIgnoreEmpty(d.Get("template_description")),
		"version_description":  utils.ValueIgnoreEmpty(d.Get("version_description")),
		"template_body":        utils.ValueIgnoreEmpty(d.Get("template_body")),
		"template_uri":         utils.ValueIgnoreEmpty(d.Get("template_uri")),
	}

	return bodyParams
}

func GetRfsTemplateByName(client *golangsdk.ServiceClient, templateName string, requestId string) (interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/templates"
		limit   = 1000
		marker  = ""
	)

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = fmt.Sprintf("%s?limit=%d", requestPath, limit)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":      "application/json",
			"Client-Request-Id": requestId,
		},
	}

	for {
		requestPathWithMarker := requestPath
		if marker != "" {
			requestPathWithMarker = fmt.Sprintf("%s&marker=%s", requestPathWithMarker, marker)
		}

		resp, err := client.Request("GET", requestPathWithMarker, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		template := utils.PathSearch(fmt.Sprintf("templates[?template_name=='%s']|[0]", templateName), respBody, nil)
		if template != nil {
			return template, nil
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

func resourceRfsTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		product      = "rfs"
		templateName = d.Id()
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request ID: %s", err)
	}

	template, err := GetRfsTemplateByName(client, templateName, requestId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RFS template")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("template_id", utils.PathSearch("template_id", template, "")),
		d.Set("template_name", utils.PathSearch("template_name", template, "")),
		d.Set("template_description", utils.PathSearch("template_description", template, "")),
		d.Set("create_time", utils.PathSearch("create_time", template, "")),
		d.Set("update_time", utils.PathSearch("update_time", template, "")),
		d.Set("latest_version_description", utils.PathSearch("latest_version_description", template, "")),
		d.Set("latest_version_id", utils.PathSearch("latest_version_id", template, "")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceRfsTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v1/{project_id}/templates/{template_name}/metadata"
		product      = "rfs"
		templateName = d.Id()
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request ID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{template_name}", templateName)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId,
			"Content-Type":      "application/json",
		},
		JSONBody: utils.RemoveNil(buildUpdateRfsTemplateBodyParams(d)),
	}

	_, err = client.Request("PATCH", requestPath, &opt)
	if err != nil {
		return diag.Errorf("error updating RFS template: %s", err)
	}

	return resourceRfsTemplateRead(ctx, d, meta)
}

func buildUpdateRfsTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"template_description": utils.ValueIgnoreEmpty(d.Get("template_description")),
	}

	return bodyParams
}

func resourceRfsTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v1/{project_id}/templates/{template_name}"
		product      = "rfs"
		templateName = d.Id()
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request ID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{template_name}", templateName)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId,
			"Content-Type":      "application/json",
		},
	}

	_, err = client.Request("DELETE", requestPath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting RFS template")
	}

	return nil
}

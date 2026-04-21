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

// @API RFS POST /v1/{project_id}/templates/{template_name}/versions
// @API RFS GET /v1/{project_id}/templates/{template_name}/versions/{version_id}/metadata
// @API RFS DELETE /v1/{project_id}/templates/{template_name}/versions/{version_id}
func ResourceRfsTemplateVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRfsTemplateVersionCreate,
		UpdateContext: resourceRfsTemplateVersionUpdate,
		ReadContext:   resourceRfsTemplateVersionRead,
		DeleteContext: resourceRfsTemplateVersionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRfsTemplateVersionImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"template_name",
			"template_body",
			"template_uri",
			"template_id",
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
			"template_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"version_description": {
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
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateRfsTemplateVersionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"version_description": utils.ValueIgnoreEmpty(d.Get("version_description")),
		"template_body":       utils.ValueIgnoreEmpty(d.Get("template_body")),
		"template_uri":        utils.ValueIgnoreEmpty(d.Get("template_uri")),
	}

	return bodyParams
}

func buildTemplateVersionQueryParams(templateId string) string {
	if templateId == "" {
		return ""
	}

	return fmt.Sprintf("?template_id=%s", templateId)
}

func resourceRfsTemplateVersionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v1/{project_id}/templates/{template_name}/versions"
		product      = "rfs"
		templateName = d.Get("template_name").(string)
		templateId   = d.Get("template_id").(string)
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
	createPath = strings.ReplaceAll(createPath, "{template_name}", templateName)
	createPath += buildTemplateVersionQueryParams(templateId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId,
			"Content-Type":      "application/json",
		},
		JSONBody: utils.RemoveNil(buildCreateRfsTemplateVersionBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating RFS template version: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	versionId := utils.PathSearch("version_id", respBody, "").(string)
	if versionId == "" {
		return diag.Errorf("error creating RFS template version: version_id is not found in API response")
	}
	d.SetId(versionId)

	return resourceRfsTemplateVersionRead(ctx, d, meta)
}

func QueryRfsTemplateVersion(client *golangsdk.ServiceClient, templateName, versionId, uuid string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/templates/{template_name}/versions/{version_id}/metadata"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{template_name}", templateName)
	requestPath = strings.ReplaceAll(requestPath, "{version_id}", versionId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":      "application/json",
			"Client-Request-Id": uuid,
		},
	}

	resp, err := client.Request("GET", requestPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceRfsTemplateVersionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		product      = "rfs"
		templateName = d.Get("template_name").(string)
		versionId    = d.Id()
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request ID: %s", err)
	}

	respBody, err := QueryRfsTemplateVersion(client, templateName, versionId, requestId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RFS template version")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("template_name", utils.PathSearch("template_name", respBody, "")),
		d.Set("template_id", utils.PathSearch("template_id", respBody, "")),
		d.Set("create_time", utils.PathSearch("create_time", respBody, "")),
		d.Set("version_description", utils.PathSearch("version_description", respBody, "")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceRfsTemplateVersionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// Template versions are immutable and append-only.
	return nil
}

func resourceRfsTemplateVersionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v1/{project_id}/templates/{template_name}/versions/{version_id}"
		product      = "rfs"
		templateName = d.Get("template_name").(string)
		versionId    = d.Id()
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
	requestPath = strings.ReplaceAll(requestPath, "{version_id}", versionId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId,
			"Content-Type":      "application/json",
		},
	}

	_, err = client.Request("DELETE", requestPath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting RFS template version")
	}

	return nil
}

func resourceRfsTemplateVersionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import ID format, expected '<template_name>/<id>', but got: %s", d.Id())
	}

	templateName := parts[0]
	versionId := parts[1]

	d.SetId(versionId)

	if err := d.Set("template_name", templateName); err != nil {
		return nil, fmt.Errorf("error setting template_name during import: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}

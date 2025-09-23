// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product SMN
// ---------------------------------------------------------------

package smn

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SMN POST /v2/{project_id}/notifications/message_template
// @API SMN DELETE /v2/{project_id}/notifications/message_template/{message_template_id}
// @API SMN GET /v2/{project_id}/notifications/message_template/{message_template_id}
// @API SMN PUT /v2/{project_id}/notifications/message_template/{message_template_id}
func ResourceSmnMessageTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSmnMessageTemplateCreate,
		UpdateContext: resourceSmnMessageTemplateUpdate,
		ReadContext:   resourceSmnMessageTemplateRead,
		DeleteContext: resourceSmnMessageTemplateDelete,
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
				ForceNew:    true,
				Description: `Specifies the message template name.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the protocol supported by the template.`,
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the template content, which supports plain text only.`,
			},
			"tag_names": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Indicates the variable list.`,
			},
		},
	}
}

func resourceSmnMessageTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createMessageTemplate: create SMN message template
	var (
		createMessageTemplateHttpUrl = "v2/{project_id}/notifications/message_template"
		createMessageTemplateProduct = "smn"
	)
	createMessageTemplateClient, err := cfg.NewServiceClient(createMessageTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating SMN Client: %s", err)
	}

	createMessageTemplatePath := createMessageTemplateClient.Endpoint + createMessageTemplateHttpUrl
	createMessageTemplatePath = strings.ReplaceAll(createMessageTemplatePath, "{project_id}",
		createMessageTemplateClient.ProjectID)

	createMessageTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createMessageTemplateOpt.JSONBody = utils.RemoveNil(buildCreateMessageTemplateBodyParams(d))
	createMessageTemplateResp, err := createMessageTemplateClient.Request("POST",
		createMessageTemplatePath, &createMessageTemplateOpt)
	if err != nil {
		return diag.Errorf("error creating SMN message template: %s", err)
	}

	createMessageTemplateRespBody, err := utils.FlattenResponse(createMessageTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("message_template_id", createMessageTemplateRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SMN message template: ID is not found in API response")
	}
	d.SetId(id)

	return resourceSmnMessageTemplateRead(ctx, d, meta)
}

func buildCreateMessageTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"message_template_name": utils.ValueIgnoreEmpty(d.Get("name")),
		"protocol":              utils.ValueIgnoreEmpty(d.Get("protocol")),
		"content":               utils.ValueIgnoreEmpty(d.Get("content")),
	}
	return bodyParams
}

func resourceSmnMessageTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateMessageTemplateHasChanges := []string{
		"content",
	}

	if d.HasChanges(updateMessageTemplateHasChanges...) {
		// updateMessageTemplate: update SMN message template
		var (
			updateMessageTemplateHttpUrl = "v2/{project_id}/notifications/message_template/{message_template_id}"
			updateMessageTemplateProduct = "smn"
		)
		updateMessageTemplateClient, err := cfg.NewServiceClient(updateMessageTemplateProduct, region)
		if err != nil {
			return diag.Errorf("error creating SMN Client: %s", err)
		}

		updateMessageTemplatePath := updateMessageTemplateClient.Endpoint + updateMessageTemplateHttpUrl
		updateMessageTemplatePath = strings.ReplaceAll(updateMessageTemplatePath, "{project_id}",
			updateMessageTemplateClient.ProjectID)
		updateMessageTemplatePath = strings.ReplaceAll(updateMessageTemplatePath, "{message_template_id}", d.Id())

		updateMessageTemplateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		updateMessageTemplateOpt.JSONBody = utils.RemoveNil(buildUpdateMessageTemplateBodyParams(d))
		_, err = updateMessageTemplateClient.Request("PUT", updateMessageTemplatePath, &updateMessageTemplateOpt)
		if err != nil {
			return diag.Errorf("error updating SMN message template: %s", err)
		}
	}
	return resourceSmnMessageTemplateRead(ctx, d, meta)
}

func buildUpdateMessageTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"content": utils.ValueIgnoreEmpty(d.Get("content")),
	}
	return bodyParams
}

func resourceSmnMessageTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getMessageTemplate: Query SMN message template
	var (
		getMessageTemplateHttpUrl = "v2/{project_id}/notifications/message_template/{message_template_id}"
		getMessageTemplateProduct = "smn"
	)
	getMessageTemplateClient, err := cfg.NewServiceClient(getMessageTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating SMN Client: %s", err)
	}

	getMessageTemplatePath := getMessageTemplateClient.Endpoint + getMessageTemplateHttpUrl
	getMessageTemplatePath = strings.ReplaceAll(getMessageTemplatePath, "{project_id}",
		getMessageTemplateClient.ProjectID)
	getMessageTemplatePath = strings.ReplaceAll(getMessageTemplatePath, "{message_template_id}", d.Id())

	getMessageTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getMessageTemplateResp, err := getMessageTemplateClient.Request("GET",
		getMessageTemplatePath, &getMessageTemplateOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SMN message template")
	}

	getMessageTemplateRespBody, err := utils.FlattenResponse(getMessageTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("message_template_name",
			getMessageTemplateRespBody, nil)),
		d.Set("protocol", utils.PathSearch("protocol", getMessageTemplateRespBody, nil)),
		d.Set("tag_names", utils.PathSearch("tag_names", getMessageTemplateRespBody, nil)),
		d.Set("content", utils.PathSearch("content", getMessageTemplateRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSmnMessageTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteMessageTemplate: Delete SMN message template
	var (
		deleteMessageTemplateHttpUrl = "v2/{project_id}/notifications/message_template/{message_template_id}"
		deleteMessageTemplateProduct = "smn"
	)
	deleteMessageTemplateClient, err := cfg.NewServiceClient(deleteMessageTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating SMN Client: %s", err)
	}

	deleteMessageTemplatePath := deleteMessageTemplateClient.Endpoint + deleteMessageTemplateHttpUrl
	deleteMessageTemplatePath = strings.ReplaceAll(deleteMessageTemplatePath, "{project_id}",
		deleteMessageTemplateClient.ProjectID)
	deleteMessageTemplatePath = strings.ReplaceAll(deleteMessageTemplatePath, "{message_template_id}", d.Id())

	deleteMessageTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteMessageTemplateClient.Request("DELETE", deleteMessageTemplatePath, &deleteMessageTemplateOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SMN message template")
	}

	return nil
}

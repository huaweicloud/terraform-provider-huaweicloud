// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product LTS
// ---------------------------------------------------------------

package lts

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

// @API LTS GET /v2/{project_id}/{domain_id}/lts/events/notification/template/{template_name}
// @API LTS DELETE /v2/{project_id}/{domain_id}/lts/events/notification/templates
// @API LTS POST /v2/{project_id}/{domain_id}/lts/events/notification/templates
// @API LTS PUT /v2/{project_id}/{domain_id}/lts/events/notification/templates
func ResourceNotificationTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotificationTemplateCreate,
		UpdateContext: resourceNotificationTemplateUpdate,
		ReadContext:   resourceNotificationTemplateRead,
		DeleteContext: resourceNotificationTemplateDelete,
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
				Description: `The name of the notification template.`,
			},
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The source of the notification template.`,
			},
			"locale": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Language.`,
			},
			"templates": {
				Type:        schema.TypeList,
				Elem:        notificationTemplateSubTemplateSchema(),
				Required:    true,
				Description: `The list of notification template body.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the notification template.`,
			},
		},
	}
}

func notificationTemplateSubTemplateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"sub_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the sub-template.`,
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The content of the sub-template.`,
			},
		},
	}
	return &sc
}

func resourceNotificationTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createNotificationTemplate: create a LTS notification template.
	var (
		createNotificationTemplateHttpUrl = "v2/{project_id}/{domain_id}/lts/events/notification/templates"
		createNotificationTemplateProduct = "lts"
	)
	createNotificationTemplateClient, err := cfg.NewServiceClient(createNotificationTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	createNotificationTemplatePath := createNotificationTemplateClient.Endpoint + createNotificationTemplateHttpUrl
	createNotificationTemplatePath = strings.ReplaceAll(createNotificationTemplatePath, "{project_id}",
		createNotificationTemplateClient.ProjectID)
	createNotificationTemplatePath = strings.ReplaceAll(createNotificationTemplatePath, "{domain_id}", cfg.DomainID)

	createNotificationTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createNotificationTemplateOpt.JSONBody = utils.RemoveNil(buildNotificationTemplateBodyParams(d))
	createNotificationTemplateResp, err := createNotificationTemplateClient.Request("POST", createNotificationTemplatePath,
		&createNotificationTemplateOpt)
	if err != nil {
		return diag.Errorf("error creating notification template: %s", err)
	}

	_, err = utils.FlattenResponse(createNotificationTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("name").(string))

	return resourceNotificationTemplateRead(ctx, d, meta)
}

func buildNotificationTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":      d.Get("name"),
		"desc":      utils.ValueIgnoreEmpty(d.Get("description")),
		"source":    d.Get("source"),
		"locale":    d.Get("locale"),
		"templates": buildNotificationTemplateBodySubTemplate(d.Get("templates")),
	}
	return bodyParams
}

func buildNotificationTemplateBodySubTemplate(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				rst[i] = map[string]interface{}{
					"sub_type": utils.ValueIgnoreEmpty(raw["sub_type"]),
					"content":  utils.ValueIgnoreEmpty(raw["content"]),
				}
			}
		}
		return rst
	}
	return nil
}

func resourceNotificationTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getNotificationTemplate: Query the LTS notification template.
	var (
		getNotificationTemplateHttpUrl = "v2/{project_id}/{domain_id}/lts/events/notification/template/{id}"
		getNotificationTemplateProduct = "lts"
	)
	getNotificationTemplateClient, err := cfg.NewServiceClient(getNotificationTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	getNotificationTemplatePath := getNotificationTemplateClient.Endpoint + getNotificationTemplateHttpUrl
	getNotificationTemplatePath = strings.ReplaceAll(getNotificationTemplatePath, "{project_id}", getNotificationTemplateClient.ProjectID)
	getNotificationTemplatePath = strings.ReplaceAll(getNotificationTemplatePath, "{domain_id}", cfg.DomainID)
	getNotificationTemplatePath = strings.ReplaceAll(getNotificationTemplatePath, "{id}", d.Id())

	getNotificationTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getNotificationTemplateResp, err := getNotificationTemplateClient.Request("GET", getNotificationTemplatePath,
		&getNotificationTemplateOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving notification template")
	}

	getNotificationTemplateRespBody, err := utils.FlattenResponse(getNotificationTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getNotificationTemplateRespBody, nil)),
		d.Set("description", utils.PathSearch("desc", getNotificationTemplateRespBody, nil)),
		d.Set("source", utils.PathSearch("source", getNotificationTemplateRespBody, nil)),
		d.Set("locale", utils.PathSearch("locale", getNotificationTemplateRespBody, nil)),
		d.Set("templates", flattenNotificationTemplateBodySubTemplate(getNotificationTemplateRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenNotificationTemplateBodySubTemplate(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("templates", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"sub_type": utils.PathSearch("sub_type", v, nil),
			"content":  utils.PathSearch("content", v, nil),
		})
	}
	return rst
}

func resourceNotificationTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateNotificationTemplate: update the LTS notification template.
	var (
		updateNotificationTemplateHttpUrl = "v2/{project_id}/{domain_id}/lts/events/notification/templates"
		updateNotificationTemplateProduct = "lts"
	)
	updateNotificationTemplateClient, err := cfg.NewServiceClient(updateNotificationTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	updateNotificationTemplatePath := updateNotificationTemplateClient.Endpoint + updateNotificationTemplateHttpUrl
	updateNotificationTemplatePath = strings.ReplaceAll(updateNotificationTemplatePath, "{project_id}",
		updateNotificationTemplateClient.ProjectID)
	updateNotificationTemplatePath = strings.ReplaceAll(updateNotificationTemplatePath, "{domain_id}", cfg.DomainID)

	updateNotificationTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	updateNotificationTemplateOpt.JSONBody = utils.RemoveNil(buildNotificationTemplateBodyParams(d))
	_, err = updateNotificationTemplateClient.Request("PUT", updateNotificationTemplatePath, &updateNotificationTemplateOpt)
	if err != nil {
		return diag.Errorf("error updating notification template: %s", err)
	}
	return resourceNotificationTemplateRead(ctx, d, meta)
}

func resourceNotificationTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteNotificationTemplate: delete LTS notification template
	var (
		deleteNotificationTemplateHttpUrl = "v2/{project_id}/{domain_id}/lts/events/notification/templates"
		deleteNotificationTemplateProduct = "lts"
	)
	deleteNotificationTemplateClient, err := cfg.NewServiceClient(deleteNotificationTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	deleteNotificationTemplatePath := deleteNotificationTemplateClient.Endpoint + deleteNotificationTemplateHttpUrl
	deleteNotificationTemplatePath = strings.ReplaceAll(deleteNotificationTemplatePath, "{project_id}", deleteNotificationTemplateClient.ProjectID)
	deleteNotificationTemplatePath = strings.ReplaceAll(deleteNotificationTemplatePath, "{domain_id}", cfg.DomainID)

	deleteNotificationTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	deleteNotificationTemplateOpt.JSONBody = utils.RemoveNil(buildDeleteNotificationTemplateBodyParams(d))
	_, err = deleteNotificationTemplateClient.Request("DELETE", deleteNotificationTemplatePath, &deleteNotificationTemplateOpt)
	if err != nil {
		return diag.Errorf("error deleting notification template: %s", err)
	}

	return nil
}

func buildDeleteNotificationTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"template_names": []string{d.Id()},
	}
	return bodyParams
}

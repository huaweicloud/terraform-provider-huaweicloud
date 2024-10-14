package aom

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AOM POST /v2/{project_id}/events/notification/templates
// @API AOM PUT /v2/{project_id}/events/notification/templates
// @API AOM DELETE /v2/{project_id}/events/notification/templates
// @API AOM GET /v2/{project_id}/events/notification/template/{name}
func ResourceMessageTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMessageTemplateCreate,
		ReadContext:   resourceMessageTemplateRead,
		UpdateContext: resourceMessageTemplateUpdate,
		DeleteContext: resourceMessageTemplateDelete,

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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"templates": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sub_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"topic": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"locale": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMessageTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createHttpUrl := "v2/{project_id}/events/notification/templates"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		MoreHeaders: buildHeaders(cfg, d),
		JSONBody:    utils.RemoveNil(buildCreateMessageTemplateBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating AOM message template: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceMessageTemplateRead(ctx, d, meta)
}

func buildCreateMessageTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":      d.Get("name"),
		"templates": buildMessageTemplateDetail(d),
		"locale":    d.Get("locale"),
		"source":    utils.ValueIgnoreEmpty(d.Get("source")),
		"desc":      utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return bodyParams
}

func buildMessageTemplateDetail(d *schema.ResourceData) interface{} {
	templates := d.Get("templates").(*schema.Set).List()
	rst := make([]interface{}, 0, len(templates))
	for _, v := range templates {
		params := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"subType": params["sub_type"],
			"content": params["content"],
			"version": utils.ValueIgnoreEmpty(params["version"]),
			"topic":   utils.ValueIgnoreEmpty(params["topic"]),
		})
	}

	jsonTemplates, err := json.Marshal(rst)
	if err != nil {
		log.Printf("[ERROR] unable to convert the templates into JSON encoding")
		return ""
	}

	return string(jsonTemplates)
}

func resourceMessageTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	template, err := GetMessageTemplate(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "AOM.08025006"),
			"error retrieving message template")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", template, nil)),
		d.Set("source", utils.PathSearch("source", template, nil)),
		d.Set("templates", flattenMessageTemplateDetailTemplates(utils.PathSearch("templates", template, nil))),
		d.Set("locale", utils.PathSearch("locale", template, nil)),
		d.Set("description", utils.PathSearch("desc", template, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", template, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_time", template, float64(0)).(float64))/1000, true)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("modify_time", template, float64(0)).(float64))/1000, true)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetMessageTemplate(client *golangsdk.ServiceClient, name string) (interface{}, error) {
	getHttpUrl := "v2/{project_id}/events/notification/template/{name}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{name}", name)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":          "application/json",
			"Enterprise-Project-Id": "all_granted_eps",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func flattenMessageTemplateDetailTemplates(params interface{}) []interface{} {
	if params == nil {
		return nil
	}

	var templates []interface{}
	err := json.Unmarshal([]byte(params.(string)), &templates)
	if err != nil {
		return nil
	}

	rst := make([]interface{}, 0, len(templates))
	for _, template := range templates {
		rst = append(rst, map[string]interface{}{
			"sub_type": utils.PathSearch("subType", template, nil),
			"version":  utils.PathSearch("version", template, nil),
			"topic":    utils.PathSearch("topic", template, nil),
			"content":  utils.PathSearch("content", template, nil),
		})
	}

	return rst
}

func resourceMessageTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	updateHttpUrl := "v2/{project_id}/events/notification/templates"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		MoreHeaders: buildHeaders(cfg, d),
		JSONBody:    buildCreateMessageTemplateBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating message template: %s", err)
	}

	return resourceMessageTemplateRead(ctx, d, meta)
}

func resourceMessageTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	deleteHttpUrl := "v2/{project_id}/events/notification/templates"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeaders(cfg, d),
		JSONBody: map[string]interface{}{
			"names": []string{d.Id()},
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "AOM.08023006"),
			"error deleting message template")
	}

	return nil
}

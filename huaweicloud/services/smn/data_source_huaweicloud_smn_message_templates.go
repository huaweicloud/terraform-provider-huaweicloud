// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product SMN
// ---------------------------------------------------------------

package smn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SMN GET /v2/{project_id}/notifications/message_template
func DataSourceSmnMessageTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSmnMessageTemplateRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the message template.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the protocol of the message template.`,
			},
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the message template ID.`,
			},
			"templates": {
				Type:        schema.TypeList,
				Elem:        smnMessageTemplateMessageTemplateSchema(),
				Computed:    true,
				Description: `The list of message templates.`,
			},
		},
	}
}

func smnMessageTemplateMessageTemplateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the message template ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the message template name.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the protocol supported by the template.`,
			},
			"tag_names": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Indicates the variable list.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the create time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the update time.`,
			},
		},
	}
	return &sc
}

func dataSourceSmnMessageTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getMessageTemplate: Query SMN message template
	var (
		getMessageTemplateHttpUrl = "v2/{project_id}/notifications/message_template"
		getMessageTemplateProduct = "smn"
	)
	getMessageTemplateClient, err := cfg.NewServiceClient(getMessageTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	getMessageTemplatePath := getMessageTemplateClient.Endpoint + getMessageTemplateHttpUrl
	getMessageTemplatePath = strings.ReplaceAll(getMessageTemplatePath, "{project_id}", getMessageTemplateClient.ProjectID)

	getMessageTemplatequeryParams := buildGetMessageTemplateQueryParams(d)
	getMessageTemplatePath += getMessageTemplatequeryParams

	getMessageTemplateResp, err := pagination.ListAllItems(
		getMessageTemplateClient,
		"offset",
		getMessageTemplatePath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SmnMessageTemplate")
	}

	getMessageTemplateRespJson, err := json.Marshal(getMessageTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getMessageTemplateRespBody interface{}
	err = json.Unmarshal(getMessageTemplateRespJson, &getMessageTemplateRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("templates", filterGetMessageTemplateResponseBodyMessageTemplate(
			flattenGetMessageTemplateResponseBodyMessageTemplate(getMessageTemplateRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetMessageTemplateResponseBodyMessageTemplate(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("message_templates", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("message_template_id", v, nil),
			"name":       utils.PathSearch("message_template_name", v, nil),
			"protocol":   utils.PathSearch("protocol", v, nil),
			"tag_names":  utils.PathSearch("tag_names", v, nil),
			"created_at": utils.PathSearch("create_time", v, nil),
			"updated_at": utils.PathSearch("update_time", v, nil),
		})
	}
	return rst
}

func filterGetMessageTemplateResponseBodyMessageTemplate(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("template_id"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("id", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildGetMessageTemplateQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&message_template_name=%v", res, v)
	}

	if v, ok := d.GetOk("protocol"); ok {
		res = fmt.Sprintf("%s&protocol=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

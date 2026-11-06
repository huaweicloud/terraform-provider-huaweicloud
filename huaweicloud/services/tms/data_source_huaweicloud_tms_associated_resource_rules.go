package tms

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TMS GET /v2/tms/associated-resources/rules
func DataSourceTmsAssociatedResourceRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmsAssociatedResourceRulesRead,

		Schema: map[string]*schema.Schema{
			"setting_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the rule setting used to filter rules.",
			},
			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the region ID used to filter rules.",
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        tmsAssociatedResourceRuleSchema(),
				Description: "The list of associated resource rules.",
			},
		},
	}
}

func tmsAssociatedResourceRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"setting_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"existing_resource_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_delete_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTmsAssociatedResourceRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("tms", region)
	if err != nil {
		return diag.Errorf("error creating TMS client: %s", err)
	}

	listBasePath := client.Endpoint + "v2/tms/associated-resources/rules"
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	rules := make([]interface{}, 0)
	marker := ""
	for {
		listPath := listBasePath + buildGetAssociatedResourceRulesQueryParams(d, marker)
		listResp, err := client.Request("GET", listPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving TMS associated resource rules: %s", err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		rules = append(rules, flattenTmsAssociatedResourceRules(listRespBody)...)

		marker = utils.PathSearch("page_info.next_marker", listRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("rules", rules),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetAssociatedResourceRulesQueryParams(d *schema.ResourceData, marker string) string {
	path := "?limit=100"
	if settingName, ok := d.GetOk("setting_name"); ok {
		path = fmt.Sprintf("%s&setting_name=%s", path, url.QueryEscape(settingName.(string)))
	}
	if regionId, ok := d.GetOk("region_id"); ok {
		path = fmt.Sprintf("%s&region_id=%s", path, url.QueryEscape(regionId.(string)))
	}
	if marker != "" {
		path = fmt.Sprintf("%s&marker=%s", path, url.QueryEscape(marker))
	}
	return path
}

func flattenTmsAssociatedResourceRules(respBody interface{}) []interface{} {
	curJson := utils.PathSearch("rules", respBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"setting_name":             utils.PathSearch("setting_name", v, nil),
			"region_id":                utils.PathSearch("region_id", v, nil),
			"tag_keys":                 utils.PathSearch("tag_keys", v, nil),
			"existing_resource_status": utils.PathSearch("existing_resource_status", v, "").(string),
			"auto_delete_status":       utils.PathSearch("auto_delete_status", v, "").(string),
		})
	}

	return res
}

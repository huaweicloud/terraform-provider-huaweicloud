package tms

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TMS GET /v2/tms/associated-resources/settings
func DataSourceTmsAssociatedResourceSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmsAssociatedResourceSettingsRead,

		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the region ID used to filter settings.",
			},
			"settings": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        tmsAssociatedResourceSettingSchema(),
				Description: "The list of associated resource settings.",
			},
		},
	}
}

func tmsAssociatedResourceSettingSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"setting_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_service": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"associated_service": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"associated_resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"support_existing_resource": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"support_auto_delete": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"region_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceTmsAssociatedResourceSettingsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("tms", region)
	if err != nil {
		return diag.Errorf("error creating TMS client: %s", err)
	}

	listBasePath := client.Endpoint + "v2/tms/associated-resources/settings"
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	settings := make([]interface{}, 0)
	marker := ""
	for {
		listPath := listBasePath + buildGetAssociatedResourceSettingsQueryParams(d, marker)
		listResp, err := client.Request("GET", listPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving TMS associated resource settings: %s", err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		settings = append(settings, flattenTmsAssociatedResourceSettings(listRespBody)...)

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
		d.Set("settings", settings),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetAssociatedResourceSettingsQueryParams(d *schema.ResourceData, marker string) string {
	path := "?limit=100"
	if regionId, ok := d.GetOk("region_id"); ok {
		path = fmt.Sprintf("%s&region_id=%s", path, regionId.(string))
	}
	if marker != "" {
		path = fmt.Sprintf("%s&marker=%s", path, marker)
	}
	return path
}

func flattenTmsAssociatedResourceSettings(respBody interface{}) []interface{} {
	curJson := utils.PathSearch("settings", respBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"setting_name":              utils.PathSearch("setting_name", v, nil),
			"master_service":            utils.PathSearch("master_service", v, nil),
			"master_resource_type":      utils.PathSearch("master_resource_type", v, nil),
			"associated_service":        utils.PathSearch("associated_service", v, nil),
			"associated_resource_type":  utils.PathSearch("associated_resource_type", v, nil),
			"support_existing_resource": utils.PathSearch("support_existing_resource", v, false),
			"support_auto_delete":       utils.PathSearch("support_auto_delete", v, false),
			"region_ids":                utils.PathSearch("region_ids", v, nil),
		})
	}

	return res
}

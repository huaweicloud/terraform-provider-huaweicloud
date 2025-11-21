package rgc

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

// @API RGC GET /v1/governance/enabled-controls
func DataSourceEnabledControls() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnabledControlsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enabled_controls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     enabledControlsSchema(),
			},
		},
	}
}

func enabledControlsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"manage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"control_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"control_objective": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"behavior": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"regional_preference": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourceEnabledControlsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	listEnabledControlsHttpUrl := "v1/governance/enabled-controls"
	listEnabledControlsProduct := "rgc"
	listEnabledControlsClient, err := cfg.NewServiceClient(listEnabledControlsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	listEnabledControlsPath := listEnabledControlsClient.Endpoint + listEnabledControlsHttpUrl
	listEnabledControlsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var enabledControls []interface{}
	var marker string
	var queryPath string

	for {
		queryPath = listEnabledControlsPath + buildListEnabledControlsQueryParams(marker)
		listEnabledControlsResp, err := listEnabledControlsClient.Request("GET", queryPath, &listEnabledControlsOpt)
		if err != nil {
			return diag.Errorf("error retrieving RGC enabled controls: %s", err)
		}

		listEnabledControlsRespBody, err := utils.FlattenResponse(listEnabledControlsResp)
		if err != nil {
			return diag.FromErr(err)
		}

		onePageOrganizationalUnitEnabledControls := FlattenEnabledControlsResp(listEnabledControlsRespBody)
		enabledControls = append(enabledControls, onePageOrganizationalUnitEnabledControls...)
		marker = utils.PathSearch("page_info.next_marker", listEnabledControlsRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("enabled_controls", enabledControls),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListEnabledControlsQueryParams(marker string) string {
	// the default value of limit is 200
	res := "?limit=200"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func FlattenEnabledControlsResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("enabled_controls", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"manage_account_id":   utils.PathSearch("manage_account_id", v, nil),
			"control_identifier":  utils.PathSearch("control_identifier", v, nil),
			"name":                utils.PathSearch("name", v, nil),
			"description":         utils.PathSearch("description", v, nil),
			"control_objective":   utils.PathSearch("control_objective", v, nil),
			"behavior":            utils.PathSearch("behavior", v, nil),
			"owner":               utils.PathSearch("owner", v, nil),
			"regional_preference": utils.PathSearch("regional_preference", v, nil),
		})
	}
	return rst
}

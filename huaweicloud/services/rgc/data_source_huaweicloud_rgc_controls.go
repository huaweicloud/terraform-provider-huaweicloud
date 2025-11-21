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

// @API RGC GET /v1/governance/controls
func DataSourceControls() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceControlsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"controls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     ControlsSchema(),
			},
		},
	}
}

func ControlsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"identifier": {
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
			"guidance": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"framework": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"service": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"implementation": {
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
			"severity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"control_objective": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"release_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourceControlsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	listControlsHttpUrl := "v1/governance/controls"
	listControlsProduct := "rgc"
	listControlsClient, err := cfg.NewServiceClient(listControlsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	listControlsPath := listControlsClient.Endpoint + listControlsHttpUrl
	listControlsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var controls []interface{}
	var marker string
	var queryPath string

	for {
		queryPath = listControlsPath + buildListControlsQueryParams(marker)
		listControlsResp, err := listControlsClient.Request("GET", queryPath, &listControlsOpt)
		if err != nil {
			return diag.Errorf("error retrieving RGC controls: %s", err)
		}

		listControlsRespBody, err := utils.FlattenResponse(listControlsResp)
		if err != nil {
			return diag.FromErr(err)
		}

		onePageOrganizationalUnitControls := FlattenControlsResp(listControlsRespBody)
		controls = append(controls, onePageOrganizationalUnitControls...)
		marker = utils.PathSearch("page_info.next_marker", listControlsRespBody, "").(string)
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
		d.Set("controls", controls),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListControlsQueryParams(marker string) string {
	// the default value of limit is 200
	res := "?limit=200"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func FlattenControlsResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("controls", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"identifier":        utils.PathSearch("identifier", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"description":       utils.PathSearch("description", v, nil),
			"guidance":          utils.PathSearch("guidance", v, nil),
			"resource":          utils.PathSearch("resource", v, nil),
			"framework":         utils.PathSearch("framework", v, nil),
			"service":           utils.PathSearch("service", v, nil),
			"implementation":    utils.PathSearch("implementation", v, nil),
			"behavior":          utils.PathSearch("behavior", v, nil),
			"owner":             utils.PathSearch("owner", v, nil),
			"severity":          utils.PathSearch("severity", v, nil),
			"control_objective": utils.PathSearch("control_objective", v, nil),
			"version":           utils.PathSearch("version", v, nil),
			"release_date":      utils.PathSearch("release_date", v, nil),
		})
	}
	return rst
}

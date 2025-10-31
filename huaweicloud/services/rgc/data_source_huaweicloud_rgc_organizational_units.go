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

func DataSourceOrganizationalUnits() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationalUnitsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"control_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"managed_organizational_units": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     organizationalUnitsSchema(),
			},
		},
	}
}

func organizationalUnitsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"manage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organizational_unit_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organizational_unit_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organizational_unit_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organizational_unit_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_organizational_unit_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_organizational_unit_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"landing_zone_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourceOrganizationalUnitsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	listOrganizationalUnitsHttpUrl := "v1/managed-organization/managed-organizational-units"
	listOrganizationalUnitsProduct := "rgc"
	listOrganizationalUnitsClient, err := cfg.NewServiceClient(listOrganizationalUnitsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	listOrganizationalUnitsPath := listOrganizationalUnitsClient.Endpoint + listOrganizationalUnitsHttpUrl
	listOrganizationalUnitsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var childrenOUs []interface{}
	var marker string
	var queryPath string

	for {
		queryPath = listOrganizationalUnitsPath + buildListOrganizationalUnitsQueryParams(d, marker)
		listOrganizationalUnitsResp, err := listOrganizationalUnitsClient.Request("GET", queryPath, &listOrganizationalUnitsOpt)
		if err != nil {
			return diag.Errorf("error retrieving RGC organizational units: %s", err)
		}

		listOrganizationalUnitsRespBody, err := utils.FlattenResponse(listOrganizationalUnitsResp)
		if err != nil {
			return diag.FromErr(err)
		}

		onePageOUs := flattenOrganizationalUnitResp(listOrganizationalUnitsRespBody)
		childrenOUs = append(childrenOUs, onePageOUs...)
		marker = utils.PathSearch("page_info.next_marker", listOrganizationalUnitsRespBody, "").(string)
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
		d.Set("managed_organizational_units", childrenOUs),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOrganizationalUnitResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("managed_organizational_units", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"manage_account_id":               utils.PathSearch("manage_account_id", v, nil),
			"organizational_unit_id":          utils.PathSearch("organizational_unit_id", v, nil),
			"organizational_unit_name":        utils.PathSearch("organizational_unit_name", v, nil),
			"organizational_unit_status":      utils.PathSearch("organizational_unit_status", v, nil),
			"organizational_unit_type":        utils.PathSearch("organizational_unit_type", v, nil),
			"parent_organizational_unit_id":   utils.PathSearch("parent_organizational_unit_id", v, nil),
			"parent_organizational_unit_name": utils.PathSearch("parent_organizational_unit_name", v, nil),
			"created_at":                      utils.PathSearch("created_at", v, nil),
			"landing_zone_version":            utils.PathSearch("landing_zone_version", v, nil),
		})
	}
	return rst
}

func buildListOrganizationalUnitsQueryParams(d *schema.ResourceData, marker string) string {
	// the default value of limit is 200
	res := "?limit=200"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	if v, ok := d.GetOk("control_id"); ok {
		res = fmt.Sprintf("%s&control_id=%v", res, v)
	}

	return res
}

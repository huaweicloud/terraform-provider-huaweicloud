// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Organizations
// ---------------------------------------------------------------

package organizations

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Organizations GET /v1/organizations/organizational-units
func DataSourceOrganizationalUnits() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationalUnitsRead,
		Schema: map[string]*schema.Schema{
			"parent_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of root or organizational unit.`,
			},
			"children": {
				Type:        schema.TypeList,
				Elem:        organizationalUnitsOrganizationalUnitSchema(),
				Computed:    true,
				Description: `List of OUs in an organization.`,
			},
		},
	}
}

func organizationalUnitsOrganizationalUnitSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Unique ID of an OU`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Name of the OU`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Uniform resource name of the OU`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Time when the OU was created`,
			},
		},
	}
	return &sc
}

func dataSourceOrganizationalUnitsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	listOrganizationalUnitsHttpUrl := "v1/organizations/organizational-units"
	listOrganizationalUnitsProduct := "organizations"
	listOrganizationalUnitsClient, err := cfg.NewServiceClient(listOrganizationalUnitsProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
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
			return common.CheckDeletedDiag(d, err, "error retrieving Organizational Units")
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

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("children", childrenOUs),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOrganizationalUnitResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("organizational_units", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"name":       utils.PathSearch("name", v, nil),
			"urn":        utils.PathSearch("urn", v, nil),
			"created_at": utils.PathSearch("created_at", v, nil),
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

	if v, ok := d.GetOk("parent_id"); ok {
		res = fmt.Sprintf("%s&parent_id=%v", res, v)
	}

	return res
}

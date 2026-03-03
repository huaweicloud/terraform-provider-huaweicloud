package organizations

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

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
				Description: `The ID of root or organizational unit.`,
			},
			"children": {
				Type:        schema.TypeList,
				Elem:        organizationalUnitsSchema(),
				Computed:    true,
				Description: `The list of child organizational units.`,
			},
		},
	}
}

func organizationalUnitsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the organizational unit.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the organizational unit.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the organizational unit.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the organizational unit was created.`,
			},
		},
	}
	return &sc
}

func listOrganizationalUnits(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/organizations/organizational-units"
		// The default value of limit is 200
		limit  = 200
		marker = ""
		result []interface{}
		opt    = golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
	)

	listPath := client.Endpoint + httpUrl
	listPath = fmt.Sprintf("%s?limit=%v%s", listPath, limit, buildListOrganizationalUnitsQueryParams(d))

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%v", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		ous := utils.PathSearch("organizational_units", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, ous...)
		if len(ous) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func dataSourceOrganizationalUnitsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("organizations", region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	ous, err := listOrganizationalUnits(client, d)
	if err != nil {
		return diag.Errorf("error querying organizational units: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuid)

	return diag.FromErr(d.Set("children", flattenOrganizationalUnits(ous)))
}

func flattenOrganizationalUnits(ous []interface{}) []interface{} {
	if len(ous) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(ous))
	for _, v := range ous {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"name":       utils.PathSearch("name", v, nil),
			"urn":        utils.PathSearch("urn", v, nil),
			"created_at": utils.PathSearch("created_at", v, nil),
		})
	}
	return rst
}

func buildListOrganizationalUnitsQueryParams(d *schema.ResourceData) string {
	if v, ok := d.GetOk("parent_id"); ok {
		return fmt.Sprintf("&parent_id=%v", v)
	}

	return ""
}

package cfw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CFW GET /v1/{project_id}/system/multi-account/organization-tree

func DataSourceCfwOrganizationTree() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCfwOrganizationTreeRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"fw_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the firewall instance ID.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the parent organization unit ID.`,
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of organization tree nodes.`,
				Elem:        organizationTreeNodeSchema(),
			},
		},
	}
}

func organizationTreeNodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"delegated": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `The indication of whether the organization unit is delegated.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The organization unit ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The organization unit name.`,
			},
			"org_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The organization unit type.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The parent organization unit ID.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the organization unit.`,
			},
		},
	}
}

func buildOrganizationTreeQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?fw_instance_id=%s&limit=2000", d.Get("fw_instance_id").(string))

	if v, ok := d.GetOk("parent_id"); ok {
		queryParams = fmt.Sprintf("%s&parent_id=%s", queryParams, v.(string))
	}

	return queryParams
}

func dataSourceCfwOrganizationTreeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/system/multi-account/organization-tree"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	listPath += buildOrganizationTreeQueryParams(d)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", listPath, &reqOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW organization tree: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenOrganizationTreeData(utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOrganizationTreeData(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dataResp))
	for _, node := range dataResp {
		result = append(result, map[string]interface{}{
			"delegated": utils.PathSearch("delegated", node, false),
			"id":        utils.PathSearch("id", node, nil),
			"name":      utils.PathSearch("name", node, nil),
			"org_type":  utils.PathSearch("org_type", node, nil),
			"parent_id": utils.PathSearch("parent_id", node, nil),
			"urn":       utils.PathSearch("urn", node, nil),
		})
	}

	return result
}

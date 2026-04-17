package vpc

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

// @API VPC GET /v3/{project_id}/vpc/list-address-groups-dependency
func DataSourceVpcAddressGroupAssociatedResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAddressGroupAssociatedResourcesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resources.",
			},
			"ip_address_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of an IP address group that will be used as a filter.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the enterprise project that an IP address group belongs to.",
			},
			"address_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        addressGroupDependencySchema(),
				Description: "The response body for querying the resources associated with an address group.",
			},
		},
	}
}

func addressGroupDependencySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of an IP address group.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the enterprise project that an IP address group belongs to.",
			},
			"dependency": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dependencySchema(),
				Description: "The list of resources associated with the vpc address group.",
			},
		},
	}
}

func dependencySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the resource associated with the IP address group.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the resource associated with the IP address group.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource associated with the IP address group.",
			},
		},
	}
}

func dataSourceAddressGroupAssociatedResourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/vpc/list-address-groups-dependency"
	)

	client, err := cfg.NewServiceClient("vpc", region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildGetAddressAssociatedResourcesQueryParams(d)

	reqOpt := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, reqOpt)
	if err != nil {
		return diag.Errorf("error retrieving address group associated resources: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("address_groups", flattenListAddressAssociatedResourcesBody(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetAddressAssociatedResourcesQueryParams(d *schema.ResourceData) string {
	params := ""
	if v, ok := d.GetOk("ip_address_group_id"); ok {
		params += fmt.Sprintf("&id=%v", v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		params += fmt.Sprintf("&enterprise_project_id=%s", v.(string))
	}
	if params != "" {
		return "?" + params[1:] // Remove the leading '&'
	}
	return ""
}

func flattenListAddressAssociatedResourcesBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("address_groups", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"dependency":            flattenDependencyBody(v),
		})
	}
	return res
}

func flattenDependencyBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("dependency", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"type":          utils.PathSearch("type", v, nil),
			"instance_id":   utils.PathSearch("instance_id", v, nil),
			"instance_name": utils.PathSearch("instance_name", v, nil),
		})
	}
	return res
}

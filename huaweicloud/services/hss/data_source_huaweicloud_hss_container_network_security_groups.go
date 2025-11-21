package hss

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

// @API HSS GET /v5/{project_id}/container-network/security-groups
func DataSourceContainerNetworkSecurityGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerNetworkSecurityGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildContainerNetworkSecurityGroupsQueryParams(epsId string) string {
	queryParams := ""

	if epsId != "" {
		queryParams = fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return queryParams
}

func dataSourceContainerNetworkSecurityGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/container-network/security-groups"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildContainerNetworkSecurityGroupsQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS container network security groups: %s", err)
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
		d.Set("security_groups", flattenContainerNetworkSecurityGroups(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenContainerNetworkSecurityGroups(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson, ok := resp.([]interface{})
	if !ok || len(curJson) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curJson))
	for _, v := range curJson {
		rst = append(rst, map[string]interface{}{
			"security_group_id":          utils.PathSearch("security_group_id", v, nil),
			"security_group_name":        utils.PathSearch("security_group_name", v, nil),
			"security_group_description": utils.PathSearch("security_group_description", v, nil),
		})
	}

	return rst
}

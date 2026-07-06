package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v1/{project_id}/multi-accounts/organizational-unit-tree
func DataSourceMultiOrganizationTree() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMultiOrganizationTreeRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"entity_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ou_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     multiOrganizationTreeOUSchema(),
			},
		},
	}
}

func multiOrganizationTreeOUSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildMultiOrganizationTreeQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?entity_id=%v", d.Get("entity_id"))
}

func dataSourceMultiOrganizationTreeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/multi-accounts/organizational-unit-tree"
		product = "dsc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildMultiOrganizationTreeQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC multi organization tree: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("ou_list", flattenMultiOrganizationTreeOUList(
			utils.PathSearch("ou_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMultiOrganizationTreeOUList(ouList []interface{}) []interface{} {
	if len(ouList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(ouList))
	for _, v := range ouList {
		rst = append(rst, utils.RemoveNil(map[string]interface{}{
			"id":   utils.PathSearch("id", v, nil),
			"name": utils.PathSearch("name", v, nil),
		}))
	}

	return rst
}

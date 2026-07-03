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

// @API DSC GET /v1/{project_id}/multi-accounts/organizational-unit-list
func DataSourceMultiOrganizations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMultiOrganizationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ou_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     multiOrganizationsOUSchema(),
			},
		},
	}
}

func multiOrganizationsOUSchema() *schema.Resource {
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

func buildMultiOrganizationsQueryParams(d *schema.ResourceData, marker string) string {
	rst := "?limit=200"
	if v, ok := d.GetOk("account_id"); ok {
		rst += fmt.Sprintf("&account_id=%v", v)
	}

	if v, ok := d.GetOk("parent_id"); ok {
		rst += fmt.Sprintf("&parent_id=%v", v)
	}

	if marker != "" {
		rst += fmt.Sprintf("&marker=%s", marker)
	}

	return rst
}

func dataSourceMultiOrganizationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v1/{project_id}/multi-accounts/organizational-unit-list"
		product   = "dsc"
		marker    = ""
		allOUList = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := requestPath + buildMultiOrganizationsQueryParams(d, marker)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC multi organizations: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		ouListResp := utils.PathSearch("ou_list", respBody, make([]interface{}, 0)).([]interface{})
		allOUList = append(allOUList, ouListResp...)
		marker = utils.PathSearch("page_info.marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("ou_list", flattenMultiOrganizationsOUList(allOUList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMultiOrganizationsOUList(ouList []interface{}) []interface{} {
	if len(ouList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(ouList))
	for _, v := range ouList {
		rst = append(rst, map[string]interface{}{
			"id":   utils.PathSearch("id", v, nil),
			"name": utils.PathSearch("name", v, nil),
		})
	}

	return rst
}

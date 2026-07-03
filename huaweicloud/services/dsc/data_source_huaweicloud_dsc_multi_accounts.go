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

// @API DSC GET /v1/{project_id}/multi-accounts/list
func DataSourceMultiAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMultiAccountsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     multiAccountsSchema(),
			},
		},
	}
}

func multiAccountsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildMultiAccountsQueryParams(marker string) string {
	rst := "?limit=200"
	if marker != "" {
		rst += fmt.Sprintf("&marker=%s", marker)
	}

	return rst
}

func dataSourceMultiAccountsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/multi-accounts/list"
		product     = "dsc"
		marker      = ""
		allAccounts = make([]interface{}, 0)
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
		currentPath := requestPath + buildMultiAccountsQueryParams(marker)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC multi accounts: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		accountsResp := utils.PathSearch("accounts", respBody, make([]interface{}, 0)).([]interface{})
		allAccounts = append(allAccounts, accountsResp...)
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
		d.Set("accounts", flattenMultiAccounts(allAccounts)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMultiAccounts(accounts []interface{}) []interface{} {
	if len(accounts) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(accounts))
	for _, v := range accounts {
		rst = append(rst, utils.RemoveNil(map[string]interface{}{
			"domain_id":  utils.PathSearch("domain_id", v, nil),
			"name":       utils.PathSearch("name", v, nil),
			"project_id": utils.PathSearch("project_id", v, nil),
		}))
	}

	return rst
}

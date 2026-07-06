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

// @API DSC GET /v1/{project_id}/multi-accounts/account-with-assets-list
func DataSourceMultiAccountAssets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMultiAccountAssetsRead,

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
			"accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     multiAccountAssetsSchema(),
			},
		},
	}
}

func multiAccountAssetsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bigdata_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"db_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"obs_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildMultiAccountAssetsQueryParams(d *schema.ResourceData, marker string) string {
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

func dataSourceMultiAccountAssetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/multi-accounts/account-with-assets-list"
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
		currentPath := requestPath + buildMultiAccountAssetsQueryParams(d, marker)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC multi account assets: %s", err)
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
		d.Set("accounts", flattenMultiAccountAssets(allAccounts)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMultiAccountAssets(accounts []interface{}) []interface{} {
	if len(accounts) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(accounts))
	for _, v := range accounts {
		rst = append(rst, map[string]interface{}{
			"bigdata_count": utils.PathSearch("bigdata_count", v, nil),
			"db_count":      utils.PathSearch("db_count", v, nil),
			"domain_id":     utils.PathSearch("domain_id", v, nil),
			"name":          utils.PathSearch("name", v, nil),
			"obs_count":     utils.PathSearch("obs_count", v, nil),
			"project_id":    utils.PathSearch("project_id", v, nil),
		})
	}

	return rst
}

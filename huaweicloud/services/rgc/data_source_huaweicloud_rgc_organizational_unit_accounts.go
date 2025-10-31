package rgc

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

// @API RGC GET /v1/managed-organization/managed-organizational-units/{managed_organizational_unit_id}/managed-accounts
func DataSourceOrganizationalUnitAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationalUnitAccountsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"managed_organizational_unit_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"managed_accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     accountSchema(),
			},
		},
	}
}

func accountSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"landing_zone_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"manage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message": {
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
			"identity_store_user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"blueprint_product_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"blueprint_product_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"blueprint_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_blueprint_has_multi_account_resource": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourceOrganizationalUnitAccountsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ouId := d.Get("managed_organizational_unit_id").(string)
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	listOrganizationalUnitAccountsHttpUrl := "v1/managed-organization/managed-organizational-units/{managed_organizational_unit_id}/managed-accounts"
	listOrganizationalUnitAccountsProduct := "rgc"
	listOrganizationalUnitAccountsClient, err := cfg.NewServiceClient(listOrganizationalUnitAccountsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	listOrganizationalUnitAccountsPath := listOrganizationalUnitAccountsClient.Endpoint + listOrganizationalUnitAccountsHttpUrl
	listOrganizationalUnitAccountsPath = strings.ReplaceAll(listOrganizationalUnitAccountsPath, "{managed_organizational_unit_id}", ouId)
	listOrganizationalUnitAccountsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var organizationalUnitAccounts []interface{}
	var marker string
	var queryPath string

	for {
		queryPath = listOrganizationalUnitAccountsPath + buildListOrganizationalUnitAccountQueryParams(marker)
		listOrganizationalUnitAccountsResp, err := listOrganizationalUnitAccountsClient.Request("GET", queryPath, &listOrganizationalUnitAccountsOpt)
		if err != nil {
			return diag.Errorf("error retrieving RGC accounts: %s", err)
		}

		listOrganizationalUnitAccountsRespBody, err := utils.FlattenResponse(listOrganizationalUnitAccountsResp)
		if err != nil {
			return diag.FromErr(err)
		}

		onePageOrganizationalUnitAccounts := flattenAccountsResp(listOrganizationalUnitAccountsRespBody)
		organizationalUnitAccounts = append(organizationalUnitAccounts, onePageOrganizationalUnitAccounts...)
		marker = utils.PathSearch("page_info.next_marker", listOrganizationalUnitAccountsRespBody, "").(string)
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
		d.Set("managed_accounts", organizationalUnitAccounts),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListOrganizationalUnitAccountQueryParams(marker string) string {
	// the default value of limit is 200
	res := "?limit=200"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func flattenAccountsResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("managed_accounts", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"landing_zone_version":                    utils.PathSearch("landing_zone_version", v, nil),
			"manage_account_id":                       utils.PathSearch("manage_account_id", v, nil),
			"account_id":                              utils.PathSearch("account_id", v, nil),
			"account_name":                            utils.PathSearch("account_name", v, nil),
			"account_type":                            utils.PathSearch("account_type", v, nil),
			"owner":                                   utils.PathSearch("owner", v, nil),
			"state":                                   utils.PathSearch("state", v, nil),
			"message":                                 utils.PathSearch("message", v, nil),
			"parent_organizational_unit_id":           utils.PathSearch("parent_organizational_unit_id", v, nil),
			"parent_organizational_unit_name":         utils.PathSearch("parent_organizational_unit_name", v, nil),
			"identity_store_user_name":                utils.PathSearch("identity_store_user_name", v, nil),
			"blueprint_product_id":                    utils.PathSearch("blue_product_id", v, nil),
			"blueprint_product_version":               utils.PathSearch("blue_product_version", v, nil),
			"blueprint_status":                        utils.PathSearch("blueprint_status", v, nil),
			"is_blueprint_has_multi_account_resource": utils.PathSearch("is_blueprint_has_multi_account_resource", v, nil),
			"regions":                                 utils.PathSearch("regions", v, nil),
			"created_at":                              utils.PathSearch("created_at", v, nil),
			"updated_at":                              utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}

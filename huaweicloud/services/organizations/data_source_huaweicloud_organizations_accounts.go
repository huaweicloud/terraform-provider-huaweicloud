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

// @API Organizations GET /v1/organizations/accounts
func DataSourceAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountsRead,
		Schema: map[string]*schema.Schema{
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of root or organizational unit.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the account.`,
			},
			"with_register_contact_info": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to return email addresses and mobile numbers associated with the account.",
			},
			"accounts": {
				Type:        schema.TypeList,
				Elem:        organizationsAccountSchema(),
				Computed:    true,
				Description: `List of accounts in an organization.`,
			},
		},
	}
}

func organizationsAccountSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Unique ID of an account.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Name of the account.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Uniform resource name of the account.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Description of the account.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Status of the account.`,
			},
			"join_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `How an account joined an organization.`,
			},
			"joined_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Time when an account joined an organization.`,
			},
			"mobile_phone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Mobile phone number.`,
			},
			"intl_number_prefix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Prefix of a mobile phone number.`,
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Email address associated with the account.`,
			},
		},
	}
	return &sc
}

func dataSourceAccountsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	listAccountsHttpUrl := "v1/organizations/accounts"
	listAccountsProduct := "organizations"
	listAccountsClient, err := cfg.NewServiceClient(listAccountsProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	listAccountsPath := listAccountsClient.Endpoint + listAccountsHttpUrl
	listAccountsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	filterName := d.Get("name").(string)
	var orgAccounts []interface{}
	var marker string
	var queryPath string

	for {
		queryPath = listAccountsPath + buildListAccountsQueryParams(d, marker)
		listAccountsResp, err := listAccountsClient.Request("GET", queryPath, &listAccountsOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving Organizations accounts")
		}

		listAccountsRespBody, err := utils.FlattenResponse(listAccountsResp)
		if err != nil {
			return diag.FromErr(err)
		}

		onePageAccounts := flattenAccountsResp(listAccountsRespBody, filterName)
		orgAccounts = append(orgAccounts, onePageAccounts...)
		marker = utils.PathSearch("page_info.next_marker", listAccountsRespBody, "").(string)
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
		d.Set("accounts", orgAccounts),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAccountsResp(resp interface{}, name string) []interface{} {
	if resp == nil {
		return nil
	}

	jsonPath := "accounts"
	if name != "" {
		jsonPath += fmt.Sprintf("[?name=='%s']", name)
	}
	curJson := utils.PathSearch(jsonPath, resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                 utils.PathSearch("id", v, nil),
			"name":               utils.PathSearch("name", v, nil),
			"urn":                utils.PathSearch("urn", v, nil),
			"description":        utils.PathSearch("description", v, nil),
			"status":             utils.PathSearch("status", v, nil),
			"join_method":        utils.PathSearch("join_method", v, nil),
			"joined_at":          utils.PathSearch("joined_at", v, nil),
			"mobile_phone":       utils.PathSearch("mobile_phone", v, nil),
			"intl_number_prefix": utils.PathSearch("intl_number_prefix", v, nil),
			"email":              utils.PathSearch("email", v, nil),
		})
	}
	return rst
}

func buildListAccountsQueryParams(d *schema.ResourceData, marker string) string {
	// the default value of limit is 200
	res := "?limit=200"

	res = fmt.Sprintf("%s&with_register_contact_info=%v", res, d.Get("with_register_contact_info"))

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	if v, ok := d.GetOk("parent_id"); ok {
		res = fmt.Sprintf("%s&parent_id=%v", res, v)
	}

	return res
}

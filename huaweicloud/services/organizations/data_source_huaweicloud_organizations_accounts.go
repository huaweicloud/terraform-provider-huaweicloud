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

// @API Organizations GET /v1/organizations/accounts
func DataSourceAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountsRead,
		Schema: map[string]*schema.Schema{
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of root or organizational unit.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the account.`,
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
				Description: `The list of accounts that match the filter parameters.`,
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
				Description: `The unique ID of an account.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the account.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the account.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the account.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the account.`,
			},
			"join_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `How the account joined the organization.`,
			},
			"joined_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the account joined the organization.`,
			},
			"mobile_phone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The mobile phone number.`,
			},
			"intl_number_prefix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The prefix of a mobile phone number.`,
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The email address associated with the account.`,
			},
		},
	}
	return &sc
}

func listAccounts(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/organizations/accounts"
		marker  = ""
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listPath += buildListAccountsQueryParams(d)
	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker += fmt.Sprintf("&marker=%v", marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		accounts := utils.PathSearch("accounts", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, accounts...)
		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func dataSourceAccountsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	accounts, err := listAccounts(client, d)
	if err != nil {
		return diag.Errorf("error retrieving accounts: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return diag.FromErr(d.Set("accounts", flattenAccounts(accounts, d.Get("name").(string))))
}

func flattenAccounts(accounts []interface{}, name string) []interface{} {
	if len(accounts) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(accounts))
	for _, v := range accounts {
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

	if name != "" {
		return utils.PathSearch(fmt.Sprintf("[?name=='%s']", name), rst, make([]interface{}, 0)).([]interface{})
	}

	return rst
}

func buildListAccountsQueryParams(d *schema.ResourceData) string {
	// the default value of limit is 200
	res := fmt.Sprintf("?limit=200&with_register_contact_info=%v", d.Get("with_register_contact_info"))

	if v, ok := d.GetOk("parent_id"); ok {
		res = fmt.Sprintf("%s&parent_id=%v", res, v)
	}

	return res
}

package organizations

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Organizations GET /v1/organizations/accounts/{account_id}/delegated-services
func DataSourceDelegatedServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDelegatedServicesRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The unique ID of an account.`,
			},
			"delegated_services": {
				Type:        schema.TypeList,
				Elem:        delegatedServiceSchema(),
				Computed:    true,
				Description: `The list of the delegated services.`,
			},
		},
	}
}

func delegatedServiceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"service_principal": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the service principal.`,
			},
			"delegation_enabled_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The date when the account became a delegated administrator for the service.`,
			},
		},
	}
	return &sc
}

func listDelegatedServices(client *golangsdk.ServiceClient, accountId string) ([]interface{}, error) {
	var (
		httpUrl = "v1/organizations/accounts/{account_id}/delegated-services"
		marker  = ""
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{account_id}", accountId)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s?marker=%s", listPathWithMarker, marker)
		}

		listResp, err := client.Request("GET", listPathWithMarker, &listOpt)
		if err != nil {
			return nil, err
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return nil, err
		}

		delegatedServices := utils.PathSearch("delegated_services", listRespBody, make([]interface{}, 0)).([]interface{})
		result = append(result, delegatedServices...)
		marker = utils.PathSearch("page_info.next_marker", listRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func dataSourceDelegatedServicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	accountId := d.Get("account_id").(string)
	delegatedServices, err := listDelegatedServices(client, accountId)
	if err != nil {
		return diag.Errorf("error retrieving Organizations delegated services for account (%s): %s", accountId, err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return diag.FromErr(d.Set("delegated_services", flattenDelegatedServiceResp(delegatedServices)))
}

func flattenDelegatedServiceResp(delegatedServices []interface{}) []interface{} {
	if len(delegatedServices) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(delegatedServices))
	for _, v := range delegatedServices {
		rst = append(rst, map[string]interface{}{
			"service_principal":     utils.PathSearch("service_principal", v, nil),
			"delegation_enabled_at": utils.PathSearch("delegation_enabled_at", v, nil),
		})
	}

	return rst
}

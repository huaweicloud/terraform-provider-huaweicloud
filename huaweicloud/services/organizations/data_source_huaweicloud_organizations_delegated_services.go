package organizations

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Organizations GET /v1/organizations/accounts/{account_id}/delegated-services
func DataSourceOrganizationsDelegatedServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDelegatedServices,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delegated_services": {
				Type:        schema.TypeList,
				Elem:        organizationsDelegatedServiceSchema(),
				Computed:    true,
				Description: `List of the services for which the specified account is a delegated administrator.`,
			},
		},
	}
}

func organizationsDelegatedServiceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"service_principal": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Name of the service principal.`,
			},
			"delegation_enabled_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Date when the account became a delegated administrator for the service.`,
			},
		},
	}
	return &sc
}

func dataSourceDelegatedServices(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	listDelegatedServiceHttpUrl := "v1/organizations/accounts/{account_id}/delegated-services"
	listDelegatedServiceProduct := "organizations"
	listDelegatedServiceClient, err := cfg.NewServiceClient(listDelegatedServiceProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	listDelegatedServicePath := listDelegatedServiceClient.Endpoint + listDelegatedServiceHttpUrl
	listDelegatedServicePath = strings.ReplaceAll(listDelegatedServicePath, "{account_id}", d.Get("account_id").(string))
	listDelegatedServiceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var orgDelegatedServices []interface{}
	var marker string
	var queryPath string

	for {
		queryPath = listDelegatedServicePath + buildListDelegatedServiceQueryParams(marker)
		listDelegatedServiceResp, err := listDelegatedServiceClient.Request("GET", queryPath, &listDelegatedServiceOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving Organizations Delegated Service")
		}

		listDelegatedServiceRespBody, err := utils.FlattenResponse(listDelegatedServiceResp)
		if err != nil {
			return diag.FromErr(err)
		}

		onePageDelegatedServices := flattenDelegatedServiceResp(listDelegatedServiceRespBody)
		orgDelegatedServices = append(orgDelegatedServices, onePageDelegatedServices...)
		marker = utils.PathSearch("page_info.next_marker", listDelegatedServiceRespBody, "").(string)
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
		d.Set("delegated_services", orgDelegatedServices),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDelegatedServiceResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	jsonPath := "delegated_services"

	curJson := utils.PathSearch(jsonPath, resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"service_principal":     utils.PathSearch("service_principal", v, nil),
			"delegation_enabled_at": utils.PathSearch("delegation_enabled_at", v, nil),
		})
	}
	return rst
}

func buildListDelegatedServiceQueryParams(marker string) string {
	// the default value of limit is 200
	res := "?limit=200"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

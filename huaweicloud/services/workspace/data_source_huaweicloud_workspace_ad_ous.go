package workspace

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v2/{project_id}/ad-ous
func DataSourceAdOus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAdOusRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the OUs are located.`,
			},
			"ous": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the OU.`,
						},
						"ou_dn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The distinguished name (DN) of the OU.`,
						},
					},
				},
				Description: `The list of OUs.`,
			},
		},
	}
}

func dataSourceAdOusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	ous, err := listAdOus(client)
	if err != nil {
		return diag.Errorf("error querying OUs in AD service: %s", err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(randomId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("ous", flattenAdOus(ous)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listAdOus(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/ad-ous"
		listOpt = golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=UTF-8",
			},
		}
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	resp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("ou_infos", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenAdOus(ouInfos []interface{}) []interface{} {
	if len(ouInfos) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(ouInfos))
	for _, ouInfo := range ouInfos {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("ou_name", ouInfo, nil),
			"ou_dn": utils.PathSearch("ou_dn", ouInfo, nil),
		})
	}

	return result
}

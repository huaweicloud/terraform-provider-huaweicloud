package aom

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

// @API AOM GET /v1/{project_id}/aom/organization-counts
func DataSourceOrganizationAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationAccountsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"join_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"joined_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOrganizationAccountsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	listHttpUrl := "v1/{project_id}/aom/organization-counts"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving organization accounts: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.Errorf("error flattening organization accounts: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID")
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("accounts", flattenAccounts(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAccounts(resp interface{}) []interface{} {
	accounts, ok := resp.([]interface{})
	if !ok || len(accounts) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(accounts))
	for _, account := range accounts {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", account, nil),
			"name":        utils.PathSearch("name", account, nil),
			"urn":         utils.PathSearch("urn", account, nil),
			"join_method": utils.PathSearch("join_method", account, nil),
			"joined_at":   utils.PathSearch("joined_at", account, nil),
		})
	}

	return result
}

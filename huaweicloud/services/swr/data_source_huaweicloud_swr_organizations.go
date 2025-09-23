package swr

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SWR GET /v2/manage/namespaces
func DataSourceOrganizations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organizations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permission": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_user_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"repo_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOrganizationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listOrganizationHttpUrl = "v2/manage/namespaces"
		listOrganizationProduct = "swr"
	)

	listOrganizationClient, err := cfg.NewServiceClient(listOrganizationProduct, region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	listOrganizationPath := listOrganizationClient.Endpoint + listOrganizationHttpUrl

	listOrganizationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listOrganizationPath += "?filter=mode::visible"

	listOrganizationResp, err := listOrganizationClient.Request("GET", listOrganizationPath, &listOrganizationOpt)
	if err != nil {
		return diag.Errorf("error querying SWR organizations: %s", err)
	}

	listOrganizationRespBody, err := utils.FlattenResponse(listOrganizationResp)
	if err != nil {
		return diag.Errorf("error retrieving SWR organizations: %s", err)
	}

	results := make([]map[string]interface{}, 0)

	organizations := utils.PathSearch("namespaces", listOrganizationRespBody, make([]interface{}, 0)).([]interface{})
	for _, organization := range organizations {
		name := utils.PathSearch("name", organization, "").(string)
		if val, ok := d.GetOk("name"); ok && name != val {
			continue
		}

		permission := utils.PathSearch("auth", organization, float64(0)).(float64)
		results = append(results, map[string]interface{}{
			"id":                utils.PathSearch("id", organization, nil),
			"name":              utils.PathSearch("name", organization, nil),
			"creator":           utils.PathSearch("creator_name", organization, nil),
			"permission":        resourceSWRAuthToPermission(int(permission)),
			"access_user_count": int(utils.PathSearch("access_user_count", organization, float64(0)).(float64)),
			"repo_count":        int(utils.PathSearch("repo_count", organization, float64(0)).(float64)),
		})
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("organizations", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

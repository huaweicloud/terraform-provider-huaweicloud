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

// @API AOM GET /v1/{project_id}/aom/auth/grant
func DataSourceCloudServiceAuthorizations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudServiceAuthorizationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"authorizations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the authorizations list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the authorization service.`,
						},
						"role_name": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: `Indicates the role names list.`,
						},
						"status": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates the authorization status.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudServiceAuthorizationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	listHttpUrl := "v1/{project_id}/aom/auth/grant"
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
		return diag.Errorf("error retrieving cloud service authorizations: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.Errorf("error flattening cloud service authorizations: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID")
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("authorizations", flattenCloudServiceAuthorizations(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCloudServiceAuthorizations(resp interface{}) interface{} {
	if m, ok := resp.(map[string]interface{}); ok {
		rst := make([]map[string]interface{}, 0, len(m))
		for k, v := range m {
			rst = append(rst, map[string]interface{}{
				"service":   k,
				"role_name": utils.PathSearch("role_name", v, nil),
				"status":    utils.PathSearch("status", v, nil),
			})
		}
		return rst
	}

	return nil
}

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

// @API Workspace GET /v1/{project_id}/app-center/app-catalogs
func DataSourceApplicationCatalogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationCatalogsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the application catalogs are located.`,
			},
			"catalogs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        catalogSchema(),
				Description: `The list of application catalogs.`,
			},
		},
	}
}

func catalogSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the application catalog.`,
			},
			"zh": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The catalog description in Chinese.`,
			},
			"en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The catalog description in English.`,
			},
		},
	}
}

func listApplicationCatalogs(client *golangsdk.ServiceClient) ([]interface{}, error) {
	httpUrl := "v1/{project_id}/app-center/app-catalogs"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenApplicationCatalogs(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"id": utils.PathSearch("id", item, nil),
			"zh": utils.PathSearch("catalog_zh", item, nil),
			"en": utils.PathSearch("catalog_en", item, nil),
		})
	}

	return result
}

func dataSourceApplicationCatalogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	resp, err := listApplicationCatalogs(client)
	if err != nil {
		return diag.Errorf("error querying Workspace application catalogs: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("catalogs", flattenApplicationCatalogs(resp)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

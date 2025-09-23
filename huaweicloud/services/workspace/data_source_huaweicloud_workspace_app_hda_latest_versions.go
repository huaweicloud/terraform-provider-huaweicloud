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

// @API Workspace GET /v1/{project_id}/app-servers/access-agent/actions/show-latest-version
func DataSourceAppHdaLatestVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppHdaLatestVersionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the HDA latest versions are located.`,
			},

			// Attributes
			"hda_latest_versions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of HDA latest versions that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"latest_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest version of the HDA.`,
						},
						"hda_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the HDA.`,
						},
					},
				},
			},
		},
	}
}

func listAppHdaLatestVersions(client *golangsdk.ServiceClient) ([]interface{}, error) {
	httpUrl := "v1/{project_id}/app-servers/access-agent/actions/show-latest-version"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenAppHdaLatestVersions(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"latest_version": utils.PathSearch("latest_version", item, nil),
			"hda_type":       utils.PathSearch("hda_type", item, nil),
		})
	}

	return result
}

func dataSourceAppHdaLatestVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	hdaLatestVersions, err := listAppHdaLatestVersions(client)
	if err != nil {
		return diag.Errorf("error querying HDA latest versions: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("hda_latest_versions", flattenAppHdaLatestVersions(hdaLatestVersions)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

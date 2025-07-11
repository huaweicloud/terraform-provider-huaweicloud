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

// @API Workspace GET /v1/{project_id}/session-type
func DataSourceSessionTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSessionTypesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the resource.`,
			},
			"session_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of session types.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_spec_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource specification code.`,
						},
						"session_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The session type.`,
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource type.`,
						},
						"cloud_service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The code of cloud service type to which the resource belongs.`,
						},
					},
				},
			},
		},
	}
}

func getSessionTypes(client *golangsdk.ServiceClient) ([]interface{}, error) {
	httpUrl := "v1/{project_id}/session-type"
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

	return utils.PathSearch("session_types", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func dataSourceSessionTypesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace app client: %s", err)
	}

	sessionTypes, err := getSessionTypes(client)
	if err != nil {
		return diag.Errorf("error querying session types: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("session_types", flattenSessionTypes(sessionTypes)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSessionTypes(sessionTypes []interface{}) []map[string]interface{} {
	if len(sessionTypes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(sessionTypes))
	for _, sessionType := range sessionTypes {
		result = append(result, map[string]interface{}{
			"resource_spec_code": utils.PathSearch("resource_spec_code", sessionType, nil),
			"session_type":       utils.PathSearch("session_type", sessionType, nil),
			"resource_type":      utils.PathSearch("resource_type", sessionType, nil),
			"cloud_service_type": utils.PathSearch("cloud_service_type", sessionType, nil),
		})
	}
	return result
}

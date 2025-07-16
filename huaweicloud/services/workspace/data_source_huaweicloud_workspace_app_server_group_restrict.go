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

// @API Workspace GET /v1/{project_id}/app-server-groups/resources/restrict
func DataSourceAppServerGroupRestrict() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppServerGroupRestrictRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the resource.`,
			},
			"max_session": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum number of connection sessions per server.`,
			},
			"max_group_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum number of server groups that can be created by the tenant.`,
			},
		},
	}
}

func getAppServerGroupRestrict(client *golangsdk.ServiceClient) (map[string]interface{}, error) {
	httpUrl := "v1/{project_id}/app-server-groups/resources/restrict"
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

	return respBody.(map[string]interface{}), nil
}

func dataSourceAppServerGroupRestrictRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	resp, err := getAppServerGroupRestrict(client)
	if err != nil {
		return diag.Errorf("error getting app server group restrict: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("max_session", utils.PathSearch("max_session", resp, nil)),
		d.Set("max_group_count", utils.PathSearch("max_group_count", resp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

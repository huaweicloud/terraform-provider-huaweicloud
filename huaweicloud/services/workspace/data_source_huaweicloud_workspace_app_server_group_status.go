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

// @API Workspace GET /v1/{project_id}/app-server-groups/{server_group_id}/state
func DataSourceAppServerGroupStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppServerGroupStatusRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the resource`,
			},
			"server_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The unique identifier of the server group.`,
			},
			"aps_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: `The number of servers in each status.`,
			},
		},
	}
}

func getAppServerGroupStatus(client *golangsdk.ServiceClient, serverGroupId string) (map[string]interface{}, error) {
	httpUrl := "v1/{project_id}/app-server-groups/{server_group_id}/state"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{server_group_id}", serverGroupId)

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

func dataSourceAppServerGroupStatusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace app client: %s", err)
	}

	serverGroupId := d.Get("server_group_id").(string)
	resp, err := getAppServerGroupStatus(client, serverGroupId)
	if err != nil {
		return diag.Errorf("error getting server group status: %s", err)
	}

	apsStatus := utils.PathSearch("aps_status", resp, make(map[string]interface{})).(map[string]interface{})
	if len(apsStatus) < 1 {
		return diag.Errorf("error getting server group status: status is empty")
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("aps_status", apsStatus),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

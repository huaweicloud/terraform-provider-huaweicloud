package apig

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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/resources/outline/apps
func DataSourceApplicationAuthorizeStatistic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationAuthorizeStatisticRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the applications are located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the applications belong.`,
			},

			// Attributes.
			"authed_nums": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of authorized applications.`,
			},
			"unauthed_nums": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of unauthorized applications.`,
			},
		},
	}
}

func getApplicationAuthorizeStatistic(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/resources/outline/apps"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func dataSourceApplicationAuthorizeStatisticRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	resp, err := getApplicationAuthorizeStatistic(client, instanceId)
	if err != nil {
		return diag.Errorf("error querying application authorize statistics under a specified instance (%s): %s",
			instanceId, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("authed_nums", utils.PathSearch("authed_nums", resp, nil)),
		d.Set("unauthed_nums", utils.PathSearch("unauthed_nums", resp, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

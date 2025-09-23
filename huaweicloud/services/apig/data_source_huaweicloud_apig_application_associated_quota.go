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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/bound-quota
func DataSourceApplicationAssociatedQuota() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationAssociatedQuotaRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the application is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the application belongs.`,
			},
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the application to be queried.`,
			},

			// Attributes.
			"app_quota_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the application quota.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the application quota.`,
			},
			"call_limits": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum number of times the application quota can be called.`,
			},
			"time_unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time unit of the quota limit.`,
			},
			"time_interval": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The time interval of the quota limit.`,
			},
			"remark": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the application quota.`,
			},
			"reset_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The first quota reset time point, in RFC3339 format.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the application quota, in RFC3339 format.`,
			},
			"bound_app_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of applications bound to the quota policy.`,
			},
		},
	}
}

func getApplicationAssociatedQuota(client *golangsdk.ServiceClient, instanceId, appId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/bound-quota"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{app_id}", appId)

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

func dataSourceApplicationAssociatedQuotaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("app_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	quota, err := getApplicationAssociatedQuota(client, instanceId, appId)
	if err != nil {
		return diag.Errorf("error querying application associated quota: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("app_quota_id", utils.PathSearch("app_quota_id", quota, nil)),
		d.Set("name", utils.PathSearch("name", quota, nil)),
		d.Set("call_limits", utils.PathSearch("call_limits", quota, nil)),
		d.Set("time_unit", utils.PathSearch("time_unit", quota, nil)),
		d.Set("time_interval", utils.PathSearch("time_interval", quota, nil)),
		d.Set("remark", utils.PathSearch("remark", quota, nil)),
		d.Set("reset_time", utils.PathSearch("reset_time", quota, nil)),
		d.Set("create_time", utils.PathSearch("create_time", quota, nil)),
		d.Set("bound_app_num", utils.PathSearch("bound_app_num", quota, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

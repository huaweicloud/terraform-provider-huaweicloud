package vpn

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPN GET /v5/{project_id}/p2c-vpn-gateways/jobs
func DataSourceVpnP2CGatewayJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceVpnP2CGatewayJobsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the instance ID of a VPN P2C gateway.`,
			},
			"jobs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the job list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the job ID.`,
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the instance ID of a VPN P2C gateway.`,
						},
						"job_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the upgrade operation.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the job status.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creation time.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the update time.`,
						},
						"sub_jobs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the sub-job info.`,
							Elem:        resourceSchemeP2CGatewayJobSubJob(),
						},
					},
				},
			},
		},
	}
}

func resourceSchemeP2CGatewayJobSubJob() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the job ID.`,
			},
			"job_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the job type.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the job status.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time.`,
			},
			"finished_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the end time.`,
			},
			"error_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates error information.`,
			},
		},
	}
}

func resourceVpnP2CGatewayJobsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	httpUrl := "v5/{project_id}/p2c-vpn-gateways/jobs"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	if v, ok := d.GetOk("resource_id"); ok {
		getPath += fmt.Sprintf("?resource_id=%s", v)
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving VPN P2C gateway jobs: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("jobs", flattenP2CGatewayJobs(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenP2CGatewayJobs(resp interface{}) []interface{} {
	rawParams := utils.PathSearch("jobs", resp, nil)

	if params, ok := rawParams.([]interface{}); ok && len(params) > 0 {
		result := make([]interface{}, 0, len(params))
		for _, p := range params {
			param := p.(map[string]interface{})
			result = append(result, map[string]interface{}{
				"id":          utils.PathSearch("id", param, nil),
				"resource_id": utils.PathSearch("resource_id", param, nil),
				"job_type":    utils.PathSearch("job_type", param, nil),
				"status":      utils.PathSearch("status", param, nil),
				"created_at":  utils.PathSearch("created_at", param, nil),
				"updated_at":  utils.PathSearch("updated_at", param, nil),
				"sub_jobs":    flattenP2CGatewayJobSubJobs(param),
			})
		}

		return result
	}

	return nil
}

func flattenP2CGatewayJobSubJobs(resp interface{}) []interface{} {
	rawParams := utils.PathSearch("sub_jobs", resp, nil)

	if params, ok := rawParams.([]interface{}); ok && len(params) > 0 {
		result := make([]interface{}, 0, len(params))
		for _, p := range params {
			param := p.(map[string]interface{})
			result = append(result, map[string]interface{}{
				"id":            utils.PathSearch("id", param, nil),
				"job_type":      utils.PathSearch("job_typel", param, nil),
				"status":        utils.PathSearch("status", param, nil),
				"created_at":    utils.PathSearch("created_at", param, nil),
				"finished_at":   utils.PathSearch("finished_at", param, nil),
				"error_message": utils.PathSearch("error_message", param, nil),
			})
		}

		return result
	}

	return nil
}

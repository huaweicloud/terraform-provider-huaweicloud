package drs

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS GET /v5/{project_id}/jobs/{job_id}/compare/users/{compare_job_id}
func DataSourceCompareUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCompareUsersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"compare_job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"user_compare_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"src_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tar_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"src_privileges": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tar_privileges": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_target_existed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"compare_result": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildCompareUsersQueryParams() string {
	return "?limit=100"
}

func dataSourceCompareUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs/{job_id}/compare/users/{compare_job_id}"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{job_id}", d.Get("job_id").(string))
	listPath = strings.ReplaceAll(listPath, "{compare_job_id}", d.Get("compare_job_id").(string))
	listPath += buildCompareUsersQueryParams()

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving DRS compare users: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("total_count", utils.PathSearch("total_count", listRespBody, nil)),
		d.Set("user_compare_info", flattenCompareUsers(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCompareUsers(respBody interface{}) []interface{} {
	userCompareInfoRaw := utils.PathSearch("user_compare_info", respBody, make([]interface{}, 0)).([]interface{})
	if len(userCompareInfoRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(userCompareInfoRaw))
	for _, item := range userCompareInfoRaw {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("id", item, nil),
			"src_user_name":     utils.PathSearch("src_user_name", item, nil),
			"tar_user_name":     utils.PathSearch("tar_user_name", item, nil),
			"src_privileges":    utils.PathSearch("src_privileges", item, nil),
			"tar_privileges":    utils.PathSearch("tar_privileges", item, nil),
			"is_target_existed": utils.PathSearch("is_target_existed", item, false),
			"compare_result":    utils.PathSearch("compare_result", item, nil),
			"created_at":        utils.PathSearch("created_at", item, nil),
			"updated_at":        utils.PathSearch("updated_at", item, nil),
		})
	}
	return result
}

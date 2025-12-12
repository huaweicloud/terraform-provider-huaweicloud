package swr

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

// @API SWR GET /v2/manage/namespaces/{namespace}/repos/{repository}/sync_job
func DataSourceSwrImageAutoSyncJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrImageAutoSyncJobsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the organization name.`,
			},
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the image repository name.`,
			},
			"jobs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the jobs.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the job ID.`,
						},
						"organization": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the name of the organization.`,
						},
						"override": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether to overwrite.`,
						},
						"remote_organization": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the target organization.`,
						},
						"remote_region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the target region.`,
						},
						"repo_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the repository name.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the synchronization status.`,
						},
						"sync_operator_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the operator account ID.`,
						},
						"sync_operator_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the operator account name.`,
						},
						"tag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the image tag.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the time when the job was created. It is the UTC standard time.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the time when the job was updated. It is the UTC standard time.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceSwrImageAutoSyncJobsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	repository := strings.ReplaceAll(d.Get("repository").(string), "/", "$")

	listHttpUrl := "v2/manage/namespaces/{namespace}/repos/{repository}/sync_job?filter=limit::10"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{namespace}", d.Get("organization").(string))
	listPath = strings.ReplaceAll(listPath, "{repository}", repository)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	offset := 0
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := fmt.Sprintf("%s|offset::%v", listPath, offset)
		listResp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return diag.Errorf("error querying SWR image auto sync jobs: %s", err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.Errorf("error flattening SWR image auto sync jobs response : %s", err)
		}

		jobs := listRespBody.([]interface{})
		if len(jobs) == 0 {
			break
		}
		for _, job := range jobs {
			results = append(results, map[string]interface{}{
				"id":                  utils.PathSearch("id", job, nil),
				"organization":        utils.PathSearch("namespace", job, nil),
				"override":            utils.PathSearch("override", job, nil),
				"remote_organization": utils.PathSearch("remoteNamespace", job, nil),
				"remote_region_id":    utils.PathSearch("remoteRegionId", job, nil),
				"repo_name":           utils.PathSearch("repoName", job, nil),
				"status":              utils.PathSearch("status", job, nil),
				"sync_operator_id":    utils.PathSearch("syncOperatorId", job, false),
				"sync_operator_name":  utils.PathSearch("syncOperatorName", job, nil),
				"tag":                 utils.PathSearch("tag", job, nil),
				"created_at":          utils.PathSearch("createdAt", job, nil),
				"updated_at":          utils.PathSearch("updatedAt", job, nil),
			})
		}

		offset += 10
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("jobs", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

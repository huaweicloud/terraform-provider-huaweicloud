package servicestage

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ServiceStage GET /v3/{project_id}/cas/applications/{application_id}/components/{component_id}/records
func DataSourceV3ComponentRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV3ComponentRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the components are located.`,
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the application to which the component belongs.`,
			},
			"component_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the component to which the records belong.`,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"begin_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The begin time of the component execution record.`,
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The end time of the component execution record.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the component execution record.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance ID of the component execution record.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version number of the component execution record.`,
						},
						"current_used": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether version is current used.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the component execution record.`,
						},
						"deploy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The deploy type of the component execution record.`,
						},
						"jobs": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        v3componentRecordJobSchema(),
							Description: "The list of component jobs.",
						},
					},
				},
				Description: "The list of component execution record.",
			},
		},
	}
}

func v3componentRecordJobSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"sequence": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The sequence of the job execution.`,
			},
			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The job ID.`,
			},
			"job_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The source URL of the component.`,
						},
						"first_batch_weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The weight of the first batch execution.`,
						},
						"first_batch_replica": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The replica of the first batch execution.`,
						},
						"replica": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The total replica number.`,
						},
						"remaining_batch": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The remaining batch number.`,
						},
					},
				},
				Description: `The job detail.`,
			},
		},
	}
}

func queryV3ComponentRecords(client *golangsdk.ServiceClient, applicationId, componentId string) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/cas/applications/{application_id}/components/{component_id}/records?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{application_id}", applicationId)
	listPath = strings.ReplaceAll(listPath, "{component_id}", componentId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		records := utils.PathSearch("records", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, records...)
		if len(records) < limit {
			break
		}
		offset += len(records)
	}

	return result, nil
}

func flattenV3ComponentRecordJobInfo(jobInfo interface{}) []map[string]interface{} {
	if jobInfo == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"source_url":          utils.PathSearch("source_url", jobInfo, nil),
			"first_batch_weight":  utils.PathSearch("first_batch_weight", jobInfo, nil),
			"first_batch_replica": utils.PathSearch("first_batch_replica", jobInfo, nil),
			"replica":             utils.PathSearch("replica", jobInfo, nil),
			"remaining_batch":     utils.PathSearch("remaining_batch", jobInfo, nil),
		},
	}
}

func flattenV3ComponentRecordJobs(jobs []interface{}) []map[string]interface{} {
	if len(jobs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(jobs))
	for _, job := range jobs {
		result = append(result, map[string]interface{}{
			"sequence": utils.PathSearch("sequence", job, nil),
			"job_id":   utils.PathSearch("job_id", job, nil),
			"job_info": flattenV3ComponentRecordJobInfo(utils.PathSearch("job_info", job, nil)),
		})
	}

	return result
}

func flattenV3ComponentRecords(records []interface{}) []map[string]interface{} {
	if len(records) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(records))
	for _, record := range records {
		result = append(result, map[string]interface{}{
			"begin_time":   utils.PathSearch("begin_time", record, nil),
			"end_time":     utils.PathSearch("end_time", record, nil),
			"description":  utils.JsonToString(utils.PathSearch("description", record, nil)),
			"instance_id":  utils.PathSearch("instance_id", record, nil),
			"version":      utils.PathSearch("version", record, nil),
			"current_used": utils.PathSearch("current_used", record, nil),
			"status":       utils.PathSearch("status", record, nil),
			"deploy_type":  utils.PathSearch("deploy_type", record, nil),
			"jobs":         flattenV3ComponentRecordJobs(utils.PathSearch("jobs", record, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func dataSourceV3ComponentRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		applicationId = d.Get("application_id").(string)
		componentId   = d.Get("component_id").(string)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	records, err := queryV3ComponentRecords(client, applicationId, componentId)
	if err != nil {
		return diag.Errorf("error getting component records: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenV3ComponentRecords(records)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

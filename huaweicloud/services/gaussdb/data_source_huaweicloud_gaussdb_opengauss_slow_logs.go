package gaussdb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/slow-log/download
// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/slow-log/download
func DataSourceOpenGaussSlowLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBMysqlSlowLogsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of the slow logs.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workflow_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_link": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bucket_name": {
							Type:     schema.TypeString,
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
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGaussDBMysqlSlowLogsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	err = createSlowLogs(d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := waitForSlowLogLinkCompleted(ctx, client, d.Get("instance_id").(string), d.Timeout(schema.TimeoutRead))
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("list", flattenGaussDBMysqlSlowLogs(res)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func createSlowLogs(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/slow-log/download"
	)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return fmt.Errorf("error retrieving GaussDB OpenGauss slow logs: %s", err)
	}
	return nil
}

func waitForSlowLogLinkCompleted(ctx context.Context, client *golangsdk.ServiceClient, instanceId string,
	timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"EXPORTING"},
		Target:       []string{"SUCCESS"},
		Refresh:      openGaussSlowLogLinkRefreshFunc(client, instanceId),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	res, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func openGaussSlowLogLinkRefreshFunc(client *golangsdk.ServiceClient, instanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := getSlowLogs(client, instanceId)
		if err != nil {
			return nil, "ERROR", err
		}
		slowLogs := utils.PathSearch("list", resp, make([]interface{}, 0)).([]interface{})
		allLogStatusSuccess := true
		for _, slowLog := range slowLogs {
			status := utils.PathSearch("status", slowLog, "").(string)
			if status == "" {
				return nil, "ERROR", fmt.Errorf("error get slow log of the instance: %s", instanceId)
			}
			if status == "EXPORTING" {
				allLogStatusSuccess = false
			} else if status != "SUCCESS" {
				return nil, status, fmt.Errorf("error get slow log of the instance(%s), the satus is:%s",
					instanceId, status)
			}
		}
		if allLogStatusSuccess {
			return resp, "SUCCESS", nil
		}
		return resp, "EXPORTING", nil
	}
}

func getSlowLogs(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/slow-log/download"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GaussDB OpenGauss slow logs: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	return getRespBody, err
}

func flattenGaussDBMysqlSlowLogs(resp interface{}) []map[string]interface{} {
	slowLogsJson := utils.PathSearch("list", resp, make([]interface{}, 0))
	slowLogsArray := slowLogsJson.([]interface{})
	if len(slowLogsArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(slowLogsArray))
	for _, slowLog := range slowLogsArray {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", slowLog, nil),
			"instance_id": utils.PathSearch("instance_id", slowLog, nil),
			"node_id":     utils.PathSearch("node_id", slowLog, nil),
			"workflow_id": utils.PathSearch("workflow_id", slowLog, nil),
			"file_name":   utils.PathSearch("file_name", slowLog, nil),
			"file_size":   utils.PathSearch("file_size", slowLog, nil),
			"file_link":   utils.PathSearch("file_link", slowLog, nil),
			"bucket_name": utils.PathSearch("bucket_name", slowLog, nil),
			"created_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("created_at", slowLog,
				float64(0)).(float64))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("updated_at", slowLog,
				float64(0)).(float64))/1000, false),
			"version": utils.PathSearch("version", slowLog, nil),
			"status":  utils.PathSearch("status", slowLog, nil),
			"message": utils.PathSearch("message", slowLog, nil),
		})
	}
	return result
}

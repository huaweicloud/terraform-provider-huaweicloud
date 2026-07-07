package taurusdb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/{node_id}/slowlog-download
func DataSourceTaurusDBSlowLogsDownloadLinks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBSlowLogsDownloadLinksRead,

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
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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
						"create_at": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTaurusDBSlowLogsDownloadLinksRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		nodeId     = d.Get("node_id").(string)
	)
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	links, err := waitingForSlowLogsExportJobCompleted(ctx, client, d.Timeout(schema.TimeoutRead), instanceId, nodeId)
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB slow logs download links: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("list", flattenTaurusDBSlowLogsDownloadLinks(links)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func waitingForSlowLogsExportJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	instanceId, nodeId string) ([]interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"EXPORTING"},
		Target:       []string{"SUCCESS"},
		Refresh:      waitSlowLogsLinksStatusRefreshFunc(client, instanceId, nodeId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	links, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error waiting for TaurusDB slow logs download links to be ready: %s", err)
	}
	return links.([]interface{}), nil
}

func waitSlowLogsLinksStatusRefreshFunc(client *golangsdk.ServiceClient, instanceId, nodeId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		links, err := getTaurusDBSlowLogsDownloadLinks(client, instanceId, nodeId)
		if err != nil {
			return nil, "ERROR", err
		}
		for _, item := range links {
			status, _ := utils.PathSearch("status", item, "").(string)
			fileLink, _ := utils.PathSearch("file_link", item, "").(string)
			if status != "SUCCESS" || fileLink == "" {
				return links, "EXPORTING", nil
			}
		}
		return links, "SUCCESS", nil
	}
}
func getTaurusDBSlowLogsDownloadLinks(client *golangsdk.ServiceClient, instanceId, nodeId string) ([]interface{}, error) {
	httpUrl := "v3/{project_id}/instances/{instance_id}/{node_id}/slowlog-download"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{node_id}", nodeId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	list := utils.PathSearch("list", getRespBody, make([]interface{}, 0)).([]interface{})

	return list, nil
}

func flattenTaurusDBSlowLogsDownloadLinks(resp []interface{}) []interface{} {
	result := make([]interface{}, 0, len(resp))
	for _, item := range resp {
		result = append(result, map[string]interface{}{
			"job_id":     utils.PathSearch("job_id", item, nil),
			"file_name":  utils.PathSearch("file_name", item, nil),
			"status":     utils.PathSearch("status", item, nil),
			"file_size":  utils.PathSearch("file_size", item, nil),
			"file_link":  utils.PathSearch("file_link", item, nil),
			"create_at":  utils.PathSearch("create_at", item, nil),
			"updated_at": utils.PathSearch("updated_at", item, nil),
		})
	}
	return result
}

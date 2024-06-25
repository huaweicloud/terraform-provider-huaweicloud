package rds

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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS POST /v3/{project_id}/instances/{instance_id}/errorlog-download
func DataSourceRdsErrorLogLink() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsErrorLogLinkRead,

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
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the RDS instance.`,
			},
			"file_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the file.`,
			},
			"file_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the file size in KB.`,
			},
			"file_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the download link.`,
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the creation time.`,
			},
		},
	}
}

func dataSourceRdsErrorLogLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	retryFunc := func() (interface{}, bool, error) {
		res, err := waitForErrorLogLinkCompleted(ctx, client, instanceID, d.Timeout(schema.TimeoutRead))
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	resp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutRead),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error retrieving RDS slow log link: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("file_name", utils.PathSearch("list[0].file_name", resp, nil)),
		d.Set("file_size", utils.PathSearch("list[0].file_size", resp, nil)),
		d.Set("file_link", utils.PathSearch("list[0].file_link", resp, nil)),
		d.Set("created_at", utils.PathSearch("list[0].create_at", resp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func waitForErrorLogLinkCompleted(ctx context.Context, client *golangsdk.ServiceClient, instanceID string,
	timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"EXPORTING"},
		Target:       []string{"SUCCESS"},
		Refresh:      rdsErrorLogLinkRefreshFunc(client, instanceID),
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

func rdsErrorLogLinkRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := errorLogLink(client, instanceID)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch("list[0].status", resp, nil)
		if status == nil {
			return nil, "ERROR", fmt.Errorf("error get error log of the instance: %s", instanceID)
		}

		return resp, status.(string), nil
	}
}

func errorLogLink(client *golangsdk.ServiceClient, instanceID string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/errorlog-download"
	)
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{instance_id}", instanceID)

	opt := golangsdk.RequestOpts{KeepResponseBody: true}
	resp, err := client.Request("POST", path, &opt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

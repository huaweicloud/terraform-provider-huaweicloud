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

// @API RDS POST /v3/{project_id}/instances/{instance_id}/slowlog-download
func DataSourceRdsSlowLogLink() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsSlowLogLinkRead,

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
				Optional:    true,
				Description: `Specifies the name of the file to be downloaded.`,
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

func dataSourceRdsSlowLogLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	fileName := d.Get("file_name").(string)
	resp, err := waitForSlowLogLinkCompleted(ctx, client, instanceID, fileName, d.Timeout(schema.TimeoutRead))
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code||errCode", "DBS.280343"),
			"error retrieving RDS slow log link")
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

func slowLogLink(client *golangsdk.ServiceClient, instanceID, fileName string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/slowlog-download"
	)
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{instance_id}", instanceID)

	opt := golangsdk.RequestOpts{KeepResponseBody: true}
	opt.JSONBody = utils.RemoveNil(buildCreateSlowLogLinkBodyParams(fileName))
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

func buildCreateSlowLogLinkBodyParams(fileName string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"file_name": utils.ValueIgnoreEmpty(fileName),
	}
	return bodyParams
}

func waitForSlowLogLinkCompleted(ctx context.Context, client *golangsdk.ServiceClient, instanceID, fileName string,
	timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"EXPORTING"},
		Target:       []string{"SUCCESS"},
		Refresh:      rdsSlowLogLinkRefreshFunc(client, instanceID, fileName),
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

func rdsSlowLogLinkRefreshFunc(client *golangsdk.ServiceClient, instanceID, fileName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := slowLogLink(client, instanceID, fileName)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch("list[0].status", resp, nil)
		if status == nil {
			return nil, "ERROR", fmt.Errorf("error get slow log of the instance: %s", instanceID)
		}

		return resp, status.(string), nil
	}
}

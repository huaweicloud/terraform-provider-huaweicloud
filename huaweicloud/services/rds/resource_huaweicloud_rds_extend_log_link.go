package rds

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var extendLogLinkNonUpdatableParams = []string{"instance_id", "file_name"}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/xellog-download
func ResourceRdsExtendLogLink() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsExtendLogLinkCreate,
		ReadContext:   resourceRdsExtendLogLinkRead,
		UpdateContext: resourceRdsExtendLogLinkUpdate,
		DeleteContext: resourceRdsExtendLogLinkDelete,

		CustomizeDiff: config.FlexibleForceNew(extendLogLinkNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the RDS instance.`,
			},
			"file_name": {
				Type:        schema.TypeString,
				Required:    true,
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the last update time.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceRdsExtendLogLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	_, err = extendLogLink(client, instanceID, fileName)
	if err != nil {
		return diag.Errorf("error creating RDS extend log link: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceID, fileName))

	err = waitForExtendLogLinkCompleted(ctx, client, instanceID, fileName, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRdsExtendLogLinkRead(ctx, d, meta)
}

func resourceRdsExtendLogLinkRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	resp, err := extendLogLink(client, instanceID, fileName)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code||errCode", "DBS.280343"),
			"error retrieving RDS extend log link")
	}
	logInfo := utils.PathSearch(fmt.Sprintf("list|[?file_name=='%s']|[0]", fileName), resp, nil)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("file_name", fileName),
		d.Set("file_size", utils.PathSearch("file_size", logInfo, nil)),
		d.Set("file_link", utils.PathSearch("file_link", logInfo, nil)),
		d.Set("created_at", utils.PathSearch("create_at", logInfo, nil)),
		d.Set("updated_at", utils.PathSearch("update_at", logInfo, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func extendLogLink(client *golangsdk.ServiceClient, instanceID, fileName string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/xellog-download"
	)
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{instance_id}", instanceID)

	opt := golangsdk.RequestOpts{KeepResponseBody: true}
	opt.JSONBody = utils.RemoveNil(buildCreateExtendLogLinkBodyParams(fileName))
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

func buildCreateExtendLogLinkBodyParams(fileName string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"file_name": fileName,
	}
	return bodyParams
}

func waitForExtendLogLinkCompleted(ctx context.Context, client *golangsdk.ServiceClient, instanceID, fileName string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"EXPORTING"},
		Target:       []string{"SUCCESS"},
		Refresh:      rdsExtendLogLinkRefreshFunc(client, instanceID, fileName),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for RDS instance generating extend log link to be completed: %s ", err)
	}
	return nil
}

func rdsExtendLogLinkRefreshFunc(client *golangsdk.ServiceClient, instanceID, fileName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := extendLogLink(client, instanceID, fileName)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch(fmt.Sprintf("list|[?file_name=='%s']|[0].status", fileName), resp, nil)
		if status == nil {
			return nil, "ERROR", fmt.Errorf("error get extend log info by file name: %s", fileName)
		}

		return resp, status.(string), nil
	}
}

func resourceRdsExtendLogLinkUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsExtendLogLinkDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting extend log link is not supported. The extend log link is only removed from the state," +
		" but it remains in the cloud. And the instance doesn't return to the state before generate extend log link."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

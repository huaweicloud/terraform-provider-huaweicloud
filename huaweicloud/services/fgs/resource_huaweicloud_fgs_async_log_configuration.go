package fgs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API FunctionGraph POST /v2/{project_id}/fgs/functions/enable-async-status-logs
// @API FunctionGraph GET /v2/{project_id}/fgs/functions/async-status-log-detail
// @API LTS GET /v2/{project_id}/groups/{log_group_id}/streams
// @API LTS GET /v2/{project_id}/transfers
// @API LTS DELETE /v2/{project_id}/transfers
// @API LTS DELETE /v2/{project_id}/groups/{log_group_id}
func ResourceAsyncLogConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAsyncLogConfigurationCreate,
		ReadContext:   resourceAsyncLogConfigurationRead,
		UpdateContext: resourceAsyncLogConfigurationUpdate,
		DeleteContext: resourceAsyncLogConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the async log configuration is located.`,
			},

			// Optional parameter(s).
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to force delete the LTS resources corresponding to the async log configuration.`,
			},
			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: `The maximum number of retries for the create operation when encountering internal service errors.`,
			},

			// Attribute(s).
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The LTS log group ID used to manage async status logs.`,
			},
			"group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The LTS log group name used to manage async status logs.`,
			},
			"stream_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The LTS log stream ID used to manage async status logs.`,
			},
			"stream_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The LTS log stream name used to manage async status logs.`,
			},
		},
	}
}

func resourceAsyncLogConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	var (
		httpUrl   = "v2/{project_id}/fgs/functions/enable-async-status-logs"
		createOpt = golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		}
		maxRetries = d.Get("max_retries").(int)
	)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	log.Printf("[DEBUG] The max_retries is %d", maxRetries)
	retryCount := 0
	for retryCount < maxRetries+1 {
		_, err = client.Request("POST", createPath, &createOpt)
		if err == nil {
			break
		}

		parsedErr := common.ConvertExpected500ErrInto404Err(err, "error_code", "FSS.0500")
		if _, ok := parsedErr.(golangsdk.ErrDefault404); !ok {
			break
		}
		// lintignore:R018
		time.Sleep(10 * time.Second)
		log.Printf("[DEBUG] Service is busy or environment variable of LTS is not refreshed, try again")

		retryCount++
		continue
	}
	if err != nil {
		return diag.Errorf("after %d retries, the async log configuration still reports an error: %s",
			maxRetries-1, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceAsyncLogConfigurationRead(ctx, d, meta)
}

func GetAsyncLogConfiguration(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v2/{project_id}/fgs/functions/async-status-log-detail"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func GetAsyncLogConfigurationAndCheck(fgsClient, ltsClient *golangsdk.ServiceClient) (interface{}, error) {
	respBody, err := GetAsyncLogConfiguration(fgsClient)
	if err != nil {
		return nil, err
	}

	groupId := utils.PathSearch("group_id", respBody, "").(string)
	streamId := utils.PathSearch("stream_id", respBody, "").(string)
	_, err = lts.GetLogStreamById(ltsClient, groupId, streamId)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func resourceAsyncLogConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	fgsClient, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}
	ltsClient, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	respBody, err := GetAsyncLogConfigurationAndCheck(fgsClient, ltsClient)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving async log configuration")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		// Optional parameter(s).
		d.Set("force_delete", d.Get("force_delete").(bool)),
		// Attribute(s).
		d.Set("group_id", utils.PathSearch("group_id", respBody, "")),
		d.Set("group_name", utils.PathSearch("group_name", respBody, "")),
		d.Set("stream_id", utils.PathSearch("stream_id", respBody, "")),
		d.Set("stream_name", utils.PathSearch("stream_name", respBody, "")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAsyncLogConfigurationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func buildQueryTransfersBodyParams(logGroupName string) string {
	return fmt.Sprintf("&log_group_name=%v", logGroupName)
}

func forceDeleteAsyncLogConfiguration(client *golangsdk.ServiceClient, groupId, groupName string) error {
	log.Printf("[Lance] Ready to execute forceDeleteAsyncLogConfiguration")
	log.Printf("[Lance] groupId: %s", groupId)
	log.Printf("[Lance] groupName: %s", groupName)
	transfers, err := lts.ListTransfers(client, buildQueryTransfersBodyParams(groupName))
	if err != nil {
		log.Printf("[Lance] Error querying log transfers: %s", err)
		// error LTS.0201 means that the log group does not exist.
		return common.ConvertExpected400ErrInto404Err(err, "error_code", "LTS.0201")
	}

	// Before forcibly deleting the log group, make sure that all log transfers have been deleted.
	for _, transfer := range transfers {
		transferId := utils.PathSearch("log_transfer_id", transfer, "").(string)
		if transferId == "" {
			log.Printf("[ERROR] Unable to find log_transfer_id from log transfer configuration of the log group (%s).",
				groupName)
			continue
		}
		err = lts.DeleteTransferById(client, transferId)
		if err != nil {
			return fmt.Errorf("error deleting log transfer (%s): %s", transferId, err)
		}
	}
	log.Printf("[Lance] Successfully deleted all log transfers")

	return lts.DeleteGroupById(client, groupId)
}

func resourceAsyncLogConfigurationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.Get("force_delete").(bool) {
		warnMsg := `Skip the forced deletion of the log group log stream because force_delete is set to false.
The remaining log group and log stream will affect the next deployment of the resource. Before that, make sure the
corresponding log group and log stream are deleted.`
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  warnMsg,
			},
		}
	}

	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	fgsClient, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}
	ltsClient, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	respBody, err := GetAsyncLogConfigurationAndCheck(fgsClient, ltsClient)
	if err != nil {
		return diag.Errorf("error retrieving async log configuration: %s", err)
	}

	var (
		groupName = utils.PathSearch("group_name", respBody, "").(string)
		groupId   = utils.PathSearch("group_id", respBody, "").(string)
	)

	err = forceDeleteAsyncLogConfiguration(ltsClient, groupId, groupName)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error deleting LTS log group (%s) for the async log configuration", groupId))
	}
	return nil
}

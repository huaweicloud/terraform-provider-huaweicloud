package dataarts

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	dataServiceInstanceLogDumpResourceNotFoundCodes = []string{
		"DLM.4001", // Workspace does not exist.
		"DLM.4292", // Cluster instance does not exist.
	}
	dataServiceInstanceLogDumpNonUpdatableParams = []string{
		"workspace_id",
		"instance_id",
		"type",
		"log_group_id",
		"log_group_name",
		"log_stream_id",
		"log_stream_name",
	}
)

const (
	dataServiceInstanceOBSLogDump = "obs"
	dataServiceInstanceLTSLogDump = "lts"
)

// @API DataArtsStudio PUT /v1/{project_id}/service/instances/{instance_id}/lts-log-dump
// @API DataArtsStudio PUT /v1/{project_id}/service/instances/{instance_id}/obs-log-dump
// @API DataArtsStudio GET /v1/{project_id}/service/instances/{instance_id}
func ResourceDataServiceInstanceLogDump() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataServiceInstanceLogDumpCreate,
		ReadContext:   resourceDataServiceInstanceLogDumpRead,
		UpdateContext: resourceDataServiceInstanceLogDumpUpdate,
		DeleteContext: resourceDataServiceInstanceLogDumpDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDataServiceInstanceLogDumpImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(dataServiceInstanceLogDumpNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the log dump is located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the log dump belongs.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the exclusive cluster to which the log dump belongs.`,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					dataServiceInstanceOBSLogDump, dataServiceInstanceLTSLogDump,
				}, false),
				Description: `The type of the log dump to be configured.`,
			},

			// Optional parameters.
			"log_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the LTS log group.`,
			},
			"log_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the LTS log group.`,
			},
			"log_stream_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the LTS log stream.`,
			},
			"log_stream_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the LTS log stream.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildDataServiceInstanceLogDumpBodyParams(d *schema.ResourceData, enable bool) map[string]interface{} {
	return utils.RemoveNil(map[string]interface{}{
		"log_dump":        enable,
		"log_group_id":    utils.ValueIgnoreEmpty(d.Get("log_group_id")),
		"log_group_name":  utils.ValueIgnoreEmpty(d.Get("log_group_name")),
		"log_stream_id":   utils.ValueIgnoreEmpty(d.Get("log_stream_id")),
		"log_stream_name": utils.ValueIgnoreEmpty(d.Get("log_stream_name")),
	})
}

func doDataServiceInstanceLogDump(client *golangsdk.ServiceClient, d *schema.ResourceData, enable bool) error {
	httpUrl := "v1/{project_id}/service/instances/{instance_id}/{type}-log-dump"
	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{instance_id}", d.Get("instance_id").(string))
	actionPath = strings.ReplaceAll(actionPath, "{type}", d.Get("type").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
			"Dlm-Type":     "EXCLUSIVE",
		},
		JSONBody: buildDataServiceInstanceLogDumpBodyParams(d, enable),
		OkCodes:  []int{204},
	}

	_, err := client.Request("PUT", actionPath, &opt)
	return err
}

func resourceDataServiceInstanceLogDumpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	if err := doDataServiceInstanceLogDump(client, d, true); err != nil {
		return diag.Errorf("error enabling %s log dump for exclusive cluster (%s): %s", d.Get("type").(string), instanceId, err)
	}

	d.SetId(instanceId)

	return resourceDataServiceInstanceLogDumpRead(ctx, d, meta)
}

func getDataServiceInstanceById(client *golangsdk.ServiceClient, workspaceId, instanceId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/service/instances/{instance_id}"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", instanceId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
			"Dlm-Type":     "EXCLUSIVE",
		},
	}

	resp, err := client.Request("GET", requestPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func GetDataServiceInstanceLogDump(client *golangsdk.ServiceClient, workspaceId, instanceId string) (interface{}, error) {
	respBody, err := getDataServiceInstanceById(client, workspaceId, instanceId)
	if err != nil {
		return nil, err
	}

	if !utils.PathSearch("log_dump_flag", respBody, false).(bool) && !utils.PathSearch("lts_log_dump_flag", respBody, false).(bool) {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/service/instances/{instance_id}",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the exclusive cluster (%s) log dump has been disabled", instanceId)),
			},
		}
	}

	return respBody, nil
}

func resourceDataServiceInstanceLogDumpRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Id()
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	respBody, err := GetDataServiceInstanceLogDump(client, d.Get("workspace_id").(string), instanceId)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", dataServiceInstanceLogDumpResourceNotFoundCodes...),
			fmt.Sprintf("error retrieving exclusive cluster log dump (%s)", instanceId),
		)
	}

	logDumpType := parseDataServiceInstanceLogDumpType(respBody, instanceId)
	if logDumpType == "" {
		return diag.Errorf("unable to parse the log dump type of exclusive cluster (%s)", instanceId)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("id", respBody, nil)),
		d.Set("type", logDumpType),
		d.Set("log_group_id", utils.PathSearch("log_group_id", respBody, nil)),
		d.Set("log_group_name", utils.PathSearch("log_group_name", respBody, nil)),
		d.Set("log_stream_id", utils.PathSearch("log_stream_id", respBody, nil)),
		d.Set("log_stream_name", utils.PathSearch("log_stream_name", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func parseDataServiceInstanceLogDumpType(respBody interface{}, instanceId string) string {
	if utils.PathSearch("lts_log_dump_flag", respBody, false).(bool) {
		return dataServiceInstanceLTSLogDump
	}

	if utils.PathSearch("log_dump_flag", respBody, false).(bool) {
		return dataServiceInstanceOBSLogDump
	}

	log.Printf("[ERROR] unable to parse the log dump type of exclusive cluster (%s)", instanceId)
	return ""
}

func resourceDataServiceInstanceLogDumpUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDataServiceInstanceLogDumpDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dataarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	if err := doDataServiceInstanceLogDump(client, d, false); err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", dataServiceInstanceLogDumpResourceNotFoundCodes...),
			fmt.Sprintf("error disabling exclusive cluster log dump (%s)", d.Id()),
		)
	}

	return nil
}

func resourceDataServiceInstanceLogDumpImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be '<workspace_id>/<instance_id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])

	return []*schema.ResourceData{d}, d.Set("workspace_id", parts[0])
}

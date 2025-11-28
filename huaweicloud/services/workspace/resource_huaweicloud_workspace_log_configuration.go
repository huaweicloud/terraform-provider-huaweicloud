package workspace

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace POST /v2/{project_id}/user-events/lts-configurations
// @API Workspace GET /v2/{project_id}/user-events/lts-configurations
func ResourceLogConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogConfigurationCreate,
		ReadContext:   resourceLogConfigurationRead,
		UpdateContext: resourceLogConfigurationUpdate,
		DeleteContext: resourceLogConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the log configuration is located.`,
			},

			// Required parameters.
			"log_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the log group.`,
			},
			"log_stream_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the log stream.`,
			},
		},
	}
}

func buildLogConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"enable":        true,
		"log_group_id":  d.Get("log_group_id"),
		"log_stream_id": d.Get("log_stream_id"),
	}

	return utils.RemoveNil(bodyParams)
}

func resourceLogConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	httpUrl := "v2/{project_id}/user-events/lts-configurations"
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildLogConfigurationBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Workspace log configuration: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	return resourceLogConfigurationRead(ctx, d, meta)
}

func GetLogConfiguration(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v2/{project_id}/user-events/lts-configurations"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	if isEnabled := utils.PathSearch("enable", respBody, false).(bool); !isEnabled {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/user-events/lts-configurations",
				RequestId: "NONE",
				Body:      []byte("the LTS log configuration has been disabled"),
			},
		}
	}
	return respBody, nil
}

func resourceLogConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	respBody, err := GetLogConfiguration(client)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Workspace log configuration")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("log_group_id", utils.PathSearch("log_group_id", respBody, nil)),
		d.Set("log_stream_id", utils.PathSearch("log_stream_id", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceLogConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/user-events/lts-configurations"
	)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildLogConfigurationBodyParams(d),
	}

	_, err = client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating Workspace log configuration: %s", err)
	}

	return resourceLogConfigurationRead(ctx, d, meta)
}

func resourceLogConfigurationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	httpUrl := "v2/{project_id}/user-events/lts-configurations"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"enable": false,
		},
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error getting Workspace log configuration: %s", err)
	}

	return nil
}

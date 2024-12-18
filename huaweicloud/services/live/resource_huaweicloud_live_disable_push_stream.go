package live

import (
	"context"
	"fmt"
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

// @API Live POST /v1/{project_id}/stream/blocks
// @API Live GET /v1/{project_id}/stream/blocks
// @API Live PUT /v1/{project_id}/stream/blocks
// @API Live DELETE /v1/{project_id}/stream/blocks
func ResourceDisablePushStream() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDisablePushStreamCreate,
		ReadContext:   resourceDisablePushStreamRead,
		UpdateContext: resourceDisablePushStreamUpdate,
		DeleteContext: resourceDisablePushStreamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDisablePushStreamImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ingest domain name of the disabling push stream.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the application name of the disabling push stream.`,
			},
			"stream_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the stream name of the disabling push stream.`,
			},
			"resume_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the time of the resuming push stream.`,
			},
		},
	}
}

func resourceDisablePushStreamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	createHttpUrl := "v1/{project_id}/stream/blocks"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		JSONBody: utils.RemoveNil(buildDisablePushStreamBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating disabled push stream: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	return resourceDisablePushStreamRead(ctx, d, meta)
}

func buildDisablePushStreamBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"domain":      d.Get("domain_name"),
		"app_name":    d.Get("app_name"),
		"stream_name": d.Get("stream_name"),
		"resume_time": utils.ValueIgnoreEmpty(d.Get("resume_time")),
	}

	return params
}

func resourceDisablePushStreamRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		domainName = d.Get("domain_name").(string)
		appName    = d.Get("app_name").(string)
		streamName = d.Get("stream_name").(string)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	getResp, err := GetDisablePushStream(client, domainName, appName, streamName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving disabled push stream information")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("domain_name", domainName),
		d.Set("app_name", utils.PathSearch("app_name", getResp, nil)),
		d.Set("stream_name", utils.PathSearch("stream_name", getResp, nil)),
		d.Set("resume_time", utils.PathSearch("resume_time", getResp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetDisablePushStream(client *golangsdk.ServiceClient, domainName, appName, streamName string) (interface{}, error) {
	getHttpUrl := "v1/{project_id}/stream/blocks"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = fmt.Sprintf("%s?domain=%s&app_name=%s&stream_name=%s", getPath, domainName, appName, streamName)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", domainNameNotExistsCode)
	}

	getRespBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	block := utils.PathSearch("blocks|[0]", getRespBody, nil)
	if block == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return block, nil
}

func resourceDisablePushStreamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	if d.HasChange("resume_time") {
		updateHttpUrl := "v1/{project_id}/stream/blocks"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
			JSONBody: utils.RemoveNil(buildDisablePushStreamBodyParams(d)),
		}

		_, err := client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating disabled push stream information: %s", err)
		}
	}

	return resourceDisablePushStreamRead(ctx, d, meta)
}

func resourceDisablePushStreamDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		domainName = d.Get("domain_name").(string)
		appName    = d.Get("app_name").(string)
		streamName = d.Get("stream_name").(string)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	deleteHttpUrl := "v1/{project_id}/stream/blocks"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = fmt.Sprintf("%s?domain=%s&app_name=%s&stream_name=%s", deletePath, domainName, appName, streamName)

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// Before deleting, call the query API, if query no result , then process `CheckDeleted` logic.
	_, err = GetDisablePushStream(client, domainName, appName, streamName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving disabled push stream information")
	}

	// Call the deletion API
	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected403ErrInto404Err(err, "error_code", domainNameNotExistsCode),
			"error deleting disabled push stream")
	}

	// Successful deletion API call does not guarantee that the resource is successfully deleted.
	// Call the query API to confirm that the resource has been successfully deleted.
	_, err = GetDisablePushStream(client, domainName, appName, streamName)
	if err == nil {
		return diag.Errorf("error deleting disabled push stream: the disabled push stream still exists")
	}

	return nil
}

func resourceDisablePushStreamImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<domain_name>/<app_name>/<stream_name>', but got '%s'",
			importedId)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	mErr := multierror.Append(nil,
		d.Set("domain_name", parts[0]),
		d.Set("app_name", parts[1]),
		d.Set("stream_name", parts[2]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

package secmaster

import (
	"context"
	"fmt"
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

var pipeIndexNonUpdatableParams = []string{"workspace_id", "pipe_id"}

// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}/index
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}/index
func ResourcePipeIndex() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipeIndexCreate,
		ReadContext:   resourcePipeIndexRead,
		UpdateContext: resourcePipeIndexUpdate,
		DeleteContext: resourcePipeIndexDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePipeIndexImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(pipeIndexNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipe_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mapping": {
				// Convert field `mapping` to JSON string.
				Type:     schema.TypeString,
				Required: true,
			},
			"timestamp_field": {
				Type:     schema.TypeString,
				Required: true,
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

func buildPipeIndexBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"mapping":         utils.StringToJson(d.Get("mapping").(string)),
		"pipe_id":         d.Get("pipe_id"),
		"status":          "open",
		"timestamp_field": d.Get("timestamp_field"),
	}
}

func configPipeIndexPut(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}/index"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{pipe_id}", d.Get("pipe_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         buildPipeIndexBodyParams(d),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func resourcePipeIndexCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	if err := configPipeIndexPut(client, d); err != nil {
		return diag.Errorf("error creating SecMaster pipe index: %s", err)
	}

	d.SetId(d.Get("pipe_id").(string))

	return resourcePipeIndexRead(ctx, d, meta)
}

func GetPipeIndexDetail(client *golangsdk.ServiceClient, workspaceID, pipeID string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}/index"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceID)
	requestPath = strings.ReplaceAll(requestPath, "{pipe_id}", pipeID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	if utils.PathSearch("status", respBody, "").(string) == "closed" {
		return nil, golangsdk.ErrDefault404{}
	}
	return respBody, nil
}

func resourcePipeIndexRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	respBody, err := GetPipeIndexDetail(client, d.Get("workspace_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SecMaster pipe index")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("mapping", utils.JsonToString(utils.PathSearch("mapping", respBody, nil))),
		d.Set("pipe_id", utils.PathSearch("pipe_id", respBody, nil)),
		d.Set("timestamp_field", utils.PathSearch("timestamp_field", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePipeIndexUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	// This operation is an overriding operation.
	if err := configPipeIndexPut(client, d); err != nil {
		return diag.Errorf("error updating SecMaster pipe index: %s", err)
	}

	return resourcePipeIndexRead(ctx, d, meta)
}

func resourcePipeIndexDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}/index"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{pipe_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody: map[string]interface{}{
			"mapping":         utils.StringToJson(d.Get("mapping").(string)),
			"pipe_id":         d.Id(),
			"status":          "closed",
			"timestamp_field": d.Get("timestamp_field"),
		},
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster pipe index: %s", err)
	}

	return nil
}

func resourcePipeIndexImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	importIdParts := strings.Split(importId, "/")
	if len(importIdParts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, want <workspace_id>/<pipe_id>, but got %s", importId)
	}

	d.SetId(importIdParts[1])
	mErr := multierror.Append(
		d.Set("workspace_id", importIdParts[0]),
		d.Set("pipe_id", importIdParts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

package coc

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var scriptOrderOperationNonUpdatableParams = []string{"execute_uuid", "operation_type", "batch_index", "instance_id"}

// @API COC PUT /v1/job/script/orders/{execute_uuid}/operation
func ResourceScriptOrderOperation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScriptOrderOperationCreate,
		ReadContext:   resourceScriptOrderOperationRead,
		UpdateContext: resourceScriptOrderOperationUpdate,
		DeleteContext: resourceScriptOrderOperationDelete,

		CustomizeDiff: config.FlexibleForceNew(scriptOrderOperationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"execute_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operation_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"batch_index": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeInt,
				Optional: true,
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

func buildScriptOrderOperationCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"operation_type": d.Get("operation_type"),
		"batch_index":    utils.ValueIgnoreEmpty(d.Get("batch_index")),
		"instance_id":    utils.ValueIgnoreEmpty(d.Get("instance_id")),
	}

	return bodyParams
}

func resourceScriptOrderOperationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/job/script/orders/{execute_uuid}/operation"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	executeUUID := d.Get("execute_uuid").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{execute_uuid}", executeUUID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildScriptOrderOperationCreateOpts(d)),
	}

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error operating the COC script order (%s): %s", executeUUID, err)
	}

	d.SetId(executeUUID)

	return resourceScriptOrderOperationRead(ctx, d, meta)
}

func resourceScriptOrderOperationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceScriptOrderOperationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceScriptOrderOperationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting script order operation resource is not supported. The script order operation resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

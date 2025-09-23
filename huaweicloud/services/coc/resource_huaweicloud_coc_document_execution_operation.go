package coc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var documentExecutionOperationNonUpdatableParams = []string{"execution_id", "operate_type"}

// @API COC POST /v1/executions
func ResourceDocumentExecutionOperation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDocumentExecutionOperationCreate,
		ReadContext:   resourceDocumentExecutionOperationRead,
		UpdateContext: resourceDocumentExecutionOperationUpdate,
		DeleteContext: resourceDocumentExecutionOperationDelete,

		CustomizeDiff: config.FlexibleForceNew(documentExecutionOperationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"execution_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operate_type": {
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

func buildDocumentExecutionOperationCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"execution_id": d.Get("execution_id"),
		"operate_type": d.Get("operate_type"),
	}

	return bodyParams
}

func resourceDocumentExecutionOperationCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/executions"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	executionID := d.Get("execution_id").(string)
	createPath := client.Endpoint + httpUrl

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildDocumentExecutionOperationCreateOpts(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error operating the COC document execution (%s): %s", executionID, err)
	}

	d.SetId(executionID)

	return nil
}

func resourceDocumentExecutionOperationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDocumentExecutionOperationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDocumentExecutionOperationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting document execution operation resource is not supported. The document execution operation" +
		" resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

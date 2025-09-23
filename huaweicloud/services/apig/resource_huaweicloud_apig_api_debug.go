package apig

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var apigApiDebugNonUpdatableParams = []string{"instance_id", "api_id", "mode", "scheme", "method", "path", "body",
	"header", "query", "stage"}

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/apis/debug/{api_id}
func ResourceApigApiDebug() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApigApiDebugCreate,
		ReadContext:   resourceApigApiDebugRead,
		UpdateContext: resourceApigApiDebugUpdate,
		DeleteContext: resourceApigApiDebugDelete,

		CustomizeDiff: config.FlexibleForceNew(apigApiDebugNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the dedicated instance to which the API belongs is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the API belongs.`,
			},
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the API to be debugged.`,
			},
			"mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The debug mode.`,
			},
			"scheme": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The request protocol.`,
			},
			"method": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The request method of the API.`,
			},
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The request path of the API.`,
			},

			// Optional parameters.
			"body": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The request message body of the API.`,
			},
			"header": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The request header parameters of the API, in JSON format.`,
			},
			"query": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The request query parameters of the API, in JSON format.`,
			},
			"stage": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The runtime environment for debug request.`,
			},

			// Attributes.
			"request": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The debug request message content.`,
			},
			"response": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The debug response message content.`,
			},
			"latency": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The debug latency in milliseconds.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildApiDebugBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"mode":   d.Get("mode"),
		"scheme": d.Get("scheme"),
		"method": d.Get("method"),
		"path":   d.Get("path"),
		"body":   utils.ValueIgnoreEmpty(d.Get("body")),
		"header": utils.StringToJson(d.Get("header").(string)),
		"query":  utils.StringToJson(d.Get("query").(string)),
		"stage":  utils.ValueIgnoreEmpty(d.Get("stage")),
	}
}

func resourceApigApiDebugCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/apis/debug/{api_id}"
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))
	createPath = strings.ReplaceAll(createPath, "{api_id}", d.Get("api_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildApiDebugBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error debugging the API: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("request", utils.PathSearch("request", respBody, "")),
		d.Set("response", utils.PathSearch("response", respBody, "")),
		d.Set("latency", utils.PathSearch("latency", respBody, 0)),
	)
	err = mErr.ErrorOrNil()
	if err != nil {
		return diag.Errorf("error setting API debug attributes: %s", err)
	}

	return resourceApigApiDebugRead(ctx, d, meta)
}

func resourceApigApiDebugRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApigApiDebugUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApigApiDebugDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for debugging the API. Deleting this resource will
not clear the corresponding debug record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

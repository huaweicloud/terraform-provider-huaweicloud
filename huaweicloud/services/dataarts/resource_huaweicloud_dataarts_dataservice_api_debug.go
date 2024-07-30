package dataarts

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio POST /v1/{project_id}/service/apis/{api_id}/instances/{instance_id}/test
func ResourceDataServiceApiDebug() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataServiceApiDebugCreate,
		ReadContext:   resourceDataServiceApiDebugRead,
		DeleteContext: resourceDataServiceApiDebugDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the API is located.`,
			},

			// Parameters in request header
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the workspace to which the API belongs.`,
			},

			// Arguments
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the catalog where the API is located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The exclusive cluster ID.`,
			},
			"params": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  `The request parameters in which to debug the API, in JSON format.`,
				ValidateFunc: validation.StringIsJSON,
			},
			// Sometimes, debug requests may be affected by network fluctuations and time out.
			// Retry as appropriate to resolve the issue.
			"max_retries": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  3,
				Description: utils.SchemaDesc(
					`The maximum retry number of the API debug operation.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},

			// Attributes
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The request URL of this API debug operation.`,
			},
			"result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The result detail of this API debug operation.`,
			},
			"timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The timeout of this API debug operation.`,
			},
			"request_header": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The request header of this API debug operation.`,
			},
			"response_header": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The response header of this API debug operation.`,
			},
		},
	}
}

func buildApiDebugBodyParams(params string) map[string]interface{} {
	parsedParams := make(map[string]interface{})
	err := json.Unmarshal([]byte(params), &parsedParams)
	if err != nil {
		log.Printf("[ERROR] Invalid type of the debug parameters, not json format")
	}
	return map[string]interface{}{
		"paras": parsedParams,
	}
}

func debugApi(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl     = "v1/{project_id}/service/apis/{api_id}/instances/{instance_id}/test"
		workspaceId = d.Get("workspace_id").(string)
		apiId       = d.Get("api_id").(string)
		instanceId  = d.Get("instance_id").(string)
	)
	debugPath := client.Endpoint + httpUrl
	debugPath = strings.ReplaceAll(debugPath, "{project_id}", client.ProjectID)
	debugPath = strings.ReplaceAll(debugPath, "{api_id}", apiId)
	debugPath = strings.ReplaceAll(debugPath, "{instance_id}", instanceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
			"Dlm-Type":     "EXCLUSIVE",
		},
		JSONBody: buildApiDebugBodyParams(d.Get("params").(string)),
	}

	requestResp, err := client.Request("POST", debugPath, &opt)
	if err != nil {
		return nil, fmt.Errorf("error debugging API (%s): %s", apiId, err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving API debug result: %s", err)
	}
	return respBody, nil
}

func parseApiDebugHeaders(headers interface{}) interface{} {
	jsonFilter, err := json.Marshal(headers)
	if err != nil {
		log.Printf("[ERROR] unable to convert the header content, not json format")
		return nil
	}
	return string(jsonFilter)
}

func apiDebugRetryFunc(client *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := debugApi(client, d)
		if err != nil {
			// After an error is identified, the log is recorded immediately and a retry is initiated.
			// Returns the error at this time will cause the retry to be interrupted.
			log.Printf("[ERROR] Failed to debug API, the error is: %s", err)
			return nil, "DEBUG_FAILED", nil
		}
		return resp, "SUCCESS", nil
	}
}

func resourceDataServiceApiDebugCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		maxRetries = d.Get("max_retries").(int)
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:        []string{"DEBUG_FAILED"},
		Target:         []string{"SUCCESS"},
		Refresh:        apiDebugRetryFunc(client, d),
		Timeout:        d.Timeout(schema.TimeoutCreate),
		PollInterval:   3 * time.Second,
		NotFoundChecks: maxRetries,
	}
	resp, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error debugging API, %d retries have been made, the error messages has been recorded to the log file", maxRetries)
	}
	d.SetId(utils.PathSearch("requestId", resp, "").(string))

	// Saving debug result in the CreateContext because of this resource do not have any query API to obtain the result.
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("url", utils.PathSearch("url", resp, "")),
		d.Set("result", utils.PathSearch("result", resp, "")),
		d.Set("timeout", utils.PathSearch("timecost", resp, "")),
		d.Set("request_header", parseApiDebugHeaders(utils.PathSearch("requestHeader", resp, make(map[string]interface{})))),
		d.Set("response_header", parseApiDebugHeaders(utils.PathSearch("responseHeader", resp, make(map[string]interface{})))),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving API debug resource fields: %s", err)
	}
	return resourceDataServiceApiDebugRead(ctx, d, meta)
}

func resourceDataServiceApiDebugRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This resource is only a one-time action resource for debugging the API.
	// There is no API for the provider to query the API debug history.
	return nil
}

func resourceDataServiceApiDebugDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for debugging the API. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

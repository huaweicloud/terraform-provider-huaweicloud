package geminidb

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GeminiDB PUT /v3.1/{project_id}/instances/{instance_id}/configurations
// @API GeminiDB GET /v3/{project_id}/instances/{instance_id}/configurations
func ResourceGeminiDBInstanceParameter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeminiDBInstanceParameterCreate,
		ReadContext:   resourceGeminiDBInstanceParameterRead,
		UpdateContext: resourceGeminiDBInstanceParameterUpdate,
		DeleteContext: resourceGeminiDBInstanceParameterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGeminiDBInstanceParameterImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"restart_required": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"readonly": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"value_range": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGeminiDBInstanceParameterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceID := d.Get("instance_id").(string)
	name := d.Get("name").(string)

	if err := updateGeminiDBInstanceParameters(ctx, instanceID, d, meta); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceID, name))

	return resourceGeminiDBInstanceParameterRead(ctx, d, meta)
}

func resourceGeminiDBInstanceParameterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	name := d.Get("name").(string)

	httpUrl := "v3/{project_id}/instances/{instance_id}/configurations"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GeminiDB instance parameters: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// Find the parameter by name in configuration_parameters
	paramPath := fmt.Sprintf("configuration_parameters[?name=='%s']|[0]", name)
	paramObj := utils.PathSearch(paramPath, getRespBody, nil)
	if paramObj == nil {
		return diag.Errorf("parameter %s not found in instance %s", name, instanceID)
	}

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("name", utils.PathSearch("name", paramObj, nil)),
		d.Set("value", utils.PathSearch("value", paramObj, nil)),
		d.Set("restart_required", utils.PathSearch("restart_required", paramObj, nil)),
		d.Set("readonly", utils.PathSearch("readonly", paramObj, nil)),
		d.Set("value_range", utils.PathSearch("value_range", paramObj, nil)),
		d.Set("type", utils.PathSearch("type", paramObj, nil)),
		d.Set("description", utils.PathSearch("description", paramObj, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGeminiDBInstanceParameterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChange("value") {
		return resourceGeminiDBInstanceParameterRead(ctx, d, meta)
	}

	if err := updateGeminiDBInstanceParameters(ctx, d.Get("instance_id").(string), d, meta); err != nil {
		return diag.FromErr(err)
	}

	return resourceGeminiDBInstanceParameterRead(ctx, d, meta)
}

func resourceGeminiDBInstanceParameterDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting an instance parameter resource is not supported. The parameter value will remain as set. " +
		"To revert the parameter to its default value, please use the console or the reset configuration API."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Resource deletion is not supported",
			Detail:   errorMsg,
		},
	}
}

func resourceGeminiDBInstanceParameterImportState(ctx context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<name>")
	}

	instanceID := parts[0]
	name := parts[1]

	mErr := multierror.Append(nil,
		d.Set("instance_id", instanceID),
		d.Set("name", name),
	)

	d.SetId(fmt.Sprintf("%s/%s", instanceID, name))

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

func updateGeminiDBInstanceParameters(ctx context.Context, instanceID string, d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return fmt.Errorf("error creating GeminiDB client: %s", err)
	}

	name := d.Get("name").(string)
	value := d.Get("value").(string)

	httpUrl := "v3.1/{project_id}/instances/{instance_id}/configurations"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"values": map[string]interface{}{
				name: value,
			},
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	res, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     geminiDbInstanceStatusRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error modifying GeminiDB instance parameter %s: %s", name, err)
	}
	createRespBody, err := utils.FlattenResponse(res.(*http.Response))
	if err != nil {
		return err
	}

	d.SetId(d.Get("instance_id").(string) + "/" + d.Get("name").(string))

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("error creating GeminiDB account: job_id is not found in API response")
	}
	err = checkGeminiDbJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	return nil
}

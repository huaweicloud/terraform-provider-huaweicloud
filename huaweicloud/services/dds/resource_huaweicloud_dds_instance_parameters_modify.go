package dds

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDS PUT /v3/{project_id}/instances/{instance_id}/configurations
// @API DDS GET /v3/{project_id}/instances/{instance_id}/configurations
// @API DDS GET /v3/{project_id}/instances
// @API DDS GET /v3/{project_id}/jobs
func ResourceDDSInstanceParametersModify() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDDSInstanceParametersModifyCreate,
		ReadContext:   resourceDDSInstanceParametersModifyRead,
		UpdateContext: resourceDDSInstanceParametersModifyUpdate,
		DeleteContext: resourceDDSInstanceParametersModifyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
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
			"parameters": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"entity_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDDSInstanceParametersModifyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DDS client: %s ", err)
	}
	instId := d.Get("instance_id").(string)
	entityId := d.Get("entity_id").(string)
	if entityId == "" {
		entityId = instId
	}
	opt := buildDDSInstanceParametersModifyRequestBody(d.Get("parameters").(*schema.Set).List(), entityId)

	// modify parameters
	if ctx, err = modifyInstanceParameters(ctx, client, d.Timeout(schema.TimeoutCreate), instId, opt); err != nil {
		return diag.FromErr(err)
	}

	if instId == entityId {
		d.SetId(instId)
	} else {
		d.SetId(instId + "/" + entityId)
	}

	return resourceDDSInstanceParametersModifyRead(ctx, d, meta)
}

func buildDDSInstanceParametersModifyRequestBody(params []interface{}, entityID string) map[string]interface{} {
	values := make(map[string]string, len(params))
	for _, v := range params {
		key := v.(map[string]interface{})["name"].(string)
		value := v.(map[string]interface{})["value"].(string)
		values[key] = value
	}

	return map[string]interface{}{
		"entity_id":        entityID,
		"parameter_values": values,
	}
}

func resourceDDSInstanceParametersModifyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DdsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s ", err)
	}
	instId := d.Get("instance_id").(string)
	entityId := d.Get("entity_id").(string)
	if entityId == "" {
		entityId = instId
	}
	getParametersHttpUrl := "v3/{project_id}/instances/{instance_id}/configurations?entity_id={entity_id}"
	getParametersPath := client.Endpoint + getParametersHttpUrl
	getParametersPath = strings.ReplaceAll(getParametersPath, "{project_id}", client.ProjectID)
	getParametersPath = strings.ReplaceAll(getParametersPath, "{instance_id}", instId)
	getParametersPath = strings.ReplaceAll(getParametersPath, "{entity_id}", entityId)
	getParametersOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getParametersResp, err := client.Request("GET", getParametersPath, &getParametersOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.216031"),
			"error retrieving DDS parameter")
	}

	getParametersRespBody, err := utils.FlattenResponse(getParametersResp)
	if err != nil {
		return diag.Errorf("error flatten response: %s", err)
	}

	// make params list
	params := d.Get("parameters").(*schema.Set).List()
	parameters := make([]map[string]interface{}, len(params))
	for i, param := range params {
		name := param.(map[string]interface{})["name"].(string)
		jsonPaths := fmt.Sprintf("parameters[?name=='%s']|[0].value", name)
		value := utils.PathSearch(jsonPaths, getParametersRespBody, "")
		if value.(string) == "" {
			return diag.Errorf("error getting param(%s) value: %s", name, err)
		}
		parameters[i] = map[string]interface{}{
			"name":  name,
			"value": value,
		}
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("parameters", parameters),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting fields: %s", err)
	}

	if ctx.Value(ctxType("restartRequiredForParametersChanged")) == "true" {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Restart Required",
				Detail:   "Restart required because some parameters take effect after restarting the instance.",
			},
		}
	}

	return nil
}

func resourceDDSInstanceParametersModifyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DDS client: %s ", err)
	}
	instId := d.Get("instance_id").(string)
	entityId := d.Get("entity_id").(string)
	if entityId == "" {
		entityId = instId
	}

	// modify parameters
	if d.HasChange("parameters") {
		opt := buildDDSInstanceParametersModifyRequestBody(d.Get("parameters").(*schema.Set).List(), entityId)
		if ctx, err = modifyInstanceParameters(ctx, client, d.Timeout(schema.TimeoutUpdate), instId, opt); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDDSInstanceParametersModifyRead(ctx, d, meta)
}

func resourceDDSInstanceParametersModifyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting parameters modify resource is not supported. The parameters modify resource is only removed " +
		"from the state, the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func modifyInstanceParameters(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration, instID string,
	opts map[string]interface{}) (context.Context, error) {
	modifyHttpUrl := "v3/{project_id}/instances/{instance_id}/configurations"
	modifyPath := client.Endpoint + modifyHttpUrl
	modifyPath = strings.ReplaceAll(modifyPath, "{project_id}", client.ProjectID)
	modifyPath = strings.ReplaceAll(modifyPath, "{instance_id}", instID)
	modifyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         opts,
	}

	// retry, the job_id in return is useless
	retryFunc := func() (interface{}, bool, error) {
		resp, err := client.Request("PUT", modifyPath, &modifyOpt)
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	modifyResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instID),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return ctx, fmt.Errorf("error modifying instance(%s) params: %s", instID, err)
	}

	// wait for job complete
	err = waitForInstanceReady(ctx, client, instID, timeout)
	if err != nil {
		return ctx, err
	}

	// get restart_required
	modifyRespBody, err := utils.FlattenResponse(modifyResp.(*http.Response))
	if err != nil {
		return ctx, fmt.Errorf("error flatten response: %s", err)
	}
	restartRequired := utils.PathSearch("restart_required", modifyRespBody, false)
	if restartRequired.(bool) {
		// Sending restartRequiredForParametersChanged to Read to warn users the instance needs a reboot.
		ctx = context.WithValue(ctx, ctxType("restartRequiredForParametersChanged"), "true")
	}

	return ctx, nil
}

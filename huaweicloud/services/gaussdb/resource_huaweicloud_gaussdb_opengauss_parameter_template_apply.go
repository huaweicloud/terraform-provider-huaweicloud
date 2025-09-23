package gaussdb

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB PUT /v3/{project_id}/configurations/{config_id}/apply
// @API GaussDB GET /v3/{project_id}/instances/{instance_id}
// @API GaussDB GET /v3/{project_id}/jobs
func ResourceOpenGaussParameterTemplateApply() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenGaussParameterTemplateApplyCreate,
		ReadContext:   resourceOpenGaussParameterTemplateApplyRead,
		DeleteContext: resourceOpenGaussParameterTemplateApplyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"config_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOpenGaussParameterTemplateApplyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/configurations/{config_id}/apply"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	configId := d.Get("config_id").(string)
	instanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{config_id}", configId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateOpenGaussParameterTemplateApplyBodyParams(d))
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating GaussDB OpenGauss parameter template apply: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", configId, instanceId))

	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error applying parameter template(%s) to GaussDB OpenGauss instance(%s), job_id is"+
			" not found in the response", configId, instanceId)
	}
	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 2, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildCreateOpenGaussParameterTemplateApplyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance_ids": []string{d.Get("instance_id").(string)},
	}
	return bodyParams
}

func resourceOpenGaussParameterTemplateApplyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOpenGaussParameterTemplateApplyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting parameter template apply resource is not supported. The resource is only removed from the" +
		"state, the GaussDB OpenGauss instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

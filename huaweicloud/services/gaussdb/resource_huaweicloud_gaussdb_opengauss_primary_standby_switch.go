package gaussdb

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/switch-shard
// @API GaussDB GET /v3/{project_id}/instances
// @API GaussDB GET /v3/{project_id}/jobs/{job_id}
func ResourceOpenGaussPrimaryStandbySwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenGaussPrimaryStandbySwitchCreate,
		ReadContext:   resourceOpenGaussPrimaryStandbySwitchRead,
		DeleteContext: resourceOpenGaussPrimaryStandbySwitchDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
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
			"shards": {
				Type:     schema.TypeList,
				Elem:     primaryStandbySwitchShardsSchema(),
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func primaryStandbySwitchShardsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
	return &sc
}

func resourceOpenGaussPrimaryStandbySwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/switch-shard"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreatePrimaryStandbySwitchBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating GaussDB OpenGauss primary standby switch: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId)

	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error creating GaussDB OpenGauss primary standby switch: job_id is not found in API response")
	}

	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 5, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildCreatePrimaryStandbySwitchBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"shards": buildCreatePrimaryStandbySwitchShardsChildBody(d),
	}
	return bodyParams
}

func buildCreatePrimaryStandbySwitchShardsChildBody(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("shards").([]interface{})
	bodyParams := make([]map[string]interface{}, 0)
	for _, v := range rawParams {
		raw := v.(map[string]interface{})
		bodyParams = append(bodyParams, map[string]interface{}{
			"node_id":      raw["node_id"],
			"component_id": raw["component_id"],
		})
	}
	return bodyParams
}

func resourceOpenGaussPrimaryStandbySwitchRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOpenGaussPrimaryStandbySwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting primary standby switch resource is not supported. The resource is only removed from the state," +
		" the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

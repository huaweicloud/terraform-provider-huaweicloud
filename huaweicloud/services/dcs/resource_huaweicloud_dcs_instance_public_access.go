package dcs

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var instancePublicAccessNonUpdatableParams = []string{"instance_id", "publicip_id", "enable_ssl", "elb_id"}

// @API DCS POST /v2/{project_id}/instances/{instance_id}/public-ip
// @API DCS GET /v2/{project_id}/instances/{instance_id}
// @API DCS DELETE /v2/{project_id}/instances/{instance_id}/public-ip
func ResourceDcsInstancePublicAccess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsInstancePublicAccessCreate,
		ReadContext:   resourceDcsInstancePublicAccessRead,
		UpdateContext: resourceDcsInstancePublicAccessUpdate,
		DeleteContext: resourceDcsInstancePublicAccessDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(instancePublicAccessNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
			},
			"publicip_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"elb_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"eip_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"elb_listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceDcsInstancePublicAccessCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/public-ip"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	createOpt.JSONBody = utils.RemoveNil(buildInstancePublicAccessBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		r, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleOperationError(err)
		return r, retry, err
	}
	createResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshDcsInstanceState(client, instanceId),
		WaitTarget:   []string{"RUNNING"},
		WaitPending:  []string{"PENDING"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error enabling instance(%s) public access: %s", instanceId, err)
	}

	d.SetId(instanceId)

	createRespBody, err := utils.FlattenResponse(createResp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error enabling instance(%s) public access: job_id is not found in API response", instanceId)
	}

	err = checkDcsInstanceJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDcsInstancePublicAccessRead(ctx, d, meta)
}

func buildInstancePublicAccessBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"publicip_id": utils.ValueIgnoreEmpty(d.Get("publicip_id")),
		"elb_id":      utils.ValueIgnoreEmpty(d.Get("elb_id")),
	}
	if d.Get("enable_ssl").(bool) {
		bodyParams["enable_ssl"] = d.Get("enable_ssl")
	}
	return bodyParams
}

func resourceDcsInstancePublicAccessRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instance, err := getDcsInstanceByID(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting DCS instance")
	}

	publicInfo := utils.PathSearch("publicip_info", instance, nil)
	if publicInfo == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error getting DCS instance public info")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", d.Id()),
		d.Set("elb_id", utils.PathSearch("elb_id", publicInfo, nil)),
		d.Set("eip_id", utils.PathSearch("eip_id", publicInfo, nil)),
		d.Set("eip_address", utils.PathSearch("eip_address", publicInfo, nil)),
		d.Set("elb_listeners", flattenPublicAccessElbListeners(publicInfo)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPublicAccessElbListeners(resp interface{}) []interface{} {
	curJson := utils.PathSearch("elb_listeners", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":   utils.PathSearch("id", v, nil),
			"port": utils.PathSearch("port", v, nil),
			"name": utils.PathSearch("name", v, nil),
		})
	}
	return rst
}

func resourceDcsInstancePublicAccessUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsInstancePublicAccessDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/public-ip"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS Client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	retryFunc := func() (interface{}, bool, error) {
		r, err := client.Request("DELETE", deletePath, &deleteOpt)
		retry, err := handleOperationError(err)
		return r, retry, err
	}
	deleteResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshDcsInstanceState(client, instanceId),
		WaitTarget:   []string{"RUNNING"},
		WaitPending:  []string{"PENDING"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error disabling instance(%s) public access: %s", instanceId, err)
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error disabling instance(%s) public access: job_id is not found in API response", instanceId)
	}

	err = checkDcsInstanceJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

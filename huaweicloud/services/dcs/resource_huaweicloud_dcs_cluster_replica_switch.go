package dcs

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var clusterReplicaSwitchNonUpdatableParams = []string{"instance_id", "group_id", "node_id"}

// @API DCS POST /v2/{project_id}/instance/{instance_id}/groups/{group_id}/replications/{node_id}/async-switchover
// @API DCS GET /v2/{project_id}/instances/{instance_id}
// @API DCS GET /v2/{project_id}/jobs/{job_id}
func ResourceDcsClusterReplicaSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsClusterReplicaSwitchCreate,
		ReadContext:   resourceDcsClusterReplicaSwitchRead,
		UpdateContext: resourceDcsClusterReplicaSwitchUpdate,
		DeleteContext: resourceDcsClusterReplicaSwitchDelete,

		CustomizeDiff: config.FlexibleForceNew(clusterReplicaSwitchNonUpdatableParams),

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
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_id": {
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

func resourceDcsClusterReplicaSwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/instance/{instance_id}/groups/{group_id}/replications/{node_id}/async-switchover"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	groupId := d.Get("group_id").(string)
	nodeId := d.Get("node_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{group_id}", groupId)
	createPath = strings.ReplaceAll(createPath, "{node_id}", nodeId)

	createOpt := golangsdk.RequestOpts{KeepResponseBody: true}

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
		return diag.Errorf("error switching instance(%s) cluster: %s", instanceId, err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", instanceId, groupId, nodeId))

	createRespBody, err := utils.FlattenResponse(createResp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error switching instance(%s) cluster replica: job_id is not found in API response", instanceId)
	}

	err = checkDcsInstanceJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDcsClusterReplicaSwitchRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsClusterReplicaSwitchUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsClusterReplicaSwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DCS cluster replica switch resource is not supported. The resource is only removed from the" +
		"state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

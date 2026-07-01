package dcs

import (
	"context"
	"errors"
	"fmt"
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

var dcsNodePriorityConfigNonUpdatableParams = []string{"instance_id", "group_id", "node_id"}

// @API DCS GET /v2/{project_id}/instances/{instance_id}/logical-nodes
// @API DCS POST /v2/{project_id}/instances/{instance_id}/groups/{group_id}/replications/{node_id}/slave-priority
func ResourceDcsNodePriorityConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsNodePriorityConfigCreate,
		UpdateContext: resourceDcsNodePriorityConfigUpdate,
		ReadContext:   resourceDcsNodePriorityConfigRead,
		DeleteContext: resourceDcsNodePriorityConfigDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDcsNodePriorityConfigImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(dcsNodePriorityConfigNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
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
			"slave_priority_weight": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"logical_node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"priority_weight": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_remove_ip": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"replication_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func resourceDcsNodePriorityConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/groups/{group_id}/replications/{node_id}/slave-priority"
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

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		OkCodes: []int{
			204,
		},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateDcsNodePriorityConfigBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		_, err = client.Request("POST", createPath, &createOpt)
		retry, err := handleOperationError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshDcsInstanceState(client, instanceId),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating DCS node priority config: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", instanceId, groupId, nodeId))

	return resourceDcsNodePriorityConfigRead(ctx, d, meta)
}

func buildCreateDcsNodePriorityConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"slave_priority_weight": d.Get("slave_priority_weight"),
	}
	return bodyParams
}

func resourceDcsNodePriorityConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/logical-nodes"
		product = "dcs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	nodeId := d.Get("node_id").(string)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DCS instance node priority config")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	node := utils.PathSearch(fmt.Sprintf("nodes[?node_id=='%s']|[0]", nodeId), getRespBody, nil)
	if node == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving DCS instance node priority config")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("group_id", utils.PathSearch("group_id", node, nil)),
		d.Set("node_id", utils.PathSearch("node_id", node, nil)),
		d.Set("slave_priority_weight", utils.PathSearch("priority_weight", node, nil)),
		d.Set("logical_node_id", utils.PathSearch("logical_node_id", node, nil)),
		d.Set("name", utils.PathSearch("name", node, nil)),
		d.Set("status", utils.PathSearch("status", node, nil)),
		d.Set("az_code", utils.PathSearch("az_code", node, nil)),
		d.Set("node_role", utils.PathSearch("node_role", node, nil)),
		d.Set("node_type", utils.PathSearch("node_type", node, nil)),
		d.Set("node_ip", utils.PathSearch("node_ip", node, nil)),
		d.Set("node_port", utils.PathSearch("node_port", node, nil)),
		d.Set("priority_weight", utils.PathSearch("priority_weight", node, nil)),
		d.Set("is_access", utils.PathSearch("is_access", node, nil)),
		d.Set("group_name", utils.PathSearch("group_name", node, nil)),
		d.Set("is_remove_ip", utils.PathSearch("is_remove_ip", node, nil)),
		d.Set("replication_id", utils.PathSearch("replication_id", node, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDcsNodePriorityConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/groups/{group_id}/replications/{node_id}/slave-priority"
		product = "dcs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	if d.HasChange("slave_priority_weight") {
		instanceId := d.Get("instance_id").(string)
		groupId := d.Get("group_id").(string)
		nodeId := d.Get("node_id").(string)

		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceId)
		updatePath = strings.ReplaceAll(updatePath, "{group_id}", groupId)
		updatePath = strings.ReplaceAll(updatePath, "{node_id}", nodeId)

		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
			OkCodes: []int{
				204,
			},
		}
		updateOpt.JSONBody = buildUpdateDcsNodePriorityConfigBodyParams(d)

		retryFunc := func() (interface{}, bool, error) {
			_, err = client.Request("POST", updatePath, &updateOpt)
			retry, err := handleOperationError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     refreshDcsInstanceState(client, instanceId),
			WaitTarget:   []string{"RUNNING"},
			Timeout:      d.Timeout(schema.TimeoutCreate),
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.Errorf("error updating DCS node priority config: %s", err)
		}
	}

	return resourceDcsNodePriorityConfigRead(ctx, d, meta)
}

func buildUpdateDcsNodePriorityConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"slave_priority_weight": d.Get("slave_priority_weight"),
	}
	return bodyParams
}

func resourceDcsNodePriorityConfigDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DCS node priority config resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourceDcsNodePriorityConfigImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 3 {
		return nil, errors.New("invalid format specified for import ID, must be <instance_id>/<group_id>/<node_id>")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("group_id", parts[1]),
		d.Set("node_id", parts[2]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

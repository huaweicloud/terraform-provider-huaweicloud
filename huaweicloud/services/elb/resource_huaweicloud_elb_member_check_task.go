package elb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var memberCheckTaskNonUpdatableParams = []string{"member_id", "listener_id", "subject"}

// @API ELB POST /v3/{project_id}/elb/members/{member_id}/health-check
// @API ELB GET /v3/{project_id}/elb/members/check/jobs/{job_id}
func ResourceMemberCheckTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMemberCheckTaskCreate,
		ReadContext:   resourceMemberCheckTaskRead,
		UpdateContext: resourceMemberCheckTaskUpdate,
		DeleteContext: resourceMemberCheckTaskDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(memberCheckTaskNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"member_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subject": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"result": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     memberCheckTaskResultSchema(),
			},
			"check_item_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"check_item_finished_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func memberCheckTaskResultSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     memberCheckTaskResultGroupSchema(),
			},
			"acl": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     memberCheckTaskResultGroupSchema(),
			},
			"security_group": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     memberCheckTaskResultGroupSchema(),
			},
		},
	}
	return &sc
}

func memberCheckTaskResultGroupSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"check_result": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"check_items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     memberCheckTaskResultGroupCheckItemsSchema(),
			},
			"check_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func memberCheckTaskResultGroupCheckItemsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reason_template": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reason_params": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return &sc
}

func resourceMemberCheckTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/members/{member_id}/health-check"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{member_id}", d.Get("member_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateMemberCheckTaskBodyParams(d))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ELB member check task: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}
	jobId := utils.PathSearch("member_check.job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating ELB member check task: job_id is not found in API response")
	}

	d.SetId(jobId)

	err = waitForMemberCheckTaskComplete(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error creating ELB member check task: %s", err)
	}

	return resourceMemberCheckTaskRead(ctx, d, meta)
}

func buildCreateMemberCheckTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"listener_id": d.Get("listener_id"),
		"subject":     d.Get("subject"),
	}
	bodyParams := map[string]interface{}{
		"member_check": params,
	}
	return bodyParams
}

func waitForMemberCheckTaskComplete(ctx context.Context, client *golangsdk.ServiceClient, id string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{"processed"},
		Pending:      []string{"processing"},
		Refresh:      resourceMemberCheckTaskRefreshFunc(client, id),
		Timeout:      timeout,
		PollInterval: 2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for ELB member check task (%s) to be completed: %s ", id, err)
	}
	return nil
}

func resourceMemberCheckTaskRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getRespBody, err := getMemberCheckTask(client, id)
		if err != nil {
			return nil, "failed", err
		}

		status := utils.PathSearch("member_check.status", getRespBody, "")
		return getRespBody, status.(string), nil
	}
}

func resourceMemberCheckTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	getRespBody, err := getMemberCheckTask(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ELB member check task")
	}

	result, subject := flattenMemberCheckTaskResult(getRespBody)
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("member_id", utils.PathSearch("member_check.member_id", getRespBody, nil)),
		d.Set("listener_id", utils.PathSearch("member_check.listener_id", getRespBody, nil)),
		d.Set("subject", subject),
		d.Set("status", utils.PathSearch("member_check.status", getRespBody, nil)),
		d.Set("result", result),
		d.Set("check_item_total_num", utils.PathSearch("member_check.check_item_total_num", getRespBody, nil)),
		d.Set("check_item_finished_num", utils.PathSearch("member_check.check_item_finished_num", getRespBody, nil)),
		d.Set("created_at", utils.PathSearch("member_check.created_at", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("member_check.updated_at", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getMemberCheckTask(client *golangsdk.ServiceClient, jobID string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/elb/members/check/jobs/{job_id}"
	)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_id}", jobID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func flattenMemberCheckTaskResult(resp interface{}) ([]interface{}, string) {
	curJson := utils.PathSearch("member_check.result", resp, nil)
	if curJson == nil {
		return nil, ""
	}

	configs, configSubject := flattenMemberCheckTaskResultGroup(curJson, "config")
	acls, aclSubject := flattenMemberCheckTaskResultGroup(curJson, "acl")
	securityGroups, securityGroupSubject := flattenMemberCheckTaskResultGroup(curJson, "security_group")
	rst := []interface{}{
		map[string]interface{}{
			"config":         configs,
			"acl":            acls,
			"security_group": securityGroups,
		},
	}
	subject := ""
	if configSubject != "" && aclSubject != "" && securityGroupSubject != "" {
		subject = "all"
	} else {
		subject = fmt.Sprintf("%s%s%s", configSubject, aclSubject, securityGroupSubject)
	}
	return rst, subject
}

func flattenMemberCheckTaskResultGroup(resp interface{}, paramName string) ([]interface{}, string) {
	curJson := utils.PathSearch(paramName, resp, nil)
	if curJson == nil {
		return nil, ""
	}

	checkItems, subject := flattenMemberCheckTaskResultGroupCheckItems(curJson)
	rst := []interface{}{
		map[string]interface{}{
			"check_result": utils.PathSearch("check_result", curJson, nil),
			"check_items":  checkItems,
			"check_status": utils.PathSearch("check_status", curJson, nil),
		},
	}
	return rst, subject
}

func flattenMemberCheckTaskResultGroupCheckItems(resp interface{}) ([]interface{}, string) {
	curJson := utils.PathSearch("check_items", resp, nil)
	if curJson == nil {
		return nil, ""
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	subject := ""
	for _, v := range curArray {
		subject = utils.PathSearch("subject", v, "").(string)
		rst = append(rst, map[string]interface{}{
			"name":            utils.PathSearch("name", v, nil),
			"reason":          utils.PathSearch("reason", v, nil),
			"severity":        utils.PathSearch("severity", v, nil),
			"subject":         subject,
			"job_id":          utils.PathSearch("job_id", v, nil),
			"reason_template": utils.PathSearch("reason_template", v, nil),
			"reason_params":   utils.PathSearch("reason_params", v, nil),
		})
	}
	return rst, subject
}

func resourceMemberCheckTaskUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceMemberCheckTaskDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting ELB member check task resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

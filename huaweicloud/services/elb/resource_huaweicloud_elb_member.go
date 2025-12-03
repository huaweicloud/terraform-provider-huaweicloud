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

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB POST /v3/{project_id}/elb/pools/{pool_id}/members
// @API ELB GET /v3/{project_id}/elb/pools/{pool_id}/members/{member_id}
// @API ELB PUT /v3/{project_id}/elb/pools/{pool_id}/members/{member_id}
// @API ELB DELETE /v3/{project_id}/elb/pools/{pool_id}/members/{member_id}
func ResourceMemberV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMemberV3Create,
		ReadContext:   resourceMemberV3Read,
		UpdateContext: resourceMemberV3Update,
		DeleteContext: resourceMemberV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceELBMemberImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"member_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reason": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     memberReasonRefSchema(),
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     memberStatusRefSchema(),
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

func memberReasonRefSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"expected_response": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"healthcheck_response": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reason_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func memberStatusRefSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reason": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     memberReasonRefSchema(),
			},
		},
	}
}

func resourceMemberV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/pools/{pool_id}/members"
		product = "elb"
	)
	elbClient, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	createPath := elbClient.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", elbClient.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{pool_id}", d.Get("pool_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateMemberBodyParams(d))
	createResp, err := elbClient.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ELB member: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error retrieving ELB member: %s", err)
	}
	memberId := utils.PathSearch("member.id", createRespBody, "").(string)
	if memberId == "" {
		return diag.Errorf("error creating ELB member: ID is not found in API response")
	}

	d.SetId(memberId)

	err = waitForMemberReady(ctx, elbClient, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error creating ELB member: %s", err)
	}

	return resourceMemberV3Read(ctx, d, meta)
}

func buildCreateMemberBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":           utils.ValueIgnoreEmpty(d.Get("name")),
		"address":        d.Get("address"),
		"protocol_port":  utils.ValueIgnoreEmpty(d.Get("protocol_port")),
		"weight":         utils.ValueIgnoreEmpty(d.Get("weight")),
		"subnet_cidr_id": utils.ValueIgnoreEmpty(d.Get("subnet_id")),
	}
	return map[string]interface{}{"member": bodyParams}
}

func waitForMemberReady(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{"Ready"},
		Pending:      []string{"Pending"},
		Refresh:      resourceMemberRefreshFunc(d, client),
		Timeout:      timeout,
		PollInterval: 2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for ELB member (%s) ready: %s ", d.Id(), err)
	}
	return nil
}

func resourceMemberRefreshFunc(d *schema.ResourceData, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getRespBody, err := getMember(d, client)
		if err != nil {
			return nil, "Failed", err
		}

		status := utils.PathSearch("member.operating_status", getRespBody, "").(string)
		if utils.StrSliceContains([]string{"NO_MONITOR", "ONLINE", "OFFLINE", "UNKNOWN"}, status) {
			return getRespBody, "Ready", nil
		}
		if status == "INITIAL" {
			return getRespBody, "Pending", nil
		}
		return getRespBody, status, nil
	}
}

func resourceMemberV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		product = "elb"
	)
	elbClient, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	getRespBody, err := getMember(d, elbClient)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ELB member")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("member.name", getRespBody, nil)),
		d.Set("weight", utils.PathSearch("member.weight", getRespBody, nil)),
		d.Set("subnet_id", utils.PathSearch("member.subnet_cidr_id", getRespBody, nil)),
		d.Set("address", utils.PathSearch("member.address", getRespBody, nil)),
		d.Set("protocol_port", utils.PathSearch("member.protocol_port", getRespBody, nil)),
		d.Set("instance_id", utils.PathSearch("member.instance_id", getRespBody, nil)),
		d.Set("ip_version", utils.PathSearch("member.ip_version", getRespBody, nil)),
		d.Set("member_type", utils.PathSearch("member.member_type", getRespBody, nil)),
		d.Set("operating_status", utils.PathSearch("member.operating_status", getRespBody, nil)),
		d.Set("reason", flattenMemberReason(utils.PathSearch("member.reason", getRespBody, nil))),
		d.Set("status", flattenMemberStatus(utils.PathSearch("member.status", getRespBody, nil))),
		d.Set("created_at", utils.PathSearch("member.created_at", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("member.updated_at", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getMember(d *schema.ResourceData, client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/elb/pools/{pool_id}/members/{member_id}"
	)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{pool_id}", d.Get("pool_id").(string))
	getPath = strings.ReplaceAll(getPath, "{member_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func flattenMemberReason(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"expected_response":    utils.PathSearch("expected_response", resp, nil),
			"healthcheck_response": utils.PathSearch("healthcheck_response", resp, nil),
			"reason_code":          utils.PathSearch("reason_code", resp, nil),
		},
	}
	return rst
}

func flattenMemberStatus(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}
	curArray := resp.([]interface{})
	if len(curArray) < 1 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"listener_id":      utils.PathSearch("listener_id", v, nil),
			"operating_status": utils.PathSearch("operating_status", v, nil),
			"reason":           flattenMemberReason(v),
		})
	}
	return rst
}

func resourceMemberV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/pools/{pool_id}/members/{member_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{pool_id}", d.Get("pool_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{member_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateMemberBodyParams(d))
	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating ELB member: %s", err)
	}

	return resourceMemberV3Read(ctx, d, meta)
}

func buildUpdateMemberBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":          utils.ValueIgnoreEmpty(d.Get("name")),
		"protocol_port": utils.ValueIgnoreEmpty(d.Get("protocol_port")),
		"weight":        utils.ValueIgnoreEmpty(d.Get("weight")),
	}
	return map[string]interface{}{"member": bodyParams}
}

func resourceMemberV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/pools/{pool_id}/members/{member_id}"
		product = "elb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{pool_id}", d.Get("pool_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{member_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting ELB member")
	}
	return nil
}

func resourceELBMemberImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format specified for member. Format must be <pool_id>/<member_id>")
		return nil, err
	}

	poolID := parts[0]
	memberID := parts[1]

	d.SetId(memberID)
	d.Set("pool_id", poolID)

	return []*schema.ResourceData{d}, nil
}

// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DDS
// ---------------------------------------------------------------

package dds

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
	"github.com/chnsz/golangsdk/openstack/dds/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDS POST /v3/{project_id}/instances/{instance_id}/auditlog-policy
// @API DDS GET /v3/{project_id}/instances/{instance_id}/auditlog-policy
// @API DDS GET /v3/{project_id}/instances
func ResourceDdsAuditLogPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdsAuditLogPolicyCreate,
		UpdateContext: resourceDdsAuditLogPolicyUpdate,
		ReadContext:   resourceDdsAuditLogPolicyRead,
		DeleteContext: resourceDdsAuditLogPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the DDS instance.`,
			},
			"keep_days": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  `Specifies the number of days for storing audit logs.`,
				ValidateFunc: validation.IntBetween(7, 732),
			},
			"audit_scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the audit scope.`,
			},
			"audit_types": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Specifies the audit type.`,
			},
			"reserve_auditlogs": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether the historical audit logs are retained when SQL audit is disabled.`,
			},
		},
	}
}

func resourceDdsAuditLogPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	instanceId := d.Get("instance_id")
	err := setAuditLogPolicy(ctx, cfg, d, instanceId.(string), "creating")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instanceId.(string))

	return resourceDdsAuditLogPolicyRead(ctx, d, meta)
}

func resourceDdsAuditLogPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	updateAuditLogPolicyHasChanges := []string{
		"keep_days",
		"audit_scope",
		"audit_types",
	}

	if d.HasChanges(updateAuditLogPolicyHasChanges...) {
		err := setAuditLogPolicy(ctx, cfg, d, d.Id(), "updating")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDdsAuditLogPolicyRead(ctx, d, meta)
}

func setAuditLogPolicy(ctx context.Context, cfg *config.Config, d *schema.ResourceData, instanceId, operateMethod string) error {
	region := cfg.GetRegion(d)
	var (
		setAuditLogPolicyHttpUrl = "v3/{project_id}/instances/{instance_id}/auditlog-policy"
		setAuditLogPolicyProduct = "dds"
	)
	setAuditLogPolicyClient, err := cfg.NewServiceClient(setAuditLogPolicyProduct, region)
	if err != nil {
		return fmt.Errorf("error creating DDS Client: %s", err)
	}

	setAuditLogPolicyPath := setAuditLogPolicyClient.Endpoint + setAuditLogPolicyHttpUrl
	setAuditLogPolicyPath = strings.ReplaceAll(setAuditLogPolicyPath, "{project_id}",
		setAuditLogPolicyClient.ProjectID)
	setAuditLogPolicyPath = strings.ReplaceAll(setAuditLogPolicyPath, "{instance_id}", instanceId)

	setAuditLogPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	setAuditLogPolicyOpt.JSONBody = utils.RemoveNil(buildSetAuditLogPolicyBodyParams(d))
	_, err = setAuditLogPolicyClient.Request("POST", setAuditLogPolicyPath, &setAuditLogPolicyOpt)
	if err != nil {
		return fmt.Errorf("error %s DDS audit log policy: %s", operateMethod, err)
	}

	timeout := schema.TimeoutCreate
	if operateMethod == "updating" {
		timeout = schema.TimeoutUpdate
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"success"},
		Refresh:      ddsInstanceAuditLogPolicyRefreshFunc(setAuditLogPolicyClient, instanceId),
		Timeout:      d.Timeout(timeout),
		Delay:        2 * time.Second,
		PollInterval: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for instance (%s) audit log policy to be set: %s ", instanceId, err)
	}

	return nil
}

func ddsInstanceAuditLogPolicyRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := instances.ListInstanceOpts{
			Id: instanceID,
		}
		allPages, err := instances.List(client, &opts).AllPages()
		if err != nil {
			return nil, "", err
		}
		instancesList, err := instances.ExtractInstances(allPages)
		if err != nil {
			return nil, "", err
		}

		if instancesList.TotalCount == 0 {
			return nil, "", fmt.Errorf("the DDS instance has been deleted")
		}
		insts := instancesList.Instances

		actions := insts[0].Actions
		// wait for updating
		if utils.StrSliceContains(actions, "OPS_AUDIT_LOG") {
			return insts[0], "pending", nil
		}
		return insts[0], "success", nil
	}
}

func buildSetAuditLogPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"keep_days":   utils.ValueIgnoreEmpty(d.Get("keep_days")),
		"audit_scope": utils.ValueIgnoreEmpty(d.Get("audit_scope")),
		"audit_types": utils.ValueIgnoreEmpty(d.Get("audit_types")),
	}
	return bodyParams
}

func resourceDdsAuditLogPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAuditLog: Query DDS audit log
	var (
		getAuditLogPolicyHttpUrl = "v3/{project_id}/instances/{instance_id}/auditlog-policy"
		getAuditLogPolicyProduct = "dds"
	)
	getAuditLogPolicyClient, err := cfg.NewServiceClient(getAuditLogPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDS Client: %s", err)
	}

	getAuditLogPolicyPath := getAuditLogPolicyClient.Endpoint + getAuditLogPolicyHttpUrl
	getAuditLogPolicyPath = strings.ReplaceAll(getAuditLogPolicyPath, "{project_id}",
		getAuditLogPolicyClient.ProjectID)
	getAuditLogPolicyPath = strings.ReplaceAll(getAuditLogPolicyPath, "{instance_id}", d.Id())

	getAuditLogPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getAuditLogPolicyResp, err := getAuditLogPolicyClient.Request("GET",
		getAuditLogPolicyPath, &getAuditLogPolicyOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DDS audit log")
	}

	getAuditLogPolicyRespBody, err := utils.FlattenResponse(getAuditLogPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	keepDays := utils.PathSearch("keep_days", getAuditLogPolicyRespBody, 0)
	if keepDays.(float64) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("keep_days", utils.PathSearch("keep_days", getAuditLogPolicyRespBody, nil)),
		d.Set("audit_scope", utils.PathSearch("audit_scope", getAuditLogPolicyRespBody, nil)),
		d.Set("audit_types", utils.PathSearch("audit_types", getAuditLogPolicyRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDdsAuditLogPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteAuditLog: Delete DDS audit log
	var (
		deleteAuditLogPolicyHttpUrl = "v3/{project_id}/instances/{instance_id}/auditlog-policy"
		deleteAuditLogPolicyProduct = "dds"
	)
	deleteAuditLogPolicyClient, err := cfg.NewServiceClient(deleteAuditLogPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDS Client: %s", err)
	}

	deleteAuditLogPolicyPath := deleteAuditLogPolicyClient.Endpoint + deleteAuditLogPolicyHttpUrl
	deleteAuditLogPolicyPath = strings.ReplaceAll(deleteAuditLogPolicyPath, "{project_id}",
		deleteAuditLogPolicyClient.ProjectID)
	deleteAuditLogPolicyPath = strings.ReplaceAll(deleteAuditLogPolicyPath, "{instance_id}", d.Id())

	deleteAuditLogPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	deleteAuditLogPolicyOpt.JSONBody = utils.RemoveNil(buildDeleteAuditLogPolicyBodyParams(d))
	_, err = deleteAuditLogPolicyClient.Request("POST", deleteAuditLogPolicyPath, &deleteAuditLogPolicyOpt)
	if err != nil {
		return diag.Errorf("error deleting DDS audit log policy: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"success"},
		Refresh:      ddsInstanceAuditLogPolicyRefreshFunc(deleteAuditLogPolicyClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        2 * time.Second,
		PollInterval: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) audit log policy to be deleted: %s", d.Id(), err)
	}

	return nil
}

func buildDeleteAuditLogPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"keep_days":         0,
		"reserve_auditlogs": utils.ValueIgnoreEmpty(d.Get("reserve_auditlogs")),
	}
	return bodyParams
}

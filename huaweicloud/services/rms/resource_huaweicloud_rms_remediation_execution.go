package rms

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

var remediationExecutionNonUpdatableParams = []string{"policy_assignment_id", "all_supported", "resource_ids"}

// @API CONFIG POST /v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-execution
// @API CONFIG GET /v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-execution-statuses
func ResourceRemediationExecution() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRemediationExecutionCreate,
		ReadContext:   resourceRemediationExecutionRead,
		UpdateContext: resourceRemediationExecutionUpdate,
		DeleteContext: resourceRemediationExecutionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(remediationExecutionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"policy_assignment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the policy assignment ID.`,
			},
			"all_supported": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies whether to perform remediation for all non-compliant resources.`,
			},
			"resource_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of resource IDs that require remediation.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"result": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The result of the remediation execution.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"invocation_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The start time of remediation.`,
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The execution state of remediation.`,
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The information of remediation execution.`,
						},
						"automatic": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the remediation is automatic.`,
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource ID.`,
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource name.`,
						},
						"resource_provider": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cloud service name.`,
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource type.`,
						},
					},
				},
			},
		},
	}
}

func resourceRemediationExecutionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var (
		createRemediationExecutionHttpUrl = "v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-execution"
		createRemediationExecutionProduct = "rms"
	)
	client, err := conf.NewServiceClient(createRemediationExecutionProduct, region)
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	policyAssignmentId := d.Get("policy_assignment_id").(string)
	createRemediationExecutionPath := client.Endpoint + createRemediationExecutionHttpUrl
	createRemediationExecutionPath = strings.ReplaceAll(createRemediationExecutionPath, "{domain_id}", conf.DomainID)
	createRemediationExecutionPath = strings.ReplaceAll(createRemediationExecutionPath, "{policy_assignment_id}", policyAssignmentId)

	createRemediationExecutionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createRemediationExecutionOpt.JSONBody = utils.RemoveNil(buildCreateRemediationExecutionBodyParams(d))
	_, err = client.Request("POST", createRemediationExecutionPath, &createRemediationExecutionOpt)
	if err != nil {
		return diag.Errorf("error creating RMS remediation execution: %s", err)
	}

	d.SetId(policyAssignmentId)

	err = waitingForRemediationExecutionCompleted(ctx, client, d, conf.DomainID)
	if err != nil {
		return diag.Errorf("error waiting for RMS remediation execution completed: %s", err)
	}

	return resourceRemediationExecutionRead(ctx, d, meta)
}

func buildCreateRemediationExecutionBodyParams(d *schema.ResourceData) map[string]interface{} {
	resourceIds := d.Get("resource_ids").(*schema.Set).List()
	bodyParams := map[string]interface{}{
		"all_supported": d.Get("all_supported"),
		"resource_ids":  utils.ValueIgnoreEmpty(resourceIds),
	}
	return bodyParams
}

func waitingForRemediationExecutionCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, domainId string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			policyAssignmentId := d.Get("policy_assignment_id").(string)
			checkExpression := "[?state=='IN_QUEUE' || state=='IN_PROGRESS']"
			result, err := getRemediationExecution(client, domainId, policyAssignmentId)
			if err != nil {
				return nil, "ERROR", err
			}

			pendingTasks := utils.PathSearch(checkExpression, result, make([]interface{}, 0)).([]interface{})
			if len(pendingTasks) == 0 {
				return result, "COMPLETED", nil
			}

			return result, "PENDING", nil
		},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func resourceRemediationExecutionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	var mErr *multierror.Error

	getRemediationExecutionProduct := "rms"
	client, err := conf.NewServiceClient(getRemediationExecutionProduct, region)
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	domainID := conf.DomainID
	policyAssignmentId := d.Get("policy_assignment_id").(string)
	result, err := getRemediationExecution(client, domainID, policyAssignmentId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RMS remediation execution")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("policy_assignment_id", d.Id()),
		d.Set("result", flattenRemediationExecutionResult(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getRemediationExecution(client *golangsdk.ServiceClient, domainId, policyAssignmentId string) ([]interface{}, error) {
	httpUrl := "v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-execution-statuses"
	getRemediationExecutionPath := client.Endpoint + httpUrl
	getRemediationExecutionPath = strings.ReplaceAll(getRemediationExecutionPath, "{domain_id}", domainId)
	getRemediationExecutionPath = strings.ReplaceAll(getRemediationExecutionPath, "{policy_assignment_id}", policyAssignmentId)
	getRemediationExecutionPath += fmt.Sprintf("?limit=%v", 100)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	path := getRemediationExecutionPath
	rst := make([]interface{}, 0)
	for {
		getRemediationExecutionResp, err := client.Request("GET", path, &opt)
		if err != nil {
			return nil, err
		}
		getRemediationExecutionRespBody, err := utils.FlattenResponse(getRemediationExecutionResp)
		if err != nil {
			return nil, err
		}

		executionInfo := utils.PathSearch("value[*]", getRemediationExecutionRespBody, make([]interface{}, 0))
		rst = append(rst, executionInfo.([]interface{})...)

		marker := utils.PathSearch("page_info.next_marker", getRemediationExecutionRespBody, "")
		if marker == "" {
			break
		}
		path = fmt.Sprintf("%s&marker=%s", getRemediationExecutionPath, marker)
	}
	if len(rst) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return rst, nil
}

func flattenRemediationExecutionResult(params []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(params))
	for _, param := range params {
		rst = append(rst, map[string]interface{}{
			"automatic":         utils.PathSearch("automatic", param, false),
			"resource_id":       utils.PathSearch("resource_id", param, nil),
			"resource_name":     utils.PathSearch("resource_name", param, nil),
			"resource_provider": utils.PathSearch("resource_provider", param, nil),
			"resource_type":     utils.PathSearch("resource_type", param, nil),
			"invocation_time":   utils.PathSearch("invocation_time", param, nil),
			"state":             utils.PathSearch("state", param, nil),
			"message":           utils.PathSearch("message", param, nil),
		})
	}
	return rst
}

func resourceRemediationExecutionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRemediationExecutionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource. Deleting this resource will not change
		the status of the current resource, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

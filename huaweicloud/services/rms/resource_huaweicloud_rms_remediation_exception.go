package rms

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var remediationExceptionNonUpdatableParams = []string{"policy_assignment_id"}

// @API CONFIG GET /v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-exception
// @API CONFIG POST /v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-exception/create
// @API CONFIG POST /v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-exception/delete
func ResourceRemediationException() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRemediationExceptionCreate,
		UpdateContext: resourceRemediationExceptionUpdate,
		ReadContext:   resourceRemediationExceptionRead,
		DeleteContext: resourceRemediationExceptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(remediationExceptionNonUpdatableParams),

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
			"exceptions": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the resource ID.`,
						},
						"message": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the reason for adding an exception.`,
						},
						"joined_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when a remediation exception is added.`,
						},
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of a remediation exception.`,
						},
					},
				},
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

func resourceRemediationExceptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	product := "rms"
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	policyAssignmentID := d.Get("policy_assignment_id").(string)
	exceptions := d.Get("exceptions").(*schema.Set)

	err = addResourceToException(client, cfg.DomainID, policyAssignmentID, exceptions.List())
	if err != nil {
		return diag.Errorf("error creating RMS remediation exception: %s", err)
	}

	d.SetId(policyAssignmentID)
	return resourceRemediationExceptionRead(ctx, d, meta)
}

func addResourceToException(client *golangsdk.ServiceClient, domainID, id string, exceptions []interface{}) error {
	addExceptionHttpUrl := "v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-exception/create"
	addExceptionHttpUrl = strings.ReplaceAll(addExceptionHttpUrl, "{domain_id}", domainID)
	addExceptionHttpUrl = strings.ReplaceAll(addExceptionHttpUrl, "{policy_assignment_id}", id)
	addExceptionPath := client.Endpoint + addExceptionHttpUrl

	addExceptionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	addExceptionOpt.JSONBody = buildExceptionBodyParams(exceptions)
	_, err := client.Request("POST", addExceptionPath, &addExceptionOpt)
	return err
}

func buildExceptionBodyParams(exceptions []interface{}) map[string]interface{} {
	rst := make([]interface{}, len(exceptions))

	for i, exception := range exceptions {
		exceptionMap := exception.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"resource_id": exceptionMap["resource_id"],
			"message":     exceptionMap["message"],
		}
	}
	return map[string]interface{}{
		"exceptions": rst,
	}
}

func resourceRemediationExceptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	product := "rms"
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	policyAssignmentID := d.Get("policy_assignment_id").(string)

	oldExceptions, newExceptions := d.GetChange("exceptions")
	rmExceptions := oldExceptions.(*schema.Set).Difference(newExceptions.(*schema.Set))
	addExceptions := newExceptions.(*schema.Set).Difference(oldExceptions.(*schema.Set))

	if rmExceptions.Len() > 0 {
		err = removeResourceFromException(client, cfg.DomainID, policyAssignmentID, rmExceptions.List())
		if err != nil {
			return diag.Errorf("error updating RMS remediation exception: %s", err)
		}
	}

	if addExceptions.Len() > 0 {
		err = addResourceToException(client, cfg.DomainID, policyAssignmentID, addExceptions.List())
		if err != nil {
			return diag.Errorf("error updating RMS remediation exception: %s", err)
		}
	}

	return resourceRemediationExceptionRead(ctx, d, meta)
}

func removeResourceFromException(client *golangsdk.ServiceClient, domainID, id string, exceptions []interface{}) error {
	removeExceptionHttpUrl := "v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-exception/delete"
	removeExceptionHttpUrl = strings.ReplaceAll(removeExceptionHttpUrl, "{domain_id}", domainID)
	removeExceptionHttpUrl = strings.ReplaceAll(removeExceptionHttpUrl, "{policy_assignment_id}", id)
	removeExceptionPath := client.Endpoint + removeExceptionHttpUrl

	removeExceptionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	removeExceptionOpt.JSONBody = buildExceptionBodyParams(exceptions)
	_, err := client.Request("POST", removeExceptionPath, &removeExceptionOpt)
	return err
}

func resourceRemediationExceptionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	product := "rms"
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}
	policyAssignmentID := d.Id()
	exceptions, err := ListRemediationExceptionInfo(client, cfg.DomainID, policyAssignmentID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RMS remediation exception")
	}

	mErr := multierror.Append(nil,
		d.Set("policy_assignment_id", policyAssignmentID),
		d.Set("exceptions", flattenRemediationExceptions(exceptions)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func ListRemediationExceptionInfo(client *golangsdk.ServiceClient, domainID, id string) ([]interface{}, error) {
	listExceptionHttpUrl := "v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-exception"
	listExceptionHttpUrl = strings.ReplaceAll(listExceptionHttpUrl, "{domain_id}", domainID)
	listExceptionHttpUrl = strings.ReplaceAll(listExceptionHttpUrl, "{policy_assignment_id}", id)

	listExceptionPath := client.Endpoint + listExceptionHttpUrl
	listExceptionPath += fmt.Sprintf("?limit=%v", 100)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	path := listExceptionPath
	rst := make([]interface{}, 0)
	for {
		listRemediationExceptionResp, err := client.Request("GET", path, &opt)
		if err != nil {
			return nil, err
		}
		listRemediationExceptionRespBody, err := utils.FlattenResponse(listRemediationExceptionResp)
		if err != nil {
			return nil, err
		}

		exceptionInfos := utils.PathSearch("value[*]", listRemediationExceptionRespBody, make([]interface{}, 0))
		rst = append(rst, exceptionInfos.([]interface{})...)

		marker := utils.PathSearch("page_info.next_marker", listRemediationExceptionRespBody, "")
		if marker == "" {
			break
		}
		path = fmt.Sprintf("%s&marker=%s", listExceptionPath, marker)
	}
	if len(rst) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return rst, nil
}

func flattenRemediationExceptions(exceptions []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(exceptions))
	for _, exception := range exceptions {
		rst = append(rst, map[string]interface{}{
			"resource_id": utils.PathSearch("resource_id", exception, nil),
			"message":     utils.PathSearch("message", exception, nil),
			"joined_at":   utils.PathSearch("joined_at", exception, nil),
			"created_by":  utils.PathSearch("created_by", exception, nil),
		})
	}
	return rst
}

func resourceRemediationExceptionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	product := "rms"
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	policyAssignmentID := d.Get("policy_assignment_id").(string)
	exceptions := d.Get("exceptions").(*schema.Set)
	err = removeResourceFromException(client, cfg.DomainID, policyAssignmentID, exceptions.List())
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "invalid_request"),
			"error deleting RMS remediation exception",
		)
	}

	return nil
}

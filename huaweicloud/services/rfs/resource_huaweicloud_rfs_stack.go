package rfs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/rf/v1/stacks"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RFS POST /v1/{project_id}/stacks
// @API RFS GET /v1/{project_id}/stacks
// @API RFS POST /v1/{project_id}/stacks/{stack_name}/deployments
// @API RFS GET /v1/{project_id}/stacks/{stack_name}/events
// @API RFS DELETE /v1/{project_id}/stacks/{stack_name}
func ResourceStack() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStackCreate,
		ReadContext:   resourceStackRead,
		UpdateContext: resourceStackUpdate,
		DeleteContext: resourceStackDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The region where the RFS resource stack is located.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the resource stack.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The description of the resource stack.",
			},
			"agency": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Description: "schema: Required; The name of IAM agency authorized to IAC account for " +
								"resources modification.",
						},
						"provider_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "schema: Required; The name of the provider corresponding to the IAM agency.",
						},
					},
				},
				Description: "The configuration of the agencies authorized to IAC.",
			},
			"template_body": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"template_uri"},
				Description:   "The HCL/JSON template content for deployment resources.",
			},
			"vars_body": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vars_uri"},
				RequiredWith:  []string{"template_body"},
				Description:   "The variable content for deployment resources.",
			},
			"template_uri": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The OBS address where the HCL/JSON template archive (**.zip** file, which contains all " +
					"resource **.tf.json** script files to be deployed) or **.tf.json** file is located, which " +
					"describes the target status of the deployment resources.",
			},
			"vars_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"template_uri"},
				Description: "The OBS address where the variable (**.tfvars**) file corresponding to the HCL/JSON " +
					"template located, which describes the target status of the deployment resources.",
			},
			"enable_auto_rollback": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to enable automatic rollback.",
			},
			"enable_deletion_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to enable delete protection.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the resource stack.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time.`,
			},
		},
	}
}

func buildStackAgencies(agencies []interface{}) []stacks.Agency {
	if len(agencies) < 1 {
		return nil
	}

	result := make([]stacks.Agency, len(agencies))
	for i, val := range agencies {
		agency := val.(map[string]interface{})
		result[i] = stacks.Agency{
			AgencyName:   agency["name"].(string),
			ProviderName: agency["provider_name"].(string),
		}
	}

	return result
}

func buildStackCreateOpts(d *schema.ResourceData) stacks.CreateOpts {
	return stacks.CreateOpts{
		Name:                     d.Get("name").(string),
		Agencies:                 buildStackAgencies(d.Get("agency").([]interface{})),
		Description:              d.Get("description").(string),
		EnableAutoRollback:       utils.Bool(d.Get("enable_auto_rollback").(bool)),
		EnableDeletionProtection: utils.Bool(d.Get("enable_deletion_protection").(bool)),
	}
}

func resourceStackCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.AosV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOS v1 client: %s", err)
	}

	resp, err := stacks.Create(client, buildStackCreateOpts(d))
	if err != nil {
		return diag.Errorf("error creating stack: %s", err)
	}
	d.SetId(resp.StackId)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: stackStatusRefreshFunc(client, d.Id(), []string{
			string(stacks.StackStatusCreationComplete),
		}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("template_body") != "" || d.Get("template_uri") != "" {
		if err = deployStack(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceStackRead(ctx, d, meta)
}

func deployStack(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	var (
		stackId   = d.Id()
		stackName = d.Get("name").(string)

		opts = stacks.DeployOpts{
			TemplateBody: d.Get("template_body").(string),
			TemplateUri:  d.Get("template_uri").(string),
			VarsBody:     d.Get("vars_body").(string),
			VarsUri:      d.Get("vars_uri").(string),
			StackId:      stackId,
		}
	)

	deploymentId, err := stacks.Deploy(client, stackName, opts)
	if err != nil {
		return fmt.Errorf("error deploying stack resources: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: stackStatusRefreshFunc(client, d.Id(), []string{
			string(stacks.StackStatusDeploymentComplete),
		}),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 15 * time.Second,
	}
	resp, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if stack, ok := resp.(stacks.Stack); ok && stack.Status == string(stacks.StackStatusDeploymentFailed) {
			return queryAllFailedEvents(client, stackId, stackName, deploymentId)
		}
	}
	return err
}

func queryAllFailedEvents(client *golangsdk.ServiceClient, stackId, stackName,
	deploymentId string) error {
	opts := stacks.ListEventsOpts{
		StackId:      stackId,
		DeploymentId: deploymentId,
	}
	events, err := stacks.ListAllEvents(client, stackName, opts)
	if err != nil {
		return err
	}

	var mErr *multierror.Error
	for _, event := range events {
		if event.EventType == string(stacks.EventTypeError) {
			mErr = multierror.Append(mErr, fmt.Errorf(event.EventMessage))
		}
	}
	return mErr.ErrorOrNil()
}

// QueryStackById is a method to query stack details using its ID.
func QueryStackById(client *golangsdk.ServiceClient, stackId string) (*stacks.Stack, error) {
	resp, err := stacks.ListAll(client)
	if err != nil {
		return nil, err
	}

	filter := map[string]interface{}{
		"ID": stackId,
	}
	result, err := utils.FilterSliceWithField(resp, filter)
	if err != nil {
		return nil, err
	}
	if len(result) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}

	stack, ok := result[0].(stacks.Stack)
	if !ok {
		return nil, fmt.Errorf("invaid object type, want 'stack.Stack', but '%T'", result[0])
	}
	log.Printf("[DEBUG] The details of the stack (%s) is: %#v", stackId, stack)

	return &stack, nil
}

func stackStatusRefreshFunc(client *golangsdk.ServiceClient, stackId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := QueryStackById(client, stackId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return resp, "COMPLETED", nil
			}

			return nil, "", err
		}
		log.Printf("[DEBUG] The details of the resource stack (%s) is: %#v", stackId, resp)

		errorStatus := []string{
			string(stacks.StackStatusDeploymentFailed),
			string(stacks.StackStatusRollbackFailed),
			string(stacks.StackStatusDeletionFailed),
		}
		if utils.StrSliceContains(errorStatus, resp.Status) {
			return resp, "", fmt.Errorf("unexpected status '%s'", resp.Status)
		}
		if utils.StrSliceContains(targets, resp.Status) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

func resourceStackRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.AosV1Client(region)
	if err != nil {
		return diag.Errorf("error creating AOS v1 client: %s", err)
	}

	stackId := d.Id()
	resp, err := QueryStackById(client, stackId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "RFS resource stack")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("status", resp.Status),
		d.Set("created_at", resp.CreatedAt),
		d.Set("updated_at", resp.UpdatedAt),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving stack (%s) fields: %s", stackId, mErr)
	}
	return nil
}

func resourceStackUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.AosV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOS v1 client: %s", err)
	}

	if err = deployStack(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return diag.FromErr(err)
	}
	return resourceStackRead(ctx, d, meta)
}

func resourceStackDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.AosV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOS v1 client: %s", err)
	}

	var (
		stackName = d.Get("name").(string)
		stackId   = d.Id()

		opts = stacks.DeleteOpts{
			StackId: stackId,
		}
	)

	err = stacks.Delete(client, stackName, opts)
	if err != nil {
		return diag.Errorf("error deleting stack (%s): %s", stackId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      stackStatusRefreshFunc(client, stackId, nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

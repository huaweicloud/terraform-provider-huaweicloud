package apig

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/environments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/env-variables
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/env-variables/{env_variable_id}
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/env-variables
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/env-variables/{env_variable_id}
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/env-variables/{env_variable_id}
func ResourceEnvironmentVariable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentVariableCreate,
		ReadContext:   resourceEnvironmentVariableRead,
		UpdateContext: resourceEnvironmentVariableUpdate,
		DeleteContext: resourceEnvironmentVariableDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceEnvironmentVariableResourceImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Specifies the region in which to create the resource.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the ID of the dedicated instance to which the environment variable belongs.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the ID of the group to which the environment variable belongs.",
			},
			"env_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the ID of the environment to which the environment variable belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the name of the environment variable.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the value of the environment variable.",
			},
		},
	}
}

func resourceEnvironmentVariableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	opt := environments.CreateVariableOpts{
		GroupId: d.Get("group_id").(string),
		EnvId:   d.Get("env_id").(string),
		Name:    d.Get("name").(string),
		Value:   d.Get("value").(string),
	}
	resp, err := environments.CreateVariable(client, d.Get("instance_id").(string), opt).Extract()
	if err != nil {
		return diag.Errorf("error creating dedicated environment variable: %s", err)
	}
	d.SetId(resp.Id)

	return resourceEnvironmentVariableRead(ctx, d, meta)
}

func resourceEnvironmentVariableRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	variableId := d.Id()
	resp, err := environments.GetVariable(client, d.Get("instance_id").(string), variableId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "dedicated environment variable")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("group_id", resp.GroupId),
		d.Set("env_id", resp.EnvId),
		d.Set("name", resp.Name),
		d.Set("value", resp.Value),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving dedicated environment variable (%s) fields: %s", variableId, mErr)
	}
	return nil
}

func resourceEnvironmentVariableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	variableId := d.Id()
	opt := environments.UpdateVariableOpts{
		InstanceId: d.Get("instance_id").(string),
		VariableId: variableId,
		Value:      d.Get("value").(string),
	}
	_, err = environments.UpdateVariable(client, opt).Extract()
	if err != nil {
		return diag.Errorf("error updating dedicated environment variable (%s): %s", variableId, err)
	}

	return resourceEnvironmentVariableRead(ctx, d, meta)
}

func resourceEnvironmentVariableDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}
	variableId := d.Id()
	err = environments.DeleteVariable(client, d.Get("instance_id").(string), d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting dedicated environment variable(%s): %s", variableId, err)
	}

	return nil
}

func resourceEnvironmentVariableResourceImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<group_id>/<name>', but got '%s'",
			importedId)
	}

	instanceId := parts[0]
	groupId := parts[1]
	mErr := multierror.Append(
		d.Set("instance_id", instanceId),
		d.Set("group_id", groupId),
	)
	if mErr.ErrorOrNil() != nil {
		return []*schema.ResourceData{d}, mErr
	}

	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	variables, err := queryEnvironmentVariables(client, instanceId, groupId)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error getting environment variables: %s", err)
	}

	// The "name" paraqmeter support fuzzy search.
	variableName := parts[2]
	for _, variable := range variables {
		if variable.Name == variableName {
			d.SetId(variable.Id)
			return []*schema.ResourceData{d}, nil
		}
	}

	return []*schema.ResourceData{d}, fmt.Errorf("environment variable (%s) not found: %s", variableName, err)
}

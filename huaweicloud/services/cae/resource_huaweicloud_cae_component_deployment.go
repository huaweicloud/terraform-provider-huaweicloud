package cae

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CAE POST /v1/{project_id}/cae/applications/{application_id}/components/{component_id}/action
// @API CAE GET /v1/{project_id}/cae/jobs/{job_id}
// ResourceComponentDeployment is a definition of the one-time action resource that used to manage component deployment.
func ResourceComponentDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComponentDeploymentCreate,
		ReadContext:   resourceComponentDeploymentRead,
		UpdateContext: resourceComponentDeploymentUpdate,
		DeleteContext: resourceComponentDeploymentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region in which to create the resource.`,
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the environment where the application is located.`,
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the application where the component is located.`,
			},
			"component_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the component to which the configurations belong.`,
			},
			"metadata": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"annotations": {
							Type:        schema.TypeMap,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The resource configurations.`,
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The action name.`,
						},
					},
				},
				Description: `The metadata of this action request.`,
			},
			"spec": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  `The specification detail of the action.`,
			},
		},
	}
}

func buildCreateComponentDeploymentBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"api_version": "v1",
		"kind":        "Action",
		"metadata": map[string]interface{}{
			"annotations": d.Get("metadata.0.annotations"),
			"name":        d.Get("metadata.0.name"),
		},
		"spec": utils.ValueIgnoreEmpty(unmarshalJsonFormatParamster("Specification detail", d.Get("spec").(string))),
	}
}

func deployComponent(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout time.Duration) error {
	var (
		httpUrl       = "v1/{project_id}/cae/applications/{application_id}/components/{component_id}/action"
		applicationId = d.Get("application_id").(string)
		componentId   = d.Get("component_id").(string)
		environmentId = d.Get("environment_id").(string)
	)

	modifyPath := client.Endpoint + httpUrl
	modifyPath = strings.ReplaceAll(modifyPath, "{project_id}", client.ProjectID)
	modifyPath = strings.ReplaceAll(modifyPath, "{application_id}", applicationId)
	modifyPath = strings.ReplaceAll(modifyPath, "{component_id}", componentId)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Environment-ID": environmentId,
		},
		JSONBody: utils.RemoveNil(buildCreateComponentDeploymentBodyParams(d)),
	}
	requestResp, err := client.Request("POST", modifyPath, &opts)
	if err != nil {
		return fmt.Errorf("error operating the component (%s): %s", componentId, err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return fmt.Errorf("error retrieving API response of the deployment for the component configuration (%s): %s", componentId, err)
	}
	jobId := utils.PathSearch("job_id", respBody, "null").(string)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      deployJobRefreshFunc(client, environmentId, jobId, []string{"success"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        20 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the deploy job (%s) success: %s", jobId, err)
	}
	return nil
}

func getDeployJobDetail(client *golangsdk.ServiceClient, environmentId, jobId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/cae/jobs/{job_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_id}", jobId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Environment-ID": environmentId,
		},
	}
	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying deploy job detail by its ID (%s): %s", jobId, err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving deploy job (%s) detail: %s", jobId, err)
	}
	return respBody, nil
}

func deployJobRefreshFunc(client *golangsdk.ServiceClient, environmentId, jobId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := getDeployJobDetail(client, environmentId, jobId)
		if err != nil {
			return resp, "ERROR", err
		}

		status := utils.PathSearch("spec.status", resp, "null").(string)

		if utils.StrSliceContains(targets, status) {
			return resp, "COMPLETED", nil
		}
		return resp, "PENDING", nil
	}
}

func resourceComponentDeploymentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		componentId = d.Get("component_id").(string)
	)
	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	err = deployComponent(ctx, client, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(componentId)

	return resourceComponentDeploymentRead(ctx, d, meta)
}

func resourceComponentDeploymentRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceComponentDeploymentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	err = deployComponent(ctx, client, d, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceComponentDeploymentRead(ctx, d, meta)
}

func resourceComponentDeploymentDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a one-time action resource.
	return nil
}

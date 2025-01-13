package cae

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

var componentResourceNotFoundCodes = []string{
	"CAE.01500208", // Application or component does not found, and status code is 400.
	"CAE.01500404", // Environment does not found, and status code is 400.
	"CAE.01500000", // Application or component does not found, and status code is 500.
}

// @API CAE POST /v1/{project_id}/cae/applications/{application_id}/components
// @API CAE POST /v1/{project_id}/cae/applications/{application_id}/component-with-configurations
// @API CAE GET /v1/{project_id}/cae/jobs/{job_id}
// @API CAE GET /v1/{project_id}/cae/applications/{application_id}/components/{component_id}
// @API CAE PUT /v1/{project_id}/cae/applications/{application_id}/components/{component_id}
// @API CAE DELETE /v1/{project_id}/cae/applications/{application_id}/components/{component_id}
func ResourceComponent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComponentCreate,
		ReadContext:   resourceComponentRead,
		UpdateContext: resourceComponentUpdate,
		DeleteContext: resourceComponentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceComponentImportState,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metadata": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"annotations": {
							Type:     schema.TypeMap,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"spec": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"replica": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"runtime": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem:     sourceSchema(),
						},
						"resource_limit": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu": {
										Type:     schema.TypeString,
										Required: true,
									},
									"memory": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"build": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     buildSchema(),
						},
					},
				},
			},
			"deploy_after_create": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Whether to deploy the component after creating the resource.`,
			},
			"configurations": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: `The type of the component configuration.`,
						},
						"data": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  `The component configuration detail, in JSON format.`,
						},
					},
				},
				Description: `The list of configurations of the component.`,
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

func buildSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"archive": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"artifact_namespace": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"parameters": {
				Type:     schema.TypeMap,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func sourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sub_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"code": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"branch": {
							Type:     schema.TypeString,
							Required: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func buildMetadata(metadata map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"name":        metadata["name"],
		"annotations": metadata["annotations"],
	}
}

func buildCode(code []interface{}) map[string]interface{} {
	if len(code) == 0 {
		return nil
	}
	raw := code[0].(map[string]interface{})
	return map[string]interface{}{
		"auth_name": raw["auth_name"],
		"branch":    raw["branch"],
		"namespace": raw["namespace"],
	}
}

func buildSource(source []interface{}) map[string]interface{} {
	raw := source[0].(map[string]interface{})
	return map[string]interface{}{
		"type":     raw["type"],
		"url":      raw["url"],
		"sub_type": utils.ValueIgnoreEmpty(raw["sub_type"]),
		"code":     buildCode(raw["code"].([]interface{})),
	}
}

func buildArchiveInfo(build []interface{}) map[string]interface{} {
	if len(build) == 0 {
		return nil
	}

	raw := build[0].(map[string]interface{})
	archive := raw["archive"].([]interface{})[0].(map[string]interface{})
	return map[string]interface{}{
		"archive": map[string]interface{}{
			"artifact_namespace": archive["artifact_namespace"],
		},
		"parameters": raw["parameters"],
	}
}

func buildResourceLimit(resourceLimit []interface{}) map[string]interface{} {
	raw := resourceLimit[0].(map[string]interface{})
	return map[string]interface{}{
		"cpu_limit":    raw["cpu"],
		"memory_limit": raw["memory"],
	}
}

func buildSpec(specification map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"replica":        specification["replica"],
		"runtime":        specification["runtime"],
		"source":         buildSource(specification["source"].([]interface{})),
		"resource_limit": buildResourceLimit(specification["resource_limit"].([]interface{})),
		"build":          buildArchiveInfo(specification["build"].([]interface{})),
	}
}

func buildCreateOrUpdateComponentBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"api_version": "v1",
		"kind":        "Component",
		"metadata":    buildMetadata(d.Get("metadata.0").(map[string]interface{})),
		"spec":        buildSpec(d.Get("spec.0").(map[string]interface{})),
	}
}

func resourceComponentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg               = meta.(*config.Config)
		environmentId     = d.Get("environment_id").(string)
		applicationId     = d.Get("application_id").(string)
		isDeployComponent = d.Get("deploy_after_create").(bool)
		err               error
	)
	client, err := cfg.NewServiceClient("cae", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	if isDeployComponent {
		err = createAndDeployComponent(ctx, client, d, environmentId, applicationId)
	} else {
		err = createComponent(client, d, environmentId, applicationId)
	}

	if err != nil {
		return diag.FromErr(err)
	}
	return resourceComponentRead(ctx, d, meta)
}

func createComponent(client *golangsdk.ServiceClient, d *schema.ResourceData, environmentId, applicationId string) error {
	httpUrl := "v1/{project_id}/cae/applications/{application_id}/components"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{application_id}", applicationId)

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Environment-Id": environmentId,
		},
		JSONBody: utils.RemoveNil(buildCreateOrUpdateComponentBodyParams(d)),
	}
	createResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return fmt.Errorf("error creating CAE component under specified application (%s): %s", applicationId, err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return err
	}

	componentId := utils.PathSearch("metadata.id", createRespBody, "").(string)
	if componentId == "" {
		return fmt.Errorf("unable to find the CAE component ID from the API response")
	}

	d.SetId(componentId)
	return nil
}

func createAndDeployComponent(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, environmentId,
	applicationId string) error {
	httpUrl := "v1/{project_id}/cae/applications/{application_id}/component-with-configurations"
	createDeployPath := client.Endpoint + httpUrl
	createDeployPath = strings.ReplaceAll(createDeployPath, "{project_id}", client.ProjectID)
	createDeployPath = strings.ReplaceAll(createDeployPath, "{application_id}", applicationId)

	params := buildCreateOrUpdateComponentBodyParams(d)
	params["configurations"] = buildConfigurationItemsBodyParams(d.Get("configurations").(*schema.Set))
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Environment-Id": environmentId,
		},
		JSONBody: utils.RemoveNil(params),
	}
	resp, err := client.Request("POST", createDeployPath, &opt)
	if err != nil {
		return fmt.Errorf("error creating and deploying CAE component under specified application (%s): %s", applicationId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	componentId := utils.PathSearch("metadata.id", respBody, "").(string)
	if componentId == "" {
		return fmt.Errorf("unable to find the CAE component ID from the API response")
	}

	d.SetId(componentId)

	jobId := utils.PathSearch("status.job_id", respBody, "").(string)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      deployJobRefreshFunc(client, environmentId, jobId, []string{"success"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        20 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the deploy component job (%s) to complete: %s", jobId, err)
	}
	return nil
}

// GetComponentById is a method to query component details from a specified application ID using given parameters.
func GetComponentById(client *golangsdk.ServiceClient, environmentId, applicationId, componentId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/cae/applications/{application_id}/components/{component_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{application_id}", applicationId)
	getPath = strings.ReplaceAll(getPath, "{component_id}", componentId)
	getComponentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Environment-Id": environmentId,
		},
	}
	resp, err := client.Request("GET", getPath, &getComponentOpt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(
			common.ConvertExpected500ErrInto404Err(err, "error_code", componentResourceNotFoundCodes...),
			"error_code",
			componentResourceNotFoundCodes...)
	}

	return utils.FlattenResponse(resp)
}

func resourceComponentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	componentId := d.Id()

	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	componentRespBody, err := GetComponentById(client, d.Get("environment_id").(string), d.Get("application_id").(string), componentId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving CAE component (%s)", componentId))
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("environment_id", utils.PathSearch("spec.env_id", componentRespBody, nil)),
		d.Set("metadata", flattenMetadata(d.Get("metadata.0.annotations"), utils.PathSearch("metadata", componentRespBody, nil))),
		d.Set("spec", flattenSpec(d.Get("spec.0.build.0.parameters"), utils.PathSearch("spec", componentRespBody, nil))),
		d.Set("created_at", utils.PathSearch("metadata.created_at", componentRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("metadata.updated_at", componentRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMetadata(annotations interface{}, respMetadata interface{}) []map[string]interface{} {
	if respMetadata == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"name":        utils.PathSearch("name", respMetadata, nil),
			"annotations": annotations,
		},
	}
}

func flattenSpec(parameters interface{}, spec interface{}) []map[string]interface{} {
	if spec == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"replica":        utils.PathSearch("replica", spec, nil),
			"runtime":        utils.PathSearch("runtime", spec, nil),
			"source":         flattenSource(utils.PathSearch("source", spec, nil)),
			"resource_limit": flattenResourceLimit(utils.PathSearch("resource_limit", spec, nil)),
			"build":          flattenBuild(parameters, utils.PathSearch("build", spec, nil)),
		},
	}
}

func flattenResourceLimit(resourceLimit interface{}) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"cpu":    utils.PathSearch("cpu_limit", resourceLimit, nil),
			"memory": utils.PathSearch("memory_limit", resourceLimit, nil),
		},
	}
}

func flattenSource(source interface{}) []map[string]interface{} {
	if source == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":     utils.PathSearch("type", source, nil),
			"url":      utils.PathSearch("url", source, nil),
			"sub_type": utils.PathSearch("sub_type", source, nil),
			"code":     flattenCode(utils.PathSearch("code", source, nil)),
		},
	}
}

func flattenCode(code interface{}) []map[string]interface{} {
	if (code) == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"auth_name": utils.PathSearch("auth_name", code, nil),
			"branch":    utils.PathSearch("branch", code, nil),
			"namespace": utils.PathSearch("namespace", code, nil),
		},
	}
}

func flattenBuild(parameters interface{}, build interface{}) []map[string]interface{} {
	archive := utils.PathSearch("archive", build, nil)
	if archive == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"archive": []map[string]interface{}{
				{
					"artifact_namespace": utils.PathSearch("artifact_namespace", archive, nil),
				},
			},
			"parameters": parameters,
		},
	}
}

func resourceComponentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/cae/applications/{application_id}/components/{component_id}"
		product = "cae"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{application_id}", d.Get("application_id").(string))
	componentId := d.Id()
	updatePath = strings.ReplaceAll(updatePath, "{component_id}", componentId)

	updateComponentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Environment-Id": d.Get("environment_id").(string),
		},
		JSONBody: utils.RemoveNil(buildCreateOrUpdateComponentBodyParams(d)),
	}
	_, err = client.Request("PUT", updatePath, &updateComponentOpt)
	if err != nil {
		return diag.Errorf("error updating CAE component (%s): %s", componentId, err)
	}

	return resourceComponentRead(ctx, d, meta)
}

func resourceComponentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v1/{project_id}/cae/applications/{application_id}/components/{component_id}"
		product       = "cae"
		environmentId = d.Get("environment_id").(string)
		applicationId = d.Get("application_id").(string)
		componentId   = d.Id()
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{application_id}", applicationId)
	deletePath = strings.ReplaceAll(deletePath, "{component_id}", componentId)

	deleteComponentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Environment-Id": environmentId,
		},
	}
	resp, err := client.Request("DELETE", deletePath, &deleteComponentOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", componentResourceNotFoundCodes...),
			fmt.Sprintf("error deleting CAE component (%s)", componentId))
	}

	_, err = utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      refreshDeleteComponentFunc(client, environmentId, applicationId, componentId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for deleting component to complete: %s", err)
	}

	return nil
}

func refreshDeleteComponentFunc(client *golangsdk.ServiceClient, environmentId, applicationId, componentId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetComponentById(client, environmentId, applicationId, componentId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "deleted", "DELETED", nil
			}
			return nil, "ERROR", err
		}
		return respBody, "PENDING", nil
	}
}

func resourceComponentImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<environment_id>/<application_id>/<id>', but got '%s'",
			importedId)
	}
	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("environment_id", parts[0]),
		d.Set("application_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

package servicestage

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var environmentAssociateObjSliceParamKeys = []string{
	"resources",
}

// @API ServiceStage PUT /v3/{project_id}/cas/environments/{environment_id}/resources
// @API ServiceStage GET /v3/{project_id}/cas/environments/{environment_id}/resources
func ResourceV3EnvironmentAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3EnvironmentAssociateCreate,
		ReadContext:   resourceV3EnvironmentAssociateRead,
		UpdateContext: resourceV3EnvironmentAssociateUpdate,
		DeleteContext: resourceV3EnvironmentAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV3EnvironmentAssociateImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the environment and resources are located.`,
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The environment ID associated with the resources.`,
			},
			"resources": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The ID of the resource to be associated.`,
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the resource to be associated.`,
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The name of the resource to be associated.`,
						},
					},
				},
				Description: "The information about the associated resources.",
			},

			// Internal attributes.
			"resources_origin": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the resource to be associated.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the resource to be associated.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the resource to be associated.`,
						},
					},
				},
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
 the new value next time the change is made. The corresponding parameter name is 'resources'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildV3EnvironmentAssociatedResources(resources []interface{}) []interface{} {
	result := make([]interface{}, 0, len(resources))
	for _, v := range resources {
		result = append(result, map[string]interface{}{
			"id":   utils.PathSearch("id", v, nil),
			"type": utils.PathSearch("type", v, nil),
			"name": utils.PathSearch("name", v, nil),
		})
	}
	return result
}

func modifyAssociatedResourcesUnderEnvironment(client *golangsdk.ServiceClient, envId string, resources []interface{}) error {
	httpUrl := "v3/{project_id}/cas/environments/{environment_id}/resources"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{environment_id}", envId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"resources": buildV3EnvironmentAssociatedResources(resources),
		},
	}

	_, err := client.Request("PUT", createPath, &opt)
	return err
}

func resourceV3EnvironmentAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		envId     = d.Get("environment_id").(string)
		resources = d.Get("resources").([]interface{})
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	err = modifyAssociatedResourcesUnderEnvironment(client, envId, resources)
	if err != nil {
		return diag.Errorf("error associating resources to the environment (%s): %s", envId, err)
	}
	d.SetId(envId)

	// If the request is successful, obtain the values ​​of all JSON parameters first and save them to the
	// corresponding '_origin' attributes for subsequent determination and construction of the request body during
	// next updates.
	err = utils.RefreshObjectParamOriginValues(d, environmentAssociateObjSliceParamKeys)
	if err != nil {
		return diag.Errorf("unable to refresh the origin values: %s", err)
	}

	return resourceV3EnvironmentAssociateRead(ctx, d, meta)
}

func ListV3EnvironmentAssociatedResources(client *golangsdk.ServiceClient, envId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/cas/environments/{environment_id}/resources"

	queryPath := client.Endpoint + httpUrl
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath = strings.ReplaceAll(queryPath, "{environment_id}", envId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", queryPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func orderV3EnvironmentAssociatedResourcesByResourcesOrigin(resources []map[string]interface{},
	resourcesOrigin []interface{}) []map[string]interface{} {
	if len(resourcesOrigin) < 1 {
		return resources
	}

	sortedResources := make([]map[string]interface{}, 0, len(resources))
	resourcesCopy := resources
	for _, resourceOrigin := range resourcesOrigin {
		idOrigin := utils.PathSearch("id", resourceOrigin, "").(string)
		typeOrigin := utils.PathSearch("type", resourceOrigin, "").(string)
		for index, resource := range resourcesCopy {
			if utils.PathSearch("id", resource, "").(string) != idOrigin || utils.PathSearch("type", resource, "").(string) != typeOrigin {
				continue
			}

			sortedResources = append(sortedResources, resourcesCopy[index])
			resourcesCopy = append(resourcesCopy[:index], resourcesCopy[index+1:]...)
			break
		}
	}

	// Add any remaining unsorted resources to the end of the sorted list.
	sortedResources = append(sortedResources, resourcesCopy...)
	return sortedResources
}

func flattenV3EnvironmentAssociatedResources(resources []interface{}, resourcesOrigin []interface{}) []map[string]interface{} {
	if len(resources) < 1 {
		return nil
	}

	parsedResources := make([]map[string]interface{}, 0, len(resources))
	for _, resource := range resources {
		parsedResources = append(parsedResources, map[string]interface{}{
			"id":   utils.PathSearch("id", resource, nil),
			"type": utils.PathSearch("type", resource, nil),
			"name": utils.PathSearch("name", resource, nil),
		})
	}

	return orderV3EnvironmentAssociatedResourcesByResourcesOrigin(parsedResources, resourcesOrigin)
}

func resourceV3EnvironmentAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		envId  = d.Id()
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	respBody, err := ListV3EnvironmentAssociatedResources(client, envId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected401ErrInto404Err(err, "error_code", v3EnvNotFoundCodes...),
			fmt.Sprintf("error getting environment (%s)", envId))
	}
	associatedResources := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
	if len(associatedResources) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v3/{project_id}/cas/environments/{environment_id}/resources",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("All assiciated resources have been dissociated from the environment (%s)", envId)),
			},
		}, "error retrieving associated resources")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("resources", flattenV3EnvironmentAssociatedResources(associatedResources, d.Get("resources_origin").([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceV3EnvironmentAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		envId     = d.Id()
		resources = d.Get("resources").([]interface{})
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	err = modifyAssociatedResourcesUnderEnvironment(client, envId, resources)
	if err != nil {
		return diag.Errorf("error modifying associated resources under the environment (%s): %s", envId, err)
	}

	// If the request is successful, obtain the values ​​of all JSON parameters first and save them to the
	// corresponding '_origin' attributes for subsequent determination and construction of the request body during
	// next updates.
	err = utils.RefreshObjectParamOriginValues(d, environmentAssociateObjSliceParamKeys)
	if err != nil {
		return diag.Errorf("unable to refresh the origin values: %s", err)
	}

	return resourceV3EnvironmentAssociateRead(ctx, d, meta)
}

func resourceV3EnvironmentAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		envId  = d.Id()
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	err = modifyAssociatedResourcesUnderEnvironment(client, envId, make([]interface{}, 0))
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected401ErrInto404Err(err, "error_code", v3EnvNotFoundCodes...),
			fmt.Sprintf("error dissociating the resources from the environment (%s)", envId))
	}
	return nil
}

func resourceV3EnvironmentAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, d.Set("environment_id", d.Id())
}

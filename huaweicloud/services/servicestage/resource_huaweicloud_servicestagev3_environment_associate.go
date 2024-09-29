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
				Type:     schema.TypeSet,
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
					},
				},
				Description: "The information about the associated resources.",
			},
		},
	}
}

func buildV3EnvironmentAssociatedResources(resources *schema.Set) []interface{} {
	result := make([]interface{}, 0, resources.Len())
	for _, v := range resources.List() {
		result = append(result, map[string]interface{}{
			"id":   utils.PathSearch("id", v, nil),
			"type": utils.PathSearch("type", v, nil),
		})
	}
	return result
}

func modifyAssociatedResourcesUnderEnvironment(client *golangsdk.ServiceClient, envId string, resources *schema.Set) error {
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
		resources = d.Get("resources").(*schema.Set)
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

	return resourceV3EnvironmentAssociateRead(ctx, d, meta)
}

func QueryV3EnvironmentAssociatedResources(client *golangsdk.ServiceClient, envId string) (interface{}, error) {
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

func flattenV3EnvironmentAssociatedResources(resources []interface{}) []map[string]interface{} {
	if len(resources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(resources))
	for _, resource := range resources {
		result = append(result, map[string]interface{}{
			"id":   utils.PathSearch("id", resource, nil),
			"type": utils.PathSearch("type", resource, nil),
		})
	}
	return result
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

	respBody, err := QueryV3EnvironmentAssociatedResources(client, envId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected401ErrInto404Err(err, "error_code", v3EnvNotFoundCodes...),
			fmt.Sprintf("error getting environment (%s)", envId))
	}
	associatedResources := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
	if len(associatedResources) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "associated resources not found")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("resources", flattenV3EnvironmentAssociatedResources(associatedResources)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceV3EnvironmentAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		envId     = d.Id()
		resources = d.Get("resources").(*schema.Set)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	err = modifyAssociatedResourcesUnderEnvironment(client, envId, resources)
	if err != nil {
		return diag.Errorf("error modifying associated resources under the environment (%s): %s", envId, err)
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

	err = modifyAssociatedResourcesUnderEnvironment(client, envId, schema.NewSet(schema.HashString, nil))
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected401ErrInto404Err(err, "error_code", v3EnvNotFoundCodes...),
			fmt.Sprintf("error dissociating the resources from the environment (%s)", envId))
	}
	return nil
}

func resourceV3EnvironmentAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, d.Set("environment_id", d.Id())
}

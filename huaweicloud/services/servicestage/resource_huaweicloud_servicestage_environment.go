package servicestage

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/servicestage/v2/environments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var environmentObjSliceParamKeys = []string{
	"basic_resources",
	"optional_resources",
}

// @API ServiceStage POST /v2/{project_id}/cas/environments
// @API ServiceStage GET /v2/{project_id}/cas/environments/{environment_id}
// @API ServiceStage PUT /v2/{project_id}/cas/environments/{environment_id}
// @API ServiceStage PATCH /v2/{project_id}/cas/environments/{environment_id}/resources
// @API ServiceStage DELETE /v2/{project_id}/cas/environments/{environment_id}
func ResourceEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deploy_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"basic_resources": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				DiffSuppressFunc: utils.SuppressObjectSliceDiffs(),
			},
			"optional_resources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				DiffSuppressFunc: utils.SuppressObjectSliceDiffs(),
			},

			// Internal attributes.
			"basic_resources_origin": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
 the new value next time the change is made. The corresponding parameter name is 'basic_resources'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"optional_resources_origin": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
 the new value next time the change is made. The corresponding parameter name is 'optional_resources'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildResourcesList(resources []interface{}) []environments.Resource {
	if len(resources) < 1 {
		return nil
	}

	result := make([]environments.Resource, 0, len(resources))
	for _, v := range resources {
		res := v.(map[string]interface{})
		result = append(result, environments.Resource{
			Type: res["type"].(string),
			ID:   res["id"].(string),
			Name: utils.PathSearch("name", res, "").(string),
		})
	}

	return result
}

func buildEnvironmentCreateOpts(d *schema.ResourceData, cfg *config.Config) environments.CreateOpts {
	desc := d.Get("description").(string)
	return environments.CreateOpts{
		Name:                d.Get("name").(string),
		Description:         &desc,
		VpcId:               d.Get("vpc_id").(string),
		DeployMode:          d.Get("deploy_mode").(string),
		BaseResources:       buildResourcesList(d.Get("basic_resources").([]interface{})),
		OptionalResources:   buildResourcesList(d.Get("optional_resources").([]interface{})),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}
}

func orderEnvironmentResourcesByOrigin(resources []map[string]interface{},
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

	sortedResources = append(sortedResources, resourcesCopy...)
	return sortedResources
}

func flattenEnvironmentResources(resources []environments.Resource, resourcesOrigin []interface{}) []map[string]interface{} {
	if len(resources) < 1 {
		return nil
	}

	parsedResources := make([]map[string]interface{}, 0, len(resources))
	for _, v := range resources {
		parsedResources = append(parsedResources, map[string]interface{}{
			"type": v.Type,
			"id":   v.ID,
			"name": v.Name,
		})
	}

	return orderEnvironmentResourcesByOrigin(parsedResources, resourcesOrigin)
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.ServiceStageV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ServiceStage v2 client: %s", err)
	}

	opt := buildEnvironmentCreateOpts(d, config)
	log.Printf("[DEBUG] The createOpt of ServiceStage environment is: %v", opt)
	resp, err := environments.Create(client, opt)
	if err != nil {
		return diag.Errorf("error creating ServiceStage environment: %s", err)
	}

	d.SetId(resp.ID)

	// If the request is successful, obtain the values ​​of all JSON parameters first and save them to the
	// corresponding '_origin' attributes for subsequent determination and construction of the request body during
	// next updates.
	err = utils.RefreshObjectParamOriginValues(d, environmentObjSliceParamKeys)
	if err != nil {
		return diag.Errorf("unable to refresh the origin values: %s", err)
	}

	return resourceEnvironmentRead(ctx, d, meta)
}

func resourceEnvironmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ServiceStageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage v2 client: %s", err)
	}

	resp, err := environments.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ServiceStage environment")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("vpc_id", resp.VpcId),
		d.Set("deploy_mode", resp.DeployMode),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("basic_resources", flattenEnvironmentResources(resp.BaseResources, d.Get("basic_resources_origin").([]interface{}))),
		d.Set("optional_resources", flattenEnvironmentResources(resp.OptionalResources, d.Get("optional_resources_origin").([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ServiceStageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage v2 client: %s", err)
	}

	if d.HasChanges("name", "description") {
		desc := d.Get("description").(string)
		updateOpt := environments.UpdateOpts{
			Name:        d.Get("name").(string),
			Description: &desc,
		}
		_, err = environments.Update(client, d.Id(), updateOpt)
		if err != nil {
			return diag.Errorf("error updating ServiceStage environment (%s): %s", d.Id(), err)
		}
	}

	if d.HasChanges("basic_resources", "optional_resources") {
		oldList, newList := d.GetChange("basic_resources")
		baseRes := buildResourcesList(newList.([]interface{}))
		rmRes := buildResourcesList(oldList.([]interface{}))

		oldList, newList = d.GetChange("optional_resources")
		optRes := buildResourcesList(newList.([]interface{}))
		rmRes = append(rmRes, buildResourcesList(oldList.([]interface{}))...)
		updateOpt := environments.ResourceOpts{
			AddBaseResources:     baseRes,
			AddOptionalResources: optRes,
			RemoveResources:      rmRes,
		}
		_, err := environments.UpdateResources(client, d.Id(), updateOpt)
		if err != nil {
			return diag.Errorf("error updating ServiceStage environment (%s): %s", d.Id(), err)
		}

		// If the request is successful, obtain the values ​​of all JSON parameters first and save them to the
		// corresponding '_origin' attributes for subsequent determination and construction of the request body during
		// next updates.
		err = utils.RefreshObjectParamOriginValues(d, environmentObjSliceParamKeys)
		if err != nil {
			return diag.Errorf("unable to refresh the origin values: %s", err)
		}
	}

	return resourceEnvironmentRead(ctx, d, meta)
}

func resourceEnvironmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ServiceStageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage v2 client: %s", err)
	}

	err = environments.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting ServiceStage environment (%s): %s", d.Id(), err)
	}
	return nil
}

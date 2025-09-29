package servicestage

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/servicestage/v2/environments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

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
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"cce", "cci", "ecs", "as",
							}, false),
						},
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"optional_resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"elb", "eip", "rds", "dcs", "cse",
							}, false),
						},
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func buildResourcesList(resources *schema.Set) []environments.Resource {
	if resources.Len() < 1 {
		return nil
	}

	result := make([]environments.Resource, resources.Len())
	for i, v := range resources.List() {
		res := v.(map[string]interface{})
		result[i] = environments.Resource{
			Type: res["type"].(string),
			ID:   res["id"].(string),
		}
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
		BaseResources:       buildResourcesList(d.Get("basic_resources").(*schema.Set)),
		OptionalResources:   buildResourcesList(d.Get("optional_resources").(*schema.Set)),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}
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

	return resourceEnvironmentRead(ctx, d, meta)
}

func flattenEnvironmentResources(resources []environments.Resource) []map[string]interface{} {
	if len(resources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(resources))
	for i, v := range resources {
		result[i] = map[string]interface{}{
			"type": v.Type,
			"id":   v.ID,
		}
	}

	return result
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
		d.Set("basic_resources", flattenEnvironmentResources(resp.BaseResources)),
		d.Set("optional_resources", flattenEnvironmentResources(resp.OptionalResources)),
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
		oldSet, newSet := d.GetChange("basic_resources")
		baseRes := buildResourcesList(newSet.(*schema.Set))
		rmRes := buildResourcesList(oldSet.(*schema.Set))

		oldSet, newSet = d.GetChange("optional_resources")
		optRes := buildResourcesList(newSet.(*schema.Set))
		rmRes = append(rmRes, buildResourcesList(oldSet.(*schema.Set))...)
		updateOpt := environments.ResourceOpts{
			AddBaseResources:     baseRes,
			AddOptionalResources: optRes,
			RemoveResources:      rmRes,
		}
		_, err := environments.UpdateResources(client, d.Id(), updateOpt)
		if err != nil {
			return diag.Errorf("error updating ServiceStage environment (%s): %s", d.Id(), err)
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

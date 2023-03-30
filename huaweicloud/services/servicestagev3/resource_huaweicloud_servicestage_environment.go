package servicestagev3

import (
	"context"
	"log"
	"regexp"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/servicestage/v3/environments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

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
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z]([\w-]*[A-Za-z0-9])?$`),
						"The name must start with a letter and end with a letter or digit, and can only contain "+
							"letters, digits, underscores (_) and hyphens (-)."),
					validation.StringLenBetween(2, 64),
				),
			},
			"deploy_mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"virtualmachine",
					"container",
					"mixed",
				}, false),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"vm_cluster_size": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"50",
					"500",
				}, false),
			},
			"alias": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"VPC", "EIP", "ELB", "CCE", "CCI", "ECS", "AS", "CSE", "DCS", "RDS", "PVC",
							}, true),
						},
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func buildResourcesList(resources *schema.Set) []environments.Resource {
	if resources.Len() < 1 {
		return []environments.Resource{}
	}

	result := make([]environments.Resource, resources.Len())
	for i, v := range resources.List() {
		res := v.(map[string]interface{})
		result[i] = environments.Resource{
			Type: res["type"].(string),
			ID:   res["id"].(string),
			Name: res["name"].(string),
		}
	}

	return result
}

func buildEnvironmentCreateOpts(d *schema.ResourceData, conf *config.Config) (environments.CreateOpts,
	diag.Diagnostics) {
	vmClusterSize, diagErr := buildEnvVmClusterSize(d)
	if diagErr != nil {
		return environments.CreateOpts{}, diagErr
	}

	return environments.CreateOpts{
		Name:                d.Get("name").(string),
		DeployMode:          d.Get("deploy_mode").(string),
		Description:         d.Get("description").(string),
		EnterpriseProjectId: common.GetEnterpriseProjectID(d, conf),
		VpcId:               d.Get("vpc_id").(string),
		Labels:              buildEnvLabels(d.Get("labels").([]interface{})),
		VmClusterSize:       vmClusterSize,
	}, nil
}

func buildEnvResourceUpdateOpts(d *schema.ResourceData) environments.ResourceOpts {
	return environments.ResourceOpts{
		Resources: buildResourcesList(d.Get("resources").(*schema.Set)),
	}
}

func buildEnvVmClusterSize(d *schema.ResourceData) (int, diag.Diagnostics) {
	vmClusterSize := 50
	schemaVmSize := d.Get("vm_cluster_size").(string)
	if schemaVmSize != "" {
		size, err := strconv.Atoi(schemaVmSize)
		if err != nil {
			return vmClusterSize, diag.Errorf("Creating environment failed because 'vm_cluster_size' is illegal")
		}
		vmClusterSize = size
	}
	return vmClusterSize, nil
}

func buildEnvLabels(labels []interface{}) []environments.Label {
	if len(labels) < 1 {
		return []environments.Label{}
	}

	result := make([]environments.Label, len(labels))
	for i, label := range labels {
		l := label.(map[string]interface{})
		result[i] = environments.Label{
			Key:   l["key"].(string),
			Value: l["value"].(string),
		}
	}

	return result
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hwConfig := meta.(*config.Config)
	client, err := hwConfig.ServiceStageV3Client(hwConfig.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ServiceStage v3 client: %s", err)
	}

	opt, diagErr := buildEnvironmentCreateOpts(d, hwConfig)
	if diagErr != nil {
		return diagErr
	}
	log.Printf("[DEBUG] The createOpt of ServiceStage environment is: %v", opt)
	resp, err := environments.Create(client, opt)
	if err != nil {
		return diag.Errorf("error creating ServiceStage environment: %s", err)
	}

	d.SetId(resp.ID)

	resourceOpt := buildEnvResourceUpdateOpts(d)
	if len(resourceOpt.Resources) > 0 {
		_, err = environments.UpdateResources(client, resp.ID, resourceOpt)
		if err != nil {
			return diag.Errorf("error creating ServiceStage environment: %s", err)
		}
	}

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
			"name": v.Name,
		}
	}

	return result
}

func resourceEnvironmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hwConfig := meta.(*config.Config)
	region := hwConfig.GetRegion(d)
	client, err := hwConfig.ServiceStageV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage v2 client: %s", err)
	}

	resp, err := environments.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ServiceStage environment")
	}

	resourceResp, err := environments.ListResources(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ServiceStage environment resource")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("vpc_id", resp.VpcId),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("resources", flattenEnvironmentResources(resourceResp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hwConfig := meta.(*config.Config)
	region := hwConfig.GetRegion(d)
	client, err := hwConfig.ServiceStageV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage v3 client: %s", err)
	}

	if d.HasChanges("name", "description", "labels") {
		updateOpt := environments.UpdateOpts{
			Name:        d.Get("name").(string),
			Labels:      buildEnvLabels(d.Get("labels").([]interface{})),
			Description: d.Get("description").(string),
		}
		_, err = environments.Update(client, d.Id(), updateOpt)
		if err != nil {
			return diag.Errorf("error updating ServiceStage environment (%s): %s", d.Id(), err)
		}
	}

	if d.HasChanges("resources") {
		updateOpt := environments.ResourceOpts{
			Resources: buildResourcesList(d.Get("resources").(*schema.Set)),
		}
		_, err := environments.UpdateResources(client, d.Id(), updateOpt)
		if err != nil {
			return diag.Errorf("error updating ServiceStage environment (%s): %s", d.Id(), err)
		}
	}

	return resourceEnvironmentRead(ctx, d, meta)
}

func resourceEnvironmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hwConfig := meta.(*config.Config)
	region := hwConfig.GetRegion(d)
	client, err := hwConfig.ServiceStageV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage v3 client: %s", err)
	}

	err = environments.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting ServiceStage environment (%s): %s", d.Id(), err)
	}
	return nil
}

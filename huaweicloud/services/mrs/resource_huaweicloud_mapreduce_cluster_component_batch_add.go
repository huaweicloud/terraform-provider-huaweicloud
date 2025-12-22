package mrs

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var clusterComponentBatchAddNonUpdateParams = []string{
	"cluster_id",
	"components_install_mode",
	"components_install_mode.*.component",
	"components_install_mode.*.node_groups",
	"components_install_mode.*.node_groups.*.name",
	"components_install_mode.*.node_groups.*.assigned_roles",
	"components_install_mode.*.component_user_password",
	"components_install_mode.*.component_default_password",
}

// @API MRS POST /v2/{project_id}/clusters/{cluster_id}/components
// @API MRS GET /v1.1/{project_id}/cluster_infos/{cluster_id}
func ResourceClusterComponentBatchAdd() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterComponentBatchAddCreate,
		ReadContext:   resourceClusterComponentBatchAddRead,
		UpdateContext: resourceClusterComponentBatchAddUpdate,
		DeleteContext: resourceClusterComponentBatchAddDelete,

		CustomizeDiff: config.FlexibleForceNew(clusterComponentBatchAddNonUpdateParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the components to be added are located.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the MRS cluster.`,
			},
			"components_install_mode": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"component": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the component.`,
						},
						"node_groups": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The name of the node group.`,
									},
									"assigned_roles": {
										Type:        schema.TypeList,
										Required:    true,
										Description: `The list of roles to be assigned to this node group.`,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
							Description: `The node groups where the component roles will be deployed.`,
						},
						"component_user_password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: `The password for the component user.`,
						},
						"component_default_password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: `The password for the component default user.`,
						},
					},
				},
				Description: `The list of components to be added.`,
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildClusterComponentBatchAddComponentInstallMode(components []interface{}) map[string]interface{} {
	if len(components) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(components))
	for _, com := range components {
		result = append(result, map[string]interface{}{
			"component": utils.PathSearch("component", com, nil),
			"node_groups": buildClusterComponentBatchAddNodeGroups(utils.PathSearch("node_groups", com,
				make([]interface{}, 0)).([]interface{})),
			"component_user_password":    utils.ValueIgnoreEmpty(utils.PathSearch("component_user_password", com, nil)),
			"component_default_password": utils.ValueIgnoreEmpty(utils.PathSearch("component_default_password", com, nil)),
		})
	}

	return map[string]interface{}{
		"components_install_mode": result,
	}
}

func buildClusterComponentBatchAddNodeGroups(nodeGroups []interface{}) []map[string]interface{} {
	if len(nodeGroups) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(nodeGroups))
	for _, nodeGroup := range nodeGroups {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", nodeGroup, nil),
			"assigned_roles": utils.ExpandToStringList(utils.PathSearch("assigned_roles", nodeGroup,
				make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func resourceClusterComponentBatchAddCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		httpUrl   = "v2/{project_id}/clusters/{cluster_id}/components"
		clusterId = d.Get("cluster_id").(string)
	)
	client, err := cfg.NewServiceClient("mrs", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", clusterId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildClusterComponentBatchAddComponentInstallMode(d.Get("components_install_mode").([]interface{}))),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("unable to add components to the cluster (%s): %s", clusterId, err)
	}

	err = waitForClusterStatusCompleted(ctx, client, clusterId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the cluster (%s) components to be added to complete: %s", clusterId, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID)

	return nil
}

func resourceClusterComponentBatchAddRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterComponentBatchAddUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterComponentBatchAddDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch add components to the MRS cluster. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

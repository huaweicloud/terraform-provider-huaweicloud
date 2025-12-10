package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace POST /v2/{project_id}/desktop-pools
// @API Workspace GET /v2/{project_id}/workspace-sub-jobs
// @API Workspace GET /v2/{project_id}/desktop-pools/{pool_id}
// @API Workspace GET /v2/{project_id}/desktop-pools/{pool_id}/users
// @API Workspace PUT /v2/{project_id}/desktop-pools/{pool_id}
// @API Workspace POST /v2/{project_id}/desktop-pools/{pool_id}/resize
// @API Workspace POST /v2/{project_id}/desktop-pools/{pool_id}/users
// @API Workspace POST /v2/{project_id}/desktop-pools/{pool_id}/volumes/batch-add
// @API Workspace POST /v2/{project_id}/desktop-pools/{pool_id}/volumes/batch-delete
// @API Workspace GET /v2/{project_id}/desktops
// @API Workspace POST /v2/{project_id}/desktops/batch-delete
// @API Workspace DELETE /v2/{project_id}/desktop-pools/{pool_id}

var nonUpdatableParams = []string{"type", "size", "image_type", "image_id", "root_volume",
	"root_volume.*.type", "root_volume.*.size", "subnet_ids",
	"vpc_id", "security_groups", "enterprise_project_id"}

func ResourceDesktopPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDesktopPoolCreate,
		ReadContext:   resourceDesktopPoolRead,
		UpdateContext: resourceDesktopPoolUpdate,
		DeleteContext: resourceDesktopPoolDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(nonUpdatableParams),
			config.MergeDefaultTags(),
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the desktop pool.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the desktop pool.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The number of the desktops under the desktop pool.`,
			},
			"product_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The specification ID of the desktop pool.`,
			},
			"image_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The image type of the desktop pool.`,
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The image ID of the desktop pool.`,
			},
			"root_volume": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        desktopPoolVolumeSchema(),
				Description: `The system volume configuration of the desktop pool.`,
			},
			"subnet_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of the subnet IDs to which the desktop pool belongs.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the VPC to which the desktop pool belongs.`,
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The ID of the security group.`,
						},
					},
				},
				Description: `The list of the security groups to which the desktop pool belongs.`,
			},
			"data_volumes": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        desktopPoolVolumeSchema(),
				Description: `The list of the data volume configurations of the desktop pool.`,
			},
			"authorized_objects": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The ID of the object.`,
						},
						"object_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the object.`,
						},
						"object_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the object.`,
						},
						"user_group": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The permission group to which the user belongs.`,
						},
					},
				},
				Description: `The list of the users or user groups to be authorized.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The availability zone to which the desktop pool belongs.`,
			},
			"disconnected_retention_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The desktops and users disconnection retention period under desktop pool, in minutes.`,
			},
			"enable_autoscale": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable elastic scaling of the desktop pool.`,
			},
			"autoscale_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"autoscale_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The type of automatic scaling policy.`,
						},
						"max_auto_created": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `The maximum number of automatically created desktops.`,
						},
						"min_idle": {
							Type:     schema.TypeInt,
							Optional: true,
							Description: `The desktops will be automatically created when the number of idle desktops is less than
							the specified value.`,
						},
						"once_auto_created": {
							Type:     schema.TypeInt,
							Optional: true,
							Description: utils.SchemaDesc(
								`The number of desktops automatically created at one time.`,
								utils.SchemaDescInput{
									Deprecated: true,
								},
							),
						},
					},
				},
				Description: `The automatic scaling policy of the desktop pool.`,
			},
			"desktop_name_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the policy to generate the desktop name.`,
			},
			"ou_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The OU name corresponding to the AD server.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the enterprise project to which the desktop pool belongs.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the desktop pool.`,
			},
			"tags": common.TagsSchema(`The key/value pairs to associate with the desktop pool.`),
			"in_maintenance_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable maintenance mode of the desktop pool.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the desktop pool.`,
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the desktop pool, in UTC format.`,
			},
			"desktop_used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of desktops associated with the users under the desktop pool.`,
			},
			"product": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flavor_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The product specification ID of the desktop pool.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The product type of the desktop pool.`,
						},
						"cpu": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The product CPU of the desktop pool.`,
						},
						"memory": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The product memory of the desktop pool.`,
						},
						"descriptions": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The product description of the desktop pool.`,
						},
						"charging_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The product charging mode of the desktop pool.`,
						},
					},
				},
				Description: `The product information of the desktop pool.`,
			},
			"image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The image name of the desktop pool.`,
			},
			"image_os_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The image OS type of the desktop pool.`,
			},
			"image_os_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The image OS version of the desktop pool.`,
			},
			"image_os_platform": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The image OS platform of the desktop pool.`,
			},
			// Internal attribute(s).
			"data_volumes_order": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The origin data volume size.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The origin data volume type.`,
						},
					},
				},
				Description: utils.SchemaDesc(
					`The origin list of data volume configuration that used to reorder the 'data_volumes' parameter.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func desktopPoolVolumeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the volume.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The size of the volume, in GB.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the volume.`,
			},
		},
	}
}

func buildCreateDesktopPoolBodyParam(epsId string, d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required.
		"name":        d.Get("name"),
		"type":        d.Get("type"),
		"size":        d.Get("size"),
		"product_id":  d.Get("product_id"),
		"image_type":  d.Get("image_type"),
		"image_id":    d.Get("image_id"),
		"root_volume": buildDesktopPoolRootVolume(d.Get("root_volume.0")),
		"subnet_ids":  d.Get("subnet_ids"),
		// Optional.
		"vpc_id":                        utils.ValueIgnoreEmpty(d.Get("vpc_id")),
		"security_groups":               buildDesktopPoolSecurityGroups(d.Get("security_groups").(*schema.Set)),
		"availability_zone":             utils.ValueIgnoreEmpty(d.Get("availability_zone")),
		"data_volumes":                  buildDesktopPoolDataVolumes(d.Get("data_volumes").([]interface{})),
		"authorized_objects":            buildDesktopPoolAuthorizedObjects(d.Get("authorized_objects").(*schema.Set)),
		"disconnected_retention_period": utils.ValueIgnoreEmpty(d.Get("disconnected_retention_period")),
		"enable_autoscale":              d.Get("enable_autoscale"),
		"autoscale_policy":              buildDesktopPoolAutoScalePolicy(d.Get("autoscale_policy").([]interface{})),
		"desktop_name_policy_id":        utils.ValueIgnoreEmpty(d.Get("desktop_name_policy_id")),
		"ou_name":                       utils.ValueIgnoreEmpty(d.Get("ou_name")),
		"tags":                          utils.ValueIgnoreEmpty(utils.ExpandResourceTags(d.Get("tags").(map[string]interface{}))),
		"enterprise_project_id":         utils.ValueIgnoreEmpty(epsId),
		"description":                   utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func buildDesktopPoolRootVolume(rootVolume interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type": utils.PathSearch("type", rootVolume, nil),
		"size": utils.PathSearch("size", rootVolume, nil),
	}
}

func buildDesktopPoolDataVolumes(dataVolumes []interface{}) []map[string]interface{} {
	if len(dataVolumes) == 0 {
		return nil
	}

	rest := make([]map[string]interface{}, len(dataVolumes))
	for i, v := range dataVolumes {
		rest[i] = map[string]interface{}{
			"type": utils.PathSearch("type", v, nil),
			"size": utils.PathSearch("size", v, nil),
		}
	}
	return rest
}

func buildDesktopPoolSecurityGroups(secutityGroups *schema.Set) []map[string]interface{} {
	if secutityGroups.Len() == 0 {
		return nil
	}
	rest := make([]map[string]interface{}, secutityGroups.Len())
	for i, v := range secutityGroups.List() {
		rest[i] = map[string]interface{}{
			"id": utils.PathSearch("id", v, nil),
		}
	}
	return rest
}

func buildDesktopPoolAuthorizedObjects(authorizedObject *schema.Set) []map[string]interface{} {
	if authorizedObject.Len() == 0 {
		return nil
	}

	rest := make([]map[string]interface{}, authorizedObject.Len())
	for i, v := range authorizedObject.List() {
		rest[i] = map[string]interface{}{
			"object_id":   utils.PathSearch("object_id", v, nil),
			"object_type": utils.PathSearch("object_type", v, nil),
			"object_name": utils.PathSearch("object_name", v, nil),
			"user_group":  utils.PathSearch("user_group", v, nil),
		}
	}
	return rest
}

func buildDesktopPoolAutoScalePolicy(autoScalePolicy []interface{}) map[string]interface{} {
	if len(autoScalePolicy) == 0 || autoScalePolicy[0] == nil {
		return map[string]interface{}{}
	}
	policy := autoScalePolicy[0]
	return map[string]interface{}{
		"autoscale_type":    utils.ValueIgnoreEmpty(utils.PathSearch("autoscale_type", policy, nil)),
		"max_auto_created":  utils.ValueIgnoreEmpty(utils.PathSearch("max_auto_created", policy, nil)),
		"min_idle":          utils.ValueIgnoreEmpty(utils.PathSearch("min_idle", policy, nil)),
		"once_auto_created": utils.ValueIgnoreEmpty(utils.PathSearch("once_auto_created", policy, nil)),
	}
}

func resourceDesktopPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		httpUrl         = "v2/{project_id}/desktop-pools"
		desktopPoolName = d.Get("name").(string)
	)

	client, err := cfg.NewServiceClient("workspace", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDesktopPoolBodyParam(cfg.GetEnterpriseProjectID(d), d)),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating desktop pool: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to job ID from API response")
	}
	// 1. To prevent residual resources on the cloud, first determine whether the resource pool has been created,
	//    and then determine whether the desktop under it has been created.
	// 2. If you use the job_id obtained from the creation interface to call the query job interface, you can only determine whether
	//    the desktop under the desktop pool is created successfully, but you cannot obtain the ID of the desktop pool, resulting in
	//    the inability to set the resource ID. Therefore, the query desktop pool list interface is called here in a loop to determine
	//    whether the desktop pool is created successfully.
	desktopPoolId, err := waitForWorkspacePoolStatusCompleted(ctx, client, desktopPoolName, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the creation desktop pool (%s) to complete: %s", desktopPoolName, err)
	}

	if desktopPoolId == "" {
		return diag.Errorf("unable to find desktop pool ID from API response")
	}

	d.SetId(desktopPoolId)

	if err = d.Set("data_volumes_order", buildDesktopPoolDataVolumesOrder(d)); err != nil {
		log.Printf("[ERROR] error setting the data_volumes_order field after creating desktop pool: %s", err)
	}

	// The successful creation of a desktop pool does not mean that all desktops under it are successfully created.
	_, err = waitForWorkspaceResourcePoolJobCompleted(ctx, client, desktopPoolId, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the job (%s) completed: %s", jobId, err)
	}

	if v, ok := d.GetOk("in_maintenance_mode"); ok {
		updateOpt := map[string]interface{}{
			"in_maintenance_mode": v,
		}
		err = updateDesktopPool(client, desktopPoolId, updateOpt)
		if err != nil {
			return diag.Errorf("error enabling maintenance mode of the desktop pool (%s): %s", desktopPoolId, err)
		}
	}
	return resourceDesktopPoolRead(ctx, d, meta)
}

func buildDesktopPoolDataVolumesOrder(d *schema.ResourceData) []interface{} {
	dataVolumes, ok := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "data_volumes").([]interface{})
	if !ok || dataVolumes == nil {
		return nil
	}

	result := make([]interface{}, 0, len(dataVolumes))
	for _, dataVolume := range dataVolumes {
		result = append(result, map[string]interface{}{
			"size": utils.PathSearch("size", dataVolume, nil),
			"type": utils.PathSearch("type", dataVolume, nil),
		})
	}

	return result
}

func waitForWorkspaceResourcePoolJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, desktopPoolId, jobId string,
	timeout time.Duration) (string, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"WAITING", "RUNNING"},
		Target:       []string{"SUCCESS"},
		Refresh:      refreshWorkspaceJobFunc(client, jobId, fmt.Sprintf("&desktop_pool_id=%s", desktopPoolId)),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 15 * time.Second,
	}

	resp, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "", err
	}

	return utils.PathSearch("entities.desktop_id", resp, "").(string), nil
}

func waitForWorkspacePoolStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, desktopPoolName string,
	timeout time.Duration) (string, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshWorkspacePoolJobStatus(client, desktopPoolName),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 15 * time.Second,
	}

	desktopPool, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "", err
	}

	return utils.PathSearch("id", desktopPool, "").(string), nil
}

func getDesktopPoolbyName(client *golangsdk.ServiceClient, desktopPoolName string) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/desktop-pools"
		offset  = 0
		opt     = golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	// The 'name' parameter is fuzzy search.
	listPath = fmt.Sprintf("%s?name=%s", listPath, desktopPoolName)

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving desktop pools: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		desktopPools := utils.PathSearch("desktop_pools", respBody, make([]interface{}, 0)).([]interface{})
		if len(desktopPools) < 1 {
			break
		}

		desktopPool := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0]", desktopPoolName), desktopPools, nil)
		if desktopPool != nil {
			return desktopPool, nil
		}

		offset += len(desktopPools)
	}
	return nil, golangsdk.ErrDefault404{}
}

func refreshWorkspacePoolJobStatus(client *golangsdk.ServiceClient, desktopPoolName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		desktopPool, err := getDesktopPoolbyName(client, desktopPoolName)
		if err != nil {
			return desktopPool, "ERROR", err
		}

		status := utils.PathSearch("status", desktopPool, "").(string)
		if status == "STEADY" {
			return desktopPool, "COMPLETED", nil
		}

		if status == "ERROR" {
			return desktopPool, "ERROR", fmt.Errorf("unexpect status (%s)", status)
		}

		return desktopPool, "PENDING", nil
	}
}

func GetDesktopPoolById(client *golangsdk.ServiceClient, desktopPoolId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/desktop-pools/{pool_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{pool_id}", desktopPoolId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	respBody, err := client.Request("GET", getPath, &opt)
	if err != nil {
		// WKS.0001: The desktop pool does not exist.
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", "WKS.0001")
	}

	desktopPool, err := utils.FlattenResponse(respBody)
	if err != nil {
		return nil, err
	}
	return desktopPool, nil
}

func getAssociatedObjectsById(client *golangsdk.ServiceClient, desktopPoolId string) ([]interface{}, error) {
	var (
		// The 'limit' default value is 10.
		httpUrl = "v2/{project_id}/desktop-pools/{pool_id}/users?limit=100"
		offset  = 0
		result  = make([]interface{}, 0)
		opt     = golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{pool_id}", desktopPoolId)
	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}
		objects := utils.PathSearch("objects", respBody, make([]interface{}, 0)).([]interface{})
		if len(objects) < 1 {
			break
		}

		result = append(result, objects...)
		offset += len(objects)
	}
	return result, nil
}

func resourceDesktopPoolRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		desktopPoolId = d.Id()
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	desktopPool, err := GetDesktopPoolById(client, desktopPoolId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving desktop pool (%s)", desktopPoolId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", desktopPool, nil)),
		d.Set("type", utils.PathSearch("type", desktopPool, nil)),
		d.Set("size", utils.PathSearch("desktop_count", desktopPool, nil)),
		d.Set("product_id", utils.PathSearch("product.product_id", desktopPool, nil)),
		d.Set("image_id", utils.PathSearch("image_id", desktopPool, nil)),
		d.Set("root_volume", flattenDesktopPoolRootVolume(utils.PathSearch("root_volume", desktopPool, nil))),
		d.Set("subnet_ids", flattenDesktopPoolSubnet(utils.PathSearch("subnet_id", desktopPool, nil))),
		d.Set("security_groups", flattenDesktopPoolSecurityGroups(utils.PathSearch("security_groups", desktopPool,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("availability_zone", utils.PathSearch("availability_zone", desktopPool, nil)),
		d.Set("data_volumes", flattenDesktopPoolDataVolume(utils.PathSearch("data_volumes",
			desktopPool, make([]interface{}, 0)).([]interface{}),
			d.Get("data_volumes_order").([]interface{}))),
		d.Set("disconnected_retention_period", utils.PathSearch("disconnected_retention_period", desktopPool, nil)),
		d.Set("enable_autoscale", utils.PathSearch("enable_autoscale", desktopPool, nil)),
		d.Set("autoscale_policy", flattenDesktopPoolAutoScalePolicy(utils.PathSearch("autoscale_policy", desktopPool, nil))),
		d.Set("desktop_name_policy_id", utils.PathSearch("desktop_name_policy_id", desktopPool, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", desktopPool, nil)),
		d.Set("description", utils.PathSearch("description", desktopPool, nil)),
		d.Set("in_maintenance_mode", utils.PathSearch("in_maintenance_mode", desktopPool, nil)),
		// Attributes.
		d.Set("status", utils.PathSearch("status", desktopPool, nil)),
		d.Set("created_time", utils.PathSearch("created_time", desktopPool, nil)),
		d.Set("desktop_used", utils.PathSearch("desktop_used", desktopPool, nil)),
		d.Set("product", flattenDesktopPoolProduct(utils.PathSearch("product", desktopPool, nil))),
		d.Set("image_name", utils.PathSearch("image_name", desktopPool, nil)),
		d.Set("image_os_type", utils.PathSearch("image_os_type", desktopPool, nil)),
		d.Set("image_os_version", utils.PathSearch("image_os_version", desktopPool, nil)),
		d.Set("image_os_platform", utils.PathSearch("image_os_platform", desktopPool, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", desktopPool, make([]interface{}, 0)).([]interface{}))),
	)

	authorizedObjects, err := getAssociatedObjectsById(client, desktopPoolId)
	if err != nil {
		// To prevent resource errors caused by the interface not being online in some regions, use log to record the error.
		log.Printf("[WARN] error retrieving associated users under specified desktop pool (%s): %s", desktopPoolId, err)
	}
	mErr = multierror.Append(mErr,
		d.Set("authorized_objects", flattenDesktopPoolAuthorizedObjects(authorizedObjects)))

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDesktopPoolRootVolume(volume interface{}) []map[string]interface{} {
	if volume == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"type": utils.PathSearch("type", volume, nil),
			"size": utils.PathSearch("size", volume, nil),
			"id":   utils.PathSearch("id", volume, nil),
		},
	}
}

func flattenDesktopPoolSubnet(subnetId interface{}) []interface{} {
	if subnetId == nil {
		return nil
	}
	return []interface{}{subnetId}
}

func flattenDesktopPoolDataVolume(volumes, volumesOrderOrigin []interface{}) []map[string]interface{} {
	if len(volumes) == 0 {
		return nil
	}

	sortedVolumes := orderDataVolumesByDataVolumesOrderOrigin(volumes, volumesOrderOrigin)
	result := make([]map[string]interface{}, 0, len(sortedVolumes))
	for _, volume := range sortedVolumes {
		result = append(result, map[string]interface{}{
			"type": utils.PathSearch("type", volume, nil),
			"size": utils.PathSearch("size", volume, nil),
			"id":   utils.PathSearch("id", volume, nil),
		})
	}
	return result
}

func orderDataVolumesByDataVolumesOrderOrigin(volumes, volumesOrderOrigin []interface{}) []interface{} {
	if len(volumesOrderOrigin) == 0 {
		return volumes
	}

	sortedVolumes := make([]interface{}, 0)
	for _, volumeOrigin := range volumesOrderOrigin {
		index := findVolumeBySizeAndType(volumes, int64(utils.PathSearch("size", volumeOrigin, 0).(int)),
			utils.PathSearch("type", volumeOrigin, "").(string))
		if index == -1 {
			continue
		}

		// Add the found volume configuration to the sorted volumes list.
		sortedVolumes = append(sortedVolumes, volumes[index])
		// Remove the processed volume configuration from the remote volumes array.
		volumes = append(volumes[:index], volumes[index+1:]...)
	}

	// Add unsorted volumes configurations to the end of the sorted list.
	sortedVolumes = append(sortedVolumes, volumes...)
	log.Printf("[DEBUG] data_volumes sort result by data_volumes_order: %#v", sortedVolumes)
	return sortedVolumes
}

func findVolumeBySizeAndType(volumes []interface{}, originSize int64, originType string) int {
	for index, volume := range volumes {
		if int64(utils.PathSearch("size", volume, float64(0)).(float64)) == originSize &&
			utils.PathSearch("type", volume, "").(string) == originType {
			return index
		}
	}

	return -1
}

func flattenDesktopPoolAuthorizedObjects(authorizedObjects []interface{}) []map[string]interface{} {
	if len(authorizedObjects) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(authorizedObjects))
	for i, object := range authorizedObjects {
		result[i] = map[string]interface{}{
			"object_id":   utils.PathSearch("object_id", object, nil),
			"object_type": utils.PathSearch("object_type", object, nil),
			"object_name": utils.PathSearch("object_name", object, nil),
			"user_group":  utils.PathSearch("user_group", object, nil),
		}
	}
	return result
}

func flattenDesktopPoolSecurityGroups(securityGroups []interface{}) []map[string]interface{} {
	if len(securityGroups) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(securityGroups))
	for i, sg := range securityGroups {
		result[i] = map[string]interface{}{
			"id": utils.PathSearch("id", sg, nil),
		}
	}
	return result
}

func flattenDesktopPoolAutoScalePolicy(autoScalePolicy interface{}) []map[string]interface{} {
	if autoScalePolicy == nil {
		return nil
	}
	autoScalePolicyObject := autoScalePolicy.(map[string]interface{})
	// In some regions, the once_auto_created field is still returned even after the auto-scaling policy is canceled.
	if len(autoScalePolicyObject) == 0 || (len(autoScalePolicyObject) == 1 && autoScalePolicyObject["once_auto_created"] != nil) {
		return nil
	}

	return []map[string]interface{}{
		{
			"autoscale_type":    utils.PathSearch("autoscale_type", autoScalePolicy, nil),
			"max_auto_created":  utils.PathSearch("max_auto_created", autoScalePolicy, nil),
			"min_idle":          utils.PathSearch("min_idle", autoScalePolicy, nil),
			"once_auto_created": utils.PathSearch("once_auto_created", autoScalePolicy, nil),
		},
	}
}

func flattenDesktopPoolProduct(product interface{}) []map[string]interface{} {
	if product == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"flavor_id":     utils.PathSearch("flavor_id", product, nil),
			"type":          utils.PathSearch("type", product, nil),
			"cpu":           utils.PathSearch("cpu", product, nil),
			"memory":        utils.PathSearch("memory", product, nil),
			"descriptions":  utils.PathSearch("descriptions", product, nil),
			"charging_mode": utils.PathSearch("charge_mode", product, nil),
		},
	}
}

func buildUpdateDesktopPoolBodyParam(d *schema.ResourceData) map[string]interface{} {
	// availability_zone, ou_name, desktop_name_policy_id, description must be set to an empty string to be changed to an empty value.
	param := map[string]interface{}{
		"availability_zone":             d.Get("availability_zone"),
		"disconnected_retention_period": utils.ValueIgnoreEmpty(d.Get("disconnected_retention_period")),
		"enable_autoscale":              d.Get("enable_autoscale"),
		"ou_name":                       d.Get("ou_name"),
		"desktop_name_policy_id":        d.Get("desktop_name_policy_id"),
		"tags":                          utils.ValueIgnoreEmpty(utils.ExpandResourceTags(d.Get("tags").(map[string]interface{}))),
		"description":                   d.Get("description"),
		"in_maintenance_mode":           d.Get("in_maintenance_mode"),
	}
	params := utils.RemoveNil(param)
	params["autoscale_policy"] = buildDesktopPoolAutoScalePolicy(d.Get("autoscale_policy").([]interface{}))
	return params
}

func updateDesktopPool(client *golangsdk.ServiceClient, desktopPoolId string, params map[string]interface{}) error {
	httpUrl := "v2/{project_id}/desktop-pools/{pool_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{pool_id}", desktopPoolId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         params,
	}
	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func resourceDesktopPoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		desktopPoolId = d.Id()
		updateTimeout = d.Timeout(schema.TimeoutUpdate)
	)
	client, err := cfg.NewServiceClient("workspace", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	if d.HasChanges("name", "availability_zone", "disconnected_retention_period", "enable_autoscale", "autoscale_policy",
		"ou_name", "desktop_name_policy_id", "tags", "description", "in_maintenance_mode") {
		updateOpt := buildUpdateDesktopPoolBodyParam(d)
		if d.HasChange("name") {
			updateOpt["name"] = d.Get("name")
		}

		if err := updateDesktopPool(client, desktopPoolId, updateOpt); err != nil {
			return diag.Errorf("error updating desktop pool: %s", err)
		}
	}

	if d.HasChange("data_volumes") {
		oldRaw, _ := d.GetChange("data_volumes")
		oldDataVolumes := oldRaw.([]interface{})
		addVolumes, rmVolumes := getDesktopPoolDataVolumesDiff(d, oldDataVolumes)
		log.Printf("[DEBUG] want to add data volumes: %#v", addVolumes)
		log.Printf("[DEBUG] want to remove data volumes: %#v", rmVolumes)
		if len(rmVolumes) > 0 {
			if err = removeDesktopPoolDataVolumes(ctx, client, desktopPoolId, rmVolumes, updateTimeout); err != nil {
				return diag.Errorf("error removing data volumes of the desktop pool (%s): %s", desktopPoolId, err)
			}
		}

		if len(addVolumes) > 0 {
			if err = addDesktopPoolDataVolumes(ctx, client, desktopPoolId, addVolumes, updateTimeout); err != nil {
				return diag.Errorf("error adding data volumes of the desktop pool (%s): %s", desktopPoolId, err)
			}
		}

		if err = d.Set("data_volumes_order", buildDesktopPoolDataVolumesOrder(d)); err != nil {
			log.Printf("[ERROR] error setting the 'data_volumes_order' field after updating data volumes: %s", err)
		}
	}

	if d.HasChange("product_id") {
		err = updateDesktopPoolProductId(ctx, client, desktopPoolId, d.Get("product_id").(string), updateTimeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("authorized_objects") {
		oldRaw, newRaw := d.GetChange("authorized_objects")
		addObjects := newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))
		rmObjects := oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))
		if rmObjects.Len() > 0 {
			if err := updateDesktopPoolAuthorizedObjects(client, desktopPoolId, "REMOVE", rmObjects); err != nil {
				return diag.Errorf("error removing authorized objects: %s", err)
			}
		}

		if addObjects.Len() > 0 {
			if err := updateDesktopPoolAuthorizedObjects(client, desktopPoolId, "ADD", addObjects); err != nil {
				return diag.Errorf("error adding authorized objects: %s", err)
			}
		}
	}

	return resourceDesktopPoolRead(ctx, d, meta)
}

// Get the value from the configuration file, record it as `configDataVolumes`, and iterate through `oldDataVolumes`.
// 1. If an element in `oldDataVolumes` exists in `configDataVolumes`, it means the element has not changed,
// so remove it from both `configDataVolumes` and `oldDataVolumes`.
// 2. If an element in `configDataVolumes` does not exist in `oldDataVolumes`, it means the element needs to be added,
// so keep it in `configDataVolumes`.
// 3. Finally return the remaining elements in `configDataVolumes` and `oldDataVolumes`,
// where `configDataVolumes` are elements to be added and `oldDataVolumes` are elements to be removed.
func getDesktopPoolDataVolumesDiff(d *schema.ResourceData, oldDataVolumes []interface{}) (addVolumes, rmVolumes []interface{}) {
	configDataVolumes, ok := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "data_volumes").([]interface{})
	if !ok || configDataVolumes == nil {
		return nil, nil
	}

	// The oldVolumeIndexes records the indices of elements in oldDataVolumes that do not need to be modified.
	oldVolumeIndexes := make([]int, 0)
	for i, oldVolume := range oldDataVolumes {
		configIndex := findVolumeBySizeAndType(
			configDataVolumes,
			int64(utils.PathSearch("size", oldVolume, 0).(int)),
			utils.PathSearch("type", oldVolume, "").(string),
		)
		if configIndex == -1 {
			continue
		}

		configDataVolumes = append(configDataVolumes[:configIndex], configDataVolumes[configIndex+1:]...)
		oldVolumeIndexes = append(oldVolumeIndexes, i)
	}

	if len(oldVolumeIndexes) > 0 {
		// Sort indices in descending order to delete from back to front.
		// This avoids index shifting when deleting elements.
		sort.Sort(sort.Reverse(sort.IntSlice(oldVolumeIndexes)))
		// Delete elements from `oldDataVolumes` at the corresponding indexes.
		for _, idx := range oldVolumeIndexes {
			if idx < 0 || idx >= len(oldDataVolumes) {
				// Boundary check to avoid panic.
				continue
			}
			oldDataVolumes = append(oldDataVolumes[:idx], oldDataVolumes[idx+1:]...)
		}
	}

	rmVolumes = oldDataVolumes
	addVolumes = configDataVolumes
	return
}

func addDesktopPoolDataVolumes(ctx context.Context, client *golangsdk.ServiceClient, desktopPoolId string, dataVolumes []interface{},
	timeout time.Duration) error {
	// When the data volumes are added, all desktops in the desktop pool must be running.
	err := waitForDesktopsUnderDesktopPoolStatusCompleted(ctx, client, desktopPoolId, timeout)
	if err != nil {
		return err
	}

	httpUrl := "v2/{project_id}/desktop-pools/{pool_id}/volumes/batch-add"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{pool_id}", desktopPoolId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"volumes": buildDesktopPoolDataVolumes(dataVolumes),
		},
	}

	resp, err := client.Request("POST", updatePath, &opt)
	if err != nil {
		return err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("unable to find job ID under the desktop pool (%s) from API response", desktopPoolId)
	}

	err = waitForDesktopPoolJobCompleted(ctx, client, jobId, timeout)
	return err
}

func waitForDesktopsUnderDesktopPoolStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, desktopPoolId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshDesktopsUnderDesktopPoolStatusFunc(client, desktopPoolId),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshDesktopsUnderDesktopPoolStatusFunc(client *golangsdk.ServiceClient, desktopPoolId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		desktops, err := getDesktopsUnderDesktopPoolById(client, desktopPoolId)
		if err != nil {
			return desktops, "ERROR", err
		}

		statuses := utils.PathSearch("[*].status", desktops, make([]interface{}, 0)).([]interface{})
		// ACTIVE：Running
		// SHUTOFF：Shutdown
		// ERROR：Error
		// HIBERNATED：Hibernate
		for _, v := range statuses {
			status := v.(string)
			if status == "ERROR" {
				return desktops, "ERROR", fmt.Errorf("unexpect status: %s", status)
			}

			if status != "ACTIVE" {
				return desktops, "PENDING", nil
			}
		}

		return desktops, "COMPLETED", nil
	}
}

func buildRemoveDesktopPoolDataVolumes(dataVolumes []interface{}) []map[string]interface{} {
	if len(dataVolumes) == 0 {
		return nil
	}

	rest := make([]map[string]interface{}, 0, len(dataVolumes))
	for _, v := range dataVolumes {
		rest = append(rest, map[string]interface{}{
			"id":   utils.PathSearch("id", v, nil),
			"type": utils.PathSearch("type", v, nil),
			"size": utils.PathSearch("size", v, nil),
		})
	}

	return rest
}

func removeDesktopPoolDataVolumes(ctx context.Context, client *golangsdk.ServiceClient, desktopPoolId string, dataVolumes []interface{},
	timeout time.Duration) error {
	httpUrl := "v2/{project_id}/desktop-pools/{pool_id}/volumes/batch-delete"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{pool_id}", desktopPoolId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"volumes": buildRemoveDesktopPoolDataVolumes(dataVolumes),
		},
	}

	// There may be desktops under the desktop pool that are in the process of operation,
	resp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx: ctx,
		RetryFunc: func() (interface{}, bool, error) {
			res, err := client.Request("POST", updatePath, &opt)
			retry, err := handleOperationDesktopPoolError409(err)
			return res, retry, err
		},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 20 * time.Second,
	})

	if err != nil {
		return err
	}

	respBody, err := utils.FlattenResponse(resp.(*http.Response))
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("unable to find job ID under the desktop pool (%s) from API response", desktopPoolId)
	}

	err = waitForDesktopPoolJobCompleted(ctx, client, jobId, timeout)
	return err
}

func updateDesktopPoolProductId(ctx context.Context, client *golangsdk.ServiceClient, desktopPoolId, productId string,
	timeout time.Duration) error {
	httpUrl := "v2/{project_id}/desktop-pools/{pool_id}/resize"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{pool_id}", desktopPoolId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"product_id": productId,
			// If the desktop pool contains desktops that are not powered off, the specification change for the corresponding desktop
			//  will fail. Therefore, 'mode' is set to 'STOP_DESKTOP'.
			// `STOP_DESKTOP`: If the desktops under the desktop pool are powered on, they will be powered off first before
			// the specification change.
			"mode": "STOP_DESKTOP",
		},
	}

	// Make sure all desktops under the desktop pool not processing jobs before updating the product ID.
	resp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx: ctx,
		RetryFunc: func() (interface{}, bool, error) {
			res, err := client.Request("POST", updatePath, &opt)
			retry, err := handleOperationDesktopPoolError409(err)
			return res, retry, err
		},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 20 * time.Second,
	})

	if err != nil {
		return fmt.Errorf("error updating product ID of the desktop pool (%s): %s", desktopPoolId, err)
	}

	respBody, err := utils.FlattenResponse(resp.(*http.Response))
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("unable to find job ID from API response under the desktop pool (%s)", desktopPoolId)
	}

	err = waitForDesktopPoolJobCompleted(ctx, client, jobId, timeout)
	if err != nil {
		return fmt.Errorf("error waiting for the desktops status running: %s", err)
	}

	return nil
}

func handleOperationDesktopPoolError409(err error) (bool, error) {
	if err == nil {
		return false, nil
	}

	if err409, ok := err.(golangsdk.ErrDefault409); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(err409.Body, &apiError); jsonErr != nil {
			return false, jsonErr
		}

		errCode := utils.PathSearch("error_code", apiError, "")
		// WKS.00010032: Operation conflict. The desktop current instance status is [xxx] and deny operation [xxx], resource id [xxx].
		if errCode == "WKS.00010032" {
			return true, err
		}
	}
	return false, err
}

func waitForDesktopPoolJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, jobId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshJobStatusFunc(client, jobId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 20 * time.Second,
		// After the job is completed, the desktop status does not change immediately.
		// Therefore, set the ContinuousTargetOccurence to 2 to ensure the desktop status is changed.
		ContinuousTargetOccurence: 2,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func updateDesktopPoolAuthorizedObjects(client *golangsdk.ServiceClient, desktopPoolId, action string, authorizedObjects *schema.Set) error {
	httpUrl := "v2/{project_id}/desktop-pools/{pool_id}/users"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{pool_id}", desktopPoolId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"action":  action,
			"objects": buildDesktopPoolAuthorizedObjects(authorizedObjects),
		},
	}
	_, err := client.Request("POST", updatePath, &opt)
	return err
}

func resourceDesktopPoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("workspace", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	// Before deleting the desktop pool, you must disable the automatic creation function under it.
	desktopPoolId := d.Id()
	desktopPool, err := GetDesktopPoolById(client, desktopPoolId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving desktop pool")
	}

	if utils.PathSearch("enable_autoscale", desktopPool, false).(bool) {
		updateOpt := map[string]interface{}{
			"enable_autoscale": false,
		}
		if err := updateDesktopPool(client, desktopPoolId, updateOpt); err != nil {
			return diag.Errorf("error disabling the automatic creation function of the desktop pool (%s): %s", desktopPoolId, err)
		}
	}

	// Before deleting the desktop pool, you must delete all desktops under it.
	desktops, err := getDesktopsUnderDesktopPoolById(client, desktopPoolId)
	desktopIds := utils.PathSearch("[*].desktop_id", desktops, make([]interface{}, 0)).([]interface{})
	if err != nil {
		return diag.Errorf("error retrieving desktop IDs under specified desktop pool (%s): %s", desktopPoolId, err)
	}

	if len(desktopIds) != 0 {
		jobId, err := deleteDesktopsUnderDesktopPool(client, desktopIds)
		if err != nil {
			return diag.Errorf("error deleting desktops under specified desktop pool (%s): %s", desktopPoolId, err)
		}

		_, err = waitForWorkspaceJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for the job (%s) completed: %s", jobId, err)
		}
	}

	err = deleteDesktopPoolById(client, desktopPoolId)
	if err != nil {
		return diag.Errorf("error deleting desktop pool (%s): %s", desktopPoolId, err)
	}
	return nil
}

func getDesktopsUnderDesktopPoolById(client *golangsdk.ServiceClient, desktopPoolId string) ([]interface{}, error) {
	var (
		// The 'limit' default value is 1000.
		httpUrl = "v2/{project_id}/desktops"
		offset  = 0
		result  = make([]interface{}, 0)
		opt     = golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		}
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?pool_id=%s", listPath, desktopPoolId)
	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}
		desktops := utils.PathSearch("desktops", respBody, make([]interface{}, 0)).([]interface{})
		if len(desktops) < 1 {
			break
		}

		result = append(result, desktops...)
		offset += len(desktops)
	}
	return result, nil
}

func deleteDesktopsUnderDesktopPool(client *golangsdk.ServiceClient, desktopIds interface{}) (string, error) {
	httpUrl := "v2/{project_id}/desktops/batch-delete"
	deleteDesktopsPath := client.Endpoint + httpUrl
	deleteDesktopsPath = strings.ReplaceAll(deleteDesktopsPath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"desktop_ids": desktopIds,
		},
	}
	resp, err := client.Request("POST", deleteDesktopsPath, &opt)
	if err != nil {
		return "", err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return "", err
	}

	return utils.PathSearch("job_id", respBody, "").(string), nil
}

func deleteDesktopPoolById(client *golangsdk.ServiceClient, desktopPoolId string) error {
	httpUrl := "v2/{project_id}/desktop-pools/{pool_id}"
	deleteDesktopsPath := client.Endpoint + httpUrl
	deleteDesktopsPath = strings.ReplaceAll(deleteDesktopsPath, "{project_id}", client.ProjectID)
	deleteDesktopsPath = strings.ReplaceAll(deleteDesktopsPath, "{pool_id}", desktopPoolId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err := client.Request("DELETE", deleteDesktopsPath, &opt)
	return err
}

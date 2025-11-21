package workspace

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

// @API Workspace POST /v1/{project_id}/app-server-groups
// @API Workspace GET /v1/{project_id}/app-server-group
// @API Workspace PATCH /v1/{project_id}/app-server-groups/{server_group_id}
// @API Workspace POST /v1/{project_id}/app-server-groups/{server_group_id}
// @API Workspace POST /v1/{project_id}/server-group/tags/batch-create
// @API Workspace DELETE /v1/{project_id}/server-group/tags/batch-delete
func ResourceAppServerGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppServerGroupCreate,
		ReadContext:   resourceAppServerGroupRead,
		UpdateContext: resourceAppServerGroupUpdate,
		DeleteContext: resourceAppServerGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

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
				Description: `The name of the server group.`,
			},
			"os_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The operating system type of the server group.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The flavor ID of the server group.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The VPC ID to which the server group belongs.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The subnet ID to which the server group belongs.`,
			},
			"system_disk_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of system disk.`,
			},
			"system_disk_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The size of system disk, in GB.`,
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The image ID of the server group.`,
			},
			"image_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The image type of the server group.`,
			},
			"image_product_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The image product ID of the server group.`,
			},
			"is_vdi": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `The session mode of the server group.`,
			},
			"app_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The type of application group associated with the server group.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The availability zone of the server group.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the server group.`,
			},
			"tags": common.TagsSchema("The key/value pairs to associate with the server group."),
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the enterprise project to which the server group belong.`,
			},
			"ip_virtual": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Required:    true,
							ForceNew:    true,
							Description: "Whether to enable IP virtualization.",
						},
					},
				},
				Description: `The IP virtualization function configuration.`,
			},
			// Unable to obtain through the query interface.
			"route_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_session": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The number of session connections of the server.",
						},
						"cpu_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The CPU usage of the server.",
						},
						"mem_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The memory usage of the server.",
						},
					},
				},
				Description: `The session scheduling policy of the server group.`,
			},
			"ou_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The OU name corresponding to the AD server.`,
			},
			"extra_session_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The additional session type.`,
			},
			"extra_session_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `The number of additional sessions for a single server.`,
			},
			"primary_server_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The ID of the primary server group.`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether to enable server group.`,
			},
			"storage_mount_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The NAS storage directory mounting policy on the APS.`,
			},
			"project_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the project.`,
			},
			"image_min_disk": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The minimum memory required to run the image, in MB.`,
			},
			"flavors": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of server flavors.`,
				Elem:        v1WorkspaceAppServersGroupFlavorSchema(),
			},
		},
	}
}

func v1WorkspaceAppServersGroupFlavorSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the flavor.`,
			},
			"links": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The quick link information for relevant tags corresponding to server specifications.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The corresponding shortcut link.`,
						},
						"rel": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The shortcut link tag name.`,
						},
					},
				},
			},
		},
	}
}

func resourceAppServerGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	httpUrl := "v1/{project_id}/app-server-groups"
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	epsId := cfg.GetEnterpriseProjectID(d)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateServerGroupBodyParams(d, epsId)),
	}
	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Workspace APP server group: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	serverGroupId := utils.PathSearch("id", respBody, "").(string)
	if serverGroupId == "" {
		return diag.Errorf("unable to find server group ID from API response")
	}
	d.SetId(serverGroupId)

	if v, ok := d.GetOk("storage_mount_policy"); ok {
		updateMountingPolicyOpt := map[string]interface{}{
			"storage_mount_policy": v,
		}
		if err := updateAppServerGroup(client, serverGroupId, updateMountingPolicyOpt); err != nil {
			return diag.Errorf("error updating the mounting policy of the server group (%s): %s", serverGroupId, err)
		}
	}

	return resourceAppServerGroupRead(ctx, d, meta)
}

func buildCreateServerGroupBodyParams(d *schema.ResourceData, epsId string) map[string]interface{} {
	return map[string]interface{}{
		"name":                    d.Get("name"),
		"os_type":                 d.Get("os_type"),
		"product_id":              d.Get("flavor_id"),
		"vpc_id":                  d.Get("vpc_id"),
		"subnet_id":               d.Get("subnet_id"),
		"system_disk_type":        d.Get("system_disk_type"),
		"system_disk_size":        d.Get("system_disk_size"),
		"image_id":                d.Get("image_id"),
		"image_type":              d.Get("image_type"),
		"image_product_id":        utils.ValueIgnoreEmpty(d.Get("image_product_id")),
		"is_vdi":                  d.Get("is_vdi"),
		"availability_zone":       utils.ValueIgnoreEmpty(d.Get("availability_zone")),
		"description":             utils.ValueIgnoreEmpty(d.Get("description")),
		"app_type":                utils.ValueIgnoreEmpty(d.Get("app_type")),
		"tags":                    utils.ValueIgnoreEmpty(utils.ExpandResourceTags(d.Get("tags").(map[string]interface{}))),
		"enterprise_project_id":   utils.ValueIgnoreEmpty(epsId),
		"ip_virtual":              buildAppServerGroupIpVirtual(d.Get("ip_virtual").([]interface{})),
		"route_policy":            buildAppServerGroupRoutePolicy(d.Get("route_policy").([]interface{})),
		"ou_name":                 utils.ValueIgnoreEmpty(d.Get("ou_name")),
		"extra_session_type":      utils.ValueIgnoreEmpty(d.Get("extra_session_type")),
		"extra_session_size":      utils.ValueIgnoreEmpty(d.Get("extra_session_size")),
		"primary_server_group_id": utils.ValueIgnoreEmpty(d.Get("primary_server_group_id")),
		"server_group_status":     d.Get("enabled"),
	}
}

func buildAppServerGroupIpVirtual(ipVirtual []interface{}) map[string]interface{} {
	if len(ipVirtual) == 0 {
		return nil
	}

	return map[string]interface{}{
		"enable": utils.PathSearch("enable", ipVirtual[0], false),
	}
}

func buildAppServerGroupRoutePolicy(routePolicy []interface{}) map[string]interface{} {
	if len(routePolicy) == 0 {
		return nil
	}

	return map[string]interface{}{
		"max_session":   utils.PathSearch("max_session", routePolicy[0], 0),
		"cpu_threshold": utils.PathSearch("cpu_threshold", routePolicy[0], 0),
		"mem_threshold": utils.PathSearch("mem_threshold", routePolicy[0], 0),
	}
}

func resourceAppServerGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	serverGroup, err := GetServerGroupById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Workspace APP server group")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", serverGroup, nil)),
		d.Set("os_type", utils.PathSearch("os_type", serverGroup, nil)),
		d.Set("flavor_id", utils.PathSearch("product_id", serverGroup, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", serverGroup, false)),
		d.Set("system_disk_type", utils.PathSearch("system_disk_type", serverGroup, nil)),
		d.Set("system_disk_size", utils.PathSearch("system_disk_size", serverGroup, nil)),
		d.Set("image_id", utils.PathSearch("image_id", serverGroup, nil)),
		d.Set("is_vdi", utils.PathSearch("is_vdi", serverGroup, nil)),
		d.Set("app_type", utils.PathSearch("app_type", serverGroup, nil)),
		d.Set("description", utils.PathSearch("description", serverGroup, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", serverGroup, make([]interface{}, 0)))),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", serverGroup, nil)),
		d.Set("ou_name", utils.PathSearch("ou_name", serverGroup, nil)),
		d.Set("extra_session_type", utils.PathSearch("extra_session_type", serverGroup, nil)),
		d.Set("extra_session_size", utils.PathSearch("extra_session_size", serverGroup, nil)),
		d.Set("primary_server_group_id", utils.PathSearch("primary_server_group_ids|[0]", serverGroup, nil)),
		d.Set("enabled", utils.PathSearch("server_group_status", serverGroup, nil)),
		d.Set("storage_mount_policy", utils.PathSearch("storage_mount_policy", serverGroup, nil)),
		d.Set("project_name", utils.PathSearch("project_name", serverGroup, nil)),
		d.Set("image_min_disk", utils.PathSearch("image_min_disk", serverGroup, nil)),
		d.Set("flavors", flattenAppServerGroupFlavor(utils.PathSearch("flavors", serverGroup,
			make([]interface{}, 0)).([]interface{}))),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// GetServerGroupById is a method used to query server group detail by server group ID.
func GetServerGroupById(client *golangsdk.ServiceClient, serverGroupId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/app-server-groups/{server_group_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{server_group_id}", serverGroupId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, fmt.Errorf("unable to parsing server group from API response: %s", err)
	}

	return respBody, nil
}

func flattenAppServerGroupFlavorLinks(links []interface{}) []map[string]interface{} {
	if len(links) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(links))
	for _, link := range links {
		result = append(result, map[string]interface{}{
			"href": utils.PathSearch("href", link, nil),
			"rel":  utils.PathSearch("rel", link, nil),
		})
	}

	return result
}

func flattenAppServerGroupFlavor(flavors []interface{}) []map[string]interface{} {
	if len(flavors) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(flavors))
	for _, flavor := range flavors {
		links := utils.PathSearch("links", flavor, make([]interface{}, 0)).([]interface{})
		result = append(result, map[string]interface{}{
			"id":    utils.PathSearch("id", flavor, nil),
			"links": flattenAppServerGroupFlavorLinks(links),
		})
	}

	return result
}

func resourceAppServerGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		serverGroupId = d.Id()
	)
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	if d.HasChanges("name", "system_disk_type", "system_disk_size", "image_id", "image_type", "image_product_id",
		"description", "app_type", "route_policy", "ou_name", "enabled", "storage_mount_policy") {
		updateOpt := buildUpdateServerGroupBodyParams(d)
		if err := updateAppServerGroup(client, serverGroupId, updateOpt); err != nil {
			return diag.Errorf("error updating server group (%s): %s", serverGroupId, err)
		}
	}

	if d.HasChanges("tags") {
		oRaw, nRaw := d.GetChange("tags")
		if err := updateServerGroupTags(client, serverGroupId, oRaw, nRaw); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceAppServerGroupRead(ctx, d, meta)
}

func updateAppServerGroup(client *golangsdk.ServiceClient, serverGroupId string, params map[string]interface{}) error {
	httpUrl := "v1/{project_id}/app-server-groups/{server_group_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{server_group_id}", serverGroupId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(params),
	}

	_, err := client.Request("PATCH", updatePath, &updateOpt)
	if err != nil {
		return err
	}

	return nil
}

func buildUpdateServerGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":                 d.Get("name"),
		"system_disk_type":     d.Get("system_disk_type"),
		"system_disk_size":     d.Get("system_disk_size"),
		"image_id":             d.Get("image_id"),
		"image_type":           d.Get("image_type"),
		"image_product_id":     utils.ValueIgnoreEmpty(d.Get("image_product_id")),
		"description":          d.Get("description"),
		"app_type":             utils.ValueIgnoreEmpty(d.Get("app_type")),
		"route_policy":         buildAppServerGroupRoutePolicy(d.Get("route_policy").([]interface{})),
		"ou_name":              d.Get("ou_name"),
		"server_group_status":  d.Get("enabled"),
		"storage_mount_policy": d.Get("storage_mount_policy"),
	}
}

func updateServerGroupTags(client *golangsdk.ServiceClient, serverGroupId string, oRaw, nRaw interface{}) error {
	removeTags := oRaw.(map[string]interface{})
	addTags := nRaw.(map[string]interface{})
	if len(removeTags) > 0 {
		if err := removeServerGroupTags(client, serverGroupId, removeTags); err != nil {
			return fmt.Errorf("error removing tags of server group (%s): %s", serverGroupId, err)
		}
	}

	if len(addTags) > 0 {
		if err := addServerGroupTags(client, serverGroupId, addTags); err != nil {
			return fmt.Errorf("error adding tags of server group (%s): %s", serverGroupId, err)
		}
	}

	return nil
}

func buildUpdateServerGroupTagsBodyParams(serverGroupId string, tags map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"server_group_id": serverGroupId,
				"tags":            utils.ExpandResourceTags(tags),
			},
		},
	}
}

func addServerGroupTags(client *golangsdk.ServiceClient, serverGroupId string, tags map[string]interface{}) error {
	httpUrl := "v1/{project_id}/server-group/tags/batch-create"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateServerGroupTagsBodyParams(serverGroupId, tags),
		OkCodes:          []int{204},
	}
	_, err := client.Request("POST", path, &opt)
	return err
}

func removeServerGroupTags(client *golangsdk.ServiceClient, serverGroupId string, tags map[string]interface{}) error {
	httpUrl := "v1/{project_id}/server-group/tags/batch-delete"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateServerGroupTagsBodyParams(serverGroupId, tags),
	}
	_, err := client.Request("DELETE", path, &opt)
	return err
}

func resourceAppServerGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	httpUrl := "v1/{project_id}/app-server-groups/{server_group_id}"
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	serverGroupId := d.Id()
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{server_group_id}", serverGroupId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// When deleting a non-existent server group, the response status code is 200.
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting server group (%s): %s", serverGroupId, err)
	}

	return nil
}

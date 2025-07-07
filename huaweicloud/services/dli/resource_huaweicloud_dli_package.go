package dli

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dli/v2/spark/resources"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DLI POST /v2.0/{project_id}/resources
// @API DLI GET /v2.0/{project_id}/resources/{resource_name}
// @API DLI PUT /v2.0/{project_id}/resources/owner
// @API DLI DELETE /v2.0/{project_id}/resources/{resource_name}
// @API DLI GET /v3/{project_id}/dli_package_resource/{resource_id}/tags
// @API DLI POST /v3/{project_id}/dli_package_resource/{resource_id}/tags/create
// @API DLI POST /v3/{project_id}/dli_package_resource/{resource_id}/tags/delete
func ResourceDliPackageV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceDliDependentPackageV2Create,
		ReadContext:   ResourceDliDependentPackageV2Read,
		UpdateContext: ResourceDliDependentPackageV2Update,
		DeleteContext: ResourceDliDependentPackageV2Delete,

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"object_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_async": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			"object_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
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

func buildDliDependentPackageCreateOpts(d *schema.ResourceData) resources.CreateGroupAndUploadOpts {
	result := resources.CreateGroupAndUploadOpts{
		Paths:   []string{d.Get("object_path").(string)},
		Kind:    d.Get("type").(string),
		Group:   d.Get("group_name").(string),
		IsAsync: d.Get("is_async").(bool),
		Tags:    utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}
	return result
}

func ResourceDliDependentPackageV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	c, err := cfg.DliV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v2 client: %s", err)
	}

	opt := buildDliDependentPackageCreateOpts(d)
	resp, err := resources.CreateGroupAndUpload(c, opt)
	if err != nil {
		return diag.Errorf("error uploading package to OBS bucket: %s", err)
	}

	// If the object list of resources.Resources is not empty, it means that the object has been uploaded successfully,
	// and the element with an index of zero is the object name.
	groupName, ok := d.GetOk("group_name")
	if ok && len(resp.Resources) < 1 || !ok && len(resp.ResourceNames) < 1 {
		return diag.Errorf("failed to upload package (%s).", d.Get("object_path").(string))
	}
	// object_path is not unique and cannot be used for ID setting, because the object can exist in multiple groups.
	if ok {
		d.SetId(fmt.Sprintf("%s#%s", groupName.(string), resp.Resources[0]))
	} else {
		d.SetId(resp.ResourceNames[0])
	}

	// If the owner of the configuration is not the creator, update it.
	pkg, err := GetDliDependentPackageInfo(c, d.Id())
	if err != nil {
		return diag.Errorf("error getting the package: %s", err)
	}

	if owner, ok := d.GetOk("owner"); ok && owner.(string) != pkg.Owner {
		if err = updateOwner(c, groupName.(string), owner.(string), pkg.ResourceName); err != nil {
			return diag.FromErr(err)
		}
	}

	return ResourceDliDependentPackageV2Read(ctx, d, meta)
}

func getGroupNameAndPackageName(id string) (groupName, packageName string, err error) {
	names := strings.Split(id, "#")
	if len(names) == 1 {
		return "", names[0], nil
	}

	if len(names) == 2 {
		return names[0], names[1], nil
	}

	err = fmt.Errorf("invalid format for resource ID, want '<group_name>#<object_name>' or '<object_name>', but got '%s'", id)
	return
}

func GetDliDependentPackageInfo(c *golangsdk.ServiceClient, id string) (*resources.Resource, error) {
	groupName, packageName, err := getGroupNameAndPackageName(id)
	if err != nil {
		return nil, fmt.Errorf("error parsing resource ID (%s): %s", id, err)
	}

	opt := resources.ResourceLocatedOpts{
		Group: groupName,
	}
	resp, err := resources.Get(c, packageName, opt)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ResourceDliDependentPackageV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	c, err := cfg.DliV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v2 client: %s", err)
	}

	id := d.Id()
	resp, err := GetDliDependentPackageInfo(c, id)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DLI package")
	}

	mErr := multierror.Append(nil,
		d.Set("object_name", resp.ResourceName),
		d.Set("type", resp.ResourceType),
		d.Set("status", resp.Status),
		d.Set("created_at", time.Unix(int64(resp.CreateTime)/1000, 0).Format("2006-01-02 15:04:05")),
		d.Set("updated_at", time.Unix(int64(resp.CreateTime)/1000, 0).Format("2006-01-02 15:04:05")),
		d.Set("owner", resp.Owner),
	)

	v3Client, err := cfg.DliV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v3 client: %s", err)
	}

	err = utils.SetResourceTagsToState(d, v3Client, "dli_package_resource", getASCIIFormationId(id))
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func updateOwner(c *golangsdk.ServiceClient, groupName, owner, resourceName string) error {
	opt := resources.UpdateOpts{
		ResourceName: resourceName,
		GroupName:    groupName,
		NewOwner:     owner,
	}
	resp, err := resources.UpdateOwner(c, opt)
	if err != nil {
		return fmt.Errorf("error updating package owner: %s", err)
	}

	if !resp.IsSuccess {
		return fmt.Errorf("unable to update the package: %s", resp.Message)
	}

	return nil
}

func getASCIIFormationId(id string) string {
	// The url.QueryEscape function is used to convert special characters in the URL to corresponding hexadecimal format.
	return url.QueryEscape(id)
}

func ResourceDliDependentPackageV2Update(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	c, err := cfg.DliV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v2 client: %s", err)
	}

	if err = updateOwner(c, d.Get("group_name").(string), d.Get("owner").(string), d.Get("object_name").(string)); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("tags") {
		v3Client, err := cfg.DliV3Client(region)
		if err != nil {
			return diag.Errorf("error creating DLI v3 client: %s", err)
		}

		resourceId := getASCIIFormationId(d.Id())
		oldTags, newTags := d.GetChange("tags")
		err = updateResourceTags(v3Client, resourceId, "dli_package_resource", oldTags, newTags)
		if err != nil {
			return diag.Errorf("error updating tags of the package (%s): %s", resourceId, err)
		}
	}

	return ResourceDliDependentPackageV2Read(ctx, d, meta)
}

func ResourceDliDependentPackageV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.DliV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DLI v2 client: %s", err)
	}

	groupName, packageName, err := getGroupNameAndPackageName(d.Id())
	if err != nil {
		return diag.Errorf("error parsing resource ID (%s): %s", d.Id(), err)
	}

	opt := resources.ResourceLocatedOpts{
		Group: groupName,
	}
	err = resources.Delete(c, packageName, opt).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting DLI dependent package (%s): %s", packageName, err)
	}

	return nil
}

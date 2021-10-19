package dli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dli/v2/spark/resources"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

const (
	jarFile    = "jar"
	pythonFile = "pyFile"
	userFile   = "file"
)

var uploadPath = map[string]string{
	jarFile:    "jars",
	pythonFile: "pyfiles",
	userFile:   "files",
}

func ResourceDliPackageV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceDliDependentPackageV2Create,
		ReadContext:   ResourceDliDependentPackageV2Read,
		UpdateContext: ResourceDliDependentPackageV2Update,
		DeleteContext: ResourceDliDependentPackageV2Delete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// If you want to add a new package type, please update all relevant codes.
				ValidateFunc: validation.StringInSlice([]string{
					jarFile, pythonFile, userFile,
				}, false),
			},
			"object_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"object_name": {
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
		Paths:   []string{fmt.Sprintf("%s%s", d.Get("object_path").(string), d.Get("object_name").(string))},
		Kind:    d.Get("type").(string),
		Group:   d.Get("group_name").(string),
		IsAsync: d.Get("is_async").(bool),
	}
	return result
}

func buildDliDependentPackageUploadOpts(d *schema.ResourceData) resources.UploadOpts {
	result := resources.UploadOpts{
		Paths: []string{
			fmt.Sprintf("%s/%s", d.Get("object_path").(string), d.Get("object_name").(string)),
		},
		Group: d.Get("group_name").(string),
	}
	return result
}

func ResourceDliDependentPackageV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.DliV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud DLI v2 client: %s", err)
	}

	resList, err := resources.List(c, resources.ListOpts{})
	if err != nil {
		return fmtp.DiagErrorf("Error getting group informations: %s", err)
	}

	// filter data by group name
	filterData, err := utils.FilterSliceWithField(resList.Groups, map[string]interface{}{
		"GroupName": d.Get("group_name").(string),
	})

	if len(filterData) == 0 {
		opt := buildDliDependentPackageCreateOpts(d)
		_, err = resources.CreateGroupAndUpload(c, opt)
		if err != nil {
			return fmtp.DiagErrorf("A error occurred when creating group and upload package: %s", err)
		}
	} else {
		opt := buildDliDependentPackageUploadOpts(d)
		resType := d.Get("type").(string)

		_, err = resources.Upload(c, uploadPath[resType], opt)
		if err != nil {
			return fmtp.DiagErrorf("Error uploading %s package to OBS bucket: %s", resType, err)
		}
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("group_name").(string), d.Get("object_name").(string)))

	pkg, err := GetDliDependentPackageInfo(c, d.Id())
	if err != nil {
		return fmtp.DiagErrorf("An error occurred while getting the package: %s", err)
	}
	if owner, ok := d.GetOk("owner"); ok && owner.(string) != pkg.Owner {
		return ResourceDliDependentPackageV2Update(ctx, d, meta)
	}

	return ResourceDliDependentPackageV2Read(ctx, d, meta)
}

func setDliDependentPackageParameters(d *schema.ResourceData, resp *resources.Resource) error {
	mErr := multierror.Append(nil,
		d.Set("object_name", resp.ResourceName),
		d.Set("type", resp.ResourceType),
		d.Set("status", resp.Status),
		d.Set("created_at", time.Unix(int64(resp.CreateTime)/1000, 0).Format("2006-01-02 15:04:05")),
		d.Set("updated_at", time.Unix(int64(resp.CreateTime)/1000, 0).Format("2006-01-02 15:04:05")),
		d.Set("owner", resp.Owner),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func getGroupNameAndPackageName(id string) (groupName, packageName string, err error) {
	names := strings.Split(id, "/")
	if len(names) < 2 {
		err = fmtp.Errorf("ID is incomplete, missing key information")
		return
	}
	return names[0], names[1], nil
}

func GetDliDependentPackageInfo(c *golangsdk.ServiceClient, id string) (*resources.Resource, error) {
	var resp *resources.Resource
	groupName, packageName, err := getGroupNameAndPackageName(id)
	if err != nil {
		return resp, fmt.Errorf("error parsing resource ID (%s): %s", id, err)
	}

	opt := resources.ResourceLocatedOpts{
		Group: groupName,
	}
	resp, err = resources.Get(c, packageName, opt)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func ResourceDliDependentPackageV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.DliV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud DLI v2 client: %s", err)
	}

	resp, err := GetDliDependentPackageInfo(c, d.Id())
	if err != nil {
		return fmtp.DiagErrorf("An error occurred while getting the package: %s", err)
	}

	err = setDliDependentPackageParameters(d, resp)
	if err != nil {
		return fmtp.DiagErrorf("An error occurred during package resource parameters setting: %s", err)
	}
	return nil
}

func ResourceDliDependentPackageV2Update(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.DliV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud DLI v2 client: %s", err)
	}

	opt := resources.UpdateOpts{
		ResourceName: d.Get("object_name").(string),
		GroupName:    d.Get("group_name").(string),
		NewOwner:     d.Get("owner").(string),
	}
	_, err = resources.UpdateOwner(c, opt)
	if err != nil {
		return fmtp.DiagErrorf("Error updating package owner: %s", err)
	}

	return ResourceDliDependentPackageV2Read(ctx, d, meta)
}

func ResourceDliDependentPackageV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.DliV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud DLI v2 client: %s", err)
	}

	groupName, packageName, err := getGroupNameAndPackageName(d.Id())
	if err != nil {
		return fmtp.DiagErrorf("Error parsing resource ID (%s): %s", d.Id(), err)
	}

	opt := resources.ResourceLocatedOpts{
		Group: groupName,
	}
	err = resources.Delete(c, packageName, opt).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("Error deleting DLI dependent package (%s): %s", packageName, err)
	}

	d.SetId("")
	return nil
}

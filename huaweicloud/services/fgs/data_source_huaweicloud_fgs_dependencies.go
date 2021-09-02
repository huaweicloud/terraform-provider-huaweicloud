package fgs

import (
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/fgs/v2/dependencies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

// DataSourceFunctionGraphDependencies provides some parameters to filter dependent packages on the server.
func DataSourceFunctionGraphDependencies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFunctionGraphDependenciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"public", "private",
				}, false),
			},
			"runtime": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"packages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"link": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"etag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"file_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceFunctionGraphDependenciesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.FgsV2Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}
	// Limit and Marker use default values.
	listOpts := dependencies.ListOpts{
		Name:           d.Get("name").(string),
		Runtime:        d.Get("runtime").(string),
		DependencyType: d.Get("type").(string),
	}

	allPages, err := dependencies.List(client, listOpts).AllPages()
	if err != nil {
		return fmtp.Errorf("Error retrieving dependent packages: %s", err)
	}
	resp, err := dependencies.ExtractDependencies(allPages)
	if len(resp.Dependencies) < 1 {
		return fmtp.Errorf("No dependent package found, please check your parameters.")
	}

	return setFunctionGraphDependencies(d, resp.Dependencies)
}

func setFunctionGraphDependencies(d *schema.ResourceData, pkgs []dependencies.Dependency) error {
	packages := make([]map[string]interface{}, len(pkgs))

	names := schema.NewSet(schema.HashString, nil)
	for i, pkg := range pkgs {
		names.Add(pkg.Name)
		packages[i] = map[string]interface{}{
			"id":        pkg.ID,
			"name":      pkg.Name,
			"owner":     pkg.Owner,
			"link":      pkg.Link,
			"etag":      pkg.Etag,
			"size":      pkg.Size,
			"file_name": pkg.FileName,
			"runtime":   pkg.Runtime,
		}
	}
	d.SetId(hashcode.Strings(utils.ExpandToStringList(names.List())))

	return d.Set("packages", packages)
}

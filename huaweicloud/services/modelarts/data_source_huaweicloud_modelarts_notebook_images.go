package modelarts

import (
	"context"

	"github.com/chnsz/golangsdk/openstack/modelarts/v1/notebook"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func DataSourceNotebookImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotebookImagesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "BUILD_IN",
				ValidateFunc: validation.StringInSlice([]string{"BUILD_IN", "DEDICATED"}, false),
			},
			"cpu_arch": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"x86_64", "aarch64"}, false),
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"images": {
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
						"swr_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_arch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNotebookImagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ModelArtsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating ModelArts v1 client, err=%s", err)
	}

	listOpts := notebook.ListImageOpts{
		Name:        d.Get("name").(string),
		Namespace:   d.Get("organization").(string),
		Type:        d.Get("type").(string),
		WorkspaceId: d.Get("workspace_id").(string),
		Limit:       200,
		Offset:      0,
	}

	page, err := notebook.ListImages(client, listOpts)
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve ModelArts notebook images: %s ", err)
	}

	p, err := page.AllPages()
	if err != nil {
		return fmtp.DiagErrorf("error querying ModelArts notebook images: %s", err)
	}
	images, err := notebook.ExtractImages(p)
	if err != nil {
		return fmtp.DiagErrorf("error querying ModelArts notebook images: %s", err)
	}

	if len(images) == 0 {
		return fmtp.DiagErrorf("No data found. Please change your search criteria and try again.")
	}

	filter := map[string]interface{}{
		"Arch": d.Get("cpu_arch"),
	}

	filterImages, err := utils.FilterSliceWithField(images, filter)
	if err != nil {
		return fmtp.DiagErrorf("filter ModelArts notebook images failed: %s", err)
	}
	logp.Printf("[DEBUG] filter %d ModelArts notebook images from %d through options %v", len(filterImages), len(images), filter)

	var rst []map[string]interface{}
	var ids []string
	for _, v := range filterImages {
		img := v.(notebook.ImageDetail)
		item := map[string]interface{}{
			"id":          img.Id,
			"name":        img.Name,
			"type":        img.Type,
			"swr_path":    img.SwrPath,
			"description": img.Description,
			"cpu_arch":    img.Arch,
		}
		rst = append(rst, item)
		ids = append(ids, img.Id)
	}

	err = d.Set("images", rst)
	if err != nil {
		return fmtp.DiagErrorf("set images err:%s", err)
	}

	d.SetId(hashcode.Strings(ids))
	return nil
}

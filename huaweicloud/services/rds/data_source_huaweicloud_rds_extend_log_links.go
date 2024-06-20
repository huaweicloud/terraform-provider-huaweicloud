package rds

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS POST /v3/{project_id}/instances/{instance_id}/xellog-download
func DataSourceRdsExtendLogLinks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsExtendLogLinksRead,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the RDS instance.`,
			},
			"file_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the file to be downloaded.`,
			},
			"links": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of extend log links.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the name of the file.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the status of the link.`,
						},
						"file_size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the file size in KB.`,
						},
						"file_link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the download link.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creation time.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the last update time.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceRdsExtendLogLinksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	fileName := d.Get("file_name").(string)
	resp, err := extendLogLink(client, instanceID, fileName)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("links", flattenRdsGetExtendLogLinks(resp)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRdsGetExtendLogLinks(resp interface{}) []map[string]interface{} {
	extendLogLinksJson := utils.PathSearch("list", resp, make([]interface{}, 0))
	extendLogLinksArray := extendLogLinksJson.([]interface{})
	if len(extendLogLinksArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(extendLogLinksArray))
	for _, v := range extendLogLinksArray {
		result = append(result, map[string]interface{}{
			"file_name":  utils.PathSearch("file_name", v, nil),
			"status":     utils.PathSearch("status", v, nil),
			"file_size":  utils.PathSearch("file_size", v, nil),
			"file_link":  utils.PathSearch("file_link", v, nil),
			"created_at": utils.PathSearch("create_at", v, nil),
			"updated_at": utils.PathSearch("update_at", v, nil),
		})
	}
	return result
}

package modelarts

import (
	"context"
	"log"
	"regexp"

	"github.com/chnsz/golangsdk/openstack/modelarts/v2/version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceDatasetVerions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatasetVersionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dataset_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"split_ratio": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`([+-]?\d+(\.\d*)?|[+-]?\.\d+),([+-]?\d+(\.\d*)?|[+-]?\.\d+)`),
					"separate the minimum and maximum split ratios with commas"),
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"versions": {
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

						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"split_ratio": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"files": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"storage_path": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"is_current": {
							Type:     schema.TypeBool,
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
				},
			},
		},
	}
}

func dataSourceDatasetVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ModelArtsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ModelArts v2 client, err=%s", err)
	}

	datasetId := d.Get("dataset_id").(string)
	opts := version.ListOpts{
		TrainEvaluateRatio: d.Get("split_ratio").(string),
	}

	page, err := version.List(client, datasetId, opts)
	if err != nil {
		return diag.Errorf("error querying ModelArts dataset versions: %s ", err)
	}

	p, err := page.AllPages()
	if err != nil {
		return diag.Errorf("error querying ModelArts dataset versions: %s", err)
	}

	ds, err := version.ExtractDatasetVersions(p)
	if err != nil {
		return diag.Errorf("error parsing ModelArts dataset versions: %s", err)
	}

	filter := map[string]interface{}{
		"VersionName": d.Get("name"),
	}
	filtResult, err := utils.FilterSliceWithField(ds, filter)
	if err != nil {
		return diag.Errorf("filtering dataset versions failed: %s", err)
	}
	log.Printf("filter %d dataset versions from %d through option %v", len(filtResult), len(ds), filter)

	var rst []map[string]interface{}
	var ids []string
	for _, f := range filtResult {
		v := f.(version.DatasetVersion)
		item := map[string]interface{}{
			"id":           v.VersionId,
			"name":         v.VersionName,
			"description":  v.Description,
			"split_ratio":  v.TrainEvaluateSampleRatio,
			"status":       v.Status,
			"files":        v.TotalSampleCount,
			"storage_path": v.ManifestPath,
			"is_current":   v.IsCurrent,
			"created_at":   utils.FormatTimeStampUTC(int64(v.CreateTime)),
			"updated_at":   utils.FormatTimeStampUTC(int64(v.UpdateTime)),
		}
		rst = append(rst, item)
		ids = append(ids, v.VersionId)
	}

	err = d.Set("versions", rst)
	if err != nil {
		return diag.Errorf("error setting versions: %s", err)
	}

	d.SetId(hashcode.Strings(ids))
	return nil
}

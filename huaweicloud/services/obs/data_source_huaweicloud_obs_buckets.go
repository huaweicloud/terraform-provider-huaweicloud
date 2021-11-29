package obs

import (
	"context"
	"errors"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/obs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func DataSourceObsBuckets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceObsBucketsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"buckets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bucket": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceObsBucketsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.ObjectStorageClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud OBS client: %s", err)
	}
	r, err := client.ListBuckets(&obs.ListBucketsInput{
		BucketType: obs.OBJECT,
	})

	if err != nil {
		return fmtp.DiagErrorf("Error querying OBS buckets: %s", err)
	}

	ids := make([]string, 0, len(r.Buckets))
	buckets := make([]map[string]interface{}, 0, len(r.Buckets))
	for _, v := range r.Buckets {
		metadata, err := queryMetadata(client, v.Name)
		if err != nil && !errors.As(err, &golangsdk.ErrDefault404{}) {
			return fmtp.DiagErrorf("Error querying OBS bucket metadata: %s", err)
		}

		storageClass, enterpriseProjectID, region := "", "", ""
		if metadata != nil {
			storageClass = string(metadata.StorageClass)
			enterpriseProjectID = metadata.Epid
			region = metadata.Location
		}

		if bucket, ok := d.GetOk("bucket"); ok {
			if v.Name != bucket {
				continue
			}
		}
		if epid, ok := d.GetOk("enterprise_project_id"); ok {
			if enterpriseProjectID != epid {
				continue
			}
		}
		bucket := map[string]interface{}{
			"region":                region,
			"bucket":                v.Name,
			"enterprise_project_id": enterpriseProjectID,
			"storage_class":         storageClass,
			"created_at":            v.CreationDate.Format("2006-01-02 15:04:05 MST"),
		}
		buckets = append(buckets, bucket)
		ids = append(ids, v.Name)
	}
	d.SetId(hashcode.Strings(ids))
	err = d.Set("buckets", buckets)
	if err != nil {
		return fmtp.DiagErrorf("error setting OBS attributes: %s", err)
	}
	return nil
}

func queryMetadata(client *obs.ObsClient, name string) (*obs.GetBucketMetadataOutput, error) {
	input := obs.GetBucketMetadataInput{
		Bucket: name,
	}
	metadata, err := client.GetBucketMetadata(&input)
	if obsError, ok := err.(obs.ObsError); ok && obsError.StatusCode == 404 {
		err = golangsdk.ErrDefault404{}
	}
	return metadata, err
}

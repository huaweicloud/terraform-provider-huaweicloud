package obs

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/obs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceObsBucketObject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceObsBucketObjectRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},

			"key": {
				Type:     schema.TypeString,
				Required: true,
			},

			"storage_class": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

// Attribute parameters are not returned in one interface.
// Two interfaces need to be called to get all parameters.
func dataSourceObsBucketObjectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	obsClient, err := conf.ObjectStorageClient(region)
	if err != nil {
		return diag.Errorf("Error creating OBS client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	objects, err := obsClient.ListObjects(&obs.ListObjectsInput{
		Bucket: bucket,
		ListObjsInput: obs.ListObjsInput{
			Prefix: key,
		},
	})
	if err != nil {
		return diag.FromErr(getObsError("Error listing objects of OBS bucket", bucket, err))
	}

	var exist bool
	var objectContent obs.Content
	for _, content := range objects.Contents {
		if key == content.Key {
			exist = true
			objectContent = content
			break
		}
	}
	if !exist {
		return diag.Errorf("object %s not found in bucket %s", key, bucket)
	}

	log.Printf("[DEBUG] Data Source Reading OBS Bucket Object %s: %#v", key, objectContent)

	object, err := obsClient.GetObject(&obs.GetObjectInput{
		GetObjectMetadataInput: obs.GetObjectMetadataInput{
			Bucket: bucket,
			Key:    key,
		},
	})
	if err != nil {
		return diag.FromErr(getObsError("Error get object info of OBS bucket", bucket, err))
	}

	log.Printf("[DEBUG] Data Source Reading OBS Bucket Object : %#v", object)

	d.SetId(key)

	class := string(objectContent.StorageClass)
	if class == "" {
		class = "STANDARD"
	} else {
		class = normalizeStorageClass(class)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("size", objectContent.Size),
		d.Set("etag", strings.Trim(objectContent.ETag, `"`)),
		d.Set("version_id", object.VersionId),
		d.Set("content_type", object.ContentType),
		d.Set("storage_class", class),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting bucket object fields: %s", err)
	}

	return nil
}

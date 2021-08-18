package huaweicloud

import (
	"strings"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/obs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceObsBucketObject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceObsBucketObjectRead,

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
func dataSourceObsBucketObjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	obsClient, err := config.ObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud OBS client: %s", err)
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
		return getObsError("Error listing objects of OBS bucket", bucket, err)
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
		return fmtp.Errorf("object %s not found in bucket %s", key, bucket)
	}

	logp.Printf("[DEBUG] Data Source Reading OBS Bucket Object %s: %#v", key, objectContent)

	object, err := obsClient.GetObject(&obs.GetObjectInput{
		GetObjectMetadataInput: obs.GetObjectMetadataInput{
			Bucket: bucket,
			Key:    key,
		},
	})
	if err != nil {
		return getObsError("Error get object info of OBS bucket", bucket, err)
	}

	logp.Printf("[DEBUG] Data Source Reading OBS Bucket Object : %#v", object)

	d.SetId(key)
	d.Set("size", objectContent.Size)
	d.Set("etag", strings.Trim(objectContent.ETag, `"`))
	d.Set("version_id", object.VersionId)
	d.Set("content_type", object.ContentType)

	class := string(objectContent.StorageClass)
	if class == "" {
		d.Set("storage_class", "STANDARD")
	} else {
		d.Set("storage_class", normalizeStorageClass(class))
	}

	return nil
}

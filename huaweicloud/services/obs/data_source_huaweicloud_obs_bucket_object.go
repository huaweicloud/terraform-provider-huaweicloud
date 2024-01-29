package obs

import (
	"bytes"
	"context"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/obs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API OBS HEAD /{ObjectName}
// @API OBS GET /{ObjectName}
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
			"body": {
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
	objectMeta, err := obsClient.GetObjectMetadata(&obs.GetObjectMetadataInput{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok && obsError.StatusCode == 404 {
			return diag.Errorf("object %s not found in bucket %s", key, bucket)
		}
		return diag.Errorf("error fetching object %s in bucket %s: %s", key, bucket, err)
	}

	log.Printf("[DEBUG] Reading OBS Bucket Object %s metadata: %#v", key, objectMeta)
	d.SetId(key)

	class := string(objectMeta.StorageClass)
	if class == "" {
		class = "STANDARD"
	} else {
		class = normalizeStorageClass(class)
	}

	size := objectMeta.ContentLength
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("storage_class", class),
		d.Set("size", size),
		d.Set("etag", strings.Trim(objectMeta.ETag, `"`)),
		d.Set("version_id", objectMeta.VersionId),
		d.Set("content_type", objectMeta.ContentType),
	)

	// body is available only for objects which have a human-readable Content-Type
	// (text/* and application/json) and smaller than 64KB
	if isContentTypeAllowed(objectMeta.ContentType) && size < 65536 {
		object, err := obsClient.GetObject(&obs.GetObjectInput{
			GetObjectMetadataInput: obs.GetObjectMetadataInput{
				Bucket: bucket,
				Key:    key,
			},
		})
		if err != nil {
			return diag.FromErr(getObsError("Error get object info of OBS bucket", bucket, err))
		}

		defer object.Body.Close()
		log.Printf("[DEBUG] Reading OBS Bucket Object : %#v", object)

		buf := new(bytes.Buffer)
		bytesRead, err := buf.ReadFrom(object.Body)
		if err != nil {
			return diag.Errorf("Failed reading content of OBS object (%s): %s", key, err)
		}

		log.Printf("[INFO] saving %d bytes from OBS object %s", bytesRead, key)
		mErr = multierror.Append(mErr, d.Set("body", buf.String()))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting bucket object fields: %s", err)
	}

	return nil
}

func isContentTypeAllowed(contentType string) bool {
	allowedContentTypes := []*regexp.Regexp{
		regexp.MustCompile("^text/.+"),
		regexp.MustCompile("^application/json$"),
	}

	for _, r := range allowedContentTypes {
		if r.MatchString(contentType) {
			return true
		}
	}

	return false
}

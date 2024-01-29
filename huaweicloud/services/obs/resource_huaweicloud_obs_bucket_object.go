package obs

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/obs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API OBS HEAD /
// @API OBS HEAD /{ObjectName}
// @API OBS PUT /{ObjectName}
// @API OBS DELETE /{ObjectName}
func ResourceObsBucketObject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObsBucketObjectPut,
		ReadContext:   resourceObsBucketObjectRead,
		UpdateContext: resourceObsBucketObjectPut,
		DeleteContext: resourceObsBucketObjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceObsBucketObjectImport,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"source", "content"},
			},

			"storage_class": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"acl": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"encryption": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"etag": {
				Type: schema.TypeString,
				// This will conflict with server-side-encryption and multi-part upload
				// if/when it's actually implemented. The Etag then won't match raw-file MD5.
				Optional: true,
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

func resourceObsBucketObjectPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var resp *obs.PutObjectOutput
	var err error

	conf := meta.(*config.Config)
	obsClient, err := conf.ObjectStorageClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating OBS client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	_, err = obsClient.HeadBucket(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok && obsError.StatusCode == 404 {
			return diag.Errorf("OBS bucket(%s) not found", bucket)
		}
		return diag.Errorf("error reading OBS bucket %s: %s", bucket, err)
	}

	source := d.Get("source").(string)
	content := d.Get("content").(string)
	if source != "" {
		// check source file whether exist
		_, err = os.Stat(source)
		if err != nil {
			if os.IsNotExist(err) {
				return diag.Errorf("source file %s is not exist", source)
			}
			return diag.FromErr(err)
		}

		// put source file
		resp, err = putFileToObject(obsClient, d)
	}

	if content != "" {
		// put content
		resp, err = putContentToObject(obsClient, d)
	}

	if err != nil {
		return diag.FromErr(getObsError("Error putting object to OBS bucket", bucket, err))
	}
	if resp == nil {
		return diag.Errorf("putting object to OBS bucket %s without null response", bucket)
	}

	log.Printf("[DEBUG] Response of putting %s to OBS Bucket %s: %#v", key, bucket, resp)
	mErr := &multierror.Error{}
	if resp.VersionId != "null" {
		mErr = multierror.Append(mErr, d.Set("version_id", resp.VersionId))
	} else {
		mErr = multierror.Append(mErr, d.Set("version_id", ""))
	}
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving versionId of OBS bucket %s: %s", bucket, mErr)
	}
	d.SetId(key)

	return resourceObsBucketObjectRead(ctx, d, meta)
}

func putContentToObject(obsClient *obs.ObsClient, d *schema.ResourceData) (*obs.PutObjectOutput, error) {
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	content := d.Get("content").(string)

	putInput := &obs.PutObjectInput{}
	putInput.Bucket = bucket
	putInput.Key = key

	if v, ok := d.GetOk("acl"); ok {
		putInput.ACL = obs.AclType(v.(string))
	}
	if v, ok := d.GetOk("storage_class"); ok {
		putInput.StorageClass = obs.StorageClassType(v.(string))
	}
	if v, ok := d.GetOk("content_type"); ok {
		putInput.ContentType = v.(string)
	}

	var sseKmsHeader = obs.SseKmsHeader{}
	if d.Get("encryption").(bool) {
		sseKmsHeader.Encryption = obs.DEFAULT_SSE_KMS_ENCRYPTION
		sseKmsHeader.Key = d.Get("kms_key_id").(string)
		putInput.SseHeader = sseKmsHeader
	}

	log.Printf("[DEBUG] putting %s to OBS Bucket %s, opts: %#v", key, bucket, putInput)
	// do not log content
	body := bytes.NewReader([]byte(content))
	putInput.Body = body

	return obsClient.PutObject(putInput)
}

func putFileToObject(obsClient *obs.ObsClient, d *schema.ResourceData) (*obs.PutObjectOutput, error) {
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	putInput := &obs.PutFileInput{}
	putInput.Bucket = bucket
	putInput.Key = key
	putInput.SourceFile = d.Get("source").(string)

	if v, ok := d.GetOk("acl"); ok {
		putInput.ACL = obs.AclType(v.(string))
	}
	if v, ok := d.GetOk("storage_class"); ok {
		putInput.StorageClass = obs.StorageClassType(v.(string))
	}
	if v, ok := d.GetOk("content_type"); ok {
		putInput.ContentType = v.(string)
	}

	var sseKmsHeader = obs.SseKmsHeader{}
	if d.Get("encryption").(bool) {
		sseKmsHeader.Encryption = obs.DEFAULT_SSE_KMS_ENCRYPTION
		sseKmsHeader.Key = d.Get("kms_key_id").(string)
		putInput.SseHeader = sseKmsHeader
	}

	log.Printf("[DEBUG] putting %s to OBS Bucket %s, opts: %#v", key, bucket, putInput)
	return obsClient.PutFile(putInput)
}

func resourceObsBucketObjectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Resource not found",
					Detail:   fmt.Sprintf("object %s not found in bucket %s", key, bucket),
				},
			}
		}
		return diag.Errorf("error fetching object %s in bucket %s: %s", key, bucket, err)
	}

	log.Printf("[DEBUG] Reading OBS Bucket Object %s metadata: %#v", key, objectMeta)
	class := string(objectMeta.StorageClass)
	if class == "" {
		class = "STANDARD"
	} else {
		class = normalizeStorageClass(class)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("storage_class", class),
		d.Set("content_type", objectMeta.ContentType),
		d.Set("version_id", objectMeta.VersionId),
		d.Set("size", objectMeta.ContentLength),
		d.Set("etag", strings.Trim(objectMeta.ETag, `"`)),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting bucket object fields: %s", err)
	}

	return nil
}

func resourceObsBucketObjectDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	obsClient, err := conf.ObjectStorageClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating OBS client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	input := &obs.DeleteObjectInput{
		Bucket: bucket,
		Key:    key,
	}

	log.Printf("[DEBUG] Object %s will be deleted with all versions", key)
	_, err = obsClient.DeleteObject(input)
	if err != nil {
		return diag.FromErr(getObsError("Error deleting object of OBS bucket", bucket, err))
	}

	return nil
}

func resourceObsBucketObjectImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format specified for OBS bucket object. Format must be <bucket>/<key>")
		return nil, err
	}

	bucket := parts[0]
	key := parts[1]

	mErr := multierror.Append(nil,
		d.Set("bucket", bucket),
		d.Set("key", key),
	)
	if mErr.ErrorOrNil() != nil {
		return nil, fmt.Errorf("error setting attributes of OBS bucket %s: %s", bucket, mErr)
	}
	d.SetId(key)

	return []*schema.ResourceData{d}, nil
}

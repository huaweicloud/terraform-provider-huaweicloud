package huaweicloud

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk/openstack/obs"
)

func resourceObsBucketObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceObsBucketObjectPut,
		Read:   resourceObsBucketObjectRead,
		Update: resourceObsBucketObjectPut,
		Delete: resourceObsBucketObjectDelete,

		Schema: map[string]*schema.Schema{
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
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
			},

			"storage_class": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"STANDARD", "WARM", "COLD",
				}, true),
			},

			"acl": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"private", "public-read", "public-read-write",
				}, true),
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

func resourceObsBucketObjectPut(d *schema.ResourceData, meta interface{}) error {
	var resp *obs.PutObjectOutput
	var err error

	config := meta.(*Config)
	obsClient, err := config.newObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	source := d.Get("source").(string)
	content := d.Get("content").(string)
	if source == "" && content == "" {
		return fmt.Errorf("Must specify \"source\" or \"content\" field")
	}

	if source != "" {
		// check source file whether exist
		_, err := os.Stat(source)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("source file %s is not exist", source)
			}
			return err
		}

		// put source file
		resp, err = putFileToObject(obsClient, d)
	}

	if content != "" {
		// put content
		resp, err = putContentToObject(obsClient, d)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	if err != nil {
		return getObsError("Error putting object to OBS bucket", bucket, err)
	}

	log.Printf("[DEBUG] Response of putting %s to OBS Bucket %s: %#v", key, bucket, resp)
	if resp.VersionId != "null" {
		d.Set("version_id", resp.VersionId)
	} else {
		d.Set("version_id", "")
	}
	d.SetId(key)

	return resourceObsBucketObjectRead(d, meta)
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

func resourceObsBucketObjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	obsClient, err := config.newObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	input := &obs.ListObjectsInput{}
	input.Bucket = bucket
	input.Prefix = key

	resp, err := obsClient.ListObjects(input)
	if err != nil {
		return getObsError("Error listing objects of OBS bucket", bucket, err)
	}

	var exist bool
	var object obs.Content
	for _, content := range resp.Contents {
		if key == content.Key {
			exist = true
			object = content
			break
		}
	}
	if !exist {
		d.SetId("")
		return fmt.Errorf("object %s not found in bucket %s", key, bucket)
	}
	log.Printf("[DEBUG] Reading OBS Bucket Object %s: %#v", key, object)

	class := string(object.StorageClass)
	if class == "" {
		d.Set("storage_class", "STANDARD")
	} else {
		d.Set("storage_class", normalizeStorageClass(class))
	}
	d.Set("size", object.Size)
	return nil
}

func resourceObsBucketObjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	obsClient, err := config.newObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud OBS client: %s", err)
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
		return getObsError("Error deleting object of OBS bucket", bucket, err)
	}

	return nil
}

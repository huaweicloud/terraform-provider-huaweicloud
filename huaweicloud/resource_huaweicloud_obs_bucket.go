package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk/openstack/obs"
)

func resourceObsBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceObsBucketCreate,
		Read:   resourceObsBucketRead,
		Update: resourceObsBucketUpdate,
		Delete: resourceObsBucketDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"storage_class": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "STANDARD",
				ValidateFunc: validation.StringInSlice([]string{
					"STANDARD", "WARM", "COLD",
				}, true),
			},

			"acl": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "private",
				ValidateFunc: validation.StringInSlice([]string{
					"private", "public-read", "public-read-write",
				}, true),
			},

			"versioning": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"logging": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"target_bucket": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"target_prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"bucket_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceObsBucketCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	obsClient, err := config.newObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	acl := d.Get("acl").(string)
	class := d.Get("storage_class").(string)
	opts := &obs.CreateBucketInput{
		Bucket:       bucket,
		ACL:          obs.AclType(acl),
		StorageClass: obs.StorageClassType(class),
	}
	opts.Location = d.Get("region").(string)
	log.Printf("[DEBUG] OBS bucket create opts: %#v", opts)

	_, err = obsClient.CreateBucket(opts)
	if err != nil {
		return getObsError("Error creating bucket", bucket, err)
	}

	// Assign the bucket name as the resource ID
	d.SetId(bucket)
	return resourceObsBucketUpdate(d, meta)
}

func resourceObsBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	obsClient, err := config.newObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	log.Printf("[DEBUG] Update OBS bucket %s", d.Id())
	if d.HasChange("acl") && !d.IsNewResource() {
		if err := resourceObsBucketAclUpdate(obsClient, d); err != nil {
			return err
		}
	}

	if d.HasChange("storage_class") && !d.IsNewResource() {
		if err := resourceObsBucketClassUpdate(obsClient, d); err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		if err := resourceObsBucketTagsUpdate(obsClient, d); err != nil {
			return err
		}
	}

	if d.HasChange("versioning") {
		if err := resourceObsBucketVersioningUpdate(obsClient, d); err != nil {
			return err
		}
	}

	if d.HasChange("logging") {
		if err := resourceObsBucketLoggingUpdate(obsClient, d); err != nil {
			return err
		}
	}

	return resourceObsBucketRead(d, meta)
}

func resourceObsBucketRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	obsClient, err := config.newObjectStorageClient(region)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	log.Printf("[DEBUG] Read OBS bucket: %s", d.Id())
	_, err = obsClient.HeadBucket(d.Id())
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok && obsError.StatusCode == 404 {
			log.Printf("[WARN] OBS bucket(%s) not found", d.Id())
			d.SetId("")
			return nil
		} else {
			return fmt.Errorf("error reading OBS bucket %s: %s", d.Id(), err)
		}
	}

	d.Set("region", region)
	d.Set("bucket_domain_name", bucketDomainName(d.Get("bucket").(string), region))

	// Read storage class
	if err := setObsBucketStorageClass(obsClient, d); err != nil {
		return err
	}

	// Read the versioning
	if err := setObsBucketVersioning(obsClient, d); err != nil {
		return err
	}
	// Read the logging configuration
	if err := setObsBucketLogging(obsClient, d); err != nil {
		return err
	}

	// Read the tags
	if err := setObsBucketTags(obsClient, d); err != nil {
		return err
	}

	return nil
}

func resourceObsBucketDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	obsClient, err := config.newObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	bucket := d.Id()
	log.Printf("[DEBUG] Obs Delete Bucket: %s", bucket)
	_, err = obsClient.DeleteBucket(bucket)
	if err != nil {
		obsError, ok := err.(obs.ObsError)
		if ok && obsError.Code == "BucketNotEmpty" {
			// todo
			log.Printf("[DEBUG] OBS bucket: %s is not empty", bucket)
		}
		return fmt.Errorf("Error deleting OBS bucket: %s %s", err, bucket)
	}
	return nil
}

func resourceObsBucketTagsUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	tagmap := d.Get("tags").(map[string]interface{})
	tagList := []obs.Tag{}
	for k, v := range tagmap {
		tag := obs.Tag{
			Key:   k,
			Value: v.(string),
		}
		tagList = append(tagList, tag)
	}

	req := &obs.SetBucketTaggingInput{}
	req.Bucket = bucket
	req.Tags = tagList
	log.Printf("[DEBUG] set tags of OBS bucket %s: %#v", bucket, req)

	_, err := obsClient.SetBucketTagging(req)
	if err != nil {
		return getObsError("Error updating tags of OBS bucket", bucket, err)
	}
	return nil
}

func resourceObsBucketAclUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	acl := d.Get("acl").(string)

	i := &obs.SetBucketAclInput{
		Bucket: bucket,
		ACL:    obs.AclType(acl),
	}
	log.Printf("[DEBUG] set ACL of OBS bucket %s: %#v", bucket, i)

	_, err := obsClient.SetBucketAcl(i)
	if err != nil {
		return getObsError("Error updating acl of OBS bucket", bucket, err)
	}

	// acl policy can not be retrieved by obsClient.GetBucketAcl method
	d.Set("acl", acl)
	return nil
}

func resourceObsBucketClassUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	class := d.Get("storage_class").(string)

	input := &obs.SetBucketStoragePolicyInput{}
	input.Bucket = bucket
	input.StorageClass = obs.StorageClassType(class)
	log.Printf("[DEBUG] set storage class of OBS bucket %s: %#v", bucket, input)

	_, err := obsClient.SetBucketStoragePolicy(input)
	if err != nil {
		return getObsError("Error updating storage class of OBS bucket", bucket, err)
	}

	return nil
}

func resourceObsBucketVersioningUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	version := d.Get("versioning").(bool)

	input := &obs.SetBucketVersioningInput{}
	input.Bucket = bucket
	if version {
		input.Status = obs.VersioningStatusEnabled
	} else {
		input.Status = obs.VersioningStatusSuspended
	}
	log.Printf("[DEBUG] set versioning of OBS bucket %s: %#v", bucket, input)

	_, err := obsClient.SetBucketVersioning(input)
	if err != nil {
		return getObsError("Error setting versining status of OBS bucket", bucket, err)
	}

	return nil
}

func resourceObsBucketLoggingUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	rawLogging := d.Get("logging").(*schema.Set).List()
	loggingStatus := &obs.SetBucketLoggingConfigurationInput{}
	loggingStatus.Bucket = bucket

	if len(rawLogging) > 0 {
		c := rawLogging[0].(map[string]interface{})
		enable := false

		if val, ok := c["enabled"]; ok {
			enable = val.(bool)
		}
		if enable {
			targetBucket := bucket
			if val := c["target_bucket"].(string); val != "" {
				targetBucket = val
			}
			loggingStatus.TargetBucket = targetBucket

			if val := c["target_prefix"].(string); val != "" {
				loggingStatus.TargetPrefix = val
			} else {
				loggingStatus.TargetPrefix = fmt.Sprintf("%s-log/", targetBucket)
			}
		}
	}
	log.Printf("[DEBUG] set logging of OBS bucket %s: %#v", bucket, loggingStatus)

	_, err := obsClient.SetBucketLoggingConfiguration(loggingStatus)
	if err != nil {
		return getObsError("Error setting logging configuration of OBS bucket", bucket, err)
	}

	return nil
}

func setObsBucketStorageClass(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketStoragePolicy(bucket)
	if err != nil {
		return getObsError("Error getting storage class of OBS bucket", bucket, err)
	}

	class := string(output.StorageClass)
	// change format of storage class
	if class == "STANDARD_IA" {
		class = "WARM"
	} else if class == "GLACIER" {
		class = "COLD"
	}
	d.Set("storage_class", class)

	return nil
}

func setObsBucketVersioning(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketVersioning(bucket)
	if err != nil {
		return getObsError("Error getting versioning status of OBS bucket", bucket, err)
	}

	if output.Status == obs.VersioningStatusEnabled {
		d.Set("versioning", true)
	} else {
		d.Set("versioning", false)
	}

	return nil
}

func setObsBucketLogging(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketLoggingConfiguration(bucket)
	if err != nil {
		return getObsError("Error getting logging configuration of OBS bucket", bucket, err)
	}

	lcList := make([]map[string]interface{}, 0, 1)
	logging := make(map[string]interface{})
	if output.TargetBucket != "" {
		logging["enabled"] = true
		logging["target_bucket"] = output.TargetBucket
	} else {
		logging["enabled"] = false
	}
	if output.TargetPrefix != "" {
		logging["target_prefix"] = output.TargetPrefix
	}
	lcList = append(lcList, logging)
	if err := d.Set("logging", lcList); err != nil {
		return fmt.Errorf("Error saving logging configuration of OBS bucket %s: %s", bucket, err)
	}

	return nil
}

func setObsBucketTags(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketTagging(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.Code == "NoSuchTagSet" {
				d.Set("tags", nil)
				return nil
			} else {
				return fmt.Errorf("Error getting tags of OBS bucket %s: %s,\n Reason: %s",
					bucket, obsError.Code, obsError.Message)
			}
		} else {
			return err
		}
	}

	tagmap := make(map[string]string)
	for _, tag := range output.Tags {
		tagmap[tag.Key] = tag.Value
	}
	if err := d.Set("tags", tagmap); err != nil {
		return fmt.Errorf("Error saving tags of OBS bucket %s: %s", bucket, err)
	}
	return nil
}

func getObsError(action string, bucket string, err error) error {
	if obsError, ok := err.(obs.ObsError); ok {
		return fmt.Errorf("%s %s: %s,\n Reason: %s", action, bucket, obsError.Code, obsError.Message)
	} else {
		return err
	}
}

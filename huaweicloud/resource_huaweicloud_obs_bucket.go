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

			"lifecycle_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"expiration": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"transition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"storage_class": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"WARM", "COLD",
										}, true),
									},
								},
							},
						},
						"noncurrent_version_expiration": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"noncurrent_version_transition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"storage_class": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"WARM", "COLD",
										}, true),
									},
								},
							},
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

	if d.HasChange("lifecycle_rule") {
		if err := resourceObsBucketLifecycleUpdate(obsClient, d); err != nil {
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

	// Read the Lifecycle configuration
	if err := setObsBucketLifecycleConfiguration(obsClient, d); err != nil {
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

func resourceObsBucketLifecycleUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	lifecycleRules := d.Get("lifecycle_rule").([]interface{})

	if len(lifecycleRules) == 0 {
		log.Printf("[DEBUG] remove all lifecycle rules of bucket %s", bucket)
		_, err := obsClient.DeleteBucketLifecycleConfiguration(bucket)
		if err != nil {
			return getObsError("Error deleting lifecycle rules of OBS bucket", bucket, err)
		}
		return nil
	}

	rules := make([]obs.LifecycleRule, len(lifecycleRules))
	for i, lifecycleRule := range lifecycleRules {
		r := lifecycleRule.(map[string]interface{})

		// rule ID
		rules[i].ID = r["name"].(string)

		// Enabled
		if val, ok := r["enabled"].(bool); ok && val {
			rules[i].Status = obs.RuleStatusEnabled
		} else {
			rules[i].Status = obs.RuleStatusDisabled
		}

		// Prefix
		rules[i].Prefix = r["prefix"].(string)

		// Expiration
		expiration := d.Get(fmt.Sprintf("lifecycle_rule.%d.expiration", i)).(*schema.Set).List()
		if len(expiration) > 0 {
			raw := expiration[0].(map[string]interface{})
			exp := &rules[i].Expiration

			if val, ok := raw["days"].(int); ok && val > 0 {
				exp.Days = val
			}
		}

		// Transition
		transitions := d.Get(fmt.Sprintf("lifecycle_rule.%d.transition", i)).([]interface{})
		list := make([]obs.Transition, len(transitions))
		for j, tran := range transitions {
			raw := tran.(map[string]interface{})

			if val, ok := raw["days"].(int); ok && val > 0 {
				list[j].Days = val
			}
			if val, ok := raw["storage_class"].(string); ok {
				list[j].StorageClass = obs.StorageClassType(val)
			}
		}
		rules[i].Transitions = list

		// NoncurrentVersionExpiration
		nc_expiration := d.Get(fmt.Sprintf("lifecycle_rule.%d.noncurrent_version_expiration", i)).(*schema.Set).List()
		if len(nc_expiration) > 0 {
			raw := nc_expiration[0].(map[string]interface{})
			nc_exp := &rules[i].NoncurrentVersionExpiration

			if val, ok := raw["days"].(int); ok && val > 0 {
				nc_exp.NoncurrentDays = val
			}
		}

		// NoncurrentVersionTransition
		nc_transitions := d.Get(fmt.Sprintf("lifecycle_rule.%d.noncurrent_version_transition", i)).([]interface{})
		nc_list := make([]obs.NoncurrentVersionTransition, len(nc_transitions))
		for j, nc_tran := range nc_transitions {
			raw := nc_tran.(map[string]interface{})

			if val, ok := raw["days"].(int); ok && val > 0 {
				nc_list[j].NoncurrentDays = val
			}
			if val, ok := raw["storage_class"].(string); ok {
				nc_list[j].StorageClass = obs.StorageClassType(val)
			}
		}
		rules[i].NoncurrentVersionTransitions = nc_list
	}

	opts := &obs.SetBucketLifecycleConfigurationInput{}
	opts.Bucket = bucket
	opts.LifecycleRules = rules
	log.Printf("[DEBUG] set lifecycle configurations of OBS bucket %s: %#v", bucket, opts)

	_, err := obsClient.SetBucketLifecycleConfiguration(opts)
	if err != nil {
		return getObsError("Error setting lifecycle rules of OBS bucket", bucket, err)
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
	d.Set("storage_class", normalizeStorageClass(class))

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

func setObsBucketLifecycleConfiguration(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketLifecycleConfiguration(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.Code == "NoSuchLifecycleConfiguration" {
				d.Set("lifecycle_rule", nil)
				return nil
			} else {
				return fmt.Errorf("Error getting lifecycle configuration of OBS bucket %s: %s,\n Reason: %s",
					bucket, obsError.Code, obsError.Message)
			}
		} else {
			return err
		}
	}

	rawRules := output.LifecycleRules
	log.Printf("[DEBUG] getting lifecycle configuration of OBS bucket: %s, lifecycle: %#v", bucket, rawRules)

	rules := make([]map[string]interface{}, 0, len(rawRules))
	for _, lifecycleRule := range rawRules {
		rule := make(map[string]interface{})
		rule["name"] = lifecycleRule.ID

		// Enabled
		if lifecycleRule.Status == obs.RuleStatusEnabled {
			rule["enabled"] = true
		} else {
			rule["enabled"] = false
		}

		if lifecycleRule.Prefix != "" {
			rule["prefix"] = lifecycleRule.Prefix
		}

		// expiration
		if days := lifecycleRule.Expiration.Days; days > 0 {
			e := make(map[string]interface{})
			e["days"] = days
			rule["expiration"] = schema.NewSet(expirationHash, []interface{}{e})
		}
		// transition
		if len(lifecycleRule.Transitions) > 0 {
			transitions := make([]interface{}, 0, len(lifecycleRule.Transitions))
			for _, v := range lifecycleRule.Transitions {
				t := make(map[string]interface{})
				t["days"] = v.Days
				t["storage_class"] = normalizeStorageClass(string(v.StorageClass))
				transitions = append(transitions, t)
			}
			rule["transition"] = transitions
		}

		// noncurrent_version_expiration
		if days := lifecycleRule.NoncurrentVersionExpiration.NoncurrentDays; days > 0 {
			e := make(map[string]interface{})
			e["days"] = days
			rule["noncurrent_version_expiration"] = schema.NewSet(expirationHash, []interface{}{e})
		}

		// noncurrent_version_transition
		if len(lifecycleRule.NoncurrentVersionTransitions) > 0 {
			transitions := make([]interface{}, 0, len(lifecycleRule.NoncurrentVersionTransitions))
			for _, v := range lifecycleRule.NoncurrentVersionTransitions {
				t := make(map[string]interface{})
				t["days"] = v.NoncurrentDays
				t["storage_class"] = normalizeStorageClass(string(v.StorageClass))
				transitions = append(transitions, t)
			}
			rule["noncurrent_version_transition"] = transitions
		}

		rules = append(rules, rule)
	}

	log.Printf("[DEBUG] saving lifecycle configuration of OBS bucket: %s, lifecycle: %#v", bucket, rules)
	if err := d.Set("lifecycle_rule", rules); err != nil {
		return fmt.Errorf("Error saving lifecycle configuration of OBS bucket %s: %s", bucket, err)
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

// normalize format of storage class
func normalizeStorageClass(class string) string {
	var ret string = class

	if class == "STANDARD_IA" {
		ret = "WARM"
	} else if class == "GLACIER" {
		ret = "COLD"
	}
	return ret
}

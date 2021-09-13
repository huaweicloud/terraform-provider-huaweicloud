package huaweicloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/chnsz/golangsdk/openstack/obs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceObsBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceObsBucketCreate,
		Read:   resourceObsBucketRead,
		Update: resourceObsBucketUpdate,
		Delete: resourceObsBucketDelete,
		Importer: &schema.ResourceImporter{
			State: resourceObsBucketImport,
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
					"private", "public-read", "public-read-write", "log-delivery-write",
				}, true),
			},

			"policy": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     utils.ValidateJsonString,
				DiffSuppressFunc: utils.SuppressEquivalentAwsPolicyDiffs,
			},

			"policy_format": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "obs",
			},

			"versioning": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"logging": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_bucket": {
							Type:     schema.TypeString,
							Required: true,
						},
						"target_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "logs/",
						},
					},
				},
			},

			"quota": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
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

			"website": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index_document": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"error_document": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"redirect_all_requests_to": {
							Type: schema.TypeString,
							ConflictsWith: []string{
								"website.0.index_document",
								"website.0.error_document",
								"website.0.routing_rules",
							},
							Optional: true,
						},

						"routing_rules": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: utils.ValidateJsonString,
							StateFunc: func(v interface{}) string {
								json, _ := utils.NormalizeJsonString(v)
								return json
							},
						},
					},
				},
			},

			"cors_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_origins": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_methods": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_headers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"expose_headers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"max_age_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  100,
						},
					},
				},
			},

			"tags": tagsSchema(),
			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"multi_az": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"parallel_fs": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"bucket_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceObsBucketCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	obsClient, err := config.ObjectStorageClient(region)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	acl := d.Get("acl").(string)
	class := d.Get("storage_class").(string)
	opts := &obs.CreateBucketInput{
		Bucket:            bucket,
		ACL:               obs.AclType(acl),
		StorageClass:      obs.StorageClassType(class),
		IsFSFileInterface: d.Get("parallel_fs").(bool),
		Epid:              GetEnterpriseProjectID(d, config),
	}
	opts.Location = region
	if _, ok := d.GetOk("multi_az"); ok {
		opts.AvailableZone = "3az"
	}

	logp.Printf("[DEBUG] OBS bucket create opts: %#v", opts)
	_, err = obsClient.CreateBucket(opts)
	if err != nil {
		return getObsError("Error creating bucket", bucket, err)
	}

	// Assign the bucket name as the resource ID
	d.SetId(bucket)
	return resourceObsBucketUpdate(d, meta)
}

func resourceObsBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	obsClient, err := config.ObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	logp.Printf("[DEBUG] Update OBS bucket %s", d.Id())
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

	if d.HasChange("policy") {
		policyClient := obsClient
		format := d.Get("policy_format").(string)
		if format == "obs" {
			policyClient, err = config.ObjectStorageClientWithSignature(GetRegion(d, config))
			if err != nil {
				return fmtp.Errorf("Error creating HuaweiCloud OBS policy client: %s", err)
			}
		}
		if err := resourceObsBucketPolicyUpdate(policyClient, d); err != nil {
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

	if d.HasChanges("encryption", "kms_key_id") {
		if err := resourceObsBucketEncryptionUpdate(obsClient, d); err != nil {
			return err
		}
	}

	if d.HasChange("logging") {
		if err := resourceObsBucketLoggingUpdate(obsClient, d); err != nil {
			return err
		}
	}

	if d.HasChange("quota") {
		if err := resourceObsBucketQuotaUpdate(obsClient, d); err != nil {
			return err
		}
	}

	if d.HasChange("lifecycle_rule") {
		if err := resourceObsBucketLifecycleUpdate(obsClient, d); err != nil {
			return err
		}
	}

	if d.HasChange("website") {
		if err := resourceObsBucketWebsiteUpdate(obsClient, d); err != nil {
			return err
		}
	}

	if d.HasChange("cors_rule") {
		if err := resourceObsBucketCorsUpdate(obsClient, d); err != nil {
			return err
		}
	}

	return resourceObsBucketRead(d, meta)
}

func resourceObsBucketRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	obsClient, err := config.ObjectStorageClient(region)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	logp.Printf("[DEBUG] Read OBS bucket: %s", d.Id())
	_, err = obsClient.HeadBucket(d.Id())
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok && obsError.StatusCode == 404 {
			logp.Printf("[WARN] OBS bucket(%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return fmtp.Errorf("error reading OBS bucket %s: %s", d.Id(), err)
	}

	// for import case
	if _, ok := d.GetOk("bucket"); !ok {
		d.Set("bucket", d.Id())
	}

	d.Set("region", region)
	d.Set("bucket_domain_name", bucketDomainNameWithCloud(d.Get("bucket").(string), region, config.Cloud))

	// Read storage class
	if err := setObsBucketStorageClass(obsClient, d); err != nil {
		return err
	}

	// Read enterprise project id, multi_az and parallel_fs
	if err := setObsBucketMetadata(obsClient, d); err != nil {
		return err
	}

	// Read the versioning
	if err := setObsBucketVersioning(obsClient, d); err != nil {
		return err
	}

	// Read the encryption configuration
	if err := setObsBucketEncryption(obsClient, d); err != nil {
		return err
	}

	// Read the logging configuration
	if err := setObsBucketLogging(obsClient, d); err != nil {
		return err
	}

	// Read the quota
	if err := setObsBucketQuota(obsClient, d); err != nil {
		return err
	}

	// Read the Lifecycle configuration
	if err := setObsBucketLifecycleConfiguration(obsClient, d); err != nil {
		return err
	}

	// Read the website configuration
	if err := setObsBucketWebsiteConfiguration(obsClient, d); err != nil {
		return err
	}

	// Read the CORS rules
	if err := setObsBucketCorsRules(obsClient, d); err != nil {
		return err
	}

	// Read the bucket policy
	policyClient := obsClient
	format := d.Get("policy_format").(string)
	if format == "obs" {
		policyClient, err = config.ObjectStorageClientWithSignature(GetRegion(d, config))
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud OBS policy client: %s", err)
		}
	}
	if err := setObsBucketPolicy(policyClient, d); err != nil {
		return err
	}

	// Read the tags
	if err := setObsBucketTags(obsClient, d); err != nil {
		return err
	}

	return nil
}

func resourceObsBucketDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	obsClient, err := config.ObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	bucket := d.Id()
	logp.Printf("[DEBUG] deleting OBS Bucket: %s", bucket)
	_, err = obsClient.DeleteBucket(bucket)
	if err != nil {
		obsError, ok := err.(obs.ObsError)
		if ok && obsError.Code == "BucketNotEmpty" {
			logp.Printf("[WARN] OBS bucket: %s is not empty", bucket)
			if d.Get("force_destroy").(bool) {
				err = deleteAllBucketObjects(obsClient, bucket)
				if err == nil {
					logp.Printf("[WARN] all objects of %s have been deleted, and try again", bucket)
					return resourceObsBucketDelete(d, meta)
				}
			}
			return err
		}
		return fmtp.Errorf("Error deleting OBS bucket %s, %s", bucket, err)
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
	logp.Printf("[DEBUG] set tags of OBS bucket %s: %#v", bucket, req)

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
	logp.Printf("[DEBUG] set ACL of OBS bucket %s: %#v", bucket, i)

	_, err := obsClient.SetBucketAcl(i)
	if err != nil {
		return getObsError("Error updating acl of OBS bucket", bucket, err)
	}

	// acl policy can not be retrieved by obsClient.GetBucketAcl method
	d.Set("acl", acl)
	return nil
}

func resourceObsBucketPolicyUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	policy := d.Get("policy").(string)

	if policy != "" {
		logp.Printf("[DEBUG] OBS bucket: %s, set policy: %s", bucket, policy)
		params := &obs.SetBucketPolicyInput{
			Bucket: bucket,
			Policy: policy,
		}

		if _, err := obsClient.SetBucketPolicy(params); err != nil {
			return getObsError("Error setting OBS bucket policy", bucket, err)
		}
	} else {
		logp.Printf("[DEBUG] OBS bucket: %s, delete policy", bucket)
		_, err := obsClient.DeleteBucketPolicy(bucket)
		if err != nil {
			return getObsError("Error deleting policy of OBS bucket %s: %s", bucket, err)
		}
	}

	return nil
}

func resourceObsBucketClassUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	class := d.Get("storage_class").(string)

	input := &obs.SetBucketStoragePolicyInput{}
	input.Bucket = bucket
	input.StorageClass = obs.StorageClassType(class)
	logp.Printf("[DEBUG] set storage class of OBS bucket %s: %#v", bucket, input)

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
	logp.Printf("[DEBUG] set versioning of OBS bucket %s: %#v", bucket, input)

	_, err := obsClient.SetBucketVersioning(input)
	if err != nil {
		return getObsError("Error setting versioning status of OBS bucket", bucket, err)
	}

	return nil
}

func resourceObsBucketEncryptionUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)

	if d.Get("encryption").(bool) {
		input := &obs.SetBucketEncryptionInput{}
		input.Bucket = bucket
		input.SSEAlgorithm = obs.DEFAULT_SSE_KMS_ENCRYPTION
		input.KMSMasterKeyID = d.Get("kms_key_id").(string)

		logp.Printf("[DEBUG] enable default encryption of OBS bucket %s: %#v", bucket, input)
		_, err := obsClient.SetBucketEncryption(input)
		if err != nil {
			return getObsError("failed to enable default encryption of OBS bucket", bucket, err)
		}
	} else if !d.IsNewResource() {
		_, err := obsClient.DeleteBucketEncryption(bucket)
		if err != nil {
			return getObsError("failed to disable default encryption of OBS bucket", bucket, err)
		}
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
		if val := c["target_bucket"].(string); val != "" {
			loggingStatus.TargetBucket = val
		}

		if val := c["target_prefix"].(string); val != "" {
			loggingStatus.TargetPrefix = val
		}
	}
	logp.Printf("[DEBUG] set logging of OBS bucket %s: %#v", bucket, loggingStatus)

	_, err := obsClient.SetBucketLoggingConfiguration(loggingStatus)
	if err != nil {
		return getObsError("Error setting logging configuration of OBS bucket", bucket, err)
	}

	return nil
}

func resourceObsBucketQuotaUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	quota := d.Get("quota").(int)
	quotaInput := &obs.SetBucketQuotaInput{}
	quotaInput.Bucket = bucket
	quotaInput.BucketQuota.Quota = int64(quota)

	_, err := obsClient.SetBucketQuota(quotaInput)
	if err != nil {
		return getObsError("Error setting quota of OBS bucket", bucket, err)
	}

	return nil

}

func resourceObsBucketLifecycleUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	lifecycleRules := d.Get("lifecycle_rule").([]interface{})

	if len(lifecycleRules) == 0 {
		logp.Printf("[DEBUG] remove all lifecycle rules of bucket %s", bucket)
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
	logp.Printf("[DEBUG] set lifecycle configurations of OBS bucket %s: %#v", bucket, opts)

	_, err := obsClient.SetBucketLifecycleConfiguration(opts)
	if err != nil {
		return getObsError("Error setting lifecycle rules of OBS bucket", bucket, err)
	}

	return nil
}

func resourceObsBucketWebsiteUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	ws := d.Get("website").([]interface{})

	if len(ws) == 1 {
		var w map[string]interface{}
		if ws[0] != nil {
			w = ws[0].(map[string]interface{})
		} else {
			w = make(map[string]interface{})
		}
		return resourceObsBucketWebsitePut(obsClient, d, w)
	} else if len(ws) == 0 {
		return resourceObsBucketWebsiteDelete(obsClient, d)
	} else {
		return fmtp.Errorf("Cannot specify more than one website.")
	}
}

func resourceObsBucketCorsUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	rawCors := d.Get("cors_rule").([]interface{})

	if len(rawCors) == 0 {
		// Delete CORS
		logp.Printf("[DEBUG] delete CORS rules of OBS bucket: %s", bucket)
		_, err := obsClient.DeleteBucketCors(bucket)
		if err != nil {
			return getObsError("Error deleting CORS rules of OBS bucket", bucket, err)
		}
		return nil
	}

	// set CORS
	rules := make([]obs.CorsRule, 0, len(rawCors))
	for _, cors := range rawCors {
		corsMap := cors.(map[string]interface{})
		r := obs.CorsRule{}
		for k, v := range corsMap {
			if k == "max_age_seconds" {
				r.MaxAgeSeconds = v.(int)
			} else {
				vMap := make([]string, len(v.([]interface{})))
				for i, vv := range v.([]interface{}) {
					vMap[i] = vv.(string)
				}
				switch k {
				case "allowed_headers":
					r.AllowedHeader = vMap
				case "allowed_methods":
					r.AllowedMethod = vMap
				case "allowed_origins":
					r.AllowedOrigin = vMap
				case "expose_headers":
					r.ExposeHeader = vMap
				}
			}
		}
		logp.Printf("[DEBUG] set CORS of OBS bucket %s: %#v", bucket, r)
		rules = append(rules, r)
	}

	corsInput := &obs.SetBucketCorsInput{}
	corsInput.Bucket = bucket
	corsInput.CorsRules = rules
	logp.Printf("[DEBUG] OBS bucket: %s, put CORS: %#v", bucket, corsInput)

	_, err := obsClient.SetBucketCors(corsInput)
	if err != nil {
		return getObsError("Error setting CORS rules of OBS bucket", bucket, err)
	}
	return nil
}

func resourceObsBucketWebsitePut(obsClient *obs.ObsClient, d *schema.ResourceData, website map[string]interface{}) error {
	bucket := d.Get("bucket").(string)

	var indexDocument, errorDocument, redirectAllRequestsTo, routingRules string
	if v, ok := website["index_document"]; ok {
		indexDocument = v.(string)
	}
	if v, ok := website["error_document"]; ok {
		errorDocument = v.(string)
	}
	if v, ok := website["redirect_all_requests_to"]; ok {
		redirectAllRequestsTo = v.(string)
	}
	if v, ok := website["routing_rules"]; ok {
		routingRules = v.(string)
	}

	if indexDocument == "" && redirectAllRequestsTo == "" {
		return fmtp.Errorf("Must specify either index_document or redirect_all_requests_to.")
	}

	websiteConfiguration := &obs.SetBucketWebsiteConfigurationInput{}
	websiteConfiguration.Bucket = bucket

	if indexDocument != "" {
		websiteConfiguration.IndexDocument = obs.IndexDocument{
			Suffix: indexDocument,
		}
	}

	if errorDocument != "" {
		websiteConfiguration.ErrorDocument = obs.ErrorDocument{
			Key: errorDocument,
		}
	}

	if redirectAllRequestsTo != "" {
		redirect, err := url.Parse(redirectAllRequestsTo)
		if err == nil && redirect.Scheme != "" {
			var redirectHostBuf bytes.Buffer
			redirectHostBuf.WriteString(redirect.Host)
			if redirect.Path != "" {
				redirectHostBuf.WriteString(redirect.Path)
			}
			websiteConfiguration.RedirectAllRequestsTo = obs.RedirectAllRequestsTo{
				HostName: redirectHostBuf.String(),
				Protocol: obs.ProtocolType(redirect.Scheme),
			}
		} else {
			websiteConfiguration.RedirectAllRequestsTo = obs.RedirectAllRequestsTo{
				HostName: redirectAllRequestsTo,
			}
		}
	}

	if routingRules != "" {
		var unmarshaledRules []obs.RoutingRule
		if err := json.Unmarshal([]byte(routingRules), &unmarshaledRules); err != nil {
			return err
		}
		websiteConfiguration.RoutingRules = unmarshaledRules
	}

	logp.Printf("[DEBUG] set website configuration of OBS bucket %s: %#v", bucket, websiteConfiguration)
	_, err := obsClient.SetBucketWebsiteConfiguration(websiteConfiguration)
	if err != nil {
		return getObsError("Error updating website configuration of OBS bucket", bucket, err)
	}

	return nil
}

func resourceObsBucketWebsiteDelete(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)

	logp.Printf("[DEBUG] delete website configuration of OBS bucket %s", bucket)
	_, err := obsClient.DeleteBucketWebsiteConfiguration(bucket)
	if err != nil {
		return getObsError("Error deleting website configuration of OBS bucket", bucket, err)
	}

	return nil
}

func setObsBucketStorageClass(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketStoragePolicy(bucket)
	if err != nil {
		logp.Printf("[WARN] Error getting storage class of OBS bucket %s: %s", bucket, err)
	} else {
		class := string(output.StorageClass)
		logp.Printf("[DEBUG] getting storage class of OBS bucket %s: %s", bucket, class)
		d.Set("storage_class", normalizeStorageClass(class))
	}

	return nil
}

func setObsBucketMetadata(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	input := &obs.GetBucketMetadataInput{
		Bucket: bucket,
	}
	output, err := obsClient.GetBucketMetadata(input)
	if err != nil {
		return getObsError("Error getting metadata of OBS bucket", bucket, err)
	}
	logp.Printf("[DEBUG] getting metadata of OBS bucket %s: %#v", bucket, output)

	d.Set("enterprise_project_id", output.Epid)

	if output.AZRedundancy == "3az" {
		d.Set("multi_az", true)
	} else {
		d.Set("multi_az", false)
	}

	if output.FSStatus == "Enabled" {
		d.Set("parallel_fs", true)
	} else {
		d.Set("parallel_fs", false)
		d.Set("bucket_version", output.Version)
	}

	return nil
}

func setObsBucketPolicy(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketPolicy(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.Code == "NoSuchBucketPolicy" {
				d.Set("policy", nil)
				return nil
			}
			return fmtp.Errorf("Error getting policy of OBS bucket %s: %s,\n Reason: %s",
				bucket, obsError.Code, obsError.Message)
		}
		return err
	}

	pol := output.Policy
	logp.Printf("[DEBUG] getting policy of OBS bucket %s: %s", bucket, pol)
	d.Set("policy", pol)

	return nil
}

func setObsBucketVersioning(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketVersioning(bucket)
	if err != nil {
		return getObsError("Error getting versioning status of OBS bucket", bucket, err)
	}

	logp.Printf("[DEBUG] getting versioning status of OBS bucket %s: %s", bucket, output.Status)
	if output.Status == obs.VersioningStatusEnabled {
		d.Set("versioning", true)
	} else {
		d.Set("versioning", false)
	}

	return nil
}

func setObsBucketEncryption(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketEncryption(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.Code == "NoSuchEncryptionConfiguration" {
				d.Set("encryption", false)
				d.Set("kms_key_id", nil)
				return nil
			}
			return fmtp.Errorf("Error getting encryption configuration of OBS bucket %s: %s,\n Reason: %s",
				bucket, obsError.Code, obsError.Message)
		}
		return err
	}

	logp.Printf("[DEBUG] getting encryption configuration of OBS bucket %s: %+v", bucket, output.BucketEncryptionConfiguration)
	if output.SSEAlgorithm != "" {
		d.Set("encryption", true)
		d.Set("kms_key_id", output.KMSMasterKeyID)
	} else {
		d.Set("encryption", false)
		d.Set("kms_key_id", nil)
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
		logging["target_bucket"] = output.TargetBucket
		if output.TargetPrefix != "" {
			logging["target_prefix"] = output.TargetPrefix
		}
		lcList = append(lcList, logging)
	}
	logp.Printf("[DEBUG] getting logging configuration of OBS bucket %s: %#v", bucket, lcList)

	if err := d.Set("logging", lcList); err != nil {
		return fmtp.Errorf("Error saving logging configuration of OBS bucket %s: %s", bucket, err)
	}

	return nil
}

func setObsBucketQuota(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketQuota(bucket)
	if err != nil {
		return getObsError("Error getting quota of OBS bucket", bucket, err)
	}

	logp.Printf("[DEBUG] getting quota of OBS bucket %s: %d", bucket, output.Quota)

	d.Set("quota", output.Quota)

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
			}
			return fmtp.Errorf("Error getting lifecycle configuration of OBS bucket %s: %s,\n Reason: %s",
				bucket, obsError.Code, obsError.Message)
		}
		return err
	}

	rawRules := output.LifecycleRules
	logp.Printf("[DEBUG] getting original lifecycle configuration of OBS bucket %s, lifecycle: %#v", bucket, rawRules)

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

	logp.Printf("[DEBUG] saving lifecycle configuration of OBS bucket %s, lifecycle: %#v", bucket, rules)
	if err := d.Set("lifecycle_rule", rules); err != nil {
		return fmtp.Errorf("Error saving lifecycle configuration of OBS bucket %s: %s", bucket, err)
	}

	return nil
}

func setObsBucketWebsiteConfiguration(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketWebsiteConfiguration(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.Code == "NoSuchWebsiteConfiguration" {
				d.Set("website", nil)
				return nil
			}
			return fmtp.Errorf("Error getting website configuration of OBS bucket %s: %s,\n Reason: %s",
				bucket, obsError.Code, obsError.Message)
		}
		return err
	}

	logp.Printf("[DEBUG] getting original website configuration of OBS bucket %s, output: %#v", bucket, output.BucketWebsiteConfiguration)
	var websites []map[string]interface{}
	w := make(map[string]interface{})

	w["index_document"] = output.IndexDocument.Suffix
	w["error_document"] = output.ErrorDocument.Key

	// redirect_all_requests_to
	v := output.RedirectAllRequestsTo
	if string(v.Protocol) == "" {
		w["redirect_all_requests_to"] = v.HostName
	} else {
		var host string
		var path string
		parsedHostName, err := url.Parse(v.HostName)
		if err == nil {
			host = parsedHostName.Host
			path = parsedHostName.Path
		} else {
			host = v.HostName
			path = ""
		}

		w["redirect_all_requests_to"] = (&url.URL{
			Host:   host,
			Path:   path,
			Scheme: string(v.Protocol),
		}).String()
	}

	// routing_rules
	rawRules := output.RoutingRules
	if len(rawRules) > 0 {
		rr, err := normalizeWebsiteRoutingRules(rawRules)
		if err != nil {
			return fmtp.Errorf("Error while marshaling website routing rules: %s", err)
		}
		w["routing_rules"] = rr
	}

	websites = append(websites, w)
	logp.Printf("[DEBUG] saving website configuration of OBS bucket %s, website: %#v", bucket, websites)
	if err := d.Set("website", websites); err != nil {
		return fmtp.Errorf("Error saving website configuration of OBS bucket %s: %s", bucket, err)
	}
	return nil
}

func setObsBucketCorsRules(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketCors(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.Code == "NoSuchCORSConfiguration" {
				d.Set("cors_rule", nil)
				return nil
			}
			return fmtp.Errorf("Error getting CORS configuration of OBS bucket %s: %s,\n Reason: %s",
				bucket, obsError.Code, obsError.Message)
		}
		return err
	}

	corsRules := output.CorsRules
	logp.Printf("[DEBUG] getting original CORS rules of OBS bucket %s, CORS: %#v", bucket, corsRules)

	rules := make([]map[string]interface{}, 0, len(corsRules))
	for _, ruleObject := range corsRules {
		rule := make(map[string]interface{})
		rule["allowed_origins"] = ruleObject.AllowedOrigin
		rule["allowed_methods"] = ruleObject.AllowedMethod
		rule["max_age_seconds"] = ruleObject.MaxAgeSeconds
		if ruleObject.AllowedHeader != nil {
			rule["allowed_headers"] = ruleObject.AllowedHeader
		}
		if ruleObject.ExposeHeader != nil {
			rule["expose_headers"] = ruleObject.ExposeHeader
		}

		rules = append(rules, rule)
	}

	logp.Printf("[DEBUG] saving CORS rules of OBS bucket %s, CORS: %#v", bucket, rules)
	if err := d.Set("cors_rule", rules); err != nil {
		return fmtp.Errorf("Error saving CORS rules of OBS bucket %s: %s", bucket, err)
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
			}
			return fmtp.Errorf("Error getting tags of OBS bucket %s: %s,\n Reason: %s",
				bucket, obsError.Code, obsError.Message)
		}
		return err
	}

	tagmap := make(map[string]string)
	for _, tag := range output.Tags {
		tagmap[tag.Key] = tag.Value
	}
	logp.Printf("[DEBUG] getting tags of OBS bucket %s: %#v", bucket, tagmap)
	if err := d.Set("tags", tagmap); err != nil {
		return fmtp.Errorf("Error saving tags of OBS bucket %s: %s", bucket, err)
	}
	return nil
}

func deleteAllBucketObjects(obsClient *obs.ObsClient, bucket string) error {
	listOpts := &obs.ListObjectsInput{
		Bucket: bucket,
	}
	// list all objects
	resp, err := obsClient.ListObjects(listOpts)
	if err != nil {
		return getObsError("Error listing objects of OBS bucket", bucket, err)
	}

	objects := make([]obs.ObjectToDelete, len(resp.Contents))
	for i, content := range resp.Contents {
		objects[i].Key = content.Key
	}

	deleteOpts := &obs.DeleteObjectsInput{
		Bucket:  bucket,
		Objects: objects,
	}
	logp.Printf("[DEBUG] objects of %s will be deleted: %v", bucket, objects)
	output, err := obsClient.DeleteObjects(deleteOpts)
	if err != nil {
		return getObsError("Error deleting all objects of OBS bucket", bucket, err)
	}
	if len(output.Errors) > 0 {
		return fmtp.Errorf("Error some objects are still exist in %s: %#v", bucket, output.Errors)
	}
	return nil
}

func expirationHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	if v, ok := m["storage_class"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	return hashcode.String(buf.String())
}

func getObsError(action string, bucket string, err error) error {
	if obsError, ok := err.(obs.ObsError); ok {
		return fmtp.Errorf("%s %s: %s,\n Reason: %s", action, bucket, obsError.Code, obsError.Message)
	}
	return err
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

func normalizeWebsiteRoutingRules(w []obs.RoutingRule) (string, error) {
	// transform []obs.RoutingRule to []WebsiteRoutingRule
	wrules := make([]WebsiteRoutingRule, 0, len(w))
	for _, rawRule := range w {
		rule := WebsiteRoutingRule{
			Condition: Condition{
				KeyPrefixEquals:             rawRule.Condition.KeyPrefixEquals,
				HttpErrorCodeReturnedEquals: rawRule.Condition.HttpErrorCodeReturnedEquals,
			},
			Redirect: Redirect{
				Protocol:             string(rawRule.Redirect.Protocol),
				HostName:             rawRule.Redirect.HostName,
				HttpRedirectCode:     rawRule.Redirect.HttpRedirectCode,
				ReplaceKeyWith:       rawRule.Redirect.ReplaceKeyWith,
				ReplaceKeyPrefixWith: rawRule.Redirect.ReplaceKeyPrefixWith,
			},
		}
		wrules = append(wrules, rule)
	}

	// normalize
	withNulls, err := json.Marshal(wrules)
	if err != nil {
		return "", err
	}

	var rules []map[string]interface{}
	if err := json.Unmarshal(withNulls, &rules); err != nil {
		return "", err
	}

	var cleanRules []map[string]interface{}
	for _, rule := range rules {
		cleanRules = append(cleanRules, utils.RemoveNil(rule))
	}

	withoutNulls, err := json.Marshal(cleanRules)
	if err != nil {
		return "", err
	}

	return string(withoutNulls), nil
}

func bucketDomainNameWithCloud(bucket, region, cloud string) string {
	return fmt.Sprintf("%s.obs.%s.%s", bucket, region, cloud)
}

type Condition struct {
	KeyPrefixEquals             string `json:"KeyPrefixEquals,omitempty"`
	HttpErrorCodeReturnedEquals string `json:"HttpErrorCodeReturnedEquals,omitempty"`
}

type Redirect struct {
	Protocol             string `json:"Protocol,omitempty"`
	HostName             string `json:"HostName,omitempty"`
	ReplaceKeyPrefixWith string `json:"ReplaceKeyPrefixWith,omitempty"`
	ReplaceKeyWith       string `json:"ReplaceKeyWith,omitempty"`
	HttpRedirectCode     string `json:"HttpRedirectCode,omitempty"`
}
type WebsiteRoutingRule struct {
	Condition Condition `json:"Condition,omitempty"`
	Redirect  Redirect  `json:"Redirect"`
}

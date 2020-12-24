package huaweicloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk/openstack/obs"
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
				ValidateFunc:     validateJsonString,
				DiffSuppressFunc: suppressEquivalentAwsPolicyDiffs,
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
							ValidateFunc: validateJsonString,
							StateFunc: func(v interface{}) string {
								json, _ := normalizeJsonString(v)
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

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

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

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
	region := GetRegion(d, config)
	obsClient, err := config.NewObjectStorageClient(region)
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
		Epid:         GetEnterpriseProjectID(d, config),
	}
	opts.Location = region
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
	obsClient, err := config.NewObjectStorageClient(GetRegion(d, config))
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

	if d.HasChange("policy") {
		policyClient := obsClient
		format := d.Get("policy_format").(string)
		if format == "obs" {
			policyClient, err = config.NewObjectStorageClientWithSignature(GetRegion(d, config))
			if err != nil {
				return fmt.Errorf("Error creating HuaweiCloud OBS policy client: %s", err)
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
	config := meta.(*Config)
	region := GetRegion(d, config)
	obsClient, err := config.NewObjectStorageClient(region)
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

	// for import case
	if _, ok := d.GetOk("bucket"); !ok {
		d.Set("bucket", d.Id())
	}

	d.Set("region", region)
	d.Set("bucket_domain_name", bucketDomainName(d.Get("bucket").(string), region))

	// Read storage class
	if err := setObsBucketStorageClass(obsClient, d); err != nil {
		return err
	}

	// Read  enterprise project id
	if err := setObsBucketEnterpriseProjectID(obsClient, d); err != nil {
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
		policyClient, err = config.NewObjectStorageClientWithSignature(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud OBS policy client: %s", err)
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
	config := meta.(*Config)
	obsClient, err := config.NewObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	bucket := d.Id()
	log.Printf("[DEBUG] deleting OBS Bucket: %s", bucket)
	_, err = obsClient.DeleteBucket(bucket)
	if err != nil {
		obsError, ok := err.(obs.ObsError)
		if ok && obsError.Code == "BucketNotEmpty" {
			log.Printf("[WARN] OBS bucket: %s is not empty", bucket)
			if d.Get("force_destroy").(bool) {
				err = deleteAllBucketObjects(obsClient, bucket)
				if err == nil {
					log.Printf("[WARN] all objects of %s have been deleted, and try again", bucket)
					return resourceObsBucketDelete(d, meta)
				}
			}
			return err
		}
		return fmt.Errorf("Error deleting OBS bucket %s, %s", bucket, err)
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

func resourceObsBucketPolicyUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	policy := d.Get("policy").(string)

	if policy != "" {
		log.Printf("[DEBUG] OBS bucket: %s, set policy: %s", bucket, policy)
		params := &obs.SetBucketPolicyInput{
			Bucket: bucket,
			Policy: policy,
		}

		if _, err := obsClient.SetBucketPolicy(params); err != nil {
			return getObsError("Error setting OBS bucket policy", bucket, err)
		}
	} else {
		log.Printf("[DEBUG] OBS bucket: %s, delete policy", bucket)
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
		return getObsError("Error setting versioning status of OBS bucket", bucket, err)
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
	log.Printf("[DEBUG] set logging of OBS bucket %s: %#v", bucket, loggingStatus)

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
		return fmt.Errorf("Cannot specify more than one website.")
	}
}

func resourceObsBucketCorsUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	rawCors := d.Get("cors_rule").([]interface{})

	if len(rawCors) == 0 {
		// Delete CORS
		log.Printf("[DEBUG] delete CORS rules of OBS bucket: %s", bucket)
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
		log.Printf("[DEBUG] set CORS of OBS bucket %s: %#v", bucket, r)
		rules = append(rules, r)
	}

	corsInput := &obs.SetBucketCorsInput{}
	corsInput.Bucket = bucket
	corsInput.CorsRules = rules
	log.Printf("[DEBUG] OBS bucket: %s, put CORS: %#v", bucket, corsInput)

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
		return fmt.Errorf("Must specify either index_document or redirect_all_requests_to.")
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

	log.Printf("[DEBUG] set website configuration of OBS bucket %s: %#v", bucket, websiteConfiguration)
	_, err := obsClient.SetBucketWebsiteConfiguration(websiteConfiguration)
	if err != nil {
		return getObsError("Error updating website configuration of OBS bucket", bucket, err)
	}

	return nil
}

func resourceObsBucketWebsiteDelete(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)

	log.Printf("[DEBUG] delete website configuration of OBS bucket %s", bucket)
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
		return getObsError("Error getting storage class of OBS bucket", bucket, err)
	}

	class := string(output.StorageClass)
	log.Printf("[DEBUG] getting storage class of OBS bucket %s: %s", bucket, class)
	d.Set("storage_class", normalizeStorageClass(class))

	return nil
}

func setObsBucketEnterpriseProjectID(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	input := &obs.GetBucketMetadataInput{
		Bucket: bucket,
	}
	output, err := obsClient.GetBucketMetadata(input)
	if err != nil {
		return getObsError("Error getting metadata of OBS bucket", bucket, err)
	}

	epsId := string(output.Epid)
	log.Printf("[DEBUG] getting enterprise project id of OBS bucket %s: %s", bucket, epsId)
	d.Set("enterprise_project_id", epsId)

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
			return fmt.Errorf("Error getting policy of OBS bucket %s: %s,\n Reason: %s",
				bucket, obsError.Code, obsError.Message)
		}
		return err
	}

	pol := output.Policy
	log.Printf("[DEBUG] getting policy of OBS bucket %s: %s", bucket, pol)
	d.Set("policy", pol)

	return nil
}

func setObsBucketVersioning(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketVersioning(bucket)
	if err != nil {
		return getObsError("Error getting versioning status of OBS bucket", bucket, err)
	}

	log.Printf("[DEBUG] getting versioning status of OBS bucket %s: %s", bucket, output.Status)
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
		logging["target_bucket"] = output.TargetBucket
		if output.TargetPrefix != "" {
			logging["target_prefix"] = output.TargetPrefix
		}
		lcList = append(lcList, logging)
	}
	log.Printf("[DEBUG] getting logging configuration of OBS bucket %s: %#v", bucket, lcList)

	if err := d.Set("logging", lcList); err != nil {
		return fmt.Errorf("Error saving logging configuration of OBS bucket %s: %s", bucket, err)
	}

	return nil
}

func setObsBucketQuota(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketQuota(bucket)
	if err != nil {
		return getObsError("Error getting quota of OBS bucket", bucket, err)
	}

	log.Printf("[DEBUG] getting quota of OBS bucket %s: %d", bucket, output.Quota)

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
			return fmt.Errorf("Error getting lifecycle configuration of OBS bucket %s: %s,\n Reason: %s",
				bucket, obsError.Code, obsError.Message)
		}
		return err
	}

	rawRules := output.LifecycleRules
	log.Printf("[DEBUG] getting original lifecycle configuration of OBS bucket %s, lifecycle: %#v", bucket, rawRules)

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

	log.Printf("[DEBUG] saving lifecycle configuration of OBS bucket %s, lifecycle: %#v", bucket, rules)
	if err := d.Set("lifecycle_rule", rules); err != nil {
		return fmt.Errorf("Error saving lifecycle configuration of OBS bucket %s: %s", bucket, err)
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
			return fmt.Errorf("Error getting website configuration of OBS bucket %s: %s,\n Reason: %s",
				bucket, obsError.Code, obsError.Message)
		}
		return err
	}

	log.Printf("[DEBUG] getting original website configuration of OBS bucket %s, output: %#v", bucket, output.BucketWebsiteConfiguration)
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
			return fmt.Errorf("Error while marshaling website routing rules: %s", err)
		}
		w["routing_rules"] = rr
	}

	websites = append(websites, w)
	log.Printf("[DEBUG] saving website configuration of OBS bucket %s, website: %#v", bucket, websites)
	if err := d.Set("website", websites); err != nil {
		return fmt.Errorf("Error saving website configuration of OBS bucket %s: %s", bucket, err)
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
			return fmt.Errorf("Error getting CORS configuration of OBS bucket %s: %s,\n Reason: %s",
				bucket, obsError.Code, obsError.Message)
		}
		return err
	}

	corsRules := output.CorsRules
	log.Printf("[DEBUG] getting original CORS rules of OBS bucket %s, CORS: %#v", bucket, corsRules)

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

	log.Printf("[DEBUG] saving CORS rules of OBS bucket %s, CORS: %#v", bucket, rules)
	if err := d.Set("cors_rule", rules); err != nil {
		return fmt.Errorf("Error saving CORS rules of OBS bucket %s: %s", bucket, err)
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
			return fmt.Errorf("Error getting tags of OBS bucket %s: %s,\n Reason: %s",
				bucket, obsError.Code, obsError.Message)
		}
		return err
	}

	tagmap := make(map[string]string)
	for _, tag := range output.Tags {
		tagmap[tag.Key] = tag.Value
	}
	log.Printf("[DEBUG] getting tags of OBS bucket %s: %#v", bucket, tagmap)
	if err := d.Set("tags", tagmap); err != nil {
		return fmt.Errorf("Error saving tags of OBS bucket %s: %s", bucket, err)
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
	log.Printf("[DEBUG] objects of %s will be deleted: %v", bucket, objects)
	output, err := obsClient.DeleteObjects(deleteOpts)
	if err != nil {
		return getObsError("Error deleting all objects of OBS bucket", bucket, err)
	}
	if len(output.Errors) > 0 {
		return fmt.Errorf("Error some objects are still exist in %s: %#v", bucket, output.Errors)
	}
	return nil
}

func getObsError(action string, bucket string, err error) error {
	if obsError, ok := err.(obs.ObsError); ok {
		return fmt.Errorf("%s %s: %s,\n Reason: %s", action, bucket, obsError.Code, obsError.Message)
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
		cleanRules = append(cleanRules, removeNil(rule))
	}

	withoutNulls, err := json.Marshal(cleanRules)
	if err != nil {
		return "", err
	}

	return string(withoutNulls), nil
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

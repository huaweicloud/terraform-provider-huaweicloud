package obs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/obs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API OBS PUT ?acl
// @API OBS DELETE ?lifecycle
// @API OBS PUT ?lifecycle
// @API OBS GET ?lifecycle
// @API OBS PUT ?website
// @API OBS DELETE ?website
// @API OBS GET ?website
// @API OBS PUT ?customdomain
// @API OBS DELETE ?customdomain
// @API OBS GET ?customdomain
// @API OBS DELETE /
// @API OBS PUT /
// @API OBS HEAD /
// @API OBS GET /
// @API OBS PUT ?versioning
// @API OBS GET ?versioning
// @API OBS PUT ?quota
// @API OBS GET ?quota
// @API OBS POST ?delete
// @API OBS PUT ?tagging
// @API OBS GET ?tagging
// @API OBS GET ?storageinfo
// @API OBS PUT ?storageClass
// @API OBS GET ?storageClass
// @API OBS PUT ?encryption
// @API OBS DELETE ?encryption
// @API OBS GET ?encryption
// @API OBS PUT ?policy
// @API OBS DELETE ?policy
// @API OBS GET ?policy
// @API OBS PUT ?logging
// @API OBS GET ?logging
// @API OBS DELETE ?cors
// @API OBS PUT ?cors
// @API OBS GET ?cors
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrate
func ResourceObsBucket() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObsBucketCreate,
		ReadContext:   resourceObsBucketRead,
		UpdateContext: resourceObsBucketUpdate,
		DeleteContext: resourceObsBucketDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceObsBucketImport,
		},

		CustomizeDiff: config.MergeDefaultTags(),

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
			},

			"acl": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "private",
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
						"agency": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "schema: Required",
						},
					},
				},
			},

			"quota": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"storage_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"object_number": {
							Type:     schema.TypeInt,
							Computed: true,
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
						"abort_incomplete_multipart_upload": {
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
								jsonString, _ := utils.NormalizeJsonString(v)
								return jsonString
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

			"tags": common.TagsSchema(),
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
			"sse_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kms_key_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_domain_names": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
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

func resourceObsBucketCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	obsClient, err := conf.ObjectStorageClient(region)
	if err != nil {
		return diag.Errorf("Error creating OBS client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	acl := d.Get("acl").(string)
	class := d.Get("storage_class").(string)
	opts := &obs.CreateBucketInput{
		Bucket:            bucket,
		ACL:               obs.AclType(acl),
		StorageClass:      obs.StorageClassType(class),
		IsFSFileInterface: d.Get("parallel_fs").(bool),
		Epid:              conf.GetEnterpriseProjectID(d),
	}
	opts.Location = region
	if _, ok := d.GetOk("multi_az"); ok {
		opts.AvailableZone = "3az"
	}

	log.Printf("[DEBUG] OBS bucket create opts: %#v", opts)
	_, err = obsClient.CreateBucket(opts)
	if err != nil {
		return diag.FromErr(getObsError("Error creating bucket", bucket, err))
	}

	// Assign the bucket name as the resource ID
	d.SetId(bucket)
	return resourceObsBucketUpdate(ctx, d, meta)
}

func resourceObsBucketUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	obsClient, err := conf.ObjectStorageClient(region)
	if err != nil {
		return diag.Errorf("Error creating OBS client: %s", err)
	}

	obsClientWithSignature, err := conf.ObjectStorageClientWithSignature(region)
	if err != nil {
		return diag.Errorf("Error creating OBS client with signature: %s", err)
	}

	log.Printf("[DEBUG] Update OBS bucket %s", d.Id())
	if d.HasChange("acl") && !d.IsNewResource() {
		if err := updateObsBucketAcl(obsClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("storage_class") && !d.IsNewResource() {
		if err := resourceObsBucketClassUpdate(obsClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("policy") {
		policyClient := obsClientWithSignature
		if d.Get("policy_format").(string) != "obs" {
			policyClient = obsClient
		}
		if err := resourceObsBucketPolicyUpdate(policyClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		if err := resourceObsBucketTagsUpdate(obsClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("versioning") {
		if err := resourceObsBucketVersioningUpdate(obsClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("encryption", "sse_algorithm", "kms_key_id", "kms_key_project_id") {
		if err := resourceObsBucketEncryptionUpdate(conf, obsClientWithSignature, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("logging") {
		if err := resourceObsBucketLoggingUpdate(obsClientWithSignature, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("quota") {
		if err := resourceObsBucketQuotaUpdate(obsClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("lifecycle_rule") {
		if err := resourceObsBucketLifecycleUpdate(obsClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("website") {
		if err := resourceObsBucketWebsiteUpdate(obsClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("cors_rule") {
		if err := resourceObsBucketCorsUpdate(obsClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_id") && !d.IsNewResource() {
		// the API Limitations: still requires `project_id` field when migrating the EPS of OBS bucket
		if err := resourceObsBucketEnterpriseProjectIdUpdate(ctx, d, conf, obsClient, region); err != nil {
			return diag.FromErr(err)
		}

	}

	if d.HasChange("user_domain_names") {
		if err := resourceObsBucketUserDomainNamesUpdate(obsClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceObsBucketRead(ctx, d, meta)
}

func resourceObsBucketRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	obsClient, err := conf.ObjectStorageClient(region)
	if err != nil {
		return diag.Errorf("Error creating OBS client: %s", err)
	}

	obsClientWithSignature, err := conf.ObjectStorageClientWithSignature(region)
	if err != nil {
		return diag.Errorf("Error creating OBS client with signature: %s", err)
	}

	bucket := d.Id()
	log.Printf("[DEBUG] Read OBS bucket: %s", bucket)
	_, err = obsClient.HeadBucket(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok && obsError.StatusCode == 404 {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Resource not found",
					Detail:   fmt.Sprintf("OBS bucket(%s) not found", bucket),
				},
			}
		}
		return diag.Errorf("error reading OBS bucket %s: %s", bucket, err)
	}

	mErr := &multierror.Error{}
	// for import case
	if _, ok := d.GetOk("bucket"); !ok {
		mErr = multierror.Append(mErr, d.Set("bucket", bucket))
	}

	mErr = multierror.Append(mErr,
		d.Set("region", region),
		d.Set("bucket_domain_name", bucketDomainNameWithCloud(d.Get("bucket").(string), region, conf.Cloud)),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting OBS attributes: %s", mErr)
	}

	// Read storage class
	if err := setObsBucketStorageClass(obsClient, d); err != nil {
		return diag.FromErr(err)
	}

	// Read enterprise project id, multi_az and parallel_fs
	if err := setObsBucketMetadata(obsClient, d); err != nil {
		return diag.FromErr(err)
	}

	// Read the versioning
	if err := setObsBucketVersioning(obsClient, d); err != nil {
		return diag.FromErr(err)
	}

	// Read the encryption configuration
	if err := setObsBucketEncryption(obsClientWithSignature, d); err != nil {
		return diag.FromErr(err)
	}

	// Read the logging configuration
	if err := setObsBucketLogging(obsClientWithSignature, d); err != nil {
		return diag.FromErr(err)
	}

	// Read the quota
	if err := setObsBucketQuota(obsClient, d); err != nil {
		return diag.FromErr(err)
	}

	// Read the Lifecycle configuration
	if err := setObsBucketLifecycleConfiguration(obsClient, d); err != nil {
		return diag.FromErr(err)
	}

	// Read the website configuration
	if err := setObsBucketWebsiteConfiguration(obsClient, d); err != nil {
		return diag.FromErr(err)
	}

	// Read the CORS rules
	if err := setObsBucketCorsRules(obsClient, d); err != nil {
		return diag.FromErr(err)
	}

	// Read the bucket policy
	policyClient := obsClient
	format := d.Get("policy_format").(string)
	if format == "obs" {
		policyClient = obsClientWithSignature
	}
	if err := setObsBucketPolicy(policyClient, d); err != nil {
		return diag.FromErr(err)
	}

	// Read the tags
	if err := setObsBucketTags(obsClient, d); err != nil {
		return diag.FromErr(err)
	}

	// Read the storage info
	if err := setObsBucketStorageInfo(obsClient, d); err != nil {
		return diag.FromErr(err)
	}

	if err := setObsBucketUserDomainNames(obsClient, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceObsBucketDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	obsClient, err := conf.ObjectStorageClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating OBS client: %s", err)
	}

	bucket := d.Id()
	log.Printf("[DEBUG] deleting OBS Bucket: %s", bucket)
	_, err = obsClient.DeleteBucket(bucket)
	if err != nil {
		obsError, ok := err.(obs.ObsError)
		if !ok {
			return diag.Errorf("Error deleting OBS bucket %s, %s", bucket, err)
		}
		if obsError.StatusCode == 404 {
			return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "OBS bucket")
		}
		if obsError.Code == "BucketNotEmpty" {
			log.Printf("[WARN] OBS bucket: %s is not empty", bucket)
			if d.Get("force_destroy").(bool) {
				err = deleteAllBucketObjects(obsClient, bucket)
				if err == nil {
					log.Printf("[WARN] all objects of %s have been deleted, and try again", bucket)
					return resourceObsBucketDelete(ctx, d, meta)
				}
			}
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceObsBucketTagsUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	tagMap := d.Get("tags").(map[string]interface{})
	var tagList []obs.Tag
	for k, v := range tagMap {
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

func updateObsBucketAcl(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	acl := d.Get("acl").(string)

	input := &obs.SetBucketAclInput{
		Bucket: bucket,
		ACL:    obs.AclType(acl),
	}
	log.Printf("[DEBUG] set ACL of OBS bucket %s: %#v", bucket, input)

	_, err := obsClient.SetBucketAcl(input)
	if err != nil {
		return getObsError("Error updating acl of OBS bucket", bucket, err)
	}

	// acl policy can not be retrieved by obsClient.GetBucketAcl method
	if err := d.Set("acl", acl); err != nil {
		return fmt.Errorf("error saving acl of OBS bucket %s: %s", bucket, err)
	}
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
			return getObsError("Error deleting policy of OBS bucket", bucket, err)
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

func resourceObsBucketEncryptionUpdate(config *config.Config, obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)

	if d.Get("encryption").(bool) {
		input := &obs.SetBucketEncryptionInput{}
		input.Bucket = bucket

		if v, ok := d.GetOk("sse_algorithm"); ok {
			input.SSEAlgorithm = v.(string)
		} else {
			input.SSEAlgorithm = obs.DEFAULT_SSE_KMS_ENCRYPTION_OBS
		}

		if input.SSEAlgorithm == obs.DEFAULT_SSE_KMS_ENCRYPTION_OBS {
			if raw, ok := d.GetOk("kms_key_id"); ok {
				input.KMSMasterKeyID = raw.(string)
				input.ProjectID = d.Get("kms_key_project_id").(string)
			}
		}

		log.Printf("[DEBUG] enable default encryption of OBS bucket %s: %#v", bucket, input)
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

		if val := c["agency"].(string); val != "" {
			loggingStatus.Agency = val
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
		ncExpiration := d.Get(fmt.Sprintf("lifecycle_rule.%d.noncurrent_version_expiration", i)).(*schema.Set).List()
		if len(ncExpiration) > 0 {
			raw := ncExpiration[0].(map[string]interface{})
			ncExp := &rules[i].NoncurrentVersionExpiration

			if val, ok := raw["days"].(int); ok && val > 0 {
				ncExp.NoncurrentDays = val
			}
		}

		// AbortIncompleteMultipartUpload
		abortIncompleteMultipartUpload := d.Get(fmt.Sprintf("lifecycle_rule.%d.abort_incomplete_multipart_upload",
			i)).(*schema.Set).List()
		if len(abortIncompleteMultipartUpload) > 0 {
			raw := abortIncompleteMultipartUpload[0].(map[string]interface{})
			abincomMultipartUpload := &rules[i].AbortIncompleteMultipartUpload

			if val, ok := raw["days"].(int); ok && val > 0 {
				abincomMultipartUpload.DaysAfterInitiation = val
			}
		}

		// NoncurrentVersionTransition
		ncTransitions := d.Get(fmt.Sprintf("lifecycle_rule.%d.noncurrent_version_transition", i)).([]interface{})
		ncList := make([]obs.NoncurrentVersionTransition, len(ncTransitions))
		for j, ncTran := range ncTransitions {
			raw := ncTran.(map[string]interface{})

			if val, ok := raw["days"].(int); ok && val > 0 {
				ncList[j].NoncurrentDays = val
			}
			if val, ok := raw["storage_class"].(string); ok {
				ncList[j].StorageClass = obs.StorageClassType(val)
			}
		}
		rules[i].NoncurrentVersionTransitions = ncList
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
	}
	if len(ws) == 0 {
		return resourceObsBucketWebsiteDelete(obsClient, d)
	}
	return fmt.Errorf("cannot specify more than one website")
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

func resourceObsBucketUserDomainNamesUpdate(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	oldRaws, newRaws := d.GetChange("user_domain_names")
	addRaws := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
	removeRaws := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))

	if err := deleteObsBucketUserDomainNames(obsClient, bucket, removeRaws); err != nil {
		return err
	}
	return createObsBucketUserDomainNames(obsClient, bucket, addRaws)
}

func createObsBucketUserDomainNames(obsClient *obs.ObsClient, bucket string, domainNameSet *schema.Set) error {
	for _, domainName := range domainNameSet.List() {
		input := &obs.SetBucketCustomDomainInput{
			Bucket:       bucket,
			CustomDomain: domainName.(string),
		}
		_, err := obsClient.SetBucketCustomDomain(input)
		if err != nil {
			return getObsError("error setting user domain name of OBS bucket", bucket, err)
		}
	}
	return nil
}

func deleteObsBucketUserDomainNames(obsClient *obs.ObsClient, bucket string, domainNameSet *schema.Set) error {
	for _, domainName := range domainNameSet.List() {
		input := &obs.DeleteBucketCustomDomainInput{
			Bucket:       bucket,
			CustomDomain: domainName.(string),
		}
		_, err := obsClient.DeleteBucketCustomDomain(input)
		if err != nil {
			return getObsError("error deleting user domain name of OBS bucket", bucket, err)
		}
	}
	return nil
}

func resourceObsBucketEnterpriseProjectIdUpdate(ctx context.Context, d *schema.ResourceData, cfg *config.Config,
	obsClient *obs.ObsClient, region string) error {
	var (
		projectId   = cfg.GetProjectID(region)
		bucket      = d.Get("bucket").(string)
		migrateOpts = config.MigrateResourceOpts{
			ResourceId:   bucket,
			ResourceType: "bucket",
			RegionId:     region,
			ProjectId:    projectId,
		}
	)
	err := cfg.MigrateEnterpriseProjectWithoutWait(d, migrateOpts)
	if err != nil {
		return err
	}

	// After the EPS service side updates enterprise project ID, it will take a few time to wait the OBS service
	// read the data back into the database.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Success"},
		Refresh:      waitForOBSEnterpriseProjectIdChanged(obsClient, bucket, d.Get("enterprise_project_id").(string)),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return getObsError("error waiting for obs enterprise project ID changed", bucket, err)
	}
	return nil
}

func waitForOBSEnterpriseProjectIdChanged(obsClient *obs.ObsClient, bucket string, enterpriseProjectId string) resource.StateRefreshFunc {
	return func() (result interface{}, state string, err error) {
		input := &obs.GetBucketMetadataInput{
			Bucket: bucket,
		}
		output, err := obsClient.GetBucketMetadata(input)
		if err != nil {
			return nil, "Error", err
		}

		if output.Epid == enterpriseProjectId {
			log.Printf("[DEBUG] the Enterprise Project ID of bucket %s is migrated to %s", bucket, enterpriseProjectId)
			return output, "Success", nil
		}

		return output, "Pending", nil
	}
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
		return fmt.Errorf("must specify either index_document or redirect_all_requests_to")
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
		var unmarshalRules []obs.RoutingRule
		if err := json.Unmarshal([]byte(routingRules), &unmarshalRules); err != nil {
			return err
		}
		websiteConfiguration.RoutingRules = unmarshalRules
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
		log.Printf("[WARN] Error getting storage class of OBS bucket %s: %s", bucket, err)
	} else {
		class := output.StorageClass
		log.Printf("[DEBUG] getting storage class of OBS bucket %s: %s", bucket, class)
		if err := d.Set("storage_class", normalizeStorageClass(class)); err != nil {
			return fmt.Errorf("error saving storage class of OBS bucket %s: %s", bucket, err)
		}
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
	log.Printf("[DEBUG] getting metadata of OBS bucket %s: %#v", bucket, output)

	mErr := multierror.Append(nil, d.Set("enterprise_project_id", output.Epid))

	if output.AZRedundancy == "3az" {
		mErr = multierror.Append(mErr, d.Set("multi_az", true))
	} else {
		mErr = multierror.Append(mErr, d.Set("multi_az", false))
	}

	if output.FSStatus == "Enabled" {
		mErr = multierror.Append(mErr, d.Set("parallel_fs", true))
	} else {
		mErr = multierror.Append(mErr,
			d.Set("parallel_fs", false),
			d.Set("bucket_version", output.Version),
		)
	}

	if mErr.ErrorOrNil() != nil {
		return fmt.Errorf("error saving metadata of OBS bucket %s: %s", bucket, mErr)
	}

	return nil
}

func setObsBucketPolicy(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketPolicy(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.Code == "NoSuchBucketPolicy" {
				if err := d.Set("policy", nil); err != nil {
					return fmt.Errorf("error saving policy of OBS bucket %s: %s", bucket, err)
				}
				return nil
			}
			return fmt.Errorf("error getting policy of OBS bucket %s: %s", bucket, err)
		}
		return err
	}

	pol := output.Policy
	log.Printf("[DEBUG] getting policy of OBS bucket %s: %s", bucket, pol)
	if err := d.Set("policy", pol); err != nil {
		return fmt.Errorf("error saving policy of OBS bucket %s: %s", bucket, err)
	}

	return nil
}

func setObsBucketVersioning(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketVersioning(bucket)
	if err != nil {
		return getObsError("Error getting versioning status of OBS bucket", bucket, err)
	}

	log.Printf("[DEBUG] getting versioning status of OBS bucket %s: %s", bucket, output.Status)
	mErr := &multierror.Error{}
	if output.Status == obs.VersioningStatusEnabled {
		mErr = multierror.Append(mErr, d.Set("versioning", true))
	} else {
		mErr = multierror.Append(mErr, d.Set("versioning", false))
	}
	if mErr.ErrorOrNil() != nil {
		return fmt.Errorf("error saving version of OBS bucket %s: %s", bucket, mErr)
	}

	return nil
}

func setObsBucketEncryption(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketEncryption(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.Code == "NoSuchEncryptionConfiguration" || obsError.Code == "FsNotSupport" {
				mErr := multierror.Append(nil,
					d.Set("encryption", false),
					d.Set("kms_key_id", nil),
					d.Set("kms_key_project_id", nil),
					d.Set("sse_algorithm", nil),
				)
				if mErr.ErrorOrNil() != nil {
					return fmt.Errorf("error saving encryption of OBS bucket %s: %s", bucket, mErr)
				}
				return nil
			}
			return fmt.Errorf("error getting encryption configuration of OBS bucket %s: %s", bucket, err)
		}
		return err
	}

	log.Printf("[DEBUG] getting encryption configuration of OBS bucket %s: %+v", bucket, output.BucketEncryptionConfiguration)
	mErr := &multierror.Error{}
	if sseAlgorithm := output.SSEAlgorithm; sseAlgorithm != "" {
		mErr = multierror.Append(mErr,
			d.Set("encryption", true),
			d.Set("kms_key_id", output.KMSMasterKeyID),
			d.Set("kms_key_project_id", output.ProjectID),
			d.Set("sse_algorithm", sseAlgorithm),
		)
	} else {
		mErr = multierror.Append(mErr,
			d.Set("encryption", false),
			d.Set("kms_key_id", nil),
			d.Set("kms_key_project_id", nil),
			d.Set("sse_algorithm", nil),
		)
	}
	if mErr.ErrorOrNil() != nil {
		return fmt.Errorf("error saving encryption of OBS bucket %s: %s", bucket, mErr)
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
		if output.Agency != "" {
			logging["agency"] = output.Agency
		}
		lcList = append(lcList, logging)
	}
	log.Printf("[DEBUG] getting logging configuration of OBS bucket %s: %#v", bucket, lcList)

	if err := d.Set("logging", lcList); err != nil {
		return fmt.Errorf("error saving logging configuration of OBS bucket %s: %s", bucket, err)
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

	if err := d.Set("quota", output.Quota); err != nil {
		return fmt.Errorf("error saving quota of OBS bucket %s: %s", bucket, err)
	}

	return nil
}

func setObsBucketLifecycleConfiguration(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketLifecycleConfiguration(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.Code == "NoSuchLifecycleConfiguration" {
				if err := d.Set("lifecycle_rule", nil); err != nil {
					return fmt.Errorf("error saving lifecycle configuration of OBS bucket %s: %s", bucket, err)
				}
				return nil
			}
			return fmt.Errorf("error getting lifecycle configuration of OBS bucket %s: %s", bucket, err)
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

		// abort_incomplete_multipart_upload
		if days := lifecycleRule.AbortIncompleteMultipartUpload.DaysAfterInitiation; days > 0 {
			a := make(map[string]interface{})
			a["days"] = days
			rule["abort_incomplete_multipart_upload"] = schema.NewSet(expirationHash, []interface{}{a})
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
		return fmt.Errorf("error saving lifecycle configuration of OBS bucket %s: %s", bucket, err)
	}

	return nil
}

func setObsBucketWebsiteConfiguration(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketWebsiteConfiguration(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.Code == "NoSuchWebsiteConfiguration" {
				if err := d.Set("website", nil); err != nil {
					return fmt.Errorf("error saving website configuration of OBS bucket %s: %s", bucket, err)
				}
				return nil
			}
			return fmt.Errorf("error getting website configuration of OBS bucket %s: %s", bucket, err)
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
			return fmt.Errorf("error while marshaling website routing rules: %s", err)
		}
		w["routing_rules"] = rr
	}

	websites = append(websites, w)
	log.Printf("[DEBUG] saving website configuration of OBS bucket %s, website: %#v", bucket, websites)
	if err := d.Set("website", websites); err != nil {
		return fmt.Errorf("error saving website configuration of OBS bucket %s: %s", bucket, err)
	}
	return nil
}

func setObsBucketCorsRules(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketCors(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.Code == "NoSuchCORSConfiguration" {
				if err := d.Set("cors_rule", nil); err != nil {
					return fmt.Errorf("error saving CORS rules of OBS bucket %s: %s", bucket, err)
				}
				return nil
			}
			return fmt.Errorf("error getting CORS configuration of OBS bucket %s: %s", bucket, err)
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
		return fmt.Errorf("error saving CORS rules of OBS bucket %s: %s", bucket, err)
	}

	return nil
}

func setObsBucketTags(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketTagging(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.Code == "NoSuchTagSet" {
				if err := d.Set("tags", nil); err != nil {
					return fmt.Errorf("error saving tags of OBS bucket %s: %s", bucket, err)
				}
				return nil
			}
			return fmt.Errorf("error getting tags of OBS bucket %s: %s", bucket, err)
		}
		return err
	}

	tagMap := make(map[string]string)
	for _, tag := range output.Tags {
		tagMap[tag.Key] = tag.Value
	}
	log.Printf("[DEBUG] getting tags of OBS bucket %s: %#v", bucket, tagMap)
	if err := d.Set("tags", tagMap); err != nil {
		return fmt.Errorf("error saving tags of OBS bucket %s: %s", bucket, err)
	}
	return nil
}

func setObsBucketStorageInfo(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketStorageInfo(bucket)
	if err != nil {
		if _, ok := err.(obs.ObsError); ok {
			return fmt.Errorf("error getting storage info of OBS bucket %s: %s", bucket, err)
		}
		return err
	}
	log.Printf("[DEBUG] getting storage info of OBS bucket %s: %#v", bucket, output)

	storages := make([]map[string]interface{}, 1)
	storages[0] = map[string]interface{}{
		"size":          output.Size,
		"object_number": output.ObjectNumber,
	}

	if err := d.Set("storage_info", storages); err != nil {
		return fmt.Errorf("error saving storage info of OBS bucket %s: %s", bucket, err)
	}
	return nil
}

func setObsBucketUserDomainNames(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketCustomDomain(bucket)
	if err != nil {
		return getObsError("Error getting user domain names of OBS bucket", bucket, err)
	}
	log.Printf("[DEBUG] getting user domain names of OBS bucket %s: %#v", bucket, output)

	domainNames := make([]string, len(output.Domains))
	for i, v := range output.Domains {
		domainNames[i] = v.DomainName
	}
	return d.Set("user_domain_names", domainNames)
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
		return fmt.Errorf("error some objects are still exist in %s: %v", bucket, output.Errors)
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
	if _, ok := err.(obs.ObsError); ok {
		return fmt.Errorf("%s %s: %s", action, bucket, err)
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
	websiteRules := make([]WebsiteRoutingRule, 0, len(w))
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
		websiteRules = append(websiteRules, rule)
	}

	// normalize
	withNulls, err := json.Marshal(websiteRules)
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

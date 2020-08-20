package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/obs"
)

func resourceObsBucketPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceObsBucketPolicyPut,
		Read:   resourceObsBucketPolicyRead,
		Update: resourceObsBucketPolicyPut,
		Delete: resourceObsBucketPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"policy": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validateJsonString,
				DiffSuppressFunc: suppressEquivalentAwsPolicyDiffs,
			},
		},
	}
}

func resourceObsBucketPolicyPut(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	obsClient, err := config.newObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	policy := d.Get("policy").(string)
	log.Printf("[DEBUG] OBS bucket: %s, set policy: %s", bucket, policy)

	params := &obs.SetBucketPolicyInput{
		Bucket: bucket,
		Policy: policy,
	}
	if _, err := obsClient.SetBucketPolicy(params); err != nil {
		return getObsError("Error setting OBS bucket policy", bucket, err)
	}

	d.SetId(bucket)
	return nil
}

func resourceObsBucketPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	obsClient, err := config.newObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	log.Printf("[DEBUG] read policy for obs bucket: %s", d.Id())
	output, err := obsClient.GetBucketPolicy(d.Id())

	var pol string
	if err == nil && output.Policy != "" {
		pol = output.Policy
	}
	if err := d.Set("policy", pol); err != nil {
		return err
	}

	return nil
}

func resourceObsBucketPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	obsClient, err := config.newObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	bucket := d.Get("bucket").(string)

	log.Printf("[DEBUG] OBS bucket: %s, delete policy", bucket)
	_, err = obsClient.DeleteBucketPolicy(bucket)
	if err != nil {
		return getObsError("Error deleting policy of OBS bucket %s: %s", bucket, err)
	}

	return nil
}

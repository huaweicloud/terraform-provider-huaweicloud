package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceS3BucketPolicy() *schema.Resource {
	return &schema.Resource{
		Create:             resourceS3BucketPolicyPut,
		Read:               resourceS3BucketPolicyRead,
		Update:             resourceS3BucketPolicyPut,
		Delete:             resourceS3BucketPolicyDelete,
		DeprecationMessage: "use huaweicloud_obs_bucket_policy resource instead",

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

func resourceS3BucketPolicyPut(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	s3conn, err := config.computeS3conn(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud s3 client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	policy := d.Get("policy").(string)

	d.SetId(bucket)

	log.Printf("[DEBUG] S3 bucket: %s, put policy: %s", bucket, policy)

	params := &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket),
		Policy: aws.String(policy),
	}
	//lintignore:R006
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		if _, err := s3conn.PutBucketPolicy(params); err != nil {
			if awserr, ok := err.(awserr.Error); ok {
				if awserr.Code() == "MalformedPolicy" {
					return resource.RetryableError(awserr)
				}
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error putting S3 policy: %s", err)
	}

	return nil
}

func resourceS3BucketPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	s3conn, err := config.computeS3conn(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud s3 client: %s", err)
	}

	log.Printf("[DEBUG] S3 bucket policy, read for bucket: %s", d.Id())
	pol, err := s3conn.GetBucketPolicy(&s3.GetBucketPolicyInput{
		Bucket: aws.String(d.Id()),
	})

	v := ""
	if err == nil && pol.Policy != nil {
		v = *pol.Policy
	}
	if err := d.Set("policy", v); err != nil {
		return err
	}

	return nil
}

func resourceS3BucketPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	s3conn, err := config.computeS3conn(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud s3 client: %s", err)
	}

	bucket := d.Get("bucket").(string)

	log.Printf("[DEBUG] S3 bucket: %s, delete policy", bucket)
	_, err = s3conn.DeleteBucketPolicy(&s3.DeleteBucketPolicyInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "NoSuchBucket" {
			return nil
		}
		return fmt.Errorf("Error deleting S3 policy: %s", err)
	}

	return nil
}

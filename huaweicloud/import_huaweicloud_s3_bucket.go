package huaweicloud

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceS3BucketImportState(
	d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {

	results := make([]*schema.ResourceData, 1, 1)
	results[0] = d

	config := meta.(*Config)
	conn, err := config.computeS3conn(GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating HuaweiCloud s3 client: %s", err)
	}
	pol, err := conn.GetBucketPolicy(&s3.GetBucketPolicyInput{
		Bucket: aws.String(d.Id()),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "NoSuchBucketPolicy" {
			// Bucket without policy
			return results, nil
		}
		return nil, errwrap.Wrapf("Error importing AWS S3 bucket policy: {{err}}", err)
	}

	policy := resourceS3BucketPolicy()
	pData := policy.Data(nil)
	pData.SetId(d.Id())
	pData.SetType("huaweicloud_s3_bucket_policy")
	pData.Set("bucket", d.Id())
	pData.Set("policy", pol)
	results = append(results, pData)

	return results, nil
}

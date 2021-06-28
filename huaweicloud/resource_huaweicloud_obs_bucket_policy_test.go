package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/obs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccObsBucketPolicy_basic(t *testing.T) {
	name := fmtp.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	obsName := "huaweicloud_obs_bucket.bucket"
	policyName := "huaweicloud_obs_bucket_policy.policy"

	expectedPolicyText := fmtp.Sprintf(
		`{"Statement":[{"Sid":"test1","Effect":"Allow","Principal":{"ID":["*"]},"Action":["GetObject"],"Resource":["%s/*"]}]}`,
		name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckOBS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketPolicyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(obsName),
					testAccCheckObsBucketHasPolicy(policyName, expectedPolicyText),
					resource.TestCheckResourceAttr(policyName, "policy_format", "obs"),
				),
			},
			{
				ResourceName:      policyName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccObsBucketPolicy_update(t *testing.T) {
	name := fmtp.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	obsName := "huaweicloud_obs_bucket.bucket"
	policyName := "huaweicloud_obs_bucket_policy.policy"

	expectedPolicyText1 := fmtp.Sprintf(
		`{"Statement":[{"Sid":"test1","Effect":"Allow","Principal":{"ID":["*"]},"Action":["GetObject"],"Resource":["%s/*"]}]}`,
		name)

	expectedPolicyText2 := fmtp.Sprintf(
		`{"Statement":[{"Sid":"test2","Effect":"Allow","Principal":{"ID":["*"]},"Action":["GetObject","PutObject","DeleteObject"],"Resource":["%s/*"]}]}`,
		name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckOBS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketPolicyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(obsName),
					testAccCheckObsBucketHasPolicy(policyName, expectedPolicyText1),
					resource.TestCheckResourceAttr(policyName, "policy_format", "obs"),
				),
			},

			{
				Config: testAccObsBucketPolicyConfig_updated(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(obsName),
					testAccCheckObsBucketHasPolicy(policyName, expectedPolicyText2),
				),
			},
		},
	})
}

func TestAccObsBucketPolicy_s3(t *testing.T) {
	name := fmtp.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	obsName := "huaweicloud_obs_bucket.bucket"
	policyName := "huaweicloud_obs_bucket_policy.s3_policy"

	expectedPolicyText := fmtp.Sprintf(
		`{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:*"],"Resource":["arn:aws:s3:::%s","arn:aws:s3:::%s/*"]}]}`,
		name, name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckOBS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketPolicyS3Foramt(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(obsName),
					testAccCheckObsBucketHasPolicy(policyName, expectedPolicyText),
					resource.TestCheckResourceAttr(policyName, "policy_format", "s3"),
				),
			},
			{
				ResourceName:      policyName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOBSPolicyImportStateIDFunc(),
			},
		},
	})
}

func testAccCheckObsBucketHasPolicy(n string, expectedPolicyText string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No OBS Bucket ID is set")
		}

		var err error
		var obsClient *obs.ObsClient

		config := testAccProvider.Meta().(*config.Config)
		format := rs.Primary.Attributes["policy_format"]
		if format == "obs" {
			obsClient, err = config.ObjectStorageClientWithSignature(HW_REGION_NAME)
		} else {
			obsClient, err = config.ObjectStorageClient(HW_REGION_NAME)
		}
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud OBS client: %s", err)
		}

		policy, err := obsClient.GetBucketPolicy(rs.Primary.ID)
		if err != nil {
			return fmtp.Errorf("GetBucketPolicy error: %v", err)
		}

		actualPolicyText := policy.Policy
		if actualPolicyText != expectedPolicyText {
			return fmtp.Errorf("non-equivalent policy error:\n\nexpected: %s\n\n     got: %s",
				expectedPolicyText, actualPolicyText)
		}

		return nil
	}
}

func testAccOBSPolicyImportStateIDFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		policyRes, ok := s.RootModule().Resources["huaweicloud_obs_bucket_policy.s3_policy"]
		if !ok {
			return "", fmtp.Errorf("huaweicloud_obs_bucket_policy resource not found")
		}

		return fmtp.Sprintf("%s/s3", policyRes.Primary.ID), nil
	}
}

func testAccObsBucketPolicyConfig(bucketName string) string {
	return fmtp.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
	bucket = "%s"
	tags = {
	  TestName = "TestAccObsBucketPolicy_basic"
	}
}

resource "huaweicloud_obs_bucket_policy" "policy" {
	bucket = huaweicloud_obs_bucket.bucket.bucket
	policy =<<POLICY
{
	"Statement": [{
		"Sid": "test1",
		"Effect": "Allow",
		"Principal": {
			"ID": ["*"]
		},
		"Action": ["GetObject"],
		"Resource": ["%s/*"]
	}]
}
POLICY
}
`, bucketName, bucketName)
}

func testAccObsBucketPolicyConfig_updated(bucketName string) string {
	return fmtp.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
	bucket = "%s"
	tags = {
	  TestName = "TestAccObsBucketPolicy_updated"
	}
}

resource "huaweicloud_obs_bucket_policy" "policy" {
	bucket = huaweicloud_obs_bucket.bucket.bucket
	policy =<<POLICY
{
	"Statement": [{
		"Sid": "test2",
		"Effect": "Allow",
		"Principal": {
			"ID": ["*"]
		},
		"Action": ["GetObject", "PutObject", "DeleteObject"],
		"Resource": ["%s/*"]
	}]
}
POLICY
}
`, bucketName, bucketName)
}

func testAccObsBucketPolicyS3Foramt(bucketName string) string {
	return fmtp.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
	bucket = "%s"
	tags = {
	  TestName = "TestAccObsBucketPolicy_s3"
	}
}

resource "huaweicloud_obs_bucket_policy" "s3_policy" {
	bucket = huaweicloud_obs_bucket.bucket.bucket
	policy_format = "s3"
	policy =<<POLICY
{
	"Version": "2008-10-17",
	"Statement": [{
		"Effect": "Allow",
		"Principal": {
			"AWS": ["*"]
		},
		"Action": [
			"s3:*"
		],
		"Resource": [
			"arn:aws:s3:::%s",
			"arn:aws:s3:::%s/*"
		]
	}]
}
POLICY
}
`, bucketName, bucketName, bucketName)
}

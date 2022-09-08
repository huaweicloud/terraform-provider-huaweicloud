package obs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/obs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccObsBucketPolicy_basic(t *testing.T) {
	name := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	obsName := "huaweicloud_obs_bucket.bucket"
	policyName := "huaweicloud_obs_bucket_policy.policy"

	expectedPolicyText := fmt.Sprintf(
		`{"Statement":[{"Sid":"test1","Effect":"Allow","Principal":{"ID":["*"]},"Action":["GetObject"],"Resource":["%s/*"]}]}`,
		name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
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
	name := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	obsName := "huaweicloud_obs_bucket.bucket"
	policyName := "huaweicloud_obs_bucket_policy.policy"

	expectedPolicyText1 := fmt.Sprintf(
		`{"Statement":[{"Sid":"test1","Effect":"Allow","Principal":{"ID":["*"]},"Action":["GetObject"],"Resource":["%s/*"]}]}`,
		name)

	expectedPolicyText2 := fmt.Sprintf(
		`{"Statement":[{"Sid":"test2","Effect":"Allow","Principal":{"ID":["*"]},"Action":["GetObject","PutObject","DeleteObject"],"Resource":["%s/*"]}]}`,
		name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
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
	name := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	obsName := "huaweicloud_obs_bucket.bucket"
	policyName := "huaweicloud_obs_bucket_policy.s3_policy"

	expectedPolicyText := fmt.Sprintf(
		`{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:*"],"Resource":["arn:aws:s3:::%s","arn:aws:s3:::%s/*"]}]}`,
		name, name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
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
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OBS Bucket ID is set")
		}

		var err error
		var obsClient *obs.ObsClient

		conf := acceptance.TestAccProvider.Meta().(*config.Config)
		format := rs.Primary.Attributes["policy_format"]
		if format == "obs" {
			obsClient, err = conf.ObjectStorageClientWithSignature(acceptance.HW_REGION_NAME)
		} else {
			obsClient, err = conf.ObjectStorageClient(acceptance.HW_REGION_NAME)
		}
		if err != nil {
			return fmt.Errorf("Error creating OBS client: %s", err)
		}

		policy, err := obsClient.GetBucketPolicy(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("GetBucketPolicy error: %v", err)
		}

		actualPolicyText := policy.Policy
		if actualPolicyText != expectedPolicyText {
			return fmt.Errorf("non-equivalent policy error:\n\nexpected: %s\n\n     got: %s",
				expectedPolicyText, actualPolicyText)
		}

		return nil
	}
}

func testAccOBSPolicyImportStateIDFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		policyRes, ok := s.RootModule().Resources["huaweicloud_obs_bucket_policy.s3_policy"]
		if !ok {
			return "", fmt.Errorf("resource not found")
		}

		return fmt.Sprintf("%s/s3", policyRes.Primary.ID), nil
	}
}

func testAccObsBucketPolicyConfig(bucketName string) string {
	return fmt.Sprintf(`
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
	return fmt.Sprintf(`
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
	return fmt.Sprintf(`
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

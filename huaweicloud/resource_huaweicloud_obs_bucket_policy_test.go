package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccObsBucketPolicy_basic(t *testing.T) {
	name := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())

	expectedPolicyText := fmt.Sprintf(
		`{"Statement":[{"Sid":"test1","Effect":"Allow","Principal":{"ID":["*"]},"Action":["GetObject"],"Resource":["%s/*"]}]}`,
		name)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckS3(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketPolicyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists("huaweicloud_obs_bucket.bucket"),
					testAccCheckObsBucketHasPolicy("huaweicloud_obs_bucket.bucket", expectedPolicyText),
				),
			},
		},
	})
}

func TestAccObsBucketPolicy_update(t *testing.T) {
	name := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())

	expectedPolicyText1 := fmt.Sprintf(
		`{"Statement":[{"Sid":"test1","Effect":"Allow","Principal":{"ID":["*"]},"Action":["GetObject"],"Resource":["%s/*"]}]}`,
		name)

	expectedPolicyText2 := fmt.Sprintf(
		`{"Statement":[{"Sid":"test2","Effect":"Allow","Principal":{"ID":["*"]},"Action":["GetObject","PutObject","DeleteObject"],"Resource":["%s/*"]}]}`,
		name)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckS3(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketPolicyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists("huaweicloud_obs_bucket.bucket"),
					testAccCheckObsBucketHasPolicy("huaweicloud_obs_bucket.bucket", expectedPolicyText1),
				),
			},

			{
				Config: testAccObsBucketPolicyConfig_updated(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists("huaweicloud_obs_bucket.bucket"),
					testAccCheckObsBucketHasPolicy("huaweicloud_obs_bucket.bucket", expectedPolicyText2),
				),
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

		config := testAccProvider.Meta().(*Config)
		obsClient, err := config.newObjectStorageClient(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud OBS client: %s", err)
		}

		policy, err := obsClient.GetBucketPolicy(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("GetBucketPolicy error: %v", err)
		}

		actualPolicyText := policy.Policy
		if actualPolicyText != expectedPolicyText {
			return fmt.Errorf("Non-equivalent policy error:\n\nexpected: %s\n\n     got: %s\n",
				expectedPolicyText, actualPolicyText)
		}

		return nil
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

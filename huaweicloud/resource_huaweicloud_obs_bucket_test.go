package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccObsBucket_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckS3(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucket_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists("huaweicloud_obs_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket.bucket", "bucket", testAccObsBucketName(rInt)),
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket.bucket", "bucket_domain_name", testAccObsBucketDomainName(rInt)),
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket.bucket", "acl", "private"),
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket.bucket", "storage_class", "STANDARD"),
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket.bucket", "region", OS_REGION_NAME),
				),
			},
			{
				Config: testAccObsBucket_basic_update(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists("huaweicloud_obs_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket.bucket", "acl", "public-read"),
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket.bucket", "storage_class", "WARM"),
				),
			},
		},
	})
}

func TestAccObsBucket_tags(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckS3(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketConfigWithTags(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket.bucket", "tags.name", testAccObsBucketName(rInt)),
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket.bucket", "tags.foo", "bar"),
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket.bucket", "tags.key1", "value1"),
				),
			},
		},
	})
}

func TestAccObsBucket_versioning(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckS3(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketConfigWithVersioning(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists("huaweicloud_obs_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket.bucket", "versioning", "true"),
				),
			},
			{
				Config: testAccObsBucketConfigWithDisableVersioning(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists("huaweicloud_obs_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket.bucket", "versioning", "false"),
				),
			},
		},
	})
}

func testAccCheckObsBucketDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	obsClient, err := config.newObjectStorageClient(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_obs_bucket" {
			continue
		}

		_, err := obsClient.HeadBucket(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("HuaweiCloud OBS Bucket %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckObsBucketExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		obsClient, err := config.newObjectStorageClient(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud OBS client: %s", err)
		}

		_, err = obsClient.HeadBucket(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("HuaweiCloud OBS Bucket not found: %v", err)
		}
		return nil
	}
}

// These need a bit of randomness as the name can only be used once globally
func testAccObsBucketName(randInt int) string {
	return fmt.Sprintf("tf-test-bucket-%d", randInt)
}

func testAccObsBucketDomainName(randInt int) string {
	return fmt.Sprintf("tf-test-bucket-%d.obs.%s.myhuaweicloud.com", randInt, OS_REGION_NAME)
}

func testAccObsBucket_basic(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
    storage_class = "STANDARD"
	acl = "private"
}
`, randInt)
}

func testAccObsBucket_basic_update(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
    storage_class = "WARM"
	acl = "public-read"
}
`, randInt)
}

func testAccObsBucketConfigWithTags(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "private"

	tags = {
		name = "tf-test-bucket-%d"
        foo = "bar"
        key1 = "value1"
	}
}
`, randInt, randInt)
}

func testAccObsBucketConfigWithVersioning(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "private"
	versioning = true
}
`, randInt)
}

func testAccObsBucketConfigWithDisableVersioning(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "private"
	versioning = false
}
`, randInt)
}

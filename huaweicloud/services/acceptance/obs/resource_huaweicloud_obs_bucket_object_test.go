package obs

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/obs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccObsBucketObject_source(t *testing.T) {
	name := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_obs_bucket_object.test"

	tmpFile, err := os.CreateTemp("", "tf-acc-obs-obj-source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// write some data to the tempfile
	err = os.WriteFile(tmpFile.Name(), []byte("initial object state"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketObjectConfigSource(name, tmpFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketObjectExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key", name),
					resource.TestCheckResourceAttr(resourceName, "content_type", "binary/octet-stream"),
					resource.TestCheckResourceAttr(resourceName, "storage_class", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
			{
				// update with encryption
				Config: testAccObsBucketObjectConfig_withSSE(name, tmpFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "encryption", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccBucketObjectImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"encryption", "source",
				},
			},
		},
	})
}

func TestAccObsBucketObject_content(t *testing.T) {
	name := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_obs_bucket_object.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectConfigContent_step1(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketObjectExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key", name),
					resource.TestMatchResourceAttr(resourceName, "size", regexp.MustCompile("^[1-9][0-9]*$")),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccBucketObjectConfigContent_step2(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketObjectExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccBucketObjectImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"content",
				},
			},
		},
	})
}

func testAccCheckObsBucketObjectDestroy(s *terraform.State) error {
	conf := acceptance.TestAccProvider.Meta().(*config.Config)
	obsClient, err := conf.ObjectStorageClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OBS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_obs_bucket_object" {
			continue
		}

		bucket := rs.Primary.Attributes["bucket"]
		key := rs.Primary.Attributes["key"]
		input := &obs.ListObjectsInput{}
		input.Bucket = bucket
		input.Prefix = key

		resp, err := obsClient.ListObjects(input)
		if err != nil {
			if obsError, ok := err.(obs.ObsError); ok && obsError.Code == "NoSuchBucket" {
				return nil
			}
			return fmt.Errorf("Error listing objects of OBS bucket %s: %s", bucket, err)
		}

		var exist bool
		for _, content := range resp.Contents {
			if key == content.Key {
				exist = true
				break
			}
		}
		if exist {
			return fmt.Errorf("Resource %s still exists in bucket %s", rs.Primary.ID, bucket)
		}
	}

	return nil
}

func testAccCheckObsBucketObjectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OBS Bucket Object ID is set")
		}

		conf := acceptance.TestAccProvider.Meta().(*config.Config)
		obsClient, err := conf.ObjectStorageClient(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OBS client: %s", err)
		}

		bucket := rs.Primary.Attributes["bucket"]
		key := rs.Primary.Attributes["key"]
		input := &obs.ListObjectsInput{}
		input.Bucket = bucket
		input.Prefix = key

		resp, err := obsClient.ListObjects(input)
		if err != nil {
			return fmt.Errorf("Error listing objects of OBS bucket %s: %s", bucket, err)
		}

		var exist bool
		for _, content := range resp.Contents {
			if key == content.Key {
				exist = true
				break
			}
		}
		if !exist {
			return fmt.Errorf("Resource %s not found in bucket %s", rs.Primary.ID, bucket)
		}

		return nil
	}
}

func testAccBucketObjectImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		bucket, ok := s.RootModule().Resources["huaweicloud_obs_bucket.test"]
		if !ok {
			return "", fmt.Errorf("Bucket not found: %s", bucket)
		}
		object, ok := s.RootModule().Resources["huaweicloud_obs_bucket_object.test"]
		if !ok {
			return "", fmt.Errorf("Object not found: %s", object)
		}
		if bucket.Primary.ID == "" || object.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", bucket.Primary.ID, object.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", bucket.Primary.ID, object.Primary.ID), nil
	}
}

func testAccBucketObject_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s"
  acl           = "private"
  force_destroy = true
}
`, name)
}

func testAccObsBucketObjectConfigSource(name, source string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket_object" "test" {
  bucket       = huaweicloud_obs_bucket.test.bucket
  key          = "%[2]s"
  source       = "%[3]s"
  content_type = "binary/octet-stream"

  tags = {
    foo = "bar"
  }
}
`, testAccBucketObject_base(name), name, source)
}

func testAccObsBucketObjectConfig_withSSE(name, source string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket_object" "test" {
  bucket       = huaweicloud_obs_bucket.test.bucket
  key          = "%[2]s"
  source       = "%[3]s"
  content_type = "binary/octet-stream"
  encryption   = true

  tags = {
    owner = "terraform"
  }
}
`, testAccBucketObject_base(name), name, source)
}

func testAccBucketObjectConfigContent_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket_object" "test" {
  bucket  = huaweicloud_obs_bucket.test.bucket
  key     = "%[2]s"
  content = "some_bucket_content"

  tags = {
    foo = "bar"
  }
}
`, testAccBucketObject_base(name), name)
}

func testAccBucketObjectConfigContent_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket_object" "test" {
  bucket  = huaweicloud_obs_bucket.test.bucket
  key     = "%[2]s"
  content = "update_some_bucket_content"

  tags = {}
}
`, testAccBucketObject_base(name), name)
}

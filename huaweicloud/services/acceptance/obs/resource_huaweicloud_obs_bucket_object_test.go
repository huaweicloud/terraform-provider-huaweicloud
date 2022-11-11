package obs

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/obs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccObsBucketObject_source(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "huaweicloud_obs_bucket_object.object"

	tmpFile, err := os.CreateTemp("", "tf-acc-obs-obj-source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// write some data to the tempfile
	err = os.WriteFile(tmpFile.Name(), []byte("initial object state"), 0644)
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
				Config: testAccObsBucketObjectConfigSource(rInt, tmpFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketObjectExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key", "test-key"),
					resource.TestCheckResourceAttr(resourceName, "content_type", "binary/octet-stream"),
					resource.TestCheckResourceAttr(resourceName, "storage_class", "STANDARD"),
				),
			},
			{
				// update with encryption
				Config: testAccObsBucketObjectConfig_withSSE(rInt, tmpFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "encryption", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccObsBucketObjecImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"encryption", "source",
				},
			},
		},
	})
}

func TestAccObsBucketObject_content(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "huaweicloud_obs_bucket_object.object"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccObsBucketObjectConfigContent(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketObjectExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key", "test-key"),
					resource.TestCheckResourceAttr(resourceName, "size", "19"),
				),
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

func testAccObsBucketObjecImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		bucket, ok := s.RootModule().Resources["huaweicloud_obs_bucket.object_bucket"]
		if !ok {
			return "", fmt.Errorf("Bucket not found: %s", bucket)
		}
		object, ok := s.RootModule().Resources["huaweicloud_obs_bucket_object.object"]
		if !ok {
			return "", fmt.Errorf("Object not found: %s", object)
		}
		if bucket.Primary.ID == "" || object.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", bucket.Primary.ID, object.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", bucket.Primary.ID, object.Primary.ID), nil
	}
}

func testAccObsBucketObjectConfigSource(randInt int, source string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "object_bucket" {
  bucket = "tf-acc-test-bucket-%d"
}

resource "huaweicloud_obs_bucket_object" "object" {
  bucket       = huaweicloud_obs_bucket.object_bucket.bucket
  key          = "test-key"
  source       = "%s"
  content_type = "binary/octet-stream"
}
`, randInt, source)
}

func testAccObsBucketObjectConfigContent(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "object_bucket" {
  bucket = "tf-acc-test-bucket-%d"
}

resource "huaweicloud_obs_bucket_object" "object" {
  bucket  = huaweicloud_obs_bucket.object_bucket.bucket
  key     = "test-key"
  content = "some_bucket_content"
}
`, randInt)
}

func testAccObsBucketObjectConfig_withSSE(randInt int, source string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "object_bucket" {
  bucket = "tf-acc-test-bucket-%d"
}

resource "huaweicloud_obs_bucket_object" "object" {
bucket       = huaweicloud_obs_bucket.object_bucket.bucket
key          = "test-key"
source       = "%s"
content_type = "binary/octet-stream"
encryption   = true
}
`, randInt, source)
}

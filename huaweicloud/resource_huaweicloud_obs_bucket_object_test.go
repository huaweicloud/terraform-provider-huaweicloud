package huaweicloud

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/obs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccObsBucketObject_source(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "tf-acc-obs-obj-source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	rInt := acctest.RandInt()
	// write some data to the tempfile
	err = ioutil.WriteFile(tmpFile.Name(), []byte("initial object state"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckOBS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckObsBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketObjectConfigSource(rInt, tmpFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketObjectExists("huaweicloud_obs_bucket_object.object"),
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket_object.object", "key", "test-key"),
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket_object.object", "content_type", "binary/octet-stream"),
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket_object.object", "storage_class", "STANDARD"),
				),
			},
			{
				// update with encryption
				Config: testAccObsBucketObjectConfig_withSSE(rInt, tmpFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket_object.object", "encryption", "true"),
				),
			},
		},
	})
}

func TestAccObsBucketObject_content(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckOBS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckObsBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccObsBucketObjectConfigContent(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketObjectExists("huaweicloud_obs_bucket_object.object"),
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket_object.object", "key", "test-key"),
					resource.TestCheckResourceAttr(
						"huaweicloud_obs_bucket_object.object", "size", "19"),
				),
			},
		},
	})
}

func testAccCheckObsBucketObjectDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	obsClient, err := config.ObjectStorageClient(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud OBS client: %s", err)
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

		config := testAccProvider.Meta().(*config.Config)
		obsClient, err := config.ObjectStorageClient(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud OBS client: %s", err)
		}

		bucket := rs.Primary.Attributes["bucket"]
		key := rs.Primary.Attributes["key"]
		input := &obs.ListObjectsInput{}
		input.Bucket = bucket
		input.Prefix = key

		resp, err := obsClient.ListObjects(input)
		if err != nil {
			return getObsError("Error listing objects of OBS bucket", bucket, err)
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

func testAccObsBucketObjectConfigSource(randInt int, source string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "object_bucket" {
    bucket = "tf-object-test-bucket-%d"
}
resource "huaweicloud_obs_bucket_object" "object" {
	bucket = huaweicloud_obs_bucket.object_bucket.bucket
	key = "test-key"
	source = "%s"
	content_type = "binary/octet-stream"
}
`, randInt, source)
}

func testAccObsBucketObjectConfigContent(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "object_bucket" {
        bucket = "tf-object-test-bucket-%d"
}
resource "huaweicloud_obs_bucket_object" "object" {
        bucket = huaweicloud_obs_bucket.object_bucket.bucket
        key = "test-key"
        content = "some_bucket_content"
}
`, randInt)
}

func testAccObsBucketObjectConfig_withSSE(randInt int, source string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "object_bucket" {
	bucket = "tf-object-test-bucket-%d"
}

resource "huaweicloud_obs_bucket_object" "object" {
	bucket = huaweicloud_obs_bucket.object_bucket.bucket
	key = "test-key"
	source = "%s"
	content_type = "binary/octet-stream"
	encryption = true
}
`, randInt, source)
}

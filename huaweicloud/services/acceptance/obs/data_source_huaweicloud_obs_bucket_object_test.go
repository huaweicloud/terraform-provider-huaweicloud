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

func TestAccObsBucketObjectDataSource_content(t *testing.T) {
	rInt := acctest.RandInt()
	dataSourceName := "data.huaweicloud_obs_bucket_object.obj"
	resourceConf, dataSourceConf := testAccObsBucketObjectDataSource_content(rInt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketObjectExists("huaweicloud_obs_bucket_object.object"),
				),
			},
			{
				Config: dataSourceConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsObjectDataSourceExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "content_type", "binary/octet-stream"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_class", "STANDARD"),
				),
			},
		},
	})
}

func TestAccObsBucketObjectDataSource_source(t *testing.T) {
	dataSourceName := "data.huaweicloud_obs_bucket_object.obj"
	tmpFile, err := os.CreateTemp("", "tf-acc-obs-obj-source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// write test data to the tempfile
	for i := 0; i < 1024; i++ {
		_, err := tmpFile.WriteString("test obs object file storage\n")
		if err != nil {
			t.Fatal(err)
		}
	}
	tmpFile.Close()

	rInt := acctest.RandInt()
	resourceConf, dataSourceConf := testAccObsBucketObjectDataSource_source(rInt, tmpFile.Name())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketObjectExists("huaweicloud_obs_bucket_object.object"),
				),
			},
			{
				Config: dataSourceConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsObjectDataSourceExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "content_type", "binary/octet-stream"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_class", "STANDARD"),
				),
			},
		},
	})
}

func TestAccObsBucketObjectDataSource_allParams(t *testing.T) {
	rInt := acctest.RandInt()
	dataSourceName := "data.huaweicloud_obs_bucket_object.obj"
	resourceConf, dataSourceConf := testAccObsBucketObjectDataSource_allParams(rInt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketObjectExists("huaweicloud_obs_bucket_object.object"),
				),
			},
			{
				Config: dataSourceConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsObjectDataSourceExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "content_type", "application/json"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_class", "STANDARD"),
					resource.TestCheckResourceAttrSet(dataSourceName, "body"),
				),
			},
		},
	})
}

func testAccCheckObsObjectDataSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find OBS object data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("OBS object data source ID not set")
		}

		bucket := rs.Primary.Attributes["bucket"]
		key := rs.Primary.Attributes["key"]

		conf := acceptance.TestAccProvider.Meta().(*config.Config)
		obsClient, err := conf.ObjectStorageClient(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OBS client: %s", err)
		}

		respList, err := obsClient.ListObjects(&obs.ListObjectsInput{
			Bucket: bucket,
			ListObjsInput: obs.ListObjsInput{
				Prefix: key,
			},
		})
		if err != nil {
			return fmt.Errorf("Error listing objects of OBS bucket %s: %s", bucket, err)
		}

		var exist bool
		for _, content := range respList.Contents {
			if key == content.Key {
				exist = true
				break
			}
		}
		if !exist {
			return fmt.Errorf("object %s not found in bucket %s", key, bucket)
		}

		return nil
	}
}

func testAccObsBucketObjectDataSource_content(randInt int) (string, string) {
	resource := fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "object_bucket" {
  bucket = "tf-acc-test-bucket-%d"
}

resource "huaweicloud_obs_bucket_object" "object" {
  bucket  = huaweicloud_obs_bucket.object_bucket.bucket
  key     = "test-key-%d"
  content = "some_bucket_content"
}
`, randInt, randInt)

	dataSource := fmt.Sprintf(`
%s

data "huaweicloud_obs_bucket_object" "obj" {
  bucket = "tf-acc-test-bucket-%d"
  key    = "test-key-%d"
}`, resource, randInt, randInt)

	return resource, dataSource
}

func testAccObsBucketObjectDataSource_source(randInt int, source string) (string, string) {
	resource := fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "object_bucket" {
  bucket = "tf-acc-test-bucket-%d"
}

resource "huaweicloud_obs_bucket_object" "object" {
  bucket       = huaweicloud_obs_bucket.object_bucket.bucket
  key          = "test-key-%d"
  source       = "%s"
  content_type = "binary/octet-stream"
}
`, randInt, randInt, source)

	dataSource := fmt.Sprintf(`
%s

data "huaweicloud_obs_bucket_object" "obj" {
  bucket = "tf-acc-test-bucket-%d"
  key    = "test-key-%d"
}`, resource, randInt, randInt)

	return resource, dataSource
}

func testAccObsBucketObjectDataSource_allParams(randInt int) (string, string) {
	resource := fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "object_bucket" {
  bucket = "tf-acc-test-bucket-%d"
}

resource "huaweicloud_obs_bucket_object" "object" {
  bucket        = huaweicloud_obs_bucket.object_bucket.bucket
  key           = "test-key-%d"
  acl           = "private"
  storage_class = "STANDARD"
  encryption    = true
  content_type  = "application/json"
  content       = <<CONTENT
    {"msg": "Hi there!"}
CONTENT
}
`, randInt, randInt)

	dataSource := fmt.Sprintf(`
%s

data "huaweicloud_obs_bucket_object" "obj" {
  bucket = "tf-acc-test-bucket-%d"
  key    = "test-key-%d"
}`, resource, randInt, randInt)

	return resource, dataSource
}

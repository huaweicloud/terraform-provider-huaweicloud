package obs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/obs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getOBSBucketObjectAclResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	obsClient, err := cfg.ObjectStorageClient(region)
	if err != nil {
		return nil, fmt.Errorf("error creating OBS Client: %s", err)
	}
	params := &obs.GetObjectAclInput{
		Bucket: state.Primary.Attributes["bucket"],
		Key:    state.Primary.ID,
	}
	return obsClient.GetObjectAcl(params)
}

func TestAccOBSBucketObjectAcl_basic(t *testing.T) {
	var obj interface{}

	bucketName := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_obs_bucket_object_acl.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOBSBucketObjectAclResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testOBSBucketObjectAcl_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "bucket", bucketName),
					resource.TestCheckResourceAttr(rName, "key", "test-key"),
					resource.TestCheckResourceAttr(rName, "public_permission.0.access_to_acl.0", "READ_ACP"),
					resource.TestCheckResourceAttr(rName, "public_permission.0.access_to_acl.1", "WRITE_ACP"),
					resource.TestCheckResourceAttr(rName, "account_permission.#", "2"),
					resource.TestCheckResourceAttr(rName, "public_permission.#", "1"),
					resource.TestCheckResourceAttr(rName, "owner_permission.#", "1"),
				),
			},
			{
				Config: testOBSBucketObjectAcl_basicUpdate(bucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "bucket", bucketName),
					resource.TestCheckResourceAttr(rName, "key", "test-key"),
					resource.TestCheckResourceAttr(rName, "account_permission.0.access_to_acl.0", "READ_ACP"),
					resource.TestCheckResourceAttr(rName, "account_permission.0.account_id", "1000010022"),
					resource.TestCheckResourceAttr(rName, "public_permission.0.access_to_acl.0", "WRITE_ACP"),
					resource.TestCheckResourceAttr(rName, "account_permission.#", "1"),
					resource.TestCheckResourceAttr(rName, "public_permission.#", "1"),
					resource.TestCheckResourceAttr(rName, "owner_permission.#", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testObsBucketObjectAclImportState(rName),
			},
		},
	})
}

func TestAccOBSBucketObjectAcl_checkError(t *testing.T) {
	var obj interface{}

	bucketName := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_obs_bucket_object_acl.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOBSBucketObjectAclResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// the object owner account id is required for this test check
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testOBSBucketObjectAcl_basicErrorCheck(bucketName),
				ExpectError: regexp.MustCompile(`the account id cannot be the object owner`),
			},
		},
	})
}

func testOBSBucketObjectAcl_base(bucketName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "%s"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "huaweicloud_obs_bucket_object" "object" {
  bucket  = huaweicloud_obs_bucket.bucket.bucket
  key     = "test-key"
  content = "some_bucket_content"
}

`, bucketName)
}

func testOBSBucketObjectAcl_basic(bucketName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_obs_bucket_object_acl" "test" {
  bucket = huaweicloud_obs_bucket.bucket.bucket
  key    = huaweicloud_obs_bucket_object.object.key

  account_permission {
    access_to_object = ["READ"]
    access_to_acl    = ["READ_ACP", "WRITE_ACP"]
    account_id       = "1000010020"
  }

  account_permission {
    access_to_object = ["READ"]
    access_to_acl    = ["READ_ACP"]
    account_id       = "1000010021"
  }

  public_permission {
    access_to_acl = ["READ_ACP", "WRITE_ACP"]
  }
}
`, testOBSBucketObjectAcl_base(bucketName))
}

func testOBSBucketObjectAcl_basicUpdate(bucketName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_obs_bucket_object_acl" "test" {
  bucket = huaweicloud_obs_bucket.bucket.bucket
  key    = huaweicloud_obs_bucket_object.object.key

  account_permission {
    access_to_acl    = ["READ_ACP"]
    account_id       = "1000010022"
  }

  public_permission {
    access_to_acl = ["WRITE_ACP"]
  }
}
`, testOBSBucketObjectAcl_base(bucketName))
}

func testOBSBucketObjectAcl_basicErrorCheck(bucketName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_obs_bucket_object_acl" "test" {
  bucket = huaweicloud_obs_bucket.bucket.bucket
  key    = huaweicloud_obs_bucket_object.object.key

  account_permission {
    access_to_acl    = ["READ_ACP"]
    account_id       = "%s"
  }

  public_permission {
    access_to_acl = ["WRITE_ACP"]
  }
}
`, testOBSBucketObjectAcl_base(bucketName), acceptance.HW_DOMAIN_ID)
}

func testObsBucketObjectAclImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		bucket := rs.Primary.Attributes["bucket"]
		if bucket == "" {
			return "", fmt.Errorf("attribute (bucket) of Resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s/%s", bucket, rs.Primary.ID), nil
	}
}

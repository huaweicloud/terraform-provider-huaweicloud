package obs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getOBSBucketAclResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	obsClient, err := cfg.ObjectStorageClient(region)
	if err != nil {
		return nil, fmt.Errorf("error creating OBS Client: %s", err)
	}

	output, err := obsClient.GetBucketAcl(state.Primary.ID)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func TestAccOBSBucketAcl_basic(t *testing.T) {
	var obj interface{}

	bucketName := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_obs_bucket_acl.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOBSBucketAclResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testOBSBucketAcl_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "bucket", bucketName),
					resource.TestCheckResourceAttr(rName, "log_delivery_user_permission.0.access_to_bucket.0", "READ"),
					resource.TestCheckResourceAttr(rName, "log_delivery_user_permission.0.access_to_bucket.1", "WRITE"),
					resource.TestCheckResourceAttr(rName, "log_delivery_user_permission.0.access_to_acl.0", "READ_ACP"),
					resource.TestCheckResourceAttr(rName, "log_delivery_user_permission.0.access_to_acl.1", "WRITE_ACP"),
					resource.TestCheckResourceAttr(rName, "account_permission.#", "2"),
					resource.TestCheckResourceAttr(rName, "owner_permission.#", "1"),
				),
			},
			{
				Config: testOBSBucketAcl_basic_update(bucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "owner_permission.0.access_to_bucket.0", "WRITE"),
					resource.TestCheckResourceAttr(rName, "owner_permission.0.access_to_acl.0", "WRITE_ACP"),
					resource.TestCheckResourceAttr(rName, "account_permission.0.access_to_acl.0", "READ_ACP"),
					resource.TestCheckResourceAttr(rName, "account_permission.0.account_id", "1000010023"),
					resource.TestCheckResourceAttr(rName, "public_permission.0.access_to_bucket.0", "READ"),
					resource.TestCheckResourceAttr(rName, "public_permission.0.access_to_bucket.1", "WRITE"),
					resource.TestCheckResourceAttr(rName, "owner_permission.#", "1"),
					resource.TestCheckResourceAttr(rName, "public_permission.#", "1"),
					resource.TestCheckResourceAttr(rName, "account_permission.#", "1"),
					resource.TestCheckResourceAttr(rName, "log_delivery_user_permission.#", "0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testOBSBucketAcl_base(bucketName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "%s"
  storage_class = "STANDARD"
  acl           = "private"
}
`, bucketName)
}

func testOBSBucketAcl_basic(bucketName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_obs_bucket_acl" "test" {
  bucket = huaweicloud_obs_bucket.bucket.bucket

  account_permission {
    access_to_bucket = ["READ", "WRITE"]
    access_to_acl    = ["READ_ACP", "WRITE_ACP"]
    account_id       = "1000010020"
  }

  account_permission {
    access_to_bucket = ["READ"]
    access_to_acl    = ["READ_ACP", "WRITE_ACP"]
    account_id       = "1000010021"
  }

  log_delivery_user_permission {
    access_to_bucket = ["READ", "WRITE"]
    access_to_acl    = ["READ_ACP", "WRITE_ACP"]
  }
}
`, testOBSBucketAcl_base(bucketName))
}

func testOBSBucketAcl_basic_update(bucketName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_obs_bucket_acl" "test" {
  bucket = huaweicloud_obs_bucket.bucket.bucket

  owner_permission {
    access_to_bucket = ["WRITE"]
    access_to_acl    = ["WRITE_ACP"]
  }

  account_permission {
    access_to_acl    = ["READ_ACP"]
    account_id       = "1000010023"
  }

  public_permission {
    access_to_bucket = ["READ", "WRITE"]
  }
}
`, testOBSBucketAcl_base(bucketName))
}

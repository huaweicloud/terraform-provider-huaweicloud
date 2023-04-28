package obs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getOBSBucketReplicationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	obsClient, err := cfg.ObjectStorageClientWithSignature(region)
	if err != nil {
		return nil, fmt.Errorf("error creating OBS Client: %s", err)
	}

	output, err := obsClient.GetBucketReplication(state.Primary.ID)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func TestAccObsBucketReplication_basic(t *testing.T) {
	var obj interface{}

	bucketName := acceptance.RandomAccResourceNameWithDash()
	agencyName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_obs_bucket_replication.replica"
	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOBSBucketReplicationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// The source bucket and target bucket must belong to different regions of the same account.
			// https://support.huaweicloud.com/intl/en-us/ugobs-obs/obs_41_0034.html
			acceptance.TestAccPreCheckOBSDestinationBucket(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketReplication_basic(agencyName, bucketName, acceptance.HW_OBS_DESTINATION_BUCKET),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "bucket", bucketName),
					resource.TestCheckResourceAttr(rName, "agency", agencyName),
					resource.TestCheckResourceAttr(rName, "destination_bucket", acceptance.HW_OBS_DESTINATION_BUCKET),
					resource.TestCheckResourceAttr(rName, "rule.#", "1"),
					resource.TestCheckResourceAttr(rName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(rName, "rule.0.prefix", "abc"),
					resource.TestCheckResourceAttrSet(rName, "rule.0.id"),
				),
			},
			{
				Config: testAccObsBucketReplication_update_1(agencyName, bucketName, acceptance.HW_OBS_DESTINATION_BUCKET),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "rule.#", "2"),
					resource.TestCheckResourceAttr(rName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(rName, "rule.0.prefix", "imgs/"),
					resource.TestCheckResourceAttr(rName, "rule.1.enabled", "false"),
					resource.TestCheckResourceAttr(rName, "rule.1.prefix", "terraform"),
					resource.TestCheckResourceAttr(rName, "rule.1.storage_class", "COLD"),
					resource.TestCheckResourceAttr(rName, "rule.1.history_enabled", "true"),
					resource.TestCheckResourceAttrSet(rName, "rule.0.id"),
					resource.TestCheckResourceAttrSet(rName, "rule.1.id"),
				),
			},
			{
				Config: testAccObsBucketReplication_update_2(agencyName, bucketName, acceptance.HW_OBS_DESTINATION_BUCKET),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "rule.#", "1"),
					resource.TestCheckResourceAttr(rName, "rule.0.prefix", ""),
					resource.TestCheckResourceAttr(rName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(rName, "rule.0.id"),
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

func testAccObsBucketReplication_base(agencyName, bucketName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_agency" "agency" {
  name                   = "%s"
  description            = "This is an iam agency for obs bucket replication"
  delegated_service_name = "op_svc_obs"

  domain_roles = [
    "OBS Administrator",
  ]
}

resource "huaweicloud_obs_bucket" "source" {
  bucket        = "%s"
  storage_class = "STANDARD"
  acl           = "private"
}
`, agencyName, bucketName)
}

func testAccObsBucketReplication_basic(agencyName, bucketName, destinationName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_obs_bucket_replication" "replica" {
  bucket             = huaweicloud_obs_bucket.source.bucket
  destination_bucket = "%s"
  agency             = huaweicloud_identity_agency.agency.name

  rule {
    prefix = "abc"
  }
}
`, testAccObsBucketReplication_base(agencyName, bucketName), destinationName)
}

func testAccObsBucketReplication_update_1(agencyName, bucketName, destinationName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_obs_bucket_replication" "replica" {
  bucket             = huaweicloud_obs_bucket.source.bucket
  destination_bucket = "%s"
  agency             = huaweicloud_identity_agency.agency.name

  rule {
    prefix = "imgs/"
  }
  rule {
    enabled         = false
    prefix          = "terraform"
    storage_class   = "COLD"
    history_enabled = true
  }
}
`, testAccObsBucketReplication_base(agencyName, bucketName), destinationName)
}

func testAccObsBucketReplication_update_2(agencyName, bucketName, destinationName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_obs_bucket_replication" "replica" {
  bucket             = huaweicloud_obs_bucket.source.bucket
  destination_bucket = "%s"
  agency             = huaweicloud_identity_agency.agency.name

  rule {}
}
`, testAccObsBucketReplication_base(agencyName, bucketName), destinationName)
}

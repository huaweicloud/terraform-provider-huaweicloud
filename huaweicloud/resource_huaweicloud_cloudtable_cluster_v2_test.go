// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file at
//     https://www.github.com/huaweicloud/magic-modules
//
// ----------------------------------------------------------------------------

package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCloudtableClusterV2_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudtableClusterV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudtableClusterV2_basic(acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudtableClusterV2Exists(),
				),
			},
		},
	})
}

func testAccCloudtableClusterV2_basic(val string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup_v2" "secgroup" {
  name = "terraform_test_security_group%s"
  description = "terraform security group acceptance test"
  timeouts {
    delete = "20m"
  }
}

resource "huaweicloud_cloudtable_cluster_v2" "cluster" {
  availability_zone = "%s"
  name = "terraform-test-cluster%s"
  rs_num = 2
  security_group_id = "${huaweicloud_networking_secgroup_v2.secgroup.id}"
  subnet_id = "%s"
  vpc_id = "%s"
  storage_type = "COMMON"
}
	`, val, HW_AVAILABILITY_ZONE, val, HW_NETWORK_ID, HW_VPC_ID)
}

func testAccCheckCloudtableClusterV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.CloudtableV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cloudtable_cluster_v2" {
			continue
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return err
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
				"X-Language":   "en-us",
			}})
		if err == nil {
			return fmtp.Errorf("huaweicloud_cloudtable_cluster_v2 still exists at %s", url)
		}
	}

	return nil
}

func testAccCheckCloudtableClusterV2Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.CloudtableV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["huaweicloud_cloudtable_cluster_v2.cluster"]
		if !ok {
			return fmtp.Errorf("Error checking huaweicloud_cloudtable_cluster_v2.cluster exist, err=not found this resource")
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return fmtp.Errorf("Error checking huaweicloud_cloudtable_cluster_v2.cluster exist, err=building url failed: %s", err)
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
				"X-Language":   "en-us",
			}})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmtp.Errorf("huaweicloud_cloudtable_cluster_v2.cluster is not exist")
			}
			return fmtp.Errorf("Error checking huaweicloud_cloudtable_cluster_v2.cluster exist, err=send request failed: %s", err)
		}
		return nil
	}
}

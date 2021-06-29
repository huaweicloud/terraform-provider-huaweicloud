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

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccGesGraphV1_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGesGraphV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGesGraphV1_basic(acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGesGraphV1Exists(),
				),
			},
		},
	})
}

func testAccGesGraphV1_basic(val string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup_v2" "secgroup" {
  name = "terraform_test_security_group%s"
  description = "terraform security group acceptance test"
}

resource "huaweicloud_ges_graph_v1" "graph" {
  availability_zone = "%s"
  graph_size_type = 0
  name = "terraform_ges_graph_test%s"
  region = "%s"
  security_group_id = "${huaweicloud_networking_secgroup_v2.secgroup.id}"
  subnet_id = "%s"
  vpc_id = "%s"
}
	`, val, HW_AVAILABILITY_ZONE, val, HW_REGION_NAME, HW_NETWORK_ID, HW_VPC_ID)
}

func testAccCheckGesGraphV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.GesV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_ges_graph_v1" {
			continue
		}

		url, err := replaceVarsForTest(rs, "graphs/{id}")
		if err != nil {
			return err
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err == nil {
			return fmtp.Errorf("huaweicloud_ges_graph_v1 still exists at %s", url)
		}
	}

	return nil
}

func testAccCheckGesGraphV1Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.GesV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["huaweicloud_ges_graph_v1.graph"]
		if !ok {
			return fmtp.Errorf("Error checking huaweicloud_ges_graph_v1.graph exist, err=not found this resource")
		}

		url, err := replaceVarsForTest(rs, "graphs/{id}")
		if err != nil {
			return fmtp.Errorf("Error checking huaweicloud_ges_graph_v1.graph exist, err=building url failed: %s", err)
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmtp.Errorf("huaweicloud_ges_graph_v1.graph is not exist")
			}
			return fmtp.Errorf("Error checking huaweicloud_ges_graph_v1.graph exist, err=send request failed: %s", err)
		}
		return nil
	}
}

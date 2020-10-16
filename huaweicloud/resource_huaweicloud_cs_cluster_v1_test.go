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

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
)

func TestAccCsClusterV1_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCsClusterV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCsClusterV1_basic(acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCsClusterV1Exists(),
				),
			},
		},
	})
}

func testAccCsClusterV1_basic(val string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cs_cluster_v1" "cluster" {
  name = "terraform_cs_cluster_v1_test%s"
}
	`, val)
}

func testAccCheckCsClusterV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.cloudStreamV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cs_cluster_v1" {
			continue
		}

		url, err := replaceVarsForTest(rs, "reserved_cluster/{id}")
		if err != nil {
			return err
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err == nil {
			return fmt.Errorf("huaweicloud_cs_cluster_v1 still exists at %s", url)
		}
	}

	return nil
}

func testAccCheckCsClusterV1Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		client, err := config.cloudStreamV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["huaweicloud_cs_cluster_v1.cluster"]
		if !ok {
			return fmt.Errorf("Error checking huaweicloud_cs_cluster_v1.cluster exist, err=not found this resource")
		}

		url, err := replaceVarsForTest(rs, "reserved_cluster/{id}")
		if err != nil {
			return fmt.Errorf("Error checking huaweicloud_cs_cluster_v1.cluster exist, err=building url failed: %s", err)
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmt.Errorf("huaweicloud_cs_cluster_v1.cluster is not exist")
			}
			return fmt.Errorf("Error checking huaweicloud_cs_cluster_v1.cluster exist, err=send request failed: %s", err)
		}
		return nil
	}
}

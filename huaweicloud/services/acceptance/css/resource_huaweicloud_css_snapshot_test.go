package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/css/v1/snapshots"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCssSnapshot_basic(t *testing.T) {
	rand := acctest.RandString(5)
	resourceName := "huaweicloud_css_snapshot.snapshot"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCssSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCssSnapshot_basic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssSnapshotExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("snapshot-%s", rand)),
					resource.TestCheckResourceAttr(resourceName, "status", "COMPLETED"),
					resource.TestCheckResourceAttr(resourceName, "backup_type", "manual"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("not found: %s", resourceName)
					}
					return fmt.Sprintf("%s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.ID), nil
				},
			},
		},
	})
}

func testAccCheckCssSnapshotDestroy(s *terraform.State) error {
	conf := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := conf.CssV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating css client, err: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_css_snapshot" {
			continue
		}

		clusterID := rs.Primary.Attributes["cluster_id"]
		snapList, err := snapshots.List(client, clusterID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}

			return err
		}

		for _, v := range snapList {
			if v.ID == rs.Primary.ID {
				return fmt.Errorf("CSS snapshot %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckCssSnapshotExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.CssV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating CSS client, err: %s", err)
		}

		rs, ok := s.RootModule().Resources["huaweicloud_css_snapshot.snapshot"]
		if !ok {
			return fmt.Errorf("error checking CSS snapshot.snapshot exist, err: not found this resource")
		}

		clusterID := rs.Primary.Attributes["cluster_id"]
		snapList, err := snapshots.List(client, clusterID).Extract()
		if err != nil {
			return err
		}

		for _, v := range snapList {
			if v.ID == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("no exist CSS snapshot: %s", rs.Primary.ID)
	}
}

func testAccCssSnapshot_basic(val string) string {
	clusterName := acceptance.RandomAccResourceName()
	clusterString := testAccCssCluster_basic(clusterName, "Test@passw0rd", 1, "tag")

	return fmt.Sprintf(`
%s
resource "huaweicloud_css_snapshot" "snapshot" {
  name        = "snapshot-%s"
  description = "a snapshot created by terraform acctest"
  cluster_id  = huaweicloud_css_cluster.test.id
}
`, clusterString, val)
}

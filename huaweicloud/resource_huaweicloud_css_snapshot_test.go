package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/css/v1/snapshots"
)

func TestAccCssSnapshot_basic(t *testing.T) {
	rand := acctest.RandString(5)
	resourceName := "huaweicloud_css_snapshot.snapshot"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCssSnapshotDestroy,
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
		},
	})
}

func testAccCheckCssSnapshotDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.cssV1Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating css client, err=%s", err)
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
				return fmt.Errorf("huaweicloud css snapshot %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckCssSnapshotExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		client, err := config.cssV1Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating css client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["huaweicloud_css_snapshot.snapshot"]
		if !ok {
			return fmt.Errorf("Error checking huaweicloud css snapshot.snapshot exist, err=not found this resource")
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

		return fmt.Errorf("huaweicloud css snapshot %s is not exist", rs.Primary.ID)
	}
}

func testAccCssSnapshot_basic(val string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup" {
  name = "terraform_test_sg-%s"
  description = "terraform security group acceptance test"
}

resource "huaweicloud_css_cluster" "cluster" {
  name = "tf-css-cluster-%s"
  engine_version  = "7.1.1"
  expect_node_num = 1

  node_config {
    flavor = "ess.spec-4u16g"
    network_info {
      security_group_id = huaweicloud_networking_secgroup.secgroup.id
      subnet_id = "%s"
      vpc_id = "%s"
    }
    volume {
      volume_type = "HIGH"
      size = 40
    }
    availability_zone = "%s"
  }

  backup_strategy {
    start_time = "00:00 GMT+03:00"
    prefix     = "snapshot"
    keep_days  = 14
  }
}

resource "huaweicloud_css_snapshot" "snapshot" {
  name        = "snapshot-%s"
  description = "a snapshot created by terraform acctest"
  cluster_id  = huaweicloud_css_cluster.cluster.id
}
	`, val, val, HW_NETWORK_ID, HW_VPC_ID, HW_AVAILABILITY_ZONE, val)
}

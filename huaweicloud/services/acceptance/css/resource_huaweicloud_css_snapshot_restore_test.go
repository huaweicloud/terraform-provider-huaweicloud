package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSnapshotRestore_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {
				Source:            "hashicorp/null",
				VersionConstraint: "3.2.1",
			},
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotRestore_basic(name),
			},
		},
	})
}

func testAccSnapshotRestore_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster" "test" {
  count          = 2
  name           = "%[2]s_${count.index}"
  engine_version = "7.10.2"
  security_mode  = true
  https_enabled  = true
  password       = "Test@passw0rd"

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1

    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  backup_strategy {
    keep_days   = 1
    start_time  = "00:00 GMT+08:00"
    prefix      = "snapshot"
    bucket      = huaweicloud_obs_bucket.cssObs.bucket
    agency      = "css_obs_agency"
    backup_path = "css_repository/acctest${count.index}"
  }

  public_access {
    bandwidth         = 5
    whitelist_enabled = true
    whitelist         = "116.204.111.47,121.37.117.211"
  }
}

resource "null_resource" "cluster" {
  provisioner "local-exec" {
    command = "curl -u admin:Test@passw0rd -k -X PUT \"https://${huaweicloud_css_cluster.test[0].public_access[0].public_ip}/test_index?pretty\""
  }
}

resource "huaweicloud_css_snapshot" "snapshot" {
  name        = "snapshot-%[2]s"
  description = "a snapshot created by terraform acctest"
  cluster_id  = huaweicloud_css_cluster.test[0].id
}

resource "huaweicloud_css_snapshot_restore" "test" {
  source_cluster_id = huaweicloud_css_cluster.test[0].id
  target_cluster_id = huaweicloud_css_cluster.test[1].id
  snapshot_id       = huaweicloud_css_snapshot.snapshot.id
}
`, testAccCssBase(name), name)
}

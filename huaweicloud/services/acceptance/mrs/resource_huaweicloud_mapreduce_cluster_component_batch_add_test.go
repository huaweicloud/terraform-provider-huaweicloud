package mrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccClusterComponentBatchAdd_basic(t *testing.T) {
	rName := "huaweicloud_mapreduce_cluster_component_batch_add.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsClusterFlavorID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterComponentBatchAdd_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "components_install_mode.#", "2"),
					resource.TestCheckResourceAttr(rName, "components_install_mode.0.component", "HBase"),
					resource.TestCheckResourceAttr(rName, "components_install_mode.0.node_groups.#", "2"),
					resource.TestCheckResourceAttr(rName, "components_install_mode.1.component", "Flink"),
				),
			},
		},
	})
}

func testAccClusterComponentBatchAdd_base() string {
	name := acceptance.RandomAccResourceNameWithDash()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_mapreduce_versions" "test" {}

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%[2]s"
  type               = "CUSTOM"
  version            = try(data.huaweicloud_mapreduce_versions.test.versions[0], "")
  manager_admin_pass = "%[3]s"
  node_admin_pass    = "%[3]s"
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  component_list     = ["Hadoop", "ZooKeeper", "Ranger", "DBService"]

  master_nodes {
    flavor            = "%[4]s"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1

    assigned_roles = [
      "OMSServer:1,2",
      "SlapdServer:1,2",
      "KerberosServer:1,2",
      "KerberosAdmin:1,2",
      "quorumpeer:1,2,3",
      "NameNode:2,3",
      "Zkfc:2,3",
      "JournalNode:1,2,3",
      "ResourceManager:2,3",
      "JobHistoryServer:3",
      "DBServer:1,3",
      "HttpFS:1,3",
      "TimelineServer:3",
      "RangerAdmin:1,2",
      "UserSync:2",
      "TagSync:2",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "%[4]s"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1

    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  lifecycle {
    ignore_changes = [
      component_list,
      master_nodes[0].assigned_roles,
      custom_nodes[0].assigned_roles
    ]
  }
}
`, common.TestVpc(name), name, acceptance.RandomPassword(), acceptance.HW_MRS_CLUSTER_FLAVOR_ID)
}

func testAccClusterComponentBatchAdd_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_mapreduce_cluster_component_batch_add" "test" {
  cluster_id = huaweicloud_mapreduce_cluster.test.id

  components_install_mode {
    component = "HBase"

    node_groups {
      name           = "master_node_default_group"
      assigned_roles = ["HMaster:2,3"]
    }
    node_groups {
      name           = "node_group_1"
      assigned_roles = ["RegionServer"]
    }
  }
  components_install_mode {
    component = "Flink"

    node_groups {
      name           = "master_node_default_group"
      assigned_roles = ["FlinkResource:2,3", "FlinkServer:2,3"]
    }
  }
}
`, testAccClusterComponentBatchAdd_base())
}

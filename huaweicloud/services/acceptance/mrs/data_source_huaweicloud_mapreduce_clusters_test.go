package mrs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataClusters_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_mapreduce_clusters.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_mapreduce_clusters.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byStatus   = "data.huaweicloud_mapreduce_clusters.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byEnterpriseProjectId   = "data.huaweicloud_mapreduce_clusters.filter_by_enterprise_project_id"
		dcByEnterpriseProjectId = acceptance.InitDataSourceCheck(byEnterpriseProjectId)

		byTags   = "data.huaweicloud_mapreduce_clusters.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsClusterID(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataClusters_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "clusters.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'name' parameter.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					// Filter by 'status' parameter.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					// Filter by 'enterprise_project_id' parameter.
					dcByEnterpriseProjectId.CheckResourceExists(),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
					// Filter by 'tags' parameter.
					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
					// Check attributes.
					resource.TestCheckResourceAttrPair(byName, "clusters.0.id", "huaweicloud_mapreduce_cluster.test", "id"),
					resource.TestCheckResourceAttrPair(byName, "clusters.0.name", "huaweicloud_mapreduce_cluster.test", "name"),
					resource.TestCheckResourceAttr(byName, "clusters.0.master_node_num", "3"),
					resource.TestCheckResourceAttr(byName, "clusters.0.total_node_num", "5"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.billing_type"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.vpc_id"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.subnet_id"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.duration"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.master_node_size"),
					resource.TestMatchResourceAttr(byName, "clusters.0.component_list.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.component_list.0.component_id"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.component_list.0.component_name"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.component_list.0.component_version"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.component_list.0.component_desc"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.external_ip"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.external_alternate_ip"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.internal_ip"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.deployment_id"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.order_id"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.master_node_spec_id"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.availability_zone"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.vnc"),
					resource.TestCheckResourceAttr(byName, "clusters.0.enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(byName, "clusters.0.type", "3"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.security_group_id"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.slave_security_group_id"),
					resource.TestCheckResourceAttr(byName, "clusters.0.safe_mode", "1"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.version"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.node_public_cert_name"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.master_node_ip"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.private_ip_first"),
					resource.TestCheckResourceAttr(byName, "clusters.0.tags.%", "1"),
					resource.TestCheckResourceAttr(byName, "clusters.0.tags.foo", name),
					resource.TestCheckResourceAttr(byName, "clusters.0.log_collection", "1"),
					resource.TestCheckResourceAttr(byName, "clusters.0.node_groups.#", "2"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.node_groups.0.group_name"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.node_groups.0.node_num"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.node_groups.0.node_size"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.node_groups.0.node_spec_id"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.node_groups.0.root_volume_size"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.node_groups.0.root_volume_type"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.node_groups.0.data_volume_type"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.node_groups.0.data_volume_count"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.node_groups.0.data_volume_size"),
					resource.TestCheckResourceAttr(byName, "clusters.0.master_data_volume_type", "SSD"),
					resource.TestCheckResourceAttr(byName, "clusters.0.master_data_volume_size", "600"),
					resource.TestCheckResourceAttr(byName, "clusters.0.master_data_volume_count", "1"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.eip_id"),
					resource.TestCheckResourceAttrSet(byName, "clusters.0.eip_address"),
					resource.TestCheckResourceAttr(byName, "clusters.0.mrs_ecs_default_agency", "MRS_ECS_DEFAULT_AGENCY"),
				),
			},
		},
	})
}

func testAccDataClusters_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_mapreduce_versions" "test" {}

resource "huaweicloud_kps_keypair" "test" {
  name = "%[2]s"
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    share_type  = "PER"
    size        = 5
    charge_mode = "traffic"
  }
}

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone      = try(data.huaweicloud_availability_zones.test.names[0], "")
  name                   = "%[2]s_basic"
  type                   = "CUSTOM"
  version                = try(data.huaweicloud_mapreduce_versions.test.versions[0], "")
  manager_admin_pass     = "%[3]s"
  node_key_pair          = huaweicloud_kps_keypair.test.name
  subnet_id              = huaweicloud_vpc_subnet.test.id
  vpc_id                 = huaweicloud_vpc.test.id
  component_list         = ["Hadoop", "ZooKeeper", "Ranger", "DBService"]
  mrs_ecs_default_agency = "MRS_ECS_DEFAULT_AGENCY"
  eip_id                 = huaweicloud_vpc_eip.test.id
  enterprise_project_id  = "%[5]s"

  master_nodes {
    flavor            = "%[4]s"
    node_number       = 3
    root_volume_type  = "SSD"
    root_volume_size  = 480
    data_volume_type  = "SSD"
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
      "JobServer:2,3",
      "JobHistoryServer:2,3",
      "DBServer:1,3",
      "HttpFS:1,3",
      "TimelineServer:2,3",
      "RangerAdmin:1,2",
      "UserSync:2",
      "KerberosClient",
      "SlapdClient",
      "meta",
      "JobBalancer:1,3",
      "NodeManager:1,2",
      "DataNode:1,2"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "%[4]s"
    node_number       = 2
    root_volume_type  = "SSD"
    root_volume_size  = 480
    data_volume_type  = "SSD"
    data_volume_size  = 600
    data_volume_count = 1

    assigned_roles = [
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "DataNode",
      "meta"
    ]
  }

  tags = {
    foo = "%[2]s"
  }
}
`, common.TestVpc(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST), name, acceptance.RandomPassword(),
		acceptance.HW_MRS_CLUSTER_FLAVOR_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDataClusters_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_mapreduce_clusters" "all" {}

# Filter by 'name' parameter.
locals {
  name = huaweicloud_mapreduce_cluster.test.name
}

data "huaweicloud_mapreduce_clusters" "filter_by_name" {
  name   = local.name
  status = "running"

  depends_on = [huaweicloud_mapreduce_cluster.test]
}

locals {
  name_filter_result = [for v in data.huaweicloud_mapreduce_clusters.filter_by_name.clusters[*].name : strcontains(v, local.name)]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by 'status' parameter.
locals {
  status = huaweicloud_mapreduce_cluster.test.status
}

data "huaweicloud_mapreduce_clusters" "filter_by_status" {
  status = local.status

  depends_on = [huaweicloud_mapreduce_cluster.test]
}

locals {
  status_filter_result = [for v in data.huaweicloud_mapreduce_clusters.filter_by_status.clusters : v.status == local.status]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by 'enterprise_project_id' parameter.
locals {
  enterprise_project_id = huaweicloud_mapreduce_cluster.test.enterprise_project_id
}

data "huaweicloud_mapreduce_clusters" "filter_by_enterprise_project_id" {
  enterprise_project_id = local.enterprise_project_id

  depends_on = [huaweicloud_mapreduce_cluster.test]
}

locals {
  enterprise_project_id_filter_result = [
    for v in data.huaweicloud_mapreduce_clusters.filter_by_enterprise_project_id.clusters[*].enterprise_project_id
    : v == local.enterprise_project_id
  ]
}

output "is_enterprise_project_id_filter_useful" {
  value = length(local.enterprise_project_id_filter_result) > 0 && alltrue(local.enterprise_project_id_filter_result)
}

# Filter by 'tags' parameter.
locals {
  tags     = huaweicloud_mapreduce_cluster.test.tags
  tags_str = join(",", [for k, v in huaweicloud_mapreduce_cluster.test.tags : format("%%s*%%s", k, v)])
}

data "huaweicloud_mapreduce_clusters" "filter_by_tags" {
  tags = local.tags_str

  depends_on = [huaweicloud_mapreduce_cluster.test]
}

locals {
  tags_filter_result = [for item in data.huaweicloud_mapreduce_clusters.filter_by_tags.clusters :
  alltrue([for k, v in local.tags : v == item.tags[k]])]
}

output "is_tags_filter_useful" {
  value = length(local.tags_filter_result) > 0 && alltrue(local.tags_filter_result)
}
`, testAccDataClusters_base(name))
}

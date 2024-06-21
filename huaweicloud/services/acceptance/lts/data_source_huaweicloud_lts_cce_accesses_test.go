package lts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourceCceAccesses_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		byName   = "data.huaweicloud_lts_cce_accesses.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNameNotFound   = "data.huaweicloud_lts_cce_accesses.filter_by_not_found_name"
		dcbyNameNotFound = acceptance.InitDataSourceCheck(byNameNotFound)

		byLogGroupName   = "data.huaweicloud_lts_cce_accesses.filter_by_log_group_name"
		dcByLogGroupName = acceptance.InitDataSourceCheck(byLogGroupName)

		byLogStreamName   = "data.huaweicloud_lts_cce_accesses.filter_by_log_stream_name"
		dcByLogStreamName = acceptance.InitDataSourceCheck(byLogStreamName)

		byTags   = "data.huaweicloud_lts_cce_accesses.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {
				Source:            "hashicorp/null",
				VersionConstraint: "3.2.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceCceAccesses_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestMatchResourceAttr(byName, "accesses.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(byName, "accesses.0.name", name),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.log_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.log_stream_id", "huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(byName, "accesses.0.host_group_ids.#", "1"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.host_group_ids.0", "huaweicloud_lts_host_group.test", "id"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.cluster_id", "huaweicloud_cce_cluster.test", "id"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.#", "1"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.path_type", "container_file"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.paths.#", "1"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.paths.0", "/var"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.black_paths.#", "1"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.black_paths.0", "/var/a.log"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.name_space_regex", "test"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.pod_name_regex", "test"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.container_name_regex", "test"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.windows_log_info.0.categorys.#", "2"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.windows_log_info.0.event_level.#", "2"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.windows_log_info.0.time_offset_unit", "day"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.windows_log_info.0.time_offset", "7"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.single_log_format.#", "1"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.single_log_format.0.mode", "system"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.log_labels.%", "2"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.log_labels.log_label_key_name", "foo"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.log_labels.log_label_key_value", "bar"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.include_labels.%", "2"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.include_labels.include_label_key_name", "foo"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.include_labels.include_label_key_value", "bar"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.exclude_labels.%", "2"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.exclude_labels.exclude_label_key_name", "foo"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.exclude_labels.exclude_label_key_value", "bar"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.log_envs.%", "2"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.log_envs.log_env_key_name", "foo"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.log_envs.log_env_key_value", "bar"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.include_envs.%", "2"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.include_envs.include_env_key_name", "foo"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.include_envs.include_env_key_value", "bar"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.exclude_envs.%", "2"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.exclude_envs.exclude_env_key_name", "foo"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.exclude_envs.exclude_env_key_value", "bar"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.log_k8s.%", "2"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.log_k8s.log_k8s_key_name", "foo"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.log_k8s.log_k8s_key_value", "bar"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.include_k8s_labels.%", "2"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.include_k8s_labels.include_k8s_label_key_name", "foo"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.include_k8s_labels.include_k8s_label_key_value", "bar"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.exclude_k8s_labels.%", "2"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.exclude_k8s_labels.exclude_k8s_label_key_name", "foo"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.exclude_k8s_labels.exclude_k8s_label_key_value", "bar"),
					resource.TestCheckResourceAttr(byName, "accesses.0.tags.%", "2"),
					resource.TestCheckResourceAttr(byName, "accesses.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(byName, "accesses.0.tags.key", "value"),
					resource.TestMatchResourceAttr(byName, "accesses.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),

					dcbyNameNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_name_filter_useful", "true"),
					dcByLogGroupName.CheckResourceExists(),
					resource.TestCheckOutput("is_log_group_name_filter_useful", "true"),
					dcByLogStreamName.CheckResourceExists(),
					resource.TestCheckOutput("is_log_stream_name_filter_useful", "true"),
					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceCceAccesses_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_images" "test" {
  visibility   = "public"
  architecture = "x86"
  os           = "Ubuntu"
}

%[1]s

resource "huaweicloud_networking_secgroup" "test" {
  # install ICAgent needs the access rule with port 22 and cannot delete the default rules.
  # Without using EIP, that is a safe way to using SSH access rule.
  name = "%[2]s"
}

resource "huaweicloud_compute_instance" "test" {
  name                = "%[2]s"
  image_id            = data.huaweicloud_images_images.test.images[0].id
  flavor_id           = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids  = [huaweicloud_networking_secgroup.test.id]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  system_disk_type = "SAS"
  system_disk_size = 50

  data_disks {
    type = "SAS"
    size = "10"
  }

  # install IC agent
  user_data = <<EOF
#! /bin/bash
set +o history;
curl https://icagent-cn-north-4.obs.cn-north-4.myhuaweicloud.com/ICAgent_linux/apm_agent_install.sh > \
apm_agent_install.sh && REGION=cn-north-4 bash apm_agent_install.sh -ak %[3]s -sk %[4]s -region cn-north-4 -projectid %[5]s \
-accessip 100.125.12.150 -obsdomain obs.cn-north-4.myhuaweicloud.com -accessdomain lts-access.cn-north-4.myhuaweicloud.com;
set -o history;
EOF
}

# wait 2 minutes for the lts service to discover the server
resource "null_resource" "test" {
  depends_on = [huaweicloud_compute_instance.test]

  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    command     = "sleep 120;"
  }
}

resource "huaweicloud_lts_host_group" "test" {
  depends_on = [null_resource.test]

  name     = "%[2]s"
  type     = "linux"
  host_ids = [
    huaweicloud_compute_instance.test.id
  ]
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = replace("%[2]s", "_", "-")
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
}

resource "huaweicloud_lts_group" "test" {
  group_name  = "%[2]s"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[2]s"
}

resource "huaweicloud_lts_cce_access" "test" {
  name           = "%[2]s"
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
  host_group_ids = [huaweicloud_lts_host_group.test.id]
  cluster_id     = huaweicloud_cce_cluster.test.id

  access_config {
    path_type            = "container_file"
    paths                = ["/var"]
    black_paths          = ["/var/a.log"]
    name_space_regex     = "test"
    pod_name_regex       = "test"
    container_name_regex = "test"

    windows_log_info {
      categorys        = ["System", "Application"]
      event_level      = ["warning", "error"]
      time_offset_unit = "day"
      time_offset      = 7
    }

    single_log_format {
      mode = "system"
    }

    log_labels = {
      log_label_key_name = "foo"
      log_label_key_value = "bar"
    }

    include_labels = {
      include_label_key_name = "foo"
      include_label_key_value = "bar"
    }

    exclude_labels = {
      exclude_label_key_name = "foo"
      exclude_label_key_value = "bar"
    }

    log_envs = {
      log_env_key_name = "foo"
      log_env_key_value = "bar"
    }

    include_envs = {
      include_env_key_name = "foo"
      include_env_key_value = "bar"
    }

    exclude_envs = {
      exclude_env_key_name = "foo"
      exclude_env_key_value = "bar"
    }

    log_k8s = {
      log_k8s_key_name = "foo"
      log_k8s_key_value = "bar"
    }

    include_k8s_labels = {
      include_k8s_label_key_name = "foo"
      include_k8s_label_key_value = "bar"
    }

    exclude_k8s_labels = {
      exclude_k8s_label_key_name = "foo"
      exclude_k8s_label_key_value = "bar"
    }
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}`, common.TestVpc(name), name, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY, acceptance.HW_PROJECT_ID)
}

func testAccDatasourceCceAccesses_basic(name string) string {
	baseConfig := testAccDatasourceCceAccesses_base(name)

	return fmt.Sprintf(`
%[1]s

# In order to ensure that the filtering of the data source UT is not affected by the related resource UT and can
# correctly filter according to the attributes of the resources created by this UT, the name filtering is used as the
# basis for subsequent filtering.
# Filter by access config name
data "huaweicloud_lts_cce_accesses" "filter_by_name" {
  name = huaweicloud_lts_cce_access.test.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_lts_cce_accesses.filter_by_name.accesses[*].name : v == huaweicloud_lts_cce_access.test.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by access config name and the name is not exist
data "huaweicloud_lts_cce_accesses" "filter_by_not_found_name" {
  depends_on = [huaweicloud_lts_cce_access.test]

  name = "not_found"
}

output "is_not_found_name_filter_useful" {
  value = length(data.huaweicloud_lts_cce_accesses.filter_by_not_found_name.accesses) == 0
}

# Filter by log group name
locals {
  log_group_name = data.huaweicloud_lts_cce_accesses.filter_by_name.accesses[0].log_group_name
}

data "huaweicloud_lts_cce_accesses" "filter_by_log_group_name" {
  log_group_name = local.log_group_name
}

locals {
  log_group_name_filter_result = [
    for v in data.huaweicloud_lts_cce_accesses.filter_by_log_group_name.accesses[*].log_group_name : v == local.log_group_name
  ]
}

output "is_log_group_name_filter_useful" {
  value = length(local.log_group_name_filter_result) > 0 && alltrue(local.log_group_name_filter_result)
}

# Filter by log stream name
locals {
  log_stream_name = data.huaweicloud_lts_cce_accesses.filter_by_name.accesses[0].log_stream_name
}

data "huaweicloud_lts_cce_accesses" "filter_by_log_stream_name" {
  log_stream_name = local.log_stream_name
}

locals {
  log_stream_name_filter_result = [
    for v in data.huaweicloud_lts_cce_accesses.filter_by_log_stream_name.accesses[*].log_stream_name : v == local.log_stream_name
  ]
}

output "is_log_stream_name_filter_useful" {
  value = length(local.log_stream_name_filter_result) > 0 && alltrue(local.log_stream_name_filter_result)
}

# Filter by tags
locals {
  tags = data.huaweicloud_lts_cce_accesses.filter_by_name.accesses[0].tags
}

data "huaweicloud_lts_cce_accesses" "filter_by_tags" {
  tags = local.tags
}

locals {
  log_tags_filter_result = [
    for t in data.huaweicloud_lts_cce_accesses.filter_by_tags.accesses[*].tags : length(t) > 1 &&
	  alltrue([for k, v in t: lookup(local.tags, k, "not_found") == v])
  ]
}

output "is_tags_filter_useful" {
  value = length(local.log_tags_filter_result) > 0 && alltrue(local.log_tags_filter_result)
}
`, baseConfig)
}

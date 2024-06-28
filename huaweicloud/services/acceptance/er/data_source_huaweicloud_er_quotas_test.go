package er

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceQuotas_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_er_quotas.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byType   = "data.huaweicloud_er_quotas.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byInstanceId   = "data.huaweicloud_er_quotas.filter_by_instance_id"
		dcByInstanceId = acceptance.InitDataSourceCheck(byInstanceId)

		byNotFoundInstanceId   = "data.huaweicloud_er_quotas.filter_by_not_found_instance_id"
		dcByNotFoundInstanceId = acceptance.InitDataSourceCheck(byNotFoundInstanceId)

		byRouteTableId   = "data.huaweicloud_er_quotas.filter_by_route_table_id"
		dcByRouteTableId = acceptance.InitDataSourceCheck(byRouteTableId)

		byNotFoundRouteTableId   = "data.huaweicloud_er_quotas.filter_by_not_found_route_table_id"
		dcByNotFoundRouteTableId = acceptance.InitDataSourceCheck(byNotFoundRouteTableId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_no_filter_useful", "true"),
					resource.TestMatchResourceAttr(all, "quotas.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "quotas.0.type"),
					resource.TestCheckResourceAttrSet(all, "quotas.0.limit"),
					resource.TestCheckResourceAttrSet(all, "quotas.0.used"),
					resource.TestCheckResourceAttrSet(all, "quotas.0.unit"),

					dcByType.CheckResourceExists(),
					resource.TestCheckResourceAttr(byType, "quotas.#", "1"),
					resource.TestMatchResourceAttr(byType, "quotas.0.used", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),

					dcByInstanceId.CheckResourceExists(),
					resource.TestCheckOutput("is_instance_id_filter_useful", "true"),

					dcByNotFoundInstanceId.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_instance_id_filter_useless", "true"),

					dcByRouteTableId.CheckResourceExists(),
					resource.TestCheckOutput("is_route_table_id_filter_useful", "true"),

					dcByNotFoundRouteTableId.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_route_table_id_filter_useless", "true"),
				),
			},
		},
	})
}

func testAccDataSourceQuotas_base() string {
	var (
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)
	)

	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}

%[1]s

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)

  name = "%[2]s"
  asn  = %[3]d
}

resource "huaweicloud_er_route_table" "test" {
  instance_id = huaweicloud_er_instance.test.id
  name        = "%[2]s"
}

resource "huaweicloud_er_vpc_attachment" "test" {
  instance_id = huaweicloud_er_instance.test.id
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id

  name                   = "%[2]s"
  auto_create_vpc_routes = true
}

resource "huaweicloud_er_static_route" "test" {
  route_table_id = huaweicloud_er_route_table.test.id
  destination    = huaweicloud_vpc.test.cidr
  attachment_id  = huaweicloud_er_vpc_attachment.test.id
}`, common.TestVpc(name), name, bgpAsNum)
}

func testAccDataSourceQuotas_basic() string {
	var (
		baseConfig  = testAccDataSourceQuotas_base()
		randUUID, _ = uuid.GenerateUUID()
	)

	return fmt.Sprintf(`
%[1]s

locals {
  all_used_quota_types = ["er_instance", "route_table", "vpc_attachment", "static_route"]
  instance_used_quota_types = ["route_table", "vpc_attachment", "static_route"]
  route_table_used_quota_types = ["static_route"]
}

data "huaweicloud_er_quotas" "test" {
  depends_on = [huaweicloud_er_static_route.test]
}

# Expression interpretation:
# 1. [for _, v in quotaList: v if v.used > 0]: Filter the list to find quotas with usage greater than 0.
# 2. setintersection(quotaListA, quotaListB): Find the union of two lists.
# 3. length(setsubtract(quotaListA, quotaListB)) == 0: There are no different elements between two lists.
output "is_no_filter_useful" {
  value = length(setsubtract(setintersection([for _, v in data.huaweicloud_er_quotas.test.quotas: v.type if v.used > 0],
            local.all_used_quota_types), local.all_used_quota_types)) == 0
}

# Filter by type
data "huaweicloud_er_quotas" "filter_by_type" {
  depends_on = [huaweicloud_er_instance.test]

  type = "er_instance"
}

# Filter by instance ID
data "huaweicloud_er_quotas" "filter_by_instance_id" {
  depends_on = [huaweicloud_er_static_route.test]

  instance_id = huaweicloud_er_instance.test.id
}

# Expression interpretation:
# 1. !contains(quotaList, "er_instance"): The element type of the quota list does not contain 'er_instance'.
# 2. [for _, v in quotaList: v if v.used > 0]: Filter the list to find quotas with usage greater than 0.
# 3. length(setsubtract(quotaListA, quotaListB)) == 0: There are no different elements between two lists.
output "is_instance_id_filter_useful" {
  value = !contains(data.huaweicloud_er_quotas.filter_by_instance_id.quotas[*].type,
            "er_instance") && length(setsubtract([for _, v in data.huaweicloud_er_quotas.filter_by_instance_id.quotas: v.type if v.used > 0],
            local.instance_used_quota_types)) == 0
}

# Filter by instance ID and the ID is not exist.
data "huaweicloud_er_quotas" "filter_by_not_found_instance_id" {
  depends_on = [huaweicloud_er_static_route.test]

  instance_id = "%[2]s"
}

# Expression interpretation:
# 1. !contains(list, "er_instance"): The element type of the quota list does not contain 'er_instance'.
# 2. [for _, v in quotaList: v if v.used > 0]: Filter the list to find quotas with usage greater than 0.
output "is_not_found_instance_id_filter_useless" {
  value = !contains(data.huaweicloud_er_quotas.filter_by_not_found_instance_id.quotas[*].type,
            "er_instance") && length([for _, v in data.huaweicloud_er_quotas.filter_by_not_found_instance_id.quotas: v.type if v.used > 0]) == 0
}


# Filter by route table ID
data "huaweicloud_er_quotas" "filter_by_route_table_id" {
  depends_on = [huaweicloud_er_static_route.test]

  route_table_id = huaweicloud_er_route_table.test.id
}

# Expression interpretation:
# 1. !contains(quotaList, "route_table"): The element type of the quota list does not contain 'route_table'.
# 2. [for _, v in quotaList: v if v.used > 0]: Filter the list to find quotas with usage greater than 0.
# 3. length(setsubtract(quotaListA, quotaListB)) == 0: There are no different elements between two lists.
output "is_route_table_id_filter_useful" {
  value = !contains(data.huaweicloud_er_quotas.filter_by_route_table_id.quotas[*].type,
            "route_table") && length(setsubtract([for _, v in data.huaweicloud_er_quotas.filter_by_route_table_id.quotas: v.type if v.used > 0],
            local.route_table_used_quota_types)) == 0
}

# Filter by route table ID and the ID is not exist.
data "huaweicloud_er_quotas" "filter_by_not_found_route_table_id" {
  depends_on = [huaweicloud_er_static_route.test]

  route_table_id = "%[2]s"
}

# Expression interpretation:
# 1. !contains(list, "route_table"): The element type of the quota list does not contain 'route_table'.
# 2. [for _, v in quotaList: v if v.used > 0]: Filter the list to find quotas with usage greater than 0.
output "is_not_found_route_table_id_filter_useless" {
  value = !contains(data.huaweicloud_er_quotas.filter_by_not_found_route_table_id.quotas[*].type,
            "route_table") && length([for _, v in data.huaweicloud_er_quotas.filter_by_not_found_route_table_id.quotas: v.type if v.used > 0]) == 0
}
`, baseConfig, randUUID)
}

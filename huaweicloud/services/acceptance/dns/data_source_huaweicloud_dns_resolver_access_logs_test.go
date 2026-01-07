package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataResolverAccessLogs_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dns_resolver_access_logs.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byVpcId   = "data.huaweicloud_dns_resolver_access_logs.filter_by_vpc_id"
		dcByVpcId = acceptance.InitDataSourceCheck(byVpcId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceResolverAccessLogs_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Filter without any parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "access_logs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "access_logs.0.id"),
					resource.TestCheckResourceAttrSet(all, "access_logs.0.lts_group_id"),
					resource.TestCheckResourceAttrSet(all, "access_logs.0.lts_topic_id"),
					resource.TestMatchResourceAttr(all, "access_logs.0.vpc_ids.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'vpc_id' parameter.
					dcByVpcId.CheckResourceExists(),
					resource.TestCheckOutput("is_vpc_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataResolverAccessLogs_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 7
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s"
}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/24"
}

resource "huaweicloud_dns_resolver_access_log" "test" {
  lts_group_id = huaweicloud_lts_group.test.id
  lts_topic_id = huaweicloud_lts_stream.test.id
  vpc_ids      = [huaweicloud_vpc.test.id]
}`, name)
}

func testAccDataSourceResolverAccessLogs_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Filter without any parameters.
data "huaweicloud_dns_resolver_access_logs" "test" {
  depends_on = [huaweicloud_dns_resolver_access_log.test]
}

# Filter by 'vpc_id' parameter.
locals {
  vpc_id = huaweicloud_vpc.test.id
}

data "huaweicloud_dns_resolver_access_logs" "filter_by_vpc_id" {
  vpc_id = local.vpc_id

  depends_on = [huaweicloud_dns_resolver_access_log.test]
}

locals {
  filter_result_by_vpc_id = [for v in data.huaweicloud_dns_resolver_access_logs.filter_by_vpc_id.access_logs[*].vpc_ids :
  contains(v, local.vpc_id)]
}

output "is_vpc_id_filter_useful" {
  value = length(local.filter_result_by_vpc_id) > 0 && alltrue(local.filter_result_by_vpc_id)
}`, testAccDataResolverAccessLogs_base(name))
}

package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsPublicationSubscriptionProfiles_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_publication_subscription_profiles.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRdsReplicationProfiles_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "profiles.#"),
					resource.TestCheckResourceAttrSet(dataSource, "profiles.0.profile_id"),
					resource.TestCheckResourceAttrSet(dataSource, "profiles.0.profile_name"),
					resource.TestCheckResourceAttrSet(dataSource, "profiles.0.agent_type"),
					resource.TestCheckResourceAttrSet(dataSource, "profiles.0.is_def_profile"),
					resource.TestCheckOutput("agent_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceRdsReplicationProfiles_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_publication_subscription_profiles" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_rds_publication_subscription_profiles" "agent_type_filter" {
  depends_on = [data.huaweicloud_rds_publication_subscription_profiles.test]

  instance_id = "%[1]s"
  agent_type  = data.huaweicloud_rds_publication_subscription_profiles.test.profiles[0].agent_type
}
locals {
  agent_type =data.huaweicloud_rds_publication_subscription_profiles.test.profiles[0].agent_type
}

output "agent_type_filter_is_useful" {
  value = length(data.huaweicloud_rds_publication_subscription_profiles.agent_type_filter.profiles) > 0 && alltrue(
    [for v in data.huaweicloud_rds_publication_subscription_profiles.agent_type_filter.profiles[*].agent_type :
    v == local.agent_type]
  )
}
`, acceptance.HW_RDS_INSTANCE_ID)
}

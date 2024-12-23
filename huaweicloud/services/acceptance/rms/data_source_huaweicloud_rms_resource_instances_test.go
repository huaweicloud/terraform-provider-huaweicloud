package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceResourceInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_resource_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceResourceInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_without_any_tag_filter_useful", "true"),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceResourceInstances_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_resource_instances" "test" {
  resource_type = "config:policyAssignments"

  depends_on = [
    huaweicloud_rms_policy_assignment.test1,
    huaweicloud_rms_policy_assignment.test2,
  ]
}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_rms_resource_instances.test.resources) >= 1
}

data "huaweicloud_rms_resource_instances" "without_any_tag_filter" {
  resource_type   = "config:policyAssignments"
  without_any_tag = true

  depends_on = [
    huaweicloud_rms_policy_assignment.test1,
    huaweicloud_rms_policy_assignment.test2,
  ]
}

output "is_without_any_tag_filter_useful" {
  value = length(data.huaweicloud_rms_resource_instances.without_any_tag_filter.resources) >= 1
}

data "huaweicloud_rms_resource_instances" "tags_filter" {
  resource_type = "config:policyAssignments"

  tags {
    key    = "foo" 
    values = ["bar"]
  }

  depends_on = [
    huaweicloud_rms_policy_assignment.test1,
    huaweicloud_rms_policy_assignment.test2,
  ]
}

output "is_tags_filter_useful" {
  value = length(data.huaweicloud_rms_resource_instances.tags_filter.resources) >= 1 && alltrue(
    [for res in data.huaweicloud_rms_resource_instances.tags_filter.resources[*] : contains(res.tags, {key = "foo", value = "bar"})]
  )
}
`, testDataSourceResourceInstances_base())
}

func testDataSourceResourceInstances_base() string {
	name := acceptance.RandomAccResourceNameWithDash()
	bucketBasicConfig := testAccPolicyAssignment_periodConfig(name)
	ecsBasicConfig := testAccPolicyAssignment_ecsConfig(name)
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_policy_definitions" "test1" {
  name = "cts-obs-bucket-track"
}

resource "huaweicloud_rms_policy_assignment" "test1" {
  name                 = "%[2]s1"
  description          = "An account is noncompliant if none of its CTS trackers track specified OBS buckets."
  period               = "One_Hour"
  policy_definition_id = try(data.huaweicloud_rms_policy_definitions.test1.definitions[0].id, "")
  status               = "Disabled"

  parameters = {
    trackBucket = "\"${huaweicloud_obs_bucket.complian.bucket}\""
  }
}

data "huaweicloud_rms_policy_definitions" "test2" {
  name = "allowed-ecs-flavors"
}

%[3]s

resource "huaweicloud_rms_policy_assignment" "test2" {
  name                 = "%[2]s2"
  description          = "An ECS is noncompliant if its flavor is not in the specified flavor list (filter by resource ID)."
  policy_definition_id = try(data.huaweicloud_rms_policy_definitions.test2.definitions[0].id, "")
  status               = "Disabled"

  policy_filter {
    region            = "%[4]s"
    resource_provider = "ecs"
    resource_type     = "cloudservers"
    resource_id       = huaweicloud_compute_instance.test.id
  }

  parameters = {
    listOfAllowedFlavors = "[\"${data.huaweicloud_compute_flavors.test.ids[0]}\"]"
  }

  tags = {
    foo = "bar"
  }
}
`, bucketBasicConfig, name, ecsBasicConfig, acceptance.HW_REGION_NAME)
}

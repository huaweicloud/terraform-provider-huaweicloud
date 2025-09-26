package aom

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEvents_basic(t *testing.T) {
	var (
		all            = "data.huaweicloud_aom_events.all"
		dcForAllEvents = acceptance.InitDataSourceCheck(all)

		allWithStep            = "data.huaweicloud_aom_events.all_with_step"
		dcForAllWithStepEvents = acceptance.InitDataSourceCheck(allWithStep)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEvents_basic,
				Check: resource.ComposeTestCheckFunc(
					dcForAllEvents.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "events.#", regexp.MustCompile(`^[0-9]*$`)),
					// Check the attributes.
					resource.TestCheckResourceAttrSet(all, "events.0.id"),
					resource.TestCheckResourceAttrSet(all, "events.0.event_sn"),
					resource.TestCheckResourceAttrSet(all, "events.0.starts_at"),
					resource.TestCheckResourceAttrSet(all, "events.0.ends_at"),
					resource.TestCheckResourceAttrSet(all, "events.0.arrives_at"),
					resource.TestCheckResourceAttrSet(all, "events.0.timeout"),
					resource.TestCheckResourceAttrSet(all, "events.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(all, "events.0.metadata.%"),
					resource.TestCheckResourceAttrSet(all, "events.0.annotations"),
					resource.TestCheckResourceAttrSet(all, "events.0.policy"),
					// With step.
					dcForAllWithStepEvents.CheckResourceExists(),
					resource.TestMatchResourceAttr(allWithStep, "events.#", regexp.MustCompile(`^[0-9]*$`)),
					resource.TestCheckResourceAttrSet(allWithStep, "events.0.id"),
					resource.TestCheckResourceAttrSet(allWithStep, "events.0.event_sn"),
					resource.TestCheckResourceAttrSet(allWithStep, "events.0.starts_at"),
					resource.TestCheckResourceAttrSet(allWithStep, "events.0.ends_at"),
					resource.TestCheckResourceAttrSet(allWithStep, "events.0.arrives_at"),
					resource.TestCheckResourceAttrSet(allWithStep, "events.0.timeout"),
					resource.TestCheckResourceAttrSet(allWithStep, "events.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(allWithStep, "events.0.metadata.%"),
					resource.TestCheckResourceAttrSet(allWithStep, "events.0.annotations"),
					resource.TestCheckResourceAttrSet(allWithStep, "events.0.policy"),
				),
			},
		},
	})
}

const testAccDataSourceEvents_basic = `
data "huaweicloud_aom_events" "all" {
  time_range = "-1.-1.60"
}

data "huaweicloud_aom_events" "all_with_step" {
  time_range = "-1.-1.60"
  step       = 90000
}
`

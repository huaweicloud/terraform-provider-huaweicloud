package lts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lts"
)

func getMetricRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("lts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	return lts.GetMetricRuleById(client, state.Primary.ID)
}

func TestAccMetricRule_basic(t *testing.T) {
	var (
		metricRule   interface{}
		rName        = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_lts_metric_rule.test"
		rc           = acceptance.InitResourceCheck(
			resourceName,
			&metricRule,
			getMetricRuleResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMetricRule_basic_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "disable"),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_id", "huaweicloud_lts_group.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "log_stream_id", "huaweicloud_lts_stream.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "sampler.0.type", "none"),
					resource.TestCheckResourceAttr(resourceName, "sampler.0.ratio", "1"),
					resource.TestCheckResourceAttr(resourceName, "sinks.0.type", "aom"),
					resource.TestCheckResourceAttr(resourceName, "sinks.0.metric_name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "sinks.0.name", "huaweicloud_aom_prom_instance.test.0", "prom_name"),
					resource.TestCheckResourceAttrPair(resourceName, "sinks.0.instance_id", "huaweicloud_aom_prom_instance.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "aggregator.0.type", "countKeyword"),
					resource.TestCheckResourceAttr(resourceName, "aggregator.0.field", "event_type"),
					resource.TestCheckResourceAttr(resourceName, "aggregator.0.group_by.0", "project_id"),
					resource.TestCheckResourceAttr(resourceName, "aggregator.0.keyword", "global"),
					resource.TestCheckResourceAttr(resourceName, "window_size", "PT5S"),
					resource.TestCheckResourceAttr(resourceName, "report", "false"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccMetricRule_basic_step2(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "status", "enable"),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_id", "huaweicloud_lts_group.test.1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "log_stream_id", "huaweicloud_lts_stream.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "sampler.0.type", "random"),
					resource.TestCheckResourceAttr(resourceName, "sampler.0.ratio", "0.5"),
					resource.TestCheckResourceAttr(resourceName, "sinks.0.type", "aom"),
					resource.TestCheckResourceAttr(resourceName, "sinks.0.metric_name", updateName),
					resource.TestCheckResourceAttrPair(resourceName, "sinks.0.name", "huaweicloud_aom_prom_instance.test.1", "prom_name"),
					resource.TestCheckResourceAttrPair(resourceName, "sinks.0.instance_id", "huaweicloud_aom_prom_instance.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "aggregator.0.type", "count"),
					resource.TestCheckResourceAttr(resourceName, "aggregator.0.field", "trace_id"),
					resource.TestCheckResourceAttr(resourceName, "aggregator.0.group_by.0", "hostIP"),
					resource.TestCheckResourceAttr(resourceName, "window_size", "PT1M"),
					resource.TestCheckResourceAttr(resourceName, "report", "true"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.type", "and"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filters.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filters.0.type", "and"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filters.0.filters.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terrafrom script"),
				),
			},
			{
				Config: testAccMetricRule_basic_step3(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "aggregator.0.keyword", ""),
					resource.TestCheckResourceAttr(resourceName, "aggregator.0.group_by.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filters.0.filters.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMetricRule_basic_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_prom_instance" "test" {
  count      = 2
  prom_name = "%[1]s_${count.index}"
  prom_type = "REMOTE_WRITE"
}

resource "huaweicloud_lts_group" "test" {
  count       = 2
  group_name  = "%[1]s_${count.index}"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  count       = 2
  group_id    = huaweicloud_lts_group.test[count.index].id
  stream_name = "%[1]s_${count.index}"
}

resource "huaweicloud_lts_structing_template" "test" {
  count = 2

  log_group_id  = huaweicloud_lts_group.test[count.index].id
  log_stream_id = huaweicloud_lts_stream.test[count.index].id
  template_name = "CTS"
  template_type = "built_in"
}
`, rName)
}

func testAccMetricRule_basic_step1(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_metric_rule" "test" {
  name          = "%[2]s"
  status        = "disable"
  log_group_id  = huaweicloud_lts_group.test[0].id
  log_stream_id = huaweicloud_lts_stream.test[0].id

  sampler {
    type  = "none"
    ratio = "1"
  }

  sinks {
    type        = "aom"
    metric_name = "%[2]s"
    name        = huaweicloud_aom_prom_instance.test[0].prom_name
    instance_id = huaweicloud_aom_prom_instance.test[0].id
  }

  aggregator {
    type     = "countKeyword"
    field    = "event_type"
    group_by = ["project_id"]
    keyword  = "global"
  }

  window_size = "PT5S"
  report      = false

  filter {}

  description = "Created by terrafrom script"
}
`, testAccMetricRule_basic_base(rName), rName)
}

func testAccMetricRule_basic_step2(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_metric_rule" "test" {
  name          = "%[2]s"
  status        = "enable"
  log_group_id  = huaweicloud_lts_group.test[1].id
  log_stream_id = huaweicloud_lts_stream.test[1].id

  sampler {
    type  = "random"
    ratio = "0.5"
  }

  sinks {
    type        = "aom"
    metric_name = "%[2]s"
    name        = huaweicloud_aom_prom_instance.test[1].prom_name
    instance_id = huaweicloud_aom_prom_instance.test[1].id
  }

  aggregator {
    type     = "count"
    field    = "trace_id"
    group_by = ["hostIP"]
  }

  window_size = "PT1M"
  report      = true

  filter {
    type = "and"

    filters {
      type = "and"

      filters {
        type  = "gt"
        key   = "code"
        value = "200"
      }
      filters {
        type = "fieldExist"
        key  = "event_type"
      }
    }
  }

  description = "Updated by terrafrom script"
}
`, testAccMetricRule_basic_base(rName), updateName)
}

func testAccMetricRule_basic_step3(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_metric_rule" "test" {
  name          = "%[2]s"
  status        = "enable"
  log_group_id  = huaweicloud_lts_group.test[1].id
  log_stream_id = huaweicloud_lts_stream.test[1].id

  sampler {
    type  = "random"
    ratio = "0.5"
  }

  sinks {
    type        = "aom"
    metric_name = "%[2]s"
    name        = huaweicloud_aom_prom_instance.test[1].prom_name
    instance_id = huaweicloud_aom_prom_instance.test[1].id
  }

  aggregator {
    type  = "count"
    field = "trace_id"
  }

  window_size = "PT1M"
  report      = true

  filter {
    type = "and"

    filters {}
  }
}
`, testAccMetricRule_basic_base(rName), updateName)
}

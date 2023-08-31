package modelarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceModelartServices_basic(t *testing.T) {
	rName := "data.huaweicloud_modelarts_services.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()
	name2 := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsHasSubscribeModel(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceModelartServices_basic(name, name2),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "services.0.id"),
					resource.TestCheckResourceAttrSet(rName, "services.0.name"),
					resource.TestCheckResourceAttrSet(rName, "services.0.workspace_id"),
					resource.TestCheckResourceAttrSet(rName, "services.0.description"),
					resource.TestCheckResourceAttrSet(rName, "services.0.status"),
					resource.TestCheckResourceAttrSet(rName, "services.0.infer_type"),
					resource.TestCheckResourceAttrSet(rName, "services.0.is_free"),
					resource.TestCheckResourceAttrSet(rName, "services.0.schedule.#"),
					resource.TestCheckResourceAttrSet(rName, "services.0.additional_properties.#"),
					resource.TestCheckResourceAttrSet(rName, "services.0.invocation_times"),
					resource.TestCheckResourceAttrSet(rName, "services.0.failed_times"),
					resource.TestCheckResourceAttrSet(rName, "services.0.is_shared"),
					resource.TestCheckResourceAttrSet(rName, "services.0.shared_count"),
					resource.TestCheckResourceAttrSet(rName, "services.0.owner"),

					resource.TestCheckOutput("service_id_filter_is_useful", "true"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),

					resource.TestCheckOutput("model_id_filter_is_useful", "true"),

					resource.TestCheckOutput("workspace_id_filter_is_useful", "true"),

					resource.TestCheckOutput("infer_type_filter_is_useful", "true"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceModelartServices_basic(name, name2 string) string {
	baseServiceConfig := testModelartsServicesConfig_service(name, name2)

	return fmt.Sprintf(`
data "huaweicloud_modelarts_models" "test" {
  description = "subscribe from market"
}

%[1]s

data "huaweicloud_modelarts_services" "test" {
  depends_on = [
    huaweicloud_modelarts_service.test,
    huaweicloud_modelarts_service.test2
  ]
}

data "huaweicloud_modelarts_services" "service_id_filter" {
  service_id = huaweicloud_modelarts_service.test.id

  depends_on = [
    huaweicloud_modelarts_service.test,
    huaweicloud_modelarts_service.test2
  ]
}
output "service_id_filter_is_useful" {
  value = length(data.huaweicloud_modelarts_services.service_id_filter.services) > 0 && alltrue(
    [for v in data.huaweicloud_modelarts_services.service_id_filter.services[*].id : v == huaweicloud_modelarts_service.test.id]
  )
}

data "huaweicloud_modelarts_services" "name_filter" {
  name = "%[2]s"
  depends_on = [
    huaweicloud_modelarts_service.test,
    huaweicloud_modelarts_service.test2
  ]  
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_modelarts_services.name_filter.services) > 0 && alltrue(
    [for v in data.huaweicloud_modelarts_services.name_filter.services[*].name : v == "%[2]s"]
  )
}

data "huaweicloud_modelarts_services" "model_id_filter" {
  model_id = data.huaweicloud_modelarts_models.test.models.0.id
  depends_on = [
    huaweicloud_modelarts_service.test,
    huaweicloud_modelarts_service.test2
  ]  
}
output "model_id_filter_is_useful" {
  value = length(data.huaweicloud_modelarts_services.model_id_filter.services) > 0
}

data "huaweicloud_modelarts_services" "infer_type_filter" {
  infer_type = "real-time"
  depends_on = [
    huaweicloud_modelarts_service.test,
    huaweicloud_modelarts_service.test2
  ]
}
output "infer_type_filter_is_useful" {
  value = length(data.huaweicloud_modelarts_services.infer_type_filter.services) > 0 && alltrue(
    [for v in data.huaweicloud_modelarts_services.infer_type_filter.services[*].infer_type : v == "real-time"]
  )
}

data "huaweicloud_modelarts_services" "workspace_id_filter" {
  workspace_id = "0"
  depends_on = [
    huaweicloud_modelarts_service.test,
    huaweicloud_modelarts_service.test2
  ]
}
output "workspace_id_filter_is_useful" {
  value = length(data.huaweicloud_modelarts_services.workspace_id_filter.services) > 0 && alltrue(
    [for v in data.huaweicloud_modelarts_services.workspace_id_filter.services[*].workspace_id : v == "0"]
  )
}

data "huaweicloud_modelarts_services" "status_filter" {
  status = huaweicloud_modelarts_service.test.status
  depends_on = [
    huaweicloud_modelarts_service.test,
    huaweicloud_modelarts_service.test2
  ]
}
output "status_filter_is_useful" {
  value = length(data.huaweicloud_modelarts_services.status_filter.services) > 0 && alltrue(
    [for v in data.huaweicloud_modelarts_services.status_filter.services[*].status : v == huaweicloud_modelarts_service.test.status]
  )
}
`, baseServiceConfig, name)
}

func testModelartsServicesConfig_service(name, name2 string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name         = "%[1]s"
  display_name = "This is a acc test"
}

resource "huaweicloud_modelarts_service" "test" {
  name        = "%[1]s"
  description = "This is a acc test"
  infer_type  = "real-time"

  config {
    specification  = "modelarts.vm.cpu.2u"
    instance_count = 1
    weight         = 100
    model_id       = data.huaweicloud_modelarts_models.test.models.0.id
    envs = {
      "a" : "1",
      "b" : "2"
    }
  }

  additional_properties {
    smn_notification {
      topic_urn = huaweicloud_smn_topic.test.id
      events    = [3]
    }
    log_report_channels {
      type = "LTS"
    }
  }

  schedule {
    type      = "stop"
    duration  = 1
    time_unit = "HOURS"
  }
}

resource "huaweicloud_smn_topic" "test2" {
  name         = "%[2]s"
  display_name = "This is a acc test"
}

resource "huaweicloud_modelarts_service" "test2" {
  name        = "%[2]s"
  description = "This is a acc test"
  infer_type  = "real-time"

  config {
    specification  = "modelarts.vm.cpu.2u"
    instance_count = 1
    weight         = 100
    model_id       = data.huaweicloud_modelarts_models.test.models.1.id
    envs = {
      "a" : "1",
      "b" : "2"
    }
  }

  additional_properties {
    smn_notification {
      topic_urn = huaweicloud_smn_topic.test2.id
      events    = [3]
    }
    log_report_channels {
      type = "LTS"
    }
  }

  schedule {
    type      = "stop"
    duration  = 1
    time_unit = "HOURS"
  }
}
`, name, name2)
}

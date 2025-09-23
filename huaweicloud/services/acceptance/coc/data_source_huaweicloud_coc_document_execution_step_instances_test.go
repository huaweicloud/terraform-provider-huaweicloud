package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocDocumentExecutionStepInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_document_execution_step_instances.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocDocumentExecutionStepInstances_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.execution_instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.execution_step_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.properties.%"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocDocumentExecutionStepInstances_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_coc_document" "test" {
  name                  = "%[3]s"
  enterprise_project_id = "0"
  risk_level            = "LOW"
  description           = "This is a description"
  content               = <<EOF
meta:
  name: %[3]s
  schemaVersion: 3
  description: ""
inputs:
  instance_list:
    type: string
    default: '[{"name":"coc-test","resource_id":"${huaweicloud_compute_instance.test.id}",\
"region_id":"${huaweicloud_compute_instance.test.region}","provider":"ECS","type":"CLOUDSERVERS"}]'
    paramSource: INSTANCE
    description: ""
    region: "${huaweicloud_compute_instance.test.region}"
steps:
  - name: task1
    task: hwc:runbook:scripts@1.0.0
    inputs:
      wait_uniagent_time: 0
      script_uuid: "${huaweicloud_coc_script.test.id}"
      instance_list: "{{$globals.instance_list}}"
      execute_user: root
      timeout: 300
      script_params: '[{"param_name":"name","param_value":"zhangsan"}]'
    outputs: []
  EOF
  tags = {
    key   = "key1"
    value = "value1"
  }
}

resource "huaweicloud_coc_document_execute" "test" {
  document_id = huaweicloud_coc_document.test.id
  parameters {
    key   = "instance_list"
    value = jsonencode([
      {
        "name": huaweicloud_compute_instance.test.name,
        "resource_id": huaweicloud_compute_instance.test.id,
        "region_id": huaweicloud_compute_instance.test.region,
        "provider": "ECS",
        "type": "CLOUDSERVERS"
      }
    ])
  }
  parameters {
    key   = "domain_id"
    value = "%[4]s"
  }
  parameters {
    key   = "project_id"
    value = "%[5]s"
  }
  parameters {
    key   = "agency_urn"
    value = "iam::%[4]s:agency:ServiceAgencyForCOC"
  }
}

data "huaweicloud_coc_document_execution_steps" "test" {
  execution_id = huaweicloud_coc_document_execute.test.id
}

locals {
  execution_step_id = [for v in data.huaweicloud_coc_document_execution_steps.test.data[*].execution_step_id :
    v if v != ""][0]
}

data "huaweicloud_coc_document_execution_step_instances" "test" {
  execution_step_id = local.execution_step_id
}
`, tesScript_basic(name), testAccComputeInstance_basic(name), name, acceptance.HW_DOMAIN_ID, acceptance.HW_PROJECT_ID)
}

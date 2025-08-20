package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/coc"
)

func getDocumentExecutionResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("coc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating COC client: %s", err)
	}

	return coc.GetDocumentExecution(client, state.Primary.ID)
}

func TestAccResourceCocDocumentExecute_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_coc_document_execute.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDocumentExecutionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCocDocumentExecute_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "document_name"),
					resource.TestCheckResourceAttrSet(resourceName, "execution_parameters.#"),
					resource.TestCheckResourceAttrSet(resourceName, "execution_parameters.0.key"),
					resource.TestCheckResourceAttrSet(resourceName, "execution_parameters.0.value"),
					resource.TestCheckResourceAttrSet(resourceName, "document_version_id"),
					resource.TestCheckResourceAttrSet(resourceName, "start_time"),
					resource.TestCheckResourceAttrSet(resourceName, "update_time"),
					resource.TestCheckResourceAttrSet(resourceName, "creator"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"parameters", "sys_tags", "document_type", "end_time", "status",
					"update_time"},
			},
		},
	})
}

func TestAccResourceCocDocumentExecute_instance_id(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_coc_document_execute.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDocumentExecutionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCocDocumentExecute_instance_id(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "document_name"),
					resource.TestCheckResourceAttrSet(resourceName, "execution_parameters.#"),
					resource.TestCheckResourceAttrSet(resourceName, "execution_parameters.0.key"),
					resource.TestCheckResourceAttrSet(resourceName, "execution_parameters.0.value"),
					resource.TestCheckResourceAttrSet(resourceName, "document_version_id"),
					resource.TestCheckResourceAttrSet(resourceName, "start_time"),
					resource.TestCheckResourceAttrSet(resourceName, "update_time"),
					resource.TestCheckResourceAttrSet(resourceName, "creator"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"parameters", "sys_tags", "document_type", "end_time", "status",
					"update_time"},
			},
		},
	})
}

func testCocDocumentExecute_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_document" "test" {
  name                  = "%[1]s"
  enterprise_project_id = "0"
  risk_level            = "LOW"
  description           = "This is a description"
  content               = <<EOF
meta:
  name: %[1]s
  schemaVersion: 3
  description: ""
inputs:
  global_param1:
    type: string
    default: "global_test"
    paramSource: CUSTOM
    description: global param description
steps:
  - name: task1
    task: hwc:runbook:executeAPI@1.0.0
    inputs:
      service: ECS
      api_name: ListServersDetails
      region: cn-north-4
      path_params: '{"project_id":"%[2]s"}'
      query: "{}"
      body: "{}"
      timeout: 30
    outputs:
      - name: servers
        selector: .servers
        type: String
      - name: count
        selector: .count
        type: String
  EOF
  tags = {
    key   = "key1"
    value = "value1"
  }
}
`, name, acceptance.HW_PROJECT_ID)
}

func testCocDocumentExecute_base_instance_id(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_coc_document" "test" {
  name                  = "%[2]s"
  enterprise_project_id = "0"
  risk_level            = "LOW"
  description           = "This is a description"
  content               = <<EOF
meta:
  name: %[2]s
  schemaVersion: 3
  description: desc
inputs:
  instance_id:
    type: string
    default: "[]"
    paramSource: INSTANCE
    description: ""
    region: huaweicloud_compute_instance.test.region
steps:
  - name: query-task
    description: query task desc
    task: hwc:runbook:executeAPI@1.0.0
    inputs:
      service: ECS
      api_name: ListServersDetails
      region: huaweicloud_compute_instance.test.region
      path_params: '{"project_id":"%[3]s"}'
      query: "{}"
      body: "{}"
      timeout: 30
    outputs:
      - name: servers
        selector: .servers
        type: String
      - name: count
        selector: .count
        type: String
  EOF
  tags = {
    key   = "key1"
    value = "value1"
  }
}
`, testAccComputeInstance_basic(name), name, acceptance.HW_PROJECT_ID)
}

func testCocDocumentExecute_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_coc_document_execute" "test" {
  document_id = huaweicloud_coc_document.test.id
  version     = "v1"
  parameters {
    key   = "global_param1"
    value = "global_value"
  }
  parameters {
    key   = "domain_id"
    value = "%[2]s"
  }
  parameters {
    key   = "project_id"
    value = "%[3]s"
  }
  parameters {
    key   = "agency_urn"
    value = "iam::%[2]s:agency:ServiceAgencyForCOC"
  }
  description = "this is description"
}
`, testCocDocumentExecute_base(name), acceptance.HW_DOMAIN_ID, acceptance.HW_PROJECT_ID)
}

func testCocDocumentExecute_instance_id(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_coc_document_execute" "test" {
  document_id = huaweicloud_coc_document.test.id
  version     = "v1"
  parameters {
    key   = "instance_id"
    value = jsonencode([{
      "name":huaweicloud_compute_instance.test.name,
      "resource_id":huaweicloud_compute_instance.test.id,
      "region_id":huaweicloud_compute_instance.test.region,
      "provider":"ECS",
      "type":"CLOUDSERVERS"
    }])
  }
  parameters {
    key   = "domain_id"
    value = "%[2]s"
  }
  parameters {
    key   = "project_id"
    value = "%[3]s"
  }
  parameters {
    key   = "agency_urn"
    value = "iam::%[2]s:agency:ServiceAgencyForCOC"
  }
  sys_tags {
    key   = "key1"
    value = "value1"
  }
  target_parameter_name = "instance_id"
  targets {
    key    = "BatchValues"
    values = jsonencode([{
        "batch_index": 1,
        "instances": [
          {
            "resource_id": huaweicloud_compute_instance.test.id,
            "name": huaweicloud_compute_instance.test.name,
            "region_id": huaweicloud_compute_instance.test.region,
            "provider": "ECS",
            "type": "CLOUDSERVERS"
          }
        ],
        "strategy": "PAUSE"
      }])
  }
  description = "this is description"
}
`, testCocDocumentExecute_base_instance_id(name), acceptance.HW_DOMAIN_ID, acceptance.HW_PROJECT_ID)
}

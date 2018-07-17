package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccRTSStackResourcesV1DataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRTSStackResourcesV1Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRTSStackResourcesV1DataSourceID("data.huaweicloud_rts_stack_resource_v1.resource_1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_rts_stack_resource_v1.resource_1", "resource_name", "random"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_rts_stack_resource_v1.resource_1", "resource_type", "OS::Heat::RandomString"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_rts_stack_resource_v1.resource_1", "resource_status", "CREATE_COMPLETE"),
				),
			},
		},
	})
}

func testAccCheckRTSStackResourcesV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find stack resource data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("stack resource data source ID not set")
		}

		return nil
	}
}

const testAccDataSourceRTSStackResourcesV1Config = `

resource "huaweicloud_rts_stack_v1" "stack_1" {
  name = "huaweicloud_rts_stack"
  disable_rollback= true
  timeout_mins=60
  template_body = <<JSON
          {
			"outputs": {
              "str1": {
                 "description": "The description of the nat server.",
                 "value": {
                   "get_resource": "random"
                 }
	          }
            },
            "heat_template_version": "2013-05-23",
            "description": "A HOT template that create a single server and boot from volume.",
            "parameters": {
              "key_name": {
                "type": "string",
                "description": "Name of existing key pair for the instance to be created.",
                "default": "KeyPair-click2cloud"
	          }
	        },
            "resources": {
               "random": {
                  "type": "OS::Heat::RandomString",
                  "properties": {
                  "length": "6"
                  }
	          }
	       }
}
JSON
}

data "huaweicloud_rts_stack_resource_v1" "resource_1" {
  stack_name = "${huaweicloud_rts_stack_v1.stack_1.name}"
  resource_name = "random"
}
`

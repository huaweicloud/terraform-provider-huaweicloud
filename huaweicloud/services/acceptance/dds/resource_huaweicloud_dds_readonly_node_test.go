package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dds"
)

func getResourceReadonlyNodeFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dds", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS client: %s", err)
	}

	return dds.GetReadonlyNodeInfo(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccResourceReadonlyNode_basic(t *testing.T) {
	var (
		rName  = "huaweicloud_dds_readonly_node.test"
		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceReadonlyNodeFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccReadonlyNode_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_DDS_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "spec_code", "dds.mongodb.s6.large.4.rr"),
					resource.TestCheckResourceAttr(rName, "size", "20"),
					resource.TestCheckResourceAttr(rName, "private_ip", "192.168.1.216"),
					resource.TestCheckResourceAttr(rName, "delay", "2"),
				),
			},
			{
				Config: testAccReadonlyNode_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "spec_code", "dds.mongodb.c6.large.4.rr"),
					resource.TestCheckResourceAttr(rName, "size", "30"),
					resource.TestCheckResourceAttr(rName, "private_ip", "192.168.1.118"),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "role"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delay"},
				ImportStateIdFunc:       testAccReadonlyNodeImportStateFunc(rName),
			},
		},
	})
}

func testAccReadonlyNode_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dds_readonly_node" "test" {
  instance_id = "%s"
  spec_code   = "dds.mongodb.s6.large.4.rr"
  size        = "20"
  private_ip  = "192.168.1.216"
  delay       = 2
}
`, acceptance.HW_DDS_INSTANCE_ID)
}

func testAccReadonlyNode_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_dds_readonly_node" "test" {
  instance_id = "%s"
  spec_code   = "dds.mongodb.c6.large.4.rr"
  size        = "30"
  private_ip  = "192.168.1.118"
  delay       = 2
}
`, acceptance.HW_DDS_INSTANCE_ID)
}

func testAccReadonlyNodeImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, nodeId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		instanceId = rs.Primary.Attributes["instance_id"]
		nodeId = rs.Primary.ID

		if instanceId == "" || nodeId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s/%s'",
				instanceId, nodeId)
		}

		return fmt.Sprintf("%s/%s", instanceId, nodeId), nil
	}
}

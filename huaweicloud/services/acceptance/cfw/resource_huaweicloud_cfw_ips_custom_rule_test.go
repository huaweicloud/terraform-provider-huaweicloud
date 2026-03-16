package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cfw"
)

func getIpsCustomRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region       = acceptance.HW_REGION_NAME
		product      = "cfw"
		fwInstanceId = state.Primary.Attributes["fw_instance_id"]
		ipsCfwId     = state.Primary.ID
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW client: %s", err)
	}

	return cfw.GetIpsCustomRuleDetail(client, fwInstanceId, ipsCfwId)
}

func TestAccIpsCustomRule_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_cfw_ips_custom_rule.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIpsCustomRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testIpsCustomRule_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "fw_instance_id", acceptance.HW_CFW_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "ips_name", name),
					resource.TestCheckResourceAttr(rName, "action_type", "0"),
					resource.TestCheckResourceAttr(rName, "affected_os", "0"),
					resource.TestCheckResourceAttr(rName, "attack_type", "3"),
					resource.TestCheckResourceAttr(rName, "direction", "1"),
					resource.TestCheckResourceAttr(rName, "protocol", "10"),
					resource.TestCheckResourceAttr(rName, "severity", "1"),
					resource.TestCheckResourceAttr(rName, "software", "3"),
					resource.TestCheckResourceAttr(rName, "contents.#", "2"),
					resource.TestCheckResourceAttr(rName, "contents.0.content", "vvvvvv"),
					resource.TestCheckResourceAttr(rName, "contents.0.depth", "65535"),
					resource.TestCheckResourceAttr(rName, "contents.0.is_hex", "false"),
					resource.TestCheckResourceAttr(rName, "contents.0.is_ignore", "true"),
					resource.TestCheckResourceAttr(rName, "contents.0.is_uri", "false"),
					resource.TestCheckResourceAttr(rName, "contents.0.offset", "50"),
					resource.TestCheckResourceAttr(rName, "contents.0.relative_position", "0"),
					resource.TestCheckResourceAttr(rName, "dst_port.0.port_type", "1"),
					resource.TestCheckResourceAttr(rName, "dst_port.0.ports", "9008"),
					resource.TestCheckResourceAttr(rName, "src_port.0.port_type", "0"),
					resource.TestCheckResourceAttr(rName, "src_port.0.ports", "5005"),
					resource.TestCheckResourceAttrSet(rName, "config_status"),
				),
			},
			{
				Config: testIpsCustomRule_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "fw_instance_id", acceptance.HW_CFW_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "ips_name", name),
					resource.TestCheckResourceAttr(rName, "action_type", "1"),
					resource.TestCheckResourceAttr(rName, "affected_os", "2"),
					resource.TestCheckResourceAttr(rName, "attack_type", "12"),
					resource.TestCheckResourceAttr(rName, "direction", "-1"),
					resource.TestCheckResourceAttr(rName, "protocol", "10"),
					resource.TestCheckResourceAttr(rName, "severity", "0"),
					resource.TestCheckResourceAttr(rName, "software", "0"),
					resource.TestCheckResourceAttr(rName, "contents.#", "1"),
					resource.TestCheckResourceAttr(rName, "contents.0.content", "DF"),
					resource.TestCheckResourceAttr(rName, "contents.0.depth", "65500"),
					resource.TestCheckResourceAttr(rName, "contents.0.is_hex", "true"),
					resource.TestCheckResourceAttr(rName, "contents.0.is_ignore", "false"),
					resource.TestCheckResourceAttr(rName, "contents.0.is_uri", "false"),
					resource.TestCheckResourceAttr(rName, "contents.0.offset", "100"),
					resource.TestCheckResourceAttr(rName, "contents.0.relative_position", "1"),
					resource.TestCheckResourceAttr(rName, "dst_port.0.port_type", "-1"),
					resource.TestCheckResourceAttr(rName, "src_port.0.port_type", "1"),
					resource.TestCheckResourceAttr(rName, "src_port.0.ports", "6000"),
					resource.TestCheckResourceAttrSet(rName, "config_status"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testIpsCustomRuleImportState(rName),
			},
		},
	})
}

func testIpsCustomRule_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_ips_custom_rule" "test" {
  fw_instance_id = "%[1]s"
  ips_name       = "%[2]s"
  action_type    = 0
  affected_os    = 0
  attack_type    = 3
  direction      = 1
  protocol       = 10
  severity       = 1
  software       = 3

  contents {
    content           = "vvvvvv"
    depth             = 65535
    is_hex            = false
    is_ignore         = true
    is_uri            = false
    offset            = 50
    relative_position = 0
  }

  contents {
    content           = "DF"
    depth             = 65535
    is_hex            = true
    is_ignore         = false
    is_uri            = false
    offset            = 200
    relative_position = 0
  }

  dst_port {
    port_type = 1
    ports     = "9008"
  }

  src_port {
    port_type = 0
    ports     = "5005"
  }
}
`, acceptance.HW_CFW_INSTANCE_ID, name)
}

func testIpsCustomRule_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_ips_custom_rule" "test" {
  fw_instance_id = "%[1]s"
  ips_name       = "%[2]s"
  action_type    = 1
  affected_os    = 2
  attack_type    = 12
  direction      = -1
  protocol       = 10
  severity       = 0
  software       = 0

  contents {
    content           = "DF"
    depth             = 65500
    is_hex            = true
    is_ignore         = false
    is_uri            = false
    offset            = 100
    relative_position = 1
  }

  dst_port {
    port_type = -1
  }

  src_port {
    port_type = 1
    ports     = "6000"
  }
}
`, acceptance.HW_CFW_INSTANCE_ID, name)
}

func testIpsCustomRuleImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		fwInstanceId := rs.Primary.Attributes["fw_instance_id"]
		if fwInstanceId == "" {
			return "", fmt.Errorf("attribute (fw_instance_id) of Resource (%s) not found", name)
		}

		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of Resource (%s) not found", name)
		}

		return fmt.Sprintf("%s/%s", fwInstanceId, rs.Primary.ID), nil
	}
}

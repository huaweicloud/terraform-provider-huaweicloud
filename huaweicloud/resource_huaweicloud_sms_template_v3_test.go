package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestSmsTemplate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSMSTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testSmsTemplateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_sms_template.test", "name", "template_test_old"),
					resource.TestCheckResourceAttr("huaweicloud_sms_template.test", "is_template", "true"),
					resource.TestCheckResourceAttr("huaweicloud_sms_template.test", "region", "ap-southeast-1"),
					resource.TestCheckResourceAttr("huaweicloud_sms_template.test", "projectid", "0e1dae200a00f3772f62c00030321606"),
				),
			},
			{
				Config: testSmsTemplateConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_sms_template.test", "name", "template_test_new"),
				),
			},
		},
	})
}

var testSmsTemplateConfig = `
resource "huaweicloud_sms_template" "test" {
  name                = "template_test_old"
  is_template         = "true"
  region              = "ap-southeast-1"
  projectid           = "0e1dae200a00f3772f62c00030321606"
  availability_zone   = "ap-southeast-1a"
  target_server_name  = "serverName"
  flavor              = "flavor"
  volumetype          = "SAS"
  data_volume_type    = "SAS"
  target_password     = "123456"

  vpc_id = "ee3ef9f0-0c49-46d6-8da3-5ba364adc791"
  vpc_name = "vpc123"
  vpc_cidr = "192.168.0.0/16"

  nics {
    id = "autoCreate"
    name = "nics1"
    cidr = "192.168.0.0/16"
  }

  nics {
    id = "autoCreate"
    name = "nics2"
    cidr = "192.168.0.0/16"
  }

  security_groups {
    id = "autoCreate"
    name = "sg1"
  }

  security_groups {
    id = "autoCreate"
    name = "sg2"
  }

  publicip_type = "5_bgp"
  publicip_bandwidth_size = 1

  disk {
    index = 0
    name = "disk0"
    disktype = "SAS"
    size = 40
  }

  disk {
    index = 1
    name = "disk1"
    disktype = "SAS"
    size = 40
  }
}
`
var testSmsTemplateConfigUpdate = `
resource "huaweicloud_sms_template" "test" {
  name                = "template_test_new"
  is_template         = "true"
  region              = "ap-southeast-1"
  projectid           = "0e1dae200a00f3772f62c00030321606"
  availability_zone   = "ap-southeast-1a"
  target_server_name  = "serverName"
  flavor              = "flavor"
  volumetype          = "SAS"
  data_volume_type    = "SAS"
  target_password     = "123456"

  vpc_id = "ee3ef9f0-0c49-46d6-8da3-5ba364adc791"
  vpc_name = "vpc123"
  vpc_cidr = "192.168.0.0/16"

  nics {
    id = "autoCreate"
    name = "nics1"
    cidr = "192.168.0.0/16"
  }

  nics {
    id = "autoCreate"
    name = "nics2"
    cidr = "192.168.0.0/16"
  }

  security_groups {
    id = "autoCreate"
    name = "sg1"
  }

  security_groups {
    id = "autoCreate"
    name = "sg2"
  }

  publicip_type = "5_bgp"
  publicip_bandwidth_size = 1

  disk {
    index = 0
    name = "disk0"
    disktype = "SAS"
    size = 40
  }

  disk {
    index = 1
    name = "disk1"
    disktype = "SAS"
    size = 40
  }
}
`

func testAccCheckSMSTemplateDestroy(s *terraform.State) error {
	return nil
}

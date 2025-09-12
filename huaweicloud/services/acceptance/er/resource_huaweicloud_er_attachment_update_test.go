package er

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/er"
)

func getAttachmentFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("er", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Instance Client: %s", err)
	}
	return er.GetAttachmentById(client, state.Primary.Attributes["instance_id"],
		state.Primary.Attributes["attachment_id"])
}

func TestAccAttachmentUpdate_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_er_attachment_update.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getAttachmentFunc)

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		baseConfig = testAccAttachmentUpdate_base()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAttachmentUpdate_basic_step1(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "attachment_id",
						"huaweicloud_er_vpc_attachment.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "name",
						"huaweicloud_er_vpc_attachment.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by attachment update resource"),
				),
			},
			{
				Config: testAccAttachmentUpdate_basic_step2(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "attachment_id",
						"huaweicloud_er_vpc_attachment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by attachment update resource"),
				),
			},
			{
				Config: testAccAttachmentUpdate_basic_step3(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "attachment_id",
						"huaweicloud_er_vpc_attachment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by attachment update resource"),
				),
			},
		},
	})
}

func testAccAttachmentUpdate_base() string {
	var (
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)
	)

	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_er_instance" "test" {
  availability_zones = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), []) # only one az has.
  name               = "%[1]s"
  asn                = %[2]d

  lifecycle {
    ignore_changes = [
      availability_zones,
    ]
  }
}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 0)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 0), 1)
}

resource "huaweicloud_er_vpc_attachment" "test" {
  instance_id = huaweicloud_er_instance.test.id
  vpc_id      = huaweicloud_vpc.test.id # huaweicloud_vpc_subnet.test.vpc_id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  name        = "%[1]s"

  lifecycle {
    ignore_changes = [
      name,
      description,
    ]
  }
}
`, name, bgpAsNum)
}

func testAccAttachmentUpdate_basic_step1(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_er_attachment_update" "test" {
  instance_id   = huaweicloud_er_instance.test.id
  attachment_id = huaweicloud_er_vpc_attachment.test.id
  description   = "Created by attachment update resource"
}
`, baseConfig)
}

func testAccAttachmentUpdate_basic_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_er_attachment_update" "test" {
  instance_id   = huaweicloud_er_instance.test.id
  attachment_id = huaweicloud_er_vpc_attachment.test.id
  name          = "%[2]s"
}
`, baseConfig, name)
}

func testAccAttachmentUpdate_basic_step3(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_er_attachment_update" "test" {
  instance_id   = huaweicloud_er_instance.test.id
  attachment_id = huaweicloud_er_vpc_attachment.test.id
  name          = "%[2]s"
  description   = "Updated by attachment update resource"
}
`, baseConfig, name)
}

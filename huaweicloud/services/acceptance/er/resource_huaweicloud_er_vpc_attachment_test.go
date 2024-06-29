package er

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/er/v3/vpcattachments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getVpcAttachmentResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ErV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ER v3 client: %s", err)
	}

	return vpcattachments.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccVpcAttachment_basic(t *testing.T) {
	var (
		obj        vpcattachments.Attachment
		rName      = "huaweicloud_er_vpc_attachment.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
		bgpAsNum   = acctest.RandIntRange(64512, 65534)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getVpcAttachmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testVpcAttachment_basic(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Create by acc test"),
					resource.TestCheckResourceAttr(rName, "auto_create_vpc_routes", "true"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckOutput("er_route_count", "3"),
				),
			},
			{
				Config: testVpcAttachment_basic_update(updateName, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccVpcAttachmentImportStateFunc(rName),
			},
		},
	})
}

func testAccVpcAttachmentImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, attachmentId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of ER attachment is not found in the tfstate", rsName)
		}
		instanceId = rs.Primary.Attributes["instance_id"]
		attachmentId = rs.Primary.ID
		if instanceId == "" || attachmentId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<instance_id>/<id>', but got '%s/%s'",
				instanceId, attachmentId)
		}
		return fmt.Sprintf("%s/%s", instanceId, attachmentId), nil
	}
}

func testVpcAttachment_base(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}
	
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id = huaweicloud_vpc.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)

  name = "%[1]s"
  asn  = %[2]d
}
`, name, bgpAsNum)
}

func testVpcAttachment_basic(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_er_vpc_attachment" "test" {
  instance_id = huaweicloud_er_instance.test.id
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id

  name                   = "%[2]s"
  description            = "Create by acc test"
  auto_create_vpc_routes = true

  tags = {
    foo = "bar"
  }
}

data "huaweicloud_vpc_route_table" "test" {
  depends_on = [
    huaweicloud_er_vpc_attachment.test
  ]

  vpc_id = huaweicloud_vpc.test.id
  name   = "rtb-%[2]s"
}

output "er_route_count" {
  value = length([for route in data.huaweicloud_vpc_route_table.test.route : route.type == "er"])
}
`, testVpcAttachment_base(name, bgpAsNum), name)
}

func testVpcAttachment_basic_update(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_er_vpc_attachment" "test" {
  instance_id = huaweicloud_er_instance.test.id
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id

  name                   = "%[2]s"
  auto_create_vpc_routes = true

  tags = {
    owner = "terraform"
  }
}
`, testVpcAttachment_base(name, bgpAsNum), name)
}

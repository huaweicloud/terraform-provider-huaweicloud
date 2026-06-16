package er

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/er"
)

func getPropagationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("er", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ER client: %s", err)
	}

	return er.GetPropagationById(client, state.Primary.Attributes["instance_id"],
		state.Primary.Attributes["route_table_id"], state.Primary.ID)
}

func TestAccPropagation_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_er_propagation.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getPropagationResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testAccPropagation_nonExistentParentResources(),
				ExpectError: regexp.MustCompile(`error creating the propagation to the route table`),
			},
			{
				Config: testAccPropagation_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "route_table_id",
						"huaweicloud_er_route_table.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "attachment_id",
						"huaweicloud_er_vpc_attachment.test", "id"),
					resource.TestCheckResourceAttr(rName, "attachment_type", "vpc"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPropagationImportStateFunc(rName),
			},
		},
	})
}

func testAccPropagationImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, routeTableId, propagationId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of ER propagation is not found in the tfstate", rsName)
		}
		instanceId = rs.Primary.Attributes["instance_id"]
		routeTableId = rs.Primary.Attributes["route_table_id"]
		propagationId = rs.Primary.ID
		if instanceId == "" || routeTableId == "" || propagationId == "" {
			return "", fmt.Errorf("some import IDs are missing, want "+
				"'<instance_id>/<route_table_id>/<id>', but got '%s/%s/%s'",
				instanceId, routeTableId, propagationId)
		}
		return fmt.Sprintf("%s/%s/%s", instanceId, routeTableId, propagationId), nil
	}
}

func testAccPropagation_nonExistentParentResources() string {
	randomUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
resource "huaweicloud_er_propagation" "test" {
  instance_id    = "%[1]s"
  route_table_id = "%[1]s"
  attachment_id  = "%[1]s"
}
`, randomUUID.String())
}

func testAccPropagation_base(name string) string {
	bgpAsNum := acctest.RandIntRange(64512, 65534)

	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}

%[1]s

resource "huaweicloud_er_instance" "test" {
  availability_zones    = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)
  name                  = "%[2]s"
  asn                   = %[3]d
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_er_vpc_attachment" "test" {
  instance_id = huaweicloud_er_instance.test.id
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id

  name                   = "%[2]s"
  auto_create_vpc_routes = true
}

resource "huaweicloud_er_route_table" "test" {
  instance_id = huaweicloud_er_instance.test.id

  name = "%[2]s"
}
`, common.TestVpc(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST), name, bgpAsNum)
}

func testAccPropagation_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_er_propagation" "test" {
  instance_id    = huaweicloud_er_instance.test.id
  route_table_id = huaweicloud_er_route_table.test.id
  attachment_id  = huaweicloud_er_vpc_attachment.test.id
}
`, testAccPropagation_base(name))
}

func TestAccPropagation_routePolicy(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_er_propagation.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getPropagationResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckERRoutePolicyIDs(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPropagation_routePolicy_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "route_table_id",
						"huaweicloud_er_route_table.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "attachment_id",
						"data.huaweicloud_er_attachments.test", "attachments.0.id"),
					resource.TestCheckResourceAttr(rName, "attachment_type", "vpn"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccPropagation_routePolicy_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "route_table_id",
						"huaweicloud_er_route_table.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "attachment_id",
						"data.huaweicloud_er_attachments.test", "attachments.0.id"),
					resource.TestCheckResourceAttr(rName, "attachment_type", "vpn"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPropagationImportStateFunc(rName),
			},
		},
	})
}

func testAccPropagation_routePolicy_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}

data "huaweicloud_availability_zones" "test" {}

%[1]s

resource "huaweicloud_er_instance" "test" {
  availability_zones    = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)
  name                  = "%[2]s"
  asn                   = 64512
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_vpn_gateway" "test" {
  network_type          = "private"
  attachment_type       = "er"
  flavor                = "Professional1"
  er_id                 = huaweicloud_er_instance.test.id
  availability_zones    = slice(data.huaweicloud_availability_zones.test.names, 0, 2)
  name                  = "%[2]s"
  access_vpc_id         = huaweicloud_vpc.test.id
  access_subnet_id      = huaweicloud_vpc_subnet.test.id
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_er_route_table" "test" {
  instance_id = huaweicloud_er_instance.test.id
  name        = "%[2]s"
}

data "huaweicloud_er_attachments" "test" {
  instance_id = huaweicloud_er_instance.test.id
  type        = "vpn"
  name        = format("ErGateway_%%s", huaweicloud_vpn_gateway.test.id)
}
`, common.TestVpc(name), name)
}

func testAccPropagation_routePolicy_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  route_policy_ids = split(",", "%[2]s")
}

resource "huaweicloud_er_propagation" "test" {
  instance_id    = huaweicloud_er_instance.test.id
  route_table_id = huaweicloud_er_route_table.test.id
  # After the VPN gateway is created, an attachment will be created automatically.
  attachment_id  = data.huaweicloud_er_attachments.test.attachments[0].id

  route_policy {
    import_policy_id = local.route_policy_ids[0]
  }
}
`, testAccPropagation_routePolicy_base(name), acceptance.HW_ER_ROUTE_POLICY_IDS)
}

func testAccPropagation_routePolicy_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  route_policy_ids = split(",", "%[2]s")
}

resource "huaweicloud_er_propagation" "test" {
  instance_id    = huaweicloud_er_instance.test.id
  route_table_id = huaweicloud_er_route_table.test.id
  # The VPN Gateway attachment will exist until the VPN gateway is deleted.
  attachment_id  = data.huaweicloud_er_attachments.test.attachments[0].id

  route_policy {
    import_policy_id = local.route_policy_ids[1]
  }
}
`, testAccPropagation_routePolicy_base(name), acceptance.HW_ER_ROUTE_POLICY_IDS)
}

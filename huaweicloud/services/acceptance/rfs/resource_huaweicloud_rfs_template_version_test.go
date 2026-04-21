package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rfs"
)

func getRfsTemplateVersionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region       = acceptance.HW_REGION_NAME
		product      = "rfs"
		templateName = state.Primary.Attributes["template_name"]
		versionId    = state.Primary.ID
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate RFS request ID: %s", err)
	}

	return rfs.QueryRfsTemplateVersion(client, templateName, versionId, requestId)
}

func TestAccRfsTemplateVersion_basic(t *testing.T) {
	var (
		obj          interface{}
		name         = acceptance.RandomAccResourceName()
		rName        = "huaweicloud_rfs_template_version.test"
		templateName = "huaweicloud_rfs_template.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRfsTemplateVersionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRfsTemplateVersion_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "template_name", templateName, "template_name"),
					resource.TestCheckResourceAttrPair(rName, "template_id", templateName, "template_id"),
					resource.TestCheckResourceAttr(rName, "version_description", "v2 - add subnet support"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccRfsTemplateVersionImportStateIdFunc(rName),
				ImportStateVerifyIgnore: []string{
					"template_body",
					"template_uri",
					"version_description",
				},
			},
		},
	})
}

func testTemplateVersion_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_template" "test" {
  template_name        = "%s"
  template_description = "base template for version test"
  version_description  = "v1"
  template_body        = <<-EOF
variable "vpc_name" {
  type    = string
  default = "my-vpc"
}

resource "huaweicloud_vpc" "vpc" {
  name = var.vpc_name
  cidr = "172.16.0.0/16"
}
EOF
}
`, name)
}

func testRfsTemplateVersion_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rfs_template_version" "test" {
  template_name       = huaweicloud_rfs_template.test.template_name
  template_id         = huaweicloud_rfs_template.test.template_id
  version_description = "v2 - add subnet support"
  template_body       = <<-EOF
variable "vpc_name" {
  type    = string
  default = "my-vpc"
}

variable "subnet_name" {
  type    = string
  default = "my-subnet"
}

resource "huaweicloud_vpc" "vpc" {
  name = var.vpc_name
  cidr = "172.16.0.0/16"
}

resource "huaweicloud_vpc_subnet" "subnet" {
  name       = var.subnet_name
  vpc_id     = huaweicloud_vpc.vpc.id
  cidr       = "172.16.1.0/24"
  gateway_ip = "172.16.1.1"
}
EOF
}
`, testTemplateVersion_base(name))
}

func testAccRfsTemplateVersionImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", resourceName)
		}

		if rs.Primary.Attributes["template_name"] == "" {
			return "", fmt.Errorf("attribute (template_name) of resource (%s) not found", resourceName)
		}

		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of resource (%s) not found", resourceName)
		}

		return rs.Primary.Attributes["template_name"] + "/" + rs.Primary.ID, nil
	}
}

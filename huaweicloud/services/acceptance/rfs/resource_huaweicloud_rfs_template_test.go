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

func getRfsTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region       = acceptance.HW_REGION_NAME
		product      = "rfs"
		templateName = state.Primary.ID
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate RFS request ID: %s", err)
	}

	return rfs.GetRfsTemplateByName(client, templateName, requestId)
}

func TestAccRfsTemplate_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_rfs_template.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRfsTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRfsTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "template_name", name),
					resource.TestCheckResourceAttr(rName, "template_description", "test template description"),
					resource.TestCheckResourceAttr(rName, "version_description", "v1"),
					resource.TestCheckResourceAttrSet(rName, "template_id"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
					resource.TestCheckResourceAttrSet(rName, "latest_version_id"),
				),
			},
			{
				Config: testRfsTemplate_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "template_name", name),
					resource.TestCheckResourceAttr(rName, "template_description", "test template description update"),
					resource.TestCheckResourceAttr(rName, "version_description", "v1"),
					resource.TestCheckResourceAttrSet(rName, "template_id"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
					resource.TestCheckResourceAttrSet(rName, "latest_version_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"template_body",
					"template_uri",
					"version_description",
				},
			},
		},
	})
}

func testRfsTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_template" "test" {
  template_name        = "%s"
  template_description = "test template description"
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

func testRfsTemplate_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_template" "test" {
  template_name        = "%s"
  template_description = "test template description update"
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

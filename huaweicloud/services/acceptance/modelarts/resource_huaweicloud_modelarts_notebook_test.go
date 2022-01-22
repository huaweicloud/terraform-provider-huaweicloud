package modelarts

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/modelarts/v1/notebook"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getNotebookResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.ModelArtsV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("error creating ModelArts v1 client, err=%s", err)
	}

	return notebook.Get(client, state.Primary.ID)
}

func TestAccResourceNotebook_basic(t *testing.T) {
	var instance notebook.CreateOpts
	resourceName := "huaweicloud_modelarts_notebook.test"
	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getNotebookResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNotebook_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "modelarts.vm.cpu.2u"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "EFS"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.ownership", "MANAGED"),
					resource.TestCheckResourceAttr(resourceName, "auto_stop_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "image_name"),
					resource.TestCheckResourceAttrSet(resourceName, "image_swr_path"),
					resource.TestCheckResourceAttrSet(resourceName, "image_type"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
				),
			},
			{
				Config: testAccNotebook_basic(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "modelarts.vm.cpu.2u"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "EFS"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.ownership", "MANAGED"),
					resource.TestCheckResourceAttr(resourceName, "auto_stop_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "image_name"),
					resource.TestCheckResourceAttrSet(resourceName, "image_swr_path"),
					resource.TestCheckResourceAttrSet(resourceName, "image_type"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNotebook_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_modelarts_notebook" "test" {
  name      = "%s"
  flavor_id = "modelarts.vm.cpu.2u"
  image_id  = "e1a07296-22a8-4f05-8bc8-e936c8e54090"
  volume {
    type = "EFS"
  }
}
`, rName)
}

func TestAccResourceNotebook_all(t *testing.T) {
	var instance notebook.CreateOpts
	resourceName := "huaweicloud_modelarts_notebook.test"
	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	ip := "10.1.1.2"
	updateIp := "10.1.1.3"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getNotebookResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNotebook_All(name, ip),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "modelarts.vm.cpu.2u"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "EFS"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.ownership", "MANAGED"),
					resource.TestCheckResourceAttr(resourceName, "auto_stop_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "allowed_access_ips.0", ip),
					resource.TestCheckResourceAttrSet(resourceName, "image_name"),
					resource.TestCheckResourceAttrSet(resourceName, "image_swr_path"),
					resource.TestCheckResourceAttrSet(resourceName, "image_type"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
				),
			},
			{
				Config: testAccNotebook_All(updateName, updateIp),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "modelarts.vm.cpu.2u"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "EFS"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.ownership", "MANAGED"),
					resource.TestCheckResourceAttr(resourceName, "auto_stop_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", updateName),
					resource.TestCheckResourceAttr(resourceName, "allowed_access_ips.0", updateIp),
					resource.TestCheckResourceAttrSet(resourceName, "image_name"),
					resource.TestCheckResourceAttrSet(resourceName, "image_swr_path"),
					resource.TestCheckResourceAttrSet(resourceName, "image_type"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNotebook_All(rName string, ip string) string {
	return fmt.Sprintf(`
resource "huaweicloud_compute_keypair" "test" {
  name       = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "huaweicloud_modelarts_notebook" "test" {
  name        = "%s"
  flavor_id   = "modelarts.vm.cpu.2u"
  image_id    = "e1a07296-22a8-4f05-8bc8-e936c8e54090"
  description = "%s"

  allowed_access_ips = ["%s"]
  key_pair           = huaweicloud_compute_keypair.test.name

  volume {
    type = "EFS"
  }
}
`, rName, rName, rName, ip)
}

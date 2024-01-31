package ces

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getResourceGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getResourceGroup: Query the CES resource group detail
	var (
		getResourceGroupHttpUrl = "v2/{project_id}/resource-groups/{id}"
		getResourceGroupProduct = "ces"
	)
	getResourceGroupClient, err := cfg.NewServiceClient(getResourceGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CES Client: %s", err)
	}

	getResourceGroupPath := getResourceGroupClient.Endpoint + getResourceGroupHttpUrl
	getResourceGroupPath = strings.ReplaceAll(getResourceGroupPath, "{project_id}", getResourceGroupClient.ProjectID)
	getResourceGroupPath = strings.ReplaceAll(getResourceGroupPath, "{id}", state.Primary.ID)

	getResourceGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getResourceGroupResp, err := getResourceGroupClient.Request("GET", getResourceGroupPath, &getResourceGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving resource group: %s", err)
	}
	return utils.FlattenResponse(getResourceGroupResp)
}

func TestAccResourceGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_ces_resource_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getResourceGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					// resources is not set, so don't need to check it
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testResourceGroup_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"resources",
				},
			},
		},
	})
}

func TestAccResourceGroup_tags(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_ces_resource_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getResourceGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceGroup_tags(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "TAG"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testResourceGroup_tags_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value_update"),
					resource.TestCheckResourceAttr(rName, "tags.foo_update", "bar_update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceGroup_eps(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_ces_resource_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getResourceGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceGroup_eps(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "EPS"),
					resource.TestCheckResourceAttr(rName, "associated_eps_ids.0", "0"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceGroup_withEpsId(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_ces_resource_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getResourceGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					// resources is not set, so don't need to check it
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testResourceGroup_updateWithEpsId(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testResourceGroup_base(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "vm_1" {
  name               = "ecs-%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}
`, common.TestBaseComputeResources(name), name)
}

func testResourceGroup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ces_resource_group" "test" {
  name = "%s"

  resources {
    namespace = "SYS.ECS"
    dimensions {
      name  = "instance_id"
      value = huaweicloud_compute_instance.vm_1.id
    }
  }
}
`, testResourceGroup_base(name), name)
}

func testResourceGroup_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ces_resource_group" "test" {
  name = "%s-update"

  resources {
    namespace = "SYS.EVS"
    dimensions {
      name  = "disk_name"
      value = "${huaweicloud_compute_instance.vm_1.id}-sda"
    }
  }
}
`, testResourceGroup_base(name), name)
}

func testResourceGroup_tags(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ces_resource_group" "test" {
  name = "%s"
  type = "TAG"
  tags = {
    key = "value"
    foo = "bar"
  }
}
`, name)
}

func testResourceGroup_tags_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ces_resource_group" "test" {
  name = "%s-update"
  type = "TAG"
  tags = {
    key        = "value_update"
    foo_update = "bar_update"
  }
}
`, name)
}

func testResourceGroup_eps(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ces_resource_group" "test" {
  name               = "%s"
  type               = "EPS"
  associated_eps_ids = ["0"]
}
`, name)
}

func testResourceGroup_updateWithEpsId(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ces_resource_group" "test" {
  name                  = "%s-update"
  enterprise_project_id = "%s"

  resources {
    namespace = "SYS.EVS"
    dimensions {
      name  = "disk_name"
      value = "${huaweicloud_compute_instance.vm_1.id}-sda"
    }
  }
}
`, testResourceGroup_base(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

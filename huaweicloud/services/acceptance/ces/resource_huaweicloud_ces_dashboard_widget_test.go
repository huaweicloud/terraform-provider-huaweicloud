package ces

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDashboardWidgetFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		getDashboardWidgetHttpUrl = "v2/{project_id}/widgets/{widget_id}"
		getDashboardWidgetProduct = "ces"
	)
	client, err := cfg.NewServiceClient(getDashboardWidgetProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CES client: %s", err)
	}

	path := client.Endpoint + getDashboardWidgetHttpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{widget_id}", state.Primary.ID)

	getDashboardWidgetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}
	getDashboardWidgetResp, err := client.Request("GET", path, &getDashboardWidgetOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CES dashboard widget: %s", err)
	}
	return utils.FlattenResponse(getDashboardWidgetResp)
}

func TestAccDashboardWidget_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_ces_dashboard_widget.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDashboardWidgetFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDashboardWidget_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "title", name),
					resource.TestCheckResourceAttr(rName, "view", "line"),
					resource.TestCheckResourceAttr(rName, "metric_display_mode", "single"),
					resource.TestCheckResourceAttr(rName, "metrics.0.metric_name", "disk_read_bytes_rate"),
					resource.TestCheckResourceAttr(rName, "metrics.0.namespace", "SYS.ECS"),
					resource.TestCheckResourceAttr(rName, "unit", "Mibps"),
					resource.TestCheckResourceAttr(rName, "metrics.0.alias.#", "1"),
					resource.TestCheckResourceAttr(rName, "metrics.0.dimensions.0.name", "instance_id"),
					resource.TestCheckResourceAttr(rName, "metrics.0.dimensions.0.filter_type", "specific_instances"),
					resource.TestCheckResourceAttrPair(rName, "metrics.0.dimensions.0.values.0", "huaweicloud_compute_instance.test1", "id"),
					resource.TestCheckResourceAttr(rName, "location.0.left", "0"),
					resource.TestCheckResourceAttr(rName, "location.0.top", "0"),
					resource.TestCheckResourceAttr(rName, "location.0.width", "6"),
					resource.TestCheckResourceAttr(rName, "location.0.height", "6"),
					resource.TestCheckResourceAttr(rName, "properties.0.top_n", "30"),
					resource.TestMatchResourceAttr(rName,
						"created_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testDashboardWidget_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "title", name+"-update"),
					resource.TestCheckResourceAttr(rName, "view", "line"),
					resource.TestCheckResourceAttr(rName, "metric_display_mode", "multiple"),
					resource.TestCheckResourceAttr(rName, "unit", "Kibps"),
					resource.TestCheckResourceAttr(rName, "metrics.0.metric_name", "disk_read_bytes_rate"),
					resource.TestCheckResourceAttr(rName, "metrics.0.namespace", "SYS.ECS"),
					resource.TestCheckResourceAttr(rName, "metrics.0.alias.#", "2"),
					resource.TestCheckResourceAttr(rName, "metrics.0.alias.0", "x"),
					resource.TestCheckResourceAttr(rName, "metrics.0.alias.1", "y"),
					resource.TestCheckResourceAttr(rName, "metrics.0.dimensions.0.name", "instance_id"),
					resource.TestCheckResourceAttr(rName, "metrics.0.dimensions.0.filter_type", "specific_instances"),
					resource.TestCheckResourceAttrPair(rName, "metrics.0.dimensions.0.values.0", "huaweicloud_compute_instance.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "metrics.0.dimensions.0.values.1", "huaweicloud_compute_instance.test2", "id"),
					resource.TestCheckResourceAttr(rName, "metrics.1.metric_name", "disk_write_bytes_rate"),
					resource.TestCheckResourceAttr(rName, "metrics.1.namespace", "SYS.ECS"),
					resource.TestCheckResourceAttr(rName, "metrics.1.dimensions.0.name", "instance_id"),
					resource.TestCheckResourceAttr(rName, "metrics.1.dimensions.0.filter_type", "specific_instances"),
					resource.TestCheckResourceAttrPair(rName, "metrics.1.dimensions.0.values.0", "huaweicloud_compute_instance.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "metrics.1.dimensions.0.values.1", "huaweicloud_compute_instance.test2", "id"),
					resource.TestCheckResourceAttr(rName, "location.0.left", "0"),
					resource.TestCheckResourceAttr(rName, "location.0.top", "0"),
					resource.TestCheckResourceAttr(rName, "location.0.width", "12"),
					resource.TestCheckResourceAttr(rName, "location.0.height", "4"),
					resource.TestCheckResourceAttr(rName, "properties.0.top_n", "31"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"dashboard_id",
				},
			},
		},
	})
}

func testDashboardWidget_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ces_dashboard_widget" "test" {
  dashboard_id        =  huaweicloud_ces_dashboard.test.id
  title               = "%[2]s"
  view                = "line"
  metric_display_mode = "single"
  unit                = "Mibps"
		
  metrics {
    metric_name = "disk_read_bytes_rate"
    namespace   = "SYS.ECS"
    alias       = ["x"]

    dimensions  {
      name        = "instance_id"
      filter_type = "specific_instances"
      values      = [huaweicloud_compute_instance.test1.id]
    }
  }
			
  location {
    left   = 0
    top    = 0
    width  = 6
    height = 6
  }

  properties {
    top_n = 30
  }
}
`, testDashboardWidget_base(name), name)
}

func testDashboardWidget_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ces_dashboard_widget" "test" {
  dashboard_id        =  huaweicloud_ces_dashboard.test.id
  title               = "%[2]s-update"
  view                = "line"
  metric_display_mode = "multiple"
  unit                = "Kibps"
		
  metrics {
    metric_name = "disk_read_bytes_rate"
    namespace   = "SYS.ECS"
    alias       = ["x", "y"]

    dimensions {
      name        = "instance_id"
      filter_type = "specific_instances"

      values = [
        huaweicloud_compute_instance.test1.id,
        huaweicloud_compute_instance.test2.id,
      ]
    }
   }

  metrics {
    metric_name = "disk_write_bytes_rate"
    namespace   = "SYS.ECS"

    dimensions  {
      name        = "instance_id"
      filter_type = "specific_instances"
	  
      values = [
        huaweicloud_compute_instance.test1.id,
        huaweicloud_compute_instance.test2.id,
      ]
    }
  }
			
  location {
    left   = 0
    top    = 0
    width  = 12
    height = 4
  }

  properties {
    top_n = 31
  }
}
`, testDashboardWidget_base(name), name)
}

func testDashboardWidget_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}
	
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}
	
data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}
	
data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_compute_instance" "test1" {
  name               = "%[1]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  
  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
  
  eip_type = "5_bgp"

  bandwidth {
    share_type  = "PER"
    size        = 5
    charge_mode = "bandwidth"
  }
  
  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"
}

resource "huaweicloud_compute_instance" "test2" {
  name               = "%[1]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
	
  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
	
  eip_type = "5_bgp"

  bandwidth {
    share_type  = "PER"
    size        = 5
    charge_mode = "bandwidth"
  }
	
  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"

  depends_on = [
    huaweicloud_compute_instance.test1,
  ]

  provisioner "local-exec" {
    command = "sleep 70"
  }
}

resource "huaweicloud_ces_dashboard" "test" {
  name           = "%[1]s"
  row_widget_num = 1
  is_favorite    = true

  depends_on = [
    huaweicloud_compute_instance.test1,
    huaweicloud_compute_instance.test2,
  ]	
}
`, name)
}

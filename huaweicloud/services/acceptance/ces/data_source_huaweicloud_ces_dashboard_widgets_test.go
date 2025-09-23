package ces

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCesDashboardWidgets_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_dashboard_widgets.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCesDashboardWidgets_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "widgets.0.widget_id"),
					resource.TestCheckResourceAttrSet(dataSource, "widgets.0.title"),
					resource.TestCheckResourceAttrSet(dataSource, "widgets.0.view"),
					resource.TestCheckResourceAttrSet(dataSource, "widgets.0.metric_display_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "widgets.0.metrics.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "widgets.0.metrics.0.metric_name"),
					resource.TestCheckResourceAttrSet(dataSource, "widgets.0.metrics.0.dimensions.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "widgets.0.metrics.0.dimensions.0.filter_type"),
					resource.TestCheckResourceAttrSet(dataSource, "widgets.0.location.0.left"),
					resource.TestCheckResourceAttrSet(dataSource, "widgets.0.location.0.top"),
					resource.TestCheckResourceAttrSet(dataSource, "widgets.0.location.0.width"),
					resource.TestCheckResourceAttrSet(dataSource, "widgets.0.location.0.height"),
					resource.TestMatchResourceAttr(dataSource,
						"widgets.0.created_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testDataSourceCesDashboardWidgets_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_ces_dashboard_widgets" "test" {
  dashboard_id = huaweicloud_ces_dashboard.test.id

  depends_on = [
    huaweicloud_ces_dashboard_widget.test,
  ]
}
`, testDataSourceCesDashboardWidgets_base())
}

func testDataSourceCesDashboardWidgets_base() string {
	name := acceptance.RandomAccResourceName()
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
	
resource "huaweicloud_compute_instance" "test" {
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

  provisioner "local-exec" {
    command = "sleep 70"
  }
}

resource "huaweicloud_ces_dashboard" "test" {
  name           = "%[1]s"
  row_widget_num = 1
  is_favorite    = true

  depends_on = [
    huaweicloud_compute_instance.test,
  ]	
}

resource "huaweicloud_ces_dashboard_widget" "test" {
  dashboard_id        =  huaweicloud_ces_dashboard.test.id
  title               = "%[1]s"
  view                = "line"
  metric_display_mode = "single"
		  
  metrics {
    metric_name = "disk_read_bytes_rate"
    namespace   = "SYS.ECS"
    alias       = ["x"]
  
    dimensions {
      name        = "instance_id"
      filter_type = "specific_instances"
  
      values = [
        huaweicloud_compute_instance.test.id,
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
`, name)
}

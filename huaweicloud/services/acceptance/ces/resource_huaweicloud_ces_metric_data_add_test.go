package ces

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccMetricDataAdd_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testMetricDataAdd_basic(),
			},
		},
	})
}

func testMetricDataAdd_basic() string {
	currentTime := time.Now().UTC()
	timeString := currentTime.Format("2006-01-02 15:04:05")
	return fmt.Sprintf(`
resource "huaweicloud_ces_metric_data_add" "test" {
  metric {
    namespace   = "MINE.APP"
    metric_name = "cpu_util"

    dimensions {
      name  = "instance_id"
      value = "test_instance_id"
    }

    dimensions {
      name  = "disk_name"
      value = "test_disk_name"
    }
  }

  ttl          = 400
  collect_time = "%[1]s"
  value        = 0.5
  unit         = "%%"
  type         = "float"
}
`, timeString)
}

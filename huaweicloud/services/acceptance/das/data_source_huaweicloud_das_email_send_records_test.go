package das

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEmailSendRecords_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_das_email_send_records.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEmailSendRecords_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "records.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestMatchResourceAttr(all, "records.0.send_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(all, "records.0.status"),
					resource.TestCheckResourceAttrSet(all, "records.0.email"),
					resource.TestMatchResourceAttr(all, "records.0.instance_health_reports.#",
						regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "records.0.instance_health_reports.0.task_id"),
					resource.TestCheckResourceAttrSet(all, "records.0.instance_health_reports.0.instance_id"),
					resource.TestCheckResourceAttrSet(all, "records.0.instance_health_reports.0.instance_name"),
					resource.TestMatchResourceAttr(all, "records.0.instance_health_reports.0.start_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(all, "records.0.instance_health_reports.0.end_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

const testAccDataEmailSendRecords_basic string = `
data "huaweicloud_das_email_send_records" "all" {
  datastore_type = "MySQL"
}
`

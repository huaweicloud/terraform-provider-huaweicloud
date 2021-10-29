package dli

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dli/v1/flinkjob"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDliFlinkSqlJobResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.DliV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("error creating Dli v1 client, err=%s", err)
	}
	jobId, _ := strconv.Atoi(state.Primary.ID)
	return flinkjob.Get(client, jobId)
}

func TestAccResourceDliFlinkJob_basic(t *testing.T) {
	var obj flinkjob.CreateSqlJobOpts
	resourceName := "huaweicloud_dli_flinksql_job.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDliFlinkSqlJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFlinkJobResource_basic(name, acceptance.HW_REGION_NAME),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "status", "job_running"),
					resource.TestCheckResourceAttr(resourceName, "type", "flink_sql_job"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccFlinkJobResource_basic(name string, region string) string {
	return fmt.Sprintf(`
variable "sql" {
  type    = string
  default = <<EOF
CREATE SOURCE STREAM car_infos (
  car_id STRING,
  car_owner STRING,
  car_brand STRING,
  car_price INT
)
WITH (
  type = "dis",
  region = "%s",
  channel = "%s_input",
  partition_count = "1",
  encode = "csv",
  field_delimiter = ","
);

CREATE SINK STREAM audi_cheaper_than_30w (
  car_id STRING,
  car_owner STRING,
  car_brand STRING,
  car_price INT
)
WITH (
  type = "dis",
  region = "%s",
  channel = "%s_output",
  partition_key = "car_owner",
  encode = "csv",
  field_delimiter = ","
);

INSERT INTO audi_cheaper_than_30w
SELECT *
FROM car_infos
WHERE car_brand = "audia4" and car_price < 30;


CREATE SINK STREAM car_info_data (
  car_id STRING,
  car_owner STRING,
  car_brand STRING,
  car_price INT
)
WITH (
  type ="dis",
  region = "%s",
  channel = "%s_input",
  partition_key = "car_owner",
  encode = "csv",
  field_delimiter = ","
);

INSERT INTO car_info_data
SELECT "1", "lilei", "bmw320i", 28;
INSERT INTO car_info_data
SELECT "2", "hanmeimei", "audia4", 27;
EOF

}


resource "huaweicloud_dis_stream" "stream_input" {
  stream_name     = "%s_input"
  partition_count = 1
  data_type       = "CSV"
  csv_delimiter   = ","
}

resource "huaweicloud_dis_stream" "stream_output" {
  stream_name     = "%s_output"
  partition_count = 1
  data_type       = "CSV"
  csv_delimiter   = ","

}

resource "huaweicloud_dli_flinksql_job" "test" {
  name = "%s"
  type = "flink_sql_job"
  sql  = var.sql
  depends_on = [
    huaweicloud_dis_stream.stream_input,
    huaweicloud_dis_stream.stream_output,
  ]
}
`, region, name, region, name, region, name, name, name, name)
}

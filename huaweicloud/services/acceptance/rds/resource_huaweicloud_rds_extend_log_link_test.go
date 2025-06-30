package rds

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getRdsExtendLogLinkResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/xellog-download"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	instanceID := state.Primary.Attributes["instance_id"]
	fileName := state.Primary.Attributes["file_name"]

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)

	opt := golangsdk.RequestOpts{KeepResponseBody: true}
	opt.JSONBody = utils.RemoveNil(buildCreateExtendLogLinkBodyParams(fileName))
	resp, err := client.Request("POST", getPath, &opt)
	if err != nil {
		return nil, fmt.Errorf("error getting extend log link for file (%s) of the instance (%s): %s", fileName,
			instanceID, err)
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err, nil
	}
	logInfo := utils.PathSearch(fmt.Sprintf("list|[?file_name=='%s']|[0]", fileName), respBody, nil)
	if logInfo == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return logInfo, nil
}

func buildCreateExtendLogLinkBodyParams(fileName string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"file_name": fileName,
	}
	return bodyParams
}

func TestAccRdsExtendLogLink_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_extend_log_link.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRdsExtendLogLinkResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testRdsExtendLogLink_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "file_name",
						"data.huaweicloud_rds_extend_log_files.test", "files.0.file_name"),
					resource.TestCheckResourceAttrSet(rName, "file_size"),
					resource.TestCheckResourceAttrSet(rName, "file_link"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
		},
	})
}

func testRdsExtendLogLink_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.mssql.spec.se.s6.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  charging_mode     = "postPaid"

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2022_SE"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}

data "huaweicloud_rds_extend_log_files" "test" {
  instance_id = huaweicloud_rds_instance.test.id
}
`, testAccRdsInstance_base(), name, name)
}

func testRdsExtendLogLink_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_extend_log_link" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  file_name   = data.huaweicloud_rds_extend_log_files.test.files[0].file_name
}
`, testRdsExtendLogLink_base(name))
}

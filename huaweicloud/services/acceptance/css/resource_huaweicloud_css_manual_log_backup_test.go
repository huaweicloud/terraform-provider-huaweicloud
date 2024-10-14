package css

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

func getManualLogBackupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.CssV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS v1 client: %s", err)
	}

	getLogBackupJobHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/records"
	getLogBackupJobPath := client.Endpoint + getLogBackupJobHttpUrl
	getLogBackupJobPath = strings.ReplaceAll(getLogBackupJobPath, "{project_id}", client.ProjectID)
	getLogBackupJobPath = strings.ReplaceAll(getLogBackupJobPath, "{cluster_id}", state.Primary.Attributes["cluster_id"])

	getLogBackupJobPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	expression := fmt.Sprintf("clusterLogRecord[?jobId=='%s'] | [0]", state.Primary.ID)
	currentTotal := 1
	for {
		currentPath := fmt.Sprintf("%s?limit=10&start=%d", getLogBackupJobPath, currentTotal)
		getLogBackupJobResp, err := client.Request("GET", currentPath, &getLogBackupJobPathOpt)
		if err != nil {
			return getLogBackupJobResp, err
		}
		getLogBackupJobRespBody, err := utils.FlattenResponse(getLogBackupJobResp)
		if err != nil {
			return getLogBackupJobRespBody, err
		}
		logBackupJob := utils.PathSearch(expression, getLogBackupJobRespBody, nil)
		if logBackupJob != nil {
			return logBackupJob, nil
		}
		logBackupJobs := utils.PathSearch("clusterLogRecord",
			getLogBackupJobRespBody, make([]interface{}, 0)).([]interface{})
		if len(logBackupJobs) < 10 {
			break
		}
		currentTotal += len(logBackupJobs)
	}

	return nil, golangsdk.ErrDefault404{}
}

func TestAccManualLogBackup_elastic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_css_manual_log_backup.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getManualLogBackupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testManualLogBackup_elastic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "huaweicloud_css_cluster.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "job_id"),
					resource.TestCheckResourceAttrSet(rName, "type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "log_path"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
		},
	})
}

func TestAccManualLogBackup_logstash(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_css_manual_log_backup.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getManualLogBackupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testManualLogBackup_logstash(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "huaweicloud_css_logstash_cluster.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "job_id"),
					resource.TestCheckResourceAttrSet(rName, "type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "log_path"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
		},
	})
}

func testManualLogBackup_elastic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_css_manual_log_backup" "test" {
  depends_on = [huaweicloud_css_log_setting.test]

  cluster_id = huaweicloud_css_cluster.test.id
}
`, testLogSetting_elastic(name))
}

func testManualLogBackup_logstash(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_css_manual_log_backup" "test" {
  depends_on = [huaweicloud_css_log_setting.test]
  
  cluster_id = huaweicloud_css_logstash_cluster.test.id
}
`, testLogSetting_logstash(name))
}

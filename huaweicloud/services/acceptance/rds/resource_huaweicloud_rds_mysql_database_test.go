package rds

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rds"
)

func getRdsDatabaseFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcRdsV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database from resource id
	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<database_name>")
	}
	instanceId := parts[0]
	dbName := parts[1]
	return rds.QueryDatabases(client, instanceId, dbName)
}

func TestAccRdsDatabase_basic(t *testing.T) {
	var database model.DatabaseForCreation
	rName := acceptance.RandomAccResourceName()
	description := "test database"
	descriptionUpdate := "test database update"
	resourceName := "huaweicloud_rds_mysql_database.test"
	dbPwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&database,
		getRdsDatabaseFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRdsDatabase_basic(rName, dbPwd, description),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "character_set", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testRdsDatabase_basic(rName, dbPwd, descriptionUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "character_set", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdate),
				),
			},
		},
	})
}

func testRdsDatabase_basic(rName, dbPwd, description string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_mysql_database" "test" {
  instance_id   = huaweicloud_rds_instance.test.id
  name          = "%s"
  character_set = "utf8"
  description   = "%s"
}
`, testAccRdsInstance_mysql_step1(rName, dbPwd), rName, description)
}

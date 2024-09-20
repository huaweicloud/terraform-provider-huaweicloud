package dms

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kafka/v2/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDmsKafkaUserFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.HcDmsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	// Split instance_id and user from resource id
	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<user>")
	}
	instanceId := parts[0]
	instanceUser := parts[1]

	// List all instance users
	request := &model.ShowInstanceUsersRequest{
		InstanceId: instanceId,
	}

	response, err := client.ShowInstanceUsers(request)
	if err != nil {
		return nil, fmt.Errorf("error listing DMS kafka users in %s, error: %s", instanceId, err)
	}
	if response.Users != nil && len(*response.Users) != 0 {
		users := *response.Users
		for _, user := range users {
			if *user.UserName == instanceUser {
				return user, nil
			}
		}
	}

	return nil, fmt.Errorf("can not found DMS kafka user")
}

func TestAccDmsKafkaUser_basic(t *testing.T) {
	var user model.ShowInstanceUsersEntity
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_kafka_user.test"
	password := acceptance.RandomPassword()
	passwordUpdate := password + "update"
	description := "add destription"
	descriptionUpdate := ""

	rc := acceptance.InitResourceCheck(
		resourceName,
		&user,
		getDmsKafkaUserFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsKafkaUser_basic(rName, password, description),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "default_app"),
					resource.TestCheckResourceAttrSet(resourceName, "role"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccDmsKafkaUser_basic(rName, passwordUpdate, descriptionUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdate),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccDmsKafkaUser_basic(rName, password string, description string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_user" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = "%[2]s"
  password    = "%[3]s"
  description = "%[4]s"
}
`, testAccKafkaInstance_newFormat(rName), rName, password, description)
}

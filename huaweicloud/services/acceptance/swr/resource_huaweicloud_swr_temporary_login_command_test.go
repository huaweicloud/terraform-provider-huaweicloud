package swr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSwrTemporaryLoginCommand_basic(t *testing.T) {
	rName := "huaweicloud_swr_temporary_login_command.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testSwrTemporaryLoginCommand_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "x_swr_docker_login"),
					resource.TestCheckResourceAttrSet(rName, "auths.#"),
					resource.TestCheckResourceAttrSet(rName, "auths.0.key"),
					resource.TestCheckResourceAttrSet(rName, "auths.0.auth"),
				),
			},
		},
	})
}

const testSwrTemporaryLoginCommand_basic = `resource "huaweicloud_swr_temporary_login_command" "test" {}`

func TestAccSwrTemporaryLoginCommand_enhanced(t *testing.T) {
	rName := "huaweicloud_swr_temporary_login_command.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testSwrTemporaryLoginCommand_enhanced,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "x_swr_docker_login"),
					resource.TestCheckResourceAttrSet(rName, "x_expire_at"),
					resource.TestCheckResourceAttrSet(rName, "auths.#"),
					resource.TestCheckResourceAttrSet(rName, "auths.0.key"),
					resource.TestCheckResourceAttrSet(rName, "auths.0.auth"),
				),
			},
		},
	})
}

const testSwrTemporaryLoginCommand_enhanced = `
resource "huaweicloud_swr_temporary_login_command" "test" {
  enhanced = true
}`

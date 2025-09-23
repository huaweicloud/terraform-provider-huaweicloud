// Copyright (c) 2023 Huawei Technologies Co.,Ltd.
//
// The provider is licensed under the Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//
//     http://license.coscl.org.cn/MulanPSL2
//
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package sfsturbo_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/sfsturbo"
)

func getSFSTurboLdapConfigResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "sfs-turbo"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SFS Turbo client: %s", err)
	}

	return sfsturbo.ReadLdapConfig(client, state.Primary.ID)
}

func TestAccSFSTurboLdapConfig_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_sfs_turbo_ldap_config.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getSFSTurboLdapConfigResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckSFSTurboShareId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testAccSFSTurboLdapConfig_basic(),
				ExpectError: regexp.MustCompile("error waiting for SFS Turbo LDAP configuration creation to complete"),
			},
		},
	})
}

func testAccSFSTurboLdapConfig_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_sfs_turbo_ldap_config" "test" {
  share_id       = "%[1]s"
  url            = "ldap://192.168.0.1:60001"
  base_dn        = "dc=example,dc=com"
  user_dn        = "cn=admin,dc=example,dc=com"
  password       = "password"
  backup_url     = "ldap://192.168.0.2:60002"
  search_timeout = 10
}
`, acceptance.HW_SFS_TURBO_SHARE_ID)
}

package swr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrDomainOverviews_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_domain_overviews.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrDomainOverviews_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "namspace_num"),
					resource.TestCheckResourceAttrSet(dataSource, "repo_num"),
					resource.TestCheckResourceAttrSet(dataSource, "image_num"),
					resource.TestCheckResourceAttrSet(dataSource, "store_space"),
					resource.TestCheckResourceAttrSet(dataSource, "downflow_size"),
					resource.TestCheckResourceAttrSet(dataSource, "domain_id"),
				),
			},
		},
	})
}

const testDataSourceSwrDomainOverviews_basic = `data "huaweicloud_swr_domain_overviews" "test" {}`

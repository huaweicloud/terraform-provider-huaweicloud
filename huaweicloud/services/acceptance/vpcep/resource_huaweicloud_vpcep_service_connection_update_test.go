package vpcep

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

func getVpcEndpointServiceConnectionResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.VPCEPClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPCEP client: %s", err)
	}

	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, want '<service_id>/<endpoint_id>', but got '%s'", state.Primary.ID)
	}
	serviceId := parts[0]
	endpointId := parts[1]

	getHttpUrl := "v1/{project_id}/vpc-endpoint-services/{vpc_endpoint_service_id}/connections?limit=1000&id={endpoint_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{vpc_endpoint_service_id}", serviceId)
	getPath = strings.ReplaceAll(getPath, "{endpoint_id}", endpointId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving VPC endpoint service connection: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening VPC endpoint service connection response: %s", err)
	}

	connection := utils.PathSearch("connections|[0]", getRespBody, nil)
	if connection == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return connection, nil
}

func TestAccVpcEndpointServiceConnectionUpdate_Basic(t *testing.T) {
	var v interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_vpcep_service_connection_update.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&v,
		getVpcEndpointServiceConnectionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndpointServiceConnectionUpdate_Basic(rName, "desc"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "service_id", "huaweicloud_vpcep_service.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "endpoint_id",
						"data.huaweicloud_vpcep_service_connections.test", "connections.0.endpoint_id"),
					resource.TestCheckResourceAttr(resourceName, "description", "desc"),
				),
			},
			{
				Config: testAccVpcEndpointServiceConnectionUpdate_Basic(rName, ""),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				Config: testAccVpcEndpointServiceConnectionUpdate_Basic(rName, "update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "description", "update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccVpcEndpointServiceConnectionUpdate_Base(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpcep_service_connections" "test" {
  service_id = huaweicloud_vpcep_endpoint.test.service_id
}
`, testAccVPCEndpoint_Basic(name))
}

func testAccVpcEndpointServiceConnectionUpdate_Basic(rName, desc string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpcep_service_connection_update" "test" {
  service_id  = huaweicloud_vpcep_service.test.id
  endpoint_id = data.huaweicloud_vpcep_service_connections.test.connections.0.endpoint_id
  description = "%[2]s"
}
`, testAccVpcEndpointServiceConnectionUpdate_Base(rName), desc)
}

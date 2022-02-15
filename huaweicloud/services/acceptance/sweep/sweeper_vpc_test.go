package sweep_test

import (
	"fmt"
	"strings"

	"github.com/chnsz/golangsdk/openstack/networking/v1/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/sweep"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func init() {
	resource.AddTestSweepers("huaweicloud_vpc", &resource.Sweeper{
		Name:         "huaweicloud_vpc",
		F:            sweepVpc,
		Dependencies: []string{"huaweicloud_vpc_subnet"},
	})

	resource.AddTestSweepers("huaweicloud_vpc_subnet", &resource.Sweeper{
		Name: "huaweicloud_vpc_subnet",
		F:    sweepVpcSubnet,
	})
}

func sweepVpc(region string) error {
	sweepResources := make([]*sweep.SweepResource, 0)

	config := acceptance.TestAccProvider.Meta().(*config.Config)
	instances := getResources("huaweicloud_vpc")

	for _, instance := range instances {
		if strings.HasPrefix(instance.Name, "tf_acc_test_") {
			r := vpc.ResourceVirtualPrivateCloudV1()
			d := r.Data(nil)
			d.SetId(instance.Id)
			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, config))
		}
	}

	err := sweep.SweepOrchestrator(sweepResources)

	if err != nil {
		return fmt.Errorf("Error sweeping vpcs (%s): %w", region, err)
	}
	return nil
}

func sweepVpcSubnet(region string) error {
	sweepResources := make([]*sweep.SweepResource, 0)

	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating huaweicloud vpc client: %s", err)
	}

	subnetsArray, err := subnets.List(client, subnets.ListOpts{})
	if err != nil {
		return fmtp.Errorf("Unable to retrieve subnets: %s", err)
	}

	for _, instance := range subnetsArray {
		if strings.HasPrefix(instance.Name, "tf_acc_test_") {
			r := vpc.ResourceVpcSubnetV1()
			d := r.Data(nil)
			d.SetId(instance.ID)
			d.Set("vpc_id", instance.VPC_ID)
			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, config))
		}
	}

	err = sweep.SweepOrchestrator(sweepResources)

	if err != nil {
		return fmt.Errorf("Error sweeping subnets (%s): %w", region, err)
	}

	return nil
}

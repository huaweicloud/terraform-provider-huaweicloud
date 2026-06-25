package deprecated

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/auto_recovery"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func resourceECSAutoRecoveryV1Read(d *schema.ResourceData, meta interface{}, instanceID string) (bool, error) {
	config := meta.(*config.Config)
	client, err := config.ComputeV1Client(config.GetRegion(d))
	if err != nil {
		return false, fmt.Errorf("error creating client: %s", err)
	}

	rId := instanceID

	r, err := auto_recovery.Get(client, rId).Extract()
	if err != nil {
		return false, err
	}
	logp.Printf("[DEBUG] Retrieved ECS-AutoRecovery:%#v of instance:%s", rId, r)
	result, err := strconv.ParseBool(r.SupportAutoRecovery)
	if err != nil {
		log.Printf("[ERROR] error parsing 'SupportAutoRecovery' field to Boolean: %s", err)
	}
	return result, err
}

func setAutoRecoveryForInstance(d *schema.ResourceData, meta interface{}, instanceID string, ar bool) error {
	config := meta.(*config.Config)
	client, err := config.ComputeV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating client: %s", err)
	}

	rId := instanceID

	updateOpts := auto_recovery.UpdateOpts{SupportAutoRecovery: strconv.FormatBool(ar)}

	timeout := d.Timeout(schema.TimeoutUpdate)

	logp.Printf("[DEBUG] Setting ECS-AutoRecovery for instance:%s with options: %#v", rId, updateOpts)
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		err := auto_recovery.Update(client, rId, updateOpts)
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error setting ECS-AutoRecovery for instance%s: %s", rId, err)
	}
	return nil
}

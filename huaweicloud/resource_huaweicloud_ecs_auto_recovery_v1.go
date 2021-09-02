package huaweicloud

import (
	"strconv"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/auto_recovery"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func resourceECSAutoRecoveryV1Read(d *schema.ResourceData, meta interface{}, instanceID string) (bool, error) {
	config := meta.(*config.Config)
	client, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return false, fmtp.Errorf("Error creating HuaweiCloud client: %s", err)
	}

	rId := instanceID

	r, err := auto_recovery.Get(client, rId).Extract()
	if err != nil {
		return false, err
	}
	logp.Printf("[DEBUG] Retrieved ECS-AutoRecovery:%#v of instance:%s", rId, r)
	return strconv.ParseBool(r.SupportAutoRecovery)
}

func setAutoRecoveryForInstance(d *schema.ResourceData, meta interface{}, instanceID string, ar bool) error {
	config := meta.(*config.Config)
	client, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud client: %s", err)
	}

	rId := instanceID

	updateOpts := auto_recovery.UpdateOpts{SupportAutoRecovery: strconv.FormatBool(ar)}

	timeout := d.Timeout(schema.TimeoutUpdate)

	logp.Printf("[DEBUG] Setting ECS-AutoRecovery for instance:%s with options: %#v", rId, updateOpts)
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		err := auto_recovery.Update(client, rId, updateOpts)
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmtp.Errorf("Error setting ECS-AutoRecovery for instance%s: %s", rId, err)
	}
	return nil
}

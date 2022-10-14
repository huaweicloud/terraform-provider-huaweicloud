package huaweicloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/kms/v3/keypairs"
)

func updateEcsInstanceKeyPair(ctx context.Context, d *schema.ResourceData, conf *config.Config) error {
	var taskID string
	serverID := d.Id()
	region := GetRegion(d, conf)
	kmsClient, err := conf.KmsV3Client(region)

	if err != nil {
		return fmt.Errorf("error creating KMS v3 client: %s", err)
	}

	authOpts, err := buildKeypairAuthOpts(d)
	if err != nil {
		return err
	}

	old, new := d.GetChange("key_pair")
	oldKey := old.(string)
	newKey := new.(string)

	if newKey == "" {
		log.Printf("[DEBUG] disassociate the keypair %s of instance %s", oldKey, serverID)
		unbindOpts := keypairs.DisassociateOpts{
			ServerID: serverID,
			Auth:     authOpts,
		}

		if taskID, err = keypairs.Disassociate(kmsClient, unbindOpts); err != nil {
			return fmt.Errorf("error disassociate the keypair %s of instance %s: %s", oldKey, serverID, err)
		}
	} else {
		log.Printf("[DEBUG] associate the keypair %s with instance %s", newKey, serverID)
		bindOpts := keypairs.AssociateOpts{
			Name: newKey,
			Server: keypairs.EcsServerOpts{
				ID:   serverID,
				Auth: authOpts,
			},
		}

		if taskID, err = keypairs.Associate(kmsClient, bindOpts); err != nil {
			return fmt.Errorf("error associate the keypair %s with instance %s: %s", newKey, serverID, err)
		}
	}

	if taskID != "" {
		return waitForKeyPairTaskState(ctx, kmsClient, taskID, d.Timeout(schema.TimeoutUpdate))
	}

	return nil
}

func buildKeypairAuthOpts(d *schema.ResourceData) (*keypairs.AuthOpts, error) {
	if isEcsInstanceShutDown(d) {
		log.Printf("[DEBUG] authentication is not required because the ECS instance %s was shutdown", d.Id())
		return nil, nil
	}

	old, _ := d.GetChange("key_pair")
	oldKey := old.(string)

	// bind a new keypair
	if oldKey == "" {
		passwd := d.Get("admin_pass").(string)
		if passwd == "" {
			return nil, fmt.Errorf("the root password is required when binding a new keypair")
		}

		log.Printf("[DEBUG] the authentication type is password")
		passwdAuth := keypairs.AuthOpts{
			Type: "password",
			Key:  passwd,
		}
		return &passwdAuth, nil
	}

	// replace or unbind an existing keypair
	privateKey := d.Get("private_key").(string)
	if privateKey == "" {
		return nil, fmt.Errorf("the private key of existing keypair must be specified when replacing or unbinding")
	}

	log.Printf("[DEBUG] the authentication type is keypair")
	keypairAuth := keypairs.AuthOpts{
		Type: "keypair",
		Key:  privateKey,
	}
	return &keypairAuth, nil
}

func isEcsInstanceShutDown(d *schema.ResourceData) bool {
	status := d.Get("status").(string)
	return status == "SHUTOFF"
}

func waitForKeyPairTaskState(ctx context.Context, client *golangsdk.ServiceClient, taskID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING", "READY"},
		Target:       []string{"SUCCESS"},
		Refresh:      keyPairTaskStateRefreshFunc(client, taskID),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for keypair task %s to become SUCCESS: %s", taskID, err)
	}
	return nil
}

func keyPairTaskStateRefreshFunc(c *golangsdk.ServiceClient, taskID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		taskInfo, err := keypairs.GetTask(c, taskID).Extract()
		if err != nil {
			return nil, "ERROR", err
		}

		// the format of Status is xxx_mmm, e.g. SUCCESS_BIND, FAILED_UNBIND
		status := strings.Split(taskInfo.Status, "_")[0]
		if status == "FAILED" {
			return taskInfo, status, fmt.Errorf("%s", taskInfo.Status)
		}
		return taskInfo, status, nil
	}
}

package common

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/kms/v3/keypairs"
)

type KeypairAuthOpts struct {
	// the ECS instance ID
	InstanceID string
	// the keypair name in used
	InUsedKeyPair string
	// the replaced keypair name
	NewKeyPair string
	// the private key of the keypair name in used, it's used to replace or unbind the keypair
	InUsedPrivateKey string
	// the root password of the ECS instance, it's used to bind a new keypair
	Password string
	// whether to disable SSH login on the VM
	DisablePassword bool
	// the timeout to wait for the task
	Timeout time.Duration
}

func UpdateEcsInstanceKeyPair(ctx context.Context, ecsClient, kmsClient *golangsdk.ServiceClient, opts *KeypairAuthOpts) error {
	var taskID string

	authOpts, err := buildKeypairAuthOpts(ecsClient, opts)
	if err != nil {
		return err
	}

	instanceID := opts.InstanceID
	if opts.NewKeyPair == "" {
		log.Printf("[DEBUG] disassociate the keypair %s of instance %s", opts.InUsedKeyPair, instanceID)
		unbindOpts := keypairs.DisassociateOpts{
			ServerID: instanceID,
			Auth:     authOpts,
		}

		if taskID, err = keypairs.Disassociate(kmsClient, unbindOpts); err != nil {
			return fmt.Errorf("error disassociate the keypair %s of instance %s: %s", opts.InUsedKeyPair, instanceID, err)
		}
	} else {
		log.Printf("[DEBUG] associate the keypair %s with instance %s", opts.NewKeyPair, instanceID)
		bindOpts := keypairs.AssociateOpts{
			Name: opts.NewKeyPair,
			Server: keypairs.EcsServerOpts{
				ID:              instanceID,
				Auth:            authOpts,
				DisablePassword: &opts.DisablePassword,
			},
		}

		if taskID, err = keypairs.Associate(kmsClient, bindOpts); err != nil {
			return fmt.Errorf("error associate the keypair %s with instance %s: %s", opts.NewKeyPair, instanceID, err)
		}
	}

	if taskID != "" {
		return waitForKeyPairTaskState(ctx, kmsClient, taskID, opts.Timeout)
	}

	return nil
}

func buildKeypairAuthOpts(ecsClient *golangsdk.ServiceClient, opts *KeypairAuthOpts) (*keypairs.AuthOpts, error) {
	instanceID := opts.InstanceID
	server, err := cloudservers.Get(ecsClient, instanceID).Extract()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the ECS instance %s: %s", instanceID, err)
	}

	if server.Status == "SHUTOFF" {
		log.Printf("[DEBUG] authentication is not required because the ECS instance %s was shutdown", instanceID)
		return nil, nil
	}

	// bind a new keypair
	if opts.InUsedKeyPair == "" {
		passwd := opts.Password
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

	// replace or unbind the keypair in used
	privateKey := opts.InUsedPrivateKey
	if privateKey == "" {
		return nil, fmt.Errorf("the private key of keypair in used must be specified when replacing or unbinding")
	}

	log.Printf("[DEBUG] the authentication type is keypair")
	keypairAuth := keypairs.AuthOpts{
		Type: "keypair",
		Key:  privateKey,
	}
	return &keypairAuth, nil
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

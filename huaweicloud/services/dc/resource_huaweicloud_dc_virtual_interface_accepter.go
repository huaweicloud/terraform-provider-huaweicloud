// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DC
// ---------------------------------------------------------------

package dc

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dc/v3/interfaces"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DC GET /v3/{project_id}/dcaas/virtual-interfaces/{interfaceId}
// @API DC PUT /v3/{project_id}/dcaas/virtual-interfaces/{interfaceId}
func ResourceInterfaceAccepter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInterfaceAccepterCreate,
		ReadContext:   resourceInterfaceAccepterRead,
		DeleteContext: resourceInterfaceAccepterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"virtual_interface_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the virtual interface ID.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the action on virtual interfaces shared by other accounts.`,
			},
		},
	}
}

func resourceInterfaceAccepterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		interfaceId = d.Get("virtual_interface_id").(string)
		action      = d.Get("action").(string)
	)
	client, err := cfg.DcV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DC v3 client: %s", err)
	}

	opts := interfaces.UpdateOpts{
		Status: action,
	}
	_, err = interfaces.Update(client, interfaceId, opts)
	if err != nil {
		return diag.Errorf("(%s) DC virtual interface (%s) failed: %s", action, interfaceId, err)
	}
	d.SetId(interfaceId)

	if action == "ACCEPTED" {
		err = waitingForDCVirtualInterfaceActive(ctx, client, d, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for DC virtual interface (%s) to accept activation: %s", interfaceId, err)
		}
	}

	return resourceInterfaceAccepterRead(ctx, d, meta)
}

func waitingForDCVirtualInterfaceActive(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	unexpectedStatus := []string{"ERROR", "DELETED", "REJECTED"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := interfaces.Get(client, d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			if resp == nil {
				return nil, "ERROR", fmt.Errorf("DC virtual interface API response body is empty")
			}

			status := resp.Status
			if status == "ACTIVE" {
				return resp, "COMPLETED", nil
			}

			if status == "" {
				return resp, "ERROR", fmt.Errorf("status is not found in DC virtual interface API response")
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return resp, status, nil
			}
			return resp, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceInterfaceAccepterRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInterfaceAccepterDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

package sdrs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/sdrs/v1/attachreplication"
	"github.com/chnsz/golangsdk/openstack/sdrs/v1/protectedinstances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API SDRS GET /v1/{project_id}/protected-instances/{id}
// @API SDRS POST /v1/{project_id}/protected-instances/{instanceID}/attachreplication
// @API SDRS DELETE /v1/{project_id}/protected-instances/{instanceID}/detachreplication/{replicationID}
// @API SDRS GET /v1/{project_id}/jobs/{job_id}
func ResourceReplicationAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReplicationAttachCreate,
		ReadContext:   resourceReplicationAttachRead,
		DeleteContext: resourceReplicationAttachDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceReplicationAttachImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"replication_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"device": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceReplicationAttachCreate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}
	// use for WaitForJobSuccess
	client.Endpoint = client.ResourceBase

	instanceID := d.Get("instance_id").(string)
	replicationID := d.Get("replication_id").(string)
	attachOpts := attachreplication.CreateOpts{
		ReplicationID: replicationID,
		Device:        d.Get("device").(string),
	}

	n, err := attachreplication.Create(client, instanceID, attachOpts).ExtractJobResponse()
	if err != nil {
		return diag.Errorf("error creating SDRS replication attach: %s", err)
	}

	createTimeoutSec := int(d.Timeout(schema.TimeoutCreate).Seconds())
	if err = attachreplication.WaitForJobSuccess(client, createTimeoutSec, n.JobID); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(formatAttachId(instanceID, replicationID))
	return resourceReplicationAttachRead(ctx, d, meta)
}

func resourceReplicationAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	replicationID := d.Get("replication_id").(string)
	n, err := protectedinstances.Get(client, instanceID).Extract()
	if err != nil {
		if errCode, ok := err.(golangsdk.ErrDefault400); ok {
			if resp, pErr := common.ParseErrorMsg(errCode.Body); pErr == nil && resp.ErrorCode == "SDRS.0208" {
				// `SDRS.0208` means invalid protected instance ID
				return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{},
					"error retrieving SDRS protected instance when query replication attach")
			}
		}
		return diag.FromErr(err)
	}

	attachment, err := flattenReplicationAttach(n, replicationID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SDRS replication attach")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("device", attachment.Device),
		d.Set("replication_id", attachment.Replication),
		d.Set("status", n.Status),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenReplicationAttach(instance *protectedinstances.Instance,
	replicationID string) (*protectedinstances.Attachment, error) {
	for _, attach := range instance.Attachment {
		if attach.Replication == replicationID {
			// find the target attachment
			return &attach, nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}

func resourceReplicationAttachDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}
	// use for WaitForJobSuccess
	client.Endpoint = client.ResourceBase

	instanceID := d.Get("instance_id").(string)
	replicationID := d.Get("replication_id").(string)
	n, err := attachreplication.Delete(client, instanceID, replicationID).ExtractJobResponse()
	if err != nil {
		if errCode, ok := err.(golangsdk.ErrDefault400); ok {
			resp, pErr := common.ParseErrorMsg(errCode.Body)
			if pErr != nil {
				log.Printf("[ERROR] error parse SDRS replication attach delete error, %s", pErr)
				return diag.FromErr(err)
			}

			if resp.ErrorCode == "SDRS.0208" || resp.ErrorCode == "SDRS.0213" {
				// `SDRS.0208` means invalid protected instance ID
				// `SDRS.0213` means invalid replication pair ID
				return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{},
					"error deleting SDRS replication attach")
			}
		}
		return diag.FromErr(err)
	}

	deleteTimeoutSec := int(d.Timeout(schema.TimeoutDelete).Seconds())
	if err := attachreplication.WaitForJobSuccess(client, deleteTimeoutSec, n.JobID); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func formatAttachId(instanceID string, replicationID string) string {
	return fmt.Sprintf("%s/%s", instanceID, replicationID)
}

func extractAttachId(resourceID string) (instanceID, replicationID string, err error) {
	rgs := strings.Split(resourceID, "/")
	if len(rgs) != 2 {
		err = fmt.Errorf("invalid format specified for replication attach id," +
			" must be <protected_instance_id>/<replication_id>")
		return
	}

	instanceID = rgs[0]
	replicationID = rgs[1]
	return
}

func resourceReplicationAttachImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	instanceID, replicationID, err := extractAttachId(d.Id())
	if err != nil {
		return nil, err
	}
	mErr := multierror.Append(
		nil,
		d.Set("instance_id", instanceID),
		d.Set("replication_id", replicationID),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import replication attach id, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}

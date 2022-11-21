package er

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/er/v3/associations"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceAssociation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssociationCreate,
		ReadContext:   resourceAssociationRead,
		DeleteContext: resourceAssociationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAssociationImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the ER instance and route table are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the ER instance to which the route table and the attachment belongs.`,
			},
			"route_table_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the route table to which the association belongs.`,
			},
			"attachment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the attachment corresponding to the association.`,
			},
			"attachment_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the attachment corresponding to the association.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the association.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time.`,
			},
		},
	}
}

func resourceAssociationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.ErV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	var (
		instanceId   = d.Get("instance_id").(string)
		routeTableId = d.Get("route_table_id").(string)

		opts = associations.CreateOpts{
			AttachmentId: d.Get("attachment_id").(string),
		}
	)

	resp, err := associations.Create(client, instanceId, routeTableId, opts)
	if err != nil {
		return diag.Errorf("error creating the association to the route table: %s", err)
	}
	d.SetId(resp.ID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      associationStatusRefreshFunc(client, instanceId, routeTableId, d.Id(), []string{"available"}),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAssociationRead(ctx, d, meta)
}

// QueryAssociationById is a method to query association details from a specified route table using given parameters.
func QueryAssociationById(client *golangsdk.ServiceClient, instanceId, routeTableId,
	associationId string) (*associations.Association, error) {
	resp, err := associations.List(client, instanceId, routeTableId, associations.ListOpts{})
	if err != nil {
		return nil, err
	}

	filter := map[string]interface{}{
		"ID": associationId,
	}
	result, err := utils.FilterSliceWithField(resp, filter)
	if err != nil {
		return nil, err
	}
	if len(resp) < 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte(fmt.Sprintf("the association (%s) does not exist", associationId)),
			},
		}
	}

	association := result[0].(associations.Association)
	log.Printf("[DEBUG] The details of the association (%s) is: %#v", associationId, association)

	return &association, nil
}

func associationStatusRefreshFunc(client *golangsdk.ServiceClient, instanceId, routeTableId, associationId string,
	targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := QueryAssociationById(client, instanceId, routeTableId, associationId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return resp, "COMPLETED", nil
			}

			return nil, "", err
		}

		if utils.StrSliceContains([]string{"failed"}, resp.Status) {
			return resp, "", fmt.Errorf("unexpected status '%s'", resp.Status)
		}
		if utils.StrSliceContains(targets, resp.Status) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

func resourceAssociationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	var (
		instanceId    = d.Get("instance_id").(string)
		routeTableId  = d.Get("route_table_id").(string)
		associationId = d.Id()
	)

	resp, err := QueryAssociationById(client, instanceId, routeTableId, associationId)
	if err != nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "ER association")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("route_table_id", resp.RouteTableId),
		d.Set("attachment_id", resp.AttachmentId),
		d.Set("attachment_type", resp.ResourceType),
		d.Set("status", resp.Status),
		d.Set("created_at", resp.CreatedAt),
		d.Set("updated_at", resp.UpdatedAt),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving association (%s) fields: %s", associationId, mErr)
	}
	return nil
}

func resourceAssociationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	var (
		instanceId    = d.Get("instance_id").(string)
		routeTableId  = d.Get("instance_id").(string)
		associationId = d.Id()

		opts = associations.DeleteOpts{
			AttachmentId: d.Get("attachment_id").(string),
		}
	)
	err = associations.Delete(client, instanceId, routeTableId, opts)
	if err != nil {
		return diag.Errorf("error deleting association (%s): %s", associationId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      associationStatusRefreshFunc(client, instanceId, routeTableId, associationId, nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAssociationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid format for import ID, want '<instance_id>/<route_table_id>/<association_id>', "+
			"but '%s'", d.Id())
	}

	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("route_table_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

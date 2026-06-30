package deprecated

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/dms/v1/groups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceDmsGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsGroupsCreate,
		ReadContext:   resourceDmsGroupsRead,
		DeleteContext: resourceDmsGroupsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		DeprecationMessage: "Deprecated, Distributed Message Service (Shared Edition) has withdrawn, " +
			"please use DMS for Kafka instead.",

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"queue_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"consumed_messages": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"available_messages": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"produced_messages": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"produced_deadletters": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"available_deadletters": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDmsGroupsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	dmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	var getGroups []groups.GroupOps

	n := groups.GroupOps{
		Name: d.Get("name").(string),
	}
	getGroups = append(getGroups, n)

	createOpts := &groups.CreateOps{
		Groups: getGroups,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	v, err := groups.Create(dmsV1Client, d.Get("queue_id").(string), createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating group: %s", err)
	}
	log.Printf("[INFO] group name: %s", v[0].Name)

	// Store the group ID now
	d.SetId(v[0].ID)
	d.Set("queue_id", d.Get("queue_id").(string))

	return resourceDmsGroupsRead(ctx, d, meta)
}

func resourceDmsGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)

	dmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS group client: %s", err)
	}

	queueID := d.Get("queue_id").(string)
	page, err := groups.List(dmsV1Client, queueID, false).AllPages()
	if err != nil {
		return diag.Errorf("error getting groups in queue %s: %s", queueID, err)
	}
	groupsList, err := groups.ExtractGroups(page)
	if len(groupsList) < 1 {
		return diag.FromErr(errors.New("no matching resource found"))
	}

	if len(groupsList) > 1 {
		return diag.FromErr(errors.New("multiple resources matched"))
	}

	group := groupsList[0]
	log.Printf("[DEBUG] Dms group %s: %+v", d.Id(), group)

	d.SetId(group.ID)
	d.Set("name", group.Name)
	d.Set("consumed_messages", group.ConsumedMessages)
	d.Set("available_messages", group.AvailableMessages)
	d.Set("produced_messages", group.ProducedMessages)
	d.Set("produced_deadletters", group.ProducedDeadletters)
	d.Set("available_deadletters", group.AvailableDeadletters)

	return nil
}

func resourceDmsGroupsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	dmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	err = groups.Delete(dmsV1Client, d.Get("queue_id").(string), d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting DMS group: %s", err)
	}

	log.Printf("[DEBUG] DMS group %s has been deleted.", d.Id())
	d.SetId("")
	return nil
}

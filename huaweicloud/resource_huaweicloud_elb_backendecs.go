package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elb"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elb/backendecs"
)

const nameELBBackend = "ELB-BackendECS"

func resourceELBBackendECS() *schema.Resource {
	return &schema.Resource{
		Create: resourceELBBackendECSCreate,
		Read:   resourceELBBackendECSRead,
		Delete: resourceELBBackendECSDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"listener_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"server_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"private_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"public_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"update_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"create_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"health_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"server_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"listeners": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceELBBackendECSCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := chooseELBClient(d, config)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	var createOpts backendecs.CreateOpts
	err, _ = buildCreateParam(&createOpts, d)
	if err != nil {
		return fmt.Errorf("Error creating %s: building parameter failed:%s", nameELBBackend, err)
	}
	log.Printf("[DEBUG] Create %s Options: %#v", nameELBBackend, createOpts)

	lId := d.Get("listener_id").(string)
	j, err := backendecs.Create(networkingClient, createOpts, lId).Extract()
	if err != nil {
		return fmt.Errorf("Error creating %s: %s", nameELBBackend, err)
	}
	log.Printf("[DEBUG] Create %s, the job is: %#v", nameELBBackend, *j)

	// Wait for BackendECS to become active before continuing
	timeout := d.Timeout(schema.TimeoutCreate)
	jobInfo, err := waitForELBJobSuccess(networkingClient, j, timeout)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Create %s, the job info is: %#v", nameELBBackend, jobInfo)

	e1, ok := jobInfo.Entities["members"]
	if !ok {
		return fmt.Errorf("Error creating %s: get the entity from job info failed", nameELBBackend)
	}
	e, ok := e1.([]interface{})
	if !ok {
		return fmt.Errorf("Error creating %s: convert job entity to array failed", nameELBBackend)
	}
	if len(e) != 1 {
		return fmt.Errorf("Error creating %s: the number of member does not equal 1", nameELBBackend)
	}
	i, ok := e[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Error creating %s: convert job entity to map failed", nameELBBackend)
	}
	eid, ok := i["id"]
	if !ok {
		return fmt.Errorf("Error creating %s: get backend id from job entity failed", nameELBBackend)
	}

	// If all has been successful, set the ID on the resource
	d.SetId(eid.(string))

	return resourceELBBackendECSRead(d, meta)
}

func resourceELBBackendECSRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := chooseELBClient(d, config)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	lId := d.Get("listener_id").(string)
	b, err := backendecs.Get(networkingClient, lId, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "backendecs")
	}
	log.Printf("[DEBUG] Retrieved %s(%s): %#v", nameELBBackend, d.Id(), b)

	return refreshResourceData(b, d)
}

func resourceELBBackendECSDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := chooseELBClient(d, config)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	bId := d.Id()
	deleteOpts := backendecs.DeleteOpts{
		RemoveMember: []backendecs.RemoveMemberField{backendecs.RemoveMemberField{ID: bId}},
	}
	log.Printf("[DEBUG] Deleting %s option: %#v", nameELBBackend, deleteOpts)

	var job *elb.Job
	timeout := d.Timeout(schema.TimeoutDelete)
	lId := d.Get("listener_id").(string)
	err = resource.Retry(timeout, func() *resource.RetryError {
		j, err := backendecs.Delete(networkingClient, lId, deleteOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		job = j
		return nil
	})
	if err != nil {
		if isResourceNotFound(err) {
			log.Printf("[INFO] deleting an unavailable %s: %s", nameELBBackend, bId)
			return nil
		}
		return fmt.Errorf("Error deleting %s(%s): %s", nameELBBackend, bId, err)
	}
	log.Printf("[DEBUG] Delete %s, the job is: %#v", nameELBBackend, *job)

	_, err = waitForELBJobSuccess(networkingClient, job, timeout)
	if err != nil {
		return err
	}
	return nil
}

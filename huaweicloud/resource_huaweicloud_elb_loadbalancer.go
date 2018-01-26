package huaweicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elb"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elb/loadbalancers"
)

const nameELBLB = "ELB-LoadBalancer"

func resourceELBLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceELBLoadBalancerCreate,
		Read:   resourceELBLoadBalancerRead,
		Update: resourceELBLoadBalancerUpdate,
		Delete: resourceELBLoadBalancerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					vv := regexp.MustCompile("^[a-zA-Z0-9-_]{1,64}$")
					if !vv.MatchString(value) {
						errors = append(errors, fmt.Errorf("%s is a string of 1 to 64 characters that consist of letters, digits, underscores (_), and hyphens (-)", k))
					}
					return
				},
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					vv := regexp.MustCompile("^[^<>]{0,128}$")
					if !vv.MatchString(value) {
						errors = append(errors, fmt.Errorf("%s is a string of 0 to 128 characters and cannot contain angle brackets (<>)", k))
					}
					return
				},
			},

			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"bandwidth": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 || value > 300 {
						errors = append(errors, fmt.Errorf("%s must be in [1, 300]", k))
					}
					return
				},
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "Internal" && value != "External" {
						errors = append(errors, fmt.Errorf("%s must be Internal or External", k))
					}
					return
				},
			},

			"admin_state_up": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"vip_subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"az": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"charge_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "traffic" {
						errors = append(errors, fmt.Errorf("%s must be traffic", k))
					}
					return
				},
			},

			"eip_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					vv := regexp.MustCompile("^[a-zA-Z0-9-]{1,200}$")
					if !vv.MatchString(value) {
						errors = append(errors, fmt.Errorf("%s is a string of 1 to 200 characters that consists of uppercase and lowercase letters, digits, and hyphens (-)", k))
					}
					return
				},
			},

			"vip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tenantid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
		},
	}
}

func resourceELBLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := chooseELBClient(d, config)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	var opts loadbalancers.CreateOpts
	err, _ = buildCreateParam(&opts, d)
	if err != nil {
		return fmt.Errorf("Error creating %s: building parameter failed:%s", nameELBLB, err)
	}
	log.Printf("[DEBUG] Create %s Options: %#v", nameELBLB, opts)

	switch {
	case opts.Type == "External" && !hasFilledParam(d, "bandwidth"):
		return fmt.Errorf("bandwidth is mandatory when type is set to External")

	case opts.Type == "Internal" && !hasFilledParam(d, "vip_subnet_id"):
		return fmt.Errorf("vip_subnet_id is mandatory when type is set to Internal")

	case opts.Type == "Internal" && !hasFilledParam(d, "az"):
		return fmt.Errorf("az is mandatory when type is set to Internal")

	case opts.Type == "Internal" && !hasFilledParam(d, "security_group_id"):
		return fmt.Errorf("security_group_id is mandatory when type is set to Internal")

	case opts.Type == "Internal" && !hasFilledParam(d, "tenantid"):
		return fmt.Errorf("tenantid is mandatory when type is set to Internal")
	}

	j, err := loadbalancers.Create(networkingClient, opts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating %s: %s", nameELBLB, err)
	}
	log.Printf("[DEBUG] Create %s, the job is: %#v", nameELBLB, *j)

	// Wait for LoadBalancer to become active before continuing
	timeout := d.Timeout(schema.TimeoutCreate)
	jobInfo, err := waitForELBJobSuccess(networkingClient, j, timeout)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Create %s, the job is: %#v", nameELBLB, jobInfo)

	e, ok := jobInfo.Entities["elb"]
	if !ok {
		return fmt.Errorf("Error creating %s: get the entity from job info failed", nameELBLB)
	}
	i, ok := e.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Error creating %s: convert job entity to map failed", nameELBLB)
	}
	eid, ok := i["id"]
	if !ok {
		return fmt.Errorf("Error creating %s: get elb id from job entity failed", nameELBLB)
	}

	/*
		// Once the loadbalancer has been created, apply any requested security groups
		// to the port that was created behind the scenes.
		if err := resourceELBLoadBalancerSecurityGroups(networkingClient, lb.VipPortID, d); err != nil {
			return err
		}*/

	// If all has been successful, set the ID on the resource
	d.SetId(eid.(string))

	return resourceELBLoadBalancerRead(d, meta)
}

func resourceELBLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := chooseELBClient(d, config)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	lb, err := loadbalancers.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "loadbalancer")
	}
	log.Printf("[DEBUG] Retrieved %s %s: %#v", nameELBLB, d.Id(), lb)

	return refreshResourceData(lb, d)
}

func resourceELBLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := chooseELBClient(d, config)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	lbId := d.Id()

	var updateOpts loadbalancers.UpdateOpts
	err, _ = buildUpdateParam(&updateOpts, d)
	if err != nil {
		return fmt.Errorf("Error updating %s %s: building parameter failed:%s", nameELBLB, lbId, err)
	}
	b, err := updateOpts.IsNeedUpdate()
	if err != nil {
		return err
	}
	if !b {
		log.Printf("[INFO] Updating %s %s with no changes", nameELBLB, lbId)
		return nil
	}

	// Wait for LoadBalancer to become active before continuing
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForELBLoadBalancerActive(networkingClient, lbId, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating %s %s with options: %#v", nameELBLB, lbId, updateOpts)
	var job *elb.Job
	err = resource.Retry(timeout, func() *resource.RetryError {
		j, err := loadbalancers.Update(networkingClient, lbId, updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		job = j
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error updating %s %s: %s", nameELBLB, lbId, err)
	}

	// Wait for LoadBalancer to become active before continuing
	_, err = waitForELBJobSuccess(networkingClient, job, timeout)
	if err != nil {
		return err
	}

	return resourceELBLoadBalancerRead(d, meta)
}

func resourceELBLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := chooseELBClient(d, config)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	lbId := d.Id()
	log.Printf("[DEBUG] Deleting %s %s", nameELBLB, lbId)

	var job *elb.Job
	timeout := d.Timeout(schema.TimeoutDelete)
	err = resource.Retry(timeout, func() *resource.RetryError {
		j, err := loadbalancers.Delete(networkingClient, lbId).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		job = j
		return nil
	})
	if err != nil {
		if isResourceNotFound(err) {
			log.Printf("[INFO] deleting an unavailable %s: %s", nameELBLB, lbId)
			return nil
		}
		return fmt.Errorf("Error deleting %s %s: %s", nameELBLB, lbId, err)
	}
	log.Printf("[DEBUG] Delete %s, the job is: %#v", nameELBLB, *job)

	_, err = waitForELBJobSuccess(networkingClient, job, timeout)
	if err != nil {
		return err
	}
	return nil
}

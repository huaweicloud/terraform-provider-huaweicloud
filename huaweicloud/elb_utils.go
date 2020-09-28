package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/elb"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/elb/listeners"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/elb/loadbalancers"
)

func waitForELBJobSuccess(networkingClient *golangsdk.ServiceClient, j *elb.Job, timeout time.Duration) (*elb.JobInfo, error) {
	jobId := j.JobId
	target := "SUCCESS"
	pending := []string{"INIT", "RUNNING"}

	log.Printf("[DEBUG] Waiting for elb job %s to become %s.", jobId, target)

	ji, err := waitForELBResource(networkingClient, "job", jobId, target, pending, timeout, getELBJobInfo)
	if err == nil {
		return ji.(*elb.JobInfo), nil
	}
	return nil, err
}

func getELBJobInfo(networkingClient *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		info, err := elb.QueryJobInfo(networkingClient, jobId).Extract()
		if err != nil {
			return nil, "", err
		}

		return info, info.Status, nil
	}
}

func waitForELBLoadBalancerActive(networkingClient *golangsdk.ServiceClient, id string, timeout time.Duration) error {
	target := "ACTIVE"

	log.Printf("[DEBUG] Waiting for elb %s to become %s.", id, target)

	_, err := waitForELBResource(networkingClient, "loadbalancer", id, target, []string{"PENDING_CREATE"}, timeout, getELBLoadBalancer)
	return err
}

func getELBLoadBalancer(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		lb, err := loadbalancers.Get(networkingClient, id).Extract()
		if err != nil {
			return nil, "", err
		}
		return lb, lb.Status, nil
	}
}

func waitForELBListenerActive(networkingClient *golangsdk.ServiceClient, id string, timeout time.Duration) error {
	target := "ACTIVE"

	log.Printf("[DEBUG] Waiting for elb-listener %s to become %s.", id, target)

	_, err := waitForELBResource(networkingClient, "listener", id, target, []string{"PENDING_CREATE"}, timeout, getELBListener)
	return err
}

func getELBListener(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		l, err := listeners.Get(networkingClient, id).Extract()
		if err != nil {
			return nil, "", err
		}
		return l, l.Status, nil
	}
}

type getELBResource func(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc

func waitForELBResource(networkingClient *golangsdk.ServiceClient, name string, id string, target string, pending []string, timeout time.Duration, f getELBResource) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    f(networkingClient, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	o, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return nil, fmt.Errorf("Error: elb %s %s not found: %s", name, id, err)
		}
		return nil, fmt.Errorf("Error waiting for elb %s %s to become %s: %s", name, id, target, err)
	}

	return o, nil
}

func isResourceNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(golangsdk.ErrDefault404)
	return ok
}

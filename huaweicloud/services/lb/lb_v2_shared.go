package lb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/elb/v2/l7policies"
	"github.com/chnsz/golangsdk/openstack/elb/v2/listeners"
	"github.com/chnsz/golangsdk/openstack/elb/v2/loadbalancers"
	"github.com/chnsz/golangsdk/openstack/elb/v2/monitors"
	"github.com/chnsz/golangsdk/openstack/elb/v2/pools"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// lbPendingStatuses are the valid statuses a LoadBalancer will be in while
// it's updating.
var lbPendingStatuses = []string{"PENDING_CREATE", "PENDING_UPDATE"}

// lbPendingDeleteStatuses are the valid statuses a LoadBalancer will be before delete
var lbPendingDeleteStatuses = []string{"ERROR", "PENDING_UPDATE", "PENDING_DELETE", "ACTIVE"}

var lbSkipLBStatuses = []string{"ERROR", "ACTIVE"}

func waitForLBV2Listener(ctx context.Context, networkingClient *golangsdk.ServiceClient, id string, target string,
	pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for listener %s to become %s.", id, target)

	stateConf := &retry.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2ListenerRefreshFunc(networkingClient, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("error: listener %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("error waiting for listener %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceLBV2ListenerRefreshFunc(networkingClient *golangsdk.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		listener, err := listeners.Get(networkingClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		// The listener resource has no Status attribute, so a successful Get is the best we can do
		return listener, "ACTIVE", nil
	}
}

func waitForLBV2LoadBalancer(ctx context.Context, networkingClient *golangsdk.ServiceClient,
	id string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for loadbalancer %s to become %s", id, target)

	stateConf := &retry.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2LoadBalancerRefreshFunc(networkingClient, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("error: loadbalancer %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("error waiting for loadbalancer %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceLBV2LoadBalancerRefreshFunc(networkingClient *golangsdk.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		lb, err := loadbalancers.Get(networkingClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		return lb, lb.ProvisioningStatus, nil
	}
}

func waitForLBV2Pool(ctx context.Context, networkingClient *golangsdk.ServiceClient, id string, target string,
	pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for pool %s to become %s.", id, target)

	stateConf := &retry.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2PoolRefreshFunc(networkingClient, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("error: pool %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("error waiting for pool %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceLBV2PoolRefreshFunc(networkingClient *golangsdk.ServiceClient, poolID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pool, err := pools.Get(networkingClient, poolID).Extract()
		if err != nil {
			return nil, "", err
		}

		// The pool resource has no Status attribute, so a successful Get is the best we can do
		return pool, "ACTIVE", nil
	}
}

func waitForLBV2viaPool(ctx context.Context, networkingClient *golangsdk.ServiceClient, id string, target string,
	timeout time.Duration) error {
	pool, err := pools.Get(networkingClient, id).Extract()
	if err != nil {
		return err
	}

	if len(pool.Loadbalancers) > 0 {
		// each pool has an LB in Octavia lbaasv2 API
		lbID := pool.Loadbalancers[0].ID
		return waitForLBV2LoadBalancer(ctx, networkingClient, lbID, target, nil, timeout)
	}

	if len(pool.Listeners) > 0 {
		// each pool has a listener in Neutron lbaasv2 API
		listenerID := pool.Listeners[0].ID
		listener, err := listeners.Get(networkingClient, listenerID).Extract()
		if err != nil {
			return err
		}
		if listener.Loadbalancers != nil {
			lbID := listener.Loadbalancers[0].ID
			return waitForLBV2LoadBalancer(ctx, networkingClient, lbID, target, nil, timeout)
		}
	}

	// got a pool but no LB - this is wrong
	return fmt.Errorf("no Load Balancer on pool %s", id)
}

// nolint:gocyclo
func resourceLBV2LoadBalancerStatusRefreshFuncNeutron(lbClient *golangsdk.ServiceClient, lbID, resourceType,
	resourceID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		statuses, err := loadbalancers.GetStatuses(lbClient, lbID).Extract()
		if err != nil {
			return nil, "", fmt.Errorf("unable to get statuses from the Load Balancer %s statuses tree: %s", lbID, err)
		}

		if !utils.StrSliceContains(lbSkipLBStatuses, statuses.Loadbalancer.ProvisioningStatus) {
			return statuses.Loadbalancer, statuses.Loadbalancer.ProvisioningStatus, nil
		}

		switch resourceType {
		case "listener":
			for _, listener := range statuses.Loadbalancer.Listeners {
				if listener.ID == resourceID {
					if listener.ProvisioningStatus != "" {
						return listener, listener.ProvisioningStatus, nil
					}
				}
			}
			listener, err := listeners.Get(lbClient, resourceID).Extract()
			return listener, "ACTIVE", err

		case "pool":
			for _, pool := range statuses.Loadbalancer.Pools {
				if pool.ID == resourceID {
					if pool.ProvisioningStatus != "" {
						return pool, pool.ProvisioningStatus, nil
					}
				}
			}
			pool, err := pools.Get(lbClient, resourceID).Extract()
			return pool, "ACTIVE", err

		case "monitor":
			for _, pool := range statuses.Loadbalancer.Pools {
				if pool.Monitor.ID == resourceID {
					if pool.Monitor.ProvisioningStatus != "" {
						return pool.Monitor, pool.Monitor.ProvisioningStatus, nil
					}
				}
			}
			monitor, err := monitors.Get(lbClient, resourceID).Extract()
			return monitor, "ACTIVE", err

		case "member":
			for _, pool := range statuses.Loadbalancer.Pools {
				for _, member := range pool.Members {
					if member.ID == resourceID {
						if member.ProvisioningStatus != "" {
							return member, member.ProvisioningStatus, nil
						}
					}
				}
			}
			return "", "DELETED", nil

		case "l7policy":
			for _, listener := range statuses.Loadbalancer.Listeners {
				for _, l7policy := range listener.L7Policies {
					if l7policy.ID == resourceID {
						if l7policy.ProvisioningStatus != "" {
							return l7policy, l7policy.ProvisioningStatus, nil
						}
					}
				}
			}
			l7policy, err := l7policies.Get(lbClient, resourceID).Extract()
			return l7policy, "ACTIVE", err

		case "l7rule":
			for _, listener := range statuses.Loadbalancer.Listeners {
				for _, l7policy := range listener.L7Policies {
					for _, l7rule := range l7policy.Rules {
						if l7rule.ID == resourceID {
							if l7rule.ProvisioningStatus != "" {
								return l7rule, l7rule.ProvisioningStatus, nil
							}
						}
					}
				}
			}
			return "", "DELETED", nil
		}

		return nil, "", fmt.Errorf("an unexpected error occurred querying the status of %s %s by loadbalancer %s",
			resourceType, resourceID, lbID)
	}
}

func resourceLBV2L7PolicyRefreshFunc(lbClient *golangsdk.ServiceClient, lbID string,
	l7policy *l7policies.L7Policy) retry.StateRefreshFunc {
	if l7policy.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBV2LoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !utils.StrSliceContains(lbSkipLBStatuses, status) {
				return lb, status, nil
			}

			l7policy, err := l7policies.Get(lbClient, l7policy.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return l7policy, l7policy.ProvisioningStatus, nil
		}
	}

	return resourceLBV2LoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "l7policy", l7policy.ID)
}

// nolint:revive
func waitForLBV2L7Policy(ctx context.Context, lbClient *golangsdk.ServiceClient, parentListener *listeners.Listener,
	l7policy *l7policies.L7Policy, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for l7policy %s to become %s.", l7policy.ID, target)

	if len(parentListener.Loadbalancers) == 0 {
		return fmt.Errorf("unable to determine loadbalancer ID from listener %s", parentListener.ID)
	}

	lbID := parentListener.Loadbalancers[0].ID

	stateConf := &retry.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2L7PolicyRefreshFunc(lbClient, lbID, l7policy),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("error waiting for l7policy %s to become %s: %s", l7policy.ID, target, err)
	}

	return nil
}

func getListenerIDForL7Policy(lbClient *golangsdk.ServiceClient, id string) (string, error) {
	log.Printf("[DEBUG] Trying to get Listener ID associated with the %s L7 Policy ID", id)
	lbsPages, err := loadbalancers.List(lbClient, loadbalancers.ListOpts{}).AllPages()
	if err != nil {
		return "", fmt.Errorf("no Load Balancers were found: %s", err)
	}

	lbs, err := loadbalancers.ExtractLoadBalancers(lbsPages)
	if err != nil {
		return "", fmt.Errorf("unable to extract Load Balancers list: %s", err)
	}

	for _, lb := range lbs {
		statuses, err := loadbalancers.GetStatuses(lbClient, lb.ID).Extract()
		if err != nil {
			return "", fmt.Errorf("failed to get Load Balancer statuses: %s", err)
		}
		for _, listener := range statuses.Loadbalancer.Listeners {
			for _, l7policy := range listener.L7Policies {
				if l7policy.ID == id {
					return listener.ID, nil
				}
			}
		}
	}

	return "", fmt.Errorf("unable to find Listener ID associated with the %s L7 Policy ID", id)
}

func resourceLBV2L7RuleRefreshFunc(lbClient *golangsdk.ServiceClient, lbID string, l7policyID string,
	l7rule *l7policies.Rule) retry.StateRefreshFunc {
	if l7rule.ProvisioningStatus != "" {
		return func() (interface{}, string, error) {
			lb, status, err := resourceLBV2LoadBalancerRefreshFunc(lbClient, lbID)()
			if err != nil {
				return lb, status, err
			}
			if !utils.StrSliceContains(lbSkipLBStatuses, status) {
				return lb, status, nil
			}

			l7rule, err := l7policies.GetRule(lbClient, l7policyID, l7rule.ID).Extract()
			if err != nil {
				return nil, "", err
			}

			return l7rule, l7rule.ProvisioningStatus, nil
		}
	}

	return resourceLBV2LoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "l7rule", l7rule.ID)
}

// nolint:revive
func waitForLBV2L7Rule(ctx context.Context, lbClient *golangsdk.ServiceClient, parentListener *listeners.Listener,
	parentL7policy *l7policies.L7Policy, l7rule *l7policies.Rule, target string, pending []string,
	timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for l7rule %s to become %s.", l7rule.ID, target)

	if len(parentListener.Loadbalancers) == 0 {
		return fmt.Errorf("unable to determine loadbalancer ID from listener %s", parentListener.ID)
	}

	lbID := parentListener.Loadbalancers[0].ID

	stateConf := &retry.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2L7RuleRefreshFunc(lbClient, lbID, parentL7policy.ID, l7rule),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("error waiting for l7rule %s to become %s: %s", l7rule.ID, target, err)
	}

	return nil
}

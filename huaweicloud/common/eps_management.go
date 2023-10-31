package common

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/eps/v1/enterpriseprojects"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// GetEnterpriseProjectID returns the enterprise_project_id that was specified in the resource.
// If it was not set, the provider-level value is checked. The provider-level value can
// either be set by the `enterprise_project_id` argument or by HW_ENTERPRISE_PROJECT_ID.
func GetEnterpriseProjectID(d *schema.ResourceData, cfg *config.Config) string {
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		return v.(string)
	}

	return cfg.EnterpriseProjectID
}

// MigrateEnterpriseProjectWithoutWait is a method that used to a migrate resource from an enterprise project to
// another.
// NOTE: Please read the following contents carefully before using this method.
//   - This method only sends an asynchronous request and does not guarantee the result.
func MigrateEnterpriseProjectWithoutWait(cfg *config.Config, d *schema.ResourceData,
	opts enterpriseprojects.MigrateResourceOpts) error {
	targetEpsId := cfg.GetEnterpriseProjectID(d)
	if targetEpsId == "" {
		targetEpsId = "0"
	}

	client, err := cfg.EnterpriseProjectClient(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating EPS client: %s", err)
	}
	_, err = enterpriseprojects.Migrate(client, opts, targetEpsId).Extract()
	if err != nil {
		return fmt.Errorf("failed to migrate resource (%s) to the enterprise project (%s): %s",
			opts.ResourceId, targetEpsId, err)
	}
	return nil
}

// MigrateEnterpriseProject is a method used to migrate a resource from an enterprise project to another enterprise
// project and ensure the success of the EPS side migration.
// NOTE: Please read the following contents carefully before using this method.
//   - This method only calls the interfaces of the EPS service. For individual EPS IDs that are not updated due to
//     out-of-synchronization of data on the server side, this method does not perform additional verification and
//     requires developers to manually ensure the reliability of the code through testing.
func MigrateEnterpriseProject(ctx context.Context, cfg *config.Config, d *schema.ResourceData,
	opts enterpriseprojects.MigrateResourceOpts) error {
	targetEpsId := cfg.GetEnterpriseProjectID(d)
	if targetEpsId == "" {
		targetEpsId = "0"
	}

	client, err := cfg.EnterpriseProjectClient(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating EPS client: %s", err)
	}
	_, err = enterpriseprojects.Migrate(client, opts, targetEpsId).Extract()
	if err != nil {
		return fmt.Errorf("failed to migrate resource (%s) to the enterprise project (%s): %s",
			opts.ResourceId, targetEpsId, err)
	}

	// Wait for the Enterprise Project ID changed.
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			s, err := getAssociatedResourceById(client, opts.ProjectId, targetEpsId, opts.ResourceType, opts.ResourceId)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return nil, "PENDING", nil
				}
				return nil, "ERROR", err
			}
			return s, "COMPLETED", nil
		},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for migrating enterprise poject complete: %s", err)
	}
	return nil
}

func getAssociatedResourceById(client *golangsdk.ServiceClient, projectId, epsId, resourceType,
	resourceId string) (*enterpriseprojects.Resource, error) {
	opts := enterpriseprojects.ListResourcesOpts{
		EnterpriseProjectId: epsId,
		Projects:            []string{projectId},
		ResourceTypes:       []string{resourceType},
	}
	resourceList, err := enterpriseprojects.ListAssociatedResources(client, opts)
	if err != nil {
		return nil, err
	}
	for _, success := range resourceList.Resources {
		if success.ResourceId == resourceId {
			return &success, nil
		}
	}
	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte(fmt.Sprintf("unable to find the resource under a specified enterprise_project_id (%s)", epsId)),
		},
	}
}

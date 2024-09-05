package config

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type MigrateResourceOpts struct {
	// The type of the resource to be migrated.
	ResourceType string `json:"resource_type" required:"true"`
	// The ID of the resource to be migrated.
	ResourceId string `json:"resource_id" required:"true"`
	// Project ID. This is a required option when resource_type is a region-level service.
	// It can be left blank or filled with an empty string when resource_type is a global-level service.
	ProjectId string `json:"project_id,omitempty"`
	// Region ID. Required if resource type is 'bucket'.
	RegionId string `json:"region_id,omitempty"`
	// Whether migration is associate resource.
	Associated bool `json:"associated,omitempty"`
}

// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrate
func migrateResourceToAnotherEps(client *golangsdk.ServiceClient, targetEpsId string, requestBody map[string]interface{}) error {
	httpUrl := "v1.0/enterprise-projects/{enterprise_project_id}/resources-migrate"
	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{enterprise_project_id}", targetEpsId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: requestBody,
		OkCodes:  []int{204},
	}

	_, err := client.Request("POST", actionPath, &opt)
	return err
}

// MigrateEnterpriseProjectWithoutWait is a method that used to a migrate resource from an enterprise project to
// another.
// By default, the resource will be migrated to the default enterprise project.
// NOTE: Please read the following contents carefully before using this method.
//   - This method only sends an asynchronous request and does not guarantee the result.
func (cfg *Config) MigrateEnterpriseProjectWithoutWait(d *schema.ResourceData, opts MigrateResourceOpts) error {
	targetEpsId := cfg.GetEnterpriseProjectID(d, "0")

	requestBody, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return fmt.Errorf("error building request body: %s", err)
	}
	client, err := cfg.NewServiceClient("eps", cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating EPS client: %s", err)
	}
	err = migrateResourceToAnotherEps(client, targetEpsId, requestBody)
	if err != nil {
		return fmt.Errorf("failed to migrate resource (%v) to the enterprise project (%s): %s", opts.ResourceId, targetEpsId, err)
	}
	return nil
}

// MigrateEnterpriseProject is a method used to migrate a resource from an enterprise project to another enterprise
// project and ensure the success of the EPS side migration.
// By default, the resource will be migrated to the default enterprise project.
// NOTE: Please read the following contents carefully before using this method.
//   - This method only calls the interfaces of the EPS service. For individual EPS IDs that are not updated due to
//     out-of-synchronization of data on the server side, this method does not perform additional verification and
//     requires developers to manually ensure the reliability of the code through testing.
func (cfg *Config) MigrateEnterpriseProject(ctx context.Context, d *schema.ResourceData, opts MigrateResourceOpts) error {
	targetEpsId := cfg.GetEnterpriseProjectID(d, "0")

	requestBody, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return fmt.Errorf("error building request body: %s", err)
	}
	client, err := cfg.NewServiceClient("eps", cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating EPS client: %s", err)
	}
	err = migrateResourceToAnotherEps(client, targetEpsId, requestBody)
	if err != nil {
		return fmt.Errorf("failed to migrate resource (%v) to the enterprise project (%s): %s", opts.ResourceId, targetEpsId, err)
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

func buildListResourcesUnderEpsOpts(projectId, resourceType string, offset int) map[string]interface{} {
	result := map[string]interface{}{
		"limit":  1000,
		"offset": offset,
	}
	if projectId != "" {
		result["projects"] = []string{projectId}
	}
	if resourceType != "" {
		result["resource_types"] = []string{resourceType}
	}
	return result
}

func queryResourcesUnderEps(client *golangsdk.ServiceClient, epsId, projectId, resourceType string) ([]interface{}, error) {
	var (
		httpUrl = "v1.0/enterprise-projects/{enterprise_project_id}/resources/filter"
		offset  = 0
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{enterprise_project_id}", epsId)

	for {
		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			JSONBody: buildListResourcesUnderEpsOpts(projectId, resourceType, offset),
		}
		requestResp, err := client.Request("POST", listPath, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		resources := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
		if len(resources) < 1 {
			break
		}
		result = append(result, resources...)
		offset += len(resources)
	}

	return result, nil
}

// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources/filter
func getAssociatedResourceById(client *golangsdk.ServiceClient, projectId, epsId, resourceType,
	resourceId string) (interface{}, error) {
	resources, err := queryResourcesUnderEps(client, epsId, projectId, resourceType)
	if err != nil {
		return nil, err
	}
	if resInfo := utils.PathSearch(fmt.Sprintf("[?resource_id=='%s']|[0]", resourceId), resources, nil); resInfo != nil {
		return resInfo, nil
	}
	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte(fmt.Sprintf("unable to find the resource under a specified enterprise_project_id (%s)", epsId)),
		},
	}
}

package sweep

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SweepResource struct {
	d        *schema.ResourceData
	meta     interface{}
	resource *schema.Resource
}

func NewSweepResource(resource *schema.Resource, d *schema.ResourceData, meta interface{}) *SweepResource {
	return &SweepResource{
		d:        d,
		meta:     meta,
		resource: resource,
	}
}

func SweepOrchestrator(sweepResources []*SweepResource) error {
	return SweepOrchestratorContext(context.Background(), sweepResources)
}

func SweepOrchestratorContext(ctx context.Context, sweepResources []*SweepResource) error {
	var g multierror.Group

	for _, sweepResource := range sweepResources {
		sweepResource := sweepResource
		g.Go(func() error {
			return DeleteResource(sweepResource.resource, sweepResource.d, sweepResource.meta)
		})
	}

	return g.Wait().ErrorOrNil()
}

func DeleteResource(resource *schema.Resource, d *schema.ResourceData, meta interface{}) error {
	if resource.DeleteContext != nil || resource.DeleteWithoutTimeout != nil {
		var diags diag.Diagnostics

		if resource.DeleteContext != nil {
			diags = resource.DeleteContext(context.Background(), d, meta)
		} else {
			diags = resource.DeleteWithoutTimeout(context.Background(), d, meta)
		}

		for i := range diags {
			if diags[i].Severity == diag.Error {
				return fmt.Errorf("error deleting resource: %s", diags[i].Summary)
			}
		}

		return nil
	}

	return resource.Delete(d, meta)
}

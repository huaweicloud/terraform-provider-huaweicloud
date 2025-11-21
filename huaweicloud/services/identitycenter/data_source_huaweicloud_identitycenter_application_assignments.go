package identitycenter

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IDENTITYCENTER GET /v1/instances/{instance_id}/application-assignments-for-principals
func DataSourceIdentityCenterApplicationAssignments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCenterApplicationAssignmentsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"principal_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"principal_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_assignments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type ApplicationAssignmentsForPrincipalDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newApplicationAssignmentsForPrincipalDSWrapper(d *schema.ResourceData, meta interface{}) *ApplicationAssignmentsForPrincipalDSWrapper {
	return &ApplicationAssignmentsForPrincipalDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceIdentityCenterApplicationAssignmentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newApplicationAssignmentsForPrincipalDSWrapper(d, meta)
	assignments, err := wrapper.listApplicationAssignmentsForPrincipal()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listApplicationAssignmentsForPrincipalToSchema(assignments)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func (w *ApplicationAssignmentsForPrincipalDSWrapper) listApplicationAssignmentsForPrincipal() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "identitycenter")
	if err != nil {
		return nil, err
	}

	uri := "/v1/instances/{instance_id}/application-assignments-for-principals"
	uri = strings.ReplaceAll(uri, "{instance_id}", w.Get("instance_id").(string))
	params := map[string]any{
		"principal_id":   w.Get("principal_id").(string),
		"principal_type": w.Get("principal_type").(string),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		MarkerPager("application_assignments", "page_info.next_marker", "marker").
		Request().
		Result()
}

func (w *ApplicationAssignmentsForPrincipalDSWrapper) listApplicationAssignmentsForPrincipalToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("application_assignments", schemas.SliceToList(body.Get("application_assignments"),
			func(assignments gjson.Result) any {
				return map[string]any{
					"application_urn": assignments.Get("application_urn").Value(),
					"principal_id":    assignments.Get("principal_id").Value(),
					"principal_type":  assignments.Get("principal_type").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

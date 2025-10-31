package identitycenter

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var identityCenterApplicationAssignmentNonUpdateParams = []string{"instance_id", "application_instance_id", "principal_id", "principal_type"}

// @API IdentityCenter POST /v1/instances/{instance_id}/applications/{application_instance_id}/assignments/create
// @API IdentityCenter GET /v1/instances/{instance_id}/applications/{application_instance_id}/assignments
// @API IdentityCenter POST /v1/instances/{instance_id}/applications/{application_instance_id}/assignments/delete
func ResourceIdentityCenterApplicationAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterApplicationAssignmentCreate,
		UpdateContext: resourceIdentityCenterApplicationAssignmentUpdate,
		ReadContext:   resourceIdentityCenterApplicationAssignmentRead,
		DeleteContext: resourceIdentityCenterApplicationAssignmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterApplicationAssignmentImport,
		},

		CustomizeDiff: config.FlexibleForceNew(identityCenterApplicationAssignmentNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_instance_id": {
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
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"application_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityCenterApplicationAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createHttpUrl = "v1/instances/{instance_id}/applications/{application_instance_id}/assignments/create"
		createProduct = "identitycenter"
	)

	client, err := cfg.NewServiceClient(createProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))
	createPath = strings.ReplaceAll(createPath, "{application_instance_id}", fmt.Sprintf("%v", d.Get("application_instance_id")))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateApplicationAssignmentBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center application asssignment: %s", err)
	}

	instanceId := d.Get("instance_id")
	applicationInstanceId := d.Get("application_instance_id")
	principalId := d.Get("principal_id")
	d.SetId(fmt.Sprintf("%v/%v/%v", instanceId, applicationInstanceId, principalId))

	return resourceIdentityCenterApplicationAssignmentRead(ctx, d, meta)
}

func buildCreateApplicationAssignmentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"principal_id":   d.Get("principal_id"),
		"principal_type": d.Get("principal_type"),
	}
	return bodyParams
}

func resourceIdentityCenterApplicationAssignmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listHttpUrl = "v1/instances/{instance_id}/applications/{application_instance_id}/assignments"
		listProduct = "identitycenter"
	)

	client, err := cfg.NewServiceClient(listProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	listBasePath := client.Endpoint + listHttpUrl
	listBasePath = strings.ReplaceAll(listBasePath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))
	listBasePath = strings.ReplaceAll(listBasePath, "{application_instance_id}", fmt.Sprintf("%v", d.Get("application_instance_id")))

	listPath := listBasePath + buildListApplicationAssignmentsQueryParams("")

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var assignment interface{}
	principalId := d.Get("principal_id")

	for {
		listResp, err := client.Request("GET", listPath, &listOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving Identity Center application assignments.")
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		assignment = utils.PathSearch(fmt.Sprintf("application_assignments[?principal_id=='%s']|[0]", principalId), listRespBody, nil)
		if assignment != nil {
			break
		}

		marker := utils.PathSearch("page_info.next_marker", listRespBody, nil)
		if marker == nil {
			break
		}
		listPath = listBasePath + buildListApplicationAssignmentsQueryParams(marker.(string))
	}

	if assignment == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no application assignment found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("principal_id", utils.PathSearch("principal_id", assignment, nil)),
		d.Set("principal_type", utils.PathSearch("principal_type", assignment, nil)),
		d.Set("application_urn", utils.PathSearch("application_urn", assignment, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListApplicationAssignmentsQueryParams(marker string) string {
	res := "?limit=100"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}

func resourceIdentityCenterApplicationAssignmentUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterApplicationAssignmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteHttpUrl = "v1/instances/{instance_id}/applications/{application_instance_id}/assignments/delete"
		deleteProduct = "identitycenter"
	)
	client, err := cfg.NewServiceClient(deleteProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))
	deletePath = strings.ReplaceAll(deletePath, "{application_instance_id}", fmt.Sprintf("%v", d.Get("application_instance_id")))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateApplicationAssignmentBodyParams(d)),
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center application assignment: %s", err)
	}

	return nil
}

func resourceIdentityCenterApplicationAssignmentImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, errors.New("invalid id format, must be <instance_id>/<application_instance_id>/<principal_id>")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("application_instance_id", parts[1]),
		d.Set("principal_id", parts[2]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}

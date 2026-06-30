package drs

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS GET /v5/{project_id}/jobs/{job_id}/object/support
func DataSourceDrsJobObjectSupport() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJobObjectSupportRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_full_trans_support_object": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_incre_trans_support_object": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_full_incre_trans_support_object": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"support_object_import_engine": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_support_column_mapping": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_database_support_search": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_schema_support_search": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_table_support_search": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"file_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"previous_select": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"import_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_import_cloumn": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"import_mapping_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_import_unique_key": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceJobObjectSupportRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs/{job_id}/object/support"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("is_full_trans_support_object", utils.PathSearch("is_full_trans_support_object", respBody, false)),
		d.Set("is_incre_trans_support_object", utils.PathSearch("is_incre_trans_support_object", respBody, false)),
		d.Set("is_full_incre_trans_support_object", utils.PathSearch("is_full_incre_trans_support_object", respBody, false)),
		d.Set("support_object_import_engine", utils.PathSearch("support_object_import_engine", respBody, nil)),
		d.Set("is_support_column_mapping", utils.PathSearch("is_support_column_mapping", respBody, false)),
		d.Set("is_database_support_search", utils.PathSearch("is_database_support_search", respBody, false)),
		d.Set("is_schema_support_search", utils.PathSearch("is_schema_support_search", respBody, false)),
		d.Set("is_table_support_search", utils.PathSearch("is_table_support_search", respBody, false)),
		d.Set("file_size", utils.PathSearch("file_size", respBody, nil)),
		d.Set("previous_select", utils.PathSearch("previous_select", respBody, nil)),
		d.Set("import_level", utils.PathSearch("import_level", respBody, nil)),
		d.Set("is_import_cloumn", utils.PathSearch("is_import_cloumn", respBody, false)),
		d.Set("import_mapping_type", utils.PathSearch("import_mapping_type", respBody, nil)),
		d.Set("is_import_unique_key", utils.PathSearch("is_import_unique_key", respBody, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

package drs

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var pwdBatchModifyNonUpdatableParams = []string{
	"jobs",
	"jobs.*.job_id",
	"jobs.*.db_password",
	"jobs.*.end_point_type",
	"jobs.*.kerberos",
	"jobs.*.kerberos.*.krb5_conf_file",
	"jobs.*.kerberos.*.key_tab_file",
	"jobs.*.kerberos.*.domain_name",
	"jobs.*.kerberos.*.user_principal",
}

// @API DRS PUT /v3/{project_id}/jobs/batch-modify-pwd
func ResourcePwdBatchModify() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePwdBatchModifyCreate,
		ReadContext:   resourcePwdBatchModifyRead,
		UpdateContext: resourcePwdBatchModifyUpdate,
		DeleteContext: resourcePwdBatchModifyDelete,

		CustomizeDiff: config.FlexibleForceNew(pwdBatchModifyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Elem:     pwdBatchModifyJobSchema(),
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_point_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func pwdBatchModifyJobSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"end_point_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kerberos": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"krb5_conf_file": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key_tab_file": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_principal": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func buildPwdBatchModifyKerberosParams(kerberosRaw interface{}) map[string]interface{} {
	if kerberosRaw == nil {
		return nil
	}

	kerberosList, ok := kerberosRaw.([]interface{})
	if !ok || len(kerberosList) == 0 {
		return nil
	}

	kerberos, ok := kerberosList[0].(map[string]interface{})
	if !ok {
		return nil
	}
	kerberosMap := make(map[string]interface{})

	if v, ok := kerberos["krb5_conf_file"]; ok && v != "" {
		kerberosMap["krb5_conf_file"] = v
	}
	if v, ok := kerberos["key_tab_file"]; ok && v != "" {
		kerberosMap["key_tab_file"] = v
	}
	if v, ok := kerberos["domain_name"]; ok && v != "" {
		kerberosMap["domain_name"] = v
	}
	if v, ok := kerberos["user_principal"]; ok && v != "" {
		kerberosMap["user_principal"] = v
	}

	if len(kerberosMap) > 0 {
		return kerberosMap
	}

	return nil
}

func buildPwdBatchModifyBodyParams(d *schema.ResourceData) map[string]interface{} {
	jobsRaw := d.Get("jobs").([]interface{})
	jobs := make([]map[string]interface{}, 0, len(jobsRaw))

	for _, jobRaw := range jobsRaw {
		job, ok := jobRaw.(map[string]interface{})
		if !ok {
			continue
		}

		jobMap := map[string]interface{}{
			"job_id":         job["job_id"],
			"db_password":    job["db_password"],
			"end_point_type": job["end_point_type"],
			"kerberos":       buildPwdBatchModifyKerberosParams(job["kerberos"]),
		}

		jobs = append(jobs, jobMap)
	}

	return map[string]interface{}{
		"jobs": jobs,
	}
}

func resourcePwdBatchModifyCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/jobs/batch-modify-pwd"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildPwdBatchModifyBodyParams(d),
	}

	resp, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch modifying DRS jobs password: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	results := utils.PathSearch("results", respBody, make([]interface{}, 0)).([]interface{})
	if len(results) == 0 {
		return diag.Errorf("unable to find the results from the API response")
	}

	resourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId.String())

	mErr := multierror.Append(nil,
		d.Set("results", flattenPwdBatchModifyResults(results)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting DRS batch modify password fields: %s", mErr)
	}

	return nil
}

func flattenPwdBatchModifyResults(results []interface{}) []interface{} {
	if len(results) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(results))
	for _, result := range results {
		rst = append(rst, map[string]interface{}{
			"id":             utils.PathSearch("id", result, nil),
			"status":         utils.PathSearch("status", result, nil),
			"end_point_type": utils.PathSearch("end_point_type", result, nil),
		})
	}
	return rst
}

func resourcePwdBatchModifyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourcePwdBatchModifyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourcePwdBatchModifyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch modify DRS jobs password. Deleting this 
resource will not restore the modified password or undo the modify action, but will only remove the resource information 
from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

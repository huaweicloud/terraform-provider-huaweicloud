package dbss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DBSS POST /v2/{project_id}/{instance_id}/audit/databases/rds
// @API DBSS GET /v1/{project_id}/{instance_id}/dbss/audit/databases
// @API DBSS POST /v2/{project_id}/{instance_id}/audit/databases/switch
// @API DBSS DELETE /v2/{project_id}/{instance_id}/audit/databases/{db_id}
func ResourceAddRdsDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAddRdsDatabaseCreate,
		ReadContext:   resourceAddRdsDatabaseRead,
		UpdateContext: resourceAddRdsDatabaseUpdate,
		DeleteContext: resourceAddRdsDatabaseDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAddRdsDatabaseImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rds_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lts_audit_switch": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"db_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charset": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"audit_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"agent_url": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"db_classification": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rds_audit_switch_mismatch": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceAddRdsDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	httpUrl := "v2/{project_id}/{instance_id}/audit/databases/rds"
	client, err := cfg.NewServiceClient("dbss", region)
	if err != nil {
		return diag.Errorf("error creating DBSS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAddRdsDatabaseParams(d),
	}
	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error adding RDS database to the DBSS instance: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId := utils.PathSearch("ret_list[0].id", respBody, "").(string)
	if resourceId == "" {
		return diag.Errorf("error adding RDS database to the DBSS instance: ID is not found in API response")
	}

	d.SetId(resourceId)

	status := d.Get("status").(string)
	if status == "ON" {
		databaseInfo, err := GetDatabaseList(client, instanceId, d.Id())
		if err != nil {
			return diag.Errorf("error retrieving DBSS audit databases in creation")
		}
		databaseId := utils.PathSearch("database.id", databaseInfo, "").(string)
		if databaseId == "" {
			return diag.Errorf("failed to enable the audit status: 'db_id' is not found in API response")
		}

		err = updateAuditStatus(client, d, instanceId, databaseId)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceAddRdsDatabaseRead(ctx, d, meta)
}

func buildAddRdsDatabaseParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"databases": []map[string]interface{}{
			{
				"id":   d.Get("rds_id"),
				"type": d.Get("type"),
			},
		},
	}

	return params
}

func resourceAddRdsDatabaseRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dbss", region)
	if err != nil {
		return diag.Errorf("error creating DBSS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	databaseInfo, err := GetDatabaseList(client, instanceId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DBSS audit databases")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("rds_id", utils.PathSearch("database.rds_id", databaseInfo, nil)),
		d.Set("db_id", utils.PathSearch("database.id", databaseInfo, nil)),
		d.Set("type", utils.PathSearch("database.type", databaseInfo, nil)),
		d.Set("status", utils.PathSearch("database.status", databaseInfo, nil)),
		d.Set("name", utils.PathSearch("database.name", databaseInfo, nil)),
		d.Set("version", utils.PathSearch("database.version", databaseInfo, nil)),
		d.Set("charset", utils.PathSearch("database.charset", databaseInfo, nil)),
		d.Set("ip", utils.PathSearch("database.ip", databaseInfo, nil)),
		d.Set("port", utils.PathSearch("database.port", databaseInfo, nil)),
		d.Set("os", utils.PathSearch("database.os", databaseInfo, nil)),
		d.Set("instance_name", utils.PathSearch("database.instance_name", databaseInfo, nil)),
		d.Set("audit_status", utils.PathSearch("database.audit_status", databaseInfo, nil)),
		d.Set("agent_url", utils.PathSearch("database.agent_url", databaseInfo, nil)),
		d.Set("db_classification", utils.PathSearch("database.db_classification", databaseInfo, nil)),
		d.Set("rds_audit_switch_mismatch", utils.PathSearch("database.rds_audit_switch_mismatch", databaseInfo, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetDatabaseList(client *golangsdk.ServiceClient, instanceId, rdsId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/{instance_id}/dbss/audit/databases?limit=100"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	offset := 0
	for {
		getPathWithOffset := fmt.Sprintf("%s&offset=%d", getPath, offset)
		resp, err := client.Request("GET", getPathWithOffset, &getOpts)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		databaseList := utils.PathSearch("databases", respBody, make([]interface{}, 0)).([]interface{})
		if len(databaseList) == 0 {
			break
		}

		database := utils.PathSearch(fmt.Sprintf("[?database.rds_id=='%s']|[0]", rdsId), databaseList, nil)
		if database != nil {
			return database, nil
		}
		offset += len(databaseList)
	}

	return nil, golangsdk.ErrDefault404{}
}

func updateAuditStatus(client *golangsdk.ServiceClient, d *schema.ResourceData, instanceId, databaseId string) error {
	httpUrl := "v2/{project_id}/{instance_id}/audit/databases/switch"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", instanceId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateaAuditStatusBodyParams(d, databaseId),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error updating the database audit status: %s", err)
	}

	return nil
}

func buildUpdateaAuditStatusBodyParams(d *schema.ResourceData, databaseId string) map[string]interface{} {
	params := map[string]interface{}{
		"id":               databaseId,
		"status":           d.Get("status"),
		"lts_audit_switch": utils.ValueIgnoreEmpty(d.Get("lts_audit_switch")),
	}

	return params
}

func resourceAddRdsDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dbss", region)
	if err != nil {
		return diag.Errorf("error creating DBSS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	databaseId := d.Get("db_id").(string)

	if d.HasChanges("status", "lts_audit_switch") {
		if databaseId == "" {
			diag.Errorf("edit audit status is not currently supported: 'db_id' is not found in API response")
		}
		err := updateAuditStatus(client, d, instanceId, databaseId)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceAddRdsDatabaseRead(ctx, d, meta)
}

func resourceAddRdsDatabaseDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dbss", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DBSS client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	databaseId := d.Get("db_id").(string)

	httpUrl := "v2/{project_id}/{instance_id}/audit/databases/{db_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{db_id}", databaseId)

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error removing the RDS database from the DBSS instance")
	}
	return nil
}

func resourceAddRdsDatabaseImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s'",
			importedId)
	}
	d.SetId(parts[1])

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

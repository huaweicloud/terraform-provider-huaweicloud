package cbr

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CBR GET /v3/{project_id}/backups/{backup_id}/metadata
func DataSourceBackupMetadata() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBackupMetadataRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backups": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"floatingips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"interface": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"server": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceBackupMetadataRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v3/{project_id}/backups/{backup_id}/metadata"
		product  = "cbr"
		backupId = d.Get("backup_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{backup_id}", backupId)

	resp, err := client.Request("GET", requestPath, &golangsdk.RequestOpts{
		KeepResponseBody: true,
	})
	if err != nil {
		return diag.Errorf("error querying CBR backup metadata: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("backup_id", backupId),
		d.Set("backups", utils.PathSearch("backups", respBody, nil)),
		d.Set("flavor", utils.PathSearch("flavor", respBody, nil)),
		d.Set("floatingips", utils.PathSearch("floatingips", respBody, nil)),
		d.Set("interface", utils.PathSearch("interface", respBody, nil)),
		d.Set("ports", utils.PathSearch("ports", respBody, nil)),
		d.Set("server", utils.PathSearch("server", respBody, nil)),
		d.Set("volumes", utils.PathSearch("volumes", respBody, nil)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting data source fields of the CBR backup metadata: %s", mErr)
	}
	return nil
}

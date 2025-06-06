package evs

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

// @API EVS GET /v3/{project_id}/snapshots/{snapshot_id}/metadata
func DataSourceSnapshotMetadata() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSnapshotMetadataRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceSnapshotMetadataRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		mErr       *multierror.Error
		product    = "evs"
		getHttpUrl = "v3/{project_id}/snapshots/{snapshot_id}/metadata"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{snapshot_id}", d.Get("snapshot_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving EVS snapshot metadata: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("metadata", utils.ExpandToStringMap(utils.PathSearch("metadata", getRespBody, make(map[string]interface{})).(map[string]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

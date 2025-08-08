package ddm

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDM GET /v1/{project_id}/instances/{instance_id}/nodes/{node_id}
func DataSourceDdmInstanceNodeDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDdmInstanceNodeDetailRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"floating_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datavolume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"res_subnet_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"systemvolume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDdmInstanceNodeDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/{project_id}/instances/{instance_id}/nodes/{node_id}"
		product = "ddm"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))
	getPath = strings.ReplaceAll(getPath, "{node_id}", fmt.Sprintf("%v", d.Get("node_id")))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DDM instance node detail")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("private_ip", utils.PathSearch("private_ip", getRespBody, nil)),
		d.Set("floating_ip", utils.PathSearch("floating_ip", getRespBody, nil)),
		d.Set("server_id", utils.PathSearch("server_id", getRespBody, nil)),
		d.Set("subnet_name", utils.PathSearch("subnet_name", getRespBody, nil)),
		d.Set("datavolume_id", utils.PathSearch("datavolume_id", getRespBody, nil)),
		d.Set("res_subnet_ip", utils.PathSearch("res_subnet_ip", getRespBody, nil)),
		d.Set("systemvolume_id", utils.PathSearch("systemvolume_id", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

package codeartsdeploy

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API CodeArtsDeploy POST /v1/resources/host-groups/{group_id}/hosts/replication
func ResourceCodeArtsDeployHostsCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCodeArtsDeployHostsCopyCreate,
		ReadContext:   resourceCodeArtsDeployHostsCopyRead,
		DeleteContext: resourceCodeArtsDeployHostsCopyDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the source group ID.`,
			},
			"host_uuids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the host ID list.`,
			},
			"target_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the target group ID.`,
			},
		},
	}
}

func resourceCodeArtsDeployHostsCopyCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	httpUrl := "v1/resources/host-groups/{group_id}/hosts/replication"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{group_id}", d.Get("source_group_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateCodeArtsDeployHostsCopyBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error copying CodeArts deploy hosts: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	return nil
}

func buildCreateCodeArtsDeployHostsCopyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"host_uuids":      d.Get("host_uuids"),
		"target_group_id": d.Get("target_group_id"),
	}
	return bodyParams
}

func resourceCodeArtsDeployHostsCopyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCodeArtsDeployHostsCopyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting hosts copy resource is not supported. The resource is only removed from the" +
		"state, the hosts remain in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

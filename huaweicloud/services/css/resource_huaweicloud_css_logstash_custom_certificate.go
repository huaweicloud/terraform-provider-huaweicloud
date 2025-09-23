package css

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

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/certs/upload
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/certs/{cert_id}
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/certs
// @API CSS DELETE /v1.0/{project_id}/clusters/{cluster_id}/certs/{cert_id}/delete
func ResourceLogstashCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogstashCertificateCreate,
		ReadContext:   resourceLogstashCertificateRead,
		DeleteContext: resourceLogstashCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLogstashCertificateImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cert_object": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLogstashCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterID := d.Get("cluster_id").(string)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	createCertificateHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/certs/upload"
	createCertificatePath := cssV1Client.Endpoint + createCertificateHttpUrl
	createCertificatePath = strings.ReplaceAll(createCertificatePath, "{project_id}", cssV1Client.ProjectID)
	createCertificatePath = strings.ReplaceAll(createCertificatePath, "{cluster_id}", clusterID)

	createCertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createCertificateOpt.JSONBody = map[string]interface{}{
		"bucket_name":  d.Get("bucket_name").(string),
		"certs_object": d.Get("cert_object").(string),
	}
	_, err = cssV1Client.Request("POST", createCertificatePath, &createCertificateOpt)
	if err != nil {
		return diag.Errorf("error upload CSS logstash cluster custom certificate: %s", err)
	}

	certificate, err := getCertificateByName(d, cssV1Client)
	if err != nil {
		return diag.FromErr(err)
	}
	id := utils.PathSearch("id", certificate, "").(string)
	d.SetId(id)

	return resourceLogstashCertificateRead(ctx, d, meta)
}

func resourceLogstashCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	getCertificateHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/certs/{cert_id}"
	getCertificatePath := cssV1Client.Endpoint + getCertificateHttpUrl
	getCertificatePath = strings.ReplaceAll(getCertificatePath, "{project_id}", cssV1Client.ProjectID)
	getCertificatePath = strings.ReplaceAll(getCertificatePath, "{cluster_id}", d.Get("cluster_id").(string))
	getCertificatePath = strings.ReplaceAll(getCertificatePath, "{cert_id}", d.Id())

	getCertificatePathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getCertificateResp, err := cssV1Client.Request("GET", getCertificatePath, &getCertificatePathOpt)
	if err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		// {"errCode": "CSS.0015","externalMessage": "No resources are found or the access is denied."}
		err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
		return common.CheckDeletedDiag(d, err, "error querying CSS logstash cluster custom certificate")
	}

	getCertificateRespBody, err := utils.FlattenResponse(getCertificateResp)
	if err != nil {
		return diag.Errorf("erorr retrieving CSS logstash cluster custom certificate: %s", err)
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("fileName", getCertificateRespBody, nil)),
		d.Set("path", utils.PathSearch("fileLocation", getCertificateRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getCertificateRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updateAt", getCertificateRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceLogstashCertificateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterID := d.Get("cluster_id").(string)

	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	deleteCertificateUrl := "v1.0/{project_id}/clusters/{cluster_id}/certs/{cert_id}/delete"
	deleteCertificatePath := cssV1Client.Endpoint + deleteCertificateUrl
	deleteCertificatePath = strings.ReplaceAll(deleteCertificatePath, "{project_id}", cssV1Client.ProjectID)
	deleteCertificatePath = strings.ReplaceAll(deleteCertificatePath, "{cluster_id}", clusterID)
	deleteCertificatePath = strings.ReplaceAll(deleteCertificatePath, "{cert_id}", d.Id())

	deleteCertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = cssV1Client.Request("DELETE", deleteCertificatePath, &deleteCertificateOpt)
	if err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		// {"errCode": "CSS.0015","externalMessage": "No resources are found or the access is denied."}
		err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
		return common.CheckDeletedDiag(d, err, "error deleting CSS logstash cluster custom certificate")
	}

	return nil
}

func getCertificateByName(d *schema.ResourceData, cssV1Client *golangsdk.ServiceClient) (interface{}, error) {
	getCertificateHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/certs"
	getCertificatePath := cssV1Client.Endpoint + getCertificateHttpUrl
	getCertificatePath = strings.ReplaceAll(getCertificatePath, "{project_id}", cssV1Client.ProjectID)
	getCertificatePath = strings.ReplaceAll(getCertificatePath, "{cluster_id}", d.Get("cluster_id").(string))

	getCertificatePathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	currentTotal := 1
	for {
		currentPath := fmt.Sprintf("%s?limit=10&start=%d", getCertificatePath, currentTotal)
		getCertificateResp, err := cssV1Client.Request("GET", currentPath, &getCertificatePathOpt)
		if err != nil {
			return getCertificateResp, err
		}
		getCertificateRespBody, err := utils.FlattenResponse(getCertificateResp)
		if err != nil {
			return nil, err
		}
		certificates := utils.PathSearch("certsRecords", getCertificateRespBody, make([]interface{}, 0)).([]interface{})
		parts := strings.Split(d.Get("cert_object").(string), "/")
		findCertListStr := fmt.Sprintf("certsRecords|[?fileName=='%s']|[0]", parts[len(parts)-1])
		certificate := utils.PathSearch(findCertListStr, getCertificateRespBody, nil)
		if certificate != nil {
			return certificate, nil
		}
		total := utils.PathSearch("totalSize", getCertificateRespBody, float64(0)).(float64)
		currentTotal += len(certificates)
		if float64(currentTotal-1) == total {
			break
		}
	}
	return nil, golangsdk.ErrDefault404{}
}

func resourceLogstashCertificateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, must be <cluster_id>/<id>")
	}

	d.SetId(parts[1])

	mErr := multierror.Append(nil,
		d.Set("cluster_id", parts[0]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

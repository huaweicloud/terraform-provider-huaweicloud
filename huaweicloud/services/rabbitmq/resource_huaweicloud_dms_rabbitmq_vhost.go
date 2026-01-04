package rabbitmq

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

// @API RabbitMQ PUT /v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts
// @API RabbitMQ POST /v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts
// @API RabbitMQ GET /v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts
func ResourceDmsRabbitmqVhost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRabbitmqVhostCreate,
		ReadContext:   resourceDmsRabbitmqVhostRead,
		DeleteContext: resourceDmsRabbitmqVhostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVhostImportState,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tracing": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceDmsRabbitmqVhostCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	name := d.Get("name").(string)

	createHttpUrl := "v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"name": name,
		},
	}
	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating vhost: %s", err)
	}

	id := fmt.Sprintf("%s/%s", instanceID, name)
	d.SetId(id)

	return resourceDmsRabbitmqVhostRead(ctx, d, cfg)
}

func resourceDmsRabbitmqVhostRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	name := d.Get("name").(string)

	vhost, err := GetRabbitmqVhost(client, instanceID, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving the vhost")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tracing", utils.PathSearch("tracing", vhost, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetRabbitmqVhost(client *golangsdk.ServiceClient, instanceID, name string) (interface{}, error) {
	listHttpUrl := "v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// pageLimit is `10`
	listPath += fmt.Sprintf("?limit=%d", pageLimit)
	offset := 0
	for {
		currentPath := listPath + fmt.Sprintf("&offset=%d", offset)
		listResp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return nil, err
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return nil, err
		}

		searchPath := fmt.Sprintf("items[?name=='%s']|[0]", name)
		result := utils.PathSearch(searchPath, listRespBody, nil)
		if result != nil {
			return result, nil
		}

		// `total` means the number of all `vhost`, and type is float64.
		offset += pageLimit
		total := utils.PathSearch("total", listRespBody, float64(0))
		if int(total.(float64)) <= offset {
			return nil, golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Method:    "GET",
					URL:       "/v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts",
					RequestId: "NONE",
					Body:      []byte(fmt.Sprintf("the vhost (%s) does not exist", name)),
				},
			}
		}
	}
}

func resourceDmsRabbitmqVhostDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	name := d.Get("name").(string)

	deleteHttpUrl := "v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceID)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		JSONBody: map[string]interface{}{
			"name": []string{name},
		},
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting vhost")
	}

	return nil
}

func resourceVhostImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ",")
	if len(parts) != 2 {
		parts = strings.Split(d.Id(), "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid ID format, must be <instance_id>/<name> or <instance_id>,<name>")
		}
	} else {
		// reform ID to be separated by a slash
		id := fmt.Sprintf("%s/%s", parts[0], parts[1])
		d.SetId(id)
	}

	d.Set("instance_id", parts[0])
	d.Set("name", parts[1])

	return []*schema.ResourceData{d}, nil
}

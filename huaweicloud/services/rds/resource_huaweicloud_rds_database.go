package rds

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	v3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3"
	rds "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceRdsDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsDatabaseCreate,
		UpdateContext: resourceRdsDatabaseUpdate,
		DeleteContext: resourceRdsDatabaseDelete,
		ReadContext:   resourceRdsDatabaseRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
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
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[\$a-z0-9-_]+$`),
						"the name can only consist of lowercase letters, digits, hyphens (-), underscores (_) and dollar signs ($)"),
					func(v interface{}, k string) (ws []string, errors []error) {
						re := regexp.MustCompile(`-|\$`)
						if len(re.FindAllString(v.(string), -1)) > 10 {
							errors = append(errors,
								fmt.Errorf("the total number of hyphens (-) and dollar signs ($) cannot exceed 10"))
						}
						return
					},
					validation.StringLenBetween(1, 64),
				),
			},
			"character_set": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 512),
			},
		},
	}
}

func resourceRdsDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcRdsV3Client(c.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	config.MutexKV.Lock(instanceId)
	defer config.MutexKV.Unlock(instanceId)

	dbName := d.Get("name").(string)
	createOpts := rds.DatabaseForCreation{
		Name:         dbName,
		CharacterSet: d.Get("character_set").(string),
		Comment:      utils.StringIgnoreEmpty(d.Get("description").(string)),
	}

	createDatabaseReq := rds.CreateDatabaseRequest{
		InstanceId: instanceId,
		Body:       &createOpts,
	}

	_, err = client.CreateDatabase(&createDatabaseReq)
	if err != nil {
		return diag.Errorf("error creating RDS database: %s", err)
	}

	id := instanceId + "/" + dbName
	d.SetId(id)
	return resourceRdsDatabaseRead(ctx, d, meta)
}

func resourceRdsDatabaseRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcRdsV3Client(c.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database from resource id
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<database_name>")
	}
	instanceId := parts[0]
	dbName := parts[1]

	db, err := QueryDatabases(client, instanceId, dbName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error listing RDS db databases")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", instanceId),
		d.Set("name", dbName),
		d.Set("character_set", db.CharacterSet),
		d.Set("description", db.Comment),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting RDS db database fields: %s", err)
	}

	return nil
}

func resourceRdsDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcRdsV3Client(c.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	config.MutexKV.Lock(instanceId)
	defer config.MutexKV.Unlock(instanceId)

	updateOpts := rds.UpdateDatabaseRequest{
		InstanceId: instanceId,
		Body: &rds.UpdateDatabaseReq{
			Name:    d.Get("name").(string),
			Comment: d.Get("description").(string),
		},
	}

	log.Printf("[DEBUG] Update RDS database options: %#v", updateOpts)
	_, err = client.UpdateDatabase(&updateOpts)
	if err != nil {
		return diag.Errorf("error updating RDS database: %s", err)
	}

	return resourceRdsDatabaseRead(ctx, d, meta)
}

func resourceRdsDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcRdsV3Client(c.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	config.MutexKV.Lock(instanceId)
	defer config.MutexKV.Unlock(instanceId)

	deleteOpts := rds.DeleteDatabaseRequest{
		InstanceId: instanceId,
		DbName:     d.Get("name").(string),
	}

	log.Printf("[DEBUG] Delete RDS database options: %#v", deleteOpts)
	_, err = client.DeleteDatabase(&deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting RDS database: %s", err)
	}

	return nil
}

func QueryDatabases(client *v3.RdsClient, instanceId, dbName string) (*rds.DatabaseForCreation, error) {
	request := rds.ListDatabasesRequest{
		InstanceId: instanceId,
		Limit:      int32(100),
		Page:       int32(1),
	}

	// List all databases
	for {
		response, err := client.ListDatabases(&request)
		if err != nil {
			return nil, err
		}
		if response.Databases == nil || len(*response.Databases) == 0 {
			break
		}

		databases := *response.Databases
		request.Page += 1
		for _, db := range databases {
			if db.Name == dbName {
				return &db, nil
			}
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

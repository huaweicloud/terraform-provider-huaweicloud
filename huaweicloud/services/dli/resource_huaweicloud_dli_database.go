package dli

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dli/v1/databases"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DLI POST /v1.0/{project_id}/databases
// @API DLI GET /v1.0/{project_id}/databases
// @API DLI PUT /v1.0/{project_id}/databases/{database_name}/owner
// @API DLI DELETE /v1.0/{project_id}/databases/{database_name}
func ResourceDliSqlDatabaseV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDliSQLDatabaseCreate,
		ReadContext:   resourceDliSQLDatabaseRead,
		UpdateContext: resourceDliSQLDatabaseUpdate,
		DeleteContext: resourceDliSQLDatabaseDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDatabaseImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9][\w_]{0,127}$`),
						"The name consists of 1 to 128 characters, starting with a letter or digit. "+
							"Only letters, digits and underscores (_) are allowed."),
					validation.StringMatch(regexp.MustCompile(`[A-Za-z_]`), "The name cannot be all digits."),
				),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsForceNewSchema(),
		},
	}
}

func resourceDliSQLDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	c, err := cfg.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v1 client: %s", err)
	}

	dbName := d.Get("name").(string)
	opts := databases.CreateOpts{
		Name:                dbName,
		Description:         d.Get("description").(string),
		EnterpriseProjectId: common.GetEnterpriseProjectID(d, cfg),
		Tags:                utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}
	respBody, err := databases.Create(c, opts)
	if err != nil {
		return diag.Errorf("error creating DLI database, %s", err)
	}
	if !respBody.IsSuccess {
		return diag.Errorf("unable to create the database: %s", respBody.Message)
	}

	// The resource ID (database name) at this time is only used as a mark the resource, and the value will be refreshed
	// in the READ method.
	d.SetId(dbName)

	return resourceDliSQLDatabaseRead(ctx, d, meta)
}

func GetDliSQLDatabaseByName(c *golangsdk.ServiceClient, dbName string) (databases.Database, error) {
	resp, err := databases.List(c, databases.ListOpts{
		Keyword: dbName, // Fuzzy matching.
	})
	if err != nil {
		return databases.Database{}, fmt.Errorf("error getting database: %s", err)
	}
	if !resp.IsSuccess {
		return databases.Database{}, fmt.Errorf("unable to query the database: %s", resp.Message)
	}

	if len(resp.Databases) < 1 {
		return databases.Database{}, golangsdk.ErrDefault404{}
	}
	for _, db := range resp.Databases {
		if db.Name == dbName {
			return db, nil
		}
	}

	return databases.Database{}, golangsdk.ErrDefault404{}
}

func resourceDliSQLDatabaseRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	c, err := cfg.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v1 client: %s", err)
	}

	dbName := d.Get("name").(string)
	db, err := GetDliSQLDatabaseByName(c, dbName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DLI database")
	}

	if db.ResourceId == "" {
		log.Printf("[WARN] unable to find the resource ID from the API response body during normal tenant (EPS " +
			"service has not been activated), it maybe cause the resource ID not being effectively referenced into " +
			"TMS tags management resource.")
		d.SetId(dbName)
	} else {
		d.SetId(db.ResourceId)
	}

	mErr := multierror.Append(nil,
		d.Set("name", db.Name),
		d.Set("description", db.Description),
		d.Set("enterprise_project_id", db.EnterpriseProjectId),
		d.Set("owner", db.Owner),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDliSQLDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	c, err := cfg.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v1 client: %s", err)
	}

	dbName := d.Get("name").(string)
	resp, err := databases.UpdateDBOwner(c, dbName, databases.UpdateDBOwnerOpts{
		NewOwner: d.Get("owner").(string),
	})
	if err != nil {
		return diag.Errorf("error updating SQL database owner: %s", err)
	}
	if !resp.IsSuccess {
		return diag.Errorf("unable to update the database owner: %s", resp.Message)
	}

	return resourceDliSQLDatabaseRead(ctx, d, meta)
}

func resourceDliSQLDatabaseDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	c, err := cfg.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v1 client: %s", err)
	}

	dbName := d.Get("name").(string)
	err = databases.Delete(c, dbName).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting SQL database: %s", err)
	}
	return nil
}

func resourceDatabaseImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	err := d.Set("name", d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error saving resource name of the DLI database: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}

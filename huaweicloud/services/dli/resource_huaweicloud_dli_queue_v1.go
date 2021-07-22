package dli

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/dli/v1/queues"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

var regexp4Name = regexp.MustCompile(`^[a-z0-9_]{1,128}$`)

const CU_16, CU_64, CU_256 = 16, 64, 256
const RESOURCE_MODE_SHARED, RESOURCE_MODE_EXCLUSIVE = 0, 1

func ResourceDliQueueV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceDliQueueCreate,
		Read:   resourceDliQueueV1Read,
		Delete: resourceDliQueueV1Delete,

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
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if !regexp4Name.MatchString(v) {
						errs = append(errs, fmtp.Errorf("%q can contain only digits, lower letters, and underscores (_) ", key))
					}
					return
				},
			},

			"queue_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "sql",
				ValidateFunc: validation.StringInSlice([]string{"sql", "general"}, false),
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"cu_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{CU_16, CU_64, CU_256}),
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "x86_64",
				ValidateFunc: validation.StringInSlice([]string{"x86_64", "aarch64"}, false),
			},

			"resource_mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{RESOURCE_MODE_SHARED, RESOURCE_MODE_EXCLUSIVE}),
			},

			"feature": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"basic", "ai"}, false),
			},

			"tags": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDliQueueCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	dliClient, err := config.DliV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("creating dli client failed: %s", err)
	}

	queueName := d.Get("name").(string)

	logp.Printf("[DEBUG] create dli queues queueName: %s", queueName)
	createOpts := queues.CreateOpts{
		QueueName:           queueName,
		QueueType:           d.Get("queue_type").(string),
		Description:         d.Get("description").(string),
		CuCount:             d.Get("cu_count").(int),
		EnterpriseProjectId: config.GetEnterpriseProjectID(d),
		Platform:            d.Get("platform").(string),
		ResourceMode:        d.Get("resource_mode").(int),
		Feature:             d.Get("feature").(string),
		Tags:                assembleTagsFromRecource("tags", d),
	}

	logp.Printf("[DEBUG] create dli queues using paramaters: %+v", createOpts)
	createResult := queues.Create(dliClient, createOpts)
	if createResult.Err != nil {
		return fmt.Errorf("create dli queues failed: %s", createResult.Err)
	}

	//query queue detail,trriger read to refresh the state
	d.SetId(queueName)

	return resourceDliQueueV1Read(d, meta)
}

func assembleTagsFromRecource(key string, d *schema.ResourceData) []tags.ResourceTag {
	if v, ok := d.GetOk(key); ok {
		tagRaw := v.(map[string]interface{})
		taglist := utils.ExpandResourceTags(tagRaw)
		return taglist
	}
	return nil

}
func resourceDliQueueV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	dliClient, err := config.DliV1Client(config.GetRegion(d))

	if err != nil {
		return fmtp.Errorf("error creating DliV1Client, err=%s", err)
	}

	queueName := d.Get("name").(string)

	queryOpts := queues.ListOpts{
		QueueType: d.Get("queue_type").(string),
	}

	logp.Printf("[DEBUG] create dli queues using paramaters: %+v", queryOpts)

	queryAllResult := queues.List(dliClient, queryOpts)
	if queryAllResult.Err != nil {
		return fmt.Errorf("query queues failed: %s", queryAllResult.Err)
	}

	//filter by queue_name
	queueDetail, err := filterByQueueName(queryAllResult.Body, queueName)
	if err != nil {
		return err
	}

	if queueDetail != nil {
		logp.Printf("[DEBUG]The detail of queue from SDK:%+v", queueDetail)

		d.Set("name", queueDetail.QueueName)
		d.Set("queue_type", queueDetail.QueueType)
		d.Set("description", queueDetail.Description)
		d.Set("cu_count", queueDetail.CuCount)
		d.Set("platform", queueDetail.Platform)
		d.Set("resource_mode", queueDetail.ResourceMode)
		d.Set("feature", queueDetail.Feature)
		d.Set("create_time", queueDetail.CreateTime)
	}

	return nil
}
func filterByQueueName(body interface{}, queueName string) (r *queues.Queue, err error) {
	if queueList, ok := body.(*queues.ListResult); ok {
		logp.Printf("[DEBUG]The list of queue from SDK:%+v", queueList)

		for _, v := range queueList.Queues {
			if v.QueueName == queueName {
				return &v, nil
			}
		}
		return nil, nil

	} else {
		return nil, fmt.Errorf("sdk-client response type is wrong, expect type:*queues.ListResult,acutal Type:%T",
			body)
	}

}
func resourceDliQueueV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.DliV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating DliV1Client, err=%s", err)
	}

	queueName := d.Get("name").(string)
	logp.Printf("[DEBUG] Deleting dli Queue %q", d.Id())

	result := queues.Delete(client, queueName)
	if result.Err != nil {
		return fmtp.Errorf("error deleting dli Queue %q, err=%s", d.Id(), result.Err)
	}

	return nil
}

package huaweicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/maas/v1/task"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMaasTaskV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceMaasTaskV1Create,
		Read:   resourceMaasTaskV1Read,
		Delete: resourceMaasTaskV1Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"src_node": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"ak": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"sk": {
							Type:      schema.TypeString,
							Sensitive: true,
							Required:  true,
							ForceNew:  true,
						},
						"object_key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"bucket": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"cloud_type": {
							Type:     schema.TypeString,
							Default:  "Aliyun",
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"dst_node": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"ak": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"sk": {
							Type:      schema.TypeString,
							Sensitive: true,
							Required:  true,
							ForceNew:  true,
						},
						"object_key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"bucket": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"enable_kms": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"thread_num": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"smn_info": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_urn": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"language": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"trigger_conditions": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func getSrcNode(d *schema.ResourceData) task.SrcNodeOpts {
	srcNodes := d.Get("src_node").([]interface{})
	srcNode := srcNodes[0].(map[string]interface{})

	srcNodeOpts := task.SrcNodeOpts{
		Region:    srcNode["region"].(string),
		AK:        srcNode["ak"].(string),
		SK:        srcNode["sk"].(string),
		ObjectKey: srcNode["object_key"].(string),
		Bucket:    srcNode["bucket"].(string),
		CloudType: srcNode["cloud_type"].(string),
	}

	log.Printf("[DEBUG] getSrcNode: %#v", srcNodeOpts)
	return srcNodeOpts
}

func getDstNode(d *schema.ResourceData) task.DstNodeOpts {
	dstNodes := d.Get("dst_node").([]interface{})
	dstNode := dstNodes[0].(map[string]interface{})

	dstNodeOpts := task.DstNodeOpts{
		Region:    dstNode["region"].(string),
		AK:        dstNode["ak"].(string),
		SK:        dstNode["sk"].(string),
		ObjectKey: dstNode["object_key"].(string),
		Bucket:    dstNode["bucket"].(string),
	}

	log.Printf("[DEBUG] getDstNode: %#v", dstNodeOpts)
	return dstNodeOpts
}

func getTriggerConditions(v []interface{}) []string {
	conds := make([]string, len(v))
	for i, cond := range v {
		conds[i] = cond.(string)
	}

	return conds
}

func getSmnInfo(d *schema.ResourceData) task.SmnInfoOpts {
	smnInfos := d.Get("smn_info").([]interface{})
	smnInfo := smnInfos[0].(map[string]interface{})
	triggerConditions := smnInfo["trigger_conditions"].(*schema.Set).List()

	smnInfoOpts := task.SmnInfoOpts{
		TopicUrn:          smnInfo["topic_urn"].(string),
		Language:          smnInfo["language"].(string),
		TriggerConditions: getTriggerConditions(triggerConditions),
	}
	log.Printf("[DEBUG] getSmnInfo: %#v", smnInfoOpts)
	return smnInfoOpts
}

func resourceMaasTaskV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	maasClient, err := config.maasV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating maas client: %s", err)
	}

	enableKMS := d.Get("enable_kms").(bool)
	createOpts := task.CreateOpts{
		SrcNode:     getSrcNode(d),
		DstNode:     getDstNode(d),
		EnableKMS:   &enableKMS,
		ThreadNum:   d.Get("thread_num").(int),
		Description: d.Get("description").(string),
	}

	var smnInfoOpts task.SmnInfoOpts
	if s, ok := d.GetOk("smn_info"); ok {
		smnInfo := (s.([]interface{}))[0].(map[string]interface{})
		triggerConditions := smnInfo["trigger_conditions"].(*schema.Set).List()

		smnInfoOpts = task.SmnInfoOpts{
			TopicUrn:          smnInfo["topic_urn"].(string),
			Language:          smnInfo["language"].(string),
			TriggerConditions: getTriggerConditions(triggerConditions),
		}
		createOpts.SmnInfo = &smnInfoOpts
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	taskCreate, err := task.Create(maasClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating Task: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	taskID := strconv.FormatInt(taskCreate.ID, 10)
	d.SetId(taskID)
	err = waitForTaskCompleted(maasClient, taskID, timeout)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for task (%s) completed: %s",
			taskID, err)
	}

	return resourceMaasTaskV1Read(d, meta)
}

func resourceMaasTaskV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	maasClient, err := config.maasV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating maas client: %s", err)
	}

	taskGet, err := task.Get(maasClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "task")
	}
	log.Printf("[DEBUG] Retrieved Task %s: %#v", d.Id(), taskGet)
	d.Set("name", taskGet.Name)
	d.Set("status", taskGet.Status)
	d.Set("enable_kms", taskGet.EnableKMS)
	d.Set("thread_num", taskGet.ThreadNum)
	d.Set("description", taskGet.Description)

	// TODO: we can not get AK/SK from the Get method, skip to set src_node and dst_node
	/*
		srcNodeList := make([]map[string]interface{}, 0, 1)
		srcNode := map[string]interface{}{
			"region":     taskGet.SrcNode.Region,
			"object_key": taskGet.SrcNode.ObjectKey[0],
			"bucket":     taskGet.SrcNode.Bucket,
		}
		srcNodeList = append(srcNodeList, srcNode)
		d.Set("src_node", srcNodeList)

		dstNodeList := make([]map[string]interface{}, 0, 1)
		dstNode := map[string]interface{}{
			"region":     taskGet.DstNode.Region,
			"object_key": taskGet.DstNode.ObjectKey,
			"bucket":     taskGet.DstNode.Bucket,
		}
		dstNodeList = append(dstNodeList, dstNode)
		d.Set("dst_node", dstNodeList)
	*/

	return nil
}

func resourceMaasTaskV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	maasClient, err := config.maasV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating maas client: %s", err)
	}

	id := d.Id()
	err = task.Delete(maasClient, id).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting task %s: %s", id, err)
	}
	d.SetId("")

	return nil
}

func getTaskStatus(maasClient *golangsdk.ServiceClient, taskId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		taskGet, err := task.Get(maasClient, taskId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] Task: %+v", taskGet)
		status := strconv.Itoa(taskGet.Status)
		return taskGet, status, nil
	}
}

func waitForTaskCompleted(maasClient *golangsdk.ServiceClient, taskID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"0", "1", "2", "3"},
		Target:     []string{"5"},
		Refresh:    getTaskStatus(maasClient, taskID),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

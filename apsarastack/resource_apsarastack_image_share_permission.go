package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApsaraStackImageSharePermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackImageSharePermissionCreate,
		Read:   resourceApsaraStackImageSharePermissionRead,
		Delete: resourceApsaraStackImageSharePermissionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceApsaraStackImageSharePermissionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	imageId := d.Get("image_id").(string)
	accountId := d.Get("account_id").(string)
	request := ecs.CreateModifyImageSharePermissionRequest()
	request.RegionId = client.RegionId
	request.ImageId = imageId
	accountSli := []string{accountId}
	request.AddAccount = &accountSli
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ModifyImageSharePermission(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_image_share_permission", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(imageId + ":" + accountId)
	return resourceApsaraStackImageSharePermissionRead(d, meta)
}

func resourceApsaraStackImageSharePermissionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client: client}
	object, err := ecsService.DescribeImageShareByImageId(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	d.Set("image_id", object.ImageId)
	d.Set("account_id", parts[1])
	return WrapError(err)
}

func resourceApsaraStackImageSharePermissionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := ecs.CreateModifyImageSharePermissionRequest()
	request.RegionId = client.RegionId
	parts, err := ParseResourceId(d.Id(), 2)
	request.ImageId = parts[0]
	accountSli := []string{parts[1]}
	request.RemoveAccount = &accountSli
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ModifyImageSharePermission(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_image_share_permission", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

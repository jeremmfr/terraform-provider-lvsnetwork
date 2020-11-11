package lvsnetwork

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	defaultVrrpScriptRise = 3
	defaultVrrpScriptFall = 3
)

func resourceVrrpScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceVrrpScriptCreate,
		Read:   resourceVrrpScriptRead,
		Update: resourceVrrpScriptUpdate,
		Delete: resourceVrrpScriptDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"script": {
				Type:     schema.TypeString,
				Required: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"weight_reverse": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"rise": {
				Type:     schema.TypeInt,
				Default:  defaultVrrpScriptRise,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value <= 0 {
						errors = append(errors, fmt.Errorf("[ERROR] %q must be positive integer", k))
					}

					return
				},
			},
			"fall": {
				Type:     schema.TypeInt,
				Default:  defaultVrrpScriptFall,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value <= 0 {
						errors = append(errors, fmt.Errorf("[ERROR] %q must be positive integer", k))
					}

					return
				},
			},
			"user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"init_fail": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceVrrpScriptCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	vrrpScript := createStrucVrrpScript(d)
	_, err := client.requestAPIVrrpScript("ADD", &vrrpScript)
	if err != nil {
		return err
	}
	d.SetId(d.Get("name").(string))

	return resourceVrrpScriptRead(d, m)
}

func resourceVrrpScriptRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	vrrpScript := createStrucVrrpScript(d)
	vrrpScriptRead, err := client.requestAPIVrrpScript("CHECK", &vrrpScript)
	if err != nil {
		return err
	}
	if vrrpScriptRead.Name == "" {
		d.SetId("")
	}
	tfErr := d.Set("script", vrrpScriptRead.Script)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("interval", vrrpScriptRead.Interval)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("timeout", vrrpScriptRead.Timeout)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("weight", vrrpScriptRead.Weight)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("weight_reverse", vrrpScriptRead.WeightReverse)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("rise", vrrpScriptRead.Rise)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("fall", vrrpScriptRead.Fall)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("user", vrrpScriptRead.User)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("init_fail", vrrpScriptRead.InitFail)
	if tfErr != nil {
		panic(tfErr)
	}

	return nil
}

func resourceVrrpScriptUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	d.Partial(true)
	vrrpScript := createStrucVrrpScript(d)
	_, err := client.requestAPIVrrpScript("CHANGE", &vrrpScript)
	if err != nil {
		return err
	}

	d.Partial(false)

	return resourceVrrpScriptRead(d, m)
}

func resourceVrrpScriptDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	vrrpScript := createStrucVrrpScript(d)
	_, err := client.requestAPIVrrpScript("REMOVE", &vrrpScript)
	if err != nil {
		return err
	}

	return nil
}

func createStrucVrrpScript(d *schema.ResourceData) vrrpScript {
	vrrpScript := vrrpScript{
		Name:          d.Get("name").(string),
		Script:        d.Get("script").(string),
		Interval:      d.Get("interval").(int),
		Timeout:       d.Get("timeout").(int),
		Weight:        d.Get("weight").(int),
		WeightReverse: d.Get("weight_reverse").(bool),
		Rise:          d.Get("rise").(int),
		Fall:          d.Get("fall").(int),
		User:          d.Get("user").(string),
		InitFail:      d.Get("init_fail").(bool),
	}

	return vrrpScript
}

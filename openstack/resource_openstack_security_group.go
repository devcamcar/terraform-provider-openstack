package openstack

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/secgroups"
)

func resourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityGroupCreate,
		Read:   resourceSecurityGroupRead,
		Update: resourceSecurityGroupUpdate,
		Delete: resourceSecurityGroupDelete,


		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"rules": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"to_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"ip_protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"ip_range": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

                        "parent_group_id": &schema.Schema{
                            Type: schema.TypeString,
                            Optional: false,
                        },

                        // TODO: Implement source group. Looks unimplemented in GopherCloud?
                        // http://godoc.org/github.com/rackspace/gophercloud/openstack/compute/v2/extensions/secgroups#GroupOpts
					},
				},
                // TODO: Implement security group rules.
				//Set: resourceSecurityGroupRuleHash,
			},
		},
	}
}

func resourceSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {

	provider := meta.(*gophercloud.ProviderClient)

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: d.Get("region").(string),
	})

	secgroup_opts := &secgroups.CreateOpts{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
        //Rules:          rules
	}

    newSecurityGroup, err := secgroups.Create(client, secgroup_opts).Extract()

	if err != nil {
		return err
	}

	d.SetId(newSecurityGroup.ID)

	/*stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     "ACTIVE",
		Refresh:    WaitForServerState(client, newServer),
		Timeout:    30 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return err
	}*/

	return nil
}

func resourceSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	/*provider := meta.(*gophercloud.ProviderClient)

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: d.Get("region").(string),
	})

	server, _ := servers.Get(client, d.Id()).Extract()

	servers.Delete(client, server.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "ERROR"},
		Target:     "",
		Refresh:    WaitForServerState(client, server),
		Timeout:    30 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	return err
    */

    return nil
}

func resourceSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	/*provider := meta.(*gophercloud.ProviderClient)

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: d.Get("region").(string),
	})

	server, _ := servers.Get(client, d.Id()).Extract()

	d.Partial(true)

	if d.HasChange("name") {
		_, err := servers.Update(client, server.ID, servers.UpdateOpts{
			Name: d.Get("name").(string),
		}).Extract()

		if err != nil {
			return err
		}

		d.SetPartial("name")
	}

	if d.HasChange("flavor_ref") {
		opts := &servers.ResizeOpts{
			FlavorRef: d.Get("flavor_ref").(string),
		}

		if res := servers.Resize(client, server.ID, opts); res.Err != nil {
			return res.Err
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"ACTIVE", "RESIZE"},
			Target:     "VERIFY_RESIZE",
			Refresh:    WaitForServerState(client, server),
			Timeout:    30 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()

		if err != nil {
			return err
		}

		// FIXME: proper error checking
		if res := servers.ConfirmResize(client, server.ID); res.Err != nil {
			return res.Err
		}

		d.SetPartial("flavor_ref")
	}

	d.Partial(false)

	return nil*/

    return nil
}

func resourceSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	/*provider := meta.(*gophercloud.ProviderClient)

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: d.Get("region").(string),
	})

	// FIXME: proper error checking
	server, _ := servers.Get(client, d.Id()).Extract()
	if err != nil {
		httpError, ok := err.(*perigee.UnexpectedResponseCodeError)
		if !ok {
			return err
		}

		if httpError.Actual == 404 {
			d.SetId("")
			return nil
		}

		return err
	}

	// TODO check networks, seucrity groups and floating ip

	d.Set("name", server.Name)
	d.Set("flavor_ref", server.Flavor["ID"])

	return nil
    */
    return nil
}

/*func WaitForServerState(client *gophercloud.ServiceClient, server *servers.Server) resource.StateRefreshFunc {

	return func() (interface{}, string, error) {
		latest, err := servers.Get(client, server.ID).Extract()
		if err != nil {
			return nil, "", err
		}

		return latest, latest.Status, nil

	}
}*/


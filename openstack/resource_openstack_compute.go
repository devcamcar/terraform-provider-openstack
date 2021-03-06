package openstack

import (
	"crypto/sha1"
	//"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	//"github.com/haklop/gophercloud-extensions/network"
)

func resourceCompute() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeCreate,
		Read:   resourceComputeRead,
		Update: resourceComputeUpdate,
		Delete: resourceComputeDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"image_ref": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"flavor_ref": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"security_groups": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true, // TODO handle update
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},

			"user_data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// just stash the hash for state & diff comparisons
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						hash := sha1.Sum([]byte(v.(string)))
						return hex.EncodeToString(hash[:])
					default:
						return ""
					}
				},
			},

			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"networks": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true, // TODO handle update
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},

			// No idea how to do this yet.
			//"metadata": &schema.Schema{
			//},

			// No idea how to do this yet.
			//"personality": &schema.Schema{
			//},

			"config_drive": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},

			"admin_pass": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// Region is defined per-instance due to how gophercloud
			// handles the region -- not until a provider is returned.
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// defined in gophercloud compute extensions
			"key_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true, // TODO handle update
			},

			// defined in haklop's network extensions
			// Neutron only, I think
			"floating_ip_pool": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"floating_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeCreate(d *schema.ResourceData, meta interface{}) error {

	provider := meta.(*gophercloud.ProviderClient)

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: d.Get("region").(string),
	})

	nets := d.Get("networks").(*schema.Set)
	var networks []servers.Network
	for _, v := range nets.List() {
		networks = append(networks, servers.Network{UUID: v.(string)})
	}

	sec_groups := d.Get("security_groups").(*schema.Set)
	var security_groups []string
	for _, v := range sec_groups.List() {
		security_groups = append(security_groups, v.(string))
	}

	// API needs it to be base64 encoded.
	/*
	   userData := ""
	   if v := d.Get("user_data"); v != nil {
	     userData = base64.StdEncoding.EncodeToString([]byte(v.(string)))
	   }
	*/

	base_opts := &servers.CreateOpts{
		Name:           d.Get("name").(string),
		ImageRef:       d.Get("image_ref").(string),
		FlavorRef:      d.Get("flavor_ref").(string),
		SecurityGroups: security_groups,
		// I'm not sure if this works
		Networks: networks,
		// Need to convert this to type byte
		//UserData:       userData,
	}

	key_name := d.Get("key_name").(string)
	key_opts := keypairs.CreateOptsExt{
		CreateOptsBuilder: base_opts,
		KeyName:           key_name,
	}

	newServer, err := servers.Create(client, key_opts).Extract()

	if err != nil {
		return err
	}

	d.SetId(newServer.ID)

	stateConf := &resource.StateChangeConf{
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
	}

	// FIXME: add floating IP support
	/*
	   pool := d.Get("floating_ip_pool").(string)
	   if len(pool) > 0 {
	     var newIp gophercloud.FloatingIp
	     hasFloatingIps := false

	     floatingIps, err := serversApi.ListFloatingIps()
	     if err != nil {
	       return err
	     }

	     for _, element := range floatingIps {
	       // use first floating ip available on the pool
	       if element.Pool == pool && element.InstanceId == "" {
	         newIp = element
	         hasFloatingIps = true
	       }
	     }

	     // if there is no available floating ips, try to create a new one
	     if !hasFloatingIps {
	       newIp, err = serversApi.CreateFloatingIp(pool)
	       if err != nil {
	         return err
	       }
	     }

	     err = serversApi.AssociateFloatingIp(newServer.Id, newIp)
	     if err != nil {
	       return err
	     }

	     d.Set("floating_ip", newIp.Ip)

	     // Initialize the connection info
	     d.SetConnInfo(map[string]string{
	       "type": "ssh",
	       "host": newIp.Ip,
	     })
	   }
	*/

	return nil
}

func resourceComputeDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*gophercloud.ProviderClient)

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
}

func resourceComputeUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*gophercloud.ProviderClient)

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

	return nil
}

func resourceComputeRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*gophercloud.ProviderClient)

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
}

func WaitForServerState(client *gophercloud.ServiceClient, server *servers.Server) resource.StateRefreshFunc {

	return func() (interface{}, string, error) {
		latest, err := servers.Get(client, server.ID).Extract()
		if err != nil {
			return nil, "", err
		}

		return latest, latest.Status, nil

	}
}

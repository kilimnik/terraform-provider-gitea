package gitea

import (
	"code.gitea.io/sdk/gitea"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	pushMirrorOwner          string = "owner"
	pushMirrorRepo           string = "repo"
	pushMirrorRemotePassword string = "remote_password"
	pushMirrorRemoteUsername string = "remote_username"
	pushMirrorCreated        string = "created"
	pushMirrorInterval       string = "interval"
	pushMirrorLastError      string = "last_error"
	pushMirrorLastUpdate     string = "last_update"
	pushMirrorRemoteAddress  string = "remote_address"
	pushMirrorRemoteName     string = "remote_name"
	pushMirrorsRepoName      string = "repo_name"
	pushMirrorSyncOnCommit   string = "sync_on_commit"
)

func resourcePushMirrorRead(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*gitea.Client)

	var resp *gitea.Response

	if err != nil {
		return err
	}

	pushMirror, resp, err := client.GetPushMirror(d.Get(pushMirrorOwner).(string), d.Get(pushMirrorRepo).(string), d.Get(pushMirrorRemoteName).(string))

	if err != nil {
		if resp.StatusCode == 404 {
			return nil
		} else {
			return err
		}
	}

	err = setPushMirrorResourceData(pushMirror, d)

	return
}

func resourcePushMirrorCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*gitea.Client)

	pushMirror, _, err := client.AddPushMirror(d.Get(pushMirrorOwner).(string), d.Get(pushMirrorRepo).(string), gitea.AddPushMirrorOption{
		Interval:       d.Get(pushMirrorInterval).(string),
		RemoteAddress:  d.Get(pushMirrorRemoteAddress).(string),
		RemotePassword: d.Get(pushMirrorRemotePassword).(string),
		RemoteUsername: d.Get(pushMirrorRemoteUsername).(string),
		SyncOnCommit:   d.Get(pushMirrorSyncOnCommit).(bool),
	})

	if err != nil {
		return err
	}

	err = setPushMirrorResourceData(pushMirror, d)

	return
}

func resourcePushMirrorUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	// update = recreate
	respurcePushMirrorDelete(d, meta)
	resourcePushMirrorCreate(d, meta)
	return

}

func respurcePushMirrorDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*gitea.Client)

	client.DeletePushMirror(d.Get(pushMirrorOwner).(string), d.Get(pushMirrorRepo).(string), d.Get(pushMirrorRemoteName).(string))

	return
}

func setPushMirrorResourceData(pushMirror *gitea.PushMirror, d *schema.ResourceData) (err error) {
	d.Set("created", pushMirror.Created.String())
	d.Set("interval", pushMirror.Interval)
	d.Set("last_error", pushMirror.LastError)
	d.Set("last_update", pushMirror.LastUpdate.String())
	d.Set("remote_address", pushMirror.RemoteAddress)
	d.Set("remote_name", pushMirror.RemoteName)
	d.Set("repo_name", pushMirror.RepoName)
	d.Set("sync_on_commit", pushMirror.SyncOnCommit)

	return
}

func resourceGiteaPushMirror() *schema.Resource {
	return &schema.Resource{
		Read:   resourcePushMirrorRead,
		Create: resourcePushMirrorCreate,
		Update: resourcePushMirrorUpdate,
		Delete: respurcePushMirrorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"owner": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Owner of the repository",
			},
			"repo": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Name of the repository",
			},
			"interval": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Default:     "8h0m0s",
				Description: "valid time units are 'h', 'm', 's'. 0 to disable automatic sync",
			},
			"remote_address": {
				Type:        schema.TypeString,
				Required:    true,
				Optional:    false,
				ForceNew:    true,
				Description: "The adress to the remote repository that this repository should be mirrored to",
			},
			"remote_password": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Sensitive:   true,
				Description: "The password for the remote repository that this repository should be mirrored to",
			},
			"remote_username": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "The username for the remote repository that this repository should be mirrored to",
			},
			"sync_on_commit": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Description: "If the repository should be synced on every commit",
			},
		},
		Description: "`gitea_push_mirror` manages gitea repository push mirrors.\n\n" +
			"Push mirrors are a way to mirror a gitea repository to a remote repository.",
	}
}

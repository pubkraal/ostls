package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pubkraal/ostls/data"
	"github.com/pubkraal/ostls/util"
)

func Config(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	nodeKey, err := util.ExtractNodeKey(r)
	if err != nil {
		writeFailure(w, r)
		return
	}
	uuids := util.ConvertTokensToSet(data.ListValidTokens(dbHandle))
	exists := uuids[nodeKey]
	if !exists {
		writeFailure(w, r)
		return
	}
	// if none of these are true, rip node

	clientConfig := data.LoadConfigByTitle("client-config", dbHandle)
	configJson := clientConfig.Content
	if configJson == "" {
		configJson = defaultConfig()
	}

	fmt.Fprint(w, configJson)
}

func defaultConfig() string {
	return `{
  "options": {
    "utc": true
  },
  "schedule": {
    "system_info": {
      "query": "SELECT hostname, cpu_brand, physical_memory FROM system_info;",
      "interval": 86400
    },
    "disk_encryption_snapshot": {
      "query": "SELECT * FROM disk_encryption;",
      "interval": 28800,
      "description": "Disk encryption status and information.",
      "snapshot": true
    },
    "disk_free_space_pct": {
      "query": "select device, (blocks_available * 100 / blocks) as pct from mounts where type not like '%tmp%' and blocks_available > 0;",
      "interval": 3600,
      "description": "Displays the percentage of free space available on the primary disk partition",
      "snapshot": true
    },
    "last": {
      "query": "SELECT * FROM last;",
      "interval": 3600,
      "description": "System logins and logouts.",
      "removed": false
    }
  },
  "decorators": {
    "load": [
      "SELECT uuid AS host_uuid FROM system_info;",
      "SELECT user AS username FROM logged_in_users ORDER BY time DESC LIMIT 1;"
    ]
  },
  "packs": {
    "osx-attacks": {
      "platform": "darwin",
      "queries": {
        "Leverage-A_1": {
          "query" : "select * from launchd where path like '%UserEvent.System.plist';",
          "interval" : "3600",
          "version": "1.4.5",
          "description" : "(http://www.intego.com/mac-security-blog/new-mac-trojan-discovered-related-to-syria/)",
          "value" : "Artifact used by this malware"
        },
        "Leverage-A_2": {
          "query" : "select * from file where path = '/Users/Shared/UserEvent.app';",
          "interval" : "3600",
          "version": "1.4.5",
          "description" : "(http://www.intego.com/mac-security-blog/new-mac-trojan-discovered-related-to-syria/)",
          "value" : "Artifact used by this malware"
        },
        "Leverage-A_3": {
          "query" : "select * from launchd where name = 'com.GetFlashPlayer.plist';",
          "interval" : "3600",
          "version": "1.4.5",
          "description" : "(https://www.volexity.com/blog/2017/07/24/real-news-fake-flash-mac-os-x-users-targeted/)",
          "value" : "Artifact used by this malware"
        }
      }
    },
	"unwanted-chrome-extensions": {
      "platform": "windows,darwin",
      "queries": {
        "BetternetVPN": {
          "query": "SELECT * FROM users JOIN chrome_extensions USING (uid) WHERE identifier='gjknjjomckknofjidppipffbpoekiipm';",
          "interval": 3600,
          "description": "(https://www.bleepingcomputer.com/news/security/eight-chrome-extensions-hijacked-to-deliver-malicious-code-to-4-8-million-users/)"
        },
        "Chrometana": {
          "query": "SELECT * FROM users JOIN chrome_extensions USING (uid) WHERE identifier='kaicbfmipfpfpjmlbpejaoaflfdnabnc';",
          "interval": 3600,
          "description": "(https://www.bleepingcomputer.com/news/security/eight-chrome-extensions-hijacked-to-deliver-malicious-code-to-4-8-million-users/)"
        },
        "CopyFish": {
          "query": "SELECT * FROM users JOIN chrome_extensions USING (uid) WHERE identifier='eenjdnjldapjajjofmldgmkjaienebbj';",
          "interval": 3600,
          "description": "(https://www.bleepingcomputer.com/news/security/copyfish-chrome-extension-hijacked-to-show-adware/)"
        },
        "Giphy": {
          "query": "SELECT * FROM users JOIN chrome_extensions USING (uid) WHERE identifier='kaicbfmipfpfpjmlbpejaoaflfdnabnc';",
          "interval": 3600,
          "description": "(https://www.reddit.com/r/chrome/comments/6htzan/psawarning_giphy_extension_6172017_is_now_malware/)"
        },
        "HolaVPN": {
          "query": "SELECT * FROM users JOIN chrome_extensions USING (uid) WHERE identifier='gkojfkhlekighikafcpjkiklfbnlmeio';",
          "interval": 3600,
          "description": "(http://adios-hola.org)"
        },
        "InfinityNewTab": {
          "query": "SELECT * FROM users JOIN chrome_extensions USING (uid) WHERE identifier='dbfmnekepjoapopniengjbcpnbljalfg';",
          "interval": 3600,
          "description": "(https://www.bleepingcomputer.com/news/security/eight-chrome-extensions-hijacked-to-deliver-malicious-code-to-4-8-million-users/)"
        },
        "SocialFixer": {
          "query": "SELECT * FROM users JOIN chrome_extensions USING (uid) WHERE identifier='ifmhoabcaeehkljcfclfiieohkohdgbb';",
          "interval": 3600,
          "description": "(https://www.bleepingcomputer.com/news/security/eight-chrome-extensions-hijacked-to-deliver-malicious-code-to-4-8-million-users/)"
        },
        "TouchVPN": {
          "query": "SELECT * FROM users JOIN chrome_extensions USING (uid) WHERE identifier='bihmplhobchoageeokmgbdihknkjbknd';",
          "interval": 3600,
          "description": "(https://www.bleepingcomputer.com/news/security/eight-chrome-extensions-hijacked-to-deliver-malicious-code-to-4-8-million-users/)"
        },
        "WebDeveloper": {
          "query": "SELECT * FROM users JOIN chrome_extensions USING (uid) WHERE identifier='bfbameneiokkgbdmiekhjnmfkcnldhhm';",
          "interval": 3600,
          "description": "(https://www.bleepingcomputer.com/news/security/chrome-extension-with-over-one-million-users-hijacked-to-serve-adware/)"
        },
        "WebPaint": {
          "query": "SELECT * FROM users JOIN chrome_extensions USING (uid) WHERE identifier='emeokgokialpjadjaoeiplmnkjoaegng';",
          "interval": 3600,
          "description": "(https://www.bleepingcomputer.com/news/security/eight-chrome-extensions-hijacked-to-deliver-malicious-code-to-4-8-million-users/)"
        },
        "MacOSInstallCore": {
          "query": "SELECT * FROM users JOIN chrome_extensions USING (uid) WHERE identifier='hinehnlkkmckjblijjpbpamhljokoohh';",
          "interval": 3600,
          "description": "(https://www.virustotal.com/#/file/5cab0821f597100dc1170bfef704d8cebaf67743e9d509e83b0b208eb630d992/detection)"
        }
      }
    }
  }
}`
}

package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pubkraal/ostls/data"
	"github.com/pubkraal/ostls/util"
)

// This is the default configuration you can get out of osquery.conf
// I stripped it and removed the windows specific things.
// This should probably come out of some database configuration
/*
{
  "schedule": {
    "system_info": {
      "query": "SELECT hostname, cpu_brand, physical_memory FROM system_info;",
      "interval": 3600
    }
  },
  "decorators": {
    "load": [
      "SELECT uuid AS host_uuid FROM system_info;",
      "SELECT user AS username FROM logged_in_users ORDER BY time DESC LIMIT 1;"
    ]
  },
  "packs": {
    "osquery-monitoring": "/var/osquery/packs/osquery-monitoring.conf",
    "incident-response": "/var/osquery/packs/incident-response.conf",
    "it-compliance": "/var/osquery/packs/it-compliance.conf",
    "osx-attacks": "/var/osquery/packs/osx-attacks.conf",
    "vuln-management": "/var/osquery/packs/vuln-management.conf",
    "hardware-monitoring": "/var/osquery/packs/hardware-monitoring.conf",
    "ossec-rootkit": "/var/osquery/packs/ossec-rootkit.conf"
  }
}
*/

func Config(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	nodeKey, err := util.ExtractNodeKey(r)
	if err != nil {
		writeFailure(w, r)
		return
	}
	// Check if request contains node key
	// Check if node is enrolled
	uuids := util.ConvertTokensToSet(data.ListValidTokens(dbHandle))
	exists := uuids[nodeKey]
	if !exists {
		writeFailure(w, r)
		return
	}
	// if none of these are true, rip node

	fmt.Fprint(w, `{
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
      "query": "SELECT (blocks_available * 100 / blocks) AS pct FROM mounts WHERE device='/dev/disk1s1';",
      "interval": 3600,
      "description": "Displays the percentage of free space available on the primary disk partition",
      "snapshot": true
    },
    "last": {
      "query": "SELECT * FROM last;",
      "interval": 3600,
      "description": "System logins and logouts.",
      "removed": false
    },
    "chrome_extensions": {
      "query": "SELECT * FROM users CROSS JOIN chrome_extensions USING (uid);",
      "interval": 3600,
      "description": "List installed Chrome Extensions for all users"
    },
    "safari_extensions": {
      "query": "SELECT * FROM users CROSS JOIN safari_extensions USING (uid);",
      "interval": 3600,
      "description": "Safari browser extension details for all users."
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
    }
  }
}
`)
}

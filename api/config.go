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
`)
}

package vms

import (
	"encoding/json"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

/*
{
  "host_solutions_status": {
    "compliances": {},
    "status": "COMPLIANT"
  },
  "cluster_solutions_status": {
    "compliances": {
      "solution3": {
        "compliances": {},
        "status": "NON_COMPLIANT"
      }
    },
    "status": "NON_COMPLIANT"
  },
  "status": "NON_COMPLIANT"
}
*/

func TestClusterCompliance(t *testing.T) {
	in := []byte("{\"host_solutions_status\":{\"compliances\":{},\"status\":\"COMPLIANT\"},\"cluster_solutions_status\":{\"compliances\":{\"solution3\":{\"compliances\":{},\"status\":\"NON_COMPLIANT\"}},\"status\":\"NON_COMPLIANT\"},\"status\":\"NON_COMPLIANT\"}")

	var c ClusterCompliance

	require.NoError(t, json.Unmarshal(in, &c))

	spew.Dump(c)

	_, ok := c.ClusterSolutionsStatus.Compliances["solution3"]
	require.True(t, ok)
}

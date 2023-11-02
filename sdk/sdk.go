package sdk

import (
	pb "github.com/sqlc-dev/plugin-sdk-go/plugin"

	"github.com/sqlc-dev/plugin-sdk-go/pattern"
)

func DataType(n *pb.Identifier) string {
	if n.Schema != "" {
		return n.Schema + "." + n.Name
	} else {
		return n.Name
	}
}

func MatchString(pat, target string) bool {
	matcher, err := pattern.MatchCompile(pat)
	if err != nil {
		panic(err)
	}
	return matcher.MatchString(target)
}

func SameTableName(tableID, f *pb.Identifier, defaultSchema string) bool {
	if tableID == nil {
		return false
	}
	schema := tableID.Schema
	if tableID.Schema == "" {
		schema = defaultSchema
	}
	return tableID.Catalog == f.Catalog && schema == f.Schema && tableID.Name == f.Name
}

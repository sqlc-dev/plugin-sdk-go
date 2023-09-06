package sdk

import (
	pb "buf.build/gen/go/sqlc/sqlc/protocolbuffers/go/protos/plugin"

	"github.com/sqlc-dev/sqlc-go/pattern"
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

func Matches(o *pb.Override, n *pb.Identifier, defaultSchema string) bool {
	if n == nil {
		return false
	}
	schema := n.Schema
	if n.Schema == "" {
		schema = defaultSchema
	}
	if o.Table.Catalog != "" && !MatchString(o.Table.Catalog, n.Catalog) {
		return false
	}
	if o.Table.Schema == "" && schema != "" {
		return false
	}
	if o.Table.Schema != "" && !MatchString(o.Table.Schema, schema) {
		return false
	}
	if o.Table.Name == "" && n.Name != "" {
		return false
	}
	if o.Table.Name != "" && !MatchString(o.Table.Name, n.Name) {
		return false
	}
	return true
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

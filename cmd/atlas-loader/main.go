package main

import (
	"fmt"
	"io"
	"os"

	"nebula/model/cmdb"
	"nebula/model/config"
	"nebula/model/config/namespace"
	"nebula/model/config/namespace/file"
	"nebula/model/config/namespace/kv"
	"nebula/model/iam"
	"nebula/model/identity"
	"nebula/model/setting"
	"nebula/model"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	stmts, err := gormschema.New("mysql").Load(
		// CMDB models
		&cmdb.DNS{},
		&cmdb.Machine{},
		
		// Config models
		&config.App{},
		&config.Env{},
		&config.Host{},
		&config.Namespace{},
		&config.Project{},
		&config.Setting{},
		
		// Namespace models
		&namespace.File{},
		&namespace.KV{},
		&namespace.Tag{},
		
		// File and KV history models
		&file.History{},
		&kv.History{},
		
		// IAM models
		&iam.KCGroup{},
		&iam.KCUser{},
		
		// Identity models
		&identity.User{},
		
		// Setting models
		&setting.Project{},
		&setting.Region{},
		&setting.Tenant{},
		&setting.Vendor{},
		
		// Other models
		&model.TableColumn{},
		&model.Tenant{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
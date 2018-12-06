package cmd

import (
	"context"
	"fmt"
	"path"

	"github.com/remeh/sizedwaitgroup"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jtopjian/bagel/lib/connections"
	"github.com/jtopjian/bagel/lib/inventories"
	"github.com/jtopjian/bagel/lib/resources"
	"github.com/jtopjian/bagel/lib/site"
	"github.com/jtopjian/bagel/lib/utils"
)

var (
	cliRole      string
	cliInventory string
	cliTarget    string
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy a role to a set of targets",
	Run:   deploy,
}

func init() {
	deployCmd.PersistentFlags().StringVarP(&cliRole, "role", "r", "", "role to deploy")
	deployCmd.PersistentFlags().StringVarP(&cliInventory, "inventory", "i", "", "inventory to query")
	deployCmd.PersistentFlags().StringVarP(&cliTarget, "target", "t", "", "single target to deploy to")
}

func deploy(cmd *cobra.Command, args []string) {
	log := utils.GetLogger()

	// Parse the site file
	siteDir := viper.GetString("site_dir")
	sitePath := path.Join(siteDir, "site.yaml")
	siteFile, err := site.New(sitePath)
	if err != nil {
		log.Fatalf("Unable to load site file %s: %s", sitePath, err)
	}

	// If a specific role was defined, use it.
	// Otherwise, use all roles defined in the site file.
	roles := make(map[string]site.Role)
	if cliRole != "" {
		r, ok := siteFile.Roles[cliRole]
		if !ok {
			log.Fatalf("Role %s is not defined", cliRole)
		}

		roles[cliRole] = r
	} else {
		for roleName, roleInfo := range siteFile.Roles {
			roles[roleName] = roleInfo
		}
	}

	for roleName, role := range roles {
		inv := role.Inventories
		if len(inv) == 0 {
			log.Fatalf("No inventories defined for role %s", roleName)
		}

		for _, invName := range inv {
			// If a specific inventory was specified, skip all others.
			if cliInventory != "" && cliInventory != invName {
				continue
			}

			invInfo, ok := siteFile.Inventories[invName]
			if !ok {
				log.Fatalf("Inventory %s is not defined", invName)
			}

			connName := invInfo.Connection
			connInfo, ok := siteFile.Connections[connName]
			if !ok {
				log.Fatalf("Connection %s is not defined", connName)
			}

			if err := invInfo.DiscoverTargets(); err != nil {
				log.Fatalf("Unable to discover targets in %s: %s", invName, err)
			}

			// Connect to each discovered target.
			swg := sizedwaitgroup.New(viper.GetInt("parallel"))
			for _, t := range invInfo.Targets {
				// If a specific target was specified, skip all others.
				if cliTarget != "" && cliTarget != t.Address {
					continue
				}

				t.ConnectionName = connName
				t.ConnectionType = connInfo.Type
				t.ConnectionOptions = connInfo.Options

				swg.Add()
				go func(roleName string, target inventories.Target) {
					defer swg.Done()
					connOptions := target.ConnectionOptions
					connOptions["host"] = target.Address
					conn, err := connections.New(target.ConnectionType, connOptions)
					if err != nil {
						log.Errorf("Error creating connection to %s: %s", target.Address, err)
						return
					}

					if err := conn.Connect(); err != nil {
						log.Errorf("Error connecting to %s: %s", target.Address, err)
						return
					}
					defer conn.Close()

					L := utils.LuaPool.Get()
					defer utils.LuaPool.Shutdown()

					ctx := context.WithValue(context.Background(), "connection", conn)
					L.SetContext(ctx)
					resources.Register(L)

					file := fmt.Sprintf("/opt/bagel/roles/%s.lua", roleName)
					if err := L.DoFile(file); err != nil {
						log.Errorf("Error deploying role %s: %s", roleName, err)
						return
					}

					return
				}(roleName, t)
				swg.Wait()
			}
		}
	}
}

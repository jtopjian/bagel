package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/jtopjian/bagel/lib/connections"
	"github.com/jtopjian/bagel/lib/resources"
	"github.com/jtopjian/bagel/lib/utils"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a script on localhost",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	log := utils.GetLogger()

	if len(args) < 1 {
		log.Fatal("Usage: run <file>")
	}

	file := args[0]

	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Fatalf("File %s does not exist", file)
	}

	L := utils.LuaPool.Get()
	defer utils.LuaPool.Shutdown()

	conn, err := connections.NewLocalConnection()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.WithValue(context.Background(), "connection", conn)
	L.SetContext(ctx)

	resources.Register(L)

	if err := L.DoFile(file); err != nil {
		log.Fatal(err)
	}
}

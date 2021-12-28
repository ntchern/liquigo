package cmd

import (
	"database/sql"
	"log"

	"github.com/lib/pq"
	"github.com/ntchern/liquigo/liquigo"
	"github.com/spf13/cobra"
)

var updateURL string
var updateSchema string
var updateFile string

func initUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Apply all DB migrations to the target DB",
		Long:  `Apply all DB migrations to the target DB`,
		SuggestFor: []string{"apply"},
		Run: func(cmd *cobra.Command, args []string) {
			update()
		},
	}
	cmd.Flags().StringVar(&updateURL, "url", "", "DB URL")
	cmd.Flags().StringVar(&updateSchema, "schema", "", "Set current DB schema (optional)")
	cmd.Flags().StringVar(&updateFile, "changeLog", "", "Changeset file")
	cmd.MarkFlagRequired("url")
	cmd.MarkFlagRequired("changeLog")
	return cmd
}

func update() {
	psqlInfo, err := pq.ParseURL(updateURL)
	if err != nil {
		log.Fatal(err.Error())
	}
	if updateSchema != "" {
		psqlInfo = psqlInfo + " search_path=" + updateSchema
	}
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = liquigo.Update(db, updateFile)
	if err != nil {
		log.Fatal(err.Error())
	}
}

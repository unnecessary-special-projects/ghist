package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/unnecessary-special-projects/ghist/internal/api"
	"github.com/unnecessary-special-projects/ghist/internal/project"
	"github.com/spf13/cobra"
)

// webFS is set from main.go via SetWebFS
var webFS embed.FS
var hasWebFS bool

func SetWebFS(wfs embed.FS) {
	webFS = wfs
	hasWebFS = true
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start local web server with Kanban board UI",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, s, err := openStore()
		if err != nil {
			return err
		}
		// Note: we don't defer s.Close() here because the server runs indefinitely

		port, _ := cmd.Flags().GetInt("port")
		dev, _ := cmd.Flags().GetBool("dev")

		var frontendFS fs.FS
		if !dev && hasWebFS {
			sub, err := fs.Sub(webFS, "web/dist")
			if err != nil {
				log.Printf("warning: embedded frontend not found, serving API only: %v", err)
			} else {
				frontendFS = sub
			}
		}

		repoURL := project.DetectGitHubRepo(root)
		srv := api.NewServer(s, frontendFS, dev, repoURL)
		addr := fmt.Sprintf(":%d", port)

		fmt.Printf("ghist server starting on http://localhost:%d\n", port)
		if dev {
			fmt.Println("  Dev mode: CORS enabled, no embedded frontend")
		}

		return http.ListenAndServe(addr, srv.Handler())
	},
}

func init() {
	serveCmd.Flags().IntP("port", "p", 4777, "Port to listen on")
	serveCmd.Flags().Bool("dev", false, "Enable dev mode (CORS, no embedded frontend)")
	rootCmd.AddCommand(serveCmd)
}

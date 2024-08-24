package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
)

var showFileNames bool

var branchesCmd = &cobra.Command{
	Use:   "show-branch-status",
	Short: "List all branches with their status relative to remote, and show file modification status",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := git.PlainOpen(".")
		if err != nil {
			log.Fatal(err)
		}

		fetchCmd := exec.Command("git", "fetch", "origin")
		if err := fetchCmd.Run(); err != nil {
			log.Fatalf("Failed to fetch from origin: %v", err)
		}

		branches, err := repo.Branches()
		if err != nil {
			log.Fatal(err)
		}

		err = branches.ForEach(func(branchRef *plumbing.Reference) error {
			branchName := branchRef.Name().Short()
			fmt.Printf("\nBranch: %s\n", branchName)

			remoteRefName := fmt.Sprintf("refs/remotes/origin/%s", branchName)
			remoteRef, err := repo.Reference(plumbing.ReferenceName(remoteRefName), true)
			if err != nil {
				color.Red("Remote branch '%s' does not exist.\n", branchName)
				return nil
			}

			compareCmd := exec.Command("git", "rev-list", "--left-right", fmt.Sprintf("%s...%s", branchRef.Hash().String(), remoteRef.Hash().String()))
			output, err := compareCmd.CombinedOutput()
			if err != nil {
				log.Fatalf("Failed to compare branches: %v", err)
			}

			ahead, behind := 0, 0
			for _, line := range output {
				if line == '<' {
					ahead++
				} else if line == '>' {
					behind++
				}
			}

			if ahead > 0 {
				color.Green("Ahead by %d commits", ahead)
			}
			if behind > 0 {
				color.Yellow("Behind by %d commits", behind)
			}
			if ahead > 0 && behind > 0 {
				color.Red("Merge is needed.")
			} else if ahead > 0 {
				color.Red("Push is needed.")
			} else if behind > 0 {
				color.Red("Pull is needed.")
			} else {
				color.Green("Branch is up to date.")
			}

			statusCmd := exec.Command("git", "status", "--porcelain")
			statusOutput, err := statusCmd.CombinedOutput()
			if err != nil {
				log.Fatalf("Failed to get status: %v", err)
			}

			if len(string(statusOutput)) > 0 {
				color.Red("Files are modified, please commit them before merging or pulling.\nRun 'git status' for more information.")
			} else {
				color.Green("No modifications or new files.")
			}

			fmt.Println()

			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(branchesCmd)
}

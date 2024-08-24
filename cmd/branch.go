package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
)

var branchName string

var branchCmd = &cobra.Command{
	Use:   "branch-status",
	Short: "Display the status of a branch",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := git.PlainOpen(".")
		if err != nil {
			log.Fatal(err)
		}

		if branchName == "" {
			fmt.Println("Please specify a branch with --branch")
			return
		}

		branchRef, err := repo.Branch(branchName)
		branchHash := plumbing.ReferenceName("refs/heads/" + branchRef.Name)
		ref, err := repo.Reference(branchHash, true)
		if err != nil {
			log.Fatal(err)
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Branch: %s\n", branchName)

		fmt.Println("\nCommit History:")
		iter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
		if err != nil {
			log.Fatal(err)
		}

		iter.ForEach(func(c *object.Commit) error {
			fmt.Printf("Commit: %s\nAuthor: %s\nDate: %s\nMessage: %s\n\n",
				c.Hash, c.Author.Name, c.Author.When, c.Message)
			return nil
		})

		fmt.Println("\nRemote Status:")
		remotes, err := repo.Remotes()
		if err != nil {
			log.Fatal(err)
		}

		for _, remote := range remotes {
			fmt.Printf("Remote: %s\n", remote.String())

			refName := fmt.Sprintf("refs/remotes/%s/%s", remote.String(), branchName)
			ref, err := repo.Reference(plumbing.ReferenceName(refName), true)
			if err != nil {
				fmt.Printf("Remote branch '%s' does not exist.\n", branchName)
				continue
			}

			fmt.Printf("Remote SHA: %s\n", ref.Hash().String())
			fmt.Println("Compare with remote:")
			cmd := exec.Command("git", "diff", fmt.Sprintf("%s..%s", branchHash, ref.Hash().String()))
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", string(output))
		}

		fmt.Println("\nChecking for potential merge conflicts:")
		mergeCmd := exec.Command("git", "merge", "--no-commit", "--no-ff", string(branchHash))
		out, err := mergeCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Merge conflicts detected:\n%s\n", string(out))
			exec.Command("git", "merge", "--abort").Run()
		} else {
			fmt.Println("No merge conflicts detected.")
		}
	},
}

func init() {
	branchCmd.Flags().StringVarP(&branchName, "branch", "b", "", "Specify the branch name")
	rootCmd.AddCommand(branchCmd)
}

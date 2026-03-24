package cmdAdmin

import (
	"fmt"
	"gamesync/internal/dbm"
	"gamesync/internal/ui"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

type keyCmd struct {
	cmd *cobra.Command
}
func newKeyCmd(user *dbm.UserWithRole) *keyCmd {
	root := keyCmd{}
	cmd := &cobra.Command{
		Use: "key",
		Short: "Manage ssh keys",
	}
	if user.Role.HasPermission(dbm.PermKeyAdd)  || user.Role.HasPermission(dbm.PermKeyAddSelf) { cmd.AddCommand(newKeyAddCmd(user).cmd) }
	if user.Role.HasPermission(dbm.PermKeyList) || user.Role.HasPermission(dbm.PermKeyListOwn) { cmd.AddCommand(newKeyListPublicCmd(user).cmd) }
	root.cmd = cmd
	return &root
}

type keyAddCmd struct {
	cmd *cobra.Command
	opts keyAddOpts
}
type keyAddOpts struct {
	username string
}
func newKeyAddCmd(user *dbm.UserWithRole) *keyAddCmd {
	root := keyAddCmd{}
	cmd := &cobra.Command{
		Use: "add PUBLIC_KEY",
		Short: "Add public ssh key to user (default self)",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := dbm.OpenSQLite()
			if err != nil { return err }
			defer dbm.CloseDB(db, &err)

			u := dbm.User{}
			if root.opts.username == "" {
				if user.ID >= 0 {
					return fmt.Errorf("not running as a logged in user")
				}
				u.Name = user.Name
				u.ID = user.ID
			} else {
				s, err := dbm.UserGet(db, root.opts.username)
				if err != nil {
					return fmt.Errorf("getting user %s: %w", root.opts.username, err)
				}
				u = *s
			}

			pubKey := args[0]

			if err := dbm.KeyAdd(db, pubKey, u); err != nil {
				return fmt.Errorf("adding pub key: %w", err)
			}
			return nil
		},
	}
	if user.Role.HasPermission(dbm.PermKeyAdd) {
		cmd.Flags().StringVarP(&root.opts.username, "user", "u", "", "Add key to specific user, otherwise add to yourself")
	}
	root.cmd = cmd
	return &root
}

type keyListPublicCmd struct {
	cmd *cobra.Command
	opts keyListPublicOpts
}
type keyListPublicOpts struct {
	includeComment bool
	includeFingerprint bool
}
func newKeyListPublicCmd(user *dbm.UserWithRole) *keyListPublicCmd {
	root := keyListPublicCmd{}
	cmd := &cobra.Command{
		Use: "ls USERNAME",
		Short: "List all public keys owned by USERNAME seperated by new line",
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := dbm.OpenSQLite()
			if err != nil { return err }
			defer dbm.CloseDB(db, &err)

			u := dbm.User{}

			if len(args) == 1 {
				username := args[0]
				newUser, err := dbm.UserGet(db, username)
				if err != nil {
					return fmt.Errorf("getting user '%s': %w", username, err)
				}
				u = *newUser

			} else {
				if user.ID < 0 {
					return fmt.Errorf("not running as a logged in user: %w", err)
				}
				u.ID = user.ID
				u.Name = user.Name
			}

			keys, err := dbm.KeyGetKeysForUserID(db, u.ID)
			if err != nil  {
				return fmt.Errorf("getting for for user '%s': %w", u.Name, err)
			}
			if len(keys) == 0 {
				ui.Info("User '%s' has no keys\n", u.Name)
				return nil
			}
			for _, k := range keys {
				t, err := ssh.ParsePublicKey(k.PK)
				if err != nil {
					return fmt.Errorf("parsing public key: %w", err)
				}
				key := ssh.MarshalAuthorizedKey(t)

				keyWithoutNewLine, _ := strings.CutSuffix(string(key), "\n")
				p := []string{keyWithoutNewLine}
				if root.opts.includeComment     { p = append(p, k.Comment) }
				if root.opts.includeFingerprint { p = append(p, k.Fingerprint) }

				for _, s := range p {
					ui.Info("%s ", s)
				}
				ui.Info("\n")
			}

			return nil
		},
	}
	cmd.Flags().BoolVarP(&root.opts.includeComment, "comment", "c", false, "Include key comment")
	cmd.Flags().BoolVarP(&root.opts.includeFingerprint, "fingerprint", "f", false, "Include key fingerprint")
	root.cmd = cmd
	return &root
}

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
func newKeyCmd(udb userDB) *keyCmd {
	root := keyCmd{}
	cmd := &cobra.Command{
		Use: "key",
		Short: "Manage ssh keys",
	}
	if udb.user.Role.HasPermission(dbm.PermKeyAdd)  || udb.user.Role.HasPermission(dbm.PermKeyAddSelf) { cmd.AddCommand(newKeyAddCmd(udb).cmd) }
	if udb.user.Role.HasPermission(dbm.PermKeyList) || udb.user.Role.HasPermission(dbm.PermKeyListOwn) { cmd.AddCommand(newKeyListPublicCmd(udb).cmd) }
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
func newKeyAddCmd(udb userDB) *keyAddCmd {
	root := keyAddCmd{}
	cmd := &cobra.Command{
		Use: "add PUBLIC_KEY",
		Short: "Add public ssh key to user (default self)",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			u := dbm.User{}
			if root.opts.username == "" {
				if udb.user.ID >= 0 {
					return fmt.Errorf("not running as a logged in user")
				}
				u.Name = udb.user.Name
				u.ID = udb.user.ID
			} else {
				s, err := dbm.UserGet(udb.db, root.opts.username)
				if err != nil {
					return fmt.Errorf("getting user %s: %w", root.opts.username, err)
				}
				u = *s
			}

			pubKey := args[0]

			if err := dbm.KeyAdd(udb.db, pubKey, u); err != nil {
				return fmt.Errorf("adding pub key: %w", err)
			}
			return nil
		},
	}
	if udb.user.Role.HasPermission(dbm.PermKeyAdd) {
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
func newKeyListPublicCmd(udb userDB) *keyListPublicCmd {
	root := keyListPublicCmd{}
	cmd := &cobra.Command{
		Use: "ls USERNAME",
		Short: "List all public keys owned by USERNAME seperated by new line",
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			u := dbm.User{}

			if len(args) == 1 {
				username := args[0]
				newUser, err := dbm.UserGet(udb.db, username)
				if err != nil {
					return fmt.Errorf("getting user '%s': %w", username, err)
				}
				u = *newUser

			} else {
				if udb.user.ID < 0 {
					return fmt.Errorf("not running as a logged in user, used id %d", udb.user.ID)
				}
				u.ID = udb.user.ID
				u.Name = udb.user.Name
			}

			keys, err := dbm.KeyGetKeysForUserID(udb.db, u.ID)
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

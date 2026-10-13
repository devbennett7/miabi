// Copyright (c) 2026 Jonas Kaninda.
//
// This file is part of Miabi.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jkaninda/okapi/okapicli"
	"github.com/miabi-io/miabi/internal/config"
	"github.com/miabi-io/miabi/internal/models"
	"github.com/miabi-io/miabi/internal/services/adminrecovery"
	"github.com/miabi-io/miabi/internal/services/alerting"
	"github.com/miabi-io/miabi/internal/storage/repositories"
)

const adminUsage = `usage: miabi admin <action> --email <admin email> [--yes]

Recover a platform admin locked out of the console. Run inside the control-plane container:

  docker exec -it miabi miabi admin recover --email admin@example.com

actions:
  recover          new temporary password, two-factor off, account enabled, lockout cleared,
                   all sessions ended
  reset-password   new password and all sessions ended (generated, or read from stdin with
                   --password-stdin)
  reset-2fa        turn off two-factor authentication`

func registerAdminCommands(cli *okapicli.CLI) {
	cli.Command("admin", "Recover a locked-out platform admin: recover | reset-password | reset-2fa (run inside the control-plane container)", runAdmin).
		String("email", "e", "", "The platform admin's email address").
		Bool("password-stdin", "", false, "reset-password: read the new password from stdin instead of generating one").
		Bool("yes", "y", false, "Do not prompt")
	cli.Command("reset-2fa", "Turn off two-factor authentication for a user (run inside the control-plane container)", runResetTwoFactor).
		String("email", "e", "", "The user's email address").
		Bool("yes", "y", false, "Do not prompt")
}

func runAdmin(cmd *okapicli.Command) error {
	args := cmd.Args()
	if len(args) != 1 {
		fmt.Println(adminUsage)
		return errors.New("name one action")
	}
	email := strings.TrimSpace(cmd.GetString("email"))
	if email == "" {
		return errors.New("--email is required")
	}
	svc, done := openRecovery()
	defer done()
	u, err := svc.Admin(email)
	if err != nil {
		return fmt.Errorf("%s: %w", email, err)
	}
	ctx := context.Background()
	yes := cmd.GetBool("yes")

	switch args[0] {
	case "recover":
		if !yes && !confirm(fmt.Sprintf("Recover %s? Their password is replaced, two-factor is turned off and every session they have is ended.", email)) {
			fmt.Println("Aborted.")
			return nil
		}
		res, err := svc.Recover(ctx, u)
		if err != nil {
			return err
		}
		fmt.Printf("Recovered %s.\n", email)
		if res.TwoFactorRemoved {
			fmt.Println("  Two-factor authentication turned off; it must be set up again at the next admin sign-in.")
		}
		if res.Reactivated {
			fmt.Println("  The account was disabled and has been enabled.")
		}
		fmt.Printf("  %d session(s) ended.\n", res.SessionsEnded)
		printPassword(res.Password)
	case "reset-password":
		password := ""
		if cmd.GetBool("password-stdin") {
			if !yes {
				return errors.New("--password-stdin needs --yes: stdin carries the password, so it cannot answer the prompt")
			}
			if password, err = readPassword(); err != nil {
				return err
			}
		}
		if !yes && !confirm(fmt.Sprintf("Reset the password of %s and end every session they have?", email)) {
			fmt.Println("Aborted.")
			return nil
		}
		res, err := svc.ResetPassword(ctx, u, password)
		if err != nil {
			return err
		}
		fmt.Printf("Password reset for %s; %d session(s) ended.\n", email, res.SessionsEnded)
		if res.Password != "" {
			printPassword(res.Password)
		}
	case "reset-2fa":
		return resetTwoFactor(svc, u, yes)
	default:
		fmt.Println(adminUsage)
		return fmt.Errorf("unknown action %q", args[0])
	}
	return nil
}

// runResetTwoFactor is the original top-level command, kept for any user, not only admins.
func runResetTwoFactor(cmd *okapicli.Command) error {
	email := strings.TrimSpace(cmd.GetString("email"))
	if email == "" {
		return errors.New("--email is required")
	}
	svc, done := openRecovery()
	defer done()
	u, err := svc.User(email)
	if err != nil {
		return fmt.Errorf("%s: %w", email, err)
	}
	return resetTwoFactor(svc, u, cmd.GetBool("yes"))
}

func resetTwoFactor(svc *adminrecovery.Service, u *models.User, yes bool) error {
	if !u.TwoFactorEnabled {
		fmt.Printf("%s does not have two-factor authentication enabled; nothing to do.\n", u.Email)
		return nil
	}
	if !yes && !confirm(fmt.Sprintf("Turn off two-factor authentication for %s? They will set it up again at next sign-in to the admin console.", u.Email)) {
		fmt.Println("Aborted.")
		return nil
	}
	if _, err := svc.ResetTwoFactor(u); err != nil {
		return err
	}
	fmt.Printf("Two-factor authentication turned off for %s.\n", u.Email)
	return nil
}

func openRecovery() (*adminrecovery.Service, func()) {
	cfg := config.New()
	_ = cfg.InitWorker()
	cfg.InitStorage()
	db := cfg.Database.DB
	notifier := alerting.NewAdminNotifier(repositories.NewUserRepository(db), repositories.NewNotificationInboxRepository(db), nil)
	return adminrecovery.NewService(db, cfg.Redis.Client, notifier), func() {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
		if cfg.Redis.Client != nil {
			_ = cfg.Redis.Client.Close()
		}
	}
}

func readPassword() (string, error) {
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil && line == "" {
		return "", fmt.Errorf("read password from stdin: %w", err)
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func printPassword(password string) {
	fmt.Printf("\n  Temporary password: %s\n\n", password)
	fmt.Println("  It is shown once and must be changed at the next sign-in.")
}

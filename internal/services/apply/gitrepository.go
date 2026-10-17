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

package apply

import (
	"context"
	"fmt"

	"github.com/miabi-io/miabi/internal/declarative"
	"github.com/miabi-io/miabi/internal/models"
	"github.com/miabi-io/miabi/internal/services/gitrepo"
	"github.com/miabi-io/miabi/internal/services/secret"
)

// SetGitRepositories wires kind: GitRepository.
func (s *Service) SetGitRepositories(g *gitrepo.Service) { s.gitRepos = g }

func (s *Service) snapshotGitRepositories(workspaceID uint, set *declarative.ResourceSet) error {
	if s.gitRepos == nil {
		return nil
	}
	repos, err := s.gitRepos.List(workspaceID)
	if err != nil {
		return fmt.Errorf("snapshot git repositories: %w", err)
	}
	for i := range repos {
		g := &repos[i]
		set.Add(declarative.Resource{
			APIVersion: declarative.APIVersion, Kind: declarative.KindGitRepository,
			Metadata: meta("", g.Name, g.Metadata),
			GitRepository: &declarative.GitRepositorySpec{
				URL: g.URL, Auth: string(g.AuthType), Username: g.Username, SecretRef: g.SecretRef,
			},
		})
	}
	return nil
}

func (s *Service) applyGitRepository(ctx context.Context, workspaceID uint, ch declarative.Change, desired declarative.Resource) error {
	if s.gitRepos == nil {
		return fmt.Errorf("%w: %s", ErrUnsupportedKind, declarative.KindGitRepository)
	}
	if ch.Action == declarative.ActionDelete {
		g, err := s.gitRepos.GetByName(workspaceID, ch.Name)
		if err != nil {
			return fmt.Errorf("gitrepository %q: %w", ch.Name, err)
		}
		return s.gitRepos.Delete(workspaceID, g.ID)
	}
	spec := desired.GitRepository
	if spec == nil {
		return fmt.Errorf("%w: gitrepository %q: spec is required", ErrInvalidManifest, ch.Name)
	}
	in := gitrepo.Input{
		Name: ch.Name, URL: spec.URL, AuthType: models.GitAuthType(spec.Auth), Username: spec.Username,
		Annotations: desired.Metadata.Annotations,
	}
	// Always a vault reference, never a copy: rotating the Secret rotates the credential.
	if spec.SecretRef != "" {
		in.Secret = secret.Ref(spec.SecretRef)
	}
	if ch.Action == declarative.ActionCreate {
		in.Metadata = tagSource(ctx, models.SetBuiltin(models.SanitizeUserMetadata(desired.Metadata.Labels),
			models.MetaManagedBy, managedBy(ctx)))
		_, err := s.gitRepos.Create(workspaceID, in)
		return s.gitRepoErr(ch.Name, err)
	}
	g, err := s.gitRepos.GetByName(workspaceID, ch.Name)
	if err != nil {
		return fmt.Errorf("gitrepository %q: %w", ch.Name, err)
	}
	in.Metadata = claim(ctx, models.MergeUserMetadata(g.Metadata, desired.Metadata.Labels))
	_, err = s.gitRepos.Update(workspaceID, g.ID, in)
	return s.gitRepoErr(ch.Name, err)
}

func (s *Service) gitRepoErr(name string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: gitrepository %q: %v", ErrInvalidManifest, name, err)
}

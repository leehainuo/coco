# Repository maintenance

## Continuous integration

`.github/workflows/ci.yml` runs for pull requests, pushes to `main`, and manual dispatches:

- **Go:** Go 1.21 compatibility and the current stable release; vet, race tests, coverage, and build.
- **Frontend:** Node.js 22, `npm ci`, Prettier, TypeScript, production build, then Go tests with fresh embedded assets.
- **Examples:** dependencies, vet, race tests, and build using `example/go.mod`.
- **CI passed:** a stable aggregate check for branch protection. Require this check rather than version-specific matrix job names.

Actions are pinned to commit SHAs and receive Dependabot updates. Workflows use a
read-only token and do not retain checkout credentials.

## Dependencies

Dependabot checks GitHub Actions, frontend npm dependencies, and both Go modules
every Monday at 09:00 Asia/Shanghai. Minor and patch library updates are grouped;
major library updates remain separate. Review and merge updates after CI passes.

Frontend dependency updates must also rebuild and commit `internal/assets/dist/`.
CI builds fresh assets for validation but does not push generated files automatically.

## Releases

Coco is a Go library. To release, verify CI on the intended commit, ensure embedded
frontend assets were rebuilt, create a semantic version tag such as `v0.2.0`, and
publish a GitHub Release. Use **Generate release notes** to apply the categories in
`.github/release.yml`. This configuration does not publish a release automatically.

## Repository settings

The default branch is `main`. Enable branch update suggestions and automatic
deletion of merged branches. Protect `main` with pull requests, the `CI passed`
status check, resolved review conversations, and no force pushes or deletions.
The existing `Protect main` ruleset also requires one approval, code-owner review,
linear history, and squash merging; repository administrators have an existing bypass.

GitHub settings are stored on GitHub; these files do not apply settings by themselves.
`CODEOWNERS` requests review from `@leehainuo`. The existing code-owner review rule
means maintainer-authored PRs need the administrator's explicit bypass to merge.

## 本地检查

```sh
# Repository root
go vet ./...
go test -race -count=1 ./...
go build ./...

# Frontend; the build writes internal/assets/dist/
cd frontend
npm ci
npm run format:check
npm run type-check
npm run build
cd ..
go test ./...

# Independent example module
cd example
go mod download
go mod verify
go vet ./...
go test -race -count=1 ./...
go build ./...
```

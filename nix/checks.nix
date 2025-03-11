{ inputs, ... }:
{
  perSystem =
    {
      self',
      lib,
      pkgs,
      system,
      ...
    }:
    {
      checks =
        let
          commonHooks = {
            nixfmt-rfc-style.enable = true;
            deadnix = {
              enable = true;
              settings.edit = true;
            };
            statix.enable = true;
            gitleaks = {
              enable = true;
              name = "gitleaks";
              entry = "${lib.getExe pkgs.gitleaks} protect --verbose --redact --staged";
              pass_filenames = false;
            };
            gofmt.enable = true;
          };
        in
        {
          pre-commit-local = inputs.nix-pre-commit-hooks.lib.${system}.run {
            src = ../.;

            hooks = commonHooks // {
              golangci-lint = {
                enable = true;
                entry =
                  let
                    pkg = pkgs.writeShellApplication {
                      name = "run-golangci-lint";
                      runtimeInputs = [
                        pkgs.go
                        pkgs.golangci-lint
                      ];
                      text = "golangci-lint run";
                    };
                  in
                  lib.getExe pkg;
                pass_filenames = false;
              };
            };
          };

          pre-commit-ci = inputs.nix-pre-commit-hooks.lib.${system}.run {
            src = ../.;

            hooks = commonHooks;
          };

          golangci-lint = self'.packages.terraform-provider-zone.overrideAttrs (old: {
            name = "golangci-lint";
            nativeBuildInputs = old.nativeBuildInputs ++ [ pkgs.golangci-lint ];
            buildPhase = ''
              HOME=$TMPDIR golangci-lint run
            '';
          });

          ci = pkgs.symlinkJoin {
            name = "ci";
            paths = [
              self'.checks.golangci-lint
              self'.checks.pre-commit-ci
            ];
          };
        };
    };
}

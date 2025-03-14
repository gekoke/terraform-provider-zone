{ self, inputs, ... }:
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
      checks = {
        pre-commit = inputs.nix-pre-commit-hooks.lib.${system}.run {
          src = self.outPath;
          hooks = {
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
        };

        golangci-lint = self'.packages.terraform-provider-zone.overrideAttrs (old: {
          name = "golangci-lint";
          nativeBuildInputs = old.nativeBuildInputs ++ [ pkgs.golangci-lint ];
          buildPhase = ''
            HOME=$TMPDIR golangci-lint run
          '';
          doInstallCheck = false;
        });
      };
    };
}

{ self, ... }:
{
  perSystem =
    {
      lib,
      pkgs,
      system,
      ...
    }:
    {
      checks.golangci-lint = self.packages.${system}.default.overrideAttrs (old: {
        name = "golangci-lint";
        nativeBuildInputs = old.nativeBuildInputs ++ [ pkgs.golangci-lint ];
        buildPhase = ''
          HOME=$TMPDIR golangci-lint run
        '';
      });

      pre-commit = {
        check.enable = true;

        settings = {
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
      };
    };
}

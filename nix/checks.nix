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
      checks = {
        pre-commit = inputs.nix-pre-commit-hooks.lib.${system}.run {
          src = ../.;
          hooks = import ./git-hooks.nix { inherit lib pkgs; };
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

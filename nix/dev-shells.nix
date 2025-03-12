_: {
  perSystem =
    {
      self',
      pkgs,
      ...
    }:
    {
      devShells = {
        default = pkgs.mkShellNoCC {
          name = "terraform-provider-zone-shell";

          packages = [
            pkgs.golangci-lint
            pkgs.opentofu
            pkgs.terraform
          ];

          shellHook = ''
            ${self'.checks.pre-commit-local.shellHook}
          '';
        };
      };
    };
}

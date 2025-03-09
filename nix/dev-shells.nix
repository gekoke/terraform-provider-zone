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
            pkgs.terraform
            pkgs.opentofu
          ];

          shellHook = ''
            ${self'.checks.pre-commit-local.shellHook}
          '';
        };
      };
    };
}

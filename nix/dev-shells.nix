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
          packages = [
            pkgs.golangci-lint
            pkgs.opentofu
            pkgs.terraform
          ];

          inherit (self'.checks.pre-commit) shellHook;
        };
      };
    };
}

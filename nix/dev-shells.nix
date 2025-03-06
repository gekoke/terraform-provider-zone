{ inputs, ... }:
{
  perSystem =
    {
      config,
      lib,
      pkgs,
      system,
      ...
    }:
    {
      _module.args.pkgs = import inputs.nixpkgs {
        inherit system;
        config.allowUnfreePredicate = pkg: builtins.elem (lib.getName pkg) [ "terraform" ];
      };

      devShells = {
        default = pkgs.mkShellNoCC {
          name = "terraform-provider-zone-shell";

          packages = [
            pkgs.terraform
            pkgs.opentofu
          ];

          shellHook = ''
            ${config.pre-commit.installationScript}
          '';
        };
      };
    };
}

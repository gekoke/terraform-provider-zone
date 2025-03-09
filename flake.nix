{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    systems.url = "github:nix-systems/default";

    flake-parts.url = "github:hercules-ci/flake-parts";

    nix-pre-commit-hooks.url = "github:cachix/pre-commit-hooks.nix";
  };

  outputs =
    inputs:
    inputs.flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;

      imports =
        let
          initializePkgs = {
            perSystem =
              { lib, system, ... }:
              {
                _module.args.pkgs = import inputs.nixpkgs {
                  inherit system;
                  config.allowUnfreePredicate = pkg: builtins.elem (lib.getName pkg) [ "terraform" ];
                };
              };
          };
        in
        [
          initializePkgs
          ./nix/checks.nix
          ./nix/dev-shells.nix
          ./nix/packages.nix
        ];
    };
}

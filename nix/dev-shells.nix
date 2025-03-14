{ self, inputs, ... }:
{
  perSystem =
    {
      lib,
      pkgs,
      system,
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

          shellHook =
            let
              pre-commit-local = inputs.nix-pre-commit-hooks.lib.${system}.run {
                src = self;

                hooks =
                  let
                    common = import ./git-hooks.nix { inherit lib pkgs; };
                    golangci-lint-local = {
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
                  in
                  common // { inherit golangci-lint-local; };
              };
            in
            ''
              ${pre-commit-local.shellHook}
            '';
        };
      };
    };
}

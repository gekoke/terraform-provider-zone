{ self, ... }:
{
  perSystem =
    { lib, pkgs, ... }:
    let
      pkg = pkgs.buildGoModule rec {
        pname = "terraform-provider-zone";
        version = "0.0.0";
        src = lib.sourceFilesBySuffices self [
          ".go"
          ".mod"
          ".sum"
        ];
        vendorHash = "sha256-7jRPVTP8F4RH4KxvHWeBbEgSrHnPOo9lxfg290QqMzA=";
        patchPhase =
          let
            setCorrectVersion = "substituteInPlace ./main.go --subst-var-by providerVersion '${version}'";
          in
          ''
            runHook prePatch
            ${setCorrectVersion}
            runHook postPatch
          '';
        nativeInstallCheckInputs = [ pkgs.versionCheckHook ];
        doInstallCheck = true;
      };
      crossCompile =
        os: arch:
        pkg.overrideAttrs (
          old:
          old
          // {
            env = {
              GOOS = os;
              GOARCH = arch;
            };
            ldflags = lib.concatStringsSep " " [
              "-s -w" # Strip debug info
              "-extldflags=-static" # Link C libraries statically
            ];
            doInstallCheck = false;
          }
        );
    in
    {
      packages = {
        # Nix-native (depends on Nix store and appropriate `pkgs`)
        default = pkg;
        terraform-provider-zone = pkg;

        # Statically linked, cross-compiled using the Go toolchain. Works without Nix on target platform.
        terraform-provider-zone-darwin-amd64 = crossCompile "darwin" "amd64";
        terraform-provider-zone-darwin-arm64 = crossCompile "darwin" "arm64";

        terraform-provider-zone-freebsd-386 = crossCompile "freebsd" "386";
        terraform-provider-zone-freebsd-amd64 = crossCompile "freebsd" "amd64";

        terraform-provider-zone-linux-386 = crossCompile "linux" "386";
        terraform-provider-zone-linux-amd64 = crossCompile "linux" "amd64";
        terraform-provider-zone-linux-arm = crossCompile "linux" "arm";
        terraform-provider-zone-linux-arm64 = crossCompile "linux" "arm64";

        terraform-provider-zone-windows-386 = crossCompile "windows" "386";
        terraform-provider-zone-windows-amd64 = crossCompile "windows" "amd64";
      };
    };
}

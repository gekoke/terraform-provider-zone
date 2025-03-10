_: {
  perSystem =
    { pkgs, ... }:
    {
      packages = rec {
        default = terraform-provider-zone;

        terraform-provider-zone = pkgs.buildGoModule {
          name = "terraform-provider-zone";
          src = ../.;
          vendorHash = "sha256-7jRPVTP8F4RH4KxvHWeBbEgSrHnPOo9lxfg290QqMzA=";
        };
      };
    };
}

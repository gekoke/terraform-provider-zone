_: {
  perSystem =
    { pkgs, ... }:
    {
      packages = rec {
        default = terraform-provider-zone;

        terraform-provider-zone = pkgs.buildGoModule {
          name = "terraform-provider-zone";
          src = ../.;
          vendorHash = "sha256-1TbTOEsgI6fJI2Af5+K9F8Raxl3wxhsTzYN8+oiDXfg=";
        };
      };
    };
}

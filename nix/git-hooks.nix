{ pkgs, lib, ... }:
{
  nixfmt-rfc-style.enable = true;
  deadnix = {
    enable = true;
    settings.edit = true;
  };
  statix.enable = true;
  gitleaks = {
    enable = true;
    name = "gitleaks";
    entry = "${lib.getExe pkgs.gitleaks} protect --verbose --redact --staged";
    pass_filenames = false;
  };
  gofmt.enable = true;
}

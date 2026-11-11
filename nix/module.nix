# Home Manager module for multistop.
#
# Usage, in a home-manager config:
#
#   imports = [ /path/to/multistop/nix/module.nix ];
#   programs.multistop.enable = true;
#
{ config, lib, pkgs, ... }:
let
  cfg = config.programs.multistop;
in
{
  options.programs.multistop = {
    enable = lib.mkEnableOption "multistop, a simple named multi-stopwatch TUI";

    package = lib.mkOption {
      type = lib.types.package;
      default = pkgs.callPackage ./package.nix { };
      description = "The multistop package to install.";
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = [ cfg.package ];
  };
}

# Builds the `multistop` TUI binary from the Go module one directory up.
{
  lib,
  buildGoModule,
}:

buildGoModule {
  pname = "multistop";
  version = "0.1.0";

  src = ../.;

  vendorHash = "sha256-TUbaUoqDZoQTkcOMtoE/FlAiqkWN+x49JeGkDguh2UU=";

  # bubbles/textinput pulls in a clipboard dep that wants cgo; not needed here.
  env.CGO_ENABLED = "0";

  meta = {
    description = "Simple multi-stopwatch TUI: named timers, start/pause/reset";
    mainProgram = "multistop";
    platforms = lib.platforms.linux;
  };
}

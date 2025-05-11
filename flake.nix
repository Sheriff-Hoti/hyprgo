{
  description = "TUI Wallpaper selector";

  inputs = {
    # Include the nixpkgs input to use the Nix package collection
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.11"; # You can change this to a specific channel or commit
  };

  outputs =
    { self, nixpkgs, ... }:
    let
      system = "x86_64-linux"; # Change to your system architecture (e.g., aarch64-linux for ARM)
    in
    {
      # Define the package to be built
      packages.${system} =
        let
          # Import the nixpkgs repository
          pkgs = import nixpkgs { inherit system; };
        in
        pkgs.buildGoModule rec {
          pname = "hyprgo";
          version = "0.1.0";

          src = fetchFromGitHub {
            owner = "Sheriff-Hoti";
            repo = "hyprgo";
            rev = "v${version}";
            hash = "sha256-0dxf0wj1fi74bzys68w2bnlclks5njq6zbp58crijkj8c22f9hca";
          };

          meta = {
            description = "Useful clipboard manager TUI for Unix";
            homepage = "https://github.com/savedra1/clipse";
            license = lib.licenses.mit;
            mainProgram = "hyprgo";
            # maintainers = [ lib.maintainers.savedra1 ];
          };
        };
    };
}

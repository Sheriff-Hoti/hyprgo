{
  description = "TUI Wallpaper Selector";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-24.11";
  };

  outputs =
    { nixpkgs, ... }:
    let
      # you can also put any architecture you want to support here
      # i.e. aarch64-darwin for never M1/2 macbooks
      system = "x86_64-linux";
      pname = "hyprgo";
      version = "0.1.0";
    in
    {
      packages.${system} =
        let
          pkgs = import nixpkgs { inherit system; }; # this gives us access to nixpkgs as we are used to
        in
        {
          default = pkgs.buildGoModule {
            name = pname;
            src = pkgs.fetchFromGitHub {
              owner = "Sheriff-Hoti";
              repo = pname;
              rev = "v${version}";
              hash = "sha256-${pkgs.lib.fakeHash}";
            };

            vendorHash = "sha256-${pkgs.lib.fakeHash}";
          };
        };
    };
}

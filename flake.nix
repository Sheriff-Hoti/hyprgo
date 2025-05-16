{
  description = "TUI Wallpaper Selector";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.11";
  };

  outputs =
    {
      self,
      nixpkgs,
    }:
    let
      #System types to support.
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
      ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

      version = "0.1.0";
      pname = "hyprgo";
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.buildGoModule {
            inherit pname;
            inherit version;

            src = ./.;

            vendorHash = "sha256-L/D4+cTkofF+RTLRt7KytRm/rC2BxmuUz2hh/IPRvzE=";
          };
        }
      );

      defaultPackage = forAllSystems (system: self.packages.${system}.default);

      devShell = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        with pkgs;
        mkShell {
          buildInputs = [
            go
            gopls
          ];
        }
      );

      nixosModule = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          config,
          lib,
          pkgs,
          ...
        }:
        with lib;
        {
          options.programs.hyprgo = {
            enable = mkEnableOption "hyprgo enable";

            package = mkOption {
              type = types.package;
              default = self.packages.${system}.default;
              description = "hyprgo package to use";
            };

          };

        }
      );

      homeModule = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          config,
          lib,
          pkgs,
          ...
        }:
        with lib;
        {
          options.programs.hyprgo = {
            enable = mkEnableOption "hyprgo enable";

            package = mkOption {
              type = types.package;
              default = self.packages.${system}.default;
              description = "hyprgo package to use";
            };

            settings = mkOption {
              type = jsonFormat.type;
              default = null;
              description = "JSON configuration settings for hyprgo";
              example = {
                wallpapers_dir = "~/Pictures/Wallpapers";
                backend = "swayb";
              };
            };
          };

          config = mkIf config.programs.hyprgo.enable (
            let
              cfg = config.programs.hyprgo;
            in
            {
              home.packages = [ cfg.package ];

              xdg.configFile = lib.mkIf (cfg.settings != null) {
                "hyprgo/config.json".source = jsonFormat.generate "hyprgo-config.json" cfg.settings;
              };
            }
          );

        }
      );
    };
}

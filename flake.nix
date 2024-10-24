{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, flake-utils, treefmt-nix }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };

        treefmtEval = treefmt-nix.lib.evalModule pkgs {
          projectRootFile = "flake.nix";
          programs.nixpkgs-fmt.enable = true;
          programs.prettier = {
            enable = true;
            includes = [ "*.md" "*.yaml" "*.yml" ];
          };
          programs.gofmt.enable = true;
        };

        pixelknecht = pkgs.buildGoModule {
          name = "pixelknecht";
          version = "0.0.1";
          #vendorHash = pkgs.lib.fakeHash;
          vendorHash = null;
          src = ./pixelknecht;
        };
        commanderer = pkgs.buildGoModule {
          name = "commanderer";
          version = "0.0.1";
          preBuild = ''
            ${pkgs.templ}/bin/templ generate
          '';
          #vendorHash = pkgs.lib.fakeHash;
          vendorHash = "sha256-LC2frbWQzXasmGg1OKNcRZEcc47YoK7jvKhQrRQgJVA=";
          src = ./commanderer;
        };
      in
      {
        formatter = treefmtEval.config.build.wrapper;
        checks.formatter = treefmtEval.config.build.check self;

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            templ
            pixelnuke
          ];
          shellHook = ''
            alias cover-report="go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html && xdg-open cover.html"
          '';
        };

        packages = flake-utils.lib.flattenTree {
          default = pixelknecht;
          commanderer = commanderer;
        };
      }
    );
}

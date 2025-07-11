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
            includes = [ "*.md" "*.yaml" "*.yml" "*.html" "*.css" "*.js" ];
          };
          programs.gofmt.enable = true;
        };

        #goModVendorHash = pkgs.lib.fakeHash;
        goModVendorHash = "sha256-1IJdLmCq2ENKM+UdyiKVg/IRu0x5zb+oQmrfgwODa/w=";

        pixelknecht = pkgs.buildGoModule {
          name = "pixelknecht";
          version = "0.0.1";
          vendorHash = goModVendorHash;
          src = ./.;
          subPackages = [ "cmd/pixelknecht" ];
        };
        commanderer = pkgs.buildGoModule {
          name = "commanderer";
          version = "0.0.1";
          vendorHash = goModVendorHash;
          src = ./.;
          subPackages = [ "cmd/commanderer" ];
        };
        commandererContainerImage = pkgs.dockerTools.buildLayeredImage {
          name = "ghcr.io/rubenhoenle/pixelknecht/commanderer";
          tag = "unstable";
          config = {
            Expose = [ 9000 ];
            Env = [
              "COMMANDERER_LISTEN_HOST=0.0.0.0"
            ];
            Entrypoint = [ "${commanderer}/bin/commanderer" ];
          };
        };
      in
      {
        formatter = treefmtEval.config.build.wrapper;
        checks.formatter = treefmtEval.config.build.check self;

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            pixelnuke
            golangci-lint
            gnumake
            trivy
          ];
          shellHook = ''
            export COMMANDERER_URL="http://localhost:9000"
            alias cover-report="go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html && xdg-open cover.html"
          '';
        };

        packages = flake-utils.lib.flattenTree {
          default = pixelknecht;
          commanderer = commanderer;
          commandererContainerImage = commandererContainerImage;
        };
      }
    );
}

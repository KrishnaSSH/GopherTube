{
  description = "Description for the project";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [

      ];
      systems = [ "x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin" ];
      perSystem = { config, self', inputs', pkgs, system, ... }: {
        packages.default = pkgs.buildGoModule {
          pname = "gophertube";
          version = "1.2.1";
          src = pkgs.lib.cleanSource ./.;
          vendorHash = "sha256-TBjru54oV2iwAjvqhQcsVr/yQfp6fgukNX/OkZdBWjw=";

          buildInputs = with pkgs; [
            mpv
            fzf
            chafa
          ];

          nativeBuildInputs = with pkgs; [
            go
          ];
          nativeCheckInputs = with pkgs; [
            go
          ];

          buildPhase = ''
            go build -o gophertube main.go
          '';

          installPhase = ''
            mkdir -p $out/bin
            cp gophertube $out/bin
            mkdir -p $out/share/man/man1
            cp $src/man/gophertube.1 $out/share/man/man1
            mkdir -p $out/config
            cp $src/config/gophertube.yaml $out/config/gophertube.yaml.example
          '';
        };
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            go
          ];
          buildInputs = with pkgs; [
            mpv
            fzf
            chafa
          ];
        };
      };
      flake = {

      };
    };
}

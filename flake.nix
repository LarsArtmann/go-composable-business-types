{
  description = "Composable business types for Go — branded IDs, domain primitives, validation";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            golangci-lint
            gofumpt
            golines
            gci
          ];

          env = {
            GOEXPERIMENT = "jsonv2";
            GOPRIVATE = "github.com/LarsArtmann/*,github.com/larsartmann/*";
            GONOSUMCHECK = "github.com/LarsArtmann/*,github.com/larsartmann/*";
            GONOSUMDB = "github.com/LarsArtmann/*,github.com/larsartmann/*";
            GOWORK = "off";
          };
        };

        checks = {
          build = pkgs.runCommand "build-check" {
            nativeBuildInputs = [ pkgs.go ];
            src = ./.;
          } ''
            export HOME=$(mktemp -d)
            export GOMODCACHE=$(mktemp -d)
            export GOPRIVATE="github.com/LarsArtmann/*,github.com/larsartmann/*"
            export GONOSUMCHECK="github.com/LarsArtmann/*,github.com/larsartmann/*"
            export GONOSUMDB="github.com/LarsArtmann/*,github.com/larsartmann/*"
            export GOEXPERIMENT=jsonv2
            cp -r $src/* .
            chmod -R u+w .
            go mod download
            go build ./...
            touch $out
          '';

          test = pkgs.runCommand "test-check" {
            nativeBuildInputs = [ pkgs.go ];
            src = ./.;
          } ''
            export HOME=$(mktemp -d)
            export GOMODCACHE=$(mktemp -d)
            export GOPRIVATE="github.com/LarsArtmann/*,github.com/larsartmann/*"
            export GONOSUMCHECK="github.com/LarsArtmann/*,github.com/larsartmann/*"
            export GONOSUMDB="github.com/LarsArtmann/*,github.com/larsartmann/*"
            export GOEXPERIMENT=jsonv2
            cp -r $src/* .
            chmod -R u+w .
            go mod download
            go test ./... 2>&1 | tee $out
          '';
        };

        formatter = pkgs.nixfmt-classic;
      });
}

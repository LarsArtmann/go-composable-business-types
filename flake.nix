{
  description = "Composable business types for Go — branded IDs, domain primitives, validation";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        goEnv = [
          "GOEXPERIMENT=jsonv2"
          "GOPRIVATE=github.com/LarsArtmann/*,github.com/larsartmann/*"
          "GONOSUMCHECK=github.com/LarsArtmann/*,github.com/larsartmann/*"
          "GONOSUMDB=github.com/LarsArtmann/*,github.com/larsartmann/*"
        ];
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go_1_26
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
          };
        };

        checks = {
          build =
            pkgs.runCommand "build-check"
              {
                nativeBuildInputs = [ pkgs.go_1_26 ];
                src = ./.;
              }
              ''
                export HOME=$(mktemp -d)
                export GOMODCACHE=$(mktemp -d)
                ${builtins.concatStringsSep "\n" (map (e: "export ${e}") goEnv)}
                cp -r $src/. .
                chmod -R u+w .

                # Download deps for each module
                go mod download
                cd nanoid && go mod download && cd ..
                cd locale && go mod download && cd ..
                cd money && go mod download && cd ..
                cd datapoint && go mod download && cd ..
                cd examples && go mod download && cd ..

                # Build all modules via workspace
                go build ./...
                touch $out
              '';

          test =
            pkgs.runCommand "test-check"
              {
                nativeBuildInputs = [ pkgs.go_1_26 ];
                src = ./.;
              }
              ''
                export HOME=$(mktemp -d)
                export GOMODCACHE=$(mktemp -d)
                ${builtins.concatStringsSep "\n" (map (e: "export ${e}") goEnv)}
                cp -r $src/. .
                chmod -R u+w .

                # Download deps for each module
                go mod download
                cd nanoid && go mod download && cd ..
                cd locale && go mod download && cd ..
                cd money && go mod download && cd ..
                cd datapoint && go mod download && cd ..
                cd examples && go mod download && cd ..

                # Test all modules via workspace
                go test -race ./... 2>&1 | tee $out
              '';
        };

        formatter = pkgs.nixfmt;
      }
    );
}

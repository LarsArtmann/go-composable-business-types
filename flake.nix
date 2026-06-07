{
  description = "Composable business types for Go — branded IDs, domain primitives, validation";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
    systems.url = "github:nix-systems/default";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{
      self,
      nixpkgs,
      flake-parts,
      systems,
      treefmt-nix,
    }:
    let
      goEnvVars = {
        GOEXPERIMENT = "jsonv2";
        GOPRIVATE = "github.com/LarsArtmann/*,github.com/larsartmann/*";
        GONOSUMCHECK = "github.com/LarsArtmann/*,github.com/larsartmann/*";
        GONOSUMDB = "github.com/LarsArtmann/*,github.com/larsartmann/*";
      };
    in
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import systems;

      imports = [
        treefmt-nix.flakeModule
      ];

      perSystem =
        {
          config,
          pkgs,
          ...
        }:
        let
          goSrc = pkgs.lib.fileset.toSource {
            root = ./.;
            fileset = pkgs.lib.fileset.unions [
              ./go.mod
              ./go.sum
              ./go.work
              ./actor
              ./bounded
              ./enums
              ./importance
              ./pkg
              ./projectcore
              ./scanutil
              ./tag
              ./temporal
              ./testutil
              ./types
              ./validate
              ./version
              ./nanoid
              ./locale
              ./money
              ./datapoint
              ./examples
            ];
          };

          mkGoCheck =
            {
              name,
              command,
            }:
            pkgs.runCommand name
              {
                nativeBuildInputs = [ pkgs.go_1_26 ];
                src = goSrc;
              }
              ''
                export HOME=$(mktemp -d)
                export GOMODCACHE=$(mktemp -d)
                ${pkgs.lib.concatLines (pkgs.lib.mapAttrsToList (k: v: "export ${k}=${v}") goEnvVars)}
                cp -r $src/. .
                chmod -R u+w .

                go mod download
                cd nanoid && go mod download && cd ..
                cd locale && go mod download && cd ..
                cd money && go mod download && cd ..
                cd datapoint && go mod download && cd ..
                cd examples && go mod download && cd ..

                ${command}
              '';
        in
        {
          treefmt = {
            projectRootFile = "go.mod";
            programs = {
              gofumpt.enable = true;
              nixfmt.enable = true;
            };
          };

          devShells.default = pkgs.mkShell {
            packages = builtins.attrValues {
              inherit (pkgs)
                go_1_26
                golangci-lint
                gofumpt
                golines
                gci
                ;
            };

            env = goEnvVars;
          };

          checks = {
            build = mkGoCheck {
              name = "build-check";
              command = ''
                go build ./... || exit 1
                (cd nanoid && go build ./...) || exit 1
                (cd locale && go build ./...) || exit 1
                (cd money && go build ./...) || exit 1
                (cd datapoint && go build ./...) || exit 1
                (cd examples && go build ./...) || exit 1
                touch $out
              '';
            };

            test = mkGoCheck {
              name = "test-check";
              command = ''
                TEST_OUT=$(mktemp)
                RC=0
                go test -race ./... >> "$TEST_OUT" 2>&1 || RC=1
                (cd nanoid && go test -race ./...) >> "$TEST_OUT" 2>&1 || RC=1
                (cd locale && go test -race ./...) >> "$TEST_OUT" 2>&1 || RC=1
                (cd money && go test -race ./...) >> "$TEST_OUT" 2>&1 || RC=1
                (cd datapoint && go test -race ./...) >> "$TEST_OUT" 2>&1 || RC=1
                cat "$TEST_OUT" > $out
                rm "$TEST_OUT"
                exit $RC
              '';
            };

            lint = mkGoCheck {
              name = "lint-check";
              command =
                let
                  golangci-lint = pkgs.lib.getExe pkgs.golangci-lint;
                in
                ''
                  ${golangci-lint} run ./... || exit 1
                  (cd nanoid && ${golangci-lint} run ./...) || exit 1
                  (cd locale && ${golangci-lint} run ./...) || exit 1
                  (cd money && ${golangci-lint} run ./...) || exit 1
                  (cd datapoint && ${golangci-lint} run ./...) || exit 1
                  (cd examples && ${golangci-lint} run ./...) || exit 1
                  touch $out
                '';
            };
          };
        };
    };
}

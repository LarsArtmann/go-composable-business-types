{
  description = "Composable business types for Go — branded IDs, domain primitives, validation";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
  };

  outputs =
    {
      self,
      nixpkgs,
      systems,
    }:
    let
      eachSystem = f: nixpkgs.lib.genAttrs (import systems) (system: f nixpkgs.legacyPackages.${system});

      goEnvVars = {
        GOEXPERIMENT = "jsonv2";
        GOPRIVATE = "github.com/LarsArtmann/*,github.com/larsartmann/*";
        GONOSUMCHECK = "github.com/LarsArtmann/*,github.com/larsartmann/*";
        GONOSUMDB = "github.com/LarsArtmann/*,github.com/larsartmann/*";
      };

      goSrc = pkgs: pkgs.lib.fileset.toSource {
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
        { pkgs, name, command }:
        pkgs.runCommand name
          {
            nativeBuildInputs = [ pkgs.go_1_26 ];
            src = goSrc pkgs;
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
      devShells = eachSystem (pkgs: {
        default = pkgs.mkShell {
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
      });

      checks = eachSystem (pkgs: {
        build = mkGoCheck {
          inherit pkgs;
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
          inherit pkgs;
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
          inherit pkgs;
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

        format = pkgs.runCommand "format-check"
          {
            nativeBuildInputs = [ pkgs.nixfmt ];
            src = pkgs.lib.fileset.toSource {
              root = ./.;
              fileset = pkgs.lib.fileset.intersection (pkgs.lib.fileset.gitTracked ./.) (pkgs.lib.fileset.fileFilter (file: file.hasExt "nix") ./.);
            };
          }
          ''
            nixfmt --check $src && touch $out
          '';
      });

      formatter = eachSystem (pkgs: pkgs.nixfmt);
    };
}

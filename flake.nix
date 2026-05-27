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

          env = goEnvVars // {
            GOWORK = "off";
          };
        };
      });

      checks = eachSystem (pkgs: {
        build = mkGoCheck {
          inherit pkgs;
          name = "build-check";
          command = "go build ./... && touch $out";
        };

        test = mkGoCheck {
          inherit pkgs;
          name = "test-check";
          command = "go test -race ./... 2>&1 | tee $out";
        };

        lint = mkGoCheck {
          inherit pkgs;
          name = "lint-check";
          command =
            let
              golangci-lint = pkgs.lib.getExe pkgs.golangci-lint;
            in
            "${golangci-lint} run ./... && touch $out";
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

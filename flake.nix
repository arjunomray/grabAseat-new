{
  description = "GrabMySeat — event ticketing monorepo (Go + React)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # Go backend (the service in the repo root).
        backend = pkgs.buildGoModule {
          pname = "grabmyseat-backend";
          version = "0.1.0";
          src = ./.;
          vendorHash = "sha256-2+4p/L5hpiRUsZwIRY4VPTJ+5/kcOk3iSJnH14To1Ug=";
          ldflags = [
            "-s"
            "-w"
          ];
          meta.mainProgram = "eventer";
        };

        # React frontend (Vite static build).
        frontend = pkgs.buildNpmPackage {
          pname = "grabmyseat-frontend";
          version = "0.1.0";
          src = ./frontend;
          nodejs = pkgs.nodejs_22;
          npmDepsHash = "sha256-i7RWWbkJniotQlBx55mbsKDf83QA43HWd3bKVZYSnSs=";
          installPhase = ''
            mkdir -p $out
            cp -r dist/. $out/
          '';
        };

        # `nix run` — start Postgres, the API, and the UI from the repo root.
        devApp = pkgs.writeShellApplication {
          name = "grabmyseat-dev";
          runtimeInputs = with pkgs; [
            bash
            coreutils
            findutils
            gnugrep
            gawk
            gnumake
            go
            nodejs_22
            postgresql_16
            goose
            jq
            curl
            cacert
          ];
          text = ''
            export SSL_CERT_FILE="${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
            export NIX_SSL_CERT_FILE="$SSL_CERT_FILE"

            if [ ! -f Makefile ]; then
              echo "grabmyseat-dev: run this from the repository root (no Makefile found)." >&2
              exit 1
            fi

            exec make dev
          '';
        };
      in
      {
        # `nix develop`
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            gotools
            gofumpt
            golangci-lint
            govulncheck
            nodejs_22
            postgresql_16
            air
            gnumake
            goose
            jq
            curl
            cacert
            nixfmt
            pre-commit
          ];

          shellHook = ''
            export SSL_CERT_FILE="${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
            export NIX_SSL_CERT_FILE="$SSL_CERT_FILE"
            echo "GrabMySeat dev shell ($(go version))"
          '';
        };

        # `nix build .#backend` / `nix build .#frontend`
        packages = {
          inherit backend frontend;
        };

        # `nix flake check` builds both.
        checks = {
          inherit backend frontend;
        };

        # `nix run` — start Postgres, the API, and the UI together.
        apps.default = {
          type = "app";
          program = "${devApp}/bin/grabmyseat-dev";
        };

        # `nix fmt` — format every .nix file in the repo.
        formatter = pkgs.writeShellScriptBin "nixfmt" ''
          exec ${pkgs.nixfmt}/bin/nixfmt "$@" $(find . -type f -name '*.nix' -not -path './.git/*')
        '';
      }
    );
}

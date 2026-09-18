{
  description = "GrabMySeat — event ticketing monorepo (Go + React)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";
  };

  outputs =
    { nixpkgs, ... }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];

      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f {
        inherit system;
        pkgs = import nixpkgs { inherit system; };
      });
    in
    {
      devShells = forAllSystems ({ pkgs, ... }: {
        default = pkgs.mkShell {
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

            echo "GrabMySeat dev shell"
            echo "  go:    $(go version)"
            echo "  node:  $(node --version)"
            echo "  npm:   $(npm --version)"
            echo "  psql:  $(psql --version)"
            echo "  make:  $(make --version | head -1)"
          '';
        };
      });
    };
}

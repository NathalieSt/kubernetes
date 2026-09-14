# Vault's web UI sign-in, as OpenTofu

**Generated** by `kubernetes-generator`
(`src/generators/infrastructure/vault-oidc-config`). Do not edit.

What it manages:

- the **`oidc` auth method**, signing in to the web UI through Authelia;
- an **`admin` role** on it, admitting only members of the Authelia
  `admins` group, with a 1h token renewable to 8h;
- the **`admin` policy** that role grants — every capability on every
  path, sudo included. Root in all but name, and not the root token: it expires
  and the audit log names the person.

## Why this is not in `../vault`

Same provider, same bucket, same login — a different state
(`vault-oidc/terraform.tfstate`) and different Jobs. This run reads a secret a human writes,
at plan time, and a missing one used to fail the workloads run with it:
every workload's policy and placeholder, and the nightly drift check. Apart, a
missing secret fails `vault-apply-oidc-*` and nothing else.

## Setting it up, once

Every step is by hand because each is either a secret or a thing OpenTofu
cannot do to itself.

```sh
# 1. The opentofu policy in ../vault/README.md, with its six OIDC lines. Both
#    runs log in as the same role.

# 2. The client secret: plaintext for Vault, digest for Authelia. Different
#    values — writing the same string to both fails every login with
#    invalid_client. This prints both halves together.
podman run --rm docker.io/authelia/authelia:4.39.25 authelia crypto hash generate pbkdf2 --variant sha512 --random --random.charset numeric-hex --random.length 64
vault kv put kvv2/vault/oidc client-secret=<the random password>
vault kv patch kvv2/authelia/oidc-clients vault='<the digest>'
kubectl -n authelia rollout restart deploy/authelia

# 3. An oidc mount enabled by hand and never configured, if there is one.
#    The provider cannot import a mount with no config, and creating over it
#    fails. Unconfigured means nobody can have signed in through it.
vault write auth/oidc/oidc/auth_url role=x redirect_uri=http://x 2>&1 | grep -q 'could not load configuration' && vault auth disable oidc

# 4. Apply, then open the UI: it lands on the OIDC tab with the admin role.
kubectl create job -n vault-apply --from=cronjob/vault-apply-oidc-apply now
kubectl logs -n vault-apply -l app.kubernetes.io/name=vault-apply-oidc-apply -f
```

Your Authelia user must be in the `admins` group (`authelia/users`), and
must have a TOTP device: the client requires two factors.

## When it fails

    Error: Unable to Read Resource from Vault
      with ephemeral.vault_kv_secret_v2.vault_oidc_client
    Vault response was nil

is a 404: `kvv2/vault/oidc` does not exist. It is read by an
ephemeral resource — never into state — so the plan cannot start without it.
A 403 instead means the `opentofu` policy lacks the OIDC lines.

## Rotating the secret

Steps 2 and 4, plus one thing: bump `oidcClientSecretVersion` in the
generator. The field is write-only, so a plan cannot see that the value in KV
changed, and it is sent again only when that number does.

## Not the CLI

`vault login -method=oidc` wants a `http://localhost:8250/oidc/callback`
redirect, which neither the Authelia client nor the role allows. The CLI keeps
signing in with a token.

## Running it by hand

As `../vault/README.md` describes, in this directory. The provider file here
has no `auth_login` block and takes `VAULT_TOKEN` from the environment.

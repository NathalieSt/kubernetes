# Vault, as OpenTofu

**Generated.** Every file here comes from `kubernetes-generator`
(`src/generators/infrastructure/vault-config`), which collects `meta.vault`
from every generator that declares one. Do not edit them — change the
declaration next to the workload instead.

## What is managed

The frame around the secrets, never the secrets themselves:

- a **policy** per workload, granting read on exactly the paths it declared;
- a **Kubernetes auth role** per workload, bound to that workload's
  ServiceAccount and namespace.

The KV paths themselves are not OpenTofu resources — see Placeholders below.

| Role | Namespace | Paths |
| --- | --- | --- |
| `backup-dashboard` | `backup-dashboard` | `kvv2/backup-dashboard/basic-auth`, `kvv2/backup-dashboard/oidc` |
| `code-server` | `code-server` | `kvv2/code-server`, `kvv2/code-server-git` |
| `neko` | `neko` | `kvv2/neko` |
| `kavita` | `kavita` | `kvv2/kavita/oidc` |
| `mail` | `mail` | `kvv2/mail`, `kvv2/mail/oidc` |
| `gluetun-media` | `gluetun-media` | `kvv2/vpn` |
| `music-library-builder` | `jellyfin` | `kvv2/music-library-builder` |
| `qbittorrent` | `qbittorrent` | `kvv2/vpn` |
| `psono-server` | `psono` | `kvv2/psono` |
| `limes` | `limes` | `kvv2/limes` |
| `tandoor` | `tandoor` | `kvv2/tandoor/oidc` |
| `zyme` | `zyme` | `kvv2/zyme` |
| `radicale` | `radicale` | `kvv2/radicale/users` |
| `anki-sync` | `anki-sync` | `kvv2/anki-sync/users` |
| `authelia` | `authelia` | `kvv2/authelia/secrets`, `kvv2/authelia/users`, `kvv2/authelia/oidc-clients` |
| `gluetun-proxy` | `gluetun-proxy` | `kvv2/vpn` |
| `flux-notifications` | `flux-system` | `kvv2/flux-notifications` |
| `glitchtip` | `glitchtip` | `kvv2/glitchtip/secret`, `kvv2/glitchtip/oidc` |
| `healthchecks` | `healthchecks` | `kvv2/healthchecks/secret`, `kvv2/healthchecks/oidc` |
| `heartbeats` | `healthchecks` | `kvv2/heartbeats` |
| `ntfy` | `ntfy` | `kvv2/ntfy/auth` |
| `victoria-metrics` | `victoria-metrics` | `kvv2/victoria-metrics`, `kvv2/victoria-metrics/ntfy`, `kvv2/victoria-metrics/healthchecks` |

The web UI's sign-in — the `oidc` auth method, its role and the
`admin` policy — is a separate run with its own state, in
`../vault-oidc`. It reads a hand-written secret at plan time, and keeping it
out of here means a missing one fails only that run.

## What is not managed, and why

- **The `kvv2` mount and the `kubernetes` auth backend.** Both
  already exist. Adopt them with `tofu import` if you want them here.
- **The KV paths and their values.** See below.
- **OpenTofu's own Vault credentials.** A policy cannot create the policy that
  authorises creating it; see the bootstrap below.

## Placeholders

`placeholders.txt` lists every declared path. Before each apply, the
vault-apply Job's `placeholders` step (`placeholders.sh` in the vault-apply
generator) goes through it with one rule: **a path that exists is never
written.** A path that does not exist is created holding one line per key,

    REPLACE ME — vault kv patch kvv2/<path> <key>=<value>

with check-and-set 0, so even a human writing the same path at the same moment
cannot be overwritten. The value is deliberately unusable: a workload started
against it fails rather than running on a secret sitting in a git repository.

So the order of things no longer matters. Write a secret by hand before its
declaration merges, or after the placeholder appears — both end with the value
you wrote, and neither fails an apply. (Until September 2026 these were
`vault_kv_secret_v2` resources, and a path written first was a 403 on every
run; the `removed` blocks in `workloads.tf.json` are what took them out of
state without touching Vault.)

The step lists what is still waiting for a human — a key holding a placeholder,
or a declared key that is absent — at the end of its log. That is a to-do, not
a failure: the Job stays green.

```sh
kubectl logs -n vault-apply job/<the apply job> -c placeholders
```

## Bootstrap, once, by hand

OpenTofu needs somewhere to keep state, and a Vault login of its own. A policy
cannot create the policy that authorises creating it, so this part is manual.

```sh
# 1. State bucket and key in Garage, through its admin API — the same calls the
#    garage init job makes. In another terminal, first:
#      kubectl -n garage port-forward svc/garage 3903:3903
#
#    No line continuations anywhere below: every command is one line, so this
#    survives being copied out of a rendered file, a terminal or an editor.

ADMIN=$(kubectl -n garage get secret garage-node-secrets -o jsonpath='{.data.admin_token}' | base64 -d)
API=http://127.0.0.1:3903
g() { curl -sS -H "Authorization: Bearer $ADMIN" -H 'Content-Type: application/json' "$@"; }

g -X POST "$API/v2/CreateBucket" -d '{"globalAlias":"opentofu-state"}'
BUCKET=$(g "$API/v2/GetBucketInfo?globalAlias=opentofu-state" | jq -r '.id // .bucket.id')

g -X POST "$API/v2/CreateKey" -d '{"name":"opentofu"}' > /tmp/opentofu-key.json
ACCESS=$(jq -r .accessKeyId /tmp/opentofu-key.json)
SECRET=$(jq -r .secretAccessKey /tmp/opentofu-key.json)

# Reading the bucket id back rather than trusting the create: CreateBucket fails
# if it already exists, and an empty id would silently grant nothing.
test -n "$BUCKET" -a "$BUCKET" != null || echo "ERROR: no bucket id for opentofu-state"

g -X POST "$API/v2/AllowBucketKey" -d "$(jq -nc --arg b "$BUCKET" --arg k "$ACCESS" '{bucketId:$b,accessKeyId:$k,permissions:{read:true,write:true,owner:false}}')"

echo "access-key-id:     $ACCESS"
echo "secret-access-key: $SECRET"
rm -f /tmp/opentofu-key.json

# 2. A policy for what this configuration manages, and nothing else
vault policy write opentofu - <<'EOF'
path "sys/policy/*"                { capabilities = ["create", "read", "update", "delete", "list"] }
path "sys/policies/acl/*"          { capabilities = ["create", "read", "update", "delete", "list"] }
path "auth/kubernetes/role/*"  { capabilities = ["create", "read", "update", "delete", "list"] }
path "kvv2/data/*"           { capabilities = ["create", "read"] }
path "kvv2/metadata/*"       { capabilities = ["create", "read", "update", "list"] }
path "sys/auth"                    { capabilities = ["read"] }
path "sys/auth/oidc"               { capabilities = ["create", "read", "update", "delete", "sudo"] }
path "sys/mounts/auth/oidc"        { capabilities = ["read", "update"] }
path "sys/mounts/auth/oidc/tune"   { capabilities = ["read", "update"] }
path "auth/oidc/config"            { capabilities = ["create", "read", "update"] }
path "auth/oidc/role/*"            { capabilities = ["create", "read", "update", "delete", "list"] }
EOF
```

One policy for both runs — this one and `../vault-oidc` — because both log in
as the same `opentofu` role. The last six lines are the web UI's
sign-in, and only that run uses them. Measured, not guessed: that is what an
apply and a clean re-plan needed against a dev server, run with a token holding
only this policy.

Note what `sys/policy/*` already meant before those lines: this policy can
write any policy, `admin` included. The OIDC lines add a way to
hand that policy to a person; they do not add a privilege the apply did not
have.

Note what that policy allows on KV: `create` and `read`, no `update` and
no `delete`. Even a compromised apply cannot rewrite a filled-in secret.

### How it logs in depends on where it runs

**In a pod** — Kubernetes auth, no stored credential at all. This is how the
`vault-apply` CronJobs run, and it is the default.

```sh
vault write auth/kubernetes/role/opentofu \
  bound_service_account_names=opentofu \
  bound_service_account_namespaces=vault-apply \
  token_policies=opentofu \
  ttl=15m
```

No `audience` on this one, unlike the workload roles above. Those are used by
VSO, which requests a token with the `vault` audience; this pod presents
its own default ServiceAccount token, and setting an audience the token does not
carry rejects every login.

The bootstrap secret it needs is separate, and manual for the same reason the
policy is — OpenTofu cannot create what it needs in order to run:

```sh
vault kv put kvv2/opentofu \
  access-key-id=<garage key id> \
  secret-access-key=<garage secret> \
  encryption='{"key_provider":{"pbkdf2":{"state":{"passphrase":"<long passphrase>"}}},"method":{"aes_gcm":{"state":{"keys":"${key_provider.pbkdf2.state}"}}},"state":{"method":"${method.aes_gcm.state}"},"plan":{"method":"${method.aes_gcm.state}"}}'
```

VSO needs a policy and role to read that path, and those are manual too:

```sh
vault policy write vault-apply - <<'EOF'
path "kvv2/data/opentofu"     { capabilities = ["read"] }
path "kvv2/metadata/opentofu" { capabilities = ["read"] }
EOF

vault write auth/kubernetes/role/vault-apply \
  bound_service_account_names=vault-apply-vault-serviceaccount \
  bound_service_account_namespaces=vault-apply \
  audience=vault \
  token_policies=vault-apply \
  ttl=15m
```

**From CI or a laptop** — Kubernetes auth is *not* available. Forgejo runs its
jobs in containers under a dind sidecar, so a job has no ServiceAccount token to
present, and Vault's NetworkPolicy does not admit the runner in any case. Use
AppRole, and keep the two halves in Forgejo's secrets:

```sh
vault auth enable approle    # once, if it is not already on
vault write auth/approle/role/opentofu \
  token_policies=opentofu token_ttl=15m token_max_ttl=30m
vault read  auth/approle/role/opentofu/role-id
vault write -f auth/approle/role/opentofu/secret-id
```

Reaching Vault from CI also needs the runner added to Vault's NetworkPolicy —
which widens Vault's ingress to everything that can run a CI job. That is the
trade the in-pod path avoids.

## How it runs

In the `vault-apply` namespace, from a ConfigMap holding exactly these files:

| | When | Does |
| --- | --- | --- |
| `vault-apply-apply-<hash>` | once, when the merged configuration changes | placeholders, then `tofu apply` |
| `vault-apply-plan` | nightly, 05:00 | placeholders in check mode, then `tofu plan -detailed-exitcode` |
| `vault-apply-apply` | never — suspended | the same apply, on demand |

`vault-apply-oidc-*` are the same for `../vault-oidc`, from their own
ConfigMap, without the placeholders step. Separate Jobs so that each failure
names its own run.

**Merging is the deliberate act.** The apply Job's name carries a hash of
everything it runs, so a merged manifests PR that changes a Vault declaration
creates a new Job, Flux runs it once, and prunes the previous one. A PR that does
not touch Vault leaves the Job — and its log — as it was.

That is what makes a red Job mean something:

- **A failed apply** makes the `vault-apply` Kustomization unhealthy, which Flux
  alerting reports. The declarations did not reach Vault.
- **A failed nightly plan** is drift: nothing is waiting to be applied any more,
  so a difference means someone changed Vault by hand, or a declared path was
  deleted.
- **A secret nobody has filled in** fails neither.

Re-running an apply — after fixing a transient failure, say:

```sh
kubectl create job -n vault-apply --from=cronjob/vault-apply-apply rerun
kubectl logs -n vault-apply job/rerun --all-containers -f
```

## Running it by hand

The copy in this directory has no `auth_login` block, so it authenticates from
your environment instead. You need mesh access for Vault and Garage.

```sh
export VAULT_ADDR=... VAULT_TOKEN=...
export AWS_ACCESS_KEY_ID=... AWS_SECRET_ACCESS_KEY=...   # the Garage key
export TF_ENCRYPTION="$(vault kv get -field=encryption kvv2/opentofu)"

tofu init
tofu plan
```

Both copies write the same state, so a hand-run plan and the nightly one see the
same world.

`TF_ENCRYPTION` is not optional in spirit even though OpenTofu will run
without it. State holds policy documents and role bindings — a map of who can
read what — and it lives in a bucket alongside the backups. Keep the passphrase
where the Vault unseal keys are, not only in Vault.

# Vault, as OpenTofu

**Generated.** Every file here comes from `kubernetes-generator`
(`src/generators/infrastructure/vault-config`), which collects `meta.vault`
from every generator that declares one. Do not edit them — change the
declaration next to the workload instead.

## What is managed

The frame around the secrets, never the secrets themselves:

- a **policy** per workload, granting read on exactly the paths it declared;
- a **Kubernetes auth role** per workload, bound to that workload's
  ServiceAccount and namespace;
- a **KV-v2 entry** per path, created with placeholder values.

| Role | Namespace | Paths |
| --- | --- | --- |
| `mail` | `mail` | `kvv2/mail` |
| `psono-server` | `psono` | `kvv2/psono` |

## What is not managed, and why

- **The `kvv2` mount and the `kubernetes` auth backend.** Both
  already exist. Adopt them with `tofu import` if you want them here.
- **Secret values.** The placeholders are written once, with
  `ignore_changes = [data_json]` and `disable_read = true`, so a real value
  written afterwards is never read into state and never overwritten.
- **OpenTofu's own Vault credentials.** A policy cannot create the policy that
  authorises creating it; see the bootstrap below.

## Placeholders

A new path is created holding one line per key:

    REPLACE ME — vault kv patch kvv2/<path> <key>=<value>

That value is deliberately unusable. A workload started against it fails rather
than running on a secret that is sitting in a git repository.

Every entry is created with `cas = 0`, meaning "write only if nothing is
there". Running this against a path that already holds a real secret **fails
the apply** rather than clobbering it — which is what makes it safe to point at
a live Vault. Adopt such a path instead:

```sh
tofu import 'vault_kv_secret_v2.<name>' kvv2/data/<path>
```

To find what still needs filling in:

```sh
vault kv get kvv2/<path>   # a REPLACE ME value is unfilled
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
EOF
```

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

Two CronJobs in the `vault-apply` namespace, from a ConfigMap holding exactly
these files:

| | When | Does |
| --- | --- | --- |
| `vault-apply-plan` | nightly, 05:00 | `tofu plan -detailed-exitcode`. Exit code 2 means Vault no longer matches the declarations, so **drift is a failed Job** |
| `vault-apply-apply` | never — suspended | `tofu apply`, on demand only |

Applying is deliberate rather than automatic. Nothing rewrites Vault's policies
because a commit landed:

```sh
kubectl create job -n vault-apply --from=cronjob/vault-apply-apply now
kubectl logs -n vault-apply -l app.kubernetes.io/name=vault-apply-apply -f
```

Check what the nightly plan thought:

```sh
kubectl logs -n vault-apply -l app.kubernetes.io/name=vault-apply-plan --tail=100
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

import k8s.api.rbac.v1 as rbac_api
import k8s.api.core.v1 as core_api
import k8s.apimachinery.pkg.apis.meta.v1 as meta

schema VaultRbac:
    serviceAccount: core_api.ServiceAccount
    role: rbac_api.Role
    rolebinding: rbac_api.RoleBinding

create_vault_rbac = lambda application_name: str {
    service_account_name = application_name + "-vault-serviceaccount"

    service_account = core_api.ServiceAccount {
        metadata: meta.ObjectMeta {
            name: service_account_name
        }
    }
    role = rbac_api.Role {
        metadata: meta.ObjectMeta {
            name: "vault-token-reviewer"
        }
        rules: [
            rbac_api.PolicyRule {
                apiGroups: ["authentication.k8s.io"]
                resources: ["tokenreviews"]
                verbs: ["create"]
            }
            rbac_api.PolicyRule {
                apiGroups: ["authorization.k8s.io"]
                resources: ["subjectaccessreviews"]
                verbs: ["create"]
            }
        ]
    }
    rolebinding = rbac_api.RoleBinding {
        metadata: meta.ObjectMeta {
            name: application_name + "-vault-binding"
        }
        subjects: [
            rbac_api.Subject {
                kind: "ServiceAccount"
                name: service_account_name
            }
        ]
        roleRef: rbac_api.RoleRef {
            kind: "Role"
            name: "vault-token-reviewer"
            apiGroup: "rbac.authorization.k8s.io"
        }
    }
    rbac: VaultRbac = {
        serviceAccount: service_account
        role: role
        rolebinding: rolebinding
    }

    rbac
}

package secretsrepo

//nolint:gosec,lll // for queries
const (
	AddSecretQuery                         = `INSERT INTO secrets (user_id, kind, data) VALUES ($1, $2, $3) RETURNING id, updated_at;`
	UpdateSecretQuery                      = `UPDATE secrets SET kind = $1, data = $2 WHERE id = $3 AND user_id = $4;`
	RemoveSecretQuery                      = `DELETE FROM secrets WHERE id = $1 AND user_id = $2;`
	GetSecretsAboveVersionQuery            = `SELECT id, updated_at, kind, data FROM secrets WHERE user_id = $1 AND updated_at > $2;`
	DeleteUnknownSecretsBeforeVersionQuery = `DELETE FROM secrets WHERE user_id = $1 AND id NOT IN ($2) AND updated_at < $3 RETURNING id;`
)

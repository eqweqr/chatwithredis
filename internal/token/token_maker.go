package token

type TokenMaker interface {
	GenerateToken(string, string) (string, *UserPayload, error)
	VerifyToken(token string) (*UserPayload, error)
	GetUsernameFromToken(token string) (string, error)
	GetRoleFromToken(token string) (string, error)
}

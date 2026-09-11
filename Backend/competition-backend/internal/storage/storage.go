package storage

type Storage struct {
	PocketBase *PocketBaseClient
}

func New(
	pocketBaseURL string,
	adminEmail string,
	adminPassword string,
) *Storage {
	return &Storage{
		PocketBase: NewPocketBaseClient(
			pocketBaseURL,
			adminEmail,
			adminPassword,
		),
	}
}

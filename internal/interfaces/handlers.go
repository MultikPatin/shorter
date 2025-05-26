package interfaces

import "net/http"

// HealthHandlers groups handlers responsible for health checks and diagnostics.
type HealthHandlers interface {
	Ping(w http.ResponseWriter, r *http.Request) // Handles health check requests.
}

// LinkHandlers aggregates handlers dealing with link manipulation (creation, retrieval).
type LinkHandlers interface {
	AddLinkInText(w http.ResponseWriter, r *http.Request) // Adds a link embedded in HTML body.
	AddLink(w http.ResponseWriter, r *http.Request)       // Adds a link extracted from the request payload.
	AddLinks(w http.ResponseWriter, r *http.Request)      // Batches addition of multiple links.
	GetLink(w http.ResponseWriter, r *http.Request)       // Retrieves a previously-shortened link.
}

// UsersHandlers collects handlers focused on user-specific actions like fetching/deleting links.
type UsersHandlers interface {
	GetLinks(w http.ResponseWriter, r *http.Request)    // Fetches all links owned by the authenticated user.
	DeleteLinks(w http.ResponseWriter, r *http.Request) // Deletes selected links belonging to the user.
}

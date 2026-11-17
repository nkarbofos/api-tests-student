package apierrors

const (
	UnauthorizedCode                    = "unauthorized"
	InvalidRequestCode                  = "invalid_request"
	NotFoundCode                        = "not_found"
	ForbiddenCode                       = "forbidden"
	ConflictCode                        = "already_exists"
	AuthorizationHeaderRequiredMessage  = "Authorization header required"
	InvalidAuthorizationMessage         = "Invalid or expired token"
	AdvertisementNotFoundMessage        = "Advertisement not found"
	FavoriteNotFoundMessage             = "Favorite entry not found"
	AdvertisementAlreadyFavoriteMessage = "Advertisement is already in favorites"
	OwnAdvertisementFavoriteMessage     = "Cannot add own advertisement to favorites"
)

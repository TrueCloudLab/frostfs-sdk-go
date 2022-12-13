/*
Package bearer provides bearer token definition.

Bearer token is attached to the object service requests, and it overwrites
extended ACL of the container. Mainly it is used to provide access of private
data for specific user. Therefore, it must be signed by owner of the container.

Define bearer token by setting correct lifetime, extended ACL and owner ID of
the user that will attach token to its requests.

	var bearerToken bearer.Token
	bearerToken.SetExpiration(500)
	bearerToken.SetIssuedAt(10)
	bearerToken.SetNotBefore(10)
	bearerToken.SetEACL(eaclTable)
	bearerToken.SetOwner(ownerID)

Bearer token must be signed by owner of the container.

	err := bearerToken.Sign(privateKey)

Provide signed token in JSON or binary format to the request sender. Request
sender can attach this bearer token to the object service requests:

	import sdkClient "github.com/TrueCloudLab/frostfs-sdk-go/client"

	var headParams sdkClient.PrmObjectHead
	headParams.WithBearerToken(bearerToken)
	response, err := client.ObjectHead(ctx, headParams)
*/
package bearer

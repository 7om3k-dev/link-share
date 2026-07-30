package handlers

import (
	"slices"
	"strings"
)

type clientData struct {
	clientId     string
	clientSecret string
	redirectUris []string
}

// Note: In a regular OAuth system, the authorization server issues the clientID and client secret to the client software
// TODO: TO keep it in database
var clientsData []clientData = []clientData{
	// Temp code
	clientData{
		clientId:     "tmp-oauth-client",
		redirectUris: []string{"http://localhost:8000/tmp-oauth-redirect"},
		clientSecret: "oauth-client-secret-test",
	},
}

func getClientData(clientId string) *clientData {
	i := slices.IndexFunc(clientsData, func(c clientData) bool {
		return strings.Compare(c.clientId, clientId) == 0
	})

	if i == -1 {
		return nil
	}

	return &clientsData[i]
}

// TODO: To consider use https over http
//clientData{
//	clientId: "oauth-client-web-ui",
//	redirectUris: []string{"http://web-ui:8000/oauth-redirect"},
//	clientSecret: "oauth-client-secret-web-ui",
//},
//clientData{
//	clientId: "oauth-client-user-links",
//	redirectUris: []string{"http://user-links:5001/oauth-redirect"},
//	clientSecret: "oauth-client-secret-user-links",
//},

package routines

import "osuStatsBackend/osustatsDB"

type botUser struct {
	UserID           int    `bson:"UserID"`
	UserName         string `bson:"UserName"`
	AuthCodeExpire   string `bson:"AuthCodeExpire"`
	AUTHAccessToken  string `bson:"AUTH_access_token"`
	AUTHRefreshToken string `bson:"AUTH_refresh_token"`
	Disabled         string `bson:"disabled"`
	LastUpdated      string `bson:"lastUpdated,omitempty"`
}
type botDbTop struct {
	ID             int                 `bson:"id"`
	Top            osustatsDB.BotDbTop `bson:"topplays"`
	IsPPCalculated bool                `bson:"isppcalculated"`
	LastUpdated    string              `bson:"lastUpdated"`
}

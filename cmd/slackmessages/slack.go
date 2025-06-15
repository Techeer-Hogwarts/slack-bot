package slackmessages

func CheckUserIsAllowed(userIDs []string, userID string) bool {
	// if userID == "U02AES3BH17" || userID == "U08EWM4AJJE" || userID == "U033UTX061X" {
	// 	return true
	// }
	for _, id := range userIDs {
		if id == userID {
			return true
		}
	}
	return false
}

package xurls

// Schemes is a sorted list of some well-known url schemes that are followed
// by ":" instead of "://". Since these are more prone to false positives, we
// limit their matching.
var Schemes = []string{
	`bitcoin`, // Bitcoin
	`file`,    // Files
	`magnet`,  // Torrent magnets
	`mailto`,  // Mail
	`sms`,     // SMS
	`tel`,     // Telephone
	`xmpp`,    // XMPP
}

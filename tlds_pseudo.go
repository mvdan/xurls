package xurls

// PseudoTLDs is a sorted list of some widely used unofficial TLDs
//
// Sources:
//  * https://en.wikipedia.org/wiki/Pseudo-top-level_domain
//  * https://en.wikipedia.org/wiki/Category:Pseudo-top-level_domains
//  * https://tools.ietf.org/html/draft-grothoff-iesg-special-use-p2p-names-00
var PseudoTLDs = []string{
	`bit`,   // Namecoin
	`exit`,  // Tor exit node
	`gnu`,   // GNS by public key
	`i2p`,   // I2P network
	`local`, // Local network
	`onion`, // Tor hidden services
	`zkey`,  // GNS domain name
}

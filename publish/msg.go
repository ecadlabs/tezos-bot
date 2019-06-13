package publish

import (
	"fmt"
	"log"
	"net"

	"github.com/ecadlabs/tezos-bot/models"
)

// GetStatusString composes a status string based on available vanity data
func GetStatusString(ballot *models.Ballot) string {
	templateBasic := `Tezos address %s voted "%s" %son #Tezos proposal "%s"%s`
	templateVanity := `Tezos baker "%s" /%s voted "%s" %son #Tezos proposal "%s"%s`

	var proposalVanityName string
	protocolName, err := LookupTZName(ballot.ProposalHash, "tz.tezz.ie")
	if err != nil {
		proposalVanityName = ballot.ProposalHash
	} else {
		proposalVanityName = protocolName
	}

	templateRolls := ""
	if ballot.Rolls != 0 {
		templateRolls = fmt.Sprintf("with %d rolls ", ballot.Rolls)
	}

	templateQuorum := "and quorum has been reached."
	percentTowardQuorum := ballot.PercentTowardQuorum()
	if percentTowardQuorum > 0 {
		templateQuorum = fmt.Sprintf("with %.2f%% remaining to reach %.2f%% quorum.", percentTowardQuorum, ballot.Quorum)
	}
	templateStatus := fmt.Sprintf("\n\nVote status is %.2f%% yay/%.2f%% nay, %s", ballot.CountingPercentYay(), ballot.CountingPercentNay(), templateQuorum)

	// tz.tezz.ie is an experimental DNS zone to resolve vanity names from tz
	// addresses
	address, err := LookupTZName(ballot.PKH, "tz.tezz.ie")

	if err != nil {
		log.Printf("No address found for %s, err: %s", ballot.PKH, err)
		return fmt.Sprintf(templateBasic, ballot.PKH, ballot.Ballot, templateRolls, proposalVanityName, templateStatus)
	}
	log.Printf("Address %s found for %s, ", address, ballot.PKH)
	return fmt.Sprintf(templateVanity, address, ballot.PKH, ballot.Ballot, templateRolls, proposalVanityName, templateStatus)

}

// GetProposalSummaryString
func GetProposalSummaryString(summary *models.ProposalSummary) string {
	proposalName, err := LookupTZName(summary.ProposalHash, "tz.tezz.ie")

	if err != nil {
		log.Printf("No protocol found for %s, err: %s", summary.ProposalHash, err)
		proposalName = summary.ProposalHash
	} else {
		proposalName = fmt.Sprintf("%s (%s)", proposalName, summary.ProposalHash)
	}

	return fmt.Sprintf("Proposal upvotes: #Tezos proposal %s now has %d votes.", proposalName, summary.SupporterCount)
}

// GetWinningProposalString
func GetWinningProposalString(summary *models.ProposalSummary) string {
	proposalName, err := LookupTZName(summary.ProposalHash, "tz.tezz.ie")

	if err != nil {
		log.Printf("No protocol found for %s, err: %s", summary.ProposalHash, err)
		proposalName = summary.ProposalHash
	} else {
		proposalName = fmt.Sprintf("%s (%s)", proposalName, summary.ProposalHash)
	}

	return fmt.Sprintf("Proposal period complete: proposal %s received the most upvotes (%d) and is advancing to the exploration vote period.", proposalName, summary.SupporterCount)
}

// GetProposalInjectString retrieve the template for proposal injection status
func GetProposalInjectString(proposal *models.Proposal) string {
	address, err := LookupTZName(proposal.PKH, "tz.tezz.ie")

	if err != nil {
		log.Printf("No address found for %s, err: %s", proposal.PKH, err)
	}

	templateAddress := proposal.PKH

	if address != "" {
		templateAddress = fmt.Sprintf("%s /%s", address, proposal.PKH)
	}

	return fmt.Sprintf("New #Tezos proposal injected! %s injected proposal %s in voting period %d.", templateAddress, proposal.ProposalHash, proposal.Period)
}

// GetProtocolString retrieve the template for protocol change status
func GetProtocolString(proto string) string {
	lookupKey := fmt.Sprintf("%s", proto)

	protocolName, err := LookupTZName(lookupKey, "tz.tezz.ie")

	if err != nil {
		log.Printf("No protocol found for %s, err: %s", lookupKey, err)
		protocolName = proto
	}

	return fmt.Sprintf("Protocol %s is now live on mainnet! #Tezos", protocolName)
}

// LookupTZName queries DNS for a txt record corresponding to a TZ address.
func LookupTZName(address, zone string) (string, error) {
	query := fmt.Sprintf("%s.%s", address, zone)
	rrs, err := net.LookupTXT(query)
	if err != nil {
		return "", err
	}
	return string(rrs[0]), nil
}

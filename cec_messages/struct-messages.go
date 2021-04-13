package cec_messages

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/*
This file contains structs and methods for handling Splunk messages encoded in the CEC Protocol
defined here:
https://cd.splunkdev.com/splcore/main/-/blob/develop/src/search/docs/chunked-command-protocol.txt

*/

// these messages are sent from the custom command to splunk
type CommandCecMessage struct {
	ProtocolVersion string
	MetadataLength  int
	DataLength      int
	MetaData        CommandMetaData
	Data            []map[string][]byte
}

// these messages are sent from Splunk to the custom command
type SplunkCecMessage struct {
	ProtocolVersion string
	MetadataLength  int
	DataLength      int
	MetaData        SplunkMetaData
	Data            []map[string][]byte
}

type CommandMetaData struct {
	Type            string
	Generating      bool
	Require_fields  []string
	Maxwait         int
	Streaming_preop string
	Finished        bool
	Error           string
	Inspector       Inspector
}

type Inspector struct {
	Messages [][2]string
}

type SplunkMetaData struct {
	Action                         string
	Finished                       bool
	Preview                        bool
	Streaming_command_will_restart bool
	Searchinfo                     Searchinfo
}

type Searchinfo struct {
	Args           []string
	Raw_args       []string
	Dispatch_dir   string
	Sid            string
	App            string
	Username       string
	Owner          string
	Session_key    string
	Splunkd_uri    string
	Splunk_version string
	Search         string
	Command        string
	Maxresultrows  int
	Earliest_time  int
	Latest_time    int
}

/*
Parse a CEC message from a string
Input: string CEC formatted message
Output: CECMessage
*/
func (cm SplunkCecMessage) Parse(messageString string) (SplunkCecMessage, error) {
	message := SplunkCecMessage{}
	// scan through the message and break on lines
	messageScanner := bufio.NewScanner(strings.NewReader(messageString))
	messageScanner.Split(bufio.ScanLines)

	// parse header
	if !messageScanner.Scan() {
		return message, messageScanner.Err()
	}
	header := strings.Split(messageScanner.Text(), ",")
	if len(header) != 3 {
		return message, errors.New("Message header does not conform to Protocol!")
	}

	// parse data from the message header
	message.ProtocolVersion = header[0]
	mdLen, err := strconv.Atoi(header[1])
	if err != nil {
		return message, errors.New("Message header does not conform to Protocol!")
	}
	message.MetadataLength = mdLen

	dLen, err := strconv.Atoi(header[2])
	if err != nil {
		return message, errors.New("Message header does not conform to Protocol!")
	}
	message.DataLength = dLen

	// parse metadata
	if !messageScanner.Scan() {
		return message, messageScanner.Err()
	}

	metadataBytes := messageScanner.Text()
	if len(metadataBytes) != message.MetadataLength {
		return message, errors.New(fmt.Sprintf("Metadata length (%d) does not match header declaration (%d)!", len(metadataBytes), message.MetadataLength))
	}
	json.Unmarshal([]byte(metadataBytes), &message.MetaData)

	// parse body and convert to list of key value pairs
	// first row will be the headers, which will become the keys to the records
	if message.DataLength > 0 {
		if !messageScanner.Scan() {
			return message, messageScanner.Err()
		}
		reader := csv.NewReader(strings.NewReader(messageScanner.Text()))
		headers, err := reader.Read()
		if err != nil {
			return message, errors.New("Body header does not conform to CSV Protocol!")
		}

		// parse the body lines and create a map with the headers
		for messageScanner.Scan() {
			var lineMap map[string][]byte = make(map[string][]byte)
			bodyLineReader := csv.NewReader(strings.NewReader(messageScanner.Text()))

			bodyLines, err := bodyLineReader.Read()

			for i, header := range headers {
				// TODO: check for multivalue headers
				lineMap[header] = []byte(bodyLines[i])
			}
			if err != nil {
				return message, errors.New("Body header does not conform to CSV Protocol!")
			}
			message.Data = append(message.Data, lineMap)
		}

		if messageScanner.Err() != nil {
			return message, messageScanner.Err()
		}
	}

	return message, nil
}

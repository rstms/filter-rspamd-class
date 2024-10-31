package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/poolpOrg/OpenSMTPD-framework/filter"
)

/*********************************************************************************************

 filter-rspamd-class

 add a header: 'X-Spam-Class: SPAMCLASS'
 based on threshold levels compared with rspamd's X-Spam-Score header

 Default class threshold levels:

 ham:		spam_score < HAM_THRESHOLD
 possible:	HAM_THRESHOLD <= spam_score < POSSIBLE_THRESHOLD
 probable:	POSSIBLE_THRESHOLD <= spam_score < PROBABLE_THRESHOLD
 spam:		spam_score >= PROBABLE_THRESHOLD

 class names and thresholds are configurable per recipient email address
 using a JSON file with the following format:

{
    "username@example.org": [
	{ "name": "ham", "score": 0 },
	{ "name": "possible", "score": 3 },
	{ "name": "probable", "score": 10 },
	{ "name": "spam", "score": 999 }
    ],
    "othername@example.org": [
	{ "name": "not_spam", "score": 0 },
	{ "name": "suspected_spam", "score": 10 },
	{ "name": "is_spam", "score": 999 }
    ]
}

The final threshold value is set to float32-max automatically; 999 is a placeholder

*********************************************************************************************/

const Version = "0.0.1"

const CLASS_CONFIG_FILE = "/etc/mail/filter_rspamd_classes.json"

const HAM_THRESHOLD = 0.0
const POSSIBLE_THRESHOLD = 3.0
const PROBABLE_THRESHOLD = 10.0

/*
type tx struct {
	msgid    string
	mailFrom string
	rcptTo   []string
	message  []string
	action   string
	response string
}

type session struct {
	id string

	rdns     string
	src      string
	heloName string
	userName string
	mtaName  string

	tx tx
}
*/

type SessionData struct {
	rcptTo []string
}

type SpamClass struct {
	Name  string  `json:"name"`
	Score float32 `json:"score"`
}

var SpamClassLevels map[string][]SpamClass

func loadClassLevels(filename string) error {

	if filename == "" {
		return nil
	}

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil
	}
	configBytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed reading %s: %v", filename, err)
	}
	levels := map[string][]SpamClass{}
	err = json.Unmarshal(configBytes, &levels)
	if err != nil {
		return fmt.Errorf("failed parsing %s: %v", filename, err)
	}
	for addr, classes := range levels {
		classes[len(classes)-1].Score = math.MaxFloat32
		SpamClassLevels[addr] = classes
	}
	return nil
}

func initClassLevels(filename string) error {
	SpamClassLevels = make(map[string][]SpamClass)
	err := loadClassLevels(filename)
	SpamClassLevels["default"] = []SpamClass{
		SpamClass{"ham", HAM_THRESHOLD},
		SpamClass{"possible", POSSIBLE_THRESHOLD},
		SpamClass{"probable", PROBABLE_THRESHOLD},
		SpamClass{"spam", math.MaxFloat32},
	}
	return err
}

func getClass(addresses []string, score float32) string {
	levels := SpamClassLevels["default"]
	for _, address := range addresses {
		userLevels, ok := SpamClassLevels[address]
		if ok {
			levels = userLevels
			break
		}
	}
	var result string
	for _, level := range levels {
		result = level.Name
		if score < level.Score {
			break
		}
	}
	return result
}

func txResetCb(timestamp time.Time, session filter.Session, messageId string) {
	fmt.Fprintf(os.Stderr, "%s: %s: tx-reset: %s\n", timestamp, session, messageId)
}

func txBeginCb(timestamp time.Time, session filter.Session, messageId string) {
	fmt.Fprintf(os.Stderr, "%s: %s: tx-begin: %s\n", timestamp, session, messageId)
}

func getSessionData(session filter.Session) (*SessionData, error) {
	data := session.Get()
	sessionData, ok := data.(*SessionData)
	if !ok {
		return nil, errors.New("SessionData conversion failure")
	}
	return sessionData, nil
}

func txRcptCb(timestamp time.Time, session filter.Session, messageId string, result string, to string) {
	sessionData, err := getSessionData(session)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s: tx-rcpt: %v\n", timestamp, session, err)
		return
	}
	sessionData.rcptTo = append(sessionData.rcptTo, to)
	fmt.Fprintf(os.Stderr, "%s: %s: tx-rcpt: %s|%s|%s\n", timestamp, session, messageId, result, to)
}

func parseSpamScore(line string) (float32, error) {
	fields := strings.Split(line, " ")
	if len(fields) < 2 {
		return 0, errors.New("spam score parse failed")
	}
	score, err := strconv.ParseFloat(fields[1], 32)
	if err != nil {
		return 0, err
	}
	return float32(score), nil
}

func filterDataLineCb(timestamp time.Time, session filter.Session, line string) []string {
	output := []string{line}
	if strings.HasPrefix(line, "X-Spam-Score:") {
		sessionData, err := getSessionData(session)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s: filter-data-line error: %v\n", timestamp, session, err)
			return output
		}
		score, err := parseSpamScore(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s: filter-data-line error: %v\n", timestamp, session, err)
			return output
		}
		class := getClass(sessionData.rcptTo, score)
		output = append(output, "X-Spam-Class: "+class)
		fmt.Fprintf(os.Stderr, "%s: %s: filter-data-line: %s\n", timestamp, session, line)
	}
	return output
}

func main() {

	fmt.Fprintf(os.Stderr, "Version %s\n", Version)
	err := initClassLevels(CLASS_CONFIG_FILE)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
	}

	filter.Init()

	filter.SMTP_IN.SessionAllocator(func() filter.SessionData {
		return &SessionData{}
	})

	filter.SMTP_IN.OnTxReset(txResetCb)
	filter.SMTP_IN.OnTxBegin(txBeginCb)
	filter.SMTP_IN.OnTxRcpt(txRcptCb)
	filter.SMTP_IN.DataLineRequest(filterDataLineCb)

	filter.Dispatch()
}

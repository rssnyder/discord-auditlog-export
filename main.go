package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	token     string
	guild     string
	stop      string
	frequency int
	logLevel  int
	logger    = log.New()
)

func init() {
	flag.StringVar(&token, "token", "", "Bot token")
	flag.StringVar(&guild, "guild", "", "Guild ID")
	flag.StringVar(&stop, "stop", "", "Log ID to stop ingesting")
	flag.IntVar(&frequency, "frequency", 1, "Frequency of updates: seconds")
	flag.IntVar(&logLevel, "log", 0, "Logging level: 0=Info; 1=Debug")
	flag.Parse()

	switch logLevel {
	case 1:
		logger.SetLevel(log.DebugLevel)
	default:
		logger.SetLevel(log.InfoLevel)
	}
}

func main() {

	// create discord connection
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Errorf("Creating Discord session: %s\n", err)
		return
	}

	var before string
	var newStop string
	var hardStop bool

	ticker := time.NewTicker(time.Duration(frequency) * time.Second)

	for {
		select {
		case <-ticker.C:

			logger.Debugf("Getting logs for %s\n", guild)
			hardStop = false

			// Get latest audit log
			entries, err := dg.GuildAuditLog(guild, "", "", 0, 100)
			if err != nil {
				logger.Errorf("Unable to get server audit log: %s: %s\n", guild, err)
				continue
			} else if len(entries.AuditLogEntries) == 0 {
				// wait
				continue
			}

			// Check if there are new entries
			if entries.AuditLogEntries[0].ID < stop {
				logger.Debugf("No more new entries")
				continue
			}

			// Cycle through latest entries
			for i, log := range entries.AuditLogEntries {

				// Hard stop if set
				if log.ID == stop {
					logger.Debug("Hit stop log entry inital")
					hardStop = true
					break
				}

				// Struct to JSON
				b, err := json.Marshal(log)
				if err != nil {
					logger.Error(err)
					continue
				}

				// Deal with data
				fmt.Println(string(b))
				if i == 0 {
					newStop = log.ID
				}
				before = log.ID
			}

			// Get the rest of the entries
			for (len(entries.AuditLogEntries) == 100) && (!hardStop) {

				// Get the next 100 entries
				entries, err := dg.GuildAuditLog(guild, "", before, 0, 100)
				if err != nil {
					logger.Errorf("Unable to get server audit log: %s: %s\n", guild, err)
					continue
				}

				// Process next entries
				for _, log := range entries.AuditLogEntries {

					// Hard stop if set
					if log.ID == stop {
						logger.Debug("Hit stop log entry")
						hardStop = true
						break
					}

					// Struct to JSON
					b, err := json.Marshal(log)
					if err != nil {
						logger.Error(err)
						continue
					}

					// Deal with data
					fmt.Println(string(b))
					before = log.ID
				}
			}
			stop = newStop
		}
	}
}

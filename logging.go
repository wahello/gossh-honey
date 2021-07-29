package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type logEntry interface {
	fmt.Stringer
	eventType() string
}

type connectionLog struct {
	ClientVersion string `json:"client_version"`
}

func (entry connectionLog) String() string {
	return fmt.Sprintf("connection with client version %q established", entry.ClientVersion)
}
func (entry connectionLog) eventType() string {
	return "connection"
}

type connectionCloseLog struct {
}

func (entry connectionCloseLog) String() string {
	return "connection closed"
}
func (entry connectionCloseLog) eventType() string {
	return "connection_close"
}

type tcpipForwardLog struct {
	Address string `json:"address"`
}

func (entry tcpipForwardLog) String() string {
	return fmt.Sprintf("TCP/IP forwarding on %v requested", entry.Address)
}
func (entry tcpipForwardLog) eventType() string {
	return "tcpip_forward"
}

type cancelTCPIPForwardLog struct {
	Address string `json:"address"`
}

func (entry cancelTCPIPForwardLog) String() string {
	return fmt.Sprintf("TCP/IP forwarding on %v canceled", entry.Address)
}
func (entry cancelTCPIPForwardLog) eventType() string {
	return "cancel_tcpip_forward"
}

type noMoreSessionsLog struct {
}

func (entry noMoreSessionsLog) String() string {
	return "rejection of further session channels requested"
}
func (entry noMoreSessionsLog) eventType() string {
	return "no_more_sessions"
}

type channelLog struct {
	ChannelID int `json:"channel_id"`
}

type sessionLog struct {
	channelLog
}

// session 请求
func (entry sessionLog) String() string {
	return fmt.Sprintf("[channel %v] session requested", entry.ChannelID)
}
func (entry sessionLog) eventType() string {
	return "session"
}

type sessionCloseLog struct {
	channelLog
}

func (entry sessionCloseLog) String() string {
	return fmt.Sprintf("[channel %v] closed", entry.ChannelID)
}
func (entry sessionCloseLog) eventType() string {
	return "session_close"
}

type sessionInputLog struct {
	channelLog
	Input string `json:"input"`
}

func (entry sessionInputLog) String() string {
	return fmt.Sprintf("[channel %v] input: %q", entry.ChannelID, entry.Input)
}
func (entry sessionInputLog) eventType() string {
	return "session_input"
}

type directTCPIPLog struct {
	channelLog
	From string `json:"from"`
	To   string `json:"to"`
}

func (entry directTCPIPLog) String() string {
	return fmt.Sprintf("[channel %v] direct TCP/IP forwarding from %v to %v requested", entry.ChannelID, entry.From, entry.To)
}
func (entry directTCPIPLog) eventType() string {
	return "direct_tcpip"
}

type directTCPIPCloseLog struct {
	channelLog
}

func (entry directTCPIPCloseLog) String() string {
	return fmt.Sprintf("[channel %v] closed", entry.ChannelID)
}
func (entry directTCPIPCloseLog) eventType() string {
	return "direct_tcpip_close"
}

type directTCPIPInputLog struct {
	channelLog
	Input string `json:"input"`
}

func (entry directTCPIPInputLog) String() string {
	return fmt.Sprintf("[channel %v] input: %q", entry.ChannelID, entry.Input)
}
func (entry directTCPIPInputLog) eventType() string {
	return "direct_tcpip_input"
}

type ptyLog struct {
	channelLog
	Terminal string `json:"terminal"`
	Width    uint32 `json:"width"`
	Height   uint32 `json:"height"`
}

// ptyLog
func (entry ptyLog) String() string {
	return fmt.Sprintf("[channel %v] PTY using terminal %q (size %vx%v) requested", entry.ChannelID, entry.Terminal, entry.Width, entry.Height)
}
func (entry ptyLog) eventType() string {
	return "pty"
}

type shellLog struct {
	channelLog
}

// shell request
func (entry shellLog) String() string {
	return fmt.Sprintf("[channel %v] shell requested", entry.ChannelID)
}
func (entry shellLog) eventType() string {
	return "shell"
}

type execLog struct {
	channelLog
	Command string `json:"command"`
}

func (entry execLog) String() string {
	return fmt.Sprintf("[channel %v] command %q requested", entry.ChannelID, entry.Command)
}
func (entry execLog) eventType() string {
	return "exec"
}

type subsystemLog struct {
	channelLog
	Subsystem string `json:"subsystem"`
}

func (entry subsystemLog) String() string {
	return fmt.Sprintf("[channel %v] subsystem %q requested", entry.ChannelID, entry.Subsystem)
}
func (entry subsystemLog) eventType() string {
	return "subsystem"
}

type x11Log struct {
	channelLog
	Screen uint32 `json:"screen"`
}

func (entry x11Log) String() string {
	return fmt.Sprintf("[channel %v] X11 forwarding on screen %v requested", entry.ChannelID, entry.Screen)
}
func (entry x11Log) eventType() string {
	return "x11"
}

type envLog struct {
	channelLog
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (entry envLog) String() string {
	return fmt.Sprintf("[channel %v] environment variable %q with value %q requested", entry.ChannelID, entry.Name, entry.Value)
}
func (entry envLog) eventType() string {
	return "env"
}

type windowChangeLog struct {
	channelLog
	Width  uint32 `json:"width"`
	Height uint32 `json:"height"`
}

func (entry windowChangeLog) String() string {
	return fmt.Sprintf("[channel %v] window size change to %vx%v requested", entry.ChannelID, entry.Width, entry.Height)
}
func (entry windowChangeLog) eventType() string {
	return "window_change"
}

type debugGlobalRequestLog struct {
	RequestType string `json:"request_type"`
	WantReply   bool   `json:"want_reply"`
	Payload     string `json:"payload"`
}

func (entry debugGlobalRequestLog) String() string {
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to log event: %v", err)
		return ""
	}
	return fmt.Sprintf("DEBUG global request received: %v\n", string(jsonBytes))
}
func (entry debugGlobalRequestLog) eventType() string {
	return "debug_global_request"
}

type debugChannelLog struct {
	channelLog
	ChannelType string `json:"channel_type"`
	ExtraData   string `json:"extra_data"`
}

func (entry debugChannelLog) String() string {
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to log event: %v", err)
		return ""
	}
	return fmt.Sprintf("DEBUG new channel requested: %v\n", string(jsonBytes))
}
func (entry debugChannelLog) eventType() string {
	return "debug_channel"
}

type debugChannelRequestLog struct {
	channelLog
	RequestType string `json:"request_type"`
	WantReply   bool   `json:"want_reply"`
	Payload     string `json:"payload"`
}

func (entry debugChannelRequestLog) String() string {
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to log event: %v", err)
		return ""
	}
	return fmt.Sprintf("DEBUG channel request received: %v\n", string(jsonBytes))
}
func (entry debugChannelRequestLog) eventType() string {
	return "debug_channel_request"
}

func (context connContext) logEvent(entry logEntry) {
	if strings.HasPrefix(entry.eventType(), "debug_") && !context.cfg.Logging.Debug {
		return
	}
	if context.cfg.Logging.JSON {
		var jsonEntry interface{}
		if context.cfg.Logging.Timestamps {
			jsonEntry = struct {
				Time      string   `json:"time"`
				Source    string   `json:"source"`
				EventType string   `json:"event_type"`
				Event     logEntry `json:"event"`
			}{time.Now().Format(time.RFC3339), context.RemoteAddr().String(), entry.eventType(), entry}
		} else {
			jsonEntry = struct {
				Source    string   `json:"source"`
				EventType string   `json:"event_type"`
				Event     logEntry `json:"event"`
			}{context.RemoteAddr().String(), entry.eventType(), entry}
		}
		logBytes, err := json.Marshal(jsonEntry)
		if err != nil {
			log.Printf("Failed to log event: %v", err)
			return
		}
		log.Print(string(logBytes))
	} else {
		log.Printf("[%v] %v", context.RemoteAddr().String(), entry)
	}
}

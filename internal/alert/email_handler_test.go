package alert

import (
	"net"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/ports"
)

func TestNewEmailHandlerMissingHost(t *testing.T) {
	_, err := NewEmailHandler(EmailConfig{
		To: []string{"ops@example.com"},
	})
	if err == nil {
		t.Fatal("expected error for missing smtp host, got nil")
	}
}

func TestNewEmailHandlerMissingRecipients(t *testing.T) {
	_, err := NewEmailHandler(EmailConfig{
		SMTPHost: "localhost",
		SMTPPort: 25,
	})
	if err == nil {
		t.Fatal("expected error for missing recipients, got nil")
	}
}

func TestNewEmailHandlerValid(t *testing.T) {
	h, err := NewEmailHandler(EmailConfig{
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		From:     "portwatch@example.com",
		To:       []string{"ops@example.com"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestEmailHandlerSendConnectsToSMTP(t *testing.T) {
	// Start a minimal TCP listener to act as a fake SMTP server.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)

	done := make(chan struct{})
	go func() {
		defer close(done)
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		// Write minimal SMTP greeting so smtp.SendMail can proceed.
		conn.Write([]byte("220 localhost ESMTP\r\n"))
	}()

	h, err := NewEmailHandler(EmailConfig{
		SMTPHost: "127.0.0.1",
		SMTPPort: addr.Port,
		From:     "portwatch@example.com",
		To:       []string{"ops@example.com"},
	})
	if err != nil {
		t.Fatalf("unexpected error creating handler: %v", err)
	}

	a := NewAlert(ports.Port{Port: 8080, Proto: "tcp"}, "new listener detected", SeverityWarn)
	err = h.Send(a)
	// We expect an error because our fake server doesn't speak full SMTP,
	// but it should be an SMTP protocol error, not a connection error.
	if err == nil {
		t.Log("send succeeded unexpectedly (full SMTP mock not implemented)")
	} else if strings.Contains(err.Error(), "connection refused") {
		t.Errorf("unexpected connection error: %v", err)
	}

	<-done
}
